package usecases_categorytask

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByIdCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

type FindByIdCategoryTaskInputDTO struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	LIMIT  int    `json:"limit"`
	OFFSET int    `json:"offset"`
}

func NewFindByIdCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *FindByIdCategoryTaskUseCase {
	return &FindByIdCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}
func (uc *FindByIdCategoryTaskUseCase) Execute(ctx context.Context, input FindByIdCategoryTaskInputDTO) (*entities.CategoryTask, error) {
	log.Println("FindByIdCategoryTaskUseCase - Execute")
	if input.LIMIT <= 0 {
		input.LIMIT = 10 // Default limit
	}
	if input.OFFSET < 0 {
		input.OFFSET = 0
	}
	categoryTask, err := uc.CategoryTaskRepository.FindByID(ctx, input.UserID, input.ID, input.LIMIT, input.OFFSET)
	if err != nil {
		return nil, errors.New("categoryTask not found")
	}

	return categoryTask, nil
}
