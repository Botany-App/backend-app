package usecases_categorytask

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CreateCategoryTaskInputDTO struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCategoryTaskUseCase struct {
	CreateCategoryTask entities.CategoryTaskRepository
}

func NewCreateCategoryTaskUseCase(createCategoryTask entities.CategoryTaskRepository) *CreateCategoryTaskUseCase {
	return &CreateCategoryTaskUseCase{CreateCategoryTask: createCategoryTask}
}

func (c *CreateCategoryTaskUseCase) Execute(ctx context.Context, input CreateCategoryTaskInputDTO) error {
	log.Print("CreateCategoryTaskUseCase - Execute")
	categoryTask, err := entities.NewCategoryTask(input.Name, input.Description, input.UserID)
	if err != nil {
		return errors.New("error while creating category Task")
	}
	categoryTaskExists, _ := c.CreateCategoryTask.FindByName(ctx, categoryTask.Name, categoryTask.UserID)
	if categoryTaskExists != nil {
		return errors.New("category Task already exists")
	}
	err = c.CreateCategoryTask.Create(ctx, categoryTask)
	if err != nil {
		return errors.New("error while creating category Task")
	}
	return nil
}
