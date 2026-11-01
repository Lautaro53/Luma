package task

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

type Folder struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateFolderInput struct {
	Name string `json:"name"`
}

func (in *CreateFolderInput) NormalizeAndValidate() error {
	in.Name = strings.TrimSpace(in.Name)

	length := utf8.RuneCountInString(in.Name)

	if length < 1 || length > 80 {
		return errors.New(
			"el nombre debe contener entre 1 y 80 caracteres",
		)
	}

	return nil
}
