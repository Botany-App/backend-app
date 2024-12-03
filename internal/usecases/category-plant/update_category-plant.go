package usecases_categoryplant

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type UpdateCategoryPlantInputDTO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCategoryPlantUseCase struct {
	categoryPlantRepository entities.CategoryPlantRepository
}

func NewUpdateCategoryPlantUseCase(categoryPlantRepository entities.CategoryPlantRepository) *UpdateCategoryPlantUseCase {
	return &UpdateCategoryPlantUseCase{
		categoryPlantRepository: categoryPlantRepository,
	}
}

func (uc *UpdateCategoryPlantUseCase) Execute(ctx context.Context, input UpdateCategoryPlantInputDTO, userId string) (*entities.CategoryPlant, error) {
	log.Println("UpdateCategoryPlantUseCase - Execute")
	category, err := uc.categoryPlantRepository.FindById(ctx, userId, input.Id)
	if err != nil {
		return nil, errors.Wrap(err, "error while searching category")
	}
	if category == nil {
		return nil, nil
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

	err = uc.categoryPlantRepository.Update(ctx, category)
	if err != nil {
		return nil, errors.Wrap(err, "error while updating category")
	}
	categoryUpdated, err := uc.categoryPlantRepository.FindById(ctx, userId, input.Id)
	if err != nil {
		return nil, errors.Wrap(err, "error while updating category")
	}
	if categoryUpdated == nil {
		return nil, errors.Wrap(err, "error while updating category")
	}
	return categoryUpdated, nil
}
