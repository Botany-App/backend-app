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
}

func NewFindByIdCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *FindByIdCategoryTaskUseCase {
	return &FindByIdCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}
func (uc *FindByIdCategoryTaskUseCase) Execute(ctx context.Context, input FindByIdCategoryTaskInputDTO) (*entities.CategoryTask, error) {
	log.Println("FindByIdCategoryTaskUseCase - Execute")
	categoryTask, err := uc.CategoryTaskRepository.FindByID(ctx, input.UserID, input.ID)
	if err != nil {
		return nil, errors.New("categoryTask not found")
	}

	return categoryTask, nil
}
