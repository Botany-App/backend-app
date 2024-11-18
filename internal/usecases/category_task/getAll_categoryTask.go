package usecases_categorytask

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type GetAllCategoryTaskDTO struct {
	UserID string `json:"user_id"`
}

type GetAllCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

func NewGetAllCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *GetAllCategoryTaskUseCase {
	return &GetAllCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}

func (uc *GetAllCategoryTaskUseCase) Execute(ctx context.Context, input *GetAllCategoryTaskDTO) ([]entities.CategoryTask, error) {
	log.Print("--> Get All Category Task Use Case")
	categories, err := uc.CategoryTaskRepository.GetAll(ctx, input.UserID)
	if err != nil {
		return nil, errors.New("error on get all category task")
	}

	log.Print("<-- Get All Category Task Use Case")
	return categories, nil
}
