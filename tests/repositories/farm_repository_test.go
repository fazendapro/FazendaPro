package repositories

import (
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupFarmTestDB(t *testing.T) (*repository.Database, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Company{}, &models.Farm{})
	require.NoError(t, err)

	database := &repository.Database{DB: db}
	return database, db
}

func createTestCompany(t *testing.T, db *gorm.DB) *models.Company {
	company := &models.Company{
		CompanyName: "Test Company",
		Location:    "Test Location",
		FarmCNPJ:    "12345678901234",
	}
	require.NoError(t, db.Create(company).Error)
	return company
}

func TestFarmRepository_FindByID(t *testing.T) {
	database, db := setupFarmTestDB(t)
	repo := repository.NewFarmRepository(database)

	company := createTestCompany(t, db)
	farm := &models.Farm{
		CompanyID: company.ID,
		Logo:      "",
	}
	require.NoError(t, db.Create(farm).Error)

	found, err := repo.FindByID(farm.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, farm.ID, found.ID)
	assert.Equal(t, company.ID, found.CompanyID)
}

func TestFarmRepository_FindByID_NotFound(t *testing.T) {
	database, _ := setupFarmTestDB(t)
	repo := repository.NewFarmRepository(database)

	found, err := repo.FindByID(999)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestFarmRepository_Update(t *testing.T) {
	database, db := setupFarmTestDB(t)
	repo := repository.NewFarmRepository(database)

	company := createTestCompany(t, db)
	farm := &models.Farm{
		CompanyID: company.ID,
		Logo:      "",
	}
	require.NoError(t, db.Create(farm).Error)

	farm.Logo = "new-logo-url.png"
	err := repo.Update(farm)
	assert.NoError(t, err)

	var updatedFarm models.Farm
	require.NoError(t, db.First(&updatedFarm, farm.ID).Error)
	assert.Equal(t, "new-logo-url.png", updatedFarm.Logo)
}

func TestFarmRepository_LoadCompanyData(t *testing.T) {
	database, db := setupFarmTestDB(t)
	repo := repository.NewFarmRepository(database)

	company := createTestCompany(t, db)
	farm := &models.Farm{
		CompanyID: company.ID,
		Logo:      "",
	}
	require.NoError(t, db.Create(farm).Error)

	err := repo.LoadCompanyData(farm)
	assert.NoError(t, err)
	assert.NotNil(t, farm.Company)
	assert.Equal(t, company.ID, farm.Company.ID)
	assert.Equal(t, "Test Company", farm.Company.CompanyName)
}

func TestFarmRepository_LoadCompanyData_NotFound(t *testing.T) {
	database, _ := setupFarmTestDB(t)
	repo := repository.NewFarmRepository(database)

	farm := &models.Farm{
		ID: 999,
	}

	err := repo.LoadCompanyData(farm)
	assert.Error(t, err)
}


