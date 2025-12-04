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

func setupReproductionRouter(mockRepo *services.MockReproductionRepository) (*chi.Mux, *services.MockReproductionRepository) {
	reproductionService := service.NewReproductionService(mockRepo)
	reproductionHandler := handlers.NewReproductionHandler(reproductionService)
	r := chi.NewRouter()
	r.Post(tests.EndpointReproductions, reproductionHandler.CreateReproduction)
	r.Get(tests.EndpointReproductions, reproductionHandler.GetReproduction)
	r.Get(tests.EndpointReproductions+"/animal", reproductionHandler.GetReproductionByAnimal)
	r.Get(tests.EndpointReproductions+"/farm", reproductionHandler.GetReproductionsByFarm)
	r.Get(tests.EndpointReproductionsPhase, reproductionHandler.GetReproductionsByPhase)
	r.Put(tests.EndpointReproductions, reproductionHandler.UpdateReproduction)
	r.Put(tests.EndpointReproductionsPhase, reproductionHandler.UpdateReproductionPhase)
	r.Delete(tests.EndpointReproductions, reproductionHandler.DeleteReproduction)
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
	req, _ := http.NewRequest("POST", tests.EndpointReproductions, bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType(tests.TypeModelsReproduction)).Return(nil).Run(func(args mock.Arguments) {
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

	req, _ := http.NewRequest("GET", tests.EndpointReproductions, nil)
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
	req, _ := http.NewRequest("POST", tests.EndpointReproductions, bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
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

	req, _ := http.NewRequest("GET", tests.EndpointReproductionsWithID, nil)
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

	req, _ := http.NewRequest("GET", tests.EndpointReproductions, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReproductionHandler_GetReproduction_NotFound(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("GET", tests.EndpointReproductionsWithID, nil)
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

func TestReproductionHandler_GetReproductionsByFarm_WithPagination(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	expectedReproductions := []models.Reproduction{
		{ID: 1, AnimalID: 1, CurrentPhase: models.PhasePrenhas},
		{ID: 2, AnimalID: 2, CurrentPhase: models.PhaseLactacao},
	}

	req, _ := http.NewRequest("GET", "/reproductions/farm?farmId=1&page=1&limit=10", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByFarmIDWithPagination", uint(1), 1, 10).Return(expectedReproductions, int64(2), nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Success bool `json:"success"`
		Data    struct {
			Reproductions []handlers.ReproductionResponse `json:"reproductions"`
			Total         int64                           `json:"total"`
			Page          int                             `json:"page"`
			Limit         int                             `json:"limit"`
		} `json:"data"`
		Message string `json:"message"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, int64(2), response.Data.Total)
	assert.Equal(t, 1, response.Data.Page)
	assert.Equal(t, 10, response.Data.Limit)
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

func TestReproductionHandler_CreateReproduction_WithAllDates(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	inseminationDate := "2024-01-01"
	pregnancyDate := "2024-01-15"
	expectedBirthDate := "2024-10-24"
	actualBirthDate := "2024-10-20"
	lactationStartDate := "2024-10-21"
	lactationEndDate := "2024-12-01"
	dryPeriodStartDate := "2024-12-02"

	reproductionData := map[string]interface{}{
		"animal_id":               1,
		"current_phase":           3,
		"insemination_date":       inseminationDate,
		"pregnancy_date":          pregnancyDate,
		"expected_birth_date":     expectedBirthDate,
		"actual_birth_date":       actualBirthDate,
		"lactation_start_date":    lactationStartDate,
		"lactation_end_date":      lactationEndDate,
		"dry_period_start_date":   dryPeriodStartDate,
		"insemination_type":       "IA",
		"veterinary_confirmation": true,
		"observations":            "Teste completo",
	}

	jsonData, _ := json.Marshal(reproductionData)
	req, _ := http.NewRequest("POST", tests.EndpointReproductions, bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType(tests.TypeModelsReproduction)).Return(nil).Run(func(args mock.Arguments) {
		rep := args.Get(0).(*models.Reproduction)
		rep.ID = 1
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_CreateReproduction_WithInvalidDates(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	reproductionData := map[string]interface{}{
		"animal_id":         1,
		"current_phase":     3,
		"pregnancy_date":    "invalid-date",
		"insemination_date": "also-invalid",
	}

	jsonData, _ := json.Marshal(reproductionData)
	req, _ := http.NewRequest("POST", tests.EndpointReproductions, bytes.NewBuffer(jsonData))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType(tests.TypeModelsReproduction)).Return(nil).Run(func(args mock.Arguments) {
		rep := args.Get(0).(*models.Reproduction)
		rep.ID = 1
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestReproductionHandler_GetReproduction_InvalidID(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/reproductions?id=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReproductionHandler_GetReproductionByAnimal_InvalidAnimalID(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/reproductions/animal?animalId=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReproductionHandler_GetReproductionsByFarm_InvalidFarmID(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/reproductions/farm?farmId=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReproductionHandler_GetReproductionsByPhase_InvalidPhase(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/reproductions/phase?phase=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReproductionHandler_GetReproduction_WithAllDates(t *testing.T) {
	mockRepo := new(services.MockReproductionRepository)
	router, _ := setupReproductionRouter(mockRepo)

	now := time.Now()
	inseminationDate := now.AddDate(0, -9, 0)
	pregnancyDate := now.AddDate(0, -8, 0)
	expectedBirthDate := now.AddDate(0, 1, 0)
	actualBirthDate := now.AddDate(0, 0, -5)
	lactationStartDate := now.AddDate(0, 0, -4)
	lactationEndDate := now.AddDate(0, 0, -1)
	dryPeriodStartDate := now

	reproduction := &models.Reproduction{
		ID:                     1,
		AnimalID:               1,
		CurrentPhase:           models.PhaseLactacao,
		InseminationDate:       &inseminationDate,
		PregnancyDate:          &pregnancyDate,
		ExpectedBirthDate:      &expectedBirthDate,
		ActualBirthDate:        &actualBirthDate,
		LactationStartDate:     &lactationStartDate,
		LactationEndDate:       &lactationEndDate,
		DryPeriodStartDate:     &dryPeriodStartDate,
		InseminationType:       "IA",
		VeterinaryConfirmation: true,
		Observations:           "Teste completo",
		Animal: models.Animal{
			ID:                1,
			FarmID:            1,
			AnimalName:        "Vaca Teste",
			EarTagNumberLocal: 123,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	req, _ := http.NewRequest("GET", tests.EndpointReproductionsWithID, nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(reproduction, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["insemination_date"])
	assert.NotNil(t, data["pregnancy_date"])
	assert.NotNil(t, data["expected_birth_date"])
	assert.NotNil(t, data["actual_birth_date"])
	assert.NotNil(t, data["lactation_start_date"])
	assert.NotNil(t, data["lactation_end_date"])
	assert.NotNil(t, data["dry_period_start_date"])
	mockRepo.AssertExpectations(t)
}
