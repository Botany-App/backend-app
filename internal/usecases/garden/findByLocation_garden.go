package usecases_garden

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByLocationGardenUseCase struct {
	Repository entities.GardenRepository
}

type FindByLocationGardenUseCaseInputDTO struct {
	Location string `json:"location"`
	UserId   string `json:"userId"`
}

func NewFindByLocationGardenUseCase(repository entities.GardenRepository) *FindByLocationGardenUseCase {
	return &FindByLocationGardenUseCase{Repository: repository}
}

func (u *FindByLocationGardenUseCase) Execute(ctx context.Context, input FindByLocationGardenUseCaseInputDTO) ([]*entities.GardenOutputDTO, error) {
	garden, err := u.Repository.FindByLocation(ctx, input.UserId, input.Location)
	if err != nil {
		return nil, err
	}

	return garden, nil
}
