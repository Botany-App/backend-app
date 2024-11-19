package usecases_categorytask

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllCategoryTaskInputDTO struct {
	UserID string `json:"user_id"`
	LIMIT  int    `json:"limit"`
	OFFSET int    `json:"offset"`
}

type FindAllCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

func NewFindAllCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *FindAllCategoryTaskUseCase {
	return &FindAllCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}

func (uc *FindAllCategoryTaskUseCase) Execute(ctx context.Context, input FindAllCategoryTaskInputDTO) ([]*entities.CategoryTask, error) {
	log.Print("FindAllCategoryTaskUseCase - Execute")
	if input.LIMIT <= 0 {
		input.LIMIT = 10 // Default limit
	}
	if input.OFFSET < 0 {
		input.OFFSET = 0
	}
	categories, err := uc.CategoryTaskRepository.FindAll(ctx, input.UserID, input.LIMIT, input.OFFSET)
	if err != nil {
		return nil, errors.New("error on Find all category task")
	}

	return categories, nil
}
