package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"luma/internal/repository"
	"luma/internal/task"
)

type FolderStore interface {
	Create(
		ctx context.Context,
		input task.CreateFolderInput,
	) (task.Folder, error)

	List(
		ctx context.Context,
	) ([]task.Folder, error)
}

func RegisterFolderRoutes(
	mux *http.ServeMux,
	store FolderStore,
) {
	mux.HandleFunc("GET /api/folders", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		folders, err := store.List(r.Context())

		if err != nil {
			writeError(
				w,
				http.StatusInternalServerError,
				"no se pudieron consultar las carpetas",
			)
			return
		}

		writeJSON(w, http.StatusOK, folders)
	})

	mux.HandleFunc("POST /api/folders", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		r.Body = http.MaxBytesReader(w, r.Body, 8*1024)

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var input task.CreateFolderInput

		if err := decoder.Decode(&input); err != nil {
			writeError(
				w, http.StatusBadRequest,
				"JSON inválido",
			)
			return
		}

		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			writeError(
				w, http.StatusBadRequest,
				"se admite un solo objeto JSON",
			)
			return
		}

		if err := input.NormalizeAndValidate(); err != nil {
			writeError(
				w, http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		folder, err := store.Create(r.Context(), input)

		if errors.Is(err, repository.ErrFolderAlreadyExists) {
			writeError(
				w, http.StatusConflict,
				"ya existe una carpeta con ese nombre",
			)
			return
		}

		if err != nil {
			writeError(
				w, http.StatusInternalServerError,
				"no se pudo crear la carpeta",
			)
			return
		}

		writeJSON(w, http.StatusCreated, folder)
	})
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
