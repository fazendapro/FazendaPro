package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/cache"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
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

func (m *MockSaleRepository) GetByID(ctx context.Context, id uint, farmID uint) (*models.Sale, error) {
	args := m.Called(ctx, id, farmID)
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

func (m *MockSaleRepository) GetByAnimalID(ctx context.Context, animalID uint, farmID uint) ([]*models.Sale, error) {
	args := m.Called(ctx, animalID, farmID)
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

func (m *MockSaleRepository) GetMonthlySalesCount(ctx context.Context, farmID uint, startDate, endDate time.Time) (int64, error) {
	args := m.Called(ctx, farmID, startDate, endDate)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSaleRepository) GetMonthlySalesData(ctx context.Context, farmID uint, months int) ([]repository.MonthlySalesData, error) {
	args := m.Called(ctx, farmID, months)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.MonthlySalesData), args.Error(1)
}

func (m *MockSaleRepository) GetOverviewStats(ctx context.Context, farmID uint) (*repository.OverviewStats, error) {
	args := m.Called(ctx, farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.OverviewStats), args.Error(1)
}

func (m *MockSaleRepository) Update(ctx context.Context, sale *models.Sale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockSaleRepository) Delete(ctx context.Context, id uint, farmID uint) error {
	args := m.Called(ctx, id, farmID)
	return args.Error(0)
}

func TestSaleService_CreateSale_Success(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	now := time.Now()

	animal := &models.Animal{
		ID:         1,
		FarmID:     1,
		AnimalName: "Test Animal",
		Status:     0, // Active
	}

	sale := &models.Sale{
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  now,
		Notes:     "Test sale",
	}

	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)
	mockSaleRepo.On("Create", ctx, sale).Return(nil)
	mockAnimalRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(nil)
	mockCache.On("Delete", "dashboard:overview:1").Return(nil)
	for months := 6; months <= 24; months += 6 {
		mockCache.On("Delete", mock.AnythingOfType("string")).Return(nil)
	}

	err := saleService.CreateSale(ctx, sale)

	assert.NoError(t, err)
	mockAnimalRepo.AssertExpectations(t)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_CreateSale_AnimalNotFound(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

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

	mockAnimalRepo.On("FindByID", uint(1)).Return(nil, errors.New("animal not found"))

	err := saleService.CreateSale(ctx, sale)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "animal not found")
	mockAnimalRepo.AssertExpectations(t)
}

func TestSaleService_CreateSale_AnimalAlreadySold(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	ctx := context.Background()
	now := time.Now()

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

	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)

	err := saleService.CreateSale(ctx, sale)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already sold")
	mockAnimalRepo.AssertExpectations(t)
}

func TestSaleService_CreateSale_InvalidData(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

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
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

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

	mockSaleRepo.On("GetByFarmID", ctx, farmID).Return(expectedSales, nil)

	sales, err := saleService.GetSalesByFarmID(ctx, farmID)

	assert.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, expectedSales[0].ID, sales[0].ID)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetSalesByAnimalID(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	ctx := context.Background()
	animalID := uint(1)
	farmID := uint(1)

	expectedSales := []*models.Sale{
		{
			ID:        1,
			AnimalID:  animalID,
			FarmID:    farmID,
			BuyerName: "João Silva",
			Price:     1500.50,
			SaleDate:  time.Now(),
		},
	}

	animal := &models.Animal{
		ID:     animalID,
		FarmID: farmID,
	}
	mockAnimalRepo.On("FindByID", animalID).Return(animal, nil)
	mockSaleRepo.On("GetByAnimalID", ctx, animalID, farmID).Return(expectedSales, nil)

	sales, err := saleService.GetSalesByAnimalID(ctx, animalID, farmID)

	assert.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, expectedSales[0].ID, sales[0].ID)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetSalesByDateRange(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

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

	mockSaleRepo.On("GetByDateRange", ctx, farmID, startDate, endDate).Return(expectedSales, nil)

	sales, err := saleService.GetSalesByDateRange(ctx, farmID, startDate, endDate)

	assert.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, expectedSales[0].ID, sales[0].ID)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetSalesByDateRange_InvalidDateRange(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	ctx := context.Background()
	farmID := uint(1)
	startDate := time.Now()
	endDate := time.Now().Add(-7 * 24 * time.Hour)

	sales, err := saleService.GetSalesByDateRange(ctx, farmID, startDate, endDate)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "start date cannot be after end date")
	assert.Nil(t, sales)
}

func TestSaleService_UpdateSale_Success(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	farmID := uint(1)

	sale := &models.Sale{
		ID:        1,
		AnimalID:  1,
		FarmID:    farmID,
		BuyerName: "Updated Buyer",
		Price:     2000.00,
		SaleDate:  time.Now(),
		Notes:     "Updated notes",
	}

	existingSale := &models.Sale{
		ID:        1,
		AnimalID:  1,
		FarmID:    farmID,
		BuyerName: "Old Buyer",
		Price:     1000.00,
		SaleDate:  time.Now(),
	}

	mockSaleRepo.On("GetByID", ctx, sale.ID, farmID).Return(existingSale, nil)
	mockSaleRepo.On("Update", ctx, sale).Return(nil)
	mockCache.On("Delete", "dashboard:overview:1").Return(nil)
	for months := 6; months <= 24; months += 6 {
		mockCache.On("Delete", mock.AnythingOfType("string")).Return(nil)
	}

	err := saleService.UpdateSale(ctx, sale, farmID)

	assert.NoError(t, err)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_UpdateSale_InvalidData(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	ctx := context.Background()

	farmID := uint(1)
	testCases := []struct {
		name          string
		sale          *models.Sale
		expectedError string
	}{
		{
			name: "Missing Sale ID",
			sale: &models.Sale{
				AnimalID:  1,
				FarmID:    farmID,
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
				FarmID:   farmID,
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
				FarmID:    farmID,
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
				FarmID:    farmID,
				BuyerName: "João Silva",
				Price:     1500.50,
			},
			expectedError: "sale date is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.sale.ID != 0 {
				existingSale := &models.Sale{
					ID:        tc.sale.ID,
					AnimalID:  tc.sale.AnimalID,
					FarmID:    farmID,
					BuyerName: "Old Buyer",
					Price:     1000.00,
					SaleDate:  time.Now(),
				}
				mockSaleRepo.On("GetByID", ctx, tc.sale.ID, farmID).Return(existingSale, nil)
			}
			err := saleService.UpdateSale(ctx, tc.sale, farmID)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

func TestSaleService_DeleteSale_Success(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	saleID := uint(1)

	sale := &models.Sale{
		ID:        saleID,
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  time.Now(),
	}

	animal := &models.Animal{
		ID:         1,
		FarmID:     1,
		AnimalName: "Test Animal",
		Status:     1, // Sold
	}

	farmID := uint(1)
	mockSaleRepo.On("GetByID", ctx, saleID, farmID).Return(sale, nil)
	mockSaleRepo.On("Delete", ctx, saleID, farmID).Return(nil)
	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)
	mockAnimalRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(nil)
	mockCache.On("Delete", "dashboard:overview:1").Return(nil)
	for months := 6; months <= 24; months += 6 {
		mockCache.On("Delete", mock.AnythingOfType("string")).Return(nil)
	}

	err := saleService.DeleteSale(ctx, saleID, farmID)

	assert.NoError(t, err)
	mockSaleRepo.AssertExpectations(t)
	mockAnimalRepo.AssertExpectations(t)
}

func TestSaleService_DeleteSale_NotFound(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	ctx := context.Background()
	saleID := uint(1)
	farmID := uint(1)

	mockSaleRepo.On("GetByID", ctx, saleID, farmID).Return(nil, errors.New("sale not found"))

	err := saleService.DeleteSale(ctx, saleID, farmID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sale not found")
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_DeleteSale_WithAnimalUpdate(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	saleID := uint(1)
	animalID := uint(1)

	sale := &models.Sale{
		ID:        saleID,
		AnimalID:  animalID,
		FarmID:    1,
		BuyerName: "Comprador Teste",
		Price:     1000.0,
		SaleDate:  time.Now(),
	}

	animal := &models.Animal{
		ID:     animalID,
		FarmID: 1,
		Status: models.AnimalStatusSold,
	}

	farmID := uint(1)
	mockSaleRepo.On("GetByID", ctx, saleID, farmID).Return(sale, nil)
	mockSaleRepo.On("Delete", ctx, saleID, farmID).Return(nil)
	mockAnimalRepo.On("FindByID", animalID).Return(animal, nil)
	mockAnimalRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(nil)
	mockCache.On("Delete", "dashboard:overview:1").Return(nil)
	for months := 6; months <= 24; months += 6 {
		mockCache.On("Delete", mock.AnythingOfType("string")).Return(nil)
	}

	err := saleService.DeleteSale(ctx, saleID, farmID)

	assert.NoError(t, err)
	mockSaleRepo.AssertExpectations(t)
	mockAnimalRepo.AssertExpectations(t)
}

func TestSaleService_DeleteSale_AnimalNotFound(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	saleID := uint(1)
	animalID := uint(1)

	sale := &models.Sale{
		ID:        saleID,
		AnimalID:  animalID,
		FarmID:    1,
		BuyerName: "Comprador Teste",
		Price:     1000.0,
		SaleDate:  time.Now(),
	}

	farmID := uint(1)
	mockSaleRepo.On("GetByID", ctx, saleID, farmID).Return(sale, nil)
	mockSaleRepo.On("Delete", ctx, saleID, farmID).Return(nil)
	mockAnimalRepo.On("FindByID", animalID).Return(nil, errors.New("animal not found"))
	mockCache.On("Delete", "dashboard:overview:1").Return(nil)
	for months := 6; months <= 24; months += 6 {
		mockCache.On("Delete", mock.AnythingOfType("string")).Return(nil)
	}

	err := saleService.DeleteSale(ctx, saleID, farmID)

	assert.NoError(t, err)
	mockSaleRepo.AssertExpectations(t)
	mockAnimalRepo.AssertExpectations(t)
}

func TestSaleService_GetSalesHistory(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

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

	mockSaleRepo.On("GetByFarmID", ctx, farmID).Return(expectedSales, nil)

	sales, err := saleService.GetSalesHistory(ctx, farmID)

	assert.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, expectedSales[0].ID, sales[0].ID)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetOverviewStats_Success(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	farmID := uint(1)

	expectedStats := &repository.OverviewStats{
		MalesCount:   10,
		FemalesCount: 15,
		TotalSold:    5,
		TotalRevenue: 15000.50,
	}

	mockCache.On("Get", "dashboard:overview:1", mock.Anything).Return(cache.ErrCacheMiss)
	mockSaleRepo.On("GetOverviewStats", ctx, farmID).Return(expectedStats, nil)
	mockCache.On("Set", "dashboard:overview:1", expectedStats, int32(600)).Return(nil)

	stats, err := saleService.GetOverviewStats(ctx, farmID)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(10), stats.MalesCount)
	assert.Equal(t, int64(15), stats.FemalesCount)
	assert.Equal(t, int64(5), stats.TotalSold)
	assert.Equal(t, 15000.50, stats.TotalRevenue)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetOverviewStats_RepositoryError(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	farmID := uint(1)

	mockCache.On("Get", "dashboard:overview:1", mock.Anything).Return(cache.ErrCacheMiss)
	mockSaleRepo.On("GetOverviewStats", ctx, farmID).Return(nil, assert.AnError)

	stats, err := saleService.GetOverviewStats(ctx, farmID)

	assert.Error(t, err)
	assert.Nil(t, stats)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetMonthlySalesCount(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	ctx := context.Background()
	farmID := uint(1)
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	expectedCount := int64(5)

	mockSaleRepo.On("GetMonthlySalesCount", ctx, farmID, startOfMonth, endOfMonth).Return(expectedCount, nil)

	count, err := saleService.GetMonthlySalesCount(ctx, farmID, startOfMonth, endOfMonth)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetMonthlySalesCount_InvalidDateRange(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	ctx := context.Background()
	farmID := uint(1)
	now := time.Now()
	startDate := now.AddDate(0, 0, 10)
	endDate := now

	count, err := saleService.GetMonthlySalesCount(ctx, farmID, startDate, endDate)

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
	assert.Contains(t, err.Error(), "start date cannot be after end date")
	mockSaleRepo.AssertNotCalled(t, "GetMonthlySalesCount")
}

func TestSaleService_GetSaleByID_Success(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	ctx := context.Background()
	saleID := uint(1)
	expectedSale := &models.Sale{
		ID:        saleID,
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "Comprador Teste",
		Price:     1000.0,
		SaleDate:  time.Now(),
	}

	farmID := uint(1)
	mockSaleRepo.On("GetByID", ctx, saleID, farmID).Return(expectedSale, nil)

	sale, err := saleService.GetSaleByID(ctx, saleID, farmID)

	assert.NoError(t, err)
	assert.NotNil(t, sale)
	assert.Equal(t, saleID, sale.ID)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetSaleByID_NotFound(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	ctx := context.Background()
	saleID := uint(999)

	farmID := uint(1)
	mockSaleRepo.On("GetByID", ctx, saleID, farmID).Return(nil, errors.New("not found"))

	sale, err := saleService.GetSaleByID(ctx, saleID, farmID)

	assert.Error(t, err)
	assert.Nil(t, sale)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetMonthlySalesData_Success(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	farmID := uint(1)
	months := 6

	expectedData := []repository.MonthlySalesData{
		{Month: "2024-01", Count: 5},
		{Month: "2024-02", Count: 3},
	}

	mockCache.On("Get", "dashboard:monthly:1:6", mock.Anything).Return(cache.ErrCacheMiss)
	mockSaleRepo.On("GetMonthlySalesData", ctx, farmID, months).Return(expectedData, nil)
	mockCache.On("Set", "dashboard:monthly:1:6", expectedData, int32(900)).Return(nil)

	data, err := saleService.GetMonthlySalesData(ctx, farmID, months)

	assert.NoError(t, err)
	assert.Len(t, data, 2)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetMonthlySalesData_DefaultMonths(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	farmID := uint(1)

	expectedData := []repository.MonthlySalesData{}

	mockCache.On("Get", "dashboard:monthly:1:12", mock.Anything).Return(cache.ErrCacheMiss)
	mockSaleRepo.On("GetMonthlySalesData", ctx, farmID, 12).Return(expectedData, nil)
	mockCache.On("Set", "dashboard:monthly:1:12", expectedData, int32(900)).Return(nil)

	data, err := saleService.GetMonthlySalesData(ctx, farmID, 0)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetMonthlySalesData_MaxMonths(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	farmID := uint(1)

	expectedData := []repository.MonthlySalesData{}

	mockCache.On("Get", "dashboard:monthly:1:24", mock.Anything).Return(cache.ErrCacheMiss)
	mockSaleRepo.On("GetMonthlySalesData", ctx, farmID, 24).Return(expectedData, nil)
	mockCache.On("Set", "dashboard:monthly:1:24", expectedData, int32(900)).Return(nil)

	data, err := saleService.GetMonthlySalesData(ctx, farmID, 30)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_GetMonthlySalesData_RepositoryError(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	farmID := uint(1)
	months := 6

	mockCache.On("Get", "dashboard:monthly:1:6", mock.Anything).Return(cache.ErrCacheMiss)
	mockSaleRepo.On("GetMonthlySalesData", ctx, farmID, months).Return(nil, errors.New("database error"))

	data, err := saleService.GetMonthlySalesData(ctx, farmID, months)

	assert.Error(t, err)
	assert.Nil(t, data)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_CreateSale_AnimalWrongFarm(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	ctx := context.Background()
	now := time.Now()

	animal := &models.Animal{
		ID:         1,
		FarmID:     2,
		AnimalName: "Test Animal",
		Status:     0,
	}

	sale := &models.Sale{
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  now,
	}

	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)

	err := saleService.CreateSale(ctx, sale)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "animal does not belong to the specified farm")
	mockAnimalRepo.AssertExpectations(t)
	mockSaleRepo.AssertNotCalled(t, "Create")
}

func TestSaleService_CreateSale_UpdateAnimalStatusFails(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	ctx := context.Background()
	now := time.Now()

	animal := &models.Animal{
		ID:         1,
		FarmID:     1,
		AnimalName: "Test Animal",
		Status:     0,
	}

	sale := &models.Sale{
		ID:        1,
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  now,
	}

	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)
	mockSaleRepo.On("Create", ctx, sale).Return(nil)
	mockAnimalRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(errors.New("update failed"))
	mockSaleRepo.On("Delete", ctx, uint(1), uint(1)).Return(nil)

	err := saleService.CreateSale(ctx, sale)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update animal status")
	mockAnimalRepo.AssertExpectations(t)
	mockSaleRepo.AssertExpectations(t)
}

func TestSaleService_UpdateSale_RepositoryError(t *testing.T) {
	mockSaleRepo := new(MockSaleRepository)
	mockAnimalRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	ctx := context.Background()
	now := time.Now()

	sale := &models.Sale{
		ID:        1,
		FarmID:    1,
		BuyerName: "Comprador Atualizado",
		Price:     2000.0,
		SaleDate:  now,
	}

	farmID := uint(1)
	existingSale := &models.Sale{
		ID:        sale.ID,
		AnimalID:  sale.AnimalID,
		FarmID:    farmID,
		BuyerName: "Old Buyer",
		Price:     1000.00,
		SaleDate:  time.Now(),
	}
	mockSaleRepo.On("GetByID", ctx, sale.ID, farmID).Return(existingSale, nil)
	mockSaleRepo.On("Update", ctx, sale).Return(errors.New("database error"))

	err := saleService.UpdateSale(ctx, sale, farmID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockSaleRepo.AssertExpectations(t)
}
