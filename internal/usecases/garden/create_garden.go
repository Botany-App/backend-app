package usecases_garden

import (
	"context"
	"time"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CreateGardenUseCaseInputDTO struct {
	GardenName        string    `json:"garden_name"`
	UserID            string    `json:"user_id"`
	GardenDescription string    `json:"garden_description"`
	GardenLocation    string    `json:"garden_location"`
	TotalArea         float64   `json:"total_area"`
	CurrentingHeight  float64   `json:"currenting_heigth"`
	CurrentingWidth   float64   `json:"currenting_width"`
	PlantingDate      time.Time `json:"planting_date"`
	LastIrrigation    time.Time `json:"last_irrigation"`
	LastFertilization time.Time `json:"last_fertilization"`
	IrrigationWeek    int       `json:"irrigation_week"`
	SunExposure       int       `json:"sun_exposure"`
	FertilizationWeek int       `json:"fertilization_week"`
	CategoriesPlantId []string  `json:"categories_plant"`
	PlantsId          []string  `json:"plants_id"`
}

type CreateGardenUseCase struct {
	Repository entities.GardenRepository
}

func NewCreateGardenUseCase(repository entities.GardenRepository) *CreateGardenUseCase {
	return &CreateGardenUseCase{
		Repository: repository,
	}
}

func (uc *CreateGardenUseCase) Execute(ctx context.Context, input CreateGardenUseCaseInputDTO) (*entities.GardenOutputDTO, error) {
	garden, err := entities.NewGarden(
		input.GardenName,
		input.GardenLocation,
		input.GardenDescription,
		input.TotalArea,
		input.CurrentingHeight,
		input.CurrentingWidth,
		input.PlantingDate,
		input.LastIrrigation,
		input.LastFertilization,
		input.IrrigationWeek,
		input.SunExposure,
		input.FertilizationWeek,
		input.CategoriesPlantId,
		input.PlantsId,
	)
	if err != nil {
		return nil, err
	}

	id, err := uc.Repository.Create(ctx, garden)
	if err != nil {
		return nil, err
	}
	plantCreated, err := uc.Repository.FindByID(ctx, input.UserID, id)
	if err != nil {
		return nil, err
	}

	return plantCreated, nil
}
