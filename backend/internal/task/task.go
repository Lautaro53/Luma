package task

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

type Priority string

const (
	PriorityLow    Priority = "baja"
	PriorityMedium Priority = "media"
	PriorityHigh   Priority = "alta"
)

type Status string

const (
	StatusPending   Status = "pendiente"
	StatusCompleted Status = "completada"
)

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
	FolderID    *int64     `json:"folder_id"`
}

type CreateInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    Priority `json:"priority"`
	FolderID    *int64   `json:"folder_id"`
}

func (in *CreateInput) NormalizeAndValidate() error {
	in.Title = strings.TrimSpace(in.Title)
	in.Description = strings.TrimSpace(in.Description)

	length := utf8.RuneCountInString(in.Title)

	if length < 1 || length > 200 {
		return errors.New(
			"el título debe contener entre 1 y 200 caracteres",
		)
	}

	if in.Priority == "" {
		in.Priority = PriorityMedium
	}

	if in.FolderID != nil && *in.FolderID <= 0 {
		return errors.New(
			"el identificador de carpeta debe ser positivo",
		)
	}

	switch in.Priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return nil
	default:
		return errors.New("prioridad inválida")
	}
}
