package services

import (
	"errors"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/cache"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestAnimalService_CreateAnimal(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	service := service.NewAnimalService(mockRepo, mockCache)

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

	mockRepo.On("FindByEarTagNumber", uint(1), 123).Return((*models.Animal)(nil), nil)
	mockRepo.On("Create", animal).Return(nil)
	mockCache.On("Delete", mock.AnythingOfType("string")).Return(nil)

	err := service.CreateAnimal(animal)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_CreateAnimal_InvalidSex(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	animal := &models.Animal{
		FarmID:            1,
		EarTagNumberLocal: 123,
		AnimalName:        "Vaca Teste",
		Sex:               2,
		Breed:             "Holandesa",
		Type:              "Bovino",
		Confinement:       false,
		AnimalType:        0,
		Status:            0,
		Fertilization:     false,
		Castrated:         false,
		Purpose:           1,
		CurrentBatch:      1,
	}

	err := service.CreateAnimal(animal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sexo deve ser 0 (Fêmea) ou 1 (Macho)")
}

func TestAnimalService_CreateAnimal_InvalidAnimalType(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	animal := &models.Animal{
		FarmID:            1,
		EarTagNumberLocal: 123,
		AnimalName:        "Vaca Teste",
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		Confinement:       false,
		AnimalType:        15,
		Status:            0,
		Fertilization:     false,
		Castrated:         false,
		Purpose:           1,
		CurrentBatch:      1,
	}

	err := service.CreateAnimal(animal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tipo de animal inválido")
}

func TestAnimalService_CreateAnimal_InvalidPurpose(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	animal := &models.Animal{
		FarmID:            1,
		EarTagNumberLocal: 123,
		AnimalName:        "Vaca Teste",
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		Confinement:       false,
		AnimalType:        0,
		Status:            0,
		Fertilization:     false,
		Castrated:         false,
		Purpose:           5,
		CurrentBatch:      1,
	}

	err := service.CreateAnimal(animal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "propósito deve ser 0 (Carne), 1 (Leite) ou 2 (Reprodução)")
}

func TestAnimalService_GetAnimalByID(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	expectedAnimal := &models.Animal{
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

	mockRepo.On("FindByID", uint(1)).Return(expectedAnimal, nil)

	animal, err := service.GetAnimalByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedAnimal, animal)
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_GetAnimalByID_NotFound(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	mockRepo.On("FindByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	animal, err := service.GetAnimalByID(999)

	assert.Error(t, err)
	assert.Nil(t, animal)
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_GetAnimalsByFarmID(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	service := service.NewAnimalService(mockRepo, mockCache)

	expectedAnimals := []models.Animal{
		{
			ID:                1,
			FarmID:            1,
			EarTagNumberLocal: 123,
			AnimalName:        "Vaca 1",
			Sex:               0,
			Breed:             "Holandesa",
			Type:              "Bovino",
		},
		{
			ID:                2,
			FarmID:            1,
			EarTagNumberLocal: 124,
			AnimalName:        "Vaca 2",
			Sex:               0,
			Breed:             "Holandesa",
			Type:              "Bovino",
		},
	}

	mockCache.On("Get", "animals:farm:1", mock.Anything).Return(cache.ErrCacheMiss)
	mockRepo.On("FindByFarmID", uint(1)).Return(expectedAnimals, nil)
	mockCache.On("Set", "animals:farm:1", expectedAnimals, int32(300)).Return(nil)

	animals, err := service.GetAnimalsByFarmID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedAnimals, animals)
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_GetAnimalsByFarmIDAndSex(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	expectedAnimals := []models.Animal{
		{
			ID:                1,
			FarmID:            1,
			EarTagNumberLocal: 123,
			AnimalName:        "Touro 1",
			Sex:               1,
			Breed:             "Holandesa",
			Type:              "Bovino",
		},
		{
			ID:                2,
			FarmID:            1,
			EarTagNumberLocal: 124,
			AnimalName:        "Touro 2",
			Sex:               1,
			Breed:             "Holandesa",
			Type:              "Bovino",
		},
	}

	mockRepo.On("FindByFarmIDAndSex", uint(1), 1).Return(expectedAnimals, nil)

	animals, err := service.GetAnimalsByFarmIDAndSex(1, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedAnimals, animals)
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_UpdateAnimal(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	service := service.NewAnimalService(mockRepo, mockCache)

	existingAnimal := &models.Animal{
		ID:                   1,
		FarmID:               1,
		EarTagNumberLocal:    123,
		EarTagNumberRegister: 456,
		AnimalName:           "Vaca Original",
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
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	updatedAnimal := &models.Animal{
		ID:                   1,
		FarmID:               1,
		EarTagNumberLocal:    123,
		EarTagNumberRegister: 456,
		AnimalName:           "Vaca Atualizada",
		Sex:                  0,
		Breed:                "Holandesa",
		Type:                 "Bovino",
		Confinement:          true,
		AnimalType:           0,
		Status:               0,
		Fertilization:        false,
		Castrated:            true,
		Purpose:              1,
		CurrentBatch:         1,
		CreatedAt:            existingAnimal.CreatedAt,
		UpdatedAt:            time.Now(),
	}

	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(nil)
	mockCache.On("Delete", "animals:farm:1").Return(nil)

	err := service.UpdateAnimal(updatedAnimal)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_UpdateAnimal_NotFound(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	animal := &models.Animal{
		ID:     999,
		FarmID: 1,
	}

	mockRepo.On("FindByID", uint(999)).Return(nil, nil)

	err := service.UpdateAnimal(animal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "animal não encontrado")
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_UpdateAnimal_InvalidSex(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	animal := &models.Animal{
		ID:     1,
		FarmID: 1,
		Sex:    2,
	}

	existingAnimal := &models.Animal{
		ID:     1,
		FarmID: 1,
		Sex:    0,
	}
	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil)

	err := service.UpdateAnimal(animal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sexo deve ser 0 (Fêmea) ou 1 (Macho)")
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_DeleteAnimal(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	mockCache := new(MockCache)
	service := service.NewAnimalService(mockRepo, mockCache)

	existingAnimal := &models.Animal{
		ID:     1,
		FarmID: 1,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)
	mockCache.On("Delete", "animals:farm:1").Return(nil)

	err := service.DeleteAnimal(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_DeleteAnimal_NotFound(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	mockRepo.On("FindByID", uint(999)).Return(nil, nil)

	err := service.DeleteAnimal(999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "animal não encontrado")
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_DeleteAnimal_InvalidID(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	err := service.DeleteAnimal(0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID do animal é obrigatório")
}

func TestAnimalService_UpdateAnimal_InvalidID(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	animal := &models.Animal{
		ID: 0,
	}

	err := service.UpdateAnimal(animal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID do animal é obrigatório")
}

func TestAnimalService_CreateAnimal_MissingFields(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	testCases := []struct {
		name     string
		animal   *models.Animal
		errorMsg string
	}{
		{
			name:     "MissingFarmID",
			animal:   &models.Animal{EarTagNumberLocal: 123, AnimalName: "Vaca", Breed: "Holandesa", Type: "Bovino"},
			errorMsg: "farm ID é obrigatório",
		},
		{
			name:     "MissingEarTagNumber",
			animal:   &models.Animal{FarmID: 1, AnimalName: "Vaca", Breed: "Holandesa", Type: "Bovino"},
			errorMsg: "número da brinca local é obrigatório",
		},
		{
			name:     "MissingAnimalName",
			animal:   &models.Animal{FarmID: 1, EarTagNumberLocal: 123, Breed: "Holandesa", Type: "Bovino"},
			errorMsg: "nome do animal é obrigatório",
		},
		{
			name:     "MissingBreed",
			animal:   &models.Animal{FarmID: 1, EarTagNumberLocal: 123, AnimalName: "Vaca", Type: "Bovino"},
			errorMsg: "raça do animal é obrigatória",
		},
		{
			name:     "MissingType",
			animal:   &models.Animal{FarmID: 1, EarTagNumberLocal: 123, AnimalName: "Vaca", Breed: "Holandesa"},
			errorMsg: "tipo do animal é obrigatório",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.CreateAnimal(tc.animal)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMsg)
		})
	}
}

func TestAnimalService_CreateAnimal_RepositoryError(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	animal := &models.Animal{
		FarmID:            1,
		EarTagNumberLocal: 123,
		AnimalName:        "Vaca Teste",
		Sex:               0,
		Breed:             "Holandesa",
		Type:              "Bovino",
		AnimalType:        0,
		Purpose:           1,
	}

	mockRepo.On("FindByEarTagNumber", uint(1), 123).Return(nil, errors.New("database error"))

	err := service.CreateAnimal(animal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_UpdateAnimal_RepositoryError(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	animal := &models.Animal{
		ID:     1,
		FarmID: 1,
	}

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	err := service.UpdateAnimal(animal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}

func TestAnimalService_DeleteAnimal_RepositoryError(t *testing.T) {
	mockRepo := new(MockAnimalRepository)
	service := service.NewAnimalService(mockRepo, new(MockCache))

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	err := service.DeleteAnimal(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}
