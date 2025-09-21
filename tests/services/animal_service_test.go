package services

import (
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAnimalRepository struct {
	mock.Mock
}

func (m *MockAnimalRepository) Create(animal *models.Animal) error {
	args := m.Called(animal)
	return args.Error(0)
}

func (m *MockAnimalRepository) FindByID(id uint) (*models.Animal, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Animal), args.Error(1)
}

func (m *MockAnimalRepository) FindByFarmID(farmID uint) ([]models.Animal, error) {
	args := m.Called(farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Animal), args.Error(1)
}

func (m *MockAnimalRepository) FindByEarTagNumber(farmID uint, earTagNumber int) (*models.Animal, error) {
	args := m.Called(farmID, earTagNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Animal), args.Error(1)
}

func (m *MockAnimalRepository) Update(animal *models.Animal) error {
	args := m.Called(animal)
	return args.Error(0)
}

func (m *MockAnimalRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockAnimalService struct {
	mock.Mock
}

func (m *MockAnimalService) CreateAnimal(animal *models.Animal) error {
	args := m.Called(animal)
	return args.Error(0)
}

func (m *MockAnimalService) GetAnimalByID(id uint) (*models.Animal, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Animal), args.Error(1)
}

func TestAnimalService(t *testing.T) {
	t.Run("CreateAnimal_Success", func(t *testing.T) {
		mockRepo := &MockAnimalRepository{}
		animalService := service.NewAnimalService(mockRepo)

		mockRepo.On("FindByEarTagNumber", uint(1), 123).Return(nil, nil)
		mockRepo.On("Create", mock.AnythingOfType("*models.Animal")).Return(nil)

		animal := &models.Animal{
			FarmID:            1,
			EarTagNumberLocal: 123,
			AnimalName:        "Boi João",
			Sex:               1,
			Breed:             "Nelore",
			Type:              "Bovino",
			AnimalType:        0,
			Purpose:           0,
		}

		err := animalService.CreateAnimal(animal)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetAnimalByID_Success", func(t *testing.T) {
		mockRepo := &MockAnimalRepository{}
		animalService := service.NewAnimalService(mockRepo)

		expectedAnimal := &models.Animal{
			ID:                1,
			FarmID:            1,
			EarTagNumberLocal: 123,
			AnimalName:        "Boi João",
			Sex:               1,
			Breed:             "Nelore",
			Type:              "Bovino",
			AnimalType:        0,
			Purpose:           0,
		}
		mockRepo.On("FindByID", uint(1)).Return(expectedAnimal, nil)

		animal, err := animalService.GetAnimalByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, animal)
		assert.Equal(t, "Boi João", animal.AnimalName)
		assert.Equal(t, "Nelore", animal.Breed)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetAnimalByID_NotFound", func(t *testing.T) {
		mockRepo := &MockAnimalRepository{}
		animalService := service.NewAnimalService(mockRepo)

		mockRepo.On("FindByID", uint(999)).Return(nil, nil)

		animal, err := animalService.GetAnimalByID(999)

		assert.NoError(t, err)
		assert.Nil(t, animal)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetAllAnimals_Success", func(t *testing.T) {
		mockRepo := &MockAnimalRepository{}
		animalService := service.NewAnimalService(mockRepo)

		expectedAnimals := []models.Animal{
			{ID: 1, AnimalName: "Boi João", Breed: "Nelore"},
			{ID: 2, AnimalName: "Vaca Maria", Breed: "Holandesa"},
		}
		mockRepo.On("FindByFarmID", uint(1)).Return(expectedAnimals, nil)

		animals, err := animalService.GetAnimalsByFarmID(1)

		assert.NoError(t, err)
		assert.Len(t, animals, 2)
		assert.Equal(t, "Boi João", animals[0].AnimalName)
		assert.Equal(t, "Vaca Maria", animals[1].AnimalName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateAnimal_Success", func(t *testing.T) {
		mockRepo := &MockAnimalRepository{}
		animalService := service.NewAnimalService(mockRepo)

		mockRepo.On("FindByID", uint(1)).Return(&models.Animal{ID: 1}, nil)
		mockRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(nil)

		animal := &models.Animal{
			ID:                1,
			FarmID:            1,
			EarTagNumberLocal: 123,
			AnimalName:        "Boi João Atualizado",
			Breed:             "Nelore",
		}

		err := animalService.UpdateAnimal(animal)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteAnimal_Success", func(t *testing.T) {
		mockRepo := &MockAnimalRepository{}
		animalService := service.NewAnimalService(mockRepo)

		mockRepo.On("FindByID", uint(1)).Return(&models.Animal{ID: 1}, nil)
		mockRepo.On("Delete", uint(1)).Return(nil)

		err := animalService.DeleteAnimal(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
