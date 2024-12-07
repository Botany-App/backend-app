package usecases_specie

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindAllSpecieUseCase struct {
	SpecieRepository entities.SpecieRepository
}

func NewFindAllSpecieUseCase(specieRepository entities.SpecieRepository) *FindAllSpecieUseCase {
	return &FindAllSpecieUseCase{SpecieRepository: specieRepository}
}

func (uc *FindAllSpecieUseCase) Execute(ctx context.Context) ([]*entities.Specie, error) {
	log.Println("FindAllSpecieUseCase - Execute")
	species, err := uc.SpecieRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return species, nil
}
