package usecases_plant

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByCategoryNamePlantUseCase struct {
	PlantRepository entities.PlantRepository
}

type FindByCategoryNamePlantUseCaseInputDTO struct {
	UserId       string `json:"user_id"`
	CategoryName string `json:"category_name"`
}

func NewFindByCategoryNamePlantUseCase(repository entities.PlantRepository) *FindByCategoryNamePlantUseCase {
	return &FindByCategoryNamePlantUseCase{
		PlantRepository: repository,
	}
}

func (uc *FindByCategoryNamePlantUseCase) Execute(ctx context.Context, input FindByCategoryNamePlantUseCaseInputDTO) ([]*entities.PlantWithCategory, error) {
	log.Println("FindByCategoryNameUseCase - Execute")
	plants, err := uc.PlantRepository.FindByCategoryName(ctx, input.UserId, input.CategoryName)
	if err != nil {
		return nil, err
	}
	log.Println("FindByCategoryNameUseCase - Execute - plants: ", plants)
	return plants, nil
}
