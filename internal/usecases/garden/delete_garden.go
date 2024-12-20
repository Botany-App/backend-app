package usecases_garden

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type DeleteGardenUseCase struct {
	GardenRepository entities.GardenRepository
}

type DeleteGardenUseCaseInputDTO struct {
	Id     string `json:"id"`
	UserID string `json:"user_id"`
}

func NewDeleteGardenUseCase(repository entities.GardenRepository) *DeleteGardenUseCase {
	return &DeleteGardenUseCase{
		GardenRepository: repository,
	}
}

func (uc *DeleteGardenUseCase) Execute(ctx context.Context, input DeleteGardenUseCaseInputDTO) error {
	log.Println("DeleteGarden - Execute")
	Garden, err := uc.GardenRepository.FindByID(ctx, input.UserID, input.Id)
	if err != nil {
		return err
	}
	if Garden == nil {
		return err
	}

	err = uc.GardenRepository.Delete(ctx, input.UserID, input.Id)
	if err != nil {
		return err
	}
	return nil
}
