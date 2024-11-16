package usecases

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type DeleteUserUseCase struct {
	userRepository entities.UserRepository
}

func NewDeleteUserUseCase(userRepository entities.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id string) error {
	log.Println("--> Deletando usuário")
	err := uc.userRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
