package usecases_garden

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllHistoryGardenUseCaseInputDTO struct {
	GardenId string `json:"garden_id"`
}

type FindAllHistoryGardenUseCase struct {
	GardenRepo entities.GardenRepository
}

func NewFindAllHistoryGardenUseCase(gardenRepository entities.GardenRepository) *FindAllHistoryGardenUseCase {
	return &FindAllHistoryGardenUseCase{
		GardenRepo: gardenRepository,
	}
}

func (u *FindAllHistoryGardenUseCase) Execute(ctx context.Context, input FindAllHistoryGardenUseCaseInputDTO) ([]*entities.HistoryGarden, error) {
	historyGardens, err := u.GardenRepo.FindAllHistoryByGardenID(ctx, input.GardenId)
	if err != nil {
		return nil, err
	}

	return historyGardens, nil
}
