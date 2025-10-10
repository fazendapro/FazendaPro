package services

import (
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateWithPerson(user *models.User, person *models.Person) error {
	args := m.Called(user, person)
	return args.Error(0)
}

func (m *MockUserRepository) FindByPersonEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByIDWithPerson(userID uint) (*models.User, error) {
	args := m.Called(userID)
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
	return args.Get(0).([]models.Farm), args.Error(1)
}

func (m *MockUserRepository) GetUserFarmCount(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) GetUserFarmByID(userID, farmID uint) (*models.Farm, error) {
	args := m.Called(userID, farmID)
	return args.Get(0).(*models.Farm), args.Error(1)
}

func TestUserService_GetUserFarms(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uint(1)
	farms := []models.Farm{
		{ID: 1, Logo: "logo1.png"},
		{ID: 2, Logo: "logo2.png"},
	}

	mockRepo.On("GetUserFarms", userID).Return(farms, nil)

	result, err := userService.GetUserFarms(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, farms, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserFarmCount(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uint(1)
	count := int64(2)

	mockRepo.On("GetUserFarmCount", userID).Return(count, nil)

	result, err := userService.GetUserFarmCount(userID)

	assert.NoError(t, err)
	assert.Equal(t, count, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserFarmByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uint(1)
	farmID := uint(2)
	farm := &models.Farm{ID: farmID, Logo: "logo2.png"}

	mockRepo.On("GetUserFarmByID", userID, farmID).Return(farm, nil)

	result, err := userService.GetUserFarmByID(userID, farmID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, farm, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserFarmByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uint(1)
	farmID := uint(999)

	mockRepo.On("GetUserFarmByID", userID, farmID).Return((*models.Farm)(nil), gorm.ErrRecordNotFound)

	result, err := userService.GetUserFarmByID(userID, farmID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_ShouldAutoSelectFarm_SingleFarm(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uint(1)
	count := int64(1)

	mockRepo.On("GetUserFarmCount", userID).Return(count, nil)

	shouldAutoSelect, err := userService.ShouldAutoSelectFarm(userID)

	assert.NoError(t, err)
	assert.True(t, shouldAutoSelect)
	mockRepo.AssertExpectations(t)
}

func TestUserService_ShouldAutoSelectFarm_MultipleFarms(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uint(1)
	count := int64(3)

	mockRepo.On("GetUserFarmCount", userID).Return(count, nil)

	shouldAutoSelect, err := userService.ShouldAutoSelectFarm(userID)

	assert.NoError(t, err)
	assert.False(t, shouldAutoSelect)
	mockRepo.AssertExpectations(t)
}

func TestUserService_ShouldAutoSelectFarm_NoFarms(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uint(1)
	count := int64(0)

	mockRepo.On("GetUserFarmCount", userID).Return(count, nil)

	shouldAutoSelect, err := userService.ShouldAutoSelectFarm(userID)

	assert.NoError(t, err)
	assert.False(t, shouldAutoSelect)
	mockRepo.AssertExpectations(t)
}
