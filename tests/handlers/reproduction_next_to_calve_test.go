package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/fazendapro/FazendaPro-api/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReproductionRepository struct {
	mock.Mock
}

func (m *MockReproductionRepository) Create(reproduction *models.Reproduction) error {
	args := m.Called(reproduction)
	return args.Error(0)
}

func (m *MockReproductionRepository) FindByID(id uint) (*models.Reproduction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reproduction), args.Error(1)
}

func (m *MockReproductionRepository) FindByAnimalID(animalID uint) (*models.Reproduction, error) {
	args := m.Called(animalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reproduction), args.Error(1)
}

func (m *MockReproductionRepository) FindByFarmID(farmID uint) ([]models.Reproduction, error) {
	args := m.Called(farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Reproduction), args.Error(1)
}

func (m *MockReproductionRepository) FindByPhase(phase models.ReproductionPhase) ([]models.Reproduction, error) {
	args := m.Called(phase)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Reproduction), args.Error(1)
}

func (m *MockReproductionRepository) Update(reproduction *models.Reproduction) error {
	args := m.Called(reproduction)
	return args.Error(0)
}

func (m *MockReproductionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetNextToCalve(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)
	handler := handlers.NewReproductionHandler(reproductionService)

	farmID := uint(1)
	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200)
	expectedBirthDate := pregnancyDate.AddDate(0, 0, 283)

	mockReproductions := []models.Reproduction{
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                1,
				FarmID:            farmID,
				AnimalName:        tests.TestNameTataSalt,
				EarTagNumberLocal: 123,
				Photo:             tests.TestPathTataPNG,
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
				Photo:             tests.TestPathLaysPNG,
			},
		},
	}

	mockRepo.On("FindByPhase", models.PhasePrenhas).Return(mockReproductions, nil)

	req, err := http.NewRequest("GET", tests.EndpointAPIv1ReproductionsNextToCalve, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetNextToCalve(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Contains(t, response["message"], "Próximas vacas a parir encontradas com sucesso")

	data := response["data"].([]interface{})
	assert.Len(t, data, 2)

	firstAnimal := data[0].(map[string]interface{})
	assert.Equal(t, float64(1), firstAnimal["id"])
	assert.Equal(t, "Tata Salt", firstAnimal["animal_name"])
	assert.Equal(t, float64(123), firstAnimal["ear_tag_number_local"])
	assert.Equal(t, "src/assets/images/mocked/cows/tata.png", firstAnimal["photo"])
	assert.Equal(t, pregnancyDate.Format("2006-01-02"), firstAnimal["pregnancy_date"])
	assert.Equal(t, expectedBirthDate.Format("2006-01-02"), firstAnimal["expected_birth_date"])

	daysUntilBirth := int(expectedBirthDate.Sub(now).Hours() / 24)
	actualDaysUntilBirth := int(firstAnimal["days_until_birth"].(float64))
	assert.InDelta(t, float64(daysUntilBirth), float64(actualDaysUntilBirth), 1.0)

	var expectedStatus string
	if daysUntilBirth <= 30 {
		expectedStatus = "Alto"
	} else if daysUntilBirth <= 60 {
		expectedStatus = "Médio"
	} else {
		expectedStatus = "Baixo"
	}
	assert.Equal(t, expectedStatus, firstAnimal["status"])

	mockRepo.AssertExpectations(t)
}

func TestGetNextToCalve_MissingFarmID(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)
	handler := handlers.NewReproductionHandler(reproductionService)

	req, err := http.NewRequest("GET", "/api/v1/reproductions/next-to-calve", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetNextToCalve(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Equal(t, "ID da fazenda é obrigatório", response["message"])
}

func TestGetNextToCalve_InvalidFarmID(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)
	handler := handlers.NewReproductionHandler(reproductionService)

	req, err := http.NewRequest("GET", "/api/v1/reproductions/next-to-calve?farmId=invalid", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetNextToCalve(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Equal(t, "ID da fazenda inválido", response["message"])
}

func TestGetNextToCalve_ServiceError(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)
	handler := handlers.NewReproductionHandler(reproductionService)

	mockRepo.On("FindByPhase", models.PhasePrenhas).Return(nil, assert.AnError)

	req, err := http.NewRequest("GET", tests.EndpointAPIv1ReproductionsNextToCalve, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetNextToCalve(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.NotEmpty(t, response["message"])

	mockRepo.AssertExpectations(t)
}

func TestGetNextToCalve_EmptyResults(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)
	handler := handlers.NewReproductionHandler(reproductionService)

	mockRepo.On("FindByPhase", models.PhasePrenhas).Return([]models.Reproduction{}, nil)

	req, err := http.NewRequest("GET", tests.EndpointAPIv1ReproductionsNextToCalve, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetNextToCalve(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Contains(t, response["message"], "Próximas vacas a parir encontradas com sucesso (0 registros)")

	data, ok := response["data"].([]interface{})
	if !ok {
		data = []interface{}{}
	}
	assert.Len(t, data, 0)

	mockRepo.AssertExpectations(t)
}

func TestGetNextToCalve_FiltersByFarm(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)
	handler := handlers.NewReproductionHandler(reproductionService)

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
				FarmID:            1,
				AnimalName:        tests.TestNameTataSalt,
				EarTagNumberLocal: 123,
				Photo:             tests.TestPathTataPNG,
			},
		},
		{
			ID:            2,
			AnimalID:      2,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                2,
				FarmID:            2,
				AnimalName:        "Lays",
				EarTagNumberLocal: 124,
				Photo:             tests.TestPathLaysPNG,
			},
		},
	}

	mockRepo.On("FindByPhase", models.PhasePrenhas).Return(mockReproductions, nil)

	req, err := http.NewRequest("GET", tests.EndpointAPIv1ReproductionsNextToCalve, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetNextToCalve(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Len(t, data, 1)

	firstAnimal := data[0].(map[string]interface{})
	assert.Equal(t, "Tata Salt", firstAnimal["animal_name"])

	mockRepo.AssertExpectations(t)
}

func TestGetNextToCalve_SortsByDaysUntilBirth(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)
	handler := handlers.NewReproductionHandler(reproductionService)

	now := time.Now()
	pregnancyDate1 := now.AddDate(0, 0, -250)
	pregnancyDate2 := now.AddDate(0, 0, -220)
	pregnancyDate3 := now.AddDate(0, 0, -200)

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
				AnimalName:        tests.TestNameTataSalt,
				EarTagNumberLocal: 123,
				Photo:             tests.TestPathTataPNG,
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
				Photo:             tests.TestPathLaysPNG,
			},
		},
	}

	mockRepo.On("FindByPhase", models.PhasePrenhas).Return(mockReproductions, nil)

	req, err := http.NewRequest("GET", tests.EndpointAPIv1ReproductionsNextToCalve, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetNextToCalve(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))

	data := response["data"].([]interface{})
	assert.Len(t, data, 3)

	firstAnimal := data[0].(map[string]interface{})
	assert.Equal(t, "Tata Salt", firstAnimal["animal_name"])

	secondAnimal := data[1].(map[string]interface{})
	assert.Equal(t, "Lays", secondAnimal["animal_name"])

	thirdAnimal := data[2].(map[string]interface{})
	assert.Equal(t, "Matilda", thirdAnimal["animal_name"])

	mockRepo.AssertExpectations(t)
}
