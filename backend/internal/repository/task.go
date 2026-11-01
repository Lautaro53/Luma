package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"luma/internal/task"
)

var ErrTaskFolderNotFound = errors.New(
	"la carpeta indicada no existe",
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create guarda una tarea pendiente.
func (r *TaskRepository) Create(
	ctx context.Context,
	input task.CreateInput,
) (task.Task, error) {
	if err := input.NormalizeAndValidate(); err != nil {
		return task.Task{}, err
	}

	// nil representa una tarea sin carpeta.
	var folderID any
	if input.FolderID != nil {
		folderID = *input.FolderID
	}

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO tasks (
			title, description, priority, folder_id
		)
		SELECT ?, ?, ?, ?
		WHERE ? IS NULL
		   OR EXISTS (
		       SELECT 1 FROM folders WHERE id = ?
		   )
		RETURNING
			id, title, description, priority, status,
			folder_id, created_at, updated_at, completed_at
	`,
		input.Title,
		input.Description,
		input.Priority,
		folderID,
		folderID,
		folderID,
	)

	created, err := scanTask(row)

	if errors.Is(err, sql.ErrNoRows) &&
		input.FolderID != nil {
		return task.Task{}, ErrTaskFolderNotFound
	}

	if err != nil {
		return task.Task{}, fmt.Errorf(
			"crear tarea: %w", err,
		)
	}

	return created, nil

}

func (r *TaskRepository) List(
	ctx context.Context,
) ([]task.Task, error) {
	return r.ListFiltered(ctx, TaskFilter{})
}

func scanTask(
	scanner interface{ Scan(...any) error },
) (task.Task, error) {
	var item task.Task

	var description sql.NullString
	var folderID sql.NullInt64
	var completedAt sql.NullString

	var priority, status string
	var createdAt, updatedAt string

	err := scanner.Scan(
		&item.ID,
		&item.Title,
		&description,
		&priority,
		&status,
		&folderID,
		&createdAt,
		&updatedAt,
		&completedAt,
	)
	if err != nil {
		return task.Task{}, err
	}

	item.Description = description.String
	item.Priority = task.Priority(priority)
	item.Status = task.Status(status)

	if folderID.Valid {
		id := folderID.Int64
		item.FolderID = &id
	}

	item.CreatedAt, err = time.Parse(
		time.RFC3339Nano, createdAt,
	)
	if err != nil {
		return task.Task{}, fmt.Errorf(
			"fecha de creación inválida: %w", err,
		)
	}

	item.UpdatedAt, err = time.Parse(
		time.RFC3339Nano, updatedAt,
	)
	if err != nil {
		return task.Task{}, fmt.Errorf(
			"fecha de actualización inválida: %w", err,
		)
	}

	if completedAt.Valid {
		date, err := time.Parse(
			time.RFC3339Nano, completedAt.String,
		)
		if err != nil {
			return task.Task{}, fmt.Errorf(
				"fecha de finalización inválida: %w", err,
			)
		}

		item.CompletedAt = &date
	}

	return item, nil
}

func (r *TaskRepository) Update(
	ctx context.Context,
	id int64,
	input task.CreateInput,
) (task.Task, error) {
	if id <= 0 {
		return task.Task{}, fmt.Errorf("ID de tarea inválido")
	}

	if err := input.NormalizeAndValidate(); err != nil {
		return task.Task{}, err
	}

	var folderID any
	if input.FolderID != nil {
		folderID = *input.FolderID
	}

	row := r.db.QueryRowContext(ctx, `
		UPDATE tasks
		SET title = ?,
		    description = ?,
		    priority = ?,
		    folder_id = ?,
		    updated_at = strftime(
		        '%Y-%m-%dT%H:%M:%fZ', 'now'
		    )
		WHERE id = ?
		  AND (
		      ? IS NULL
		      OR EXISTS (
		          SELECT 1 FROM folders WHERE id = ?
		      )
		  )
		RETURNING
		    id, title, description, priority, status,
		    folder_id, created_at, updated_at, completed_at
	`,
		input.Title,
		input.Description,
		input.Priority,
		folderID,
		id,
		folderID,
		folderID,
	)

	updated, err := scanTask(row)

	if errors.Is(err, sql.ErrNoRows) {
		if input.FolderID != nil {
			var taskExists bool

			checkErr := r.db.QueryRowContext(ctx, `
				SELECT EXISTS(
					SELECT 1 FROM tasks WHERE id = ?
				)
			`, id).Scan(&taskExists)

			if checkErr != nil {
				return task.Task{}, fmt.Errorf(
					"comprobar tarea: %w", checkErr,
				)
			}

			if taskExists {
				return task.Task{}, ErrTaskFolderNotFound
			}
		}

		return task.Task{}, sql.ErrNoRows
	}

	if err != nil {
		return task.Task{}, fmt.Errorf(
			"actualizar tarea: %w", err,
		)
	}

	return updated, nil
}

// SetCompleted controla manualmente el estado.
// true = completada; false = pendiente.
func (r *TaskRepository) SetCompleted(
	ctx context.Context,
	id int64,
	completed bool,
) (task.Task, error) {
	if id <= 0 {
		return task.Task{}, fmt.Errorf("ID de tarea inválido")
	}

	status := task.StatusPending
	if completed {
		status = task.StatusCompleted
	}

	row := r.db.QueryRowContext(ctx, `
		UPDATE tasks
		SET status = ?,
		    completed_at = CASE
		        WHEN ? = 'completada'
		        THEN COALESCE(
		            completed_at,
		            strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
		        )
		        ELSE NULL
		    END,
		    updated_at = CASE
		        WHEN status <> ?
		        THEN strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
		        ELSE updated_at
		    END
		WHERE id = ?
		RETURNING
		    id, title, description, priority, status,
		    folder_id, created_at, updated_at, completed_at
	`, status, status, status, id)

	updated, err := scanTask(row)
	if err != nil {
		return task.Task{}, fmt.Errorf(
			"cambiar estado de tarea: %w", err,
		)
	}

	return updated, nil
}

// Delete elimina definitivamente una tarea.
func (r *TaskRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return fmt.Errorf("ID de tarea inválido")
	}

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM tasks WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("eliminar tarea: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar eliminación: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
