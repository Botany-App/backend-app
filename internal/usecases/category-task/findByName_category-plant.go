package usecases_categoryTask

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type FindByNameCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

type FindByNameCategoryTaskInputDTO struct {
	Name   string `json:"name"`
	UserId string `json:"userId"`
}

func NewFindByNameCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *FindByNameCategoryTaskUseCase {
	return &FindByNameCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}
func (uc *FindByNameCategoryTaskUseCase) Execute(ctx context.Context, input FindByNameCategoryTaskInputDTO) ([]*entities.CategoryTask, error) {
	log.Println("FindByNameCategoryTaskUseCase - Execute")
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	categoriesTask, err := uc.CategoryTaskRepository.FindByName(ctx, input.UserId, input.Name)
	if err != nil {
		return nil, errors.Wrap(err, "error finding all category Task")
	}
	return categoriesTask, nil
}
