package repositories

import (
	"fmt"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupReproductionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Company{}, &models.Farm{}, &models.Animal{}, &models.Reproduction{})
	require.NoError(t, err)

	return db
}

func createTestFarmForReproduction(t *testing.T, db *gorm.DB) *models.Farm {
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

func TestReproductionRepository_FindByPhase(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm := createTestFarmForReproduction(t, db)

	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200)

	animal1 := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Tata Salt",
		EarTagNumberLocal: 123,
		Photo:             "src/assets/images/mocked/cows/tata.png",
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
	}
	require.NoError(t, db.Create(animal1).Error)

	animal2 := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Lays",
		EarTagNumberLocal: 124,
		Photo:             "src/assets/images/mocked/cows/lays.png",
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
	}
	require.NoError(t, db.Create(animal2).Error)

	reproduction1 := &models.Reproduction{
		AnimalID:      animal1.ID,
		CurrentPhase:  models.PhasePrenhas,
		PregnancyDate: &pregnancyDate,
	}
	require.NoError(t, db.Create(reproduction1).Error)

	reproduction2 := &models.Reproduction{
		AnimalID:      animal2.ID,
		CurrentPhase:  models.PhasePrenhas,
		PregnancyDate: &pregnancyDate,
	}
	require.NoError(t, db.Create(reproduction2).Error)

	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, models.PhasePrenhas, result[0].CurrentPhase)
	assert.Equal(t, models.PhasePrenhas, result[1].CurrentPhase)
	assert.Equal(t, "Tata Salt", result[0].Animal.AnimalName)
	assert.Equal(t, "Lays", result[1].Animal.AnimalName)
}

func TestReproductionRepository_FindByPhase_Error(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestReproductionRepository_FindByPhase_EmptyResults(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	animal := &models.Animal{
		FarmID:            1,
		AnimalName:        "Tata Salt",
		EarTagNumberLocal: 123,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
	}
	require.NoError(t, db.Create(animal).Error)

	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestReproductionRepository_FindByPhase_DifferentPhases(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm := createTestFarmForReproduction(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Tata Salt",
		EarTagNumberLocal: 123,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
	}
	require.NoError(t, db.Create(animal).Error)

	reproduction := &models.Reproduction{
		AnimalID:     animal.ID,
		CurrentPhase: models.PhaseLactacao,
	}
	require.NoError(t, db.Create(reproduction).Error)

	require.NoError(t, db.Model(reproduction).Update("CurrentPhase", models.PhaseLactacao).Error)

	var checkReproduction models.Reproduction
	require.NoError(t, db.First(&checkReproduction, reproduction.ID).Error)
	assert.Equal(t, models.PhaseLactacao, checkReproduction.CurrentPhase)

	result, err := reproductionRepo.FindByPhase(models.PhaseLactacao)

	assert.NoError(t, err)
	if assert.Len(t, result, 1, "Deveria encontrar 1 reprodução na fase de lactação") {
		assert.Equal(t, models.PhaseLactacao, result[0].CurrentPhase)
	}
}

func TestReproductionRepository_FindByPhase_WithAnimalData(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm := createTestFarmForReproduction(t, db)

	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200)

	animal := &models.Animal{
		FarmID:               farm.ID,
		AnimalName:           "Tata Salt",
		EarTagNumberLocal:    123,
		EarTagNumberRegister: 456,
		Photo:                "src/assets/images/mocked/cows/tata.png",
		Sex:                  1,
		Breed:                "Holandesa",
		Type:                 "vaca",
	}
	require.NoError(t, db.Create(animal).Error)

	reproduction := &models.Reproduction{
		AnimalID:      animal.ID,
		CurrentPhase:  models.PhasePrenhas,
		PregnancyDate: &pregnancyDate,
	}
	require.NoError(t, db.Create(reproduction).Error)
	require.NoError(t, db.Model(reproduction).Update("CurrentPhase", models.PhasePrenhas).Error)

	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	reproductionResult := result[0]
	assert.Equal(t, uint(1), reproductionResult.AnimalID)
	assert.Equal(t, models.PhasePrenhas, reproductionResult.CurrentPhase)
	assert.NotNil(t, reproductionResult.PregnancyDate)
	assert.Equal(t, pregnancyDate.Format("2006-01-02 15:04:05"), reproductionResult.PregnancyDate.Format("2006-01-02 15:04:05"))

	animalResult := reproductionResult.Animal
	assert.Equal(t, uint(1), animalResult.ID)
	assert.Equal(t, uint(1), animalResult.FarmID)
	assert.Equal(t, "Tata Salt", animalResult.AnimalName)
	assert.Equal(t, 123, animalResult.EarTagNumberLocal)
	assert.Equal(t, 456, animalResult.EarTagNumberRegister)
	assert.Equal(t, "src/assets/images/mocked/cows/tata.png", animalResult.Photo)
	assert.Equal(t, 1, animalResult.Sex)
	assert.Equal(t, "Holandesa", animalResult.Breed)
	assert.Equal(t, "vaca", animalResult.Type)
}

func TestReproductionRepository_FindByPhase_WithNullPregnancyDate(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm := createTestFarmForReproduction(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Tata Salt",
		EarTagNumberLocal: 123,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
	}
	require.NoError(t, db.Create(animal).Error)

	reproduction := &models.Reproduction{
		AnimalID:      animal.ID,
		CurrentPhase:  models.PhasePrenhas,
		PregnancyDate: nil,
	}
	require.NoError(t, db.Create(reproduction).Error)

	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	reproductionResult := result[0]
	assert.Equal(t, models.PhasePrenhas, reproductionResult.CurrentPhase)
	assert.Nil(t, reproductionResult.PregnancyDate)
}

func TestReproductionRepository_FindByPhase_WithMultipleFarms(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm1 := createTestFarmForReproduction(t, db)
	company2 := &models.Company{
		CompanyName: "Test Company 2",
		FarmCNPJ:    "98765432109876",
	}
	require.NoError(t, db.Create(company2).Error)
	farm2 := &models.Farm{
		CompanyID: company2.ID,
		Logo:      "",
	}
	require.NoError(t, db.Create(farm2).Error)

	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200)

	animal1 := &models.Animal{
		FarmID:            farm1.ID,
		AnimalName:        "Tata Salt",
		EarTagNumberLocal: 123,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
	}
	require.NoError(t, db.Create(animal1).Error)

	animal2 := &models.Animal{
		FarmID:            farm2.ID,
		AnimalName:        "Lays",
		EarTagNumberLocal: 124,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
	}
	require.NoError(t, db.Create(animal2).Error)

	animal3 := &models.Animal{
		FarmID:            farm1.ID,
		AnimalName:        "Matilda",
		EarTagNumberLocal: 125,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
	}
	require.NoError(t, db.Create(animal3).Error)

	reproduction1 := &models.Reproduction{
		AnimalID:      animal1.ID,
		CurrentPhase:  models.PhasePrenhas,
		PregnancyDate: &pregnancyDate,
	}
	require.NoError(t, db.Create(reproduction1).Error)

	reproduction2 := &models.Reproduction{
		AnimalID:      animal2.ID,
		CurrentPhase:  models.PhasePrenhas,
		PregnancyDate: &pregnancyDate,
	}
	require.NoError(t, db.Create(reproduction2).Error)

	reproduction3 := &models.Reproduction{
		AnimalID:      animal3.ID,
		CurrentPhase:  models.PhasePrenhas,
		PregnancyDate: &pregnancyDate,
	}
	require.NoError(t, db.Create(reproduction3).Error)

	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	assert.NoError(t, err)
	assert.Len(t, result, 3)

	for _, reproduction := range result {
		assert.Equal(t, models.PhasePrenhas, reproduction.CurrentPhase)
		assert.NotNil(t, reproduction.PregnancyDate)
	}

	farm1Count := 0
	farm2Count := 0
	for _, reproduction := range result {
		if reproduction.Animal.FarmID == farm1.ID {
			farm1Count++
		} else if reproduction.Animal.FarmID == farm2.ID {
			farm2Count++
		}
	}
	assert.Equal(t, 2, farm1Count)
	assert.Equal(t, 1, farm2Count)
}

func TestReproductionRepository_Create(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm := createTestFarmForReproduction(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Tata Salt",
		EarTagNumberLocal: 123,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	reproduction := &models.Reproduction{
		AnimalID:     animal.ID,
		CurrentPhase: models.PhaseVazias,
	}

	err := reproductionRepo.Create(reproduction)
	assert.NoError(t, err)
	assert.NotZero(t, reproduction.ID)
}

func TestReproductionRepository_FindByID(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm := createTestFarmForReproduction(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Tata Salt",
		EarTagNumberLocal: 123,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	reproduction := &models.Reproduction{
		AnimalID:     animal.ID,
		CurrentPhase: models.PhaseVazias,
	}
	require.NoError(t, db.Create(reproduction).Error)

	found, err := reproductionRepo.FindByID(reproduction.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, reproduction.ID, found.ID)
	assert.Equal(t, animal.ID, found.AnimalID)
	assert.NotNil(t, found.Animal)
	assert.Equal(t, "Tata Salt", found.Animal.AnimalName)
}

func TestReproductionRepository_FindByID_NotFound(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	found, err := reproductionRepo.FindByID(999)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestReproductionRepository_FindByAnimalID(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm := createTestFarmForReproduction(t, db)

	animal1 := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Tata Salt",
		EarTagNumberLocal: 123,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal1).Error)

	animal2 := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Lays",
		EarTagNumberLocal: 124,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal2).Error)

	reproduction1 := &models.Reproduction{
		AnimalID:     animal1.ID,
		CurrentPhase: models.PhaseVazias,
	}
	require.NoError(t, db.Create(reproduction1).Error)

	reproduction2 := &models.Reproduction{
		AnimalID:     animal2.ID,
		CurrentPhase: models.PhasePrenhas,
	}
	require.NoError(t, db.Create(reproduction2).Error)

	found, err := reproductionRepo.FindByAnimalID(animal1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, animal1.ID, found.AnimalID)
	assert.Equal(t, models.PhaseVazias, found.CurrentPhase)
	assert.NotNil(t, found.Animal)
	assert.Equal(t, "Tata Salt", found.Animal.AnimalName)
}

func TestReproductionRepository_FindByAnimalID_NotFound(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	found, err := reproductionRepo.FindByAnimalID(999)
	assert.NoError(t, err)
	assert.Nil(t, found)
}

func TestReproductionRepository_FindByFarmID(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm1 := createTestFarmForReproduction(t, db)

	company2 := &models.Company{
		CompanyName: "Test Company 2",
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
		AnimalName:        "Tata Salt",
		EarTagNumberLocal: 123,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal1).Error)

	animal2 := &models.Animal{
		FarmID:            farm1.ID,
		AnimalName:        "Lays",
		EarTagNumberLocal: 124,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal2).Error)

	animal3 := &models.Animal{
		FarmID:            farm2.ID,
		AnimalName:        "Matilda",
		EarTagNumberLocal: 125,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal3).Error)

	reproduction1 := &models.Reproduction{
		AnimalID:     animal1.ID,
		CurrentPhase: models.PhaseVazias,
	}
	require.NoError(t, db.Create(reproduction1).Error)

	reproduction2 := &models.Reproduction{
		AnimalID:     animal2.ID,
		CurrentPhase: models.PhasePrenhas,
	}
	require.NoError(t, db.Create(reproduction2).Error)

	reproduction3 := &models.Reproduction{
		AnimalID:     animal3.ID,
		CurrentPhase: models.PhaseLactacao,
	}
	require.NoError(t, db.Create(reproduction3).Error)

	reproductions, err := reproductionRepo.FindByFarmID(farm1.ID)
	assert.NoError(t, err)
	assert.Len(t, reproductions, 2)
	for _, reproduction := range reproductions {
		assert.Equal(t, farm1.ID, reproduction.Animal.FarmID)
		assert.NotNil(t, reproduction.Animal)
	}
}

func TestReproductionRepository_FindByFarmIDWithPagination(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm := createTestFarmForReproduction(t, db)

	for i := 0; i < 10; i++ {
		animal := &models.Animal{
			FarmID:            farm.ID,
			AnimalName:        fmt.Sprintf("Animal %d", i),
			EarTagNumberLocal: 100 + i,
			Sex:               0,
			Breed:             "Holandesa",
			Type:              "Bovino",
			AnimalType:        0,
			Status:            0,
			Purpose:           0,
		}
		require.NoError(t, db.Create(animal).Error)

		reproduction := &models.Reproduction{
			AnimalID:     animal.ID,
			CurrentPhase: models.PhaseVazias,
		}
		require.NoError(t, db.Create(reproduction).Error)
	}

	reproductions, total, err := reproductionRepo.FindByFarmIDWithPagination(farm.ID, 1, 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
	assert.Len(t, reproductions, 5)

	reproductions2, total2, err := reproductionRepo.FindByFarmIDWithPagination(farm.ID, 2, 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total2)
	assert.Len(t, reproductions2, 5)

	assert.NotEqual(t, reproductions[0].ID, reproductions2[0].ID)
}

func TestReproductionRepository_FindByFarmIDWithPagination_EmptyResult(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm := createTestFarmForReproduction(t, db)

	reproductions, total, err := reproductionRepo.FindByFarmIDWithPagination(farm.ID, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, reproductions, 0)
}

func TestReproductionRepository_Update(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm := createTestFarmForReproduction(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Tata Salt",
		EarTagNumberLocal: 123,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	reproduction := &models.Reproduction{
		AnimalID:     animal.ID,
		CurrentPhase: models.PhaseVazias,
	}
	require.NoError(t, db.Create(reproduction).Error)

	reproduction.CurrentPhase = models.PhasePrenhas
	now := time.Now()
	reproduction.PregnancyDate = &now
	reproduction.Observations = "Atualizado"

	err := reproductionRepo.Update(reproduction)
	assert.NoError(t, err)

	updated, err := reproductionRepo.FindByID(reproduction.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.PhasePrenhas, updated.CurrentPhase)
	assert.NotNil(t, updated.PregnancyDate)
	assert.Equal(t, "Atualizado", updated.Observations)
}

func TestReproductionRepository_Delete(t *testing.T) {
	db := setupReproductionTestDB(t)
	reproductionRepo := repository.NewReproductionRepository(db)

	farm := createTestFarmForReproduction(t, db)

	animal := &models.Animal{
		FarmID:            farm.ID,
		AnimalName:        "Tata Salt",
		EarTagNumberLocal: 123,
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           0,
	}
	require.NoError(t, db.Create(animal).Error)

	reproduction := &models.Reproduction{
		AnimalID:     animal.ID,
		CurrentPhase: models.PhaseVazias,
	}
	require.NoError(t, db.Create(reproduction).Error)

	err := reproductionRepo.Delete(reproduction.ID)
	assert.NoError(t, err)

	found, err := reproductionRepo.FindByID(reproduction.ID)
	assert.NoError(t, err)
	assert.Nil(t, found)
}
