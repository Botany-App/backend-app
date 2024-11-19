package usecases_categorytask

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllCategoryTaskInputDTO struct {
	UserID string `json:"user_id"`
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
	categories, err := uc.CategoryTaskRepository.FindAll(ctx, input.UserID)
	if err != nil {
		return nil, errors.New("error on Find all category task")
	}

	return categories, nil
}
