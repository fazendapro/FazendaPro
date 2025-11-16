package services

import (
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService(t *testing.T) {
	t.Run("CreateUser_Success", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		mockUserRepo.On("FarmExists", uint(1)).Return(true, nil)
		mockUserRepo.On("CreateWithPerson", mock.AnythingOfType("*models.User"), mock.AnythingOfType("*models.Person")).Return(nil).Run(func(args mock.Arguments) {
			user := args.Get(0).(*models.User)
			user.ID = 1 // Simular ID retornado pelo banco
		})
		mockUserRepo.On("FindByIDWithPerson", uint(1)).Return(&models.User{ID: 1, FarmID: 1}, nil)
		mockUserRepo.On("CreateUserFarm", mock.AnythingOfType("*models.UserFarm")).Return(nil)

		user := &models.User{
			FarmID: 1,
		}
		person := &models.Person{
			FirstName: "João",
			LastName:  "Silva",
			Email:     "joao@fazenda.com",
			Password:  "senha123",
			CPF:       "12345678901",
		}

		err := userService.CreateUser(user, person)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("GetUserByEmail_Success", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		expectedUser := &models.User{
			ID:     1,
			FarmID: 1,
			Person: &models.Person{
				Email: "joao@fazenda.com",
			},
		}
		mockUserRepo.On("FindByPersonEmail", "joao@fazenda.com").Return(expectedUser, nil)

		user, err := userService.GetUserByEmail("joao@fazenda.com")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uint(1), user.FarmID)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("GetUserByEmail_NotFound", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		mockUserRepo.On("FindByPersonEmail", "inexistente@fazenda.com").Return(nil, nil)

		user, err := userService.GetUserByEmail("inexistente@fazenda.com")

		assert.NoError(t, err)
		assert.Nil(t, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ValidatePassword_Success", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		mockUserRepo.On("ValidatePassword", uint(1), "senha123").Return(true, nil)

		valid, err := userService.ValidatePassword(1, "senha123")

		assert.NoError(t, err)
		assert.True(t, valid)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ValidatePassword_InvalidPassword", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		mockUserRepo.On("ValidatePassword", uint(1), "senhaerrada").Return(false, nil)

		valid, err := userService.ValidatePassword(1, "senhaerrada")

		assert.NoError(t, err)
		assert.False(t, valid)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("UpdatePersonData_Success", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		mockUserRepo.On("UpdatePersonData", uint(1), mock.AnythingOfType("*models.Person")).Return(nil)

		personData := &models.Person{
			ID:        1,
			FirstName: "João Atualizado",
			LastName:  "Silva Santos",
			Email:     "joao.novo@fazenda.com",
			Password:  "novasenha123",
			CPF:       "12345678901",
		}

		err := userService.UpdatePersonData(1, personData)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("GetUserWithPerson_Success", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

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
		mockUserRepo.On("FindByIDWithPerson", uint(1)).Return(expectedUser, nil)

		user, err := userService.GetUserWithPerson(1)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uint(1), user.FarmID)
		assert.NotNil(t, user.Person)
		assert.Equal(t, "João", user.Person.FirstName)
		assert.Equal(t, "joao@fazenda.com", user.Person.Email)
		mockUserRepo.AssertExpectations(t)
	})
}
