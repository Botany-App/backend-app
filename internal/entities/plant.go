package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Plant struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Specie             uuid.UUID `json:"specie"`
	Location           string    `json:"location"`
	PlantingTime       time.Time `json:"planting_time"`
	IrrigationWeek     int       `json:"irrigation_week"`
	FertilizantionWeek int       `json:"fertilizantion_week"`
	SumExposition      float32   `json:"sum_exposition"`
	User_Id            uuid.UUID `json:"user_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type PlantRepository interface {
	Create(plant *Plant) error
	GetByID(id string) (*Plant, error)
	GetByUserID(id string) ([]*Plant, error)
	GetByPlantingTime(plantingTime time.Time) ([]*Plant, error)
	GetBySpecie(specie string) ([]*Plant, error)
	GetByName(name string) ([]*Plant, error)
	GetByLocation(location string) ([]*Plant, error)
	GetByCategory(category string) ([]*Plant, error)
	Update(plant *Plant) error
	Delete(id string) error
}

func NewPlant(name, location, description string, userID, specieID uuid.UUID, irrigationWeek, fertilizationWeek int, sumExpositon float32, platingTime time.Time) (*Plant, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if location == "" {
		return nil, errors.New("location is required")
	}

	if description == "" {
		return nil, errors.New("description is required")
	}

	if userID == uuid.Nil {
		return nil, errors.New("user_id is required")
	}

	if specieID == uuid.Nil {
		return nil, errors.New("specie is required")
	}

	if irrigationWeek < 0 {
		return nil, errors.New("irrigation_week is required")
	}

	if fertilizationWeek < 0 {
		return nil, errors.New("fertilization_week is required")
	}

	if sumExpositon < 0 {
		return nil, errors.New("sum_exposition is required")
	}

	if platingTime.IsZero() {
		platingTime = time.Now()
	}
	return &Plant{
		ID:                 uuid.New(),
		Name:               name,
		Description:        description,
		Specie:             specieID,
		Location:           location,
		PlantingTime:       platingTime,
		IrrigationWeek:     irrigationWeek,
		FertilizantionWeek: fertilizationWeek,
		SumExposition:      sumExpositon,
		User_Id:            userID,
	}, nil
}
