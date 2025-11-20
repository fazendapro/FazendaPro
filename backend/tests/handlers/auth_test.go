package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/fazendapro/FazendaPro-api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	user := &models.User{
		ID:     1,
		FarmID: 1,
		Person: &models.Person{
			Email: "test@example.com",
		},
	}

	refreshToken := &models.RefreshToken{
		ID:     1,
		UserID: 1,
		Token:  "valid-refresh-token",
	}
	refreshToken.User = *user

	mockRefreshTokenRepo.On("FindByToken", "valid-refresh-token").Return(refreshToken, nil)

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	reqData := map[string]string{
		"refresh_token": "valid-refresh-token",
	}
	jsonData, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.RefreshToken(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	assert.Contains(t, response, "access_token")
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_InvalidMethod(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	req := httptest.NewRequest("GET", "/api/auth/refresh", nil)
	w := httptest.NewRecorder()

	authHandler.RefreshToken(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestAuthHandler_RefreshToken_InvalidJSON(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.RefreshToken(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_RefreshToken_InvalidToken(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	mockRefreshTokenRepo.On("FindByToken", "invalid-token").Return(nil, nil)

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	reqData := map[string]string{
		"refresh_token": "invalid-token",
	}
	jsonData, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.RefreshToken(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_RepositoryError(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	mockRefreshTokenRepo.On("FindByToken", "token").Return(nil, errors.New("erro ao buscar"))

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	reqData := map[string]string{
		"refresh_token": "token",
	}
	jsonData, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.RefreshToken(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	mockRefreshTokenRepo.On("DeleteByToken", "refresh-token").Return(nil)

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	reqData := map[string]string{
		"refresh_token": "refresh-token",
	}
	jsonData, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/api/auth/logout", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Logout(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthHandler_Logout_InvalidMethod(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	req := httptest.NewRequest("GET", "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	authHandler.Logout(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestAuthHandler_Logout_InvalidJSON(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	req := httptest.NewRequest("POST", "/api/auth/logout", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Logout(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Logout_RepositoryError(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	mockRefreshTokenRepo.On("DeleteByToken", "token").Return(errors.New("erro ao deletar"))

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	reqData := map[string]string{
		"refresh_token": "token",
	}
	jsonData, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/api/auth/logout", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Logout(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthHandler_Register_InvalidData(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	reqData := map[string]interface{}{
		"user": map[string]interface{}{
			"farm_id": 1,
		},
		"person": map[string]interface{}{
			"first_name": "",
			"last_name":  "Silva",
			"email":      "invalid-email",
			"password":   "123",
		},
	}
	jsonData, _ := json.Marshal(reqData)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
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

	mockUserRepo.On("FindByPersonEmail", "test@example.com").Return(expectedUser, nil).Twice()

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "wrong-password",
	}
	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response["success"].(bool))
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Register_GetUserByEmailError(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	reqData := map[string]interface{}{
		"user": map[string]interface{}{
			"farm_id": 1,
		},
		"person": map[string]interface{}{
			"first_name": "João",
			"last_name":  "Silva",
			"email":      "joao@example.com",
			"password":   "password123",
			"cpf":        "12345678901",
		},
	}
	jsonData, _ := json.Marshal(reqData)

	mockUserRepo.On("FarmExists", uint(1)).Return(true, nil)
	mockUserRepo.On("CreateWithPerson", mock.AnythingOfType("*models.User"), mock.AnythingOfType("*models.Person")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		user.ID = 1
	})
	mockUserRepo.On("FindByIDWithPerson", uint(1)).Return(&models.User{ID: 1, FarmID: 1}, nil)
	mockUserRepo.On("CreateUserFarm", mock.AnythingOfType("*models.UserFarm")).Return(nil)
	mockUserRepo.On("FindByPersonEmail", "joao@example.com").Return(nil, errors.New("database error"))

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Register(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Register_UserNotFoundAfterCreation(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	reqData := map[string]interface{}{
		"user": map[string]interface{}{
			"farm_id": 1,
		},
		"person": map[string]interface{}{
			"first_name": "João",
			"last_name":  "Silva",
			"email":      "joao@example.com",
			"password":   "password123",
			"cpf":        "12345678901",
		},
	}
	jsonData, _ := json.Marshal(reqData)

	mockUserRepo.On("FarmExists", uint(1)).Return(true, nil)
	mockUserRepo.On("CreateWithPerson", mock.AnythingOfType("*models.User"), mock.AnythingOfType("*models.Person")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		user.ID = 1
	})
	mockUserRepo.On("FindByIDWithPerson", uint(1)).Return(&models.User{ID: 1, FarmID: 1}, nil)
	mockUserRepo.On("CreateUserFarm", mock.AnythingOfType("*models.UserFarm")).Return(nil)
	mockUserRepo.On("FindByPersonEmail", "joao@example.com").Return(nil, nil)

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Register(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Register_RefreshTokenCreationError(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	hashedPassword, _ := utils.HashPassword("password123")
	expectedUser := &models.User{
		ID:     1,
		FarmID: 1,
		Person: &models.Person{
			ID:        1,
			FirstName: "João",
			LastName:  "Silva",
			Email:     "joao@example.com",
			Password:  hashedPassword,
		},
	}

	reqData := map[string]interface{}{
		"user": map[string]interface{}{
			"farm_id": 1,
		},
		"person": map[string]interface{}{
			"first_name": "João",
			"last_name":  "Silva",
			"email":      "joao@example.com",
			"password":   "password123",
			"cpf":        "12345678901",
		},
	}
	jsonData, _ := json.Marshal(reqData)

	mockUserRepo.On("FarmExists", uint(1)).Return(true, nil)
	mockUserRepo.On("CreateWithPerson", mock.AnythingOfType("*models.User"), mock.AnythingOfType("*models.Person")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		user.ID = 1
	})
	mockUserRepo.On("FindByIDWithPerson", uint(1)).Return(&models.User{ID: 1, FarmID: 1}, nil)
	mockUserRepo.On("CreateUserFarm", mock.AnythingOfType("*models.UserFarm")).Return(nil)
	mockUserRepo.On("FindByPersonEmail", "joao@example.com").Return(expectedUser, nil)
	mockRefreshTokenRepo.On("Create", uint(1), mock.AnythingOfType("time.Time")).Return(nil, errors.New("database error"))

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Register(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidMethod(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	req := httptest.NewRequest("GET", "/api/auth/login", nil)
	w := httptest.NewRecorder()

	authHandler.Login(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Login(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_GetUserByEmailError(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	mockUserRepo.On("FindByPersonEmail", "test@example.com").Return(nil, errors.New("database error"))

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Login(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}

	mockUserRepo.On("FindByPersonEmail", "test@example.com").Return(nil, nil)

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Login_ValidatePasswordError(t *testing.T) {
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

	refreshToken := &models.RefreshToken{
		ID:     1,
		UserID: 1,
		Token:  "refresh-token",
	}

	mockUserRepo.On("FindByPersonEmail", "test@example.com").Return(expectedUser, nil).Twice()
	mockRefreshTokenRepo.On("Create", uint(1), mock.AnythingOfType("time.Time")).Return(refreshToken, nil)

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

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
	mockRefreshTokenRepo.AssertExpectations(t)
}


func TestAuthHandler_Login_RefreshTokenError(t *testing.T) {
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

	mockUserRepo.On("FindByPersonEmail", "test@example.com").Return(expectedUser, nil).Twice()
	mockRefreshTokenRepo.On("Create", uint(1), mock.AnythingOfType("time.Time")).Return(nil, errors.New("database error"))

	userService := service.NewUserService(mockUserRepo)
	authHandler := handlers.NewAuthHandler(userService, mockRefreshTokenRepo, "test-secret")

	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	authHandler.Login(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
}
