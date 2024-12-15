package usecases_task

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type DeleteTaskUseCase struct {
	Repository entities.TaskRepository
}

type DeleteTaskInputDTO struct {
	UserId string `json:"user_id"`
	TaskId string `json:"id"`
}

func NewDeleteTaskUseCase(repository entities.TaskRepository) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{Repository: repository}
}

func (u *DeleteTaskUseCase) Execute(ctx context.Context, input DeleteTaskInputDTO) error {

	err := u.Repository.Delete(ctx, input.UserId, input.TaskId)
	if err != nil {
		return err
	}

	return nil
}
