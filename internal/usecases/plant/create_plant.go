package usecases_plant

import (
	"context"
	"log"
	"time"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CreatePlantUseCase struct {
	PlantRepository  entities.PlantRepository
	SpecieRepository entities.SpecieRepository
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
	SpecieHaverstTime    int       `json:"specie_harvest_time"`
}

func NewCreatePlantUseCase(plantRepository entities.PlantRepository, specieRepository entities.SpecieRepository) *CreatePlantUseCase {
	return &CreatePlantUseCase{
		PlantRepository:  plantRepository,
		SpecieRepository: specieRepository,
	}
}

func CalculateEstimatedHarvestDate(plantingDate, estimatedHarvestDate time.Time, harvestTime int) time.Time {
	// Se a data estimada de colheita já existir, ajusta com base no HarvestTime
	if !estimatedHarvestDate.IsZero() {
		return estimatedHarvestDate.AddDate(0, 0, harvestTime)
	}
	// Caso contrário, calcula a partir da data de plantio
	return plantingDate.AddDate(0, 0, harvestTime)
}

func (uc *CreatePlantUseCase) Execute(ctx context.Context, input CreatePlantUseCaseInputDTO) (*entities.PlantWithCategory, error) {
	log.Println("CreatePlantUseCase - Execute")
	newPlant, err := entities.NewPlant(
		input.PlantName,
		input.PlantDescription,
		input.PlantingDate,
		input.EstimatedHarvestDate,
		input.PlantStatus,
		input.CurrentHeight,
		input.CurrentWidth,
		input.IrrigationWeek,
		input.HealthStatus,
		input.LastIrrigation,
		input.LastFertilization,
		input.SunExposure,
		input.FertilizationWeek,
		input.UserID,
		input.SpeciesID,
		input.CategoriesPlant,
	)

	if err != nil {
		return nil, err
	}

	estimatedHarvestDate := CalculateEstimatedHarvestDate(newPlant.PlantingDate, newPlant.EstimatedHarvestDate, input.SpecieHaverstTime)
	log.Print("Estimated Harvest Date: ", estimatedHarvestDate)
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
