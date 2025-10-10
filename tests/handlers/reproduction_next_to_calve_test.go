package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetNextToCalve(t *testing.T) {
	// Setup
	mockRepo := &repository.MockReproductionRepository{}
	mockService := &service.MockReproductionService{}
	handler := NewReproductionHandler(mockService)

	// Mock data
	farmID := uint(1)
	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200)              // 200 dias atrás
	expectedBirthDate := pregnancyDate.AddDate(0, 0, 283) // 283 dias após a prenhez (período médio de gestação)

	mockReproductions := []models.Reproduction{
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                1,
				FarmID:            farmID,
				AnimalName:        "Tata Salt",
				EarTagNumberLocal: 123,
				Photo:             "src/assets/images/mocked/cows/tata.png",
			},
		},
		{
			ID:            2,
			AnimalID:      2,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                2,
				FarmID:            farmID,
				AnimalName:        "Lays",
				EarTagNumberLocal: 124,
				Photo:             "src/assets/images/mocked/cows/lays.png",
			},
		},
	}

	// Mock service call
	mockService.On("GetReproductionsByPhase", models.PhasePrenhas).Return(mockReproductions, nil)

	// Create request
	req, err := http.NewRequest("GET", "/api/v1/reproductions/next-to-calve?farmId=1", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetNextToCalve(rr, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Contains(t, response["message"], "Próximas vacas a parir encontradas com sucesso")

	// Verify data structure
	data := response["data"].([]interface{})
	assert.Len(t, data, 2)

	// Verify first animal
	firstAnimal := data[0].(map[string]interface{})
	assert.Equal(t, float64(1), firstAnimal["id"])
	assert.Equal(t, "Tata Salt", firstAnimal["animal_name"])
	assert.Equal(t, float64(123), firstAnimal["ear_tag_number_local"])
	assert.Equal(t, "src/assets/images/mocked/cows/tata.png", firstAnimal["photo"])
	assert.Equal(t, pregnancyDate.Format("2006-01-02"), firstAnimal["pregnancy_date"])
	assert.Equal(t, expectedBirthDate.Format("2006-01-02"), firstAnimal["expected_birth_date"])

	// Verify days until birth calculation
	daysUntilBirth := int(expectedBirthDate.Sub(now).Hours() / 24)
	assert.Equal(t, float64(daysUntilBirth), firstAnimal["days_until_birth"])

	// Verify status calculation
	var expectedStatus string
	if daysUntilBirth <= 30 {
		expectedStatus = "Alto"
	} else if daysUntilBirth <= 60 {
		expectedStatus = "Médio"
	} else {
		expectedStatus = "Baixo"
	}
	assert.Equal(t, expectedStatus, firstAnimal["status"])

	mockService.AssertExpectations(t)
}

func TestGetNextToCalve_MissingFarmID(t *testing.T) {
	// Setup
	mockService := &service.MockReproductionService{}
	handler := NewReproductionHandler(mockService)

	// Create request without farmId
	req, err := http.NewRequest("GET", "/api/v1/reproductions/next-to-calve", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetNextToCalve(rr, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Equal(t, "ID da fazenda é obrigatório", response["message"])
}

func TestGetNextToCalve_InvalidFarmID(t *testing.T) {
	// Setup
	mockService := &service.MockReproductionService{}
	handler := NewReproductionHandler(mockService)

	// Create request with invalid farmId
	req, err := http.NewRequest("GET", "/api/v1/reproductions/next-to-calve?farmId=invalid", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetNextToCalve(rr, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Equal(t, "ID da fazenda inválido", response["message"])
}

func TestGetNextToCalve_ServiceError(t *testing.T) {
	// Setup
	mockService := &service.MockReproductionService{}
	handler := NewReproductionHandler(mockService)

	// Mock service error
	mockService.On("GetReproductionsByPhase", models.PhasePrenhas).Return(nil, assert.AnError)

	// Create request
	req, err := http.NewRequest("GET", "/api/v1/reproductions/next-to-calve?farmId=1", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetNextToCalve(rr, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"], "Erro ao buscar registros de reprodução")

	mockService.AssertExpectations(t)
}

func TestGetNextToCalve_EmptyResults(t *testing.T) {
	// Setup
	mockService := &service.MockReproductionService{}
	handler := NewReproductionHandler(mockService)

	// Mock empty results
	mockService.On("GetReproductionsByPhase", models.PhasePrenhas).Return([]models.Reproduction{}, nil)

	// Create request
	req, err := http.NewRequest("GET", "/api/v1/reproductions/next-to-calve?farmId=1", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetNextToCalve(rr, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Contains(t, response["message"], "Próximas vacas a parir encontradas com sucesso (0 registros)")

	// Verify empty data
	data := response["data"].([]interface{})
	assert.Len(t, data, 0)

	mockService.AssertExpectations(t)
}

func TestGetNextToCalve_FiltersByFarm(t *testing.T) {
	// Setup
	mockService := &service.MockReproductionService{}
	handler := NewReproductionHandler(mockService)

	// Mock data with different farms
	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200)

	mockReproductions := []models.Reproduction{
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                1,
				FarmID:            1, // Target farm
				AnimalName:        "Tata Salt",
				EarTagNumberLocal: 123,
				Photo:             "src/assets/images/mocked/cows/tata.png",
			},
		},
		{
			ID:            2,
			AnimalID:      2,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                2,
				FarmID:            2, // Different farm
				AnimalName:        "Lays",
				EarTagNumberLocal: 124,
				Photo:             "src/assets/images/mocked/cows/lays.png",
			},
		},
	}

	// Mock service call
	mockService.On("GetReproductionsByPhase", models.PhasePrenhas).Return(mockReproductions, nil)

	// Create request for farm 1
	req, err := http.NewRequest("GET", "/api/v1/reproductions/next-to-calve?farmId=1", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetNextToCalve(rr, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))

	// Verify only farm 1 animals are returned
	data := response["data"].([]interface{})
	assert.Len(t, data, 1)

	firstAnimal := data[0].(map[string]interface{})
	assert.Equal(t, "Tata Salt", firstAnimal["animal_name"])

	mockService.AssertExpectations(t)
}

func TestGetNextToCalve_SortsByDaysUntilBirth(t *testing.T) {
	// Setup
	mockService := &service.MockReproductionService{}
	handler := NewReproductionHandler(mockService)

	// Mock data with different pregnancy dates
	now := time.Now()
	pregnancyDate1 := now.AddDate(0, 0, -250) // 30 days until birth
	pregnancyDate2 := now.AddDate(0, 0, -220) // 60 days until birth
	pregnancyDate3 := now.AddDate(0, 0, -200) // 80 days until birth

	mockReproductions := []models.Reproduction{
		{
			ID:            3,
			AnimalID:      3,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate3,
			Animal: models.Animal{
				ID:                3,
				FarmID:            1,
				AnimalName:        "Matilda",
				EarTagNumberLocal: 125,
				Photo:             "src/assets/images/mocked/cows/matilda.png",
			},
		},
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate1,
			Animal: models.Animal{
				ID:                1,
				FarmID:            1,
				AnimalName:        "Tata Salt",
				EarTagNumberLocal: 123,
				Photo:             "src/assets/images/mocked/cows/tata.png",
			},
		},
		{
			ID:            2,
			AnimalID:      2,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate2,
			Animal: models.Animal{
				ID:                2,
				FarmID:            1,
				AnimalName:        "Lays",
				EarTagNumberLocal: 124,
				Photo:             "src/assets/images/mocked/cows/lays.png",
			},
		},
	}

	// Mock service call
	mockService.On("GetReproductionsByPhase", models.PhasePrenhas).Return(mockReproductions, nil)

	// Create request
	req, err := http.NewRequest("GET", "/api/v1/reproductions/next-to-calve?farmId=1", nil)
	assert.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetNextToCalve(rr, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))

	// Verify sorting by days until birth (ascending)
	data := response["data"].([]interface{})
	assert.Len(t, data, 3)

	// First should be Tata Salt (30 days)
	firstAnimal := data[0].(map[string]interface{})
	assert.Equal(t, "Tata Salt", firstAnimal["animal_name"])

	// Second should be Lays (60 days)
	secondAnimal := data[1].(map[string]interface{})
	assert.Equal(t, "Lays", secondAnimal["animal_name"])

	// Third should be Matilda (80 days)
	thirdAnimal := data[2].(map[string]interface{})
	assert.Equal(t, "Matilda", thirdAnimal["animal_name"])

	mockService.AssertExpectations(t)
}
