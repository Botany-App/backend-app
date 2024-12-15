package usecases_task

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllTaskUseCase struct {
	Repository entities.TaskRepository
}

type FindAllTaskInputDTO struct {
	UserId string `json:"user_id"`
}

func NewFindAllTaskUseCase(repository entities.TaskRepository) *FindAllTaskUseCase {
	return &FindAllTaskUseCase{Repository: repository}
}

func (u *FindAllTaskUseCase) Execute(ctx context.Context, input FindAllTaskInputDTO) ([]*entities.TaskOutputDTO, error) {

	task, err := u.Repository.FindAll(ctx, input.UserId)
	if err != nil {
		return nil, err
	}

	return task, nil
}
