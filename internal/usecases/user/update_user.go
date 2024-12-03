package usecases

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UpdateUserInputDTO struct {
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type UpdateUserUseCase struct {
	UserRepository entities.UserRepository
}

func NewUpdateUserUseCase(userRepo entities.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		UserRepository: userRepo,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInputDTO) (*entities.User, error) {
	log.Println("UpdateUserUseCase - Execute")
	user, err := uc.UserRepository.FindByID(ctx, input.Id)
	if err != nil {
		return nil, errors.New("erro ao buscar usuário")
	}
	if user == nil {
		return nil, nil
	}

	if input.Name == user.Name && input.Email == user.Email {
		return nil, errors.New("nenhum campo foi alterado")
	}
	if input.Name != "" && input.Name != user.Name {
		user.Name = input.Name
	}

	if input.Email != "" && input.Email != user.Email {
		user.Email = input.Email
	}

	if err := uc.UserRepository.Update(ctx, user); err != nil {
		return nil, errors.New("erro ao atualizar usuário")
	}
	user, err = uc.UserRepository.FindByID(ctx, input.Id)
	if err != nil {
		return nil, errors.New("erro ao buscar usuário")
	}
	return user, nil
}
