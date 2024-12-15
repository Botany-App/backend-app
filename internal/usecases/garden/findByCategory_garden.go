package usecases_garden

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByCategoryNameGardenUseCase struct {
	GardenRepository entities.GardenRepository
}

type FindByCategoryNameGardenUseCaseInputDTO struct {
	UserId       string `json:"user_id"`
	CategoryName string `json:"category_name"`
}

func NewFindByCategoryNameGardenUseCase(repository entities.GardenRepository) *FindByCategoryNameGardenUseCase {
	return &FindByCategoryNameGardenUseCase{
		GardenRepository: repository,
	}
}

func (uc *FindByCategoryNameGardenUseCase) Execute(ctx context.Context, input FindByCategoryNameGardenUseCaseInputDTO) ([]*entities.GardenOutputDTO, error) {
	log.Println("FindByCategoryNameUseCase - Execute")
	Gardens, err := uc.GardenRepository.FindByCategoryName(ctx, input.UserId, input.CategoryName)
	if err != nil {
		return nil, err
	}
	log.Println("FindByCategoryNameUseCase - Execute - Gardens: ", Gardens)
	return Gardens, nil
}
