package usecases_task

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByNameTaskUseCase struct {
	Repository entities.TaskRepository
}

type FindByNameTaskInputDTO struct {
	UserId   string `json:"user_id"`
	TaskName string `json:"name"`
}

func NewFindByNameTaskUseCase(repository entities.TaskRepository) *FindByNameTaskUseCase {
	return &FindByNameTaskUseCase{Repository: repository}
}

func (u *FindByNameTaskUseCase) Execute(ctx context.Context, input FindByNameTaskInputDTO) ([]*entities.TaskOutputDTO, error) {

	task, err := u.Repository.FindByName(ctx, input.UserId, input.TaskName)
	if err != nil {
		return nil, err
	}

	return task, nil
}
