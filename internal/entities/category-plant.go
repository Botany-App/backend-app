package entities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type CategoryPlant struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserId      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryPlantRepository interface {
	Create(ctx context.Context, categoryPlant *CategoryPlant) (string, error)
	FindById(ctx context.Context, userId, id string) (*CategoryPlant, error)
	FindAll(ctx context.Context, userId string) ([]*CategoryPlant, error)
	FindByName(ctx context.Context, userId, name string) ([]*CategoryPlant, error)
	Update(ctx context.Context, categoryPlant *CategoryPlant) error
	Delete(ctx context.Context, userId, id string) error
}

func NewCreateCategoryPlant(name, description, userId string) (*CategoryPlant, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if description == "" {
		description = "Sem descrição"
	}
	if userId == "" {
		return nil, errors.New("user_id is required")
	}

	return &CategoryPlant{
		Id:          uuid.New().String(),
		Name:        name,
		Description: description,
		UserId:      userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
