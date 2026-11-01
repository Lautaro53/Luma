package task

import (
	"strings"
	"testing"
)

func TestCreateFolderValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateFolderInput
		wantErr bool
	}{
		{
			name: "carpeta válida",
			input: CreateFolderInput{
				Name: "Infraestructura",
			},
		},
		{
			name: "nombre vacío",
			input: CreateFolderInput{
				Name: "   ",
			},
			wantErr: true,
		},
		{
			name: "nombre demasiado largo",
			input: CreateFolderInput{
				Name: strings.Repeat("a", 81),
			},
			wantErr: true,
		},
		{
			name: "normalización",
			input: CreateFolderInput{
				Name: "  Matemáticas  ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.NormalizeAndValidate()

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"error obtenido: %v, esperado: %v",
					err,
					tt.wantErr,
				)
			}

			if tt.name == "normalización" &&
				tt.input.Name != "Matemáticas" {
				t.Fatalf(
					"nombre inesperado: %q",
					tt.input.Name,
				)
			}
		})
	}
}

func TestCreateTaskFolderID(t *testing.T) {
	// Sin carpeta debe ser válido.
	input := CreateInput{
		Title: "Organizar apuntes",
	}

	if err := input.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}

	id := int64(1)
	input.FolderID = &id

	if err := input.NormalizeAndValidate(); err != nil {
		t.Fatal(err)
	}

	invalidID := int64(0)
	input.FolderID = &invalidID

	if err := input.NormalizeAndValidate(); err == nil {
		t.Fatal("se aceptó un identificador inválido")
	}
}
