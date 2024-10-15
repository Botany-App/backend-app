package repositories

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository é o mock da interface UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *entities.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*entities.User, error) {
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

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserRepository(t *testing.T) {
	mockRepo := new(MockUserRepository)

	// Dados de teste
	userID := uuid.New()
	user := &entities.User{
		ID:        userID.String(),
		Name:      "John Doe",
		Email:     "johndoe@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("Create User", func(t *testing.T) {
		// Configurar o mock para retornar nil ao criar um usuário
		mockRepo.On("Create", user).Return(nil)

		err := mockRepo.Create(user)
		assert.NoError(t, err)

		// Verificar se o método Create foi chamado com o usuário correto
		mockRepo.AssertCalled(t, "Create", user)
	})

	t.Run("Get User by ID", func(t *testing.T) {
		// Configurar o mock para retornar o usuário correto
		mockRepo.On("GetByID", userID).Return(user, nil)

		savedUser, err := mockRepo.GetByID(userID)
		assert.NoError(t, err)
		assert.Equal(t, user, savedUser)

		// Verificar se o método GetByID foi chamado com o ID correto
		mockRepo.AssertCalled(t, "GetByID", userID)
	})

	t.Run("Get User by Email", func(t *testing.T) {
		// Configurar o mock para retornar o usuário com o e-mail correto
		mockRepo.On("GetByEmail", "johndoe@example.com").Return(user, nil)

		savedUser, err := mockRepo.GetByEmail("johndoe@example.com")
		assert.NoError(t, err)
		assert.Equal(t, user, savedUser)

		// Verificar se o método GetByEmail foi chamado com o e-mail correto
		mockRepo.AssertCalled(t, "GetByEmail", "johndoe@example.com")
	})

	t.Run("Update User", func(t *testing.T) {
		// Atualizando o nome e e-mail do usuário
		user.Name = "John Updated"
		user.Email = "johnupdated@example.com"

		// Configurar o mock para não retornar erros ao atualizar
		mockRepo.On("Update", user).Return(nil)

		err := mockRepo.Update(user)
		assert.NoError(t, err)

		// Verificar se o método Update foi chamado com os dados corretos
		mockRepo.AssertCalled(t, "Update", user)
	})

	t.Run("Delete User", func(t *testing.T) {
		// Configurar o mock para não retornar erros ao deletar o usuário
		mockRepo.On("Delete", userID).Return(nil)

		err := mockRepo.Delete(userID)
		assert.NoError(t, err)

		// Verificar se o método Delete foi chamado com o ID correto
		mockRepo.AssertCalled(t, "Delete", userID)
	})
}
