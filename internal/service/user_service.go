package service

import (
	"errors"
	"strings"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email")
	}

	return s.repository.FindByPersonEmail(email)
}

func (s *UserService) CreateUser(user *models.User, personData *models.Person) error {
	if personData.FirstName == "" {
		return errors.New("first name is required")
	}

	if personData.LastName == "" {
		return errors.New("last name is required")
	}

	if personData.Email == "" {
		return errors.New("email is required")
	}

	if !strings.Contains(personData.Email, "@") {
		return errors.New("invalid email")
	}

	if personData.Password == "" {
		return errors.New("password is required")
	}

	if len(personData.Password) < 6 {
		return errors.New("password must have at least 6 characters")
	}

	if personData.CPF == "" {
		return errors.New("CPF is required")
	}

	if user.FarmID == 0 {
		return errors.New("farm ID is required")
	}

	return s.repository.CreateWithPerson(user, personData)
}

func (s *UserService) GetUserWithPerson(userID uint) (*models.User, error) {
	return s.repository.FindByIDWithPerson(userID)
}

func (s *UserService) UpdatePersonData(userID uint, personData *models.Person) error {
	return s.repository.UpdatePersonData(userID, personData)
}

func (s *UserService) ValidatePassword(userID uint, password string) (bool, error) {
	return s.repository.ValidatePassword(userID, password)
}
