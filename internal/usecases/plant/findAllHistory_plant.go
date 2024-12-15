package usecases_plant

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllHistoryPlantUseCase struct {
	PlantRepo entities.PlantRepository
}

type FindAllHistoryPlantUseCaseInputDTO struct {
	PlantId string `json:"plant_id"`
}

func NewFindAllHistoryPlantUseCase(plantRepository entities.PlantRepository) *FindAllHistoryPlantUseCase {
	return &FindAllHistoryPlantUseCase{
		PlantRepo: plantRepository,
	}
}

func (u *FindAllHistoryPlantUseCase) Execute(ctx context.Context, input FindAllHistoryPlantUseCaseInputDTO) ([]*entities.HistoryPlant, error) {
	historyPlants, err := u.PlantRepo.FindAllHistoryByPlantID(ctx, input.PlantId)
	if err != nil {
		return nil, err
	}

	return historyPlants, nil
}
