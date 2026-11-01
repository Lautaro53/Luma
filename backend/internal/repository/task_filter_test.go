package repository

import (
	"context"
	"testing"

	"luma/internal/task"
)

func TestTaskFilters(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()

	tasks := NewTaskRepository(db)
	folders := NewFolderRepository(db)

	folder, err := folders.Create(
		ctx,
		task.CreateFolderInput{Name: "Infraestructura"},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tasks.Create(ctx, task.CreateInput{
		Title:    "Configurar VPN",
		Priority: task.PriorityHigh,
		FolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = tasks.Create(ctx, task.CreateInput{
		Title: "Comprar cuaderno",
	})
	if err != nil {
		t.Fatal(err)
	}

	completed, err := tasks.Create(ctx, task.CreateInput{
		Title:    "Repasar OSPF",
		FolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = tasks.SetCompleted(ctx, completed.ID, true)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		filter TaskFilter
		want   int
	}{
		{"todas", TaskFilter{}, 3},
		{"pendientes", TaskFilter{
			Status: task.StatusPending,
		}, 2},
		{"completadas", TaskFilter{
			Status: task.StatusCompleted,
		}, 1},
		{"prioridad alta", TaskFilter{
			Priority: task.PriorityHigh,
		}, 1},
		{"carpeta específica", TaskFilter{
			FolderID: &folder.ID,
		}, 2},
		{"sin carpeta", TaskFilter{
			WithoutFolder: true,
		}, 1},
		{"filtros combinados", TaskFilter{
			Status:   task.StatusPending,
			Priority: task.PriorityHigh,
			FolderID: &folder.ID,
		}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tasks.ListFiltered(ctx, tt.filter)
			if err != nil {
				t.Fatal(err)
			}

			if len(result) != tt.want {
				t.Fatalf(
					"esperadas %d tareas, obtenidas %d",
					tt.want, len(result),
				)
			}
		})
	}
}

func TestTaskFilterValidation(t *testing.T) {
	id := int64(1)

	tests := []TaskFilter{
		{Status: "desconocido"},
		{Priority: "urgente"},
		{WithoutFolder: true, FolderID: &id},
	}

	for _, filter := range tests {
		if err := filter.Validate(); err == nil {
			t.Fatalf("se aceptó un filtro inválido: %+v", filter)
		}
	}
}
