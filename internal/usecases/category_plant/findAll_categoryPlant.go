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

type FindAllCategoryPlantInputDTO struct {
	UserID string `json:"user_id"`
}

func NewFindAllCategoryPlantUseCase(categoryPlantRepository entities.CategoryPlantRepository) *FindAllCategoryPlantUseCase {
	return &FindAllCategoryPlantUseCase{FindAllCategoryPlantRepository: categoryPlantRepository}
}

func (uc *FindAllCategoryPlantUseCase) Execute(ctx context.Context, input FindAllCategoryPlantInputDTO) ([]*entities.CategoryPlant, error) {
	log.Println("FindAllCategoryPlantUseCase - Execute")

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
