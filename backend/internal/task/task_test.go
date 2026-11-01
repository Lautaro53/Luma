package task

import (
	"strings"
	"testing"
)

func TestCreateInputValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateInput
		wantErr bool
	}{
		{
			name: "tarea válida",
			input: CreateInput{
				Title:    "Estudiar Go",
				Priority: PriorityHigh,
			},
		},
		{
			name: "título vacío",
			input: CreateInput{
				Title: "   ",
			},
			wantErr: true,
		},
		{
			name: "título demasiado largo",
			input: CreateInput{
				Title: strings.Repeat("a", 201),
			},
			wantErr: true,
		},
		{
			name: "prioridad inválida",
			input: CreateInput{
				Title:    "Estudiar",
				Priority: "urgente",
			},
			wantErr: true,
		},
		{
			name: "prioridad predeterminada",
			input: CreateInput{
				Title: "Organizar apuntes",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.NormalizeAndValidate()

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"error obtenido: %v, esperado error: %v",
					err,
					tt.wantErr,
				)
			}

			if err == nil && tt.input.Priority == "" {
				t.Fatal("la prioridad no fue establecida")
			}
		})
	}
}

func TestTitleNormalization(t *testing.T) {
	input := CreateInput{
		Title: "  Estudiar Go  ",
	}

	if err := input.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}

	if input.Title != "Estudiar Go" {
		t.Fatalf("título inesperado: %q", input.Title)
	}

	if input.Priority != PriorityMedium {
		t.Fatalf(
			"prioridad esperada: media, obtenida: %s",
			input.Priority,
		)
	}
}
