package usecases_categoryTask

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
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

func NewDeleteCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *DeleteCategoryTaskUseCase {
	return &DeleteCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}
func (uc *DeleteCategoryTaskUseCase) Execute(ctx context.Context, input DeleteCategoryTaskInputDTO) error {
	log.Println("DeleteCategoryTaskUseCase - Execute")

	err := uc.CategoryTaskRepository.Delete(ctx, input.UserId, input.Id)
	if err != nil {
		return errors.Wrap(err, "error delete  category Task")
	}
	return nil
}
