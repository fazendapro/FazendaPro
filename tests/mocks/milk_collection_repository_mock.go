package mocks

import (
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockMilkCollectionRepository struct {
	mock.Mock
}

func (m *MockMilkCollectionRepository) Create(milkCollection *models.MilkCollection) error {
	args := m.Called(milkCollection)
	return args.Error(0)
}

func (m *MockMilkCollectionRepository) FindByID(id uint) (*models.MilkCollection, error) {
	args := m.Called(id)
	return args.Get(0).(*models.MilkCollection), args.Error(1)
}

func (m *MockMilkCollectionRepository) FindByFarmID(farmID uint) ([]models.MilkCollection, error) {
	args := m.Called(farmID)
	return args.Get(0).([]models.MilkCollection), args.Error(1)
}

func (m *MockMilkCollectionRepository) FindByFarmIDWithDateRange(farmID uint, startDate, endDate *time.Time) ([]models.MilkCollection, error) {
	args := m.Called(farmID, startDate, endDate)
	return args.Get(0).([]models.MilkCollection), args.Error(1)
}

func (m *MockMilkCollectionRepository) FindByAnimalID(animalID uint) ([]models.MilkCollection, error) {
	args := m.Called(animalID)
	return args.Get(0).([]models.MilkCollection), args.Error(1)
}

func (m *MockMilkCollectionRepository) Update(milkCollection *models.MilkCollection) error {
	args := m.Called(milkCollection)
	return args.Error(0)
}

func (m *MockMilkCollectionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
