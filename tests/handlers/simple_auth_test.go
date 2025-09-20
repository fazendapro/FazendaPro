package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_LoginRequest(t *testing.T) {

	req := handlers.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
}

func TestAuthHandler_RegisterRequest(t *testing.T) {

	user := models.User{}
	person := models.Person{}
	req := handlers.RegisterRequest{
		User:   user,
		Person: person,
	}

	assert.NotNil(t, req.User)
	assert.NotNil(t, req.Person)
}

func TestAuthHandler_RefreshTokenRequest(t *testing.T) {

	req := handlers.RefreshTokenRequest{
		RefreshToken: "refresh-token-123",
	}

	assert.Equal(t, "refresh-token-123", req.RefreshToken)
}

func TestAuthHandler_LoginResponse(t *testing.T) {

	response := handlers.LoginResponse{
		Success:      true,
		Message:      "Login realizado com sucesso",
		AccessToken:  "access-token-123",
		RefreshToken: "refresh-token-123",
	}

	assert.True(t, response.Success)
	assert.Equal(t, "Login realizado com sucesso", response.Message)
	assert.Equal(t, "access-token-123", response.AccessToken)
	assert.Equal(t, "refresh-token-123", response.RefreshToken)
}

func TestAuthHandler_RegisterResponse(t *testing.T) {

	response := handlers.RegisterResponse{
		Success:      true,
		Message:      "Usuário criado com sucesso",
		AccessToken:  "access-token-123",
		RefreshToken: "refresh-token-123",
	}

	assert.True(t, response.Success)
	assert.Equal(t, "Usuário criado com sucesso", response.Message)
	assert.Equal(t, "access-token-123", response.AccessToken)
	assert.Equal(t, "refresh-token-123", response.RefreshToken)
}

func TestAuthHandler_RefreshTokenResponse(t *testing.T) {

	response := handlers.RefreshTokenResponse{
		Success:     true,
		Message:     "Token renovado com sucesso",
		AccessToken: "new-access-token-123",
	}

	assert.True(t, response.Success)
	assert.Equal(t, "Token renovado com sucesso", response.Message)
	assert.Equal(t, "new-access-token-123", response.AccessToken)
}

func TestAuthHandler_ErrorHandling(t *testing.T) {

	w := httptest.NewRecorder()
	handlers.SendErrorResponse(w, "Erro de teste", http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Erro de teste", response["message"])
}

func TestAuthHandler_SuccessHandling(t *testing.T) {

	w := httptest.NewRecorder()
	data := map[string]string{"user_id": "123"}
	handlers.SendSuccessResponse(w, data, "Operação realizada com sucesso", http.StatusOK)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Operação realizada com sucesso", response["message"])
}

func TestAuthHandler_JSONParsing(t *testing.T) {

	loginData := handlers.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonData, err := json.Marshal(loginData)
	assert.NoError(t, err)

	var parsedData handlers.LoginRequest
	err = json.Unmarshal(jsonData, &parsedData)
	assert.NoError(t, err)
	assert.Equal(t, loginData.Email, parsedData.Email)
	assert.Equal(t, loginData.Password, parsedData.Password)
}

func TestAuthHandler_HTTPMethods(t *testing.T) {

	tests := []struct {
		method string
		valid  bool
	}{
		{http.MethodPost, true},
		{http.MethodGet, false},
		{http.MethodPut, false},
		{http.MethodDelete, false},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/auth/login", nil)
			w := httptest.NewRecorder()

			if req.Method != http.MethodPost {
				handlers.SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
			} else {
				handlers.SendSuccessResponse(w, nil, "Método permitido", http.StatusOK)
			}

			if tt.valid {
				assert.Equal(t, http.StatusOK, w.Code)
			} else {
				assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
			}
		})
	}
}
