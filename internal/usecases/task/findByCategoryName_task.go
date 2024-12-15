package usecases_task

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByCategoryNameTaskUseCase struct {
	Repository entities.TaskRepository
}

type FindByCategoryNameTaskInputDTO struct {
	UserId           string `json:"user_id"`
	TaskCategoryName string `json:"category_name"`
}

func NewFindByCategoryNameTaskUseCase(repository entities.TaskRepository) *FindByCategoryNameTaskUseCase {
	return &FindByCategoryNameTaskUseCase{Repository: repository}
}

func (u *FindByCategoryNameTaskUseCase) Execute(ctx context.Context, input FindByCategoryNameTaskInputDTO) ([]*entities.TaskOutputDTO, error) {

	task, err := u.Repository.FindByCategoryName(ctx, input.UserId, input.TaskCategoryName)
	if err != nil {
		return nil, err
	}

	return task, nil
}
