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
	"github.com/fazendapro/FazendaPro-api/tests/mocks"
	"github.com/fazendapro/FazendaPro-api/tests/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMilkCollectionRouter(mockMilkRepo *mocks.MockMilkCollectionRepository, mockAnimalRepo *services.MockAnimalRepository) (*chi.Mux, *mocks.MockMilkCollectionRepository) {
	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)
	milkHandler := handlers.NewMilkCollectionHandler(milkService)
	r := chi.NewRouter()
	r.Post("/milk-collections", milkHandler.CreateMilkCollection)
	r.Put("/milk-collections/{id}", milkHandler.UpdateMilkCollection)
	r.Get("/milk-collections/farm/{farmId}", milkHandler.GetMilkCollectionsByFarmID)
	r.Get("/milk-collections/animal/{animalId}", milkHandler.GetMilkCollectionsByAnimalID)
	r.Get("/milk-collections/top-producers", milkHandler.GetTopMilkProducers)
	return r, mockMilkRepo
}

func TestMilkCollectionHandler_CreateMilkCollection_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	milkData := map[string]interface{}{
		"animal_id": 1,
		"liters":    35.5,
		"date":      "2024-01-15",
	}

	jsonData, _ := json.Marshal(milkData)
	req, _ := http.NewRequest("POST", "/milk-collections", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createdMilkCollection := &models.MilkCollection{
		ID:       1,
		AnimalID: 1,
		Liters:   35.5,
		Date:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	animal := &models.Animal{ID: 1, CurrentBatch: models.Batch1, FarmID: 1}
	milkCollections := []models.MilkCollection{*createdMilkCollection}

	mockMilkRepo.On("Create", mock.AnythingOfType("*models.MilkCollection")).Return(nil).Run(func(args mock.Arguments) {
		mc := args.Get(0).(*models.MilkCollection)
		mc.ID = 1
	})
	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)
	mockMilkRepo.On("FindByAnimalID", uint(1)).Return(milkCollections, nil)
	mockAnimalRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(nil).Maybe()
	mockMilkRepo.On("FindByID", uint(1)).Return(createdMilkCollection, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response handlers.MilkCollectionResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	mockMilkRepo.AssertExpectations(t)
	mockAnimalRepo.AssertExpectations(t)
}

func TestMilkCollectionHandler_CreateMilkCollection_InvalidJSON(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	req, _ := http.NewRequest("POST", "/milk-collections", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMilkCollectionHandler_CreateMilkCollection_InvalidDate(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	milkData := map[string]interface{}{
		"animal_id": 1,
		"liters":    35.5,
		"date":      "invalid-date",
	}

	jsonData, _ := json.Marshal(milkData)
	req, _ := http.NewRequest("POST", "/milk-collections", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMilkCollectionHandler_CreateMilkCollection_ServiceError(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	milkData := map[string]interface{}{
		"animal_id": 1,
		"liters":    35.5,
		"date":      "2024-01-15",
	}

	jsonData, _ := json.Marshal(milkData)
	req, _ := http.NewRequest("POST", "/milk-collections", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockMilkRepo.On("Create", mock.AnythingOfType("*models.MilkCollection")).Return(errors.New("erro ao criar"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionHandler_UpdateMilkCollection_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	milkData := map[string]interface{}{
		"animal_id": 1,
		"liters":    40.0,
		"date":      "2024-01-15",
	}

	jsonData, _ := json.Marshal(milkData)
	req, _ := http.NewRequest("PUT", "/milk-collections/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	updatedMilkCollection := &models.MilkCollection{
		ID:       1,
		AnimalID: 1,
		Liters:   40.0,
		Date:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	mockMilkRepo.On("Update", mock.AnythingOfType("*models.MilkCollection")).Return(nil)
	mockMilkRepo.On("FindByID", uint(1)).Return(updatedMilkCollection, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response handlers.MilkCollectionResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionHandler_UpdateMilkCollection_InvalidID(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	milkData := map[string]interface{}{
		"animal_id": 1,
		"liters":    40.0,
		"date":      "2024-01-15",
	}

	jsonData, _ := json.Marshal(milkData)
	req, _ := http.NewRequest("PUT", "/milk-collections/invalid", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMilkCollectionHandler_GetMilkCollectionsByFarmID_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	expectedCollections := []models.MilkCollection{
		{ID: 1, AnimalID: 1, Liters: 35.5, Date: time.Now()},
		{ID: 2, AnimalID: 2, Liters: 40.0, Date: time.Now()},
	}

	req, _ := http.NewRequest("GET", "/milk-collections/farm/1", nil)
	w := httptest.NewRecorder()

	mockMilkRepo.On("FindByFarmID", uint(1)).Return(expectedCollections, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response handlers.MilkCollectionsResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 2)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionHandler_GetMilkCollectionsByAnimalID_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	expectedCollections := []models.MilkCollection{
		{ID: 1, AnimalID: 1, Liters: 35.5, Date: time.Now()},
		{ID: 2, AnimalID: 1, Liters: 40.0, Date: time.Now()},
	}

	req, _ := http.NewRequest("GET", "/milk-collections/animal/1", nil)
	w := httptest.NewRecorder()

	mockMilkRepo.On("FindByAnimalID", uint(1)).Return(expectedCollections, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response handlers.MilkCollectionsResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 2)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionHandler_GetTopMilkProducers_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	expectedCollections := []models.MilkCollection{
		{ID: 1, AnimalID: 1, Liters: 50.0, Date: time.Now()},
		{ID: 2, AnimalID: 2, Liters: 45.0, Date: time.Now()},
	}

	req, _ := http.NewRequest("GET", "/milk-collections/top-producers?farmId=1", nil)
	w := httptest.NewRecorder()

	mockMilkRepo.On("FindByFarmIDWithDateRange", uint(1), mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time")).Return(expectedCollections, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionHandler_GetMilkCollectionsByFarmID_WithDateRange(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	expectedCollections := []models.MilkCollection{
		{ID: 1, AnimalID: 1, Liters: 35.5, Date: time.Now()},
	}

	req, _ := http.NewRequest("GET", "/milk-collections/farm/1?start_date=2024-01-01&end_date=2024-01-31", nil)
	w := httptest.NewRecorder()

	mockMilkRepo.On("FindByFarmIDWithDateRange", uint(1), mock.AnythingOfType("*time.Time"), mock.AnythingOfType("*time.Time")).Return(expectedCollections, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionHandler_GetMilkCollectionsByFarmID_InvalidFarmID(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	req, _ := http.NewRequest("GET", "/milk-collections/farm/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMilkCollectionHandler_GetMilkCollectionsByAnimalID_InvalidAnimalID(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	req, _ := http.NewRequest("GET", "/milk-collections/animal/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMilkCollectionHandler_GetTopMilkProducers_MissingFarmID(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	req, _ := http.NewRequest("GET", "/milk-collections/top-producers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMilkCollectionHandler_GetTopMilkProducers_InvalidFarmID(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	req, _ := http.NewRequest("GET", "/milk-collections/top-producers?farmId=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMilkCollectionHandler_GetMilkCollectionsByFarmID_WithAnimalData(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	birthDate := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: 1,
			Liters:   35.5,
			Date:     time.Now(),
			Animal: models.Animal{
				ID:                1,
				FarmID:            1,
				AnimalName:        "Vaca Teste",
				EarTagNumberLocal: 123,
				BirthDate:         &birthDate,
			},
		},
	}

	req, _ := http.NewRequest("GET", "/milk-collections/farm/1", nil)
	w := httptest.NewRecorder()

	mockMilkRepo.On("FindByFarmID", uint(1)).Return(expectedCollections, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response handlers.MilkCollectionsResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 1)
	assert.NotEmpty(t, response.Data[0].Animal.BirthDate)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionHandler_GetMilkCollectionsByFarmID_WithNilBirthDate(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(services.MockAnimalRepository)
	router, _ := setupMilkCollectionRouter(mockMilkRepo, mockAnimalRepo)

	expectedCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: 1,
			Liters:   35.5,
			Date:     time.Now(),
			Animal: models.Animal{
				ID:                1,
				FarmID:            1,
				AnimalName:        "Vaca Teste",
				EarTagNumberLocal: 123,
				BirthDate:         nil,
			},
		},
	}

	req, _ := http.NewRequest("GET", "/milk-collections/farm/1", nil)
	w := httptest.NewRecorder()

	mockMilkRepo.On("FindByFarmID", uint(1)).Return(expectedCollections, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response handlers.MilkCollectionsResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Len(t, response.Data, 1)
	assert.Empty(t, response.Data[0].Animal.BirthDate)
	mockMilkRepo.AssertExpectations(t)
}
