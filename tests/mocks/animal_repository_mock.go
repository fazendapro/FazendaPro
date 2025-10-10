package mocks

import (
	"github.com/fazendapro/FazendaPro-api/internal/models"
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
	return args.Get(0).(*models.Animal), args.Error(1)
}

func (m *MockAnimalRepository) FindByFarmID(farmID uint) ([]models.Animal, error) {
	args := m.Called(farmID)
	return args.Get(0).([]models.Animal), args.Error(1)
}

func (m *MockAnimalRepository) FindByEarTagNumber(farmID uint, earTagNumber int) (*models.Animal, error) {
	args := m.Called(farmID, earTagNumber)
	return args.Get(0).(*models.Animal), args.Error(1)
}

func (m *MockAnimalRepository) FindByFarmIDAndSex(farmID uint, sex int) ([]models.Animal, error) {
	args := m.Called(farmID, sex)
	return args.Get(0).([]models.Animal), args.Error(1)
}

func (m *MockAnimalRepository) Update(animal *models.Animal) error {
	args := m.Called(animal)
	return args.Error(0)
}

func (m *MockAnimalRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
