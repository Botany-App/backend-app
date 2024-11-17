package usecasesTasks

import (
	"context"
	"errors"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type DeleteTaskUseCase struct {
	TaskRepository entities.TaskRepository
}

type DeleteTaskDTO struct {
	TaskID string
	UserID string
}

func NewDeleteTaskUseCase(taskRepository entities.TaskRepository) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{
		TaskRepository: taskRepository,
	}
}
func (uc *DeleteTaskUseCase) Execute(ctx context.Context, dto *DeleteTaskDTO) error {
	task, err := uc.TaskRepository.FindByID(dto.TaskID, dto.UserID)
	if err != nil {
		return err
	}

	if task.UserID != dto.UserID {
		return errors.New("task does not belong to the user")
	}

	err = uc.TaskRepository.Delete(task.ID.String(), dto.UserID)
	if err != nil {
		return err
	}
	return nil
}
