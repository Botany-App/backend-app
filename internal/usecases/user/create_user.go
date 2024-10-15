package usecases

import "github.com/lucasBiazon/botany-back/internal/entities"

type CreateUserInputDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserOutputDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateUserUseCase struct {
	userRepo entities.UserRepository
}

func NewCreateUserUseCase(repo entities.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: repo}
}

func (uc *CreateUserUseCase) Execute(input CreateUserInputDTO) (*CreateUserOutputDTO, error) {
	user, err := entities.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	err = uc.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &CreateUserOutputDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}
