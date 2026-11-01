package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"luma/internal/task"
)

var ErrFolderAlreadyExists = errors.New(
	"ya existe una carpeta con ese nombre",
)

type FolderRepository struct {
	db *sql.DB
}

func NewFolderRepository(db *sql.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) Create(
	ctx context.Context,
	input task.CreateFolderInput,
) (task.Folder, error) {
	if err := input.NormalizeAndValidate(); err != nil {
		return task.Folder{}, err
	}

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO folders (name)
		VALUES (?)
		ON CONFLICT(name) DO NOTHING
		RETURNING id, name, created_at, updated_at
	`, input.Name)

	folder, err := scanFolder(row)

	if errors.Is(err, sql.ErrNoRows) {
		return task.Folder{}, ErrFolderAlreadyExists
	}

	if err != nil {
		return task.Folder{}, fmt.Errorf(
			"crear carpeta: %w", err,
		)
	}

	return folder, nil

}

func (r *FolderRepository) List(
	ctx context.Context,
) ([]task.Folder, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, created_at, updated_at
		FROM folders
		ORDER BY name COLLATE NOCASE ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("listar carpetas: %w", err)
	}
	defer rows.Close()

	folders := make([]task.Folder, 0)

	for rows.Next() {
		folder, err := scanFolder(rows)
		if err != nil {
			return nil, fmt.Errorf(
				"leer carpeta: %w", err,
			)
		}

		folders = append(folders, folder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return folders, nil
}

func scanFolder(
	scanner interface{ Scan(...any) error },
) (task.Folder, error) {
	var folder task.Folder
	var createdAt, updatedAt string

	err := scanner.Scan(
		&folder.ID,
		&folder.Name,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return task.Folder{}, err
	}

	folder.CreatedAt, err = time.Parse(
		time.RFC3339Nano, createdAt,
	)
	if err != nil {
		return task.Folder{}, err
	}

	folder.UpdatedAt, err = time.Parse(
		time.RFC3339Nano, updatedAt,
	)
	if err != nil {
		return task.Folder{}, err
	}

	return folder, nil
}

func (r *FolderRepository) Rename(
	ctx context.Context,
	id int64,
	name string,
) (task.Folder, error) {
	if id <= 0 {
		return task.Folder{}, fmt.Errorf(
			"identificador de carpeta inválido",
		)
	}

	input := task.CreateFolderInput{Name: name}

	if err := input.NormalizeAndValidate(); err != nil {
		return task.Folder{}, err
	}

	row := r.db.QueryRowContext(ctx, `
		UPDATE OR IGNORE folders
		SET name = ?,
		    updated_at = strftime(
		        '%Y-%m-%dT%H:%M:%fZ', 'now'
		    )
		WHERE id = ?
		RETURNING id, name, created_at, updated_at
	`, input.Name, id)

	folder, err := scanFolder(row)

	if errors.Is(err, sql.ErrNoRows) {

		var exists int

		checkErr := r.db.QueryRowContext(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM folders WHERE id = ?
			)
		`, id).Scan(&exists)

		if checkErr != nil {
			return task.Folder{}, fmt.Errorf(
				"comprobar carpeta: %w", checkErr,
			)
		}

		if exists == 1 {
			return task.Folder{}, ErrFolderAlreadyExists
		}

		return task.Folder{}, sql.ErrNoRows
	}

	if err != nil {
		return task.Folder{}, fmt.Errorf(
			"renombrar carpeta: %w", err,
		)
	}

	return folder, nil
}

func (r *FolderRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return fmt.Errorf(
			"identificador de carpeta inválido",
		)
	}

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM folders
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("eliminar carpeta: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("comprobar eliminación: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
