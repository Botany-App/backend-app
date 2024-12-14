package usecases_categoryTask

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type CreateCategoryTaskInputDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCategoryTaskUseCase struct {
	categoryTaskRepository entities.CategoryTaskRepository
}

func NewCreateCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *CreateCategoryTaskUseCase {
	return &CreateCategoryTaskUseCase{
		categoryTaskRepository: categoryTaskRepository,
	}
}

func (uc *CreateCategoryTaskUseCase) Execute(ctx context.Context, input CreateCategoryTaskInputDTO, userId string) (*entities.CategoryTask, error) {
	log.Println("CreateCategoryTaskUseCase - Execute")
	newCategoryTask, err := entities.NewCreateCategoryTask(input.Name, input.Description, userId)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new category Task")
	}

	existingCategoryTask, err := uc.categoryTaskRepository.FindByName(ctx, userId, input.Name)
	if err != nil {
		return nil, errors.Wrap(err, "error finding category Task by name")
	}

	for _, category := range existingCategoryTask {
		if category.Name == input.Name {
			log.Println(category.Name)
			return nil, errors.New("category Task already exists")
		}
	}

	id, err := uc.categoryTaskRepository.Create(ctx, newCategoryTask)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new category Task")
	}

	categoryCreated, err := uc.categoryTaskRepository.FindById(ctx, userId, id)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new category Task")
	}
	return categoryCreated, nil

}
