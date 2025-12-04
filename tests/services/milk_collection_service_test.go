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

func TestMilkCollectionService_GetMilkCollectionByID_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	expectedMilkCollection := &models.MilkCollection{
		ID:       1,
		AnimalID: 1,
		Liters:   35.0,
		Date:     time.Now(),
	}

	mockMilkRepo.On("FindByID", uint(1)).Return(expectedMilkCollection, nil)

	result, err := milkService.GetMilkCollectionByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedMilkCollection, result)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_GetMilkCollectionByID_NotFound(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	mockMilkRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

	result, err := milkService.GetMilkCollectionByID(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_GetMilkCollectionsByFarmID_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	expectedCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: 1,
			Liters:   35.0,
			Date:     time.Now(),
		},
		{
			ID:       2,
			AnimalID: 2,
			Liters:   30.0,
			Date:     time.Now(),
		},
	}

	mockMilkRepo.On("FindByFarmID", uint(1)).Return(expectedCollections, nil)

	result, err := milkService.GetMilkCollectionsByFarmID(1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_GetMilkCollectionsByFarmIDWithDateRange_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	expectedCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: 1,
			Liters:   35.0,
			Date:     time.Now(),
		},
	}

	mockMilkRepo.On("FindByFarmIDWithDateRange", uint(1), &startDate, &endDate).Return(expectedCollections, nil)

	result, err := milkService.GetMilkCollectionsByFarmIDWithDateRange(1, &startDate, &endDate)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_GetMilkCollectionsByAnimalID_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	expectedCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: 1,
			Liters:   35.0,
			Date:     time.Now(),
		},
	}

	mockMilkRepo.On("FindByAnimalID", uint(1)).Return(expectedCollections, nil)

	result, err := milkService.GetMilkCollectionsByAnimalID(1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_UpdateMilkCollection_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	milkCollection := &models.MilkCollection{
		ID:       1,
		AnimalID: 1,
		Liters:   40.0,
		Date:     time.Now(),
	}

	mockMilkRepo.On("Update", milkCollection).Return(nil)

	err := milkService.UpdateMilkCollection(milkCollection)

	assert.NoError(t, err)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_UpdateMilkCollection_Error(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	milkCollection := &models.MilkCollection{
		ID:       1,
		AnimalID: 1,
		Liters:   40.0,
		Date:     time.Now(),
	}

	mockMilkRepo.On("Update", milkCollection).Return(errors.New("database error"))

	err := milkService.UpdateMilkCollection(milkCollection)

	assert.Error(t, err)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_DeleteMilkCollection_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	mockMilkRepo.On("Delete", uint(1)).Return(nil)

	err := milkService.DeleteMilkCollection(1)

	assert.NoError(t, err)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_DeleteMilkCollection_Error(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	mockMilkRepo.On("Delete", uint(1)).Return(errors.New("database error"))

	err := milkService.DeleteMilkCollection(1)

	assert.Error(t, err)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_GetMilkCollectionsByFarmIDWithPagination_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	expectedCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: 1,
			Liters:   35.0,
			Date:     time.Now(),
		},
		{
			ID:       2,
			AnimalID: 2,
			Liters:   30.0,
			Date:     time.Now(),
		},
	}

	mockMilkRepo.On("FindByFarmIDWithPagination", uint(1), 1, 10).Return(expectedCollections, int64(2), nil)

	result, total, err := milkService.GetMilkCollectionsByFarmIDWithPagination(1, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_GetMilkCollectionsByFarmIDWithPagination_Error(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	mockMilkRepo.On("FindByFarmIDWithPagination", uint(1), 1, 10).Return(nil, int64(0), errors.New("database error"))

	result, total, err := milkService.GetMilkCollectionsByFarmIDWithPagination(1, 1, 10)

	assert.Error(t, err)
	assert.Equal(t, int64(0), total)
	assert.Nil(t, result)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_GetMilkCollectionsByFarmIDWithDateRangePaginated_Success(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	expectedCollections := []models.MilkCollection{
		{
			ID:       1,
			AnimalID: 1,
			Liters:   35.0,
			Date:     time.Now(),
		},
	}

	mockMilkRepo.On("FindByFarmIDWithDateRangePaginated", uint(1), &startDate, &endDate, 1, 10).Return(expectedCollections, int64(1), nil)

	result, total, err := milkService.GetMilkCollectionsByFarmIDWithDateRangePaginated(1, &startDate, &endDate, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
	mockMilkRepo.AssertExpectations(t)
}

func TestMilkCollectionService_GetMilkCollectionsByFarmIDWithDateRangePaginated_Error(t *testing.T) {
	mockMilkRepo := new(mocks.MockMilkCollectionRepository)
	mockAnimalRepo := new(mocks.MockAnimalRepository)

	batchService := service.NewBatchService(mockAnimalRepo, mockMilkRepo)
	milkService := service.NewMilkCollectionService(mockMilkRepo, batchService)

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	mockMilkRepo.On("FindByFarmIDWithDateRangePaginated", uint(1), &startDate, &endDate, 1, 10).Return(nil, int64(0), errors.New("database error"))

	result, total, err := milkService.GetMilkCollectionsByFarmIDWithDateRangePaginated(1, &startDate, &endDate, 1, 10)

	assert.Error(t, err)
	assert.Equal(t, int64(0), total)
	assert.Nil(t, result)
	mockMilkRepo.AssertExpectations(t)
}
