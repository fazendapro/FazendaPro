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

func setupOverviewTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Company{}, &models.Farm{}, &models.Animal{}, &models.Sale{})
	require.NoError(t, err)

	return db
}

func TestSaleRepository_GetOverviewStats(t *testing.T) {
	db := setupOverviewTestDB(t)
	saleRepo := repository.NewSaleRepository(db)

	company := &models.Company{CompanyName: "Test Company"}
	db.Create(company)

	farm := &models.Farm{CompanyID: company.ID}
	db.Create(farm)

	maleAnimal1 := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Male Animal 1",
		Sex:        1,
		Breed:      "Holstein",
		Type:       "Cattle",
	}
	db.Create(maleAnimal1)

	maleAnimal2 := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Male Animal 2",
		Sex:        1,
		Breed:      "Holstein",
		Type:       "Cattle",
	}
	db.Create(maleAnimal2)

	femaleAnimal1 := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Female Animal 1",
		Sex:        0,
		Breed:      "Holstein",
		Type:       "Cattle",
	}
	db.Create(femaleAnimal1)

	sale1 := &models.Sale{
		AnimalID:  maleAnimal1.ID,
		FarmID:    farm.ID,
		BuyerName: "Buyer 1",
		Price:     1500.50,
		SaleDate:  time.Now(),
	}
	db.Create(sale1)

	sale2 := &models.Sale{
		AnimalID:  maleAnimal2.ID,
		FarmID:    farm.ID,
		BuyerName: "Buyer 2",
		Price:     2000.00,
		SaleDate:  time.Now(),
	}
	db.Create(sale2)

	stats, err := saleRepo.GetOverviewStats(context.Background(), farm.ID)
	require.NoError(t, err)

	assert.Equal(t, int64(2), stats.MalesCount)
	assert.Equal(t, int64(1), stats.FemalesCount)
	assert.Equal(t, int64(2), stats.TotalSold)
	assert.Equal(t, 3500.50, stats.TotalRevenue)
}

func TestSaleRepository_GetOverviewStats_NoData(t *testing.T) {
	db := setupOverviewTestDB(t)
	saleRepo := repository.NewSaleRepository(db)

	company := &models.Company{CompanyName: "Test Company"}
	db.Create(company)

	farm := &models.Farm{CompanyID: company.ID}
	db.Create(farm)

	stats, err := saleRepo.GetOverviewStats(context.Background(), farm.ID)
	require.NoError(t, err)

	assert.Equal(t, int64(0), stats.MalesCount)
	assert.Equal(t, int64(0), stats.FemalesCount)
	assert.Equal(t, int64(0), stats.TotalSold)
	assert.Equal(t, 0.0, stats.TotalRevenue)
}
