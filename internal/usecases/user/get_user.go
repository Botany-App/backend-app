package usecases

import "github.com/lucasBiazon/botany-back/internal/entities"

type GetByIdInputDTO struct {
	ID string `json:"id"`
}

type GetByEmailInputDTO struct {
	Email string `json:"email"`
}

type GetUserOutputDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetUserUseCase struct {
	userRepo entities.UserRepository
}

func NewGetUserUseCase(repo entities.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{userRepo: repo}
}

func (uc *GetUserUseCase) GetById(input GetByIdInputDTO) (*GetUserOutputDTO, error) {
	user, err := uc.userRepo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &GetUserOutputDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

func (uc *GetUserUseCase) GetByEmail(input GetByEmailInputDTO) (*GetUserOutputDTO, error) {
	user, err := uc.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}

	return &GetUserOutputDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}
