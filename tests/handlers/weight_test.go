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
	"github.com/fazendapro/FazendaPro-api/tests"
	"github.com/fazendapro/FazendaPro-api/tests/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupWeightRouter(mockRepo *services.MockWeightRepository) (*chi.Mux, *services.MockWeightRepository) {
	weightService := service.NewWeightService(mockRepo)
	weightHandler := handlers.NewWeightHandler(weightService)
	r := chi.NewRouter()
	r.Route("/api/v1/weights", func(r chi.Router) {
		r.Post("/", weightHandler.CreateOrUpdateWeight)
		r.Put("/", weightHandler.UpdateWeight)
		r.Get("/farm/{farmId}", weightHandler.GetWeightsByFarm)
		r.Get("/animal/{animalId}", weightHandler.GetWeightByAnimal)
	})
	return r, mockRepo
}

func TestWeightHandler_CreateOrUpdateWeight_Create(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	router, _ := setupWeightRouter(mockRepo)

	weightData := map[string]interface{}{
		"animal_id":     1,
		"date":          "2024-01-15",
		"animal_weight": 450.5,
	}

	jsonData, _ := json.Marshal(weightData)
	req, _ := http.NewRequest("POST", "/api/v1/weights", bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	createdWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		AnimalWeight: 450.5,
		Date:         time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		Animal: models.Animal{
			ID:                1,
			AnimalName:        "Boi João",
			EarTagNumberLocal: 123,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Primeira chamada: handler verifica se existe (retorna nil)
	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, nil).Once()
	// Segunda chamada: service verifica se existe dentro de CreateOrUpdateWeight (retorna nil)
	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, nil).Once()
	// Terceira chamada: criar o peso
	mockRepo.On("Create", mock.AnythingOfType(tests.TypeModelsWeight)).Return(nil).Run(func(args mock.Arguments) {
		w := args.Get(0).(*models.Weight)
		w.ID = 1
		w.CreatedAt = time.Now()
		w.UpdatedAt = time.Now()
		w.Animal = models.Animal{
			ID:                1,
			AnimalName:        "Boi João",
			EarTagNumberLocal: 123,
		}
	})
	// Quarta chamada: handler recupera após criar
	mockRepo.On("FindByAnimalID", uint(1)).Return(createdWeight, nil).Once()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	if response != nil {
		assert.True(t, response["success"].(bool))
	}
	mockRepo.AssertExpectations(t)
}

func TestWeightHandler_CreateOrUpdateWeight_Update(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	router, _ := setupWeightRouter(mockRepo)

	weightDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	existingWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 400.0,
		Animal: models.Animal{
			ID:                1,
			AnimalName:        "Boi João",
			EarTagNumberLocal: 123,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	weightData := map[string]interface{}{
		"animal_id":     1,
		"date":          "2024-01-15",
		"animal_weight": 450.5,
	}

	jsonData, _ := json.Marshal(weightData)
	req, _ := http.NewRequest("POST", "/api/v1/weights", bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	mockRepo.On("FindByAnimalID", uint(1)).Return(existingWeight, nil)
	mockRepo.On("Update", mock.AnythingOfType(tests.TypeModelsWeight)).Return(nil)
	mockRepo.On("FindByAnimalID", uint(1)).Return(&models.Weight{
		ID:           1,
		AnimalID:     1,
		AnimalWeight: 450.5,
		Date:         weightDate,
		Animal: models.Animal{
			ID:                1,
			AnimalName:        "Boi João",
			EarTagNumberLocal: 123,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestWeightHandler_CreateOrUpdateWeight_InvalidMethod(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	weightService := service.NewWeightService(mockRepo)
	weightHandler := handlers.NewWeightHandler(weightService)

	req, _ := http.NewRequest("GET", "/api/v1/weights", nil)
	w := httptest.NewRecorder()

	weightHandler.CreateOrUpdateWeight(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestWeightHandler_CreateOrUpdateWeight_InvalidJSON(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	router, _ := setupWeightRouter(mockRepo)

	req, _ := http.NewRequest("POST", "/api/v1/weights", bytes.NewBuffer([]byte(tests.TestErrorInvalidJSON)))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWeightHandler_CreateOrUpdateWeight_InvalidDate(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	router, _ := setupWeightRouter(mockRepo)

	weightData := map[string]interface{}{
		"animal_id":     1,
		"date":          "invalid-date",
		"animal_weight": 450.5,
	}

	jsonData, _ := json.Marshal(weightData)
	req, _ := http.NewRequest("POST", "/api/v1/weights", bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWeightHandler_CreateOrUpdateWeight_ServiceError(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	router, _ := setupWeightRouter(mockRepo)

	weightData := map[string]interface{}{
		"animal_id":     1,
		"date":          "2024-01-15",
		"animal_weight": 450.5,
	}

	jsonData, _ := json.Marshal(weightData)
	req, _ := http.NewRequest("POST", "/api/v1/weights", bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, errors.New(tests.TestErrorDatabaseError))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestWeightHandler_GetWeightByAnimal_Success(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	router, _ := setupWeightRouter(mockRepo)

	weightDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 450.5,
		Animal: models.Animal{
			ID:                1,
			AnimalName:        "Boi João",
			EarTagNumberLocal: 123,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req, _ := http.NewRequest("GET", "/api/v1/weights/animal/1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByAnimalID", uint(1)).Return(expectedWeight, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestWeightHandler_GetWeightByAnimal_NotFound(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	router, _ := setupWeightRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/api/v1/weights/animal/1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestWeightHandler_GetWeightsByFarm_Success(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	router, _ := setupWeightRouter(mockRepo)

	weightDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedWeights := []models.Weight{
		{
			ID:           1,
			AnimalID:     1,
			Date:         weightDate,
			AnimalWeight: 450.5,
			Animal: models.Animal{
				ID:                1,
				AnimalName:        "Boi João",
				EarTagNumberLocal: 123,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:           2,
			AnimalID:     2,
			Date:         weightDate,
			AnimalWeight: 500.0,
			Animal: models.Animal{
				ID:                2,
				AnimalName:        "Boi Pedro",
				EarTagNumberLocal: 124,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	req, _ := http.NewRequest("GET", "/api/v1/weights/farm/1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByFarmID", uint(1)).Return(expectedWeights, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestWeightHandler_UpdateWeight_Success(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	router, _ := setupWeightRouter(mockRepo)

	weightDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	existingWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 450.5,
		Animal: models.Animal{
			ID:                1,
			AnimalName:        "Boi João",
			EarTagNumberLocal: 123,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	weightData := map[string]interface{}{
		"id":            1,
		"animal_id":     1,
		"date":          "2024-01-15",
		"animal_weight": 480.0,
	}

	jsonData, _ := json.Marshal(weightData)
	req, _ := http.NewRequest("PUT", "/api/v1/weights", bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(existingWeight, nil)
	mockRepo.On("Update", mock.AnythingOfType(tests.TypeModelsWeight)).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestWeightHandler_UpdateWeight_InvalidMethod(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	weightService := service.NewWeightService(mockRepo)
	weightHandler := handlers.NewWeightHandler(weightService)

	req, _ := http.NewRequest("GET", "/api/v1/weights", nil)
	w := httptest.NewRecorder()

	weightHandler.UpdateWeight(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestWeightHandler_UpdateWeight_NotFound(t *testing.T) {
	mockRepo := new(services.MockWeightRepository)
	router, _ := setupWeightRouter(mockRepo)

	weightData := map[string]interface{}{
		"id":            1,
		"animal_id":     1,
		"date":          "2024-01-15",
		"animal_weight": 480.0,
	}

	jsonData, _ := json.Marshal(weightData)
	req, _ := http.NewRequest("PUT", "/api/v1/weights", bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}
