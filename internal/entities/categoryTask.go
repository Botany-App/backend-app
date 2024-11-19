package entities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type CategoryTask struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryTaskRepository interface {
	FindAll(ctx context.Context, userID string, limit, offset int) ([]*CategoryTask, error)
	FindByName(ctx context.Context, userID, name string, limit, offset int) ([]*CategoryTask, error)
	FindByID(ctx context.Context, userID, id string, limit, offset int) (*CategoryTask, error)
	Create(ctx context.Context, category *CategoryTask) error
	Update(ctx context.Context, category *CategoryTask) error
	Delete(ctx context.Context, userID, id string) error
}

func NewCategoryTask(name, description, userID string) (*CategoryTask, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if description == "" {
		description = "No description"
	}

	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	return &CategoryTask{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
