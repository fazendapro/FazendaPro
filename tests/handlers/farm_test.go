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
	"github.com/fazendapro/FazendaPro-api/tests"
	"github.com/fazendapro/FazendaPro-api/tests/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupFarmRouter(mockRepo *services.MockFarmRepository) (*chi.Mux, *services.MockFarmRepository) {
	farmService := service.NewFarmService(mockRepo)
	farmHandler := handlers.NewFarmHandler(farmService)
	r := chi.NewRouter()
	r.Get("/farm", farmHandler.GetFarm)
	r.Put("/farm", farmHandler.UpdateFarm)
	return r, mockRepo
}

func TestFarmHandler_GetFarm_Success(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	router, _ := setupFarmRouter(mockRepo)

	expectedFarm := &models.Farm{
		ID:       1,
		CompanyID: 1,
		Logo:     "logo.png",
	}

	req, _ := http.NewRequest("GET", tests.EndpointFarmWithID, nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(expectedFarm, nil)
	mockRepo.On("LoadCompanyData", expectedFarm).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestFarmHandler_GetFarm_MissingID(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	router, _ := setupFarmRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/farm", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFarmHandler_GetFarm_InvalidID(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	router, _ := setupFarmRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/farm?id=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFarmHandler_GetFarm_NotFound(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	router, _ := setupFarmRouter(mockRepo)

	req, _ := http.NewRequest("GET", tests.EndpointFarmWithID, nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestFarmHandler_GetFarm_ServiceError(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	router, _ := setupFarmRouter(mockRepo)

	req, _ := http.NewRequest("GET", tests.EndpointFarmWithID, nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("erro ao buscar"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestFarmHandler_GetFarm_LoadCompanyDataError(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	router, _ := setupFarmRouter(mockRepo)

	expectedFarm := &models.Farm{
		ID:       1,
		CompanyID: 1,
	}

	req, _ := http.NewRequest("GET", tests.EndpointFarmWithID, nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(expectedFarm, nil)
	mockRepo.On("LoadCompanyData", expectedFarm).Return(errors.New("erro ao carregar dados"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestFarmHandler_UpdateFarm_Success(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	router, _ := setupFarmRouter(mockRepo)

	farmData := map[string]interface{}{
		"logo": tests.TestFileNewLogoPNG,
	}

	jsonData, _ := json.Marshal(farmData)
	req, _ := http.NewRequest("PUT", "/farm?id=1", bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	mockRepo.On("Update", mock.AnythingOfType("*models.Farm")).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestFarmHandler_UpdateFarm_InvalidMethod(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	farmService := service.NewFarmService(mockRepo)
	farmHandler := handlers.NewFarmHandler(farmService)

	req, _ := http.NewRequest("POST", "/farm?id=1", nil)
	w := httptest.NewRecorder()

	farmHandler.UpdateFarm(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestFarmHandler_UpdateFarm_MissingID(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	router, _ := setupFarmRouter(mockRepo)

	farmData := map[string]interface{}{
		"logo": tests.TestFileNewLogoPNG,
	}

	jsonData, _ := json.Marshal(farmData)
	req, _ := http.NewRequest("PUT", "/farm", bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFarmHandler_UpdateFarm_InvalidID(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	router, _ := setupFarmRouter(mockRepo)

	farmData := map[string]interface{}{
		"logo": tests.TestFileNewLogoPNG,
	}

	jsonData, _ := json.Marshal(farmData)
	req, _ := http.NewRequest("PUT", "/farm?id=invalid", bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFarmHandler_UpdateFarm_InvalidJSON(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	router, _ := setupFarmRouter(mockRepo)

	req, _ := http.NewRequest("PUT", "/farm?id=1", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFarmHandler_UpdateFarm_ServiceError(t *testing.T) {
	mockRepo := new(services.MockFarmRepository)
	router, _ := setupFarmRouter(mockRepo)

	farmData := map[string]interface{}{
		"logo": tests.TestFileNewLogoPNG,
	}

	jsonData, _ := json.Marshal(farmData)
	req, _ := http.NewRequest("PUT", "/farm?id=1", bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	mockRepo.On("Update", mock.AnythingOfType("*models.Farm")).Return(errors.New("erro ao atualizar"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}

