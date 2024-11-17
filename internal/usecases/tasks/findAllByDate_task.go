package usecasesTasks

import (
	"context"
	"log"
	"time"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllByDateDTO struct {
	UserID string
	Date   string
}

type FindAllByDateTaskUseCase struct {
	TaskRepository entities.TaskRepository
}

func NewFindAllByDateTaskUseCase(taskRepository entities.TaskRepository) *FindAllByDateTaskUseCase {
	return &FindAllByDateTaskUseCase{
		TaskRepository: taskRepository,
	}
}

func (u *FindAllByDateTaskUseCase) Execute(ctx context.Context, input *FindAllByDateDTO) ([]entities.Task, error) {
	log.Println("--> Find all tasks by date use case")
	date, err := time.Parse(input.Date, "2006-01-02")
	if err != nil {
		return nil, err
	}
	tasks, err := u.TaskRepository.FindAllByDate(input.UserID, date)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
