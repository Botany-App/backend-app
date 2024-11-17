package usecasesTasks

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllTaskDTO struct {
	UserID string `json:"user_id"`
}

type FindAllTaskUseCase struct {
	TaskRepository entities.TaskRepository
}

func NewFindAllTaskUseCase(taskRepository entities.TaskRepository) *FindAllTaskUseCase {
	return &FindAllTaskUseCase{
		TaskRepository: taskRepository,
	}
}

func (uc *FindAllTaskUseCase) Execute(ctx context.Context, input *FindAllTaskDTO) ([]entities.Task, error) {
	log.Println("--> Find All Task")
	tasks, err := uc.TaskRepository.FindAll(input.UserID)
	if err != nil {
		return nil, err
	}
	log.Print("<- Found ", len(tasks), " tasks")
	return tasks, nil
}
