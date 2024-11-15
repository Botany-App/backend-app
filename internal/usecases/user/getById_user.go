package usecases

import (
	"context"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type GetUserUseCase struct {
	userRepository entities.UserRepository
}

type GetById struct {
	ID string `json:"id"`
}

func NewGetUserUseCase(userRepository entities.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{userRepository: userRepository}
}

func (uc *GetUserUseCase) GetUserById(ctx context.Context, id string) (*entities.User, error) {
	user, err := uc.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
