package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/utils"
)

type UserService struct {
	repository repository.UserRepositoryInterface
}

func NewUserService(repository repository.UserRepositoryInterface) *UserService {
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

	hashedPassword, err := utils.HashPassword(personData.Password)
	if err != nil {
		return errors.New("error hashing password")
	}
	personData.Password = hashedPassword

	if personData.CPF == "" {
		return errors.New("CPF is required")
	}

	if user.FarmID == 0 {
		return errors.New("farm ID is required")
	}

	if err := s.ensureFarmExists(user.FarmID); err != nil {
		return fmt.Errorf("error ensuring farm exists: %w", err)
	}

	if err := s.repository.CreateWithPerson(user, personData); err != nil {
		return err
	}

	userFarm := &models.UserFarm{
		UserID:    user.ID,
		FarmID:    user.FarmID,
		IsPrimary: true,
	}

	if err := s.repository.CreateUserFarm(userFarm); err != nil {
		return err
	}

	userWithPerson, err := s.repository.FindByIDWithPerson(user.ID)
	if err != nil {
		return err
	}

	*user = *userWithPerson

	return nil
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

func (s *UserService) ValidatePasswordByEmail(email, password string) (bool, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}

	return utils.CheckPasswordHash(password, user.Person.Password), nil
}

func (s *UserService) ensureFarmExists(farmID uint) error {
	exists, err := s.repository.FarmExists(farmID)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return s.repository.CreateDefaultFarm(farmID)
}

func (s *UserService) GetUserFarms(userID uint) ([]models.Farm, error) {
	return s.repository.GetUserFarms(userID)
}

func (s *UserService) GetUserFarmCount(userID uint) (int64, error) {
	return s.repository.GetUserFarmCount(userID)
}

func (s *UserService) GetUserFarmByID(userID, farmID uint) (*models.Farm, error) {
	return s.repository.GetUserFarmByID(userID, farmID)
}

func (s *UserService) ShouldAutoSelectFarm(userID uint) (bool, error) {
	count, err := s.GetUserFarmCount(userID)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}
