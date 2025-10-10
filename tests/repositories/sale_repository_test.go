package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSaleTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate tables
	err = db.AutoMigrate(&models.Farm{}, &models.Animal{}, &models.Sale{})
	require.NoError(t, err)

	return db
}

func TestSaleRepository_Create(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	// Create test data
	farm := &models.Farm{Name: "Test Farm"}
	db.Create(farm)

	animal := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Test Animal",
		Sex:        0,
		Breed:      "Holstein",
		Type:       "Cattle",
	}
	db.Create(animal)

	sale := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  time.Now(),
		Notes:     "Test sale",
	}

	err := repo.Create(context.Background(), sale)
	assert.NoError(t, err)
	assert.NotZero(t, sale.ID)
}

func TestSaleRepository_GetByID(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	// Create test data
	farm := &models.Farm{Name: "Test Farm"}
	db.Create(farm)

	animal := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Test Animal",
		Sex:        0,
		Breed:      "Holstein",
		Type:       "Cattle",
	}
	db.Create(animal)

	sale := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  time.Now(),
		Notes:     "Test sale",
	}
	db.Create(sale)

	retrievedSale, err := repo.GetByID(context.Background(), sale.ID)
	assert.NoError(t, err)
	assert.Equal(t, sale.ID, retrievedSale.ID)
	assert.Equal(t, "João Silva", retrievedSale.BuyerName)
	assert.Equal(t, 1500.50, retrievedSale.Price)
}

func TestSaleRepository_GetByFarmID(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	// Create test data
	farm := &models.Farm{Name: "Test Farm"}
	db.Create(farm)

	animal := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Test Animal",
		Sex:        0,
		Breed:      "Holstein",
		Type:       "Cattle",
	}
	db.Create(animal)

	sale1 := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  time.Now().Add(-24 * time.Hour),
		Notes:     "First sale",
	}
	db.Create(sale1)

	sale2 := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "Maria Santos",
		Price:     2000.00,
		SaleDate:  time.Now(),
		Notes:     "Second sale",
	}
	db.Create(sale2)

	sales, err := repo.GetByFarmID(context.Background(), farm.ID)
	assert.NoError(t, err)
	assert.Len(t, sales, 2)
	assert.Equal(t, sale2.ID, sales[0].ID) // Should be ordered by sale_date DESC
	assert.Equal(t, sale1.ID, sales[1].ID)
}

func TestSaleRepository_GetByAnimalID(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	// Create test data
	farm := &models.Farm{Name: "Test Farm"}
	db.Create(farm)

	animal := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Test Animal",
		Sex:        0,
		Breed:      "Holstein",
		Type:       "Cattle",
	}
	db.Create(animal)

	sale := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  time.Now(),
		Notes:     "Test sale",
	}
	db.Create(sale)

	sales, err := repo.GetByAnimalID(context.Background(), animal.ID)
	assert.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, sale.ID, sales[0].ID)
}

func TestSaleRepository_GetByDateRange(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	// Create test data
	farm := &models.Farm{Name: "Test Farm"}
	db.Create(farm)

	animal := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Test Animal",
		Sex:        0,
		Breed:      "Holstein",
		Type:       "Cattle",
	}
	db.Create(animal)

	now := time.Now()
	startDate := now.Add(-7 * 24 * time.Hour)
	endDate := now.Add(24 * time.Hour)

	sale1 := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  now.Add(-3 * 24 * time.Hour), // Within range
		Notes:     "Within range",
	}
	db.Create(sale1)

	sale2 := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "Maria Santos",
		Price:     2000.00,
		SaleDate:  now.Add(-10 * 24 * time.Hour), // Outside range
		Notes:     "Outside range",
	}
	db.Create(sale2)

	sales, err := repo.GetByDateRange(context.Background(), farm.ID, startDate, endDate)
	assert.NoError(t, err)
	assert.Len(t, sales, 1)
	assert.Equal(t, sale1.ID, sales[0].ID)
}

func TestSaleRepository_Update(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	// Create test data
	farm := &models.Farm{Name: "Test Farm"}
	db.Create(farm)

	animal := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Test Animal",
		Sex:        0,
		Breed:      "Holstein",
		Type:       "Cattle",
	}
	db.Create(animal)

	sale := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  time.Now(),
		Notes:     "Original notes",
	}
	db.Create(sale)

	// Update sale
	sale.BuyerName = "Updated Buyer"
	sale.Price = 2000.00
	sale.Notes = "Updated notes"

	err := repo.Update(context.Background(), sale)
	assert.NoError(t, err)

	// Verify update
	updatedSale, err := repo.GetByID(context.Background(), sale.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Buyer", updatedSale.BuyerName)
	assert.Equal(t, 2000.00, updatedSale.Price)
	assert.Equal(t, "Updated notes", updatedSale.Notes)
}

func TestSaleRepository_Delete(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	// Create test data
	farm := &models.Farm{Name: "Test Farm"}
	db.Create(farm)

	animal := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Test Animal",
		Sex:        0,
		Breed:      "Holstein",
		Type:       "Cattle",
	}
	db.Create(animal)

	sale := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  time.Now(),
		Notes:     "Test sale",
	}
	db.Create(sale)

	err := repo.Delete(context.Background(), sale.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(context.Background(), sale.ID)
	assert.Error(t, err)
}
