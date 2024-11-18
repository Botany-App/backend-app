package usecases_categoryplant

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CreateCategoryPlantDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
}

type CreateCategoryPlantUseCase struct {
	CreateCategoryPlant entities.CategoryPlantRepository
}

func NewCreateCategoryPlantUseCase(createCategoryPlant entities.CategoryPlantRepository) *CreateCategoryPlantUseCase {
	return &CreateCategoryPlantUseCase{CreateCategoryPlant: createCategoryPlant}
}

func (uc *CreateCategoryPlantUseCase) Execute(ctx context.Context, input CreateCategoryPlantDTO) error {
	log.Println("CreateCategoryPlantUseCase - Execute")
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return errors.New("error parsing user id")
	}
	category, err := entities.NewCategoryPlant(input.Name, input.Description, userID)
	if err != nil {
		return errors.New("error creating category plant")
	}
	err = uc.CreateCategoryPlant.Create(ctx, category)
	if err != nil {
		return errors.New("error creating category plant")
	}
	return nil
}
