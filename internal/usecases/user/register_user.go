package usecases

import (
	"context"
	"errors"
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
	log.Println("StartRegistrationUseCase - Execute")
	user, err := entities.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return err
	}

	userExists, err := uc.userRepository.FindByEmail(ctx, user.Email)
	if err != nil {
		log.Print(err)
		return errors.New("erro ao buscar usuário")
	}
	if userExists != nil {
		return errors.New("email já cadastrado")
	}

	if err := uc.userRepository.Create(ctx, user); err != nil {
		return err
	}

	tokenVerification, err := services.NewEmailService().GenerateCode()
	if err != nil {
		return err
	}

	if err = uc.userRepository.StoreToken(ctx, user.Email, tokenVerification); err != nil {
		return err
	}

	if err = services.NewEmailService().SendEmail(user.Email, tokenVerification); err != nil {
		return err
	}
	return nil
}

func (uc *RegisterUserUseCase) ConfirmEmail(ctx context.Context, input ConfirmEmailInputDTO) error {
	log.Println("ConfirmEmailUseCase - Execute")
	err := uc.userRepository.ActivateAccount(ctx, input.Email, input.Token)
	if err != nil {
		return err
	}
	return nil
}

func (uc *RegisterUserUseCase) ResendToken(ctx context.Context, input ResendTokenInputDTO) error {
	log.Println("ResendTokenUseCase - Execute")
	newToken, err := services.NewEmailService().GenerateCode()
	if err != nil {
		return err
	}
	token, err := uc.userRepository.ResendToken(ctx, input.Email, newToken)
	if err != nil {
		return err
	}
	err = services.NewEmailService().SendEmail(input.Email, token)
	if err != nil {
		return err
	}
	return nil
}
