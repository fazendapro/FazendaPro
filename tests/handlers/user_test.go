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
	"github.com/fazendapro/FazendaPro-api/tests/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserRouter(mockRepo *services.MockUserRepository) (*chi.Mux, *services.MockUserRepository) {
	userService := service.NewUserService(mockRepo)
	userHandler := handlers.NewUserHandler(userService)
	r := chi.NewRouter()
	r.Get("/users", userHandler.GetUser)
	r.Post("/users", userHandler.CreateUser)
	r.Get("/users/person", userHandler.GetUserWithPerson)
	r.Put("/users/person", userHandler.UpdatePersonData)
	return r, mockRepo
}

func TestUserHandler_GetUser_Success(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	expectedUser := &models.User{
		ID:     1,
		FarmID: 1,
		Person: &models.Person{
			Email:     "test@example.com",
			FirstName: "João",
			LastName:  "Silva",
		},
	}

	req, _ := http.NewRequest("GET", "/users?email=test@example.com", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByPersonEmail", "test@example.com").Return(expectedUser, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_GetUser_MissingEmail(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/users?email=notfound@example.com", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByPersonEmail", "notfound@example.com").Return(nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	userData := map[string]interface{}{
		"user": map[string]interface{}{
			"farm_id": 1,
		},
		"person": map[string]interface{}{
			"first_name": "João",
			"last_name":  "Silva",
			"email":      "joao@example.com",
			"password":   "password123",
			"cpf":        "12345678900",
		},
	}

	jsonData, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("FarmExists", uint(1)).Return(true, nil)
	personID := uint(1)
	mockRepo.On("CreateWithPerson", mock.AnythingOfType("*models.User"), mock.AnythingOfType("*models.Person")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		user.ID = 1
		user.PersonID = &personID
	})
	mockRepo.On("CreateUserFarm", mock.AnythingOfType("*models.UserFarm")).Return(nil)
	mockRepo.On("FindByIDWithPerson", uint(1)).Return(&models.User{ID: 1, PersonID: &personID, FarmID: 1}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_CreateUser_InvalidMethod(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	userService := service.NewUserService(mockRepo)
	userHandler := handlers.NewUserHandler(userService)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	userHandler.CreateUser(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestUserHandler_CreateUser_InvalidJSON(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_CreateUser_ServiceError(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	userData := map[string]interface{}{
		"user": map[string]interface{}{
			"farm_id": 1,
		},
		"person": map[string]interface{}{
			"first_name": "",
			"last_name":  "Silva",
			"email":      "joao@example.com",
			"password":   "password123",
			"cpf":        "12345678900",
		},
	}

	jsonData, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_GetUserWithPerson_Success(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	personID := uint(1)
	expectedUser := &models.User{
		ID:       1,
		PersonID: &personID,
		FarmID:   1,
		Person: &models.Person{
			Email:     "test@example.com",
			FirstName: "João",
			LastName:  "Silva",
		},
	}

	req, _ := http.NewRequest("GET", "/users/person?id=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByIDWithPerson", uint(1)).Return(expectedUser, nil).Once()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_GetUserWithPerson_MissingID(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/users/person", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_GetUserWithPerson_NotFound(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/users/person?id=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByIDWithPerson", uint(1)).Return(nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_UpdatePersonData_Success(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	personData := map[string]interface{}{
		"first_name": "João",
		"last_name":  "Silva Atualizado",
		"email":      "joao@example.com",
	}

	jsonData, _ := json.Marshal(personData)
	req, _ := http.NewRequest("PUT", "/users/person", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("UpdatePersonData", uint(1), mock.AnythingOfType("*models.Person")).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_UpdatePersonData_InvalidMethod(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	userService := service.NewUserService(mockRepo)
	userHandler := handlers.NewUserHandler(userService)

	req, _ := http.NewRequest("POST", "/users/person", nil)
	w := httptest.NewRecorder()

	userHandler.UpdatePersonData(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestUserHandler_UpdatePersonData_InvalidJSON(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	req, _ := http.NewRequest("PUT", "/users/person", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_UpdatePersonData_ServiceError(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	router, _ := setupUserRouter(mockRepo)

	personData := map[string]interface{}{
		"first_name": "João",
		"last_name":  "Silva",
	}

	jsonData, _ := json.Marshal(personData)
	req, _ := http.NewRequest("PUT", "/users/person", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("UpdatePersonData", uint(1), mock.AnythingOfType("*models.Person")).Return(errors.New("erro ao atualizar"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}
