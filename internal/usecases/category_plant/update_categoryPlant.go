package usecases_categoryplant

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UpdateCategoryPlantUseCase struct {
	UpdateCategoryPlantRepository entities.CategoryPlantRepository
}

type UpdateCategoryPlantDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
}

func NewUpdateCategoryPlantUseCase(repository entities.CategoryPlantRepository) *UpdateCategoryPlantUseCase {
	return &UpdateCategoryPlantUseCase{UpdateCategoryPlantRepository: repository}
}

func (uc *UpdateCategoryPlantUseCase) Execute(ctx context.Context, input UpdateCategoryPlantDTO) error {
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return errors.New("invalid user id")
	}
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return errors.New("invalid id")
	}

	categoryPlant, err := uc.UpdateCategoryPlantRepository.FindByID(ctx, userID, id)
	if err != nil {
		return errors.New("category plant not found")
	}

	if input.Name == "" && input.Description == "" && input.Name == categoryPlant.Name && input.Description == categoryPlant.Description {
		return errors.New("no data to update")
	}

	if input.Name != "" && input.Name != categoryPlant.Name {
		categoryPlant.Name = input.Name
	}
	if input.Description != "" && input.Description != categoryPlant.Description {
		categoryPlant.Description = input.Description
	}

	err = uc.UpdateCategoryPlantRepository.Update(ctx, categoryPlant)
	if err != nil {
		return errors.New("error updating category plant")
	}

	return nil
}
