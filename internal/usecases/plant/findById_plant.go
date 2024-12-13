package usecases_plant

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByIdPlantUseCase struct {
	PlantRepository entities.PlantRepository
}

type FindByIdPlantUseCaseInputDTO struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func NewFindByIdPlantUseCase(repository entities.PlantRepository) *FindByIdPlantUseCase {
	return &FindByIdPlantUseCase{
		PlantRepository: repository,
	}
}

func (uc *FindByIdPlantUseCase) Execute(ctx context.Context, input FindByIdPlantUseCaseInputDTO) (*entities.PlantWithCategory, error) {
	log.Println("FindByIdPlant - Execute")
	plant, err := uc.PlantRepository.FindByID(ctx, input.UserID, input.ID)
	if err != nil {
		return nil, err
	}
	return plant, nil
}
