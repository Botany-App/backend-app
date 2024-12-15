package usecases_task

import (
	"context"
	"time"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CreateTaskUseCase struct {
	Repository entities.TaskRepository
}

type CreateTaskUseCaseInputDTO struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	TaskDate     time.Time `json:"task_date"`
	UrgencyLevel int       `json:"urgency_level"`
	TaskStatus   string    `json:"task_status"`
	UserId       string    `json:"user_id"`
	CategoriesId []string  `json:"categories_id"`
	GardensId    []string  `json:"gardens_id"`
	PlantsId     []string  `json:"plants_id"`
}

func NewCreateTaskUseCase(repository entities.TaskRepository) *CreateTaskUseCase {
	return &CreateTaskUseCase{
		Repository: repository,
	}
}

func (uc *CreateTaskUseCase) Execute(ctx context.Context, input CreateTaskUseCaseInputDTO) (*entities.TaskOutputDTO, error) {
	task, err := entities.NewTask(
		input.Name,
		input.Description,
		input.TaskDate,
		input.UrgencyLevel,
		input.TaskStatus,
		input.UserId,
		input.CategoriesId,
		input.GardensId,
		input.PlantsId,
	)
	if err != nil {
		return nil, err
	}

	id, err := uc.Repository.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	taskOutput, err := uc.Repository.FindByID(ctx, input.UserId, id)
	if err != nil {
		return nil, err
	}

	return taskOutput, nil
}
