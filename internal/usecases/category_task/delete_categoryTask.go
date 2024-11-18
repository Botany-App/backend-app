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

type DeleteCategoryTaskDTO struct {
	UserID string `json:"user_id"`
	ID     string `json:"id"`
}

func NewDeleteCategoryTaskUseCase(repository entities.CategoryTaskRepository) *DeleteCategoryTaskUseCase {
	return &DeleteCategoryTaskUseCase{CategoryTaskRepository: repository}
}

func (useCase *DeleteCategoryTaskUseCase) Execute(ctx context.Context, input DeleteCategoryTaskDTO) error {
	log.Println("--> DELETE CATEGORY TASK")
	err := useCase.CategoryTaskRepository.Delete(ctx, input.UserID, input.ID)
	if err != nil {
		log.Println("Erro ao deletar categoria de tarefa")
		return errors.New("erro ao deletar categoria de tarefa")
	}

	log.Println("<-Categoria de tarefa deletada com sucesso")
	return nil
}
