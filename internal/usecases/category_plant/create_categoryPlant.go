package usecases_categoryplant

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CreateCategoryPlantInputDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	LIMIT       int    `json:"limit"`
	OFFSET      int    `json:"offset"`
}

type CreateCategoryPlantUseCase struct {
	CreateCategoryPlant entities.CategoryPlantRepository
}

func NewCreateCategoryPlantUseCase(createCategoryPlant entities.CategoryPlantRepository) *CreateCategoryPlantUseCase {
	return &CreateCategoryPlantUseCase{CreateCategoryPlant: createCategoryPlant}
}

func (uc *CreateCategoryPlantUseCase) Execute(ctx context.Context, input CreateCategoryPlantInputDTO) error {
	log.Println("CreateCategoryPlantUseCase - Execute")
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return errors.New("error parsing user id")
	}
	if input.LIMIT <= 0 {
		input.LIMIT = 10 // Default limit
	}
	if input.OFFSET < 0 {
		input.OFFSET = 0
	}
	category, err := entities.NewCategoryPlant(input.Name, input.Description, userID)
	if err != nil {
		return errors.New("error creating category plant")
	}
	categoryPlantExists, _ := uc.CreateCategoryPlant.FindByName(ctx, userID, category.Name, input.LIMIT, input.OFFSET)
	if categoryPlantExists != nil {
		return errors.New("category plant already exists")
	}
	err = uc.CreateCategoryPlant.Create(ctx, category)
	if err != nil {
		return errors.New("error creating category plant")
	}
	return nil
}
