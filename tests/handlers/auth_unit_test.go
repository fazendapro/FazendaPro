package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/fazendapro/FazendaPro-api/internal/utils"
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

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(userID uint, expiresAt time.Time) (*models.RefreshToken, error) {
	args := m.Called(userID, expiresAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) DeleteByToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteByUserID(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("Login_Success", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}

		hashedPassword, _ := utils.HashPassword("password123")
		expectedUser := &models.User{
			ID:     1,
			FarmID: 1,
			Person: &models.Person{
				Email:    "test@example.com",
				Password: hashedPassword,
			},
		}

		mockUserRepo.On("FindByPersonEmail", "test@example.com").Return(expectedUser, nil)
		mockRefreshTokenRepo.On("Create", uint(1), mock.AnythingOfType("time.Time")).Return(&models.RefreshToken{Token: "refresh_token_123"}, nil)

		userService := service.NewUserService(mockUserRepo)
		authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "secret")

		loginData := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		authHandler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Login_InvalidCredentials", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}

		mockUserRepo.On("FindByPersonEmail", "test@example.com").Return(nil, nil)

		userService := service.NewUserService(mockUserRepo)
		authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "secret")

		loginData := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		authHandler.Login(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockUserRepo.AssertExpectations(t)
	})
}
