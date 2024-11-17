package usecasesTasks

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllByNameDTO struct {
	UserID string
	Name   string
}

type FindAllByNameTaskUseCase struct {
	TaskRepository entities.TaskRepository
}

func NewFindAllByNameTaskUseCase(taskRepository entities.TaskRepository) *FindAllByNameTaskUseCase {
	return &FindAllByNameTaskUseCase{
		TaskRepository: taskRepository,
	}
}

func (u *FindAllByNameTaskUseCase) Execute(ctx context.Context, dto *FindAllByNameDTO) ([]entities.Task, error) {
	log.Println("--> Find all tasks by name")
	tasks, err := u.TaskRepository.FindAllByName(dto.UserID, dto.Name)
	if err != nil {
		return nil, err
	}
	log.Println("<-- Find all tasks by name")
	return tasks, nil
}
