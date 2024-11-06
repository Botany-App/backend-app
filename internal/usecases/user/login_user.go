package usecases

import (
	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type LoginUserInputDTO struct {
	Email    string
	Password string
}

type LoginUserOutputDTO struct {
	Token string
}

type LoginUserUseCase struct {
	userRepository entities.UserRepository
}

func NewLoginUserUseCase(userRepository entities.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepository: userRepository,
	}
}

func (u *LoginUserUseCase) Execute(input LoginUserInputDTO) (*LoginUserOutputDTO, error) {
	user, err := u.userRepository.Login(input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return nil, err
	}
	return &LoginUserOutputDTO{
		Token: token,
	}, nil
}
