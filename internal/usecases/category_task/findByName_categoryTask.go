package usecases_categorytask

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByNameCategoryTaskInputDTO struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	LIMIT  int    `json:"limit"`
	OFFSET int    `json:"offset"`
}

type FindByNameCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

func NewFindByNameCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *FindByNameCategoryTaskUseCase {
	return &FindByNameCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}

func (u *FindByNameCategoryTaskUseCase) Execute(ctx context.Context, input FindByNameCategoryTaskInputDTO) ([]*entities.CategoryTask, error) {
	log.Println("FindByNameCategoryTaskUseCase - Execute")
	if input.LIMIT <= 0 {
		input.LIMIT = 10 // Default limit
	}
	if input.OFFSET < 0 {
		input.OFFSET = 0
	}
	categoryTask, err := u.CategoryTaskRepository.FindByName(ctx, input.UserID, input.Name, input.LIMIT, input.OFFSET)
	if err != nil {
		return nil, errors.New("error Findting category task by name")
	}

	return categoryTask, nil
}
