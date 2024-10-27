package usecases

import (
	"fmt"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CreateUserInputDTO struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}

type CreateUserOutputDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateUserUseCase struct {
	userRepo entities.UserRepository
}

func NewCreateUserUseCase(repo entities.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: repo,
	}
}

func (uc *CreateUserUseCase) Execute(input CreateUserInputDTO) error {
	cacheKey := fmt.Sprintf("user:%s", input.Email)

	user, err := uc.userRepo.VerifyEmail(cacheKey, input.Code)
	if err != nil {
		return err
	}

	err = uc.userRepo.Create(user)
	if err != nil {
		return err
	}

	fmt.Println("User created successfully")
	return nil

}
