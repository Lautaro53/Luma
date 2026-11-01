package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

func Open() (*sql.DB, error) {
	// Directorio de datos locales.
	if err := os.MkdirAll("data", 0700); err != nil {
		return nil, fmt.Errorf("crear directorio de datos: %w", err)
	}

	dsn := "file:data/luma.db?_fk=1&_busy_timeout=5000"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir SQLite: %w", err)
	}

	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("conectar SQLite: %w", err)
	}

	return db, nil
}
