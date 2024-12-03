package usecases_categoryplant

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type CreateCategoryPlantInputDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCategoryPlantUseCase struct {
	categoryPlantRepository entities.CategoryPlantRepository
}

func NewCreateCategoryPlantUseCase(categoryPlantRepository entities.CategoryPlantRepository) *CreateCategoryPlantUseCase {
	return &CreateCategoryPlantUseCase{
		categoryPlantRepository: categoryPlantRepository,
	}
}

func (uc *CreateCategoryPlantUseCase) Execute(ctx context.Context, input CreateCategoryPlantInputDTO, userId string) (*entities.CategoryPlant, error) {
	log.Println("CreateCategoryPlantUseCase - Execute")
	newCategoryPlant, err := entities.NewCreateCategoryPlant(input.Name, input.Description, userId)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new category plant")
	}
	existingCategoryPlant, err := uc.categoryPlantRepository.FindByName(ctx, userId, input.Name)
	if err != nil {
		return nil, errors.Wrap(err, "error finding category plant by name")
	}
	if existingCategoryPlant != nil {
		return nil, errors.New("category with the same name already exists")
	}

	id, err := uc.categoryPlantRepository.Create(ctx, newCategoryPlant)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new category plant")
	}

	categoryCreated, err := uc.categoryPlantRepository.FindById(ctx, userId, id)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new category plant")
	}
	return categoryCreated, nil

}
