package usecases_categoryplant

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllCategoryPlantUseCase struct {
	FindAllCategoryPlantRepository entities.CategoryPlantRepository
}

type FindAllCategoryPlantDTO struct {
	UserID string `json:"user_id"`
}

func NewFindAllCategoryPlantUseCase(categoryPlantRepository entities.CategoryPlantRepository) *FindAllCategoryPlantUseCase {
	return &FindAllCategoryPlantUseCase{FindAllCategoryPlantRepository: categoryPlantRepository}
}

func (uc *FindAllCategoryPlantUseCase) Execute(ctx context.Context, input FindAllCategoryPlantDTO) ([]*entities.CategoryPlant, error) {
	log.Println("FindAllCategoryPlantUseCase: Executing...")
	defer log.Println("FindAllCategoryPlantUseCase: End")

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, err
	}

	categoryPlants, err := uc.FindAllCategoryPlantRepository.FindAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	return categoryPlants, nil
}
