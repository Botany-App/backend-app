package usecases

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UpdateUserUseCase struct {
	UserRepository entities.UserRepository
}

type UpdateUserInputDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func NewUpdateUserUseCase(userRepo entities.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		UserRepository: userRepo,
	}
}

// Execute atualiza o usuário com base no DTO fornecido
func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInputDTO) error {
	user, err := uc.UserRepository.GetByID(ctx, input.ID)
	log.Println(user.Name, user.Email)
	if err != nil {
		return errors.New("erro ao buscar usuário")
	}

	if input.Name == user.Name && input.Email == user.Email {
		return errors.New("nenhum campo foi alterado")
	}
	if input.Name != "" && input.Name != user.Name {
		user.Name = input.Name
	}

	if input.Email != "" && input.Email != user.Email {
		user.Email = input.Email
	}

	if err := uc.UserRepository.Update(ctx, user); err != nil {
		return errors.New("erro ao atualizar usuário")
	}

	return nil
}
