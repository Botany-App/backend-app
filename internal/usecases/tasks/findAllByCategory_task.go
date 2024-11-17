package usecasesTasks

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllByCategoryDTO struct {
	UserID     string
	CategoryID string
}

type FindAllByCategoryTaskUseCase struct {
	TaskRepository entities.TaskRepository
}

func NewFindAllByCategoryTaskUseCase(taskRepository entities.TaskRepository) *FindAllByCategoryTaskUseCase {
	return &FindAllByCategoryTaskUseCase{
		TaskRepository: taskRepository,
	}
}

func (u *FindAllByCategoryTaskUseCase) Execute(ctx context.Context, dto *FindAllByCategoryDTO) ([]entities.Task, error) {
	tasks, err := u.TaskRepository.FindAllByCategory(dto.UserID, dto.CategoryID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
