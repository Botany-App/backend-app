package usecases_categorytask

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type GetByNameCategoryTaskDTO struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type GetByNameCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

func NewGetByNameCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *GetByNameCategoryTaskUseCase {
	return &GetByNameCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}

func (u *GetByNameCategoryTaskUseCase) Execute(ctx context.Context, input *GetByNameCategoryTaskDTO) ([]entities.CategoryTask, error) {
	log.Println("-> Get by name category task use case")
	categoryTask, err := u.CategoryTaskRepository.GetByName(ctx, input.UserID, input.Name)
	if err != nil {
		return nil, errors.New("error getting category task by name")
	}

	log.Println("<- Get by name category task use case")
	return categoryTask, nil
}
