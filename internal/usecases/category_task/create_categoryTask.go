package usecases_categorytask

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CreateCategoryTaskDTO struct {
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

func (c *CreateCategoryTaskUseCase) Execute(ctx context.Context, input CreateCategoryTaskDTO) error {
	log.Print("--> Create Category Task Use Case")
	categoryTask, err := entities.NewCategoryTask(input.Name, input.Description, input.UserID)
	if err != nil {
		log.Print(err)
		return err
	}

	log.Println("--> Checking if category Task already exists")
	categoryTaskExists, _ := c.CreateCategoryTask.GetByName(ctx, categoryTask.Name, categoryTask.UserID)
	if categoryTaskExists != nil {
		log.Print(err)

		return errors.New("category Task already exists")
	}

	log.Println("--> Creating category Task")
	err = c.CreateCategoryTask.Create(ctx, categoryTask)
	if err != nil {
		log.Print(err)

		return errors.New("error while creating category Task")
	}

	log.Println("Category Task created successfully <-")
	return nil
}
