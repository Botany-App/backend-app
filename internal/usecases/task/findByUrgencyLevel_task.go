package usecases_task

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByUrgencyLevelTaskUseCase struct {
	Repository entities.TaskRepository
}

type FindByUrgencyLevelTaskInputDTO struct {
	UserId           string `json:"user_id"`
	TaskUrgencyLevel int    `json:"UrgencyLevel"`
}

func NewFindByUrgencyLevelTaskUseCase(repository entities.TaskRepository) *FindByUrgencyLevelTaskUseCase {
	return &FindByUrgencyLevelTaskUseCase{Repository: repository}
}

func (u *FindByUrgencyLevelTaskUseCase) Execute(ctx context.Context, input FindByUrgencyLevelTaskInputDTO) ([]*entities.TaskOutputDTO, error) {

	task, err := u.Repository.FindByUrgencyLevel(ctx, input.UserId, input.TaskUrgencyLevel)
	if err != nil {
		return nil, err
	}

	return task, nil
}
