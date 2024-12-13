package usecases_plant

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByNamePlantUseCase struct {
	PlantRepository entities.PlantRepository
}

type FindByNamePlantUseCaseInputDTO struct {
	Name   string `json:"name"`
	UserId string `json:"user_id"`
}

func NewFindByNameCategoryPlantUseCase(repository entities.PlantRepository) *FindByNamePlantUseCase {
	return &FindByNamePlantUseCase{
		PlantRepository: repository,
	}
}

func (uc *FindByNamePlantUseCase) Execute(ctx context.Context, input FindByNamePlantUseCaseInputDTO) ([]*entities.PlantWithCategory, error) {
	log.Println("FindByNamePlant - Execute")
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	plants, err := uc.PlantRepository.FindByName(ctx, input.UserId, input.Name)
	if err != nil {
		return nil, err
	}
	return plants, nil
}
