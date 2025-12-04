package repositories

import (
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupVaccineTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Company{}, &models.Farm{}, &models.Vaccine{})
	require.NoError(t, err)

	return db
}

func createTestFarmForVaccine(t *testing.T, db *gorm.DB) *models.Farm {
	company := &models.Company{
		CompanyName: "Test Company",
		Location:    "Test Location",
		FarmCNPJ:    "12345678901234",
	}
	require.NoError(t, db.Create(company).Error)

	farm := &models.Farm{
		CompanyID: company.ID,
		Logo:      "",
	}
	require.NoError(t, db.Create(farm).Error)
	return farm
}

func TestVaccineRepository_Create(t *testing.T) {
	db := setupVaccineTestDB(t)
	repo := repository.NewVaccineRepository(db)

	farm := createTestFarmForVaccine(t, db)

	vaccine := &models.Vaccine{
		FarmID:       farm.ID,
		Name:         "Vacina Aftosa",
		Description:  "Vacina contra febre aftosa",
		Manufacturer: "Fabricante XYZ",
	}

	err := repo.Create(vaccine)
	assert.NoError(t, err)
	assert.NotZero(t, vaccine.ID)
}

func TestVaccineRepository_FindByID(t *testing.T) {
	db := setupVaccineTestDB(t)
	repo := repository.NewVaccineRepository(db)

	farm := createTestFarmForVaccine(t, db)

	vaccine := &models.Vaccine{
		FarmID:       farm.ID,
		Name:         "Vacina Aftosa",
		Description:  "Vacina contra febre aftosa",
		Manufacturer: "Fabricante XYZ",
	}
	require.NoError(t, db.Create(vaccine).Error)

	found, err := repo.FindByID(vaccine.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, vaccine.ID, found.ID)
	assert.Equal(t, "Vacina Aftosa", found.Name)
}

func TestVaccineRepository_FindByID_NotFound(t *testing.T) {
	db := setupVaccineTestDB(t)
	repo := repository.NewVaccineRepository(db)

	found, err := repo.FindByID(999)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestVaccineRepository_FindByFarmID(t *testing.T) {
	db := setupVaccineTestDB(t)
	repo := repository.NewVaccineRepository(db)

	farm := createTestFarmForVaccine(t, db)

	vaccine1 := &models.Vaccine{
		FarmID: farm.ID,
		Name:   "Vacina Aftosa",
	}
	vaccine2 := &models.Vaccine{
		FarmID: farm.ID,
		Name:   "Vacina Brucelose",
	}
	require.NoError(t, db.Create(vaccine1).Error)
	require.NoError(t, db.Create(vaccine2).Error)

	vaccines, err := repo.FindByFarmID(farm.ID)
	assert.NoError(t, err)
	assert.Len(t, vaccines, 2)
}

func TestVaccineRepository_FindByFarmID_Empty(t *testing.T) {
	db := setupVaccineTestDB(t)
	repo := repository.NewVaccineRepository(db)

	farm := createTestFarmForVaccine(t, db)

	vaccines, err := repo.FindByFarmID(farm.ID)
	assert.NoError(t, err)
	assert.Len(t, vaccines, 0)
}

func TestVaccineRepository_Update(t *testing.T) {
	db := setupVaccineTestDB(t)
	repo := repository.NewVaccineRepository(db)

	farm := createTestFarmForVaccine(t, db)

	vaccine := &models.Vaccine{
		FarmID: farm.ID,
		Name:   "Vacina Aftosa",
	}
	require.NoError(t, db.Create(vaccine).Error)

	vaccine.Name = "Vacina Aftosa Atualizada"
	vaccine.Description = "Nova descrição"
	vaccine.UpdatedAt = time.Now()

	err := repo.Update(vaccine)
	assert.NoError(t, err)

	updated, err := repo.FindByID(vaccine.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Vacina Aftosa Atualizada", updated.Name)
	assert.Equal(t, "Nova descrição", updated.Description)
}

func TestVaccineRepository_Delete(t *testing.T) {
	db := setupVaccineTestDB(t)
	repo := repository.NewVaccineRepository(db)

	farm := createTestFarmForVaccine(t, db)

	vaccine := &models.Vaccine{
		FarmID: farm.ID,
		Name:   "Vacina Aftosa",
	}
	require.NoError(t, db.Create(vaccine).Error)

	err := repo.Delete(vaccine.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(vaccine.ID)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestVaccineRepository_FindByFarmID_MultipleFarms(t *testing.T) {
	db := setupVaccineTestDB(t)
	repo := repository.NewVaccineRepository(db)

	farm1 := createTestFarmForVaccine(t, db)

	company2 := &models.Company{
		CompanyName: "Test Company 2",
		Location:    "Test Location 2",
		FarmCNPJ:    "98765432109876",
	}
	require.NoError(t, db.Create(company2).Error)
	farm2 := &models.Farm{
		CompanyID: company2.ID,
		Logo:      "",
	}
	require.NoError(t, db.Create(farm2).Error)

	vaccine1 := &models.Vaccine{
		FarmID: farm1.ID,
		Name:   "Vacina Aftosa",
	}
	vaccine2 := &models.Vaccine{
		FarmID: farm2.ID,
		Name:   "Vacina Brucelose",
	}
	require.NoError(t, db.Create(vaccine1).Error)
	require.NoError(t, db.Create(vaccine2).Error)

	vaccines, err := repo.FindByFarmID(farm1.ID)
	assert.NoError(t, err)
	assert.Len(t, vaccines, 1)
	assert.Equal(t, farm1.ID, vaccines[0].FarmID)
}
