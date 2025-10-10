package services

import (
	"context"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSale_UpdatesAnimalStatus(t *testing.T) {
	// Arrange
	mockSaleRepo := &MockSaleRepository{}
	mockAnimalRepo := &MockAnimalRepository{}

	service := NewSaleService(mockSaleRepo, mockAnimalRepo)

	// Create test animal with active status
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

	// Mock expectations
	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)
	mockSaleRepo.On("Create", mock.Anything, sale).Return(nil)
	mockAnimalRepo.On("Update", mock.MatchedBy(func(a *models.Animal) bool {
		return a.Status == models.AnimalStatusSold
	})).Return(nil)

	// Act
	err := service.CreateSale(context.Background(), sale)

	// Assert
	assert.NoError(t, err)
	mockAnimalRepo.AssertExpectations(t)
	mockSaleRepo.AssertExpectations(t)
}

func TestCreateSale_PreventsSaleOfAlreadySoldAnimal(t *testing.T) {
	// Arrange
	mockSaleRepo := &MockSaleRepository{}
	mockAnimalRepo := &MockAnimalRepository{}

	service := NewSaleService(mockSaleRepo, mockAnimalRepo)

	// Create test animal with sold status
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

	// Mock expectations
	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)

	// Act
	err := service.CreateSale(context.Background(), sale)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "animal is already sold")
	mockAnimalRepo.AssertExpectations(t)
	mockSaleRepo.AssertNotCalled(t, "Create")
}

func TestDeleteSale_RevertsAnimalStatus(t *testing.T) {
	// Arrange
	mockSaleRepo := &MockSaleRepository{}
	mockAnimalRepo := &MockAnimalRepository{}

	service := NewSaleService(mockSaleRepo, mockAnimalRepo)

	// Create test sale
	sale := &models.Sale{
		ID:       1,
		AnimalID: 1,
		FarmID:   1,
	}

	// Create test animal with sold status
	animal := &models.Animal{
		ID:     1,
		FarmID: 1,
		Status: models.AnimalStatusSold,
	}

	// Mock expectations
	mockSaleRepo.On("GetByID", mock.Anything, uint(1)).Return(sale, nil)
	mockSaleRepo.On("Delete", mock.Anything, uint(1)).Return(nil)
	mockAnimalRepo.On("FindByID", uint(1)).Return(animal, nil)
	mockAnimalRepo.On("Update", mock.MatchedBy(func(a *models.Animal) bool {
		return a.Status == models.AnimalStatusActive
	})).Return(nil)

	// Act
	err := service.DeleteSale(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	mockAnimalRepo.AssertExpectations(t)
	mockSaleRepo.AssertExpectations(t)
}
