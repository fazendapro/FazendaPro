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

func setupAnimalTestDB(t *testing.T) (*repository.Database, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Company{}, &models.Farm{}, &models.Animal{})
	require.NoError(t, err)

	database := &repository.Database{DB: db}
	return database, db
}

func createTestFarmForAnimal(t *testing.T, db *gorm.DB) *models.Farm {
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

func TestAnimalRepository_Create(t *testing.T) {
	database, db := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	farm := createTestFarmForAnimal(t, db)

	birthDate := time.Now().AddDate(-2, 0, 0)
	animal := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 123,
		AnimalName:        "Boi João",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		BirthDate:         &birthDate,
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}

	err := repo.Create(animal)
	assert.NoError(t, err)
	assert.NotZero(t, animal.ID)
}

func TestAnimalRepository_FindByID(t *testing.T) {
	database, db := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	farm := createTestFarmForAnimal(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 123,
		AnimalName:        "Boi João",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	found, err := repo.FindByID(animal.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, animal.ID, found.ID)
	assert.Equal(t, "Boi João", found.AnimalName)
}

func TestAnimalRepository_FindByID_NotFound(t *testing.T) {
	database, _ := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	found, err := repo.FindByID(999)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestAnimalRepository_FindByID_WithParents(t *testing.T) {
	database, db := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	farm := createTestFarmForAnimal(t, db)

	father := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 100,
		AnimalName:        "Pai",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(father).Error)

	mother := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 101,
		AnimalName:        "Mãe",
		Sex:               0,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(mother).Error)

	child := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 102,
		AnimalName:        "Filho",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		FatherID:          &father.ID,
		MotherID:          &mother.ID,
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(child).Error)

	found, err := repo.FindByID(child.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.NotNil(t, found.Father)
	assert.NotNil(t, found.Mother)
	assert.Equal(t, father.ID, found.Father.ID)
	assert.Equal(t, mother.ID, found.Mother.ID)
}

func TestAnimalRepository_FindByFarmID(t *testing.T) {
	database, db := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	farm1 := createTestFarmForAnimal(t, db)

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

	animal1 := &models.Animal{
		FarmID:            farm1.ID,
		EarTagNumberLocal: 123,
		AnimalName:        "Animal 1",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal1).Error)

	animal2 := &models.Animal{
		FarmID:            farm1.ID,
		EarTagNumberLocal: 124,
		AnimalName:        "Animal 2",
		Sex:               0,
		Breed:             "Holstein",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal2).Error)

	animal3 := &models.Animal{
		FarmID:            farm2.ID,
		EarTagNumberLocal: 125,
		AnimalName:        "Animal 3",
		Sex:               1,
		Breed:             "Angus",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal3).Error)

	animals, err := repo.FindByFarmID(farm1.ID)
	assert.NoError(t, err)
	assert.Len(t, animals, 2)
	assert.Equal(t, farm1.ID, animals[0].FarmID)
	assert.Equal(t, farm1.ID, animals[1].FarmID)
}

func TestAnimalRepository_FindByEarTagNumber(t *testing.T) {
	database, db := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	farm := createTestFarmForAnimal(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 123,
		AnimalName:        "Boi João",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	found, err := repo.FindByEarTagNumber(farm.ID, 123)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, 123, found.EarTagNumberLocal)
	assert.Equal(t, "Boi João", found.AnimalName)
}

func TestAnimalRepository_FindByEarTagNumber_NotFound(t *testing.T) {
	database, db := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	farm := createTestFarmForAnimal(t, db)

	found, err := repo.FindByEarTagNumber(farm.ID, 999)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestAnimalRepository_FindByEarTagNumber_WrongFarm(t *testing.T) {
	database, db := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	farm1 := createTestFarmForAnimal(t, db)

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

	animal := &models.Animal{
		FarmID:            farm1.ID,
		EarTagNumberLocal: 123,
		AnimalName:        "Boi João",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	found, err := repo.FindByEarTagNumber(farm2.ID, 123)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestAnimalRepository_Update(t *testing.T) {
	database, db := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	farm := createTestFarmForAnimal(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 123,
		AnimalName:        "Boi João",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	animal.AnimalName = "Boi João Atualizado"
	animal.Breed = "Holstein"

	err := repo.Update(animal)
	assert.NoError(t, err)

	updated, err := repo.FindByID(animal.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Boi João Atualizado", updated.AnimalName)
	assert.Equal(t, "Holstein", updated.Breed)
}

func TestAnimalRepository_Delete(t *testing.T) {
	database, db := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	farm := createTestFarmForAnimal(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 123,
		AnimalName:        "Boi João",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	err := repo.Delete(animal.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(animal.ID)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestAnimalRepository_FindByFarmIDAndSex(t *testing.T) {
	database, db := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	farm := createTestFarmForAnimal(t, db)

	male1 := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 100,
		AnimalName:        "Macho 1",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(male1).Error)

	male2 := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 101,
		AnimalName:        "Macho 2",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(male2).Error)

	female1 := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 200,
		AnimalName:        "Fêmea 1",
		Sex:               0,
		Breed:             "Holstein",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(female1).Error)

	males, err := repo.FindByFarmIDAndSex(farm.ID, 1)
	assert.NoError(t, err)
	assert.Len(t, males, 2)
	for _, animal := range males {
		assert.Equal(t, 1, animal.Sex)
	}

	females, err := repo.FindByFarmIDAndSex(farm.ID, 0)
	assert.NoError(t, err)
	assert.Len(t, females, 1)
	assert.Equal(t, 0, females[0].Sex)
}

func TestAnimalRepository_CountBySex(t *testing.T) {
	database, db := setupAnimalTestDB(t)
	repo := repository.NewAnimalRepository(database)

	farm := createTestFarmForAnimal(t, db)

	male1 := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 100,
		AnimalName:        "Macho 1",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(male1).Error)

	male2 := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 101,
		AnimalName:        "Macho 2",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(male2).Error)

	female1 := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 200,
		AnimalName:        "Fêmea 1",
		Sex:               0,
		Breed:             "Holstein",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(female1).Error)

	female2 := &models.Animal{
		FarmID:            farm.ID,
		EarTagNumberLocal: 201,
		AnimalName:        "Fêmea 2",
		Sex:               0,
		Breed:             "Holstein",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(female2).Error)

	maleCount, err := repo.CountBySex(farm.ID, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), maleCount)

	femaleCount, err := repo.CountBySex(farm.ID, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), femaleCount)

	company2 := &models.Company{
		CompanyName: "Test Company 2",
		Location:    "Test Location 2",
		FarmCNPJ:    "98765432109877",
	}
	require.NoError(t, db.Create(company2).Error)
	farm2 := &models.Farm{
		CompanyID: company2.ID,
		Logo:      "",
	}
	require.NoError(t, db.Create(farm2).Error)

	emptyCount, err := repo.CountBySex(farm2.ID, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), emptyCount)
}
