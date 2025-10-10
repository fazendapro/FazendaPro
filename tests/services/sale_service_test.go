package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSaleRepository struct {
	mock.Mock
}

func (m *MockSaleRepository) Create(ctx context.Context, sale *models.Sale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockSaleRepository) GetByID(ctx context.Context, id uint) (*models.Sale, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Sale), args.Error(1)
}

func (m *MockSaleRepository) GetByFarmID(ctx context.Context, farmID uint) ([]*models.Sale, error) {
	args := m.Called(ctx, farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleRepository) GetByAnimalID(ctx context.Context, animalID uint) ([]*models.Sale, error) {
	args := m.Called(ctx, animalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleRepository) GetByDateRange(ctx context.Context, farmID uint, startDate, endDate time.Time) ([]*models.Sale, error) {
	args := m.Called(ctx, farmID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleRepository) Update(ctx context.Context, sale *models.Sale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockSaleRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockAnimalRepository struct {
	mock.Mock
}

func (m *MockAnimalRepository) GetByID(ctx context.Context, id uint) (*models.Animal, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Animal), args.Error(1)
}

func (m *MockAnimalRepository) Update(ctx context.Context, animal *models.Animal) error {
	args := m.Called(ctx, animal)
	return args.Error(0)
}

func (m *MockAnimalRepository) Create(ctx context.Context, animal *models.Animal) error {
	args := m.Called(ctx, animal)
	return args.Error(0)
}

func (m *MockAnimalRepository) GetByFarmID(ctx context.Context, farmID uint) ([]*models.Animal, error) {
	args := m.Called(ctx, farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Animal), args.Error(1)
}

func (m *MockAnimalRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSaleService_CreateSale_Success(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()
	now := time.Now()

	// Mock animal
	animal := &models.Animal{
		ID:         1,
		FarmID:     1,
		AnimalName: "Test Animal",
		Status:     0, // Active
	}

	// Mock sale
	sale := &models.Sale{
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  now,
		Notes:     "Test sale",
	}

	// Setup mocks
	mockAnimalRepo.On("GetByID", ctx, uint(1)).Return(animal, nil)
	mockSaleRepo.On("Create", ctx, sale).Return(nil)
	mockAnimalRepo.On("Update", ctx, mock.AnythingOfType("*models.Animal")).Return(nil)

	// Execute
	err := saleService.CreateSale(ctx, sale)

	// Assert
	assert.NoError(t, err)
	mockAnimalRepo.AssertExpectations(t)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_CreateSale_AnimalNotFound(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()
	now := time.Now()

	sale := &models.Sale{
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  now,
		Notes:     "Test sale",
	}

	// Setup mocks
	mockAnimalRepo.On("GetByID", ctx, uint(1)).Return(nil, errors.New("animal not found"))

	// Execute
	err := saleService.CreateSale(ctx, sale)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "animal not found")
	mockAnimalRepo.AssertExpectations(t)
}

func TestSaleService_CreateSale_AnimalAlreadySold(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()
	now := time.Now()

	// Mock animal that is already sold
	animal := &models.Animal{
		ID:         1,
		FarmID:     1,
		AnimalName: "Test Animal",
		Status:     1, // Sold
	}

	sale := &models.Sale{
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  now,
		Notes:     "Test sale",
	}

	// Setup mocks
	mockAnimalRepo.On("GetByID", ctx, uint(1)).Return(animal, nil)

	// Execute
	err := saleService.CreateSale(ctx, sale)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already sold")
	mockAnimalRepo.AssertExpectations(t)
}

func TestSaleService_CreateSale_InvalidData(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()

	testCases := []struct {
		name          string
		sale          *models.Sale
		expectedError string
	}{
		{
			name: "Missing Animal ID",
			sale: &models.Sale{
				FarmID:    1,
				BuyerName: "João Silva",
				Price:     1500.50,
				SaleDate:  time.Now(),
			},
			expectedError: "animal ID is required",
		},
		{
			name: "Missing Farm ID",
			sale: &models.Sale{
				AnimalID:  1,
				BuyerName: "João Silva",
				Price:     1500.50,
				SaleDate:  time.Now(),
			},
			expectedError: "farm ID is required",
		},
		{
			name: "Missing Buyer Name",
			sale: &models.Sale{
				AnimalID: 1,
				FarmID:   1,
				Price:    1500.50,
				SaleDate: time.Now(),
			},
			expectedError: "buyer name is required",
		},
		{
			name: "Invalid Price",
			sale: &models.Sale{
				AnimalID:  1,
				FarmID:    1,
				BuyerName: "João Silva",
				Price:     0,
				SaleDate:  time.Now(),
			},
			expectedError: "price must be greater than zero",
		},
		{
			name: "Missing Sale Date",
			sale: &models.Sale{
				AnimalID:  1,
				FarmID:    1,
				BuyerName: "João Silva",
				Price:     1500.50,
			},
			expectedError: "sale date is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := saleService.CreateSale(ctx, tc.sale)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

func TestSaleService_GetSalesByFarmID(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()
	farmID := uint(1)

	expectedSales := []*models.Sale{
		{
			ID:        1,
			AnimalID:  1,
			FarmID:    farmID,
			BuyerName: "João Silva",
			Price:     1500.50,
			SaleDate:  time.Now(),
		},
	}

	// Setup mocks
	mockSaleRepo.On("GetByFarmID", ctx, farmID).Return(expectedSales, nil)

	// Execute
	sales, err := saleService.GetSalesByFarmID(ctx, farmID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, expectedSales[0].ID, sales[0].ID)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetSalesByAnimalID(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()
	animalID := uint(1)

	expectedSales := []*models.Sale{
		{
			ID:        1,
			AnimalID:  animalID,
			FarmID:    1,
			BuyerName: "João Silva",
			Price:     1500.50,
			SaleDate:  time.Now(),
		},
	}

	// Setup mocks
	mockSaleRepo.On("GetByAnimalID", ctx, animalID).Return(expectedSales, nil)

	// Execute
	sales, err := saleService.GetSalesByAnimalID(ctx, animalID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, expectedSales[0].ID, sales[0].ID)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetSalesByDateRange(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()
	farmID := uint(1)
	startDate := time.Now().Add(-7 * 24 * time.Hour)
	endDate := time.Now()

	expectedSales := []*models.Sale{
		{
			ID:        1,
			AnimalID:  1,
			FarmID:    farmID,
			BuyerName: "João Silva",
			Price:     1500.50,
			SaleDate:  time.Now().Add(-3 * 24 * time.Hour),
		},
	}

	// Setup mocks
	mockSaleRepo.On("GetByDateRange", ctx, farmID, startDate, endDate).Return(expectedSales, nil)

	// Execute
	sales, err := saleService.GetSalesByDateRange(ctx, farmID, startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, expectedSales[0].ID, sales[0].ID)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetSalesByDateRange_InvalidDateRange(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()
	farmID := uint(1)
	startDate := time.Now()
	endDate := time.Now().Add(-7 * 24 * time.Hour) // Start after end

	// Execute
	sales, err := saleService.GetSalesByDateRange(ctx, farmID, startDate, endDate)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "start date cannot be after end date")
	assert.Nil(t, sales)
}

func TestSaleService_UpdateSale_Success(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()

	sale := &models.Sale{
		ID:        1,
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "Updated Buyer",
		Price:     2000.00,
		SaleDate:  time.Now(),
		Notes:     "Updated notes",
	}

	// Setup mocks
	mockSaleRepo.On("Update", ctx, sale).Return(nil)

	// Execute
	err := saleService.UpdateSale(ctx, sale)

	// Assert
	assert.NoError(t, err)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_UpdateSale_InvalidData(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()

	testCases := []struct {
		name          string
		sale          *models.Sale
		expectedError string
	}{
		{
			name: "Missing Sale ID",
			sale: &models.Sale{
				AnimalID:  1,
				FarmID:    1,
				BuyerName: "João Silva",
				Price:     1500.50,
				SaleDate:  time.Now(),
			},
			expectedError: "sale ID is required",
		},
		{
			name: "Missing Buyer Name",
			sale: &models.Sale{
				ID:       1,
				AnimalID: 1,
				FarmID:   1,
				Price:    1500.50,
				SaleDate: time.Now(),
			},
			expectedError: "buyer name is required",
		},
		{
			name: "Invalid Price",
			sale: &models.Sale{
				ID:        1,
				AnimalID:  1,
				FarmID:    1,
				BuyerName: "João Silva",
				Price:     0,
				SaleDate:  time.Now(),
			},
			expectedError: "price must be greater than zero",
		},
		{
			name: "Missing Sale Date",
			sale: &models.Sale{
				ID:        1,
				AnimalID:  1,
				FarmID:    1,
				BuyerName: "João Silva",
				Price:     1500.50,
			},
			expectedError: "sale date is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := saleService.UpdateSale(ctx, tc.sale)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

func TestSaleService_DeleteSale_Success(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()
	saleID := uint(1)

	// Mock sale
	sale := &models.Sale{
		ID:        saleID,
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  time.Now(),
	}

	// Mock animal
	animal := &models.Animal{
		ID:         1,
		FarmID:     1,
		AnimalName: "Test Animal",
		Status:     1, // Sold
	}

	// Setup mocks
	mockSaleRepo.On("GetByID", ctx, saleID).Return(sale, nil)
	mockSaleRepo.On("Delete", ctx, saleID).Return(nil)
	mockAnimalRepo.On("GetByID", ctx, uint(1)).Return(animal, nil)
	mockAnimalRepo.On("Update", ctx, mock.AnythingOfType("*models.Animal")).Return(nil)

	// Execute
	err := saleService.DeleteSale(ctx, saleID)

	// Assert
	assert.NoError(t, err)
	mockSaleRepo.AssertExpectations(t)
	mockAnimalRepo.AssertExpectations(t)
}

func TestSaleService_DeleteSale_NotFound(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()
	saleID := uint(1)

	// Setup mocks
	mockSaleRepo.On("GetByID", ctx, saleID).Return(nil, errors.New("sale not found"))

	// Execute
	err := saleService.DeleteSale(ctx, saleID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sale not found")
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetSalesHistory(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo)

	ctx := context.Background()
	farmID := uint(1)

	expectedSales := []*models.Sale{
		{
			ID:        1,
			AnimalID:  1,
			FarmID:    farmID,
			BuyerName: "João Silva",
			Price:     1500.50,
			SaleDate:  time.Now(),
		},
	}

	// Setup mocks
	mockSaleRepo.On("GetByFarmID", ctx, farmID).Return(expectedSales, nil)

	// Execute
	sales, err := saleService.GetSalesHistory(ctx, farmID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, expectedSales[0].ID, sales[0].ID)
	mockSaleRepo.AssertExpectations(t)
}
