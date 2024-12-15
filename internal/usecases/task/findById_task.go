package usecases_task

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByIdTaskUseCase struct {
	Repository entities.TaskRepository
}

type FindByIdTaskInputDTO struct {
	UserId string `json:"user_id"`
	TaskId string `json:"id"`
}

func NewFindByIdTaskUseCase(repository entities.TaskRepository) *FindByIdTaskUseCase {
	return &FindByIdTaskUseCase{Repository: repository}
}

func (u *FindByIdTaskUseCase) Execute(ctx context.Context, input FindByIdTaskInputDTO) (*entities.TaskOutputDTO, error) {

	task, err := u.Repository.FindByID(ctx, input.UserId, input.TaskId)
	if err != nil {
		return nil, err
	}

	return task, nil
}
