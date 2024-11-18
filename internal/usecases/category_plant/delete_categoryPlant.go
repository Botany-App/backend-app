package usecases_categoryplant

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type DeleteCategoryPlantUseCase struct {
	DeleteCategoryPlant entities.CategoryPlantRepository
}

type DeleteCategoryPlantDTO struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func NewDeleteCategoryPlantUseCase(repository entities.CategoryPlantRepository) *DeleteCategoryPlantUseCase {
	return &DeleteCategoryPlantUseCase{
		DeleteCategoryPlant: repository,
	}
}

func (uc *DeleteCategoryPlantUseCase) Execute(ctx context.Context, input DeleteCategoryPlantDTO) error {
	log.Print("DeleteCategoryPlantUseCase.Execute")
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return errors.New("invalid user id")
	}
	id, err := uuid.Parse(input.ID)
	if err != nil {
		return errors.New("invalid id")
	}

	err = uc.DeleteCategoryPlant.Delete(ctx, userID, id)
	if err != nil {
		return errors.New("error deleting category plant")
	}
	return nil
}
