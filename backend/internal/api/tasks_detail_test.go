package api

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"luma/internal/repository"
	"luma/internal/task"
)

type fakeTaskDetailStore struct {
	items map[int64]task.Task
}

func (f *fakeTaskDetailStore) Update(
	_ context.Context,
	id int64,
	input task.CreateInput,
) (task.Task, error) {
	item, exists := f.items[id]
	if !exists {
		return task.Task{}, sql.ErrNoRows
	}

	if input.FolderID != nil && *input.FolderID == 999 {
		return task.Task{}, repository.ErrTaskFolderNotFound
	}

	item.Title = input.Title
	item.Description = input.Description
	item.Priority = input.Priority
	item.FolderID = input.FolderID
	f.items[id] = item

	return item, nil
}

func (f *fakeTaskDetailStore) SetCompleted(
	_ context.Context,
	id int64,
	completed bool,
) (task.Task, error) {
	item, exists := f.items[id]
	if !exists {
		return task.Task{}, sql.ErrNoRows
	}

	if completed {
		item.Status = task.StatusCompleted

		if item.CompletedAt == nil {
			now := time.Now().UTC()
			item.CompletedAt = &now
		}
	} else {
		item.Status = task.StatusPending
		item.CompletedAt = nil
	}

	f.items[id] = item
	return item, nil
}

func (f *fakeTaskDetailStore) Delete(
	_ context.Context,
	id int64,
) error {
	if _, exists := f.items[id]; !exists {
		return sql.ErrNoRows
	}

	delete(f.items, id)
	return nil
}

func TestTaskDetailEndpoints(t *testing.T) {
	store := &fakeTaskDetailStore{
		items: map[int64]task.Task{
			1: {
				ID:       1,
				Title:    "Estudiar Go",
				Priority: task.PriorityMedium,
				Status:   task.StatusPending,
			},
		},
	}

	mux := http.NewServeMux()
	RegisterTaskDetailRoutes(mux, store)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantText   string
	}{
		{
			name:       "editar tarea",
			method:     http.MethodPut,
			path:       "/api/tasks/1",
			body:       `{"title":"Repasar Go","priority":"alta"}`,
			wantStatus: http.StatusOK,
			wantText:   "Repasar Go",
		},
		{
			name:       "completar tarea",
			method:     http.MethodPatch,
			path:       "/api/tasks/1/status",
			body:       `{"status":"completada"}`,
			wantStatus: http.StatusOK,
			wantText:   `"status":"completada"`,
		},
		{
			name:       "reabrir tarea",
			method:     http.MethodPatch,
			path:       "/api/tasks/1/status",
			body:       `{"status":"pendiente"}`,
			wantStatus: http.StatusOK,
			wantText:   `"completed_at":null`,
		},
		{
			name:       "título inválido",
			method:     http.MethodPut,
			path:       "/api/tasks/1",
			body:       `{"title":" "}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "carpeta inexistente",
			method:     http.MethodPut,
			path:       "/api/tasks/1",
			body:       `{"title":"Estudiar","folder_id":999}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "tarea inexistente",
			method:     http.MethodPut,
			path:       "/api/tasks/999",
			body:       `{"title":"Estudiar"}`,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "estado inválido",
			method:     http.MethodPatch,
			path:       "/api/tasks/1/status",
			body:       `{"status":"en_progreso"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "estado ausente",
			method:     http.MethodPatch,
			path:       "/api/tasks/1/status",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "ID inválido",
			method:     http.MethodDelete,
			path:       "/api/tasks/abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "eliminar tarea",
			method:     http.MethodDelete,
			path:       "/api/tasks/1",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "eliminar inexistente",
			method:     http.MethodDelete,
			path:       "/api/tasks/1",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				tt.method,
				tt.path,
				strings.NewReader(tt.body),
			)

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"esperado HTTP %d, obtenido %d: %s",
					tt.wantStatus,
					rec.Code,
					rec.Body.String(),
				)
			}

			if tt.wantText != "" &&
				!strings.Contains(rec.Body.String(), tt.wantText) {
				t.Fatalf(
					"respuesta inesperada: %s",
					rec.Body.String(),
				)
			}

			if tt.wantStatus == http.StatusNoContent &&
				rec.Body.Len() != 0 {
				t.Fatal("204 no debe incluir contenido")
			}
		})
	}
}
