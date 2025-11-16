package services

import (
	"errors"
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestFarmService_GetFarmByID_Success(t *testing.T) {
	mockRepo := new(MockFarmRepository)
	farmService := service.NewFarmService(mockRepo)

	expectedFarm := &models.Farm{
		ID:        1,
		CompanyID: 1,
		Logo:      "logo.png",
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedFarm, nil)

	farm, err := farmService.GetFarmByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, farm)
	assert.Equal(t, uint(1), farm.ID)
	mockRepo.AssertExpectations(t)
}

func TestFarmService_GetFarmByID_InvalidID(t *testing.T) {
	mockRepo := new(MockFarmRepository)
	farmService := service.NewFarmService(mockRepo)

	farm, err := farmService.GetFarmByID(0)

	assert.Error(t, err)
	assert.Nil(t, farm)
	assert.Contains(t, err.Error(), "farm ID is required")
}

func TestFarmService_GetFarmByID_NotFound(t *testing.T) {
	mockRepo := new(MockFarmRepository)
	farmService := service.NewFarmService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

	farm, err := farmService.GetFarmByID(1)

	assert.Error(t, err)
	assert.Nil(t, farm)
	mockRepo.AssertExpectations(t)
}

func TestFarmService_UpdateFarm_Success(t *testing.T) {
	mockRepo := new(MockFarmRepository)
	farmService := service.NewFarmService(mockRepo)

	farm := &models.Farm{
		ID:   1,
		Logo: "new-logo.png",
	}

	mockRepo.On("Update", farm).Return(nil)

	err := farmService.UpdateFarm(farm)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFarmService_UpdateFarm_InvalidID(t *testing.T) {
	mockRepo := new(MockFarmRepository)
	farmService := service.NewFarmService(mockRepo)

	farm := &models.Farm{
		ID:   0,
		Logo: "new-logo.png",
	}

	err := farmService.UpdateFarm(farm)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "farm ID is required")
}

func TestFarmService_LoadCompanyData_Success(t *testing.T) {
	mockRepo := new(MockFarmRepository)
	farmService := service.NewFarmService(mockRepo)

	farm := &models.Farm{
		ID:        1,
		CompanyID: 1,
	}

	mockRepo.On("LoadCompanyData", farm).Return(nil)

	err := farmService.LoadCompanyData(farm)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFarmService_LoadCompanyData_Error(t *testing.T) {
	mockRepo := new(MockFarmRepository)
	farmService := service.NewFarmService(mockRepo)

	farm := &models.Farm{
		ID:        1,
		CompanyID: 1,
	}

	mockRepo.On("LoadCompanyData", farm).Return(errors.New("erro ao carregar"))

	err := farmService.LoadCompanyData(farm)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

