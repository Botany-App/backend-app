package usecases_categoryTask

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type UpdateCategoryTaskInputDTO struct {
	UserId      string `json:"user_id"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCategoryTaskUseCase struct {
	categoryTaskRepository entities.CategoryTaskRepository
}

func NewUpdateCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *UpdateCategoryTaskUseCase {
	return &UpdateCategoryTaskUseCase{
		categoryTaskRepository: categoryTaskRepository,
	}
}

func (uc *UpdateCategoryTaskUseCase) Execute(ctx context.Context, input UpdateCategoryTaskInputDTO) (*entities.CategoryTask, error) {
	log.Println("UpdateCategoryTaskUseCase - Execute")

	category, err := uc.categoryTaskRepository.FindById(ctx, input.UserId, input.Id)
	if err != nil {
		return nil, errors.Wrap(err, "error while searching category")
	}
	if category == nil {
		return nil, errors.New("category not found")
	}
	existingCategoryTask, err := uc.categoryTaskRepository.FindByName(ctx, input.UserId, input.Name)
	if err != nil {
		return nil, errors.Wrap(err, "error finding category Task by name")
	}

	for _, category := range existingCategoryTask {
		if category.Name == input.Name {
			log.Println(category.Name)
			return nil, errors.New("category Task already exists")
		}
	}

	if input.Name == category.Name && input.Description == category.Description {
		return nil, errors.New("no field was changed")
	}

	if input.Name != " " && input.Name != category.Name {
		category.Name = input.Name
	}
	if input.Description != " " && input.Description != category.Description {
		category.Description = input.Description
	}

	err = uc.categoryTaskRepository.Update(ctx, category)
	if err != nil {
		return nil, errors.Wrap(err, "error while updating category")
	}
	categoryUpdated, err := uc.categoryTaskRepository.FindById(ctx, input.UserId, input.Id)
	if err != nil {
		return nil, errors.Wrap(err, "error while updating category")
	}
	if categoryUpdated == nil {
		return nil, errors.Wrap(err, "error while updating category")
	}
	return categoryUpdated, nil
}
