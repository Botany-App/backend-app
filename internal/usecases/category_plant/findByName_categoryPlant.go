package usecases_categoryplant

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByNameCategoryPlantUseCase struct {
	FindByNameCategoryPlantRepository entities.CategoryPlantRepository
}

type FindByNameCategoryPlantInputDTO struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
	LIMIT  int    `json:"limit"`
	OFFSET int    `json:"offset"`
}

func NewCategoryPlantFindByNameUseCase(repository entities.CategoryPlantRepository) *FindByNameCategoryPlantUseCase {
	return &FindByNameCategoryPlantUseCase{FindByNameCategoryPlantRepository: repository}
}

func (uc *FindByNameCategoryPlantUseCase) Execute(ctx context.Context, input FindByNameCategoryPlantInputDTO) ([]*entities.CategoryPlant, error) {
	log.Println("FindByNameCategoryPlantUseCase - Execute")
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
	categories, err := uc.FindByNameCategoryPlantRepository.FindByName(ctx, userID, input.Name, input.LIMIT, input.OFFSET)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
