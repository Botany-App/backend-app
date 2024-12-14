package usecases_categoryTask

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type FindAllCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

func NewFindAllCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *FindAllCategoryTaskUseCase {
	return &FindAllCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}
func (uc *FindAllCategoryTaskUseCase) Execute(ctx context.Context, input string) ([]*entities.CategoryTask, error) {
	log.Println("FindAllCategoryTaskUseCase - Execute")
	categoriesTask, err := uc.CategoryTaskRepository.FindAll(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "error finding by all categories Tasks")
	}
	return categoriesTask, nil
}
