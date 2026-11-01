package api

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"luma/internal/repository"
	"luma/internal/task"
)

type fakeFolderDetailStore struct {
	folders map[int64]task.Folder
}

func (f *fakeFolderDetailStore) Rename(
	_ context.Context,
	id int64,
	name string,
) (task.Folder, error) {
	folder, exists := f.folders[id]
	if !exists {
		return task.Folder{}, sql.ErrNoRows
	}

	for otherID, other := range f.folders {
		if otherID != id && strings.EqualFold(other.Name, name) {
			return task.Folder{},
				repository.ErrFolderAlreadyExists
		}
	}

	folder.Name = name
	f.folders[id] = folder

	return folder, nil
}

func (f *fakeFolderDetailStore) Delete(
	_ context.Context,
	id int64,
) error {
	if _, exists := f.folders[id]; !exists {
		return sql.ErrNoRows
	}

	delete(f.folders, id)
	return nil
}

func TestFolderDetailEndpoints(t *testing.T) {
	store := &fakeFolderDetailStore{
		folders: map[int64]task.Folder{
			1: {ID: 1, Name: "Infraestructura"},
			2: {ID: 2, Name: "Matemáticas"},
		},
	}

	mux := http.NewServeMux()
	RegisterFolderDetailRoutes(mux, store)

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		status int
	}{
		{
			name:   "renombrar",
			method: http.MethodPut,
			path:   "/api/folders/1",
			body:   `{"name":"Redes"}`,
			status: http.StatusOK,
		},
		{
			name:   "nombre duplicado",
			method: http.MethodPut,
			path:   "/api/folders/1",
			body:   `{"name":"matemáticas"}`,
			status: http.StatusConflict,
		},
		{
			name:   "nombre vacío",
			method: http.MethodPut,
			path:   "/api/folders/1",
			body:   `{"name":"  "}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "ID inválido",
			method: http.MethodPut,
			path:   "/api/folders/abc",
			body:   `{"name":"Redes"}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "carpeta inexistente",
			method: http.MethodPut,
			path:   "/api/folders/999",
			body:   `{"name":"Redes"}`,
			status: http.StatusNotFound,
		},
		{
			name:   "eliminar",
			method: http.MethodDelete,
			path:   "/api/folders/1",
			status: http.StatusNoContent,
		},
		{
			name:   "eliminar inexistente",
			method: http.MethodDelete,
			path:   "/api/folders/1",
			status: http.StatusNotFound,
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

			if rec.Code != tt.status {
				t.Fatalf(
					"esperado HTTP %d, obtenido %d: %s",
					tt.status, rec.Code, rec.Body.String(),
				)
			}
		})
	}

	if _, exists := store.folders[1]; exists {
		t.Fatal("la carpeta eliminada sigue existiendo")
	}

	if _, exists := store.folders[2]; !exists {
		t.Fatal("se eliminó una carpeta incorrecta")
	}
}
