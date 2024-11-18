package entities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type CategoryPlant struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryPlantRepository interface {
	Create(ctx context.Context, categoryPlant *CategoryPlant) error
	FindAll(ctx context.Context, userID uuid.UUID) ([]*CategoryPlant, error)
	FindByID(ctx context.Context, userID, id uuid.UUID) (*CategoryPlant, error)
	FindByName(ctx context.Context, userID uuid.UUID, name string) ([]*CategoryPlant, error)
	Update(ctx context.Context, categoryPlant *CategoryPlant) error
	Delete(ctx context.Context, userID, id uuid.UUID) error
}

func NewCategoryPlant(name, description string, userID uuid.UUID) (*CategoryPlant, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if description == "" {
		description = "No description"
	}

	if userID == uuid.Nil {
		return nil, errors.New("user_id is required")
	}

	return &CategoryPlant{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
