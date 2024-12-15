package usecases_task

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByStatusTaskUseCase struct {
	Repository entities.TaskRepository
}

type FindByStatusTaskInputDTO struct {
	UserId     string `json:"user_id"`
	TaskStatus string `json:"status"`
}

func NewFindByStatusTaskUseCase(repository entities.TaskRepository) *FindByStatusTaskUseCase {
	return &FindByStatusTaskUseCase{Repository: repository}
}

func (u *FindByStatusTaskUseCase) Execute(ctx context.Context, input FindByStatusTaskInputDTO) ([]*entities.TaskOutputDTO, error) {

	task, err := u.Repository.FindByStatus(ctx, input.UserId, input.TaskStatus)
	if err != nil {
		return nil, err
	}

	return task, nil
}
