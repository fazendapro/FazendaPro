package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
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

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
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

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Person{}, &models.Farm{}, &models.Company{})
	return db
}

func generateTestToken(userID uint) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": 1234567890,
		"exp": 1234567890 + 86400,
	})
	tokenString, _ := token.SignedString([]byte("test-secret"))
	return tokenString
}

func TestGetUserFarms_SingleFarm(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService)

	userID := uint(1)
	farm := models.Farm{
		ID:   1,
		Logo: "logo1.png",
	}

	mockRepo.On("GetUserFarms", userID).Return([]models.Farm{farm}, nil)
	mockRepo.On("GetUserFarmCount", userID).Return(int64(1), nil)

	req := httptest.NewRequest("GET", "/api/v1/farms/user", nil)
	req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
	w := httptest.NewRecorder()

	farmHandler.GetUserFarms(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response["success"].(bool))
	assert.Equal(t, 1, len(response["farms"].([]interface{})))
	assert.True(t, response["auto_select"].(bool))
	assert.Equal(t, float64(1), response["selected_farm_id"])
}

func TestGetUserFarms_MultipleFarms(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService)

	userID := uint(1)
	farms := []models.Farm{
		{ID: 1, Logo: "logo1.png"},
		{ID: 2, Logo: "logo2.png"},
	}

	mockRepo.On("GetUserFarms", userID).Return(farms, nil)
	mockRepo.On("GetUserFarmCount", userID).Return(int64(2), nil)

	req := httptest.NewRequest("GET", "/api/v1/farms/user", nil)
	req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
	w := httptest.NewRecorder()

	farmHandler.GetUserFarms(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response["success"].(bool))
	assert.Equal(t, 2, len(response["farms"].([]interface{})))
	assert.False(t, response["auto_select"].(bool))
	assert.Nil(t, response["selected_farm_id"])
}

func TestGetUserFarms_NoFarms(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService)

	userID := uint(1)

	mockRepo.On("GetUserFarms", userID).Return([]models.Farm{}, nil)
	mockRepo.On("GetUserFarmCount", userID).Return(int64(0), nil)

	req := httptest.NewRequest("GET", "/api/v1/farms/user", nil)
	req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
	w := httptest.NewRecorder()

	farmHandler.GetUserFarms(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response["success"].(bool))
	assert.Equal(t, 0, len(response["farms"].([]interface{})))
	assert.False(t, response["auto_select"].(bool))
}

func TestSelectFarm(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService)

	userID := uint(1)
	farmID := uint(2)
	farm := models.Farm{
		ID:   farmID,
		Logo: "logo2.png",
	}

	mockRepo.On("GetUserFarmByID", userID, farmID).Return(&farm, nil)

	requestBody := map[string]interface{}{
		"farm_id": farmID,
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/v1/farms/select", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	farmHandler.SelectFarm(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response["success"].(bool))
	assert.Equal(t, float64(farmID), response["farm_id"])
}

func TestSelectFarm_InvalidFarm(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService)

	userID := uint(1)
	farmID := uint(999)

	mockRepo.On("GetUserFarmByID", userID, farmID).Return((*models.Farm)(nil), gorm.ErrRecordNotFound)

	requestBody := map[string]interface{}{
		"farm_id": farmID,
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/v1/farms/select", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	farmHandler.SelectFarm(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"], "Fazenda não encontrada")
}
