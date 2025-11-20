package repository

import (
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
)

type AnimalRepositoryInterface interface {
	Create(animal *models.Animal) error
	FindByID(id uint) (*models.Animal, error)
	FindByFarmID(farmID uint) ([]models.Animal, error)
	FindByEarTagNumber(farmID uint, earTagNumber int) (*models.Animal, error)
	FindByFarmIDAndSex(farmID uint, sex int) ([]models.Animal, error)
	CountBySex(farmID uint, sex int) (int64, error)
	Update(animal *models.Animal) error
	Delete(id uint) error
}

type UserRepositoryInterface interface {
	FindByPersonEmail(email string) (*models.User, error)
	CreateWithPerson(user *models.User, personData *models.Person) error
	FindByIDWithPerson(userID uint) (*models.User, error)
	UpdatePersonData(userID uint, personData *models.Person) error
	ValidatePassword(userID uint, password string) (bool, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	FarmExists(farmID uint) (bool, error)
	CreateDefaultFarm(farmID uint) error
	GetUserFarms(userID uint) ([]models.Farm, error)
	GetUserFarmCount(userID uint) (int64, error)
	GetUserFarmByID(userID, farmID uint) (*models.Farm, error)
	CreateUserFarm(userFarm *models.UserFarm) error
}

type MilkCollectionRepositoryInterface interface {
	Create(milkCollection *models.MilkCollection) error
	FindByID(id uint) (*models.MilkCollection, error)
	FindByFarmID(farmID uint) ([]models.MilkCollection, error)
	FindByFarmIDWithDateRange(farmID uint, startDate, endDate *time.Time) ([]models.MilkCollection, error)
	FindByAnimalID(animalID uint) ([]models.MilkCollection, error)
	Update(milkCollection *models.MilkCollection) error
	Delete(id uint) error
}

type ReproductionRepositoryInterface interface {
	Create(reproduction *models.Reproduction) error
	FindByID(id uint) (*models.Reproduction, error)
	FindByAnimalID(animalID uint) (*models.Reproduction, error)
	FindByFarmID(farmID uint) ([]models.Reproduction, error)
	FindByPhase(phase models.ReproductionPhase) ([]models.Reproduction, error)
	Update(reproduction *models.Reproduction) error
	Delete(id uint) error
}

type FarmRepositoryInterface interface {
	FindByID(id uint) (*models.Farm, error)
	Update(farm *models.Farm) error
	LoadCompanyData(farm *models.Farm) error
}

type DebtRepositoryInterface interface {
	Create(debt *models.Debt) error
	FindByID(id uint) (*models.Debt, error)
	FindAllWithPagination(page, limit int, year, month *int) ([]models.Debt, int64, error)
	Delete(id uint) error
	GetTotalByPersonInMonth(year, month int) ([]PersonTotal, error)
}
