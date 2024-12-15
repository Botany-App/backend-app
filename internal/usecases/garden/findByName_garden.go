package usecases_garden

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByNameGardenUseCase struct {
	Repository entities.GardenRepository
}

type FindByNameGardenUseCaseInputDTO struct {
	Name   string `json:"name"`
	UserId string `json:"userId"`
}

func NewFindByNameGardenUseCase(repository entities.GardenRepository) *FindByNameGardenUseCase {
	return &FindByNameGardenUseCase{Repository: repository}
}

func (u *FindByNameGardenUseCase) Execute(ctx context.Context, input FindByNameGardenUseCaseInputDTO) ([]*entities.GardenOutputDTO, error) {
	garden, err := u.Repository.FindByName(ctx, input.UserId, input.Name)
	if err != nil {
		return nil, err
	}

	return garden, nil
}
