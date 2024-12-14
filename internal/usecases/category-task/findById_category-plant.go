package usecases_categoryTask

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type FindByIdCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

type FindByIdCategoryTaskInputDTO struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

func NewFindByIdCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *FindByIdCategoryTaskUseCase {
	return &FindByIdCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}
func (uc *FindByIdCategoryTaskUseCase) Execute(ctx context.Context, input FindByIdCategoryTaskInputDTO) (*entities.CategoryTask, error) {
	log.Println("FindByIdCategoryTaskUseCase - Execute")
	categoryTask, err := uc.CategoryTaskRepository.FindById(ctx, input.UserId, input.Id)
	if err != nil {
		return nil, errors.Wrap(err, "error finding by id category Task")
	}
	return categoryTask, nil
}
