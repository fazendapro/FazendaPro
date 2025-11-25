package services

import (
	"context"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSale_UpdatesAnimalStatus(t *testing.T) {
	mockSaleRepo := &MockSaleRepository{}
	mockAnimalRepo := &MockAnimalRepository{}
	mockCache := new(MockCache)

	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	animal := &models.Animal{
		ID:     1,
		FarmID: 1,
		Status: models.AnimalStatusActive,
	}

	sale := &models.Sale{
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1000.0,
		SaleDate:  time.Now(),
		Notes:     "Venda teste",
	}

	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)
	mockSaleRepo.On("Create", mock.Anything, sale).Return(nil)
	mockAnimalRepo.On("Update", mock.MatchedBy(func(a *models.Animal) bool {
		return a.Status == models.AnimalStatusSold
	})).Return(nil)
	mockCache.On("Delete", "dashboard:overview:1").Return(nil)
	for months := 6; months <= 24; months += 6 {
		mockCache.On("Delete", mock.AnythingOfType("string")).Return(nil)
	}

	err := saleService.CreateSale(context.Background(), sale)

	assert.NoError(t, err)
	mockAnimalRepo.AssertExpectations(t)
	mockSaleRepo.AssertExpectations(t)
}

func TestCreateSale_PreventsSaleOfAlreadySoldAnimal(t *testing.T) {
	mockSaleRepo := &MockSaleRepository{}
	mockAnimalRepo := &MockAnimalRepository{}

	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, new(MockCache))

	animal := &models.Animal{
		ID:     1,
		FarmID: 1,
		Status: models.AnimalStatusSold,
	}

	sale := &models.Sale{
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1000.0,
		SaleDate:  time.Now(),
		Notes:     "Venda teste",
	}

	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)

	err := saleService.CreateSale(context.Background(), sale)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "animal is already sold")
	mockAnimalRepo.AssertExpectations(t)
	mockSaleRepo.AssertNotCalled(t, "Create")
}

func TestDeleteSale_RevertsAnimalStatus(t *testing.T) {
	mockSaleRepo := &MockSaleRepository{}
	mockAnimalRepo := &MockAnimalRepository{}
	mockCache := new(MockCache)

	saleService := service.NewSaleService(mockSaleRepo, mockAnimalRepo, mockCache)

	sale := &models.Sale{
		ID:       1,
		AnimalID: 1,
		FarmID:   1,
	}

	animal := &models.Animal{
		ID:     1,
		FarmID: 1,
		Status: models.AnimalStatusSold,
	}

	farmID := uint(1)
	mockSaleRepo.On("GetByID", mock.Anything, uint(1), farmID).Return(sale, nil)
	mockSaleRepo.On("Delete", mock.Anything, uint(1), farmID).Return(nil)
	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)
	mockAnimalRepo.On("Update", mock.MatchedBy(func(a *models.Animal) bool {
		return a.Status == models.AnimalStatusActive
	})).Return(nil)
	mockCache.On("Delete", "dashboard:overview:1").Return(nil)
	for months := 6; months <= 24; months += 6 {
		mockCache.On("Delete", mock.AnythingOfType("string")).Return(nil)
	}

	err := saleService.DeleteSale(context.Background(), 1, farmID)

	assert.NoError(t, err)
	mockAnimalRepo.AssertExpectations(t)
	mockSaleRepo.AssertExpectations(t)
}
