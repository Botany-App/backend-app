package usecases_plant

import (
	"context"
	"log"
	"time"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CreatePlantUseCase struct {
	PlantRepository entities.PlantRepository
}

type CreatePlantUseCaseInputDTO struct {
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
	CategoriesPlant      []string  `json:"categories_plant"`
}

func NewCreatePlantUseCase(plantRepository entities.PlantRepository) *CreatePlantUseCase {
	return &CreatePlantUseCase{
		PlantRepository: plantRepository,
	}
}

func (uc *CreatePlantUseCase) Execute(ctx context.Context, input CreatePlantUseCaseInputDTO) (*entities.Plant, error) {
	log.Println("CreatePlantUseCase - Execute")
	newPlant, err := entities.NewPlant(input)
	if err != nil {
		return nil, err
	}

	specie, err := uc.PlantRepository.FindByID(ctx, input.UserID, input.SpeciesID)
	if err != nil {
		return nil, err
	}

	estimatedHarvestDate, err := entities.CalculateEstimatedHarvest(newPlant, specie)
	if err != nil {
		return nil, err
	}
	newPlant.EstimatedHarvestDate = estimatedHarvestDate
	id, err := uc.PlantRepository.Create(ctx, newPlant)
	if err != nil {
		return nil, err
	}

	plantCreated, err := uc.PlantRepository.FindByID(ctx, input.UserID, id)
	if err != nil {
		return nil, err
	}
	return plantCreated, nil
}
