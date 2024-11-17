package usecasesTasks

import (
	"context"
	"time"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UpdateTaskDTO struct {
	UserID         string                  `json:"user_id"`
	Name           string                  `json:"name"`
	Description    string                  `json:"description"`
	TaskDate       time.Time               `json:"task_date"`
	GardenPlantaId []string                `json:"garden_planta_id"`
	CategoriesId   []string                `json:"categories_id"`
	TaskStatus     entities.TaskStatusEnum `json:"task_status"`
}

type UpdateTaskUseCase struct {
	TaskRepository entities.TaskRepository
}

func NewUpdateTaskUseCase(taskRepo entities.TaskRepository) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{
		TaskRepository: taskRepo,
	}
}

func (uc *UpdateTaskUseCase) Execute(ctx context.Context, input *UpdateTaskDTO) error {
	task, err := uc.TaskRepository.FindByID(input.UserID, input.Name)
	if err != nil {
		return err
	}

	if input.Description != "" {
		task.Description = input.Description
	}

	if !input.TaskDate.IsZero() {
		task.TaskDate = input.TaskDate
	}

	if len(input.GardenPlantaId) > 0 {
		task.GardenPlantaId = make([]string, len(input.GardenPlantaId))
		copy(task.GardenPlantaId, input.GardenPlantaId)
	}

	if len(input.CategoriesId) > 0 {
		task.CategoriesId = make([]string, len(input.CategoriesId))
		copy(task.CategoriesId, input.CategoriesId)
	}

	if input.TaskStatus != "" {
		task.TaskStatus = input.TaskStatus
	}

	if err := uc.TaskRepository.Update(task); err != nil {
		return err
	}

	for _, categoryID := range task.CategoriesId {
		if err := uc.TaskRepository.UpdateTaskCategory(task.ID.String(), categoryID); err != nil {
			return err
		}
	}

	for _, gardenPlantaID := range task.GardenPlantaId {
		if err := uc.TaskRepository.UpdateTaskGardenPlanta(task.ID.String(), gardenPlantaID); err != nil {
			return err
		}
	}

	return nil
}
