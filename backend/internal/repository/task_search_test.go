package repository

import (
	"context"
	"slices"
	"testing"

	"luma/internal/task"
)

func TestTaskSearchAndSort(t *testing.T) {
	repo := NewTaskRepository(newTestDB(t))
	ctx := context.Background()

	create := func(
		title, description string,
		priority task.Priority,
	) task.Task {
		t.Helper()

		item, err := repo.Create(ctx, task.CreateInput{
			Title:       title,
			Description: description,
			Priority:    priority,
		})
		if err != nil {
			t.Fatal(err)
		}

		return item
	}

	low := create(
		"Estudiar redes", "Repasar OSPF",
		task.PriorityLow,
	)

	high1 := create(
		"Resolver ficha", "Matemáticas",
		task.PriorityHigh,
	)

	medium := create(
		"Revisar 100%", "Notas del curso",
		task.PriorityMedium,
	)

	high2 := create(
		"Preparar informe", "Proyecto",
		task.PriorityHigh,
	)

	tests := []struct {
		name   string
		filter TaskFilter
		want   []int64
	}{
		{
			name:   "buscar por título",
			filter: TaskFilter{Search: "FiChA"},
			want:   []int64{high1.ID},
		},
		{
			name:   "buscar por descripción",
			filter: TaskFilter{Search: "ospf"},
			want:   []int64{low.ID},
		},
		{
			name:   "buscar porcentaje literal",
			filter: TaskFilter{Search: "%"},
			want:   []int64{medium.ID},
		},
		{
			name: "más recientes",
			filter: TaskFilter{
				Sort: SortNewest,
			},
			want: []int64{
				high2.ID, medium.ID, high1.ID, low.ID,
			},
		},
		{
			name: "más antiguas",
			filter: TaskFilter{
				Sort: SortOldest,
			},
			want: []int64{
				low.ID, high1.ID, medium.ID, high2.ID,
			},
		},
		{
			name: "por prioridad",
			filter: TaskFilter{
				Sort: SortPriority,
			},
			want: []int64{
				high2.ID, high1.ID, medium.ID, low.ID,
			},
		},
		{
			name: "búsqueda y prioridad combinadas",
			filter: TaskFilter{
				Search:   "redes",
				Priority: task.PriorityLow,
			},
			want: []int64{low.ID},
		},
		{
			name: "sin coincidencias",
			filter: TaskFilter{
				Search: "inexistente",
			},
			want: []int64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.ListFiltered(
				ctx, tt.filter,
			)
			if err != nil {
				t.Fatal(err)
			}

			var ids []int64
			for _, item := range result {
				ids = append(ids, item.ID)
			}

			if !slices.Equal(ids, tt.want) {
				t.Fatalf(
					"orden esperado: %v, obtenido: %v",
					tt.want, ids,
				)
			}
		})
	}
}

func TestTaskInvalidSort(t *testing.T) {
	filter := TaskFilter{
		Sort: "random",
	}

	if err := filter.Validate(); err == nil {
		t.Fatal("se aceptó un ordenamiento inválido")
	}
}
