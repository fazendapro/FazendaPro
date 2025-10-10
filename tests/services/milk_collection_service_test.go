package services

import (
	"errors"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/fazendapro/FazendaPro-api/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMilkCollectionService_CreateMilkCollection_WithBatchUpdate(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	milkCollection := &models.MilkCollection{
		ID:       1,
		AnimalID: 1,
		Liters:   35.0,
		Date:     time.Now(),
	}

	animal := &models.Animal{
		ID:           1,
		CurrentBatch: models.Batch2,
	}

	milkCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: 1,
			Liters:   35.0,
			Date:     time.Now(),
		},
	}

	mockMilkRepo.On("Create", milkCollection).Return(nil)
	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)
	mockMilkRepo.On("FindByAnimalID", uint(1)).Return(milkCollections, nil)
	mockAnimalRepo.On("Update", mock.Anything).Return(nil)

	err := milkService.CreateMilkCollection(milkCollection)

	assert.NoError(t, err)
	mockMilkRepo.AssertExpectations(t)
	mockAnimalRepo.AssertExpectations(t)
}

func TestMilkCollectionService_CreateMilkCollection_WithoutBatchUpdate(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	milkCollection := &models.MilkCollection{
		ID:       1,
		AnimalID: 1,
		Liters:   25.0,
		Date:     time.Now(),
	}

	animal := &models.Animal{
		ID:           1,
		CurrentBatch: models.Batch2,
	}

	milkCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: 1,
			Liters:   25.0,
			Date:     time.Now(),
		},
	}

	mockMilkRepo.On("Create", milkCollection).Return(nil)
	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)
	mockMilkRepo.On("FindByAnimalID", uint(1)).Return(milkCollections, nil)

	err := milkService.CreateMilkCollection(milkCollection)

	assert.NoError(t, err)
	mockMilkRepo.AssertExpectations(t)
	mockAnimalRepo.AssertNotCalled(t, "Update")
}

func TestMilkCollectionService_CreateMilkCollection_CreateError(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	milkCollection := &models.MilkCollection{
		ID:       1,
		AnimalID: 1,
		Liters:   35.0,
		Date:     time.Now(),
	}

	mockMilkRepo.On("Create", milkCollection).Return(errors.New("database error"))

	err := milkService.CreateMilkCollection(milkCollection)

	assert.Error(t, err)
	mockMilkRepo.AssertExpectations(t)
	mockAnimalRepo.AssertNotCalled(t, "FindByID")
}

func TestMilkCollectionService_CreateMilkCollection_BatchUpdateError(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	milkCollection := &models.MilkCollection{
		ID:       1,
		AnimalID: 1,
		Liters:   35.0,
		Date:     time.Now(),
	}

	mockMilkRepo.On("Create", milkCollection).Return(nil)
	mockAnimalRepo.On("FindByID", uint(1)).Return((*models.Animal)(nil), errors.New("animal not found"))

	err := milkService.CreateMilkCollection(milkCollection)

	assert.Error(t, err)
	mockMilkRepo.AssertExpectations(t)
	mockAnimalRepo.AssertExpectations(t)
}
