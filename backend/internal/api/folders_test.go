package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"luma/internal/repository"
	"luma/internal/task"
)

type fakeFolderStore struct {
	folders []task.Folder
}

func (f *fakeFolderStore) List(
	_ context.Context,
) ([]task.Folder, error) {
	return f.folders, nil
}

func (f *fakeFolderStore) Create(
	_ context.Context,
	input task.CreateFolderInput,
) (task.Folder, error) {
	for _, folder := range f.folders {
		if strings.EqualFold(folder.Name, input.Name) {
			return task.Folder{},
				repository.ErrFolderAlreadyExists
		}
	}

	folder := task.Folder{
		ID:   int64(len(f.folders) + 1),
		Name: input.Name,
	}

	f.folders = append(f.folders, folder)
	return folder, nil
}

func TestFolderEndpoints(t *testing.T) {
	mux := http.NewServeMux()

	store := &fakeFolderStore{
		folders: make([]task.Folder, 0),
	}

	RegisterFolderRoutes(mux, store)

	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantText   string
	}{
		{
			name:       "listar vacío",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			wantText:   "[]",
		},
		{
			name:       "crear carpeta",
			method:     http.MethodPost,
			body:       `{"name":"Infraestructura"}`,
			wantStatus: http.StatusCreated,
			wantText:   "Infraestructura",
		},
		{
			name:       "listar carpeta creada",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			wantText:   "Infraestructura",
		},
		{
			name:       "nombre vacío",
			method:     http.MethodPost,
			body:       `{"name":"  "}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "carpeta duplicada",
			method:     http.MethodPost,
			body:       `{"name":"infraestructura"}`,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "campo desconocido",
			method:     http.MethodPost,
			body:       `{"name":"Redes","otro":1}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "JSON malformado",
			method:     http.MethodPost,
			body:       `{"name":`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				tt.method,
				"/api/folders",
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
				!strings.Contains(
					rec.Body.String(),
					tt.wantText,
				) {
				t.Fatalf(
					"respuesta inesperada: %s",
					rec.Body.String(),
				)
			}
		})
	}
}
