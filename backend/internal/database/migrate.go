package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

type migration struct {
	version  int
	filename string
}

var migrations = []migration{
	{1, "001_create_tasks.sql"},
	{2, "002_create_folders.sql"},
}

// Migrate aplica únicamente las migraciones pendientes.
func Migrate(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar migraciones: %w", err)
	}
	defer tx.Rollback()

	// Tabla de control de versiones.
	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TEXT NOT NULL DEFAULT
				(strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
		)
	`)
	if err != nil {
		return fmt.Errorf("crear control de migraciones: %w", err)
	}

	for _, m := range migrations {
		var count int

		err := tx.QueryRow(
			"SELECT COUNT(*) FROM schema_migrations WHERE version = ?",
			m.version,
		).Scan(&count)

		if err != nil {
			return fmt.Errorf(
				"consultar migración %d: %w", m.version, err,
			)
		}

		if count > 0 {
			continue
		}

		path := filepath.Join("migrations", m.filename)

		script, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("leer %s: %w", path, err)
		}

		if _, err := tx.Exec(string(script)); err != nil {
			return fmt.Errorf(
				"ejecutar migración %d: %w",
				m.version, err,
			)
		}

		_, err = tx.Exec(
			`INSERT INTO schema_migrations (version, name)
			 VALUES (?, ?)`,
			m.version,
			m.filename,
		)
		if err != nil {
			return fmt.Errorf(
				"registrar migración %d: %w",
				m.version, err,
			)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar migraciones: %w", err)
	}

	return nil
}
