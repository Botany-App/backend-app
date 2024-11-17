package usecasesTasks

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindTasksFarFromDeadlineDTO struct {
	UserID string
	Days   int
}

type FindTasksFarFromDeadlineUseCase struct {
	TaskRepository entities.TaskRepository
}

func NewFindTasksFarFromDeadlineUseCase(taskRepository entities.TaskRepository) *FindTasksFarFromDeadlineUseCase {
	return &FindTasksFarFromDeadlineUseCase{
		TaskRepository: taskRepository,
	}
}

func (u *FindTasksFarFromDeadlineUseCase) Execute(ctx context.Context, input *FindTasksFarFromDeadlineDTO) ([]entities.Task, error) {
	log.Println("FindTasksFarFromDeadlineUseCase")
	tasks, err := u.TaskRepository.FindTasksFarFromDeadline(input.UserID, input.Days)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("FindTasksFarFromDeadlineUseCase")
	return tasks, nil
}
