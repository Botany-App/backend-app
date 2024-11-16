package usecases

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type DeleteUserInputDTO struct {
	ID string `json:"id"`
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
	log.Println("--> DELETE USER")
	err := uc.userRepository.Delete(ctx, input.ID)
	if err != nil {
		log.Println("Erro ao deletar usuário")
		return err
	}
	log.Println("<-Usuário deletado com sucesso")
	return nil
}
