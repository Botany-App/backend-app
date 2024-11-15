package usecases

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	services "github.com/lucasBiazon/botany-back/internal/service"
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

func (uc *RegisterUserUseCase) StartRegistration(ctx context.Context, input RegisterUserInputDTO) error {
	user, err := entities.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return err
	}
	log.Println("---Criando user---")
	if err := uc.userRepository.Create(ctx, user); err != nil {
		return err
	}

	log.Println("---Gerando token---")
	tokenVerification, err := services.NewEmailService().GenerateCode()
	if err != nil {
		return err
	}

	log.Println("---Salvando token---")
	if err = uc.userRepository.StoreToken(ctx, user.Email, tokenVerification); err != nil {
		return err
	}

	log.Println("---Enviando email---")
	if err = services.NewEmailService().SendEmail(user.Email, tokenVerification); err != nil {
		return err
	}
	return nil
}

func (uc *RegisterUserUseCase) ConfirmEmail(ctx context.Context, email, token string) error {
	log.Println("---Confirming email---")
	err := uc.userRepository.ActivateAccount(ctx, email, token)
	if err != nil {
		return err
	}
	return nil
}

func (uc *RegisterUserUseCase) ResendToken(ctx context.Context, email string) error {
	log.Println("---Gerando um novo token---")
	newToken, err := services.NewEmailService().GenerateCode()
	if err != nil {
		return err
	}
	log.Println("---Resend token---")
	err = uc.userRepository.ResendToken(ctx, email, newToken)
	if err != nil {
		return err
	}
	return nil
}
