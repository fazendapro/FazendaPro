package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/fazendapro/FazendaPro-api/internal/utils"
	"github.com/fazendapro/FazendaPro-api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Login(t *testing.T) {
	t.Run("Login_Success", func(t *testing.T) {
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}

		hashedPassword, _ := utils.HashPassword("password123")
		expectedUser := &models.User{
			ID:     1,
			FarmID: 1,
			Person: &models.Person{
				Email:    tests.TestEmailExample,
				Password: hashedPassword,
			},
		}

		mockUserRepo.On("FindByPersonEmail", tests.TestEmailExample).Return(expectedUser, nil)
		mockRefreshTokenRepo.On("Create", uint(1), mock.AnythingOfType("time.Time")).Return(&models.RefreshToken{Token: "refresh_token_123"}, nil)

		userService := service.NewUserService(mockUserRepo)
		authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "secret")

		loginData := map[string]string{
			"email":    tests.TestEmailExample,
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
			"email":    tests.TestEmailExample,
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
