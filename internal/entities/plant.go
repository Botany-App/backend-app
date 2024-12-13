package entities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Plant struct {
	Id                   string    `json:"id"`
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
	UserId               string    `json:"user_id"`
	SpeciesId            string    `json:"species_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	CategoriesPlant      []string  `json:"categories_plant"`
}

type PlantWithCategory struct {
	PlantId              string
	PlantName            string
	PlantDescription     string
	PlantingDate         time.Time
	EstimatedHarvestDate time.Time
	PlantStatus          string
	CurrentHeight        float64
	CurrentWidth         float64
	IrrigationWeek       int
	HealthStatus         string
	LastIrrigation       time.Time
	LastFertilization    time.Time
	SunExposure          float64
	FertilizationWeek    float64
	UserId               string
	SpeciesId            string
	PlantCreatedAt       time.Time
	PlantUpdatedAt       time.Time
	Category             []CategoryPlant
}

type PlantRepository interface {
	Create(ctx context.Context, plant *Plant) (string, error)
	FindByID(ctx context.Context, userId, id string) (*PlantWithCategory, error)
	FindBySpeciesName(ctx context.Context, userId, speciesId string) ([]*PlantWithCategory, error)
	FindByCategoryName(ctx context.Context, userId, categoryName string) ([]*PlantWithCategory, error)
	FindByName(ctx context.Context, userId, name string) ([]*PlantWithCategory, error)
	FindAll(ctx context.Context, userId string) ([]*PlantWithCategory, error)
	Update(ctx context.Context, plant *Plant) error
	Delete(ctx context.Context, userId, id string) error
}

func NewPlant(
	plantName string,
	plantDescription string,
	plantingDate time.Time,
	estimatedHarvestDate time.Time,
	plantStatus string,
	currentHeight float64,
	currentWidth float64,
	irrigationWeek int,
	healthStatus string,
	lastIrrigation time.Time,
	lastFertilization time.Time,
	sunExposure float64,
	fertilizationWeek float64,
	userId string,
	speciesId string,
	categoriesPlant []string,
) (*Plant, error) {
	if plantName == "" {
		return nil, errors.New("expected plant name")
	}
	if plantDescription == "" {
		return nil, errors.New("expected plant description")
	}
	if plantingDate.IsZero() {
		plantingDate = time.Now()
	}
	if estimatedHarvestDate.IsZero() {
		estimatedHarvestDate = time.Now()
	}
	if plantStatus == "" {
		plantStatus = "Growing"
	}

	if currentHeight == 0 {
		currentHeight = 0
	}
	if currentWidth == 0 {
		currentWidth = 0
	}
	if irrigationWeek == 0 {
		irrigationWeek = 0
	}
	if healthStatus == "" {
		healthStatus = "Healthy"
	}
	if lastIrrigation.IsZero() {
		lastIrrigation = time.Now()
	}
	if lastFertilization.IsZero() {
		lastFertilization = time.Now()
	}
	if sunExposure == 0 {
		sunExposure = 0
	}
	if fertilizationWeek == 0 {
		fertilizationWeek = 0
	}
	if userId == "" {
		return nil, errors.New("expected user id")
	}
	if speciesId == "" {
		return nil, errors.New("expected species id")
	}

	return &Plant{
		Id:                   uuid.New().String(),
		PlantName:            plantName,
		PlantDescription:     plantDescription,
		PlantingDate:         plantingDate,
		EstimatedHarvestDate: estimatedHarvestDate,
		PlantStatus:          plantStatus,
		CurrentHeight:        currentHeight,
		CurrentWidth:         currentWidth,
		IrrigationWeek:       irrigationWeek,
		HealthStatus:         healthStatus,
		LastIrrigation:       lastIrrigation,
		LastFertilization:    lastFertilization,
		SunExposure:          sunExposure,
		FertilizationWeek:    fertilizationWeek,
		UserId:               userId,
		SpeciesId:            speciesId,
		CategoriesPlant:      categoriesPlant,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}, nil
}
