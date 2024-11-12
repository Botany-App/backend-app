package usecases

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type RegisterUserInputDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ConfirmEmailInputDTO struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type ResendTokenInputDTO struct {
	Email string `json:"email"`
}
type RegisterUserUseCase struct {
	userRepository entities.UserRepository
}

func NewRegisterUserUseCase(userRepository entities.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{userRepository: userRepository}
}

func (u *RegisterUserUseCase) StartRegistration(ctx context.Context, input RegisterUserInputDTO) error {
	user, err := entities.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return err
	}
	log.Print(user)
	if err := u.userRepository.Create(ctx, user); err != nil {
		return err
	}

	tokenVerification, err := utils.GenerateCode()
	if err != nil {
		return err
	}
	if err = u.userRepository.StoreToken(ctx, user.Email, tokenVerification); err != nil {
		return err
	}

	if err = utils.SendEmail(user.Email, tokenVerification); err != nil {
		return err
	}
	return nil
}

func (u *RegisterUserUseCase) ConfirmEmail(ctx context.Context, email, token string) error {
	err := u.userRepository.ActivateAccount(ctx, email, token)
	if err != nil {
		return err
	}
	return nil
}

func (u *RegisterUserUseCase) ResendToken(ctx context.Context, email string) error {
	newToken, err := utils.GenerateCode()
	if err != nil {
		return err
	}
	err = u.userRepository.ResendToken(ctx, email, newToken)
	if err != nil {
		return err
	}
	return nil
}
