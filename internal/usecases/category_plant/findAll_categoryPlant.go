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
	LIMIT  int    `json:"limit"`
	OFFSET int    `json:"offset"`
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
	if input.LIMIT <= 0 {
		input.LIMIT = 10 // Default limit
	}
	if input.OFFSET < 0 {
		input.OFFSET = 0
	}
	categoryPlants, err := uc.FindAllCategoryPlantRepository.FindAll(ctx, userID, input.LIMIT, input.OFFSET)
	if err != nil {
		return nil, err
	}

	return categoryPlants, nil
}
