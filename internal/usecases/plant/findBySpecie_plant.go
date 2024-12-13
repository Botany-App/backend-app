package usecases_plant

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindBySpecieNamePlantUseCase struct {
	PlantRepository entities.PlantRepository
}

type FindBySpecieNamePlantUseCaseInputDTO struct {
	UserId     string `json:"user_id"`
	SpecieName string `json:"category_name"`
}

func NewFindBySpecieNamePlantUseCase(repository entities.PlantRepository) *FindBySpecieNamePlantUseCase {
	return &FindBySpecieNamePlantUseCase{
		PlantRepository: repository,
	}
}

func (uc *FindBySpecieNamePlantUseCase) Execute(ctx context.Context, input FindBySpecieNamePlantUseCaseInputDTO) ([]*entities.PlantWithCategory, error) {
	log.Println("FindBySpecieNamePlantUseCase - Execute")
	plants, err := uc.PlantRepository.FindBySpeciesName(ctx, input.UserId, input.SpecieName)
	if err != nil {
		return nil, err
	}
	return plants, nil
}
