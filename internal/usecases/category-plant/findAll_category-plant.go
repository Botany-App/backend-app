package usecases_categoryplant

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type FindAllCategoryPlantUseCase struct {
	CategoryPlantRepository entities.CategoryPlantRepository
}

func NewFindAllCategoryPlantUseCase(categoryPlantRepository entities.CategoryPlantRepository) *FindAllCategoryPlantUseCase {
	return &FindAllCategoryPlantUseCase{
		CategoryPlantRepository: categoryPlantRepository,
	}
}
func (uc *FindAllCategoryPlantUseCase) Execute(ctx context.Context, input string) ([]*entities.CategoryPlant, error) {
	log.Println("FindAllCategoryPlantUseCase - Execute")
	categoriesPlant, err := uc.CategoryPlantRepository.FindAll(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "error finding by all categories plants")
	}
	return categoriesPlant, nil
}
