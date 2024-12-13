package usecases_plant

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type DeletePlantUseCase struct {
	PlantRepository  entities.PlantRepository
	SpecieRepository entities.SpecieRepository
}

type DeletePlantUseCaseInputDTO struct {
	Id     string `json:"id"`
	UserID string `json:"user_id"`
}

func NewDeletePlantUseCase(repository entities.PlantRepository, specieRepository entities.SpecieRepository) *DeletePlantUseCase {
	return &DeletePlantUseCase{
		PlantRepository:  repository,
		SpecieRepository: specieRepository,
	}
}

func (uc *DeletePlantUseCase) Execute(ctx context.Context, input DeletePlantUseCaseInputDTO) error {
	log.Println("DeletePlant - Execute")
	plant, err := uc.PlantRepository.FindByID(ctx, input.UserID, input.Id)
	if err != nil {
		return err
	}
	if plant == nil {
		return err
	}

	err = uc.PlantRepository.Delete(ctx, input.UserID, input.Id)
	if err != nil {
		return err
	}
	return nil
}
