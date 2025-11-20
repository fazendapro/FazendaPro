package repositories

import (
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByPersonEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CreateWithPerson(user *models.User, personData *models.Person) error {
	args := m.Called(user, personData)
	return args.Error(0)
}

func (m *MockUserRepository) FindByIDWithPerson(userID uint) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePersonData(userID uint, personData *models.Person) error {
	args := m.Called(userID, personData)
	return args.Error(0)
}

func (m *MockUserRepository) ValidatePassword(userID uint, password string) (bool, error) {
	args := m.Called(userID, password)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestUserRepository(t *testing.T) {
	t.Run("CreateUser_Success", func(t *testing.T) {
		mockRepo := &MockUserRepository{}

		mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		user := &models.User{
			FarmID: 1,
		}

		err := mockRepo.Create(user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetUserByEmail_Success", func(t *testing.T) {
		mockRepo := &MockUserRepository{}

		expectedUser := &models.User{
			ID:     1,
			FarmID: 1,
		}
		mockRepo.On("FindByEmail", "joao@fazenda.com").Return(expectedUser, nil)

		user, err := mockRepo.FindByEmail("joao@fazenda.com")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uint(1), user.FarmID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetUserByEmail_NotFound", func(t *testing.T) {
		mockRepo := &MockUserRepository{}

		mockRepo.On("FindByEmail", "inexistente@fazenda.com").Return(nil, nil)

		user, err := mockRepo.FindByEmail("inexistente@fazenda.com")

		assert.NoError(t, err)
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetUserByID_Success", func(t *testing.T) {
		mockRepo := &MockUserRepository{}

		expectedUser := &models.User{
			ID:     1,
			FarmID: 1,
		}
		mockRepo.On("FindByIDWithPerson", uint(1)).Return(expectedUser, nil)

		user, err := mockRepo.FindByIDWithPerson(1)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uint(1), user.ID)
		assert.Equal(t, uint(1), user.FarmID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateUser_Success", func(t *testing.T) {
		mockRepo := &MockUserRepository{}

		mockRepo.On("UpdatePersonData", uint(1), mock.AnythingOfType("*models.Person")).Return(nil)

		personData := &models.Person{
			ID:        1,
			FirstName: "João Atualizado",
			LastName:  "Silva Santos",
			Email:     "joao.novo@fazenda.com",
			Password:  "novasenha123",
			CPF:       "12345678901",
		}

		err := mockRepo.UpdatePersonData(1, personData)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteUser_Success", func(t *testing.T) {
		mockRepo := &MockUserRepository{}

		mockRepo.On("ValidatePassword", uint(1), "senha123").Return(true, nil)

		valid, err := mockRepo.ValidatePassword(1, "senha123")

		assert.NoError(t, err)
		assert.True(t, valid)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetUserWithPerson_Success", func(t *testing.T) {
		mockRepo := &MockUserRepository{}

		expectedUser := &models.User{
			ID:     1,
			FarmID: 1,
			Person: &models.Person{
				ID:        1,
				FirstName: "João",
				LastName:  "Silva",
				Email:     "joao@fazenda.com",
			},
		}
		mockRepo.On("FindByIDWithPerson", uint(1)).Return(expectedUser, nil)

		user, err := mockRepo.FindByIDWithPerson(1)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotNil(t, user.Person)
		assert.Equal(t, "João", user.Person.FirstName)
		assert.Equal(t, "joao@fazenda.com", user.Person.Email)
		mockRepo.AssertExpectations(t)
	})
}
