package entities

import (
	"errors"

	"github.com/google/uuid"
)

type Plant struct {
	ID                   uuid.UUID `json:"id"`
	Title                string    `json:"title"`
	Specie               string    `json:"specie"`
	Location             string    `json:"location"`
	WateringFrequency    int       `json:"watering_frequency"`
	FertilizingFrequency int       `json:"fertilizing_frequency"`
	PruningFrequency     int       `json:"pruning_frequency"`
	User_Id              uuid.UUID `json:"user_id"`
	Garden               uuid.UUID `json:"garden"`
}

type PlantRepository interface {
	Create(plant *Plant) error
	GetByID(id string) (*Plant, error)
	GetByUserID(id string) ([]*Plant, error)
	Update(plant *Plant) error
	Delete(id string) error
}

func NewPlant(title, specie, location string, wateringFrequency, fertilizingFrequency, pruningFrequency int, user_id uuid.UUID, garden uuid.UUID) (*Plant, error) {
	if title == "" {
		return nil, errors.New("expected title")
	}

	if specie == "" {
		return nil, errors.New("expected specie")
	}

	if location == "" {
		return nil, errors.New("expected location")
	}
	if wateringFrequency == 0 || wateringFrequency == ' ' {
		return nil, errors.New("expected wateringFrequency")
	}
	if fertilizingFrequency == 0 || fertilizingFrequency == ' ' {
		return nil, errors.New("expected fertilizingFrequency")
	}

	if pruningFrequency == 0 || pruningFrequency == ' ' {
		return nil, errors.New("expected pruningFrequency")
	}

	if user_id == uuid.Nil {
		return nil, errors.New("expected user_id")
	}

	return &Plant{
		ID:                   uuid.New(),
		Title:                title,
		Specie:               specie,
		Location:             location,
		WateringFrequency:    wateringFrequency,
		FertilizingFrequency: fertilizingFrequency,
		PruningFrequency:     pruningFrequency,
		User_Id:              user_id,
		Garden:               garden,
	}, nil
}
