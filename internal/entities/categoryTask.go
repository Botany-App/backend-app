package entities

import (
	"errors"

	"github.com/google/uuid"
)

type CategoryTask struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CategoryTaskRepository interface {
	GetAll(userID string) ([]CategoryTask, error)
	GetByName(userID, name string) ([]CategoryTask, error)
	GetByID(userID, id string) ([]CategoryTask, error)
	Create(userID string, category *CategoryTask) error
	Update(userID string, category *CategoryTask) error
	Delete(userID, id string) error
}

func NewCategoryTask(name, description string) (*CategoryTask, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if description == "" {
		description = "No description"
	}

	return &CategoryTask{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
	}, nil
}
