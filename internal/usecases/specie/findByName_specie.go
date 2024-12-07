package usecases_specie

import (
	"context"
	"errors"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindByNameSpecieUseCase struct {
	SpecieRepository entities.SpecieRepository
}

type FindByNameSpecieInputDTO struct {
	CommonName string `json:"common_name"`
}

func NewFindByNameSpecieUseCase(specieRepository entities.SpecieRepository) *FindByNameSpecieUseCase {
	return &FindByNameSpecieUseCase{SpecieRepository: specieRepository}
}

func (u *FindByNameSpecieUseCase) Execute(ctx context.Context, input FindByNameSpecieInputDTO) ([]*entities.Specie, error) {
	if input.CommonName == ""  {
		return nil, errors.New("common_name or scientific_name is required")
	}
	specie, err := u.SpecieRepository.FindByName(ctx, input.CommonName)

	if err != nil {
		return nil, err
	}

	return specie, nil
}
