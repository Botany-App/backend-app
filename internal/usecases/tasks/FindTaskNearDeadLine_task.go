package usecasesTasks

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindTasksNearDeadlineDTO struct {
	UserID string
	Days   int
}

type FindTaskNearDeadLineTaskUseCase struct {
	TaskRepository entities.TaskRepository
}

func NewFindTaskNearDeadLineTaskUseCase(taskRepository entities.TaskRepository) *FindTaskNearDeadLineTaskUseCase {
	return &FindTaskNearDeadLineTaskUseCase{
		TaskRepository: taskRepository,
	}
}

func (u *FindTaskNearDeadLineTaskUseCase) Execute(ctx context.Context, input *FindTasksNearDeadlineDTO) ([]entities.Task, error) {
	tasks, err := u.TaskRepository.FindTasksNearDeadline(input.UserID, input.Days)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
