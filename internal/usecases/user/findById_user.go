package usecases

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type FindUserByIdInputDTO struct {
	Id string `json:"id"`
}

type FindUserByIdUseCase struct {
	userRepository entities.UserRepository
}

func NewFindUserByIdUseCase(userRepository entities.UserRepository) *FindUserByIdUseCase {
	return &FindUserByIdUseCase{userRepository: userRepository}
}

func (uc *FindUserByIdUseCase) Execute(ctx context.Context, input FindUserByIdInputDTO) (*entities.User, error) {
	log.Println("FindUserByIdUseCase - Execute")
	user, err := uc.userRepository.FindByID(ctx, input.Id)
	if err != nil {
		log.Println("Erro ao buscar usuário pelo ID")
		return nil, err
	}
	if user == nil {
		log.Println("Usuário não encontrado")
		return nil, nil
	}
	return user, nil
}
