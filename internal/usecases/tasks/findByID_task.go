package usecasesTasks

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByIDTaskDTO struct {
	UserID string `json:"user_id"`
	ID     string `json:"id"`
}

type FindByIDTaskUseCase struct {
	TaskRepository entities.TaskRepository
}

func NewFindByIDTaskUseCase(taskRepository entities.TaskRepository) *FindByIDTaskUseCase {
	return &FindByIDTaskUseCase{TaskRepository: taskRepository}
}

func (uc *FindByIDTaskUseCase) Execute(ctx context.Context, input *FindByIDTaskDTO) (*entities.Task, error) {
	log.Println("--> FIND TASK BY ID")
	task, err := uc.TaskRepository.FindByID(input.UserID, input.ID)
	if err != nil {
		return nil, err
	}
	log.Println("<- Task found successfully")
	return task, nil
}
