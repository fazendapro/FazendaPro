package repositories

import (
	"errors"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReproductionRepository_FindByPhase(t *testing.T) {
	// Setup
	mockDB := &MockGormDB{}
	reproductionRepo := repository.NewReproductionRepository(mockDB)

	// Mock data
	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200) // 200 dias atrás

	mockReproductions := []models.Reproduction{
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                1,
				FarmID:            1,
				AnimalName:        "Tata Salt",
				EarTagNumberLocal: 123,
				Photo:             "src/assets/images/mocked/cows/tata.png",
			},
		},
		{
			ID:            2,
			AnimalID:      2,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                2,
				FarmID:            1,
				AnimalName:        "Lays",
				EarTagNumberLocal: 124,
				Photo:             "src/assets/images/mocked/cows/lays.png",
			},
		},
	}

	// Mock GORM calls
	mockDB.On("Preload", "Animal").Return(mockDB)
	mockDB.On("Where", "current_phase = ?", models.PhasePrenhas).Return(mockDB)
	mockDB.On("Find", &[]models.Reproduction{}).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Mock the Find method to populate the slice
	mockDB.On("Find", mock.AnythingOfType("*[]models.Reproduction")).Run(func(args mock.Arguments) {
		reproductions := args.Get(0).(*[]models.Reproduction)
		*reproductions = mockReproductions
	}).Return(mockDB)

	// Test
	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, models.PhasePrenhas, result[0].CurrentPhase)
	assert.Equal(t, models.PhasePrenhas, result[1].CurrentPhase)
	assert.Equal(t, "Tata Salt", result[0].Animal.AnimalName)
	assert.Equal(t, "Lays", result[1].Animal.AnimalName)

	mockDB.AssertExpectations(t)
}

func TestReproductionRepository_FindByPhase_Error(t *testing.T) {
	// Setup
	mockDB := &MockGormDB{}
	reproductionRepo := repository.NewReproductionRepository(mockDB)

	// Mock GORM calls with error
	mockDB.On("Preload", "Animal").Return(mockDB)
	mockDB.On("Where", "current_phase = ?", models.PhasePrenhas).Return(mockDB)
	mockDB.On("Find", mock.AnythingOfType("*[]models.Reproduction")).Return(mockDB)
	mockDB.On("Error").Return(errors.New("database error"))

	// Test
	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	mockDB.AssertExpectations(t)
}

func TestReproductionRepository_FindByPhase_EmptyResults(t *testing.T) {
	// Setup
	mockDB := &MockGormDB{}
	reproductionRepo := repository.NewReproductionRepository(mockDB)

	// Mock GORM calls with empty results
	mockDB.On("Preload", "Animal").Return(mockDB)
	mockDB.On("Where", "current_phase = ?", models.PhasePrenhas).Return(mockDB)
	mockDB.On("Find", mock.AnythingOfType("*[]models.Reproduction")).Run(func(args mock.Arguments) {
		reproductions := args.Get(0).(*[]models.Reproduction)
		*reproductions = []models.Reproduction{}
	}).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Test
	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 0)

	mockDB.AssertExpectations(t)
}

func TestReproductionRepository_FindByPhase_DifferentPhases(t *testing.T) {
	// Setup
	mockDB := &MockGormDB{}
	reproductionRepo := repository.NewReproductionRepository(mockDB)

	// Mock data for different phases
	mockReproductions := []models.Reproduction{
		{
			ID:           1,
			AnimalID:     1,
			CurrentPhase: models.PhaseLactacao,
			Animal: models.Animal{
				ID:                1,
				FarmID:            1,
				AnimalName:        "Tata Salt",
				EarTagNumberLocal: 123,
			},
		},
	}

	// Mock GORM calls
	mockDB.On("Preload", "Animal").Return(mockDB)
	mockDB.On("Where", "current_phase = ?", models.PhaseLactacao).Return(mockDB)
	mockDB.On("Find", mock.AnythingOfType("*[]models.Reproduction")).Run(func(args mock.Arguments) {
		reproductions := args.Get(0).(*[]models.Reproduction)
		*reproductions = mockReproductions
	}).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Test
	result, err := reproductionRepo.FindByPhase(models.PhaseLactacao)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, models.PhaseLactacao, result[0].CurrentPhase)

	mockDB.AssertExpectations(t)
}

func TestReproductionRepository_FindByPhase_WithAnimalData(t *testing.T) {
	// Setup
	mockDB := &MockGormDB{}
	reproductionRepo := repository.NewReproductionRepository(mockDB)

	// Mock data with complete animal information
	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200) // 200 dias atrás

	mockReproductions := []models.Reproduction{
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                   1,
				FarmID:               1,
				AnimalName:           "Tata Salt",
				EarTagNumberLocal:    123,
				EarTagNumberRegister: 456,
				Photo:                "src/assets/images/mocked/cows/tata.png",
				Sex:                  1,
				Breed:                "Holandesa",
				Type:                 "vaca",
			},
		},
	}

	// Mock GORM calls
	mockDB.On("Preload", "Animal").Return(mockDB)
	mockDB.On("Where", "current_phase = ?", models.PhasePrenhas).Return(mockDB)
	mockDB.On("Find", mock.AnythingOfType("*[]models.Reproduction")).Run(func(args mock.Arguments) {
		reproductions := args.Get(0).(*[]models.Reproduction)
		*reproductions = mockReproductions
	}).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Test
	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	reproduction := result[0]
	assert.Equal(t, uint(1), reproduction.AnimalID)
	assert.Equal(t, models.PhasePrenhas, reproduction.CurrentPhase)
	assert.NotNil(t, reproduction.PregnancyDate)
	assert.Equal(t, pregnancyDate, *reproduction.PregnancyDate)

	// Verify animal data
	animal := reproduction.Animal
	assert.Equal(t, uint(1), animal.ID)
	assert.Equal(t, uint(1), animal.FarmID)
	assert.Equal(t, "Tata Salt", animal.AnimalName)
	assert.Equal(t, 123, animal.EarTagNumberLocal)
	assert.Equal(t, 456, animal.EarTagNumberRegister)
	assert.Equal(t, "src/assets/images/mocked/cows/tata.png", animal.Photo)
	assert.Equal(t, 1, animal.Sex)
	assert.Equal(t, "Holandesa", animal.Breed)
	assert.Equal(t, "vaca", animal.Type)

	mockDB.AssertExpectations(t)
}

func TestReproductionRepository_FindByPhase_WithNullPregnancyDate(t *testing.T) {
	// Setup
	mockDB := &MockGormDB{}
	reproductionRepo := repository.NewReproductionRepository(mockDB)

	// Mock data with null pregnancy date
	mockReproductions := []models.Reproduction{
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: nil, // Null pregnancy date
			Animal: models.Animal{
				ID:                1,
				FarmID:            1,
				AnimalName:        "Tata Salt",
				EarTagNumberLocal: 123,
			},
		},
	}

	// Mock GORM calls
	mockDB.On("Preload", "Animal").Return(mockDB)
	mockDB.On("Where", "current_phase = ?", models.PhasePrenhas).Return(mockDB)
	mockDB.On("Find", mock.AnythingOfType("*[]models.Reproduction")).Run(func(args mock.Arguments) {
		reproductions := args.Get(0).(*[]models.Reproduction)
		*reproductions = mockReproductions
	}).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Test
	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	reproduction := result[0]
	assert.Equal(t, models.PhasePrenhas, reproduction.CurrentPhase)
	assert.Nil(t, reproduction.PregnancyDate)

	mockDB.AssertExpectations(t)
}

func TestReproductionRepository_FindByPhase_WithMultipleFarms(t *testing.T) {
	// Setup
	mockDB := &MockGormDB{}
	reproductionRepo := repository.NewReproductionRepository(mockDB)

	// Mock data with animals from different farms
	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200) // 200 dias atrás

	mockReproductions := []models.Reproduction{
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                1,
				FarmID:            1,
				AnimalName:        "Tata Salt",
				EarTagNumberLocal: 123,
			},
		},
		{
			ID:            2,
			AnimalID:      2,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                2,
				FarmID:            2,
				AnimalName:        "Lays",
				EarTagNumberLocal: 124,
			},
		},
		{
			ID:            3,
			AnimalID:      3,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                3,
				FarmID:            1,
				AnimalName:        "Matilda",
				EarTagNumberLocal: 125,
			},
		},
	}

	// Mock GORM calls
	mockDB.On("Preload", "Animal").Return(mockDB)
	mockDB.On("Where", "current_phase = ?", models.PhasePrenhas).Return(mockDB)
	mockDB.On("Find", mock.AnythingOfType("*[]models.Reproduction")).Run(func(args mock.Arguments) {
		reproductions := args.Get(0).(*[]models.Reproduction)
		*reproductions = mockReproductions
	}).Return(mockDB)
	mockDB.On("Error").Return(nil)

	// Test
	result, err := reproductionRepo.FindByPhase(models.PhasePrenhas)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// Verify all reproductions are in PhasePrenhas
	for _, reproduction := range result {
		assert.Equal(t, models.PhasePrenhas, reproduction.CurrentPhase)
		assert.NotNil(t, reproduction.PregnancyDate)
	}

	// Verify farm distribution
	farm1Count := 0
	farm2Count := 0
	for _, reproduction := range result {
		if reproduction.Animal.FarmID == 1 {
			farm1Count++
		} else if reproduction.Animal.FarmID == 2 {
			farm2Count++
		}
	}
	assert.Equal(t, 2, farm1Count)
	assert.Equal(t, 1, farm2Count)

	mockDB.AssertExpectations(t)
}

// MockGormDB is a mock implementation of gorm.DB for testing
type MockGormDB struct {
	mock.Mock
}

func (m *MockGormDB) Preload(query string, args ...interface{}) *MockGormDB {
	m.Called(query, args)
	return m
}

func (m *MockGormDB) Where(query interface{}, args ...interface{}) *MockGormDB {
	m.Called(query, args)
	return m
}

func (m *MockGormDB) Find(dest interface{}, conds ...interface{}) *MockGormDB {
	m.Called(dest, conds)
	return m
}

func (m *MockGormDB) Error() error {
	args := m.Called()
	return args.Error(0)
}
