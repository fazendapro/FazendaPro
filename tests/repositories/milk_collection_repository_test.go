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

func setupMilkCollectionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Company{}, &models.Farm{}, &models.Animal{}, &models.MilkCollection{})
	require.NoError(t, err)

	return db
}

func createTestFarmForMilkCollection(t *testing.T, db *gorm.DB) *models.Farm {
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

func createTestAnimalForMilkCollection(t *testing.T, db *gorm.DB, farmID uint) *models.Animal {
	animal := &models.Animal{
		FarmID:            farmID,
		EarTagNumberLocal: 123,
		AnimalName:        "Vaca Leiteira",
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

func TestMilkCollectionRepository_Create(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	farm := createTestFarmForMilkCollection(t, db)
	animal := createTestAnimalForMilkCollection(t, db, farm.ID)

	collectionDate := time.Now()
	milkCollection := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   25.5,
		Date:     collectionDate,
	}

	err := repo.Create(milkCollection)
	assert.NoError(t, err)
	assert.NotZero(t, milkCollection.ID)
}

func TestMilkCollectionRepository_FindByID(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	farm := createTestFarmForMilkCollection(t, db)
	animal := createTestAnimalForMilkCollection(t, db, farm.ID)

	collectionDate := time.Now()
	milkCollection := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   25.5,
		Date:     collectionDate,
	}
	require.NoError(t, db.Create(milkCollection).Error)

	found, err := repo.FindByID(milkCollection.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, milkCollection.ID, found.ID)
	assert.Equal(t, 25.5, found.Liters)
	assert.NotNil(t, found.Animal)
	assert.Equal(t, animal.ID, found.Animal.ID)
}

func TestMilkCollectionRepository_FindByID_NotFound(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	found, err := repo.FindByID(999)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestMilkCollectionRepository_FindByFarmID(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	farm1 := createTestFarmForMilkCollection(t, db)
	animal1 := createTestAnimalForMilkCollection(t, db, farm1.ID)

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
	animal2 := createTestAnimalForMilkCollection(t, db, farm2.ID)

	collection1 := &models.MilkCollection{
		AnimalID: animal1.ID,
		Liters:   25.5,
		Date:     time.Now().Add(-24 * time.Hour),
	}
	require.NoError(t, db.Create(collection1).Error)

	collection2 := &models.MilkCollection{
		AnimalID: animal1.ID,
		Liters:   30.0,
		Date:     time.Now(),
	}
	require.NoError(t, db.Create(collection2).Error)

	collection3 := &models.MilkCollection{
		AnimalID: animal2.ID,
		Liters:   20.0,
		Date:     time.Now(),
	}
	require.NoError(t, db.Create(collection3).Error)

	collections, err := repo.FindByFarmID(farm1.ID)
	assert.NoError(t, err)
	assert.Len(t, collections, 2)
	assert.True(t, collections[0].Date.After(collections[1].Date) || collections[0].Date.Equal(collections[1].Date))
}

func TestMilkCollectionRepository_FindByFarmIDWithDateRange(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	farm := createTestFarmForMilkCollection(t, db)
	animal := createTestAnimalForMilkCollection(t, db, farm.ID)

	now := time.Now()
	startDate := now.Add(-7 * 24 * time.Hour)
	endDate := now.Add(24 * time.Hour)

	collection1 := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   25.5,
		Date:     now.Add(-3 * 24 * time.Hour),
	}
	require.NoError(t, db.Create(collection1).Error)

	collection2 := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   20.0,
		Date:     now.Add(-10 * 24 * time.Hour),
	}
	require.NoError(t, db.Create(collection2).Error)

	collection3 := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   30.0,
		Date:     now,
	}
	require.NoError(t, db.Create(collection3).Error)

	collections, err := repo.FindByFarmIDWithDateRange(farm.ID, &startDate, &endDate)
	assert.NoError(t, err)
	assert.Len(t, collections, 2)
}

func TestMilkCollectionRepository_FindByFarmIDWithDateRange_StartDateOnly(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	farm := createTestFarmForMilkCollection(t, db)
	animal := createTestAnimalForMilkCollection(t, db, farm.ID)

	now := time.Now()
	startDate := now.Add(-7 * 24 * time.Hour)

	collection1 := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   25.5,
		Date:     now.Add(-3 * 24 * time.Hour),
	}
	require.NoError(t, db.Create(collection1).Error)

	collection2 := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   20.0,
		Date:     now.Add(-10 * 24 * time.Hour),
	}
	require.NoError(t, db.Create(collection2).Error)

	collections, err := repo.FindByFarmIDWithDateRange(farm.ID, &startDate, nil)
	assert.NoError(t, err)
	assert.Len(t, collections, 1)
}

func TestMilkCollectionRepository_FindByFarmIDWithDateRange_EndDateOnly(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	farm := createTestFarmForMilkCollection(t, db)
	animal := createTestAnimalForMilkCollection(t, db, farm.ID)

	now := time.Now()
	endDate := now.Add(24 * time.Hour)

	collection1 := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   25.5,
		Date:     now.Add(-3 * 24 * time.Hour),
	}
	require.NoError(t, db.Create(collection1).Error)

	collection2 := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   20.0,
		Date:     now.Add(2 * 24 * time.Hour),
	}
	require.NoError(t, db.Create(collection2).Error)

	collections, err := repo.FindByFarmIDWithDateRange(farm.ID, nil, &endDate)
	assert.NoError(t, err)
	assert.Len(t, collections, 1)
}

func TestMilkCollectionRepository_FindByFarmIDWithDateRange_NoDates(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	farm := createTestFarmForMilkCollection(t, db)
	animal := createTestAnimalForMilkCollection(t, db, farm.ID)

	collection1 := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   25.5,
		Date:     time.Now(),
	}
	require.NoError(t, db.Create(collection1).Error)

	collections, err := repo.FindByFarmIDWithDateRange(farm.ID, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, collections, 1)
}

func TestMilkCollectionRepository_FindByAnimalID(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	farm := createTestFarmForMilkCollection(t, db)
	animal1 := createTestAnimalForMilkCollection(t, db, farm.ID)

	animal2 := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 124,
		AnimalName:        "Vaca Leiteira 2",
		Sex:               0,
		Breed:             "Holstein",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal2).Error)

	collection1 := &models.MilkCollection{
		AnimalID: animal1.ID,
		Liters:   25.5,
		Date:     time.Now().Add(-24 * time.Hour),
	}
	require.NoError(t, db.Create(collection1).Error)

	collection2 := &models.MilkCollection{
		AnimalID: animal1.ID,
		Liters:   30.0,
		Date:     time.Now(),
	}
	require.NoError(t, db.Create(collection2).Error)

	collection3 := &models.MilkCollection{
		AnimalID: animal2.ID,
		Liters:   20.0,
		Date:     time.Now(),
	}
	require.NoError(t, db.Create(collection3).Error)

	collections, err := repo.FindByAnimalID(animal1.ID)
	assert.NoError(t, err)
	assert.Len(t, collections, 2)
	for _, collection := range collections {
		assert.Equal(t, animal1.ID, collection.AnimalID)
		assert.NotNil(t, collection.Animal)
	}
	assert.True(t, collections[0].Date.After(collections[1].Date) || collections[0].Date.Equal(collections[1].Date))
}

func TestMilkCollectionRepository_Update(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	farm := createTestFarmForMilkCollection(t, db)
	animal := createTestAnimalForMilkCollection(t, db, farm.ID)

	collectionDate := time.Now()
	milkCollection := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   25.5,
		Date:     collectionDate,
	}
	require.NoError(t, db.Create(milkCollection).Error)

	milkCollection.Liters = 30.0
	newDate := time.Now().Add(24 * time.Hour)
	milkCollection.Date = newDate

	err := repo.Update(milkCollection)
	assert.NoError(t, err)

	updated, err := repo.FindByID(milkCollection.ID)
	assert.NoError(t, err)
	assert.Equal(t, 30.0, updated.Liters)
	assert.True(t, updated.Date.Equal(newDate))
}

func TestMilkCollectionRepository_Update_NotFound(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	farm := createTestFarmForMilkCollection(t, db)
	animal := createTestAnimalForMilkCollection(t, db, farm.ID)

	milkCollection := &models.MilkCollection{
		ID:       999,
		AnimalID: animal.ID,
		Liters:   25.5,
		Date:     time.Now(),
	}

	err := repo.Update(milkCollection)
	assert.Error(t, err)
}

func TestMilkCollectionRepository_Delete(t *testing.T) {
	db := setupMilkCollectionTestDB(t)
	repo := repository.NewMilkCollectionRepository(db)

	farm := createTestFarmForMilkCollection(t, db)
	animal := createTestAnimalForMilkCollection(t, db, farm.ID)

	collectionDate := time.Now()
	milkCollection := &models.MilkCollection{
		AnimalID: animal.ID,
		Liters:   25.5,
		Date:     collectionDate,
	}
	require.NoError(t, db.Create(milkCollection).Error)

	err := repo.Delete(milkCollection.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(milkCollection.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}
