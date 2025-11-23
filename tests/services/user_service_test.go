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
			user.ID = 1
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

	t.Run("GetUserByEmail_InvalidEmail", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		user, err := userService.GetUserByEmail("invalid-email")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid email")
		mockUserRepo.AssertNotCalled(t, "FindByPersonEmail")
	})

	t.Run("CreateUser_InvalidEmail", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		user := &models.User{FarmID: 1}
		person := &models.Person{
			FirstName: "João",
			LastName:  "Silva",
			Email:     "invalid-email",
			Password:  "senha123",
			CPF:       "12345678901",
		}

		err := userService.CreateUser(user, person)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email")
		mockUserRepo.AssertNotCalled(t, "CreateWithPerson")
	})

	t.Run("CreateUser_ShortPassword", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		user := &models.User{FarmID: 1}
		person := &models.Person{
			FirstName: "João",
			LastName:  "Silva",
			Email:     "joao@fazenda.com",
			Password:  "12345",
			CPF:       "12345678901",
		}

		err := userService.CreateUser(user, person)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must have at least 6 characters")
		mockUserRepo.AssertNotCalled(t, "CreateWithPerson")
	})

	t.Run("CreateUser_MissingFields", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		testCases := []struct {
			name     string
			user     *models.User
			person   *models.Person
			errorMsg string
		}{
			{
				name:     "MissingFirstName",
				user:     &models.User{FarmID: 1},
				person:   &models.Person{LastName: "Silva", Email: "joao@fazenda.com", Password: "senha123", CPF: "12345678901"},
				errorMsg: "first name is required",
			},
			{
				name:     "MissingLastName",
				user:     &models.User{FarmID: 1},
				person:   &models.Person{FirstName: "João", Email: "joao@fazenda.com", Password: "senha123", CPF: "12345678901"},
				errorMsg: "last name is required",
			},
			{
				name:     "MissingEmail",
				user:     &models.User{FarmID: 1},
				person:   &models.Person{FirstName: "João", LastName: "Silva", Password: "senha123", CPF: "12345678901"},
				errorMsg: "email is required",
			},
			{
				name:     "MissingPassword",
				user:     &models.User{FarmID: 1},
				person:   &models.Person{FirstName: "João", LastName: "Silva", Email: "joao@fazenda.com", CPF: "12345678901"},
				errorMsg: "password is required",
			},
			{
				name:     "MissingCPF",
				user:     &models.User{FarmID: 1},
				person:   &models.Person{FirstName: "João", LastName: "Silva", Email: "joao@fazenda.com", Password: "senha123"},
				errorMsg: "CPF is required",
			},
			{
				name:     "MissingFarmID",
				user:     &models.User{FarmID: 0},
				person:   &models.Person{FirstName: "João", LastName: "Silva", Email: "joao@fazenda.com", Password: "senha123", CPF: "12345678901"},
				errorMsg: "farm ID is required",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := userService.CreateUser(tc.user, tc.person)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			})
		}
	})

	t.Run("CreateUser_FarmDoesNotExist", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		mockUserRepo.On("FarmExists", uint(1)).Return(false, nil)
		mockUserRepo.On("CreateDefaultFarm", uint(1)).Return(nil)
		mockUserRepo.On("CreateWithPerson", mock.AnythingOfType("*models.User"), mock.AnythingOfType("*models.Person")).Return(nil).Run(func(args mock.Arguments) {
			user := args.Get(0).(*models.User)
			user.ID = 1
		})
		mockUserRepo.On("FindByIDWithPerson", uint(1)).Return(&models.User{ID: 1, FarmID: 1}, nil)
		mockUserRepo.On("CreateUserFarm", mock.AnythingOfType("*models.UserFarm")).Return(nil)

		user := &models.User{FarmID: 1}
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

	t.Run("ValidatePasswordByEmail_Success", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		expectedUser := &models.User{
			ID:     1,
			FarmID: 1,
			Person: &models.Person{
				Email:    "joao@fazenda.com",
				Password: "$2a$10$hashedpassword",
			},
		}
		mockUserRepo.On("FindByPersonEmail", "joao@fazenda.com").Return(expectedUser, nil)

		valid, err := userService.ValidatePasswordByEmail("joao@fazenda.com", "password123")

		assert.NoError(t, err)
		assert.NotNil(t, valid)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ValidatePasswordByEmail_InvalidEmail", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		valid, err := userService.ValidatePasswordByEmail("invalid-email", "password123")

		assert.Error(t, err)
		assert.False(t, valid)
		mockUserRepo.AssertNotCalled(t, "FindByPersonEmail")
	})

	t.Run("ValidatePasswordByEmail_UserNotFound", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		mockUserRepo.On("FindByPersonEmail", "notfound@fazenda.com").Return(nil, nil)

		valid, err := userService.ValidatePasswordByEmail("notfound@fazenda.com", "password123")

		assert.NoError(t, err)
		assert.False(t, valid)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("GetUserFarms_Success", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		expectedFarms := []models.Farm{
			{ID: 1},
			{ID: 2},
		}
		mockUserRepo.On("GetUserFarms", uint(1)).Return(expectedFarms, nil)

		farms, err := userService.GetUserFarms(1)

		assert.NoError(t, err)
		assert.Len(t, farms, 2)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("GetUserFarmCount_Success", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		mockUserRepo.On("GetUserFarmCount", uint(1)).Return(int64(3), nil)

		count, err := userService.GetUserFarmCount(1)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("GetUserFarmByID_Success", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		expectedFarm := &models.Farm{ID: 1}
		mockUserRepo.On("GetUserFarmByID", uint(1), uint(1)).Return(expectedFarm, nil)

		farm, err := userService.GetUserFarmByID(1, 1)

		assert.NoError(t, err)
		assert.NotNil(t, farm)
		assert.Equal(t, uint(1), farm.ID)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ShouldAutoSelectFarm_True", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		mockUserRepo.On("GetUserFarmCount", uint(1)).Return(int64(1), nil)

		shouldAuto, err := userService.ShouldAutoSelectFarm(1)

		assert.NoError(t, err)
		assert.True(t, shouldAuto)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("ShouldAutoSelectFarm_False", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		userService := service.NewUserService(mockUserRepo)

		mockUserRepo.On("GetUserFarmCount", uint(1)).Return(int64(2), nil)

		shouldAuto, err := userService.ShouldAutoSelectFarm(1)

		assert.NoError(t, err)
		assert.False(t, shouldAuto)
		mockUserRepo.AssertExpectations(t)
	})
}
