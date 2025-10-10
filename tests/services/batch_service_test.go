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

func TestBatchService_UpdateAnimalBatch_Batch1(t *testing.T) {
	mockAnimalRepo := new(mocks.MockAnimalRepository)
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)

	animalID := uint(1)
	animal := &models.Animal{
		ID:           animalID,
		CurrentBatch: models.Batch2,
	}

	milkCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: animalID,
			Liters:   35.0,
			Date:     time.Now().Add(-24 * time.Hour),
		},
		{
			ID:       2,
			AnimalID: animalID,
			Liters:   40.0,
			Date:     time.Now(),
		},
	}

	mockAnimalRepo.On("FindByID", animalID).Return(animal, nil)
	mockMilkRepo.On("FindByAnimalID", animalID).Return(milkCollections, nil)
	mockAnimalRepo.On("Update", mock.Anything).Return(nil)

	err := batchService.UpdateAnimalBatch(animalID)

	assert.NoError(t, err)
	mockAnimalRepo.AssertExpectations(t)
	mockMilkRepo.AssertExpectations(t)
}

func TestBatchService_UpdateAnimalBatch_Batch2(t *testing.T) {
	mockAnimalRepo := new(mocks.MockAnimalRepository)
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)

	animalID := uint(1)
	animal := &models.Animal{
		ID:           animalID,
		CurrentBatch: models.Batch1,
	}

	milkCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: animalID,
			Liters:   25.0,
			Date:     time.Now(),
		},
	}

	mockAnimalRepo.On("FindByID", animalID).Return(animal, nil)
	mockMilkRepo.On("FindByAnimalID", animalID).Return(milkCollections, nil)
	mockAnimalRepo.On("Update", mock.Anything).Return(nil)

	err := batchService.UpdateAnimalBatch(animalID)

	assert.NoError(t, err)
	mockAnimalRepo.AssertExpectations(t)
	mockMilkRepo.AssertExpectations(t)
}

func TestBatchService_UpdateAnimalBatch_Batch3(t *testing.T) {
	mockAnimalRepo := new(mocks.MockAnimalRepository)
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)

	animalID := uint(1)
	animal := &models.Animal{
		ID:           animalID,
		CurrentBatch: models.Batch1,
	}

	milkCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: animalID,
			Liters:   15.0,
			Date:     time.Now(),
		},
	}

	mockAnimalRepo.On("FindByID", animalID).Return(animal, nil)
	mockMilkRepo.On("FindByAnimalID", animalID).Return(milkCollections, nil)
	mockAnimalRepo.On("Update", mock.Anything).Return(nil)

	err := batchService.UpdateAnimalBatch(animalID)

	assert.NoError(t, err)
	mockAnimalRepo.AssertExpectations(t)
	mockMilkRepo.AssertExpectations(t)
}

func TestBatchService_UpdateAnimalBatch_NoChange(t *testing.T) {
	mockAnimalRepo := new(mocks.MockAnimalRepository)
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)

	animalID := uint(1)
	animal := &models.Animal{
		ID:           animalID,
		CurrentBatch: models.Batch1,
	}

	milkCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: animalID,
			Liters:   35.0,
			Date:     time.Now(),
		},
	}

	mockAnimalRepo.On("FindByID", animalID).Return(animal, nil)
	mockMilkRepo.On("FindByAnimalID", animalID).Return(milkCollections, nil)

	err := batchService.UpdateAnimalBatch(animalID)

	assert.NoError(t, err)
	mockAnimalRepo.AssertNotCalled(t, "Update")
	mockAnimalRepo.AssertExpectations(t)
	mockMilkRepo.AssertExpectations(t)
}

func TestBatchService_UpdateAnimalBatch_NoMilkCollections(t *testing.T) {
	mockAnimalRepo := new(mocks.MockAnimalRepository)
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)

	animalID := uint(1)
	animal := &models.Animal{
		ID:           animalID,
		CurrentBatch: models.Batch1,
	}

	mockAnimalRepo.On("FindByID", animalID).Return(animal, nil)
	mockMilkRepo.On("FindByAnimalID", animalID).Return([]models.MilkCollection{}, nil)

	err := batchService.UpdateAnimalBatch(animalID)

	assert.NoError(t, err)
	mockAnimalRepo.AssertNotCalled(t, "Update")
	mockAnimalRepo.AssertExpectations(t)
	mockMilkRepo.AssertExpectations(t)
}

func TestBatchService_UpdateAnimalBatch_AnimalNotFound(t *testing.T) {
	mockAnimalRepo := new(mocks.MockAnimalRepository)
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)

	animalID := uint(1)

	mockAnimalRepo.On("FindByID", animalID).Return((*models.Animal)(nil), errors.New("animal not found"))

	err := batchService.UpdateAnimalBatch(animalID)

	assert.Error(t, err)
	mockAnimalRepo.AssertExpectations(t)
}

func TestBatchService_UpdateAnimalBatch_MilkCollectionError(t *testing.T) {
	mockAnimalRepo := new(mocks.MockAnimalRepository)
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)

	animalID := uint(1)
	animal := &models.Animal{
		ID:           animalID,
		CurrentBatch: models.Batch1,
	}

	mockAnimalRepo.On("FindByID", animalID).Return(animal, nil)
	mockMilkRepo.On("FindByAnimalID", animalID).Return([]models.MilkCollection{}, errors.New("database error"))

	err := batchService.UpdateAnimalBatch(animalID)

	assert.Error(t, err)
	mockAnimalRepo.AssertExpectations(t)
	mockMilkRepo.AssertExpectations(t)
}
