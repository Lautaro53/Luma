package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"luma/internal/repository"
	"luma/internal/task"
)

type TaskDetailStore interface {
	Update(
		ctx context.Context,
		id int64,
		input task.CreateInput,
	) (task.Task, error)

	SetCompleted(
		ctx context.Context,
		id int64,
		completed bool,
	) (task.Task, error)

	Delete(
		ctx context.Context,
		id int64,
	) error
}

func RegisterTaskDetailRoutes(
	mux *http.ServeMux,
	store TaskDetailStore,
) {
	mux.HandleFunc("PUT /api/tasks/{id}", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		id, err := parseTaskID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, "ID inválido")
			return
		}

		var input task.CreateInput

		if err := decodeTaskJSON(w, r, &input); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}

		if err := input.NormalizeAndValidate(); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		updated, err := store.Update(
			r.Context(), id, input,
		)

		switch {
		case errors.Is(err, sql.ErrNoRows):
			writeError(
				w, http.StatusNotFound,
				"tarea no encontrada",
			)
			return

		case errors.Is(err, repository.ErrTaskFolderNotFound):
			writeError(
				w, http.StatusBadRequest,
				"la carpeta indicada no existe",
			)
			return

		case err != nil:
			writeError(
				w, http.StatusInternalServerError,
				"no se pudo actualizar la tarea",
			)
			return
		}

		writeJSON(w, http.StatusOK, updated)
	})

	mux.HandleFunc("PATCH /api/tasks/{id}/status", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		id, err := parseTaskID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, "ID inválido")
			return
		}

		var input struct {
			Status task.Status `json:"status"`
		}

		if err := decodeTaskJSON(w, r, &input); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}

		switch input.Status {
		case task.StatusPending, task.StatusCompleted:
		default:
			writeError(
				w, http.StatusBadRequest,
				"estado inválido",
			)
			return
		}

		updated, err := store.SetCompleted(
			r.Context(),
			id,
			input.Status == task.StatusCompleted,
		)

		if errors.Is(err, sql.ErrNoRows) {
			writeError(
				w, http.StatusNotFound,
				"tarea no encontrada",
			)
			return
		}

		if err != nil {
			writeError(
				w, http.StatusInternalServerError,
				"no se pudo cambiar el estado",
			)
			return
		}

		writeJSON(w, http.StatusOK, updated)
	})

	mux.HandleFunc("DELETE /api/tasks/{id}", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		id, err := parseTaskID(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, "ID inválido")
			return
		}

		err = store.Delete(r.Context(), id)

		if errors.Is(err, sql.ErrNoRows) {
			writeError(
				w, http.StatusNotFound,
				"tarea no encontrada",
			)
			return
		}

		if err != nil {
			writeError(
				w, http.StatusInternalServerError,
				"no se pudo eliminar la tarea",
			)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func parseTaskID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil || id <= 0 {
		return 0, errors.New("ID inválido")
	}

	return id, nil
}

func decodeTaskJSON(
	w http.ResponseWriter,
	r *http.Request,
	destination any,
) error {
	r.Body = http.MaxBytesReader(w, r.Body, 8*1024)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("decodificar JSON: %w", err)
	}

	var extra any

	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("contenido JSON adicional o inválido")
	}

	return nil
}
