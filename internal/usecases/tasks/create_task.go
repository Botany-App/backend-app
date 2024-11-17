package usecasesTasks

import (
	"context"
	"time"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CreateTaskDTO struct {
	Name           string    `json:"name"`
	Descrition     string    `json:"description"`
	TaskDate       time.Time `json:"task_date"`
	UserID         string    `json:"user_id"`
	GardenPlantaId []string  `json:"garden_planta_id"`
	CategoriesId   []string  `json:"categories_id"`
}

type CreateTaskUseCase struct {
	TaskRepository entities.TaskRepository
}

func NewCreateTaskUseCase(taskRepository entities.TaskRepository) *CreateTaskUseCase {
	return &CreateTaskUseCase{
		TaskRepository: taskRepository,
	}
}

func (uc *CreateTaskUseCase) Execute(ctx context.Context, input *CreateTaskDTO) error {
	task, err := entities.NewTask(
		input.Name,
		input.Descrition,
		input.TaskDate,
		input.UserID,
		input.GardenPlantaId,
		input.CategoriesId,
		entities.Pending,
	)
	if err != nil {
		return err
	}
	err = uc.TaskRepository.Create(task)
	if err != nil {
		return err
	}

	if task.CategoriesId != nil {
		for _, categoryID := range task.CategoriesId {
			err = uc.TaskRepository.AddCategory(task.ID.String(), categoryID)
			if err != nil {
				return err
			}
		}
	}

	if task.GardenPlantaId != nil {
		for _, gardenPlantaID := range task.GardenPlantaId {
			err = uc.TaskRepository.AddGardenPlanta(task.ID.String(), gardenPlantaID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
