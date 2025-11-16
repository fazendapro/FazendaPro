package services

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

func (m *MockAnimalRepository) FindByFarmIDAndSex(farmID uint, sex int) ([]models.Animal, error) {
	args := m.Called(farmID, sex)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Animal), args.Error(1)
}

func (m *MockAnimalRepository) CountBySex(farmID uint, sex int) (int64, error) {
	args := m.Called(farmID, sex)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAnimalRepository) Update(animal *models.Animal) error {
	args := m.Called(animal)
	return args.Error(0)
}

func (m *MockAnimalRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByPersonEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CreateWithPerson(user *models.User, personData *models.Person) error {
	args := m.Called(user, personData)
	return args.Error(0)
}

func (m *MockUserRepository) FindByIDWithPerson(userID uint) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePersonData(userID uint, personData *models.Person) error {
	args := m.Called(userID, personData)
	return args.Error(0)
}

func (m *MockUserRepository) ValidatePassword(userID uint, password string) (bool, error) {
	args := m.Called(userID, password)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FarmExists(farmID uint) (bool, error) {
	args := m.Called(farmID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CreateDefaultFarm(farmID uint) error {
	args := m.Called(farmID)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserFarms(userID uint) ([]models.Farm, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Farm), args.Error(1)
}

func (m *MockUserRepository) GetUserFarmCount(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) GetUserFarmByID(userID, farmID uint) (*models.Farm, error) {
	args := m.Called(userID, farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Farm), args.Error(1)
}

func (m *MockUserRepository) CreateUserFarm(userFarm *models.UserFarm) error {
	args := m.Called(userFarm)
	return args.Error(0)
}

type MockReproductionRepository struct {
	mock.Mock
}

func (m *MockReproductionRepository) Create(reproduction *models.Reproduction) error {
	args := m.Called(reproduction)
	return args.Error(0)
}

func (m *MockReproductionRepository) FindByID(id uint) (*models.Reproduction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reproduction), args.Error(1)
}

func (m *MockReproductionRepository) FindByAnimalID(animalID uint) (*models.Reproduction, error) {
	args := m.Called(animalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reproduction), args.Error(1)
}

func (m *MockReproductionRepository) FindByFarmID(farmID uint) ([]models.Reproduction, error) {
	args := m.Called(farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Reproduction), args.Error(1)
}

func (m *MockReproductionRepository) FindByPhase(phase models.ReproductionPhase) ([]models.Reproduction, error) {
	args := m.Called(phase)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Reproduction), args.Error(1)
}

func (m *MockReproductionRepository) Update(reproduction *models.Reproduction) error {
	args := m.Called(reproduction)
	return args.Error(0)
}

func (m *MockReproductionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

