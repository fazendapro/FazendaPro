package services

import (
	"errors"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWeightService_CreateOrUpdateWeight_Create(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weightDate := time.Now()
	weight := &models.Weight{
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.Weight")).Return(nil).Run(func(args mock.Arguments) {
		w := args.Get(0).(*models.Weight)
		w.ID = 1
	})

	err := weightService.CreateOrUpdateWeight(weight)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWeightService_CreateOrUpdateWeight_Update(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weightDate := time.Now()
	existingWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         weightDate.AddDate(0, -1, 0),
		AnimalWeight: 400.0,
	}

	newWeight := &models.Weight{
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(existingWeight, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Weight")).Return(nil)

	err := weightService.CreateOrUpdateWeight(newWeight)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWeightService_CreateOrUpdateWeight_InvalidAnimalID(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weight := &models.Weight{
		AnimalID:     0,
		Date:         time.Now(),
		AnimalWeight: 450.5,
	}

	err := weightService.CreateOrUpdateWeight(weight)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID do animal é obrigatório")
}

func TestWeightService_CreateOrUpdateWeight_InvalidWeight(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weight := &models.Weight{
		AnimalID:     1,
		Date:         time.Now(),
		AnimalWeight: 0,
	}

	err := weightService.CreateOrUpdateWeight(weight)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "peso do animal deve ser maior que zero")
}

func TestWeightService_CreateOrUpdateWeight_RepositoryError(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weight := &models.Weight{
		AnimalID:     1,
		Date:         time.Now(),
		AnimalWeight: 450.5,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, errors.New("database error"))

	err := weightService.CreateOrUpdateWeight(weight)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}

func TestWeightService_GetWeightByID(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weightDate := time.Now()
	expectedWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedWeight, nil)

	result, err := weightService.GetWeightByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, 450.5, result.AnimalWeight)
	mockRepo.AssertExpectations(t)
}

func TestWeightService_GetWeightByAnimalID(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weightDate := time.Now()
	expectedWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(expectedWeight, nil)

	result, err := weightService.GetWeightByAnimalID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.AnimalID)
	assert.Equal(t, 450.5, result.AnimalWeight)
	mockRepo.AssertExpectations(t)
}

func TestWeightService_GetWeightsByFarmID(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weightDate := time.Now()
	expectedWeights := []models.Weight{
		{
			ID:           1,
			AnimalID:     1,
			Date:         weightDate,
			AnimalWeight: 450.5,
		},
		{
			ID:           2,
			AnimalID:     2,
			Date:         weightDate,
			AnimalWeight: 500.0,
		},
	}

	mockRepo.On("FindByFarmID", uint(1)).Return(expectedWeights, nil)

	result, err := weightService.GetWeightsByFarmID(1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 450.5, result[0].AnimalWeight)
	assert.Equal(t, 500.0, result[1].AnimalWeight)
	mockRepo.AssertExpectations(t)
}

func TestWeightService_UpdateWeight_Success(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weightDate := time.Now()
	existingWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}

	updatedWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 480.0,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingWeight, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Weight")).Return(nil)

	err := weightService.UpdateWeight(updatedWeight)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWeightService_UpdateWeight_InvalidID(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weight := &models.Weight{
		ID:           0,
		AnimalID:     1,
		Date:         time.Now(),
		AnimalWeight: 450.5,
	}

	err := weightService.UpdateWeight(weight)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID do registro de peso é obrigatório")
}

func TestWeightService_UpdateWeight_NotFound(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         time.Now(),
		AnimalWeight: 450.5,
	}

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	err := weightService.UpdateWeight(weight)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "registro de peso não encontrado")
	mockRepo.AssertExpectations(t)
}

func TestWeightService_UpdateWeight_InvalidWeight(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weightDate := time.Now()
	existingWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}

	updatedWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 0,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingWeight, nil)

	err := weightService.UpdateWeight(updatedWeight)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "peso do animal deve ser maior que zero")
	mockRepo.AssertExpectations(t)
}

func TestWeightService_DeleteWeight_Success(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	weightDate := time.Now()
	existingWeight := &models.Weight{
		ID:           1,
		AnimalID:     1,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingWeight, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := weightService.DeleteWeight(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWeightService_DeleteWeight_NotFound(t *testing.T) {
	mockRepo := &MockWeightRepository{}
	weightService := service.NewWeightService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	err := weightService.DeleteWeight(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "registro de peso não encontrado")
	mockRepo.AssertExpectations(t)
}

