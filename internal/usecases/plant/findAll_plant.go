package usecases_plant

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllPlantUseCase struct {
	PlantRepository entities.PlantRepository
}

type FindAllPlantUseCaseInputDTO struct {
	UserId string `json:"userId"`
}

func NewFindAllPlantUseCase(repository entities.PlantRepository) *FindAllPlantUseCase {
	return &FindAllPlantUseCase{
		PlantRepository: repository,
	}
}

func (uc *FindAllPlantUseCase) Execute(ctx context.Context, input FindAllPlantUseCaseInputDTO) ([]*entities.PlantWithCategory, error) {
	log.Println("FindAllPlantUseCase - Execute")
	plants, err := uc.PlantRepository.FindAll(ctx, input.UserId)
	if err != nil {
		return nil, err
	}
	return plants, nil
}
