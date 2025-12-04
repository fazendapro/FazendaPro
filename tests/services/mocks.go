package services

import (
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
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

type MockDebtRepository struct {
	mock.Mock
}

func (m *MockDebtRepository) Create(debt *models.Debt) error {
	args := m.Called(debt)
	return args.Error(0)
}

func (m *MockDebtRepository) FindByID(id uint) (*models.Debt, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Debt), args.Error(1)
}

func (m *MockDebtRepository) FindAllWithPagination(page, limit int, year, month *int) ([]models.Debt, int64, error) {
	args := m.Called(page, limit, year, month)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Debt), args.Get(1).(int64), args.Error(2)
}

func (m *MockDebtRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDebtRepository) GetTotalByPersonInMonth(year, month int) ([]repository.PersonTotal, error) {
	args := m.Called(year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.PersonTotal), args.Error(1)
}

type MockFarmRepository struct {
	mock.Mock
}

func (m *MockFarmRepository) FindByID(id uint) (*models.Farm, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Farm), args.Error(1)
}

func (m *MockFarmRepository) Update(farm *models.Farm) error {
	args := m.Called(farm)
	return args.Error(0)
}

func (m *MockFarmRepository) LoadCompanyData(farm *models.Farm) error {
	args := m.Called(farm)
	return args.Error(0)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string, dest interface{}) error {
	args := m.Called(key, dest)
	return args.Error(0)
}

func (m *MockCache) Set(key string, value interface{}, expiration int32) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockCache) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockCache) Increment(key string, delta uint64) (uint64, error) {
	args := m.Called(key, delta)
	return args.Get(0).(uint64), args.Error(1)
}

type MockVaccineRepository struct {
	mock.Mock
}

func (m *MockVaccineRepository) Create(vaccine *models.Vaccine) error {
	args := m.Called(vaccine)
	return args.Error(0)
}

func (m *MockVaccineRepository) FindByID(id uint) (*models.Vaccine, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Vaccine), args.Error(1)
}

func (m *MockVaccineRepository) FindByFarmID(farmID uint) ([]models.Vaccine, error) {
	args := m.Called(farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Vaccine), args.Error(1)
}

func (m *MockVaccineRepository) Update(vaccine *models.Vaccine) error {
	args := m.Called(vaccine)
	return args.Error(0)
}

func (m *MockVaccineRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockVaccineApplicationRepository struct {
	mock.Mock
}

func (m *MockVaccineApplicationRepository) Create(vaccineApplication *models.VaccineApplication) error {
	args := m.Called(vaccineApplication)
	return args.Error(0)
}

func (m *MockVaccineApplicationRepository) FindByID(id uint) (*models.VaccineApplication, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VaccineApplication), args.Error(1)
}

func (m *MockVaccineApplicationRepository) FindByFarmID(farmID uint) ([]models.VaccineApplication, error) {
	args := m.Called(farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.VaccineApplication), args.Error(1)
}

func (m *MockVaccineApplicationRepository) FindByFarmIDWithDateRange(farmID uint, startDate, endDate *time.Time) ([]models.VaccineApplication, error) {
	args := m.Called(farmID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.VaccineApplication), args.Error(1)
}

func (m *MockVaccineApplicationRepository) FindByAnimalID(animalID uint) ([]models.VaccineApplication, error) {
	args := m.Called(animalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.VaccineApplication), args.Error(1)
}

func (m *MockVaccineApplicationRepository) FindByVaccineID(vaccineID uint) ([]models.VaccineApplication, error) {
	args := m.Called(vaccineID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.VaccineApplication), args.Error(1)
}

func (m *MockVaccineApplicationRepository) Update(vaccineApplication *models.VaccineApplication) error {
	args := m.Called(vaccineApplication)
	return args.Error(0)
}

func (m *MockVaccineApplicationRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
