package usecases_categoryplant

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type DeleteCategoryPlantUseCase struct {
	CategoryPlantRepository entities.CategoryPlantRepository
}

type DeleteCategoryPlantInputDTO struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

func NewDeleteCategoryPlantUseCase(categoryPlantRepository entities.CategoryPlantRepository) *DeleteCategoryPlantUseCase {
	return &DeleteCategoryPlantUseCase{
		CategoryPlantRepository: categoryPlantRepository,
	}
}
func (uc *DeleteCategoryPlantUseCase) Execute(ctx context.Context, input DeleteCategoryPlantInputDTO) error {
	log.Println("DeleteCategoryPlantUseCase - Execute")

	err := uc.CategoryPlantRepository.Delete(ctx, input.UserId, input.Id)
	if err != nil {
		return errors.Wrap(err, "error delete  category plant")
	}
	return nil
}
