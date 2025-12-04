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

func setupVaccineApplicationTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Company{}, &models.Farm{}, &models.Animal{}, &models.Vaccine{}, &models.VaccineApplication{})
	require.NoError(t, err)

	return db
}

func createTestFarmForVaccineApplication(t *testing.T, db *gorm.DB) *models.Farm {
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

func createTestAnimalForVaccineApplication(t *testing.T, db *gorm.DB, farmID uint) *models.Animal {
	animal := &models.Animal{
		FarmID:            farmID,
		EarTagNumberLocal: 123,
		AnimalName:        "Vaca Teste",
		Sex:               0,
		Breed:             "Holstein",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)
	return animal
}

func createTestVaccineForVaccineApplication(t *testing.T, db *gorm.DB, farmID uint) *models.Vaccine {
	vaccine := &models.Vaccine{
		FarmID:       farmID,
		Name:         "Vacina Aftosa",
		Description:  "Vacina contra febre aftosa",
		Manufacturer: "Fabricante XYZ",
	}
	require.NoError(t, db.Create(vaccine).Error)
	return vaccine
}

func TestVaccineApplicationRepository_Create(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate,
		BatchNumber:     "LOTE123",
		Veterinarian:    "Dr. João Silva",
		Observations:    "Aplicação realizada com sucesso",
	}

	err := repo.Create(vaccineApplication)
	assert.NoError(t, err)
	assert.NotZero(t, vaccineApplication.ID)
}

func TestVaccineApplicationRepository_FindByID(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate,
		BatchNumber:     "LOTE123",
	}
	require.NoError(t, db.Create(vaccineApplication).Error)

	found, err := repo.FindByID(vaccineApplication.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, vaccineApplication.ID, found.ID)
	assert.NotNil(t, found.Animal)
	assert.Equal(t, animal.ID, found.Animal.ID)
	assert.NotNil(t, found.Vaccine)
	assert.Equal(t, vaccine.ID, found.Vaccine.ID)
}

func TestVaccineApplicationRepository_FindByID_NotFound(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	found, err := repo.FindByID(999)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestVaccineApplicationRepository_FindByFarmID(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal1 := createTestAnimalForVaccineApplication(t, db, farm.ID)
	animal2 := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	applicationDate := time.Now()
	vaccineApplication1 := &models.VaccineApplication{
		AnimalID:        animal1.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate,
	}
	vaccineApplication2 := &models.VaccineApplication{
		AnimalID:        animal2.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate,
	}
	require.NoError(t, db.Create(vaccineApplication1).Error)
	require.NoError(t, db.Create(vaccineApplication2).Error)

	applications, err := repo.FindByFarmID(farm.ID)
	assert.NoError(t, err)
	assert.Len(t, applications, 2)
}

func TestVaccineApplicationRepository_FindByFarmIDWithDateRange(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	now := time.Now()
	startDate := now.AddDate(0, 0, -30)
	endDate := now

	applicationDate1 := now.AddDate(0, 0, -10) // Dentro do range
	applicationDate2 := now.AddDate(0, 0, -40) // Fora do range

	vaccineApplication1 := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate1,
	}
	vaccineApplication2 := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate2,
	}
	require.NoError(t, db.Create(vaccineApplication1).Error)
	require.NoError(t, db.Create(vaccineApplication2).Error)

	applications, err := repo.FindByFarmIDWithDateRange(farm.ID, &startDate, &endDate)
	assert.NoError(t, err)
	assert.Len(t, applications, 1)
	assert.Equal(t, vaccineApplication1.ID, applications[0].ID)
}

func TestVaccineApplicationRepository_FindByFarmIDWithDateRange_StartDateOnly(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	now := time.Now()
	startDate := now.AddDate(0, 0, -30)

	applicationDate1 := now.AddDate(0, 0, -10) // Dentro do range
	applicationDate2 := now.AddDate(0, 0, -40) // Fora do range

	vaccineApplication1 := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate1,
	}
	vaccineApplication2 := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate2,
	}
	require.NoError(t, db.Create(vaccineApplication1).Error)
	require.NoError(t, db.Create(vaccineApplication2).Error)

	applications, err := repo.FindByFarmIDWithDateRange(farm.ID, &startDate, nil)
	assert.NoError(t, err)
	assert.Len(t, applications, 1)
	assert.Equal(t, vaccineApplication1.ID, applications[0].ID)
}

func TestVaccineApplicationRepository_FindByFarmIDWithDateRange_EndDateOnly(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	now := time.Now()
	endDate := now.AddDate(0, 0, -20)

	applicationDate1 := now.AddDate(0, 0, -10) // Dentro do range
	applicationDate2 := now.AddDate(0, 0, -30) // Fora do range

	vaccineApplication1 := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate1,
	}
	vaccineApplication2 := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate2,
	}
	require.NoError(t, db.Create(vaccineApplication1).Error)
	require.NoError(t, db.Create(vaccineApplication2).Error)

	applications, err := repo.FindByFarmIDWithDateRange(farm.ID, nil, &endDate)
	assert.NoError(t, err)
	assert.Len(t, applications, 1)
	assert.Equal(t, vaccineApplication2.ID, applications[0].ID)
}

func TestVaccineApplicationRepository_FindByFarmIDWithDateRange_NoDates(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate,
	}
	require.NoError(t, db.Create(vaccineApplication).Error)

	applications, err := repo.FindByFarmIDWithDateRange(farm.ID, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, applications, 1)
}

func TestVaccineApplicationRepository_FindByAnimalID(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	applicationDate := time.Now()
	vaccineApplication1 := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate,
	}
	vaccineApplication2 := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate,
	}
	require.NoError(t, db.Create(vaccineApplication1).Error)
	require.NoError(t, db.Create(vaccineApplication2).Error)

	applications, err := repo.FindByAnimalID(animal.ID)
	assert.NoError(t, err)
	assert.Len(t, applications, 2)
}

func TestVaccineApplicationRepository_FindByVaccineID(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal1 := createTestAnimalForVaccineApplication(t, db, farm.ID)
	animal2 := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	applicationDate := time.Now()
	vaccineApplication1 := &models.VaccineApplication{
		AnimalID:        animal1.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate,
	}
	vaccineApplication2 := &models.VaccineApplication{
		AnimalID:        animal2.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate,
	}
	require.NoError(t, db.Create(vaccineApplication1).Error)
	require.NoError(t, db.Create(vaccineApplication2).Error)

	applications, err := repo.FindByVaccineID(vaccine.ID)
	assert.NoError(t, err)
	assert.Len(t, applications, 2)
}

func TestVaccineApplicationRepository_Update(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate,
		BatchNumber:     "LOTE123",
	}
	require.NoError(t, db.Create(vaccineApplication).Error)

	vaccineApplication.BatchNumber = "LOTE456"
	vaccineApplication.Veterinarian = "Dr. Maria Silva"
	vaccineApplication.UpdatedAt = time.Now()

	err := repo.Update(vaccineApplication)
	assert.NoError(t, err)

	updated, err := repo.FindByID(vaccineApplication.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "LOTE456", updated.BatchNumber)
	assert.Equal(t, "Dr. Maria Silva", updated.Veterinarian)
}

func TestVaccineApplicationRepository_Delete(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		AnimalID:        animal.ID,
		VaccineID:       vaccine.ID,
		ApplicationDate: applicationDate,
	}
	require.NoError(t, db.Create(vaccineApplication).Error)

	err := repo.Delete(vaccineApplication.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(vaccineApplication.ID)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestVaccineApplicationRepository_FindByFarmID_MultipleFarms(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm1 := createTestFarmForVaccineApplication(t, db)

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

	animal1 := createTestAnimalForVaccineApplication(t, db, farm1.ID)
	animal2 := createTestAnimalForVaccineApplication(t, db, farm2.ID)
	vaccine1 := createTestVaccineForVaccineApplication(t, db, farm1.ID)
	vaccine2 := createTestVaccineForVaccineApplication(t, db, farm2.ID)

	applicationDate := time.Now()
	vaccineApplication1 := &models.VaccineApplication{
		AnimalID:        animal1.ID,
		VaccineID:       vaccine1.ID,
		ApplicationDate: applicationDate,
	}
	vaccineApplication2 := &models.VaccineApplication{
		AnimalID:        animal2.ID,
		VaccineID:       vaccine2.ID,
		ApplicationDate: applicationDate,
	}
	require.NoError(t, db.Create(vaccineApplication1).Error)
	require.NoError(t, db.Create(vaccineApplication2).Error)

	applications, err := repo.FindByFarmID(farm1.ID)
	assert.NoError(t, err)
	assert.Len(t, applications, 1)
	assert.Equal(t, animal1.ID, applications[0].AnimalID)
}

func TestVaccineApplicationRepository_FindByAnimalID_Empty(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal := createTestAnimalForVaccineApplication(t, db, farm.ID)

	applications, err := repo.FindByAnimalID(animal.ID)
	assert.NoError(t, err)
	assert.Len(t, applications, 0)
}

func TestVaccineApplicationRepository_FindByFarmIDWithPagination(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	for i := 0; i < 10; i++ {
		animal := createTestAnimalForVaccineApplication(t, db, farm.ID)
		applicationDate := time.Now().Add(-time.Duration(i) * 24 * time.Hour)
		vaccineApplication := &models.VaccineApplication{
			AnimalID:        animal.ID,
			VaccineID:       vaccine.ID,
			ApplicationDate: applicationDate,
		}
		require.NoError(t, db.Create(vaccineApplication).Error)
	}

	applications, total, err := repo.FindByFarmIDWithPagination(farm.ID, 1, 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
	assert.Len(t, applications, 5)

	applications2, total2, err := repo.FindByFarmIDWithPagination(farm.ID, 2, 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total2)
	assert.Len(t, applications2, 5)

	assert.NotEqual(t, applications[0].ID, applications2[0].ID)
}

func TestVaccineApplicationRepository_FindByFarmIDWithPagination_EmptyResult(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)

	applications, total, err := repo.FindByFarmIDWithPagination(farm.ID, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, applications, 0)
}

func TestVaccineApplicationRepository_FindByFarmIDWithDateRangePaginated(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	animal := createTestAnimalForVaccineApplication(t, db, farm.ID)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	now := time.Now()
	startDate := now.AddDate(0, 0, -30)
	endDate := now

	// Criar 10 aplicações, 5 dentro do range e 5 fora
	for i := 0; i < 10; i++ {
		var applicationDate time.Time
		if i < 5 {
			applicationDate = now.AddDate(0, 0, -10-i) // Dentro do range
		} else {
			applicationDate = now.AddDate(0, 0, -40-i) // Fora do range
		}

		vaccineApplication := &models.VaccineApplication{
			AnimalID:        animal.ID,
			VaccineID:       vaccine.ID,
			ApplicationDate: applicationDate,
		}
		require.NoError(t, db.Create(vaccineApplication).Error)
	}

	applications, total, err := repo.FindByFarmIDWithDateRangePaginated(farm.ID, &startDate, &endDate, 1, 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, applications, 3)

	applications2, total2, err := repo.FindByFarmIDWithDateRangePaginated(farm.ID, &startDate, &endDate, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total2)
	assert.Len(t, applications2, 2)
}

func TestVaccineApplicationRepository_FindByVaccineID_Empty(t *testing.T) {
	db := setupVaccineApplicationTestDB(t)
	repo := repository.NewVaccineApplicationRepository(db)

	farm := createTestFarmForVaccineApplication(t, db)
	vaccine := createTestVaccineForVaccineApplication(t, db, farm.ID)

	applications, err := repo.FindByVaccineID(vaccine.ID)
	assert.NoError(t, err)
	assert.Len(t, applications, 0)
}
