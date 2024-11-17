package usecases_categorytask

import (
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

func (c *CreateCategoryTaskUseCase) Execute(dto CreateCategoryTaskDTO) error {
	categoryTask, err := entities.NewCategoryTask(dto.Name, dto.Description)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Print(categoryTask)
	err = c.CreateCategoryTask.Create(dto.UserID, categoryTask)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
