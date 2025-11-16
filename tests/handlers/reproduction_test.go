package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/fazendapro/FazendaPro-api/tests/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupReproductionRouter(mockRepo *services.MockReproductionRepository) (*chi.Mux, *services.MockReproductionRepository) {
	reproductionService := service.NewReproductionService(mockRepo)
	reproductionHandler := handlers.NewReproductionHandler(reproductionService)
	r := chi.NewRouter()
	r.Post("/reproductions", reproductionHandler.CreateReproduction)
	r.Get("/reproductions", reproductionHandler.GetReproduction)
	r.Get("/reproductions/animal", reproductionHandler.GetReproductionByAnimal)
	r.Get("/reproductions/farm", reproductionHandler.GetReproductionsByFarm)
	r.Get("/reproductions/phase", reproductionHandler.GetReproductionsByPhase)
	r.Put("/reproductions", reproductionHandler.UpdateReproduction)
	r.Put("/reproductions/phase", reproductionHandler.UpdateReproductionPhase)
	r.Delete("/reproductions", reproductionHandler.DeleteReproduction)
	return r, mockRepo
}

func TestReproductionHandler_CreateReproduction_Success(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	reproductionData := map[string]interface{}{
		"animal_id":     1,
		"current_phase": 0,
	}

	jsonData, _ := json.Marshal(reproductionData)
	req, _ := http.NewRequest("POST", "/reproductions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.Reproduction")).Return(nil).Run(func(args mock.Arguments) {
		rep := args.Get(0).(*models.Reproduction)
		rep.ID = 1
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_CreateReproduction_InvalidMethod(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	reproductionService := service.NewReproductionService(mockRepo)
	reproductionHandler := handlers.NewReproductionHandler(reproductionService)

	req, _ := http.NewRequest("GET", "/reproductions", nil)
	w := httptest.NewRecorder()

	reproductionHandler.CreateReproduction(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestReproductionHandler_CreateReproduction_InvalidJSON(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("POST", "/reproductions", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReproductionHandler_CreateReproduction_ServiceError(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	reproductionData := map[string]interface{}{
		"animal_id": 1,
	}

	jsonData, _ := json.Marshal(reproductionData)
	req, _ := http.NewRequest("POST", "/reproductions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, errors.New("erro ao buscar"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_GetReproduction_Success(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	pregnancyDate := time.Now()
	expectedReproduction := &models.Reproduction{
		ID:            1,
		AnimalID:      1,
		CurrentPhase:  models.PhasePrenhas,
		PregnancyDate: &pregnancyDate,
	}

	req, _ := http.NewRequest("GET", "/reproductions?id=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(expectedReproduction, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_GetReproduction_MissingID(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/reproductions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReproductionHandler_GetReproduction_NotFound(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/reproductions?id=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_GetReproductionByAnimal_Success(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	pregnancyDate := time.Now()
	expectedReproduction := &models.Reproduction{
		ID:            1,
		AnimalID:      1,
		CurrentPhase:  models.PhasePrenhas,
		PregnancyDate: &pregnancyDate,
	}

	req, _ := http.NewRequest("GET", "/reproductions/animal?animalId=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByAnimalID", uint(1)).Return(expectedReproduction, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_GetReproductionByAnimal_MissingID(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/reproductions/animal", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReproductionHandler_GetReproductionsByFarm_Success(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	expectedReproductions := []models.Reproduction{
		{ID: 1, AnimalID: 1, CurrentPhase: models.PhasePrenhas},
		{ID: 2, AnimalID: 2, CurrentPhase: models.PhaseLactacao},
	}

	req, _ := http.NewRequest("GET", "/reproductions/farm?farmId=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByFarmID", uint(1)).Return(expectedReproductions, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_GetReproductionsByPhase_Success(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	expectedReproductions := []models.Reproduction{
		{ID: 1, AnimalID: 1, CurrentPhase: models.PhasePrenhas},
	}

	req, _ := http.NewRequest("GET", "/reproductions/phase?phase=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByPhase", models.ReproductionPhase(1)).Return(expectedReproductions, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_UpdateReproduction_Success(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	reproductionData := map[string]interface{}{
		"id":            1,
		"animal_id":     1,
		"current_phase": 1,
	}

	jsonData, _ := json.Marshal(reproductionData)
	req, _ := http.NewRequest("PUT", "/reproductions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(&models.Reproduction{ID: 1}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_UpdateReproductionPhase_Success(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	reproductionData := map[string]interface{}{
		"animal_id":       1,
		"new_phase":       2,
		"additional_data": map[string]interface{}{},
	}

	jsonData, _ := json.Marshal(reproductionData)
	req, _ := http.NewRequest("PUT", "/reproductions/phase", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("FindByAnimalID", uint(1)).Return(&models.Reproduction{ID: 1}, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_DeleteReproduction_Success(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("DELETE", "/reproductions?id=1", nil)
	w := httptest.NewRecorder()

	existingReproduction := &models.Reproduction{ID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_DeleteReproduction_MissingID(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("DELETE", "/reproductions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReproductionHandler_DeleteReproduction_NotFound(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("DELETE", "/reproductions?id=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}
