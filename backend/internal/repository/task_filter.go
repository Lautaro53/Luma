package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"luma/internal/task"
)

type TaskSort string

const (
	SortNewest   TaskSort = "newest"
	SortOldest   TaskSort = "oldest"
	SortPriority TaskSort = "priority"
)

type TaskFilter struct {
	Status        task.Status
	Priority      task.Priority
	FolderID      *int64
	WithoutFolder bool
	Search        string
	Sort          TaskSort
}

func (f TaskFilter) Validate() error {
	switch f.Status {
	case "", task.StatusPending, task.StatusCompleted:
	default:
		return errors.New("estado de filtro inválido")
	}

	switch f.Priority {
	case "", task.PriorityLow, task.PriorityMedium, task.PriorityHigh:
	default:
		return errors.New("prioridad de filtro inválida")
	}

	if f.FolderID != nil {
		if *f.FolderID <= 0 {
			return errors.New("ID de carpeta inválido")
		}

		if f.WithoutFolder {
			return errors.New("filtros de carpeta incompatibles")
		}
	}

	switch f.Sort {
	case "", SortNewest, SortOldest, SortPriority:
	default:
		return errors.New("ordenamiento inválido")
	}

	return nil
}

func (r *TaskRepository) ListFiltered(
	ctx context.Context,
	filter TaskFilter,
) ([]task.Task, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	var folderID any
	if filter.FolderID != nil {
		folderID = *filter.FolderID
	}

	query := `
		SELECT
			id, title, description, priority, status,
			folder_id, created_at, updated_at, completed_at
		FROM tasks
		WHERE (? = '' OR status = ?)
		  AND (? = '' OR priority = ?)
		  AND (? IS NULL OR folder_id = ?)
		  AND (? = 0 OR folder_id IS NULL)
	`

	args := []any{
		filter.Status, filter.Status,
		filter.Priority, filter.Priority,
		folderID, folderID,
		filter.WithoutFolder,
	}

	search := strings.TrimSpace(filter.Search)

	if search != "" {
		query += `
			AND (
				instr(lower(title), lower(?)) > 0
				OR instr(
					lower(COALESCE(description, '')),
					lower(?)
				) > 0
			)
		`

		args = append(args, search, search)
	}

	orderBy := "created_at DESC, id DESC"

	switch filter.Sort {
	case SortOldest:
		orderBy = "created_at ASC, id ASC"

	case SortPriority:
		orderBy = `
			CASE priority
				WHEN 'alta' THEN 1
				WHEN 'media' THEN 2
				ELSE 3
			END ASC,
			created_at DESC,
			id DESC
		`
	}

	query += " ORDER BY " + orderBy

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("consultar tareas: %w", err)
	}
	defer rows.Close()

	items := make([]task.Task, 0)

	for rows.Next() {
		item, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("leer tarea: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("recorrer tareas: %w", err)
	}

	return items, nil
}
