package repository

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"luma/internal/task"

	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	path := filepath.Join(t.TempDir(), "test.db")

	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}

	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })

	if _, err := db.Exec(
		"PRAGMA foreign_keys = ON",
	); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{
		"001_create_tasks.sql",
		"002_create_folders.sql",
	} {
		script, err := os.ReadFile(
			filepath.Join("..", "..", "migrations", name),
		)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := db.Exec(string(script)); err != nil {
			t.Fatalf("migración %s: %v", name, err)
		}
	}

	return db
}

func TestFolderCreateAndList(t *testing.T) {
	repo := NewFolderRepository(newTestDB(t))
	ctx := context.Background()

	created, err := repo.Create(ctx, task.CreateFolderInput{
		Name: "  Matemáticas  ",
	})
	if err != nil {
		t.Fatal(err)
	}

	if created.ID <= 0 || created.Name != "Matemáticas" {
		t.Fatalf("carpeta inesperada: %+v", created)
	}

	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Fatal("faltan las fechas de la carpeta")
	}

	_, err = repo.Create(ctx, task.CreateFolderInput{
		Name: "Infraestructura",
	})
	if err != nil {
		t.Fatal(err)
	}

	folders, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(folders) != 2 {
		t.Fatalf("esperadas 2 carpetas, obtenidas %d", len(folders))
	}

	if folders[0].Name != "Infraestructura" ||
		folders[1].Name != "Matemáticas" {
		t.Fatalf("orden inesperado: %+v", folders)
	}
}

func TestFolderInvalidAndDuplicate(t *testing.T) {
	repo := NewFolderRepository(newTestDB(t))
	ctx := context.Background()

	_, err := repo.Create(ctx, task.CreateFolderInput{
		Name: "   ",
	})
	if err == nil {
		t.Fatal("se aceptó un nombre vacío")
	}

	_, err = repo.Create(ctx, task.CreateFolderInput{
		Name: "Infraestructura",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(ctx, task.CreateFolderInput{
		Name: "infraestructura",
	})
	if err == nil {
		t.Fatal("se aceptó una carpeta duplicada")
	}
}

func TestFolderRename(t *testing.T) {
	repo := NewFolderRepository(newTestDB(t))
	ctx := context.Background()

	folder, err := repo.Create(ctx, task.CreateFolderInput{
		Name: "Infraestructura",
	})
	if err != nil {
		t.Fatal(err)
	}

	updated, err := repo.Rename(
		ctx,
		folder.ID,
		"  Redes y servidores  ",
	)
	if err != nil {
		t.Fatal(err)
	}

	if updated.ID != folder.ID {
		t.Fatal("el identificador cambió")
	}

	if updated.Name != "Redes y servidores" {
		t.Fatalf("nombre inesperado: %s", updated.Name)
	}

	folders, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(folders) != 1 ||
		folders[0].Name != "Redes y servidores" {
		t.Fatal("el cambio no se guardó")
	}
}

func TestFolderDeletePreservesTasks(t *testing.T) {
	db := newTestDB(t)
	repo := NewFolderRepository(db)
	ctx := context.Background()

	folder, err := repo.Create(ctx, task.CreateFolderInput{
		Name: "Infraestructura",
	})
	if err != nil {
		t.Fatal(err)
	}

	var foreignKeys int

	if err := db.QueryRow(
		"PRAGMA foreign_keys",
	).Scan(&foreignKeys); err != nil {
		t.Fatal(err)
	}

	if foreignKeys != 1 {
		t.Fatal("las claves externas no están activadas")
	}

	_, err = db.Exec(`
		INSERT INTO tasks (title, folder_id)
		VALUES (?, ?)
	`, "Estudiar OSPF", folder.ID)

	if err != nil {
		t.Fatal(err)
	}

	if err := repo.Delete(ctx, folder.ID); err != nil {
		t.Fatal(err)
	}

	var taskTitle string
	var folderID sql.NullInt64

	err = db.QueryRow(`
		SELECT title, folder_id
		FROM tasks
		WHERE title = ?
	`, "Estudiar OSPF").Scan(&taskTitle, &folderID)

	if err != nil {
		t.Fatalf("la tarea desapareció: %v", err)
	}

	if folderID.Valid {
		t.Fatalf(
			"la tarea conserva una carpeta eliminada: %d",
			folderID.Int64,
		)
	}

	folders, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(folders) != 0 {
		t.Fatal("la carpeta no fue eliminada")
	}
}

func TestFolderRenameDuplicate(t *testing.T) {
	repo := NewFolderRepository(newTestDB(t))
	ctx := context.Background()

	first, err := repo.Create(ctx, task.CreateFolderInput{
		Name: "Infraestructura",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Create(ctx, task.CreateFolderInput{
		Name: "Matemáticas",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Rename(
		ctx, first.ID, "matemáticas",
	)

	if err != ErrFolderAlreadyExists {
		t.Fatalf(
			"se esperaba conflicto, obtenido: %v", err,
		)
	}

	folders, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(folders) != 2 ||
		folders[0].Name != "Infraestructura" ||
		folders[1].Name != "Matemáticas" {
		t.Fatalf("se modificaron datos: %+v", folders)
	}
}
