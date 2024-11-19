package usecases_categorytask

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type DeleteCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

type DeleteCategoryTaskInputDTO struct {
	UserID string `json:"user_id"`
	ID     string `json:"id"`
}

func NewDeleteCategoryTaskUseCase(repository entities.CategoryTaskRepository) *DeleteCategoryTaskUseCase {
	return &DeleteCategoryTaskUseCase{CategoryTaskRepository: repository}
}

func (useCase *DeleteCategoryTaskUseCase) Execute(ctx context.Context, input DeleteCategoryTaskInputDTO) error {
	log.Println("DeleteCategoryTaskUseCase - Execute")
	err := useCase.CategoryTaskRepository.Delete(ctx, input.UserID, input.ID)
	if err != nil {
		return errors.New("erro ao deletar categoria de tarefa")
	}

	return nil
}
