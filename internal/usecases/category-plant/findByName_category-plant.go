package usecases_categoryplant

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type FindByNameCategoryPlantUseCase struct {
	CategoryPlantRepository entities.CategoryPlantRepository
}

type FindByNameCategoryPlantInputDTO struct {
	Name   string `json:"name"`
	UserId string `json:"userId"`
}

func NewFindByNameCategoryPlantUseCase(categoryPlantRepository entities.CategoryPlantRepository) *FindByNameCategoryPlantUseCase {
	return &FindByNameCategoryPlantUseCase{
		CategoryPlantRepository: categoryPlantRepository,
	}
}
func (uc *FindByNameCategoryPlantUseCase) Execute(ctx context.Context, input FindByNameCategoryPlantInputDTO) ([]*entities.CategoryPlant, error) {
	log.Println("FindByNameCategoryPlantUseCase - Execute")
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	categoriesPlant, err := uc.CategoryPlantRepository.FindByName(ctx, input.UserId, input.Name)
	if err != nil {
		return nil, errors.Wrap(err, "error finding all category plant")
	}
	return categoriesPlant, nil
}
