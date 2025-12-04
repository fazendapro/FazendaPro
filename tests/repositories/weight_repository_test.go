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

func setupWeightTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Company{}, &models.Farm{}, &models.Animal{}, &models.Weight{})
	require.NoError(t, err)

	return db
}

func createTestFarmForWeight(t *testing.T, db *gorm.DB) *models.Farm {
	company := &models.Company{
		CompanyName: "Test Company",
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

func TestWeightRepository_Create(t *testing.T) {
	db := setupWeightTestDB(t)
	weightRepo := repository.NewWeightRepository(db)

	farm := createTestFarmForWeight(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Boi João",
		EarTagNumberLocal: 123,
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	weightDate := time.Now()
	weight := &models.Weight{
		AnimalID:     animal.ID,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}

	err := weightRepo.Create(weight)
	assert.NoError(t, err)
	assert.NotZero(t, weight.ID)
}

func TestWeightRepository_FindByID(t *testing.T) {
	db := setupWeightTestDB(t)
	weightRepo := repository.NewWeightRepository(db)

	farm := createTestFarmForWeight(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Boi João",
		EarTagNumberLocal: 123,
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	weightDate := time.Now()
	weight := &models.Weight{
		AnimalID:     animal.ID,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}
	require.NoError(t, db.Create(weight).Error)

	result, err := weightRepo.FindByID(weight.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, weight.ID, result.ID)
	assert.Equal(t, animal.ID, result.AnimalID)
	assert.Equal(t, 450.5, result.AnimalWeight)
	assert.NotNil(t, result.Animal)
	assert.Equal(t, "Boi João", result.Animal.AnimalName)
}

func TestWeightRepository_FindByID_NotFound(t *testing.T) {
	db := setupWeightTestDB(t)
	weightRepo := repository.NewWeightRepository(db)

	result, err := weightRepo.FindByID(999)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestWeightRepository_FindByAnimalID(t *testing.T) {
	db := setupWeightTestDB(t)
	weightRepo := repository.NewWeightRepository(db)

	farm := createTestFarmForWeight(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Boi João",
		EarTagNumberLocal: 123,
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	weightDate := time.Now()
	weight := &models.Weight{
		AnimalID:     animal.ID,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}
	require.NoError(t, db.Create(weight).Error)

	result, err := weightRepo.FindByAnimalID(animal.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, animal.ID, result.AnimalID)
	assert.Equal(t, 450.5, result.AnimalWeight)
}

func TestWeightRepository_FindByAnimalID_NotFound(t *testing.T) {
	db := setupWeightTestDB(t)
	weightRepo := repository.NewWeightRepository(db)

	result, err := weightRepo.FindByAnimalID(999)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestWeightRepository_FindByAnimalID_MostRecent(t *testing.T) {
	db := setupWeightTestDB(t)
	weightRepo := repository.NewWeightRepository(db)

	farm := createTestFarmForWeight(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Boi João",
		EarTagNumberLocal: 123,
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	oldDate := time.Now().AddDate(0, -1, 0)
	oldWeight := &models.Weight{
		AnimalID:     animal.ID,
		Date:         oldDate,
		AnimalWeight: 400.0,
	}
	require.NoError(t, db.Create(oldWeight).Error)

	newDate := time.Now()
	newWeight := &models.Weight{
		AnimalID:     animal.ID,
		Date:         newDate,
		AnimalWeight: 450.5,
	}
	require.NoError(t, db.Create(newWeight).Error)

	result, err := weightRepo.FindByAnimalID(animal.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 450.5, result.AnimalWeight)
	assert.True(t, result.Date.After(oldDate))
}

func TestWeightRepository_FindByFarmID(t *testing.T) {
	db := setupWeightTestDB(t)
	weightRepo := repository.NewWeightRepository(db)

	farm := createTestFarmForWeight(t, db)

	animal1 := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Boi João",
		EarTagNumberLocal: 123,
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal1).Error)

	animal2 := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Boi Pedro",
		EarTagNumberLocal: 124,
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal2).Error)

	weightDate := time.Now()
	weight1 := &models.Weight{
		AnimalID:     animal1.ID,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}
	require.NoError(t, db.Create(weight1).Error)

	weight2 := &models.Weight{
		AnimalID:     animal2.ID,
		Date:         weightDate,
		AnimalWeight: 500.0,
	}
	require.NoError(t, db.Create(weight2).Error)

	result, err := weightRepo.FindByFarmID(farm.ID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestWeightRepository_Update(t *testing.T) {
	db := setupWeightTestDB(t)
	weightRepo := repository.NewWeightRepository(db)

	farm := createTestFarmForWeight(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Boi João",
		EarTagNumberLocal: 123,
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	weightDate := time.Now()
	weight := &models.Weight{
		AnimalID:     animal.ID,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}
	require.NoError(t, db.Create(weight).Error)

	weight.AnimalWeight = 480.0
	err := weightRepo.Update(weight)

	assert.NoError(t, err)

	updated, err := weightRepo.FindByID(weight.ID)
	assert.NoError(t, err)
	assert.Equal(t, 480.0, updated.AnimalWeight)
}

func TestWeightRepository_Delete(t *testing.T) {
	db := setupWeightTestDB(t)
	weightRepo := repository.NewWeightRepository(db)

	farm := createTestFarmForWeight(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Boi João",
		EarTagNumberLocal: 123,
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	weightDate := time.Now()
	weight := &models.Weight{
		AnimalID:     animal.ID,
		Date:         weightDate,
		AnimalWeight: 450.5,
	}
	require.NoError(t, db.Create(weight).Error)

	err := weightRepo.Delete(weight.ID)
	assert.NoError(t, err)

	result, err := weightRepo.FindByID(weight.ID)
	assert.NoError(t, err)
	assert.Nil(t, result)
}
