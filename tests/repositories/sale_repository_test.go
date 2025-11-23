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

	err = db.AutoMigrate(&models.Company{}, &models.Farm{}, &models.Animal{}, &models.Sale{})
	require.NoError(t, err)

	return db
}

func createTestFarm(t *testing.T, db *gorm.DB) *models.Farm {
	company := &models.Company{
		CompanyName: "Test Company",
	}
	require.NoError(t, db.Create(company).Error)

	farm := &models.Farm{
		CompanyID: company.ID,
		Logo:      "",
	}
	require.NoError(t, db.Create(farm).Error)
	return farm
}

func TestSaleRepository_Create(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	farm := createTestFarm(t, db)

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

	farm := createTestFarm(t, db)

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

	farm := createTestFarm(t, db)

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
	assert.Equal(t, sale2.ID, sales[0].ID)
	assert.Equal(t, sale1.ID, sales[1].ID)
}

func TestSaleRepository_GetByAnimalID(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	farm := createTestFarm(t, db)

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

	farm := createTestFarm(t, db)

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
		SaleDate:  now.Add(-3 * 24 * time.Hour),
		Notes:     "Within range",
	}
	db.Create(sale1)

	sale2 := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "Maria Santos",
		Price:     2000.00,
		SaleDate:  now.Add(-10 * 24 * time.Hour),
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

	farm := createTestFarm(t, db)

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

	sale.BuyerName = "Updated Buyer"
	sale.Price = 2000.00
	sale.Notes = "Updated notes"

	err := repo.Update(context.Background(), sale)
	assert.NoError(t, err)

	updatedSale, err := repo.GetByID(context.Background(), sale.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Buyer", updatedSale.BuyerName)
	assert.Equal(t, 2000.00, updatedSale.Price)
	assert.Equal(t, "Updated notes", updatedSale.Notes)
}

func TestSaleRepository_Delete(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	farm := createTestFarm(t, db)

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

	_, err = repo.GetByID(context.Background(), sale.ID)
	assert.Error(t, err)
}

func TestSaleRepository_GetMonthlySalesCount(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	farm := createTestFarm(t, db)

	animal := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Test Animal",
		Sex:        0,
		Breed:      "Holstein",
		Type:       "Cattle",
	}
	db.Create(animal)

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	sale1 := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  startOfMonth.Add(5 * 24 * time.Hour),
		Notes:     "Sale in current month",
	}
	db.Create(sale1)

	sale2 := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "Maria Santos",
		Price:     2000.00,
		SaleDate:  startOfMonth.Add(10 * 24 * time.Hour),
		Notes:     "Sale in current month",
	}
	db.Create(sale2)

	sale3 := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "Pedro Costa",
		Price:     1800.00,
		SaleDate:  startOfMonth.AddDate(0, -1, 0),
		Notes:     "Sale in previous month",
	}
	db.Create(sale3)

	count, err := repo.GetMonthlySalesCount(context.Background(), farm.ID, startOfMonth, endOfMonth)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestSaleRepository_GetMonthlySalesData(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	farm := createTestFarm(t, db)

	animal := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Test Animal",
		Sex:        0,
		Breed:      "Holstein",
		Type:       "Cattle",
		AnimalType: 0,
		Status:     0,
		Purpose:    0,
	}
	db.Create(animal)

	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	sale1 := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "Buyer 1",
		Price:     1500.50,
		SaleDate:  time.Date(currentYear, time.Month(currentMonth), 5, 0, 0, 0, 0, time.UTC),
	}
	db.Create(sale1)

	sale2 := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "Buyer 2",
		Price:     2000.00,
		SaleDate:  time.Date(currentYear, time.Month(currentMonth), 15, 0, 0, 0, 0, time.UTC),
	}
	db.Create(sale2)

	lastMonth := currentMonth - 1
	if lastMonth == 0 {
		lastMonth = 12
		currentYear--
	}
	sale3 := &models.Sale{
		AnimalID:  animal.ID,
		FarmID:    farm.ID,
		BuyerName: "Buyer 3",
		Price:     1800.00,
		SaleDate:  time.Date(currentYear, time.Month(lastMonth), 10, 0, 0, 0, 0, time.UTC),
	}
	db.Create(sale3)

	monthlyData, err := repo.GetMonthlySalesData(context.Background(), farm.ID, 6)
	assert.NoError(t, err)
	assert.Len(t, monthlyData, 6)

	for i := 0; i < len(monthlyData)-1; i++ {
		if monthlyData[i].Year == monthlyData[i+1].Year {
			assert.True(t, monthlyData[i].Month <= monthlyData[i+1].Month)
		} else {
			assert.True(t, monthlyData[i].Year < monthlyData[i+1].Year)
		}
	}

	validMonths := []string{"Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"}
	for _, data := range monthlyData {
		assert.Contains(t, validMonths, data.Month)
	}
}

func TestSaleRepository_GetMonthlySalesData_Empty(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	farm := createTestFarm(t, db)

	monthlyData, err := repo.GetMonthlySalesData(context.Background(), farm.ID, 6)
	assert.NoError(t, err)
	assert.Len(t, monthlyData, 6)

	for _, data := range monthlyData {
		assert.Equal(t, 0.0, data.Sales)
		assert.Equal(t, int64(0), data.Count)
	}
}

func TestSaleRepository_GetMonthlySalesData_MultipleMonths(t *testing.T) {
	db := setupSaleTestDB(t)
	repo := repository.NewSaleRepository(db)

	farm := createTestFarm(t, db)

	animal := &models.Animal{
		FarmID:     farm.ID,
		AnimalName: "Test Animal",
		Sex:        0,
		Breed:      "Holstein",
		Type:       "Cattle",
		AnimalType: 0,
		Status:     0,
		Purpose:    0,
	}
	db.Create(animal)

	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	for i := 0; i < 3; i++ {
		month := currentMonth - i
		year := currentYear
		if month <= 0 {
			month += 12
			year--
		}

		sale := &models.Sale{
			AnimalID:  animal.ID,
			FarmID:    farm.ID,
			BuyerName: "Buyer",
			Price:     1000.0 * float64(i+1),
			SaleDate:  time.Date(year, time.Month(month), 10, 0, 0, 0, 0, time.UTC),
		}
		db.Create(sale)
	}

	monthlyData, err := repo.GetMonthlySalesData(context.Background(), farm.ID, 6)
	assert.NoError(t, err)
	assert.Len(t, monthlyData, 6)

	hasData := false
	for _, data := range monthlyData {
		if data.Sales > 0 || data.Count > 0 {
			hasData = true
			break
		}
	}
	assert.True(t, hasData)
}
