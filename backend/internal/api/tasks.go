package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"luma/internal/repository"
	"luma/internal/task"
)

type TaskStore interface {
	Create(
		ctx context.Context,
		input task.CreateInput,
	) (task.Task, error)

	ListFiltered(
		ctx context.Context,
		filter repository.TaskFilter,
	) ([]task.Task, error)
}

func RegisterTaskRoutes(
	mux *http.ServeMux,
	store TaskStore,
) {
	mux.HandleFunc("GET /api/tasks", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		filter, err := parseTaskFilter(r)
		if err != nil {
			writeError(
				w, http.StatusBadRequest, err.Error(),
			)
			return
		}

		items, err := store.ListFiltered(
			r.Context(), filter,
		)
		if err != nil {
			writeError(
				w,
				http.StatusInternalServerError,
				"no se pudieron consultar las tareas",
			)
			return
		}

		writeJSON(w, http.StatusOK, items)
	})

	mux.HandleFunc("POST /api/tasks", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		r.Body = http.MaxBytesReader(
			w, r.Body, 8*1024,
		)

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var input task.CreateInput

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

		created, err := store.Create(
			r.Context(), input,
		)

		switch {
		case errors.Is(
			err, repository.ErrTaskFolderNotFound,
		):
			writeError(
				w, http.StatusBadRequest,
				"la carpeta indicada no existe",
			)
			return

		case err != nil:
			writeError(
				w,
				http.StatusInternalServerError,
				"no se pudo crear la tarea",
			)
			return
		}

		writeJSON(w, http.StatusCreated, created)
	})
}

func parseTaskFilter(
	r *http.Request,
) (repository.TaskFilter, error) {
	query := r.URL.Query()

	filter := repository.TaskFilter{
		Status:   task.Status(query.Get("status")),
		Priority: task.Priority(query.Get("priority")),
		Search:   query.Get("search"),
		Sort: repository.TaskSort(
			query.Get("sort"),
		),
	}

	if values, ok := query["folder_id"]; ok {
		if len(values) != 1 {
			return filter, fmt.Errorf(
				"folder_id debe indicarse una sola vez",
			)
		}

		id, err := strconv.ParseInt(
			values[0], 10, 64,
		)
		if err != nil {
			return filter, fmt.Errorf(
				"folder_id inválido",
			)
		}

		filter.FolderID = &id
	}

	if values, ok := query["without_folder"]; ok {
		if len(values) != 1 {
			return filter, fmt.Errorf(
				"without_folder debe indicarse una sola vez",
			)
		}

		value, err := strconv.ParseBool(values[0])
		if err != nil {
			return filter, fmt.Errorf(
				"without_folder debe ser true o false",
			)
		}

		filter.WithoutFolder = value
	}

	if err := filter.Validate(); err != nil {
		return filter, err
	}

	return filter, nil
}
