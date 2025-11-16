package services

import (
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

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
