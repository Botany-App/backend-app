package usecases

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type DeleteUserInputDTO struct {
	Id string `json:"id"`
}

type DeleteUserUseCase struct {
	userRepository entities.UserRepository
}

func NewDeleteUserUseCase(userRepository entities.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, input DeleteUserInputDTO) error {
	log.Println("DeleteUserUseCase - Execute")
	err := uc.userRepository.Delete(ctx, input.Id)
	if err != nil {
		log.Println("Erro ao deletar usuário")
		return err
	}
	return nil
}
