package usecasesTasks

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllByStatusDTO struct {
	UserID string
	Status entities.TaskStatusEnum
}

type FindAllByStatusTaskUseCase struct {
	TaskRepository entities.TaskRepository
}

func NewFindAllByStatusTaskUseCase(taskRepository entities.TaskRepository) *FindAllByStatusTaskUseCase {
	return &FindAllByStatusTaskUseCase{
		TaskRepository: taskRepository,
	}
}

func (u *FindAllByStatusTaskUseCase) Execute(ctx context.Context, input *FindAllByStatusDTO) ([]entities.Task, error) {
	log.Println("--> Find all by status task use case")
	tasks, err := u.TaskRepository.FindAllByStatus(input.UserID, input.Status)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
