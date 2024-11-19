package usecases

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindUserByIdInputDTO struct {
	ID string `json:"id"`
}

type FindUserByIdUseCase struct {
	userRepository entities.UserRepository
}

func NewFindUserByIdUseCase(userRepository entities.UserRepository) *FindUserByIdUseCase {
	return &FindUserByIdUseCase{userRepository: userRepository}
}

func (uc *FindUserByIdUseCase) Execute(ctx context.Context, input FindUserByIdInputDTO) (*entities.User, error) {
	log.Println("FindUserByIdUseCase - Execute")
	user, err := uc.userRepository.FindByID(ctx, input.ID)
	if err != nil {
		log.Println("Erro ao buscar usuário pelo ID")
		return nil, err
	}
	return user, nil
}
