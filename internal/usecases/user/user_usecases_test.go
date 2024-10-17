package usecases

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

// Implementação dos métodos do repositório mockado
func (m *MockUserRepository) Create(user *entities.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id string) (*entities.User, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*entities.User, error) {
	args := m.Called(email)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *entities.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateUserUseCase_Execute(t *testing.T) {
	mockRepo := new(MockUserRepository)
	createUserUseCase := NewCreateUserUseCase(mockRepo)

	// Dados de entrada
	input := CreateUserInputDTO{
		Name:     "John Doe",
		Email:    "johndoe@example.com",
		Password: "password123",
	}

	// Mock do retorno de entities.NewUser
	user := &entities.User{
		ID:        uuid.New(),
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Configurando o comportamento do mock
	mockRepo.On("Create", mock.Anything).Return(nil)

	// Executa o caso de uso
	output, err := createUserUseCase.Execute(input)

	// Verificações
	assert.NoError(t, err)
	assert.Equal(t, user.Name, output.Name)
	assert.Equal(t, user.Email, output.Email)

	// Verifica se o método Create foi chamado corretamente
	mockRepo.AssertCalled(t, "Create", mock.Anything)
}

// Teste para o método GetById
func TestGetUserUseCase_GetById(t *testing.T) {
	mockRepo := new(MockUserRepository)
	getUserUseCase := NewGetByIdUserUseCase(mockRepo)

	// Dados de entrada
	userID := uuid.New() // Exemplo de ID em formato string
	input := GetByIdInputDTO{ID: userID.String()}

	// Dados simulados
	user := &entities.User{
		ID:        userID,
		Name:      "John Doe",
		Email:     "johndoe@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Configuração do mock
	mockRepo.On("GetByID", userID).Return(user, nil)

	// Chama o método GetById
	output, err := getUserUseCase.Execute(input)

	// Verificações
	assert.NoError(t, err)
	assert.Equal(t, user.ID, output.ID)
	assert.Equal(t, user.Name, output.Name)
	assert.Equal(t, user.Email, output.Email)

	// Verifica se o método GetByID foi chamado corretamente
	mockRepo.AssertCalled(t, "GetByID", userID)
}

// Teste para o método GetByEmail
func TestGetUserUseCase_GetByEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	getUserUseCase := NewGetByEmailUserUseCase(mockRepo)

	// Dados de entrada
	email := "johndoe@example.com"
	input := GetByEmailInputDTO{Email: email}

	// Dados simulados
	user := &entities.User{
		ID:        uuid.New(),
		Name:      "John Doe",
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Configuração do mock
	mockRepo.On("GetByEmail", email).Return(user, nil)

	// Chama o método GetByEmail
	output, err := getUserUseCase.Execute(input)

	// Verificações
	assert.NoError(t, err)
	assert.Equal(t, user.ID, output.ID)
	assert.Equal(t, user.Name, output.Name)
	assert.Equal(t, user.Email, output.Email)

	// Verifica se o método GetByEmail foi chamado corretamente
	mockRepo.AssertCalled(t, "GetByEmail", email)
}
