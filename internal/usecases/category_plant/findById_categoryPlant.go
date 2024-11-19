package usecases_categoryplant

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByIdCategoryPlantUseCase struct {
	FindByIDCategoryPlantRepository entities.CategoryPlantRepository
}

type FindByIdCategoryPlantInputDTO struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	LIMIT  int    `json:"limit"`
	OFFSET int    `json:"offset"`
}

func NewFindByIdCategoryPlantUseCase(repository entities.CategoryPlantRepository) *FindByIdCategoryPlantUseCase {
	return &FindByIdCategoryPlantUseCase{FindByIDCategoryPlantRepository: repository}
}

func (uc *FindByIdCategoryPlantUseCase) Execute(ctx context.Context, input FindByIdCategoryPlantInputDTO) (*entities.CategoryPlant, error) {
	log.Println("FindByIdCategoryPlantUseCase - Execute")
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return nil, errors.New("ID is not a valid UUID")
	}
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, errors.New("UserID is not a valid UUID")
	}
	if input.LIMIT <= 0 {
		input.LIMIT = 10 // Default limit
	}
	if input.OFFSET < 0 {
		input.OFFSET = 0
	}
	categoryPlant, err := uc.FindByIDCategoryPlantRepository.FindByID(ctx, userID, id, input.LIMIT, input.OFFSET)
	if err != nil {
		return nil, errors.New("CategoryPlant not found")
	}

	return categoryPlant, nil
}
