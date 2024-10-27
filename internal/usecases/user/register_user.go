package usecases

import (
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type RegisterUserInputDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserUseCase struct {
	userRepo entities.UserRepository
}

func NewRegisterUserUseCase(repo entities.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo: repo,
	}
}

func (uc *RegisterUserUseCase) Execute(input RegisterUserInputDTO) (string, error) {
	user, err := entities.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return "", err
	}
	err = uc.userRepo.HostCode(user)
	if err != nil {
		return "", err
	}
	return input.Email, nil
}
