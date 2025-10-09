package repository

import (
	"fmt"

	"github.com/fazendapro/FazendaPro-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *Database
}

func NewUserRepository(db *Database) UserRepositoryInterface {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByPersonEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.DB.Preload("Person").Joins("JOIN people ON users.person_id = people.id").Where("people.email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) CreateWithPerson(user *models.User, personData *models.Person) error {
	tx := r.db.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("error starting transaction: %w", tx.Error)
	}

	if err := tx.Create(personData).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error creating person: %w", err)
	}

	user.PersonID = &personData.ID

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error creating user: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByIDWithPerson(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.DB.Preload("Person").Preload("Farm").Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) UpdatePersonData(userID uint, personData *models.Person) error {
	var user models.User
	if err := r.db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}

	if err := r.db.DB.Model(&models.Person{}).Where("id = ?", user.PersonID).Updates(personData).Error; err != nil {
		return fmt.Errorf("error updating person data: %w", err)
	}

	return nil
}

func (r *UserRepository) ValidatePassword(userID uint, password string) (bool, error) {
	var user models.User
	if err := r.db.DB.Preload("Person").Where("id = ?", userID).First(&user).Error; err != nil {
		return false, fmt.Errorf("error finding user: %w", err)
	}

	return user.Person.Password == password, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	return r.FindByPersonEmail(email)
}

func (r *UserRepository) Create(user *models.User) error {
	return fmt.Errorf("use CreateWithPerson instead")
}

func (r *UserRepository) FarmExists(farmID uint) (bool, error) {
	var count int64
	if err := r.db.DB.Model(&models.Farm{}).Where("id = ?", farmID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("error checking farm existence: %w", err)
	}
	return count > 0, nil
}

func (r *UserRepository) CreateDefaultFarm(farmID uint) error {
	company := &models.Company{
		CompanyName: "FazendaPro Demo",
	}
	if err := r.db.DB.Create(company).Error; err != nil {
		return fmt.Errorf("error creating company: %w", err)
	}

	farm := &models.Farm{
		ID:        farmID,
		CompanyID: company.ID,
		Logo:      "",
	}
	if err := r.db.DB.Create(farm).Error; err != nil {
		return fmt.Errorf("error creating farm: %w", err)
	}

	return nil
}
