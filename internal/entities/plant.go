package entities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	usecases_plant "github.com/lucasBiazon/botany-back/internal/usecases/plant"
)

type Plant struct {
	ID                   string    `json:"id"`
	PlantName            string    `json:"plant_name"`
	PlantDescription     string    `json:"plant_description"`
	PlantingDate         time.Time `json:"planting_date"`
	EstimatedHarvestDate time.Time `json:"estimated_harvest_date"`
	PlantStatus          string    `json:"plant_status"`
	CurrentHeight        float64   `json:"current_height"`
	CurrentWidth         float64   `json:"current_width"`
	IrrigationWeek       int       `json:"irrigation_week"`
	HealthStatus         string    `json:"health_status" default:"Healthy"`
	LastIrrigation       time.Time `json:"last_irrigation"`
	LastFertilization    time.Time `json:"last_fertilization"`
	SunExposure          float64   `json:"sun_exposure"`
	FertilizationWeek    float64   `json:"fertilization_week"`
	UserID               string    `json:"user_id"`
	SpeciesID            string    `json:"species_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	CategoriesPlant      []string  `json:"categories_plant"`
}

type PlantRepository interface {
	Create(ctx context.Context, plant *Plant) (string, error)
	FindByID(ctx context.Context, userId, id string) (*Plant, error)
	FindBySpeciesName(ctx context.Context, userId, speciesId string) ([]*Plant, error)
	FindByCategoryName(ctx context.Context, userId, categoryName string) ([]*Plant, error)
	FindByName(ctx context.Context, userId, name string) ([]*Plant, error)
	FindAll(ctx context.Context, userId string) ([]*Plant, error)
	Update(ctx context.Context, plant *Plant) error
	Delete(ctx context.Context, userId, id string) error
}

func NewPlant(input usecases_plant.CreatePlantUseCaseInputDTO) (*Plant, error) {
	if input.PlantName == "" {
		return nil, errors.New("expected plant name")
	}
	if input.PlantDescription == "" {
		return nil, errors.New("expected plant description")
	}
	if input.PlantingDate.IsZero() {
		input.PlantingDate = time.Now()
	}
	if input.EstimatedHarvestDate.IsZero() {
		input.EstimatedHarvestDate = time.Now()
	}
	if input.PlantStatus == "" {
		input.PlantStatus = "Growing"
	}

	if input.CurrentHeight == 0 {
		input.CurrentHeight = 0
	}
	if input.CurrentWidth == 0 {
		input.CurrentWidth = 0
	}
	if input.IrrigationWeek == 0 {
		input.IrrigationWeek = 0
	}
	if input.HealthStatus == "" {
		input.HealthStatus = "Healthy"
	}
	if input.LastIrrigation.IsZero() {
		input.LastIrrigation = time.Now()
	}
	if input.LastFertilization.IsZero() {
		input.LastFertilization = time.Now()
	}
	if input.SunExposure == 0 {
		input.SunExposure = 0
	}
	if input.FertilizationWeek == 0 {
		input.FertilizationWeek = 0
	}
	if input.UserID == "" {
		return nil, errors.New("expected user id")
	}
	if input.SpeciesID == "" {
		return nil, errors.New("expected species id")
	}
	if len(input.CategoriesPlant) == 0 {
		input.CategoriesPlant = []string{}
	}

	return &Plant{
		ID:                   uuid.New().String(),
		PlantName:            input.PlantName,
		PlantDescription:     input.PlantDescription,
		PlantingDate:         input.PlantingDate,
		EstimatedHarvestDate: input.EstimatedHarvestDate,
		PlantStatus:          input.PlantStatus,
		CurrentHeight:        input.CurrentHeight,
		CurrentWidth:         input.CurrentWidth,
		IrrigationWeek:       input.IrrigationWeek,
		HealthStatus:         input.HealthStatus,
		LastIrrigation:       input.LastIrrigation,
		LastFertilization:    input.LastFertilization,
		SunExposure:          input.SunExposure,
		FertilizationWeek:    input.FertilizationWeek,
		UserID:               input.UserID,
		SpeciesID:            input.SpeciesID,
		CategoriesPlant:      input.CategoriesPlant,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}, nil
}

func CalculateEstimatedHarvest(p *Plant, specie *Specie) (time.Time, error) {
	if p.PlantingDate.IsZero() {
		return time.Time{}, errors.New("PlantingDate is not set")
	}

	if specie == nil {
		return time.Time{}, errors.New("Specie is nil")
	}

	if specie.HarvestTime <= 0 {
		return time.Time{}, errors.New("Specie HarvestTime is invalid")
	}

	// Calcula a data estimada de colheita
	estimatedHarvestDate := p.PlantingDate.AddDate(0, 0, specie.HarvestTime)

	// Atualiza a planta com a nova data estimada de colheita
	p.EstimatedHarvestDate = estimatedHarvestDate

	return estimatedHarvestDate, nil
}
