package repository

import (
	"context"
	"testing"

	"luma/internal/task"
)

func TestTaskCreateAndList(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	folders := NewFolderRepository(db)
	tasks := NewTaskRepository(db)

	folder, err := folders.Create(
		ctx,
		task.CreateFolderInput{
			Name: "Infraestructura",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	linked, err := tasks.Create(ctx, task.CreateInput{
		Title:       "  Configurar VPN  ",
		Description: "  Laboratorio de redes  ",
		Priority:    task.PriorityHigh,
		FolderID:    &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if linked.Title != "Configurar VPN" ||
		linked.Description != "Laboratorio de redes" {
		t.Fatal("la tarea no fue normalizada")
	}

	if linked.FolderID == nil ||
		*linked.FolderID != folder.ID {
		t.Fatal("asociación de carpeta incorrecta")
	}

	if linked.Status != task.StatusPending ||
		linked.CompletedAt != nil {
		t.Fatal("la tarea no comenzó pendiente")
	}

	if linked.CreatedAt.IsZero() ||
		linked.UpdatedAt.IsZero() {
		t.Fatal("faltan las fechas de creación")
	}

	free, err := tasks.Create(ctx, task.CreateInput{
		Title: "Comprar cuaderno",
	})
	if err != nil {
		t.Fatal(err)
	}

	if free.FolderID != nil {
		t.Fatal("la tarea debería estar sin carpeta")
	}

	if free.Priority != task.PriorityMedium {
		t.Fatal("prioridad predeterminada incorrecta")
	}

	list, err := tasks.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 2 {
		t.Fatalf(
			"esperadas 2 tareas, obtenidas %d",
			len(list),
		)
	}

	if list[0].ID != free.ID ||
		list[1].ID != linked.ID {
		t.Fatal("orden de tareas incorrecto")
	}
}

func TestTaskInvalidFolder(t *testing.T) {
	db := newTestDB(t)
	repo := NewTaskRepository(db)
	ctx := context.Background()

	invalidID := int64(9999)

	_, createErr := repo.Create(ctx, task.CreateInput{
		Title:    "Estudiar matemáticas",
		FolderID: &invalidID,
	})

	if createErr != ErrTaskFolderNotFound {
		t.Fatalf(
			"se esperaba ErrTaskFolderNotFound, obtenido: %v",
			createErr,
		)
	}

	list, listErr := repo.List(ctx)
	if listErr != nil {
		t.Fatal(listErr)
	}

	if len(list) != 0 {
		t.Fatal("se guardó una tarea con carpeta inválida")
	}
}

func TestTaskUpdateInvalidFolder(t *testing.T) {
	repo := NewTaskRepository(newTestDB(t))
	ctx := context.Background()

	created, err := repo.Create(ctx, task.CreateInput{
		Title: "Estudiar Go",
	})
	if err != nil {
		t.Fatal(err)
	}

	invalidID := int64(9999)

	_, err = repo.Update(ctx, created.ID, task.CreateInput{
		Title:    "Título modificado",
		FolderID: &invalidID,
	})

	if err != ErrTaskFolderNotFound {
		t.Fatalf(
			"se esperaba ErrTaskFolderNotFound, obtenido: %v",
			err,
		)
	}

	items, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 ||
		items[0].Title != "Estudiar Go" ||
		items[0].FolderID != nil {
		t.Fatal("la edición inválida modificó la tarea")
	}
}

func TestTaskUpdateAndCompletion(t *testing.T) {
	db := newTestDB(t)
	repo := NewTaskRepository(db)
	folders := NewFolderRepository(db)
	ctx := context.Background()

	created, err := repo.Create(ctx, task.CreateInput{
		Title: "Estudiar Go",
	})
	if err != nil {
		t.Fatal(err)
	}

	folder, err := folders.Create(
		ctx,
		task.CreateFolderInput{Name: "Programación"},
	)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := repo.Update(ctx, created.ID, task.CreateInput{
		Title:       "Repasar Go",
		Description: "Ejercicios de interfaces",
		Priority:    task.PriorityHigh,
		FolderID:    &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if updated.Title != "Repasar Go" ||
		updated.Priority != task.PriorityHigh ||
		updated.FolderID == nil ||
		*updated.FolderID != folder.ID {
		t.Fatal("la edición no se guardó correctamente")
	}

	completed, err := repo.SetCompleted(
		ctx, created.ID, true,
	)
	if err != nil {
		t.Fatal(err)
	}

	if completed.Status != task.StatusCompleted ||
		completed.CompletedAt == nil {
		t.Fatal("la tarea no fue completada")
	}

	again, err := repo.SetCompleted(
		ctx, created.ID, true,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !again.CompletedAt.Equal(*completed.CompletedAt) {
		t.Fatal("la finalización cambió al repetir la acción")
	}

	reopened, err := repo.SetCompleted(
		ctx, created.ID, false,
	)
	if err != nil {
		t.Fatal(err)
	}

	if reopened.Status != task.StatusPending ||
		reopened.CompletedAt != nil {
		t.Fatal("la tarea no volvió a pendiente")
	}
}

func TestTaskDelete(t *testing.T) {
	repo := NewTaskRepository(newTestDB(t))
	ctx := context.Background()

	created, err := repo.Create(ctx, task.CreateInput{
		Title: "Tarea temporal",
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatal(err)
	}

	items, err := repo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 0 {
		t.Fatal("la tarea continúa almacenada")
	}

	if err := repo.Delete(ctx, created.ID); err == nil {
		t.Fatal("se aceptó eliminar una tarea inexistente")
	}
}
