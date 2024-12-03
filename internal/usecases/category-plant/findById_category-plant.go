package usecases_categoryplant

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type FindByIdCategoryPlantUseCase struct {
	CategoryPlantRepository entities.CategoryPlantRepository
}

type FindByIdCategoryPlantInputDTO struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

func NewFindByIdCategoryPlantUseCase(categoryPlantRepository entities.CategoryPlantRepository) *FindByIdCategoryPlantUseCase {
	return &FindByIdCategoryPlantUseCase{
		CategoryPlantRepository: categoryPlantRepository,
	}
}
func (uc *FindByIdCategoryPlantUseCase) Execute(ctx context.Context, input FindByIdCategoryPlantInputDTO) (*entities.CategoryPlant, error) {
	log.Println("FindByIdCategoryPlantUseCase - Execute")
	categoryPlant, err := uc.CategoryPlantRepository.FindById(ctx, input.UserId, input.Id)
	if err != nil {
		return nil, errors.Wrap(err, "error finding by id category plant")
	}
	return categoryPlant, nil
}
