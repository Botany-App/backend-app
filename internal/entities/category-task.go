package entities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type CategoryTask struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserId      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryTaskRepository interface {
	Create(ctx context.Context, categoryTask *CategoryTask) (string, error)
	FindById(ctx context.Context, userId, id string) (*CategoryTask, error)
	FindAll(ctx context.Context, userId string) ([]*CategoryTask, error)
	FindByName(ctx context.Context, userId, name string) ([]*CategoryTask, error)
	Update(ctx context.Context, categoryTask *CategoryTask) error
	Delete(ctx context.Context, userId, id string) error
}

func NewCreateCategoryTask(name, description, userId string) (*CategoryTask, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if description == "" {
		description = "Sem descrição"
	}
	if userId == "" {
		return nil, errors.New("user_id is required")
	}

	return &CategoryTask{
		Id:          uuid.New().String(),
		Name:        name,
		Description: description,
		UserId:      userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
