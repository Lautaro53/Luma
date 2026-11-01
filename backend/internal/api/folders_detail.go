package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"luma/internal/repository"
	"luma/internal/task"
)

type FolderDetailStore interface {
	Rename(
		ctx context.Context,
		id int64,
		name string,
	) (task.Folder, error)

	Delete(
		ctx context.Context,
		id int64,
	) error
}

func RegisterFolderDetailRoutes(
	mux *http.ServeMux,
	store FolderDetailStore,
) {
	mux.HandleFunc("PUT /api/folders/{id}", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		id, err := parseFolderID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, "ID inválido")
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 8*1024)

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var input task.CreateFolderInput

		if err := decoder.Decode(&input); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}

		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			writeError(
				w,
				http.StatusBadRequest,
				"se admite un solo objeto JSON",
			)
			return
		}

		if err := input.NormalizeAndValidate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		folder, err := store.Rename(
			r.Context(), id, input.Name,
		)

		switch {
		case errors.Is(err, sql.ErrNoRows):
			writeError(
				w, http.StatusNotFound,
				"carpeta no encontrada",
			)
			return

		case errors.Is(err, repository.ErrFolderAlreadyExists):
			writeError(
				w, http.StatusConflict,
				"ya existe una carpeta con ese nombre",
			)
			return

		case err != nil:
			writeError(
				w, http.StatusInternalServerError,
				"no se pudo renombrar la carpeta",
			)
			return
		}

		writeJSON(w, http.StatusOK, folder)
	})

	mux.HandleFunc("DELETE /api/folders/{id}", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		id, err := parseFolderID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, "ID inválido")
			return
		}

		err = store.Delete(r.Context(), id)

		switch {
		case errors.Is(err, sql.ErrNoRows):
			writeError(
				w, http.StatusNotFound,
				"carpeta no encontrada",
			)
			return

		case err != nil:
			writeError(
				w, http.StatusInternalServerError,
				"no se pudo eliminar la carpeta",
			)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func parseFolderID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil || id <= 0 {
		return 0, errors.New("ID inválido")
	}

	return id, nil
}
