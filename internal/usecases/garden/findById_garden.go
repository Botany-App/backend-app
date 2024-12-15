package usecases_garden

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByIdGardenUseCase struct {
	Repository entities.GardenRepository
}

type FindByIdGardenUseCaseInputDTO struct {
	ID     string `json:"id"`
	UserId string `json:"user_id"`
}

func NewFindByIdGardenUseCase(repository entities.GardenRepository) *FindByIdGardenUseCase {
	return &FindByIdGardenUseCase{
		Repository: repository,
	}
}

func (uc *FindByIdGardenUseCase) Execute(ctx context.Context, input FindByIdGardenUseCaseInputDTO) (*entities.GardenOutputDTO, error) {
	garden, err := uc.Repository.FindByID(ctx, input.UserId, input.ID)
	if err != nil {
		return nil, err
	}

	return garden, nil
}
