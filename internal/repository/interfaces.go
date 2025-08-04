package repository

import (
	"github.com/fazendapro/FazendaPro-api/internal/models"
)

type AnimalRepositoryInterface interface {
	Create(animal *models.Animal) error
	FindByID(id uint) (*models.Animal, error)
	FindByFarmID(farmID uint) ([]models.Animal, error)
	FindByEarTagNumber(farmID uint, earTagNumber int) (*models.Animal, error)
	Update(animal *models.Animal) error
	Delete(id uint) error
}

// UserRepositoryInterface define os métodos que um repositório de usuários deve implementar
type UserRepositoryInterface interface {
	FindByPersonEmail(email string) (*models.User, error)
	CreateWithPerson(user *models.User, personData *models.Person) error
	FindByIDWithPerson(userID uint) (*models.User, error)
	UpdatePersonData(userID uint, personData *models.Person) error
	ValidatePassword(userID uint, password string) (bool, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}
