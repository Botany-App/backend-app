package usecases_categorytask

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type GetByIdCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

type GetByIdCategoryTaskDTO struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func NewGetByIdCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *GetByIdCategoryTaskUseCase {
	return &GetByIdCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}
func (uc *GetByIdCategoryTaskUseCase) Execute(ctx context.Context, input *GetByIdCategoryTaskDTO) (*entities.CategoryTask, error) {
	log.Println("-> GetByIdCategoryTaskUseCase - Execute")
	categoryTask, err := uc.CategoryTaskRepository.GetByID(ctx, input.UserID, input.ID)
	if err != nil {
		return nil, errors.New("categoryTask not found")
	}

	log.Println("<- GetByIdCategoryTaskUseCase - Execute")
	return categoryTask, nil
}
