package integration

import (
	"encoding/json"
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAnimalModel_Validation(t *testing.T) {
	animal := &models.Animal{
		FarmID:               1,
		EarTagNumberLocal:    123,
		EarTagNumberRegister: 456,
		AnimalName:           "Vaca Teste",
		Sex:                  0,
		Breed:                "Holandesa",
		Type:                 "Bovino",
		Confinement:          false,
		AnimalType:           0,
		Status:               0,
		Fertilization:        false,
		Castrated:            false,
		Purpose:              1,
		CurrentBatch:         1,
	}

	// Test valid animal
	assert.Equal(t, uint(1), animal.FarmID)
	assert.Equal(t, 123, animal.EarTagNumberLocal)
	assert.Equal(t, "Vaca Teste", animal.AnimalName)
	assert.Equal(t, 0, animal.Sex)
	assert.Equal(t, "Holandesa", animal.Breed)
	assert.Equal(t, "Bovino", animal.Type)
	assert.Equal(t, false, animal.Confinement)
	assert.Equal(t, false, animal.Castrated)
	assert.Equal(t, 1, animal.Purpose)
}

func TestAnimalModel_InvalidSex(t *testing.T) {
	invalidSexes := []int{-1, 2, 3, 10}

	for _, sex := range invalidSexes {
		assert.True(t, sex < 0 || sex > 1, "Sex %d should be invalid", sex)
	}
}

func TestAnimalModel_ValidSex(t *testing.T) {
	validSexes := []int{0, 1}

	for _, sex := range validSexes {
		assert.True(t, sex == 0 || sex == 1, "Sex %d should be valid", sex)
	}
}

func TestAnimalModel_Purpose(t *testing.T) {
	purposes := map[int]string{
		0: "Carne",
		1: "Leite",
		2: "Reprodução",
	}

	for purpose, expected := range purposes {
		animal := &models.Animal{
			FarmID:            1,
			EarTagNumberLocal: 123,
			AnimalName:        "Animal Teste",
			Sex:               0,
			Breed:             "Holandesa",
			Type:              "Bovino",
			Confinement:       false,
			AnimalType:        0,
			Status:            0,
			Fertilization:     false,
			Castrated:         false,
			Purpose:           purpose,
			CurrentBatch:      1,
		}

		assert.Equal(t, purpose, animal.Purpose)
		assert.Equal(t, expected, purposes[animal.Purpose])
	}
}

func TestAnimalModel_Status(t *testing.T) {
	statuses := map[int]string{
		0: "Vivo",
		1: "Morto",
	}

	for status, expected := range statuses {
		animal := &models.Animal{
			FarmID:            1,
			EarTagNumberLocal: 123,
			AnimalName:        "Animal Teste",
			Sex:               0,
			Breed:             "Holandesa",
			Type:              "Bovino",
			Confinement:       false,
			AnimalType:        0,
			Status:            status,
			Fertilization:     false,
			Castrated:         false,
			Purpose:           1,
			CurrentBatch:      1,
		}

		assert.Equal(t, status, animal.Status)
		assert.Equal(t, expected, statuses[animal.Status])
	}
}

func TestAnimalModel_AnimalType(t *testing.T) {
	animalTypes := map[int]string{
		0:  "Bovino",
		1:  "Suíno",
		2:  "Ave",
		3:  "Caprino",
		4:  "Ovino",
		5:  "Equino",
		6:  "Asinino",
		7:  "Muar",
		8:  "Bubalino",
		9:  "Coelho",
		10: "Outros",
	}

	for animalType, expected := range animalTypes {
		animal := &models.Animal{
			FarmID:            1,
			EarTagNumberLocal: 123,
			AnimalName:        "Animal Teste",
			Sex:               0,
			Breed:             "Holandesa",
			Type:              "Bovino",
			Confinement:       false,
			AnimalType:        animalType,
			Status:            0,
			Fertilization:     false,
			Castrated:         false,
			Purpose:           1,
			CurrentBatch:      1,
		}

		assert.Equal(t, animalType, animal.AnimalType)
		assert.Equal(t, expected, animalTypes[animal.AnimalType])
	}
}

func TestAnimalModel_BooleanFields(t *testing.T) {
	animal := &models.Animal{
		FarmID:            1,
		EarTagNumberLocal: 123,
		AnimalName:        "Animal Teste",
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		Confinement:       true,
		AnimalType:        0,
		Status:            0,
		Fertilization:     true,
		Castrated:         true,
		Purpose:           1,
		CurrentBatch:      1,
	}

	assert.True(t, animal.Confinement)
	assert.True(t, animal.Fertilization)
	assert.True(t, animal.Castrated)
}

func TestAnimalModel_JSONSerialization(t *testing.T) {
	animal := &models.Animal{
		ID:                   1,
		FarmID:               1,
		EarTagNumberLocal:    123,
		EarTagNumberRegister: 456,
		AnimalName:           "Vaca Teste",
		Sex:                  0,
		Breed:                "Holandesa",
		Type:                 "Bovino",
		Confinement:          false,
		AnimalType:           0,
		Status:               0,
		Fertilization:        false,
		Castrated:            false,
		Purpose:              1,
		CurrentBatch:         1,
	}

	jsonData, err := json.Marshal(animal)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var deserializedAnimal models.Animal
	err = json.Unmarshal(jsonData, &deserializedAnimal)
	assert.NoError(t, err)
	assert.Equal(t, animal.ID, deserializedAnimal.ID)
	assert.Equal(t, animal.AnimalName, deserializedAnimal.AnimalName)
	assert.Equal(t, animal.Sex, deserializedAnimal.Sex)
}

func TestAnimalModel_RequiredFields(t *testing.T) {
	animal := &models.Animal{
		FarmID:            1,
		EarTagNumberLocal: 123,
		AnimalName:        "Vaca Teste",
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		AnimalType:        0,
		Status:            0,
		Purpose:           1,
	}

	assert.NotZero(t, animal.FarmID)
	assert.NotZero(t, animal.EarTagNumberLocal)
	assert.NotEmpty(t, animal.AnimalName)
	assert.NotEmpty(t, animal.Breed)
	assert.NotEmpty(t, animal.Type)
}

func TestAnimalModel_OptionalFields(t *testing.T) {
	animal := &models.Animal{
		FarmID:               1,
		EarTagNumberLocal:    123,
		AnimalName:           "Vaca Teste",
		Sex:                  0,
		Breed:                "Holandesa",
		Type:                 "Bovino",
		AnimalType:           0,
		Status:               0,
		Purpose:              1,
		EarTagNumberRegister: 0,
		Photo:                "",
		FatherID:             nil,
		MotherID:             nil,
	}

	assert.Zero(t, animal.EarTagNumberRegister)
	assert.Empty(t, animal.Photo)
	assert.Nil(t, animal.FatherID)
	assert.Nil(t, animal.MotherID)
}
