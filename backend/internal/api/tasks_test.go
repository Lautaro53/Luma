package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"luma/internal/repository"
	"luma/internal/task"
)

type fakeTaskStore struct {
	lastFilter repository.TaskFilter
	lastInput  task.CreateInput
	createErr  error
	listErr    error
}

func (f *fakeTaskStore) ListFiltered(
	_ context.Context,
	filter repository.TaskFilter,
) ([]task.Task, error) {
	f.lastFilter = filter

	if f.listErr != nil {
		return nil, f.listErr
	}

	return []task.Task{}, nil
}

func (f *fakeTaskStore) Create(
	_ context.Context,
	input task.CreateInput,
) (task.Task, error) {
	f.lastInput = input

	if f.createErr != nil {
		return task.Task{}, f.createErr
	}

	return task.Task{
		ID:          1,
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Status:      task.StatusPending,
		FolderID:    input.FolderID,
	}, nil
}

func TestTaskListEndpoint(t *testing.T) {
	store := &fakeTaskStore{}
	mux := http.NewServeMux()
	RegisterTaskRoutes(mux, store)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/tasks?status=pendiente&priority=alta"+
			"&folder_id=2&search=OSPF&sort=priority",
		nil,
	)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("HTTP inesperado: %d", rec.Code)
	}

	if strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Fatalf("respuesta inesperada: %s", rec.Body.String())
	}

	filter := store.lastFilter

	if filter.Status != task.StatusPending ||
		filter.Priority != task.PriorityHigh ||
		filter.Search != "OSPF" ||
		filter.Sort != repository.SortPriority ||
		filter.FolderID == nil ||
		*filter.FolderID != 2 {
		t.Fatalf("filtros incorrectos: %+v", filter)
	}
}

func TestTaskListInvalidFilters(t *testing.T) {
	store := &fakeTaskStore{}
	mux := http.NewServeMux()
	RegisterTaskRoutes(mux, store)

	queries := []string{
		"?status=activo",
		"?priority=urgente",
		"?folder_id=abc",
		"?folder_id=0",
		"?without_folder=incorrecto",
		"?folder_id=1&without_folder=true",
		"?sort=random",
	}

	for _, query := range queries {
		t.Run(query, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/api/tasks"+query,
				nil,
			)

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf(
					"esperado 400, obtenido %d",
					rec.Code,
				)
			}
		})
	}
}

func TestTaskCreateEndpoint(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		createErr  error
		wantStatus int
	}{
		{
			name:       "crear tarea sin carpeta",
			body:       `{"title":"Estudiar Go"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name: "crear tarea con carpeta",
			body: `{
				"title":"Configurar VPN",
				"priority":"alta",
				"folder_id":2
			}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "título vacío",
			body:       `{"title":" "}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "prioridad inválida",
			body: `{
				"title":"Estudiar",
				"priority":"urgente"
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "carpeta inexistente",
			body: `{
				"title":"Estudiar",
				"folder_id":999
			}`,
			createErr:  repository.ErrTaskFolderNotFound,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "campo desconocido",
			body: `{
				"title":"Estudiar",
				"unexpected":true
			}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "JSON malformado",
			body:       `{"title":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "dos objetos JSON",
			body:       `{"title":"A"}{"title":"B"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error interno",
			body:       `{"title":"Estudiar"}`,
			createErr:  errors.New("base no disponible"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeTaskStore{
				createErr: tt.createErr,
			}

			mux := http.NewServeMux()
			RegisterTaskRoutes(mux, store)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/tasks",
				strings.NewReader(tt.body),
			)

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"esperado %d, obtenido %d: %s",
					tt.wantStatus,
					rec.Code,
					rec.Body.String(),
				)
			}

			if rec.Code == http.StatusCreated {
				if !strings.Contains(
					rec.Body.String(), `"status":"pendiente"`,
				) {
					t.Fatalf(
						"estado inicial incorrecto: %s",
						rec.Body.String(),
					)
				}
			}
		})
	}
}
