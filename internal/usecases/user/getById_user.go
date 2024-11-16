package usecases

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type GetUserByIdUseCase struct {
	userRepository entities.UserRepository
}

type GetUserByIdInputDTO struct {
	ID string `json:"id"`
}

func NewGetUserByIdUseCase(userRepository entities.UserRepository) *GetUserByIdUseCase {
	return &GetUserByIdUseCase{userRepository: userRepository}
}

func (uc *GetUserByIdUseCase) GetUserById(ctx context.Context, id string) (*entities.User, error) {
	log.Println("--> Pegando os dados do usuário pelo ID")
	user, err := uc.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
