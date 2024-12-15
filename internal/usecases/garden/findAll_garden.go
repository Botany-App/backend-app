package usecases_garden

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllGardenUseCaseInputDTO struct {
	UserId string `json:"user_id"`
}

type FindAllGardenUseCase struct {
	FindAllGardenRepository entities.GardenRepository
}

func NewFindAllGardenUseCase(findAllGardenRepository entities.GardenRepository) *FindAllGardenUseCase {
	return &FindAllGardenUseCase{FindAllGardenRepository: findAllGardenRepository}
}

func (useCase *FindAllGardenUseCase) Execute(ctx context.Context, input FindAllGardenUseCaseInputDTO) ([]*entities.GardenOutputDTO, error) {
	gardens, err := useCase.FindAllGardenRepository.FindAll(ctx, input.UserId)
	if err != nil {
		return nil, err
	}

	return gardens, nil
}
