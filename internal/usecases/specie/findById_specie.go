package usecases_specie

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByIdSpecieInputDTO struct {
	Id string `json:"id"`
}
type FindByIdSpecieUseCase struct {
	SpecieRepository entities.SpecieRepository
}

func NewFindByIdSpecieUseCase(specieRepository entities.SpecieRepository) *FindByIdSpecieUseCase {
	return &FindByIdSpecieUseCase{SpecieRepository: specieRepository}
}

func (f *FindByIdSpecieUseCase) Execute(ctx context.Context, input FindByIdSpecieInputDTO) (*entities.Specie, error) {
	log.Default().Println("FindByIdSpecieUseCase - Execute")
	specie, err := f.SpecieRepository.FindById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	return specie, nil
}
