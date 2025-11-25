package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/fazendapro/FazendaPro-api/tests"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Person{}, &models.Farm{}, &models.Company{})
	return db
}

func generateTestToken(userID uint) string {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": float64(userID),
		"iat": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(tests.TestSecret))
	return tokenString
}

func TestGetUserFarms_SingleFarm(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	userID := uint(1)
	farm := models.Farm{
		ID:   1,
		Logo: tests.TestFileLogo1PNG,
	}

	mockRepo.On("GetUserFarms", userID).Return([]models.Farm{farm}, nil)
	mockRepo.On("GetUserFarmCount", userID).Return(int64(1), nil)

	req := httptest.NewRequest("GET", tests.EndpointAPIv1FarmsUser, nil)
	req.Header.Set(tests.HeaderAuthorization, tests.BearerPrefix+generateTestToken(userID))
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
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	userID := uint(1)
	farms := []models.Farm{
		{ID: 1, Logo: "logo1.png"},
		{ID: 2, Logo: "logo2.png"},
	}

	mockRepo.On("GetUserFarms", userID).Return(farms, nil)
	mockRepo.On("GetUserFarmCount", userID).Return(int64(2), nil)

	req := httptest.NewRequest("GET", tests.EndpointAPIv1FarmsUser, nil)
	req.Header.Set(tests.HeaderAuthorization, tests.BearerPrefix+generateTestToken(userID))
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
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	userID := uint(1)

	mockRepo.On("GetUserFarms", userID).Return([]models.Farm{}, nil)
	mockRepo.On("GetUserFarmCount", userID).Return(int64(0), nil)

	req := httptest.NewRequest("GET", tests.EndpointAPIv1FarmsUser, nil)
	req.Header.Set(tests.HeaderAuthorization, tests.BearerPrefix+generateTestToken(userID))
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
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

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

	req := httptest.NewRequest("POST", tests.EndpointAPIv1FarmsSelect, bytes.NewBuffer(jsonBody))
	req.Header.Set(tests.HeaderAuthorization, tests.BearerPrefix+generateTestToken(userID))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
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
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	userID := uint(1)
	farmID := uint(999)

	mockRepo.On("GetUserFarmByID", userID, farmID).Return((*models.Farm)(nil), gorm.ErrRecordNotFound)

	requestBody := map[string]interface{}{
		"farm_id": farmID,
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", tests.EndpointAPIv1FarmsSelect, bytes.NewBuffer(jsonBody))
	req.Header.Set(tests.HeaderAuthorization, tests.BearerPrefix+generateTestToken(userID))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	farmHandler.SelectFarm(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"], "Fazenda não encontrada")
}

func TestGetUserFarms_InvalidMethod(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	req := httptest.NewRequest("POST", "/api/v1/farms/user", nil)
	req.Header.Set("Authorization", "Bearer "+generateTestToken(1))
	w := httptest.NewRecorder()

	farmHandler.GetUserFarms(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestGetUserFarms_InvalidToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	req := httptest.NewRequest("GET", "/api/v1/farms/user", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	farmHandler.GetUserFarms(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserFarms_MissingToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	req := httptest.NewRequest("GET", "/api/v1/farms/user", nil)
	w := httptest.NewRecorder()

	farmHandler.GetUserFarms(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserFarms_ServiceError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	userID := uint(1)
	mockRepo.On("GetUserFarms", userID).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", tests.EndpointAPIv1FarmsUser, nil)
	req.Header.Set(tests.HeaderAuthorization, tests.BearerPrefix+generateTestToken(userID))
	w := httptest.NewRecorder()

	farmHandler.GetUserFarms(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestSelectFarm_InvalidMethod(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	req := httptest.NewRequest("GET", "/api/v1/farms/select", nil)
	req.Header.Set("Authorization", "Bearer "+generateTestToken(1))
	w := httptest.NewRecorder()

	farmHandler.SelectFarm(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestSelectFarm_InvalidToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	requestBody := map[string]interface{}{
		"farm_id": 1,
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/v1/farms/select", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer invalid-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	farmHandler.SelectFarm(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSelectFarm_InvalidJSON(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	req := httptest.NewRequest("POST", "/api/v1/farms/select", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Authorization", "Bearer "+generateTestToken(1))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	farmHandler.SelectFarm(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSelectFarm_ServiceError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	userID := uint(1)
	farmID := uint(2)

	mockRepo.On("GetUserFarmByID", userID, farmID).Return(nil, gorm.ErrRecordNotFound)

	requestBody := map[string]interface{}{
		"farm_id": farmID,
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", tests.EndpointAPIv1FarmsSelect, bytes.NewBuffer(jsonBody))
	req.Header.Set(tests.HeaderAuthorization, tests.BearerPrefix+generateTestToken(userID))
	req.Header.Set(tests.HeaderContentType, tests.ContentTypeJSON)
	w := httptest.NewRecorder()

	farmHandler.SelectFarm(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetUserFarms_ShouldAutoSelectFarmError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	userID := uint(1)
	farms := []models.Farm{
		{ID: 1},
		{ID: 2},
	}

	mockRepo.On("GetUserFarms", userID).Return(farms, nil)
	mockRepo.On("GetUserFarmCount", userID).Return(int64(0), gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", tests.EndpointAPIv1Farms, nil)
	req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
	w := httptest.NewRecorder()

	farmHandler.GetUserFarms(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetUserFarms_EmptyFarmsList(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	userID := uint(1)

	mockRepo.On("GetUserFarms", userID).Return([]models.Farm{}, nil)
	mockRepo.On("GetUserFarmCount", userID).Return(int64(0), nil)

	req := httptest.NewRequest("GET", tests.EndpointAPIv1Farms, nil)
	req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
	w := httptest.NewRecorder()

	farmHandler.GetUserFarms(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response handlers.GetUserFarmsResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Len(t, response.Farms.([]interface{}), 0)
	assert.False(t, response.AutoSelect)
	assert.Nil(t, response.SelectedFarmID)
	mockRepo.AssertExpectations(t)
}

func TestExtractUserIDFromToken_InvalidClaims(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(tests.TestSecret))

	req := httptest.NewRequest("GET", tests.EndpointAPIv1Farms, nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	farmHandler.GetUserFarms(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestExtractUserIDFromToken_InvalidTokenType(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)
	farmHandler := handlers.NewFarmSelectionHandler(userService, tests.TestSecret)

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "invalid-type",
		"iat": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(tests.TestSecret))

	req := httptest.NewRequest("GET", tests.EndpointAPIv1Farms, nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	farmHandler.GetUserFarms(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
