package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user *models.User, person *models.Person) error {
	args := m.Called(user, person)
	return args.Error(0)
}

func (m *MockUserService) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserWithPerson(userID uint) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdatePersonData(userID uint, person *models.Person) error {
	args := m.Called(userID, person)
	return args.Error(0)
}

func (m *MockUserService) ValidatePasswordByEmail(email, password string) (bool, error) {
	args := m.Called(email, password)
	return args.Bool(0), args.Error(1)
}

type MockAnimalService struct {
	mock.Mock
}

func (m *MockAnimalService) CreateAnimal(animal *models.Animal) error {
	args := m.Called(animal)
	return args.Error(0)
}

func (m *MockAnimalService) GetAnimalByID(id uint) (*models.Animal, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Animal), args.Error(1)
}

func (m *MockAnimalService) GetAnimalsByFarmID(farmID uint) ([]models.Animal, error) {
	args := m.Called(farmID)
	return args.Get(0).([]models.Animal), args.Error(1)
}

func (m *MockAnimalService) UpdateAnimal(animal *models.Animal) error {
	args := m.Called(animal)
	return args.Error(0)
}

func (m *MockAnimalService) DeleteAnimal(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockMilkCollectionService struct {
	mock.Mock
}

func (m *MockMilkCollectionService) CreateMilkCollection(mc *models.MilkCollection) error {
	args := m.Called(mc)
	return args.Error(0)
}

func (m *MockMilkCollectionService) UpdateMilkCollection(mc *models.MilkCollection) error {
	args := m.Called(mc)
	return args.Error(0)
}

func (m *MockMilkCollectionService) GetMilkCollectionByID(id uint) (*models.MilkCollection, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MilkCollection), args.Error(1)
}

func (m *MockMilkCollectionService) GetMilkCollectionsByFarmID(farmID uint) ([]models.MilkCollection, error) {
	args := m.Called(farmID)
	return args.Get(0).([]models.MilkCollection), args.Error(1)
}

func (m *MockMilkCollectionService) GetMilkCollectionsByFarmIDWithDateRange(farmID uint, startDate, endDate *time.Time) ([]models.MilkCollection, error) {
	args := m.Called(farmID, startDate, endDate)
	return args.Get(0).([]models.MilkCollection), args.Error(1)
}

func (m *MockMilkCollectionService) GetMilkCollectionsByAnimalID(animalID uint) ([]models.MilkCollection, error) {
	args := m.Called(animalID)
	return args.Get(0).([]models.MilkCollection), args.Error(1)
}

type MockReproductionService struct {
	mock.Mock
}

func (m *MockReproductionService) CreateReproduction(reproduction *models.Reproduction) error {
	args := m.Called(reproduction)
	return args.Error(0)
}

func (m *MockReproductionService) GetReproductionByID(id uint) (*models.Reproduction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reproduction), args.Error(1)
}

func (m *MockReproductionService) GetReproductionByAnimalID(animalID uint) (*models.Reproduction, error) {
	args := m.Called(animalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reproduction), args.Error(1)
}

func (m *MockReproductionService) GetReproductionsByFarmID(farmID uint) ([]models.Reproduction, error) {
	args := m.Called(farmID)
	return args.Get(0).([]models.Reproduction), args.Error(1)
}

func (m *MockReproductionService) GetReproductionsByPhase(phase models.ReproductionPhase) ([]models.Reproduction, error) {
	args := m.Called(phase)
	return args.Get(0).([]models.Reproduction), args.Error(1)
}

func (m *MockReproductionService) UpdateReproduction(reproduction *models.Reproduction) error {
	args := m.Called(reproduction)
	return args.Error(0)
}

func (m *MockReproductionService) UpdateReproductionPhase(animalID uint, phase models.ReproductionPhase, additionalData map[string]interface{}) error {
	args := m.Called(animalID, phase, additionalData)
	return args.Error(0)
}

func (m *MockReproductionService) DeleteReproduction(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(userID uint, expiresAt time.Time) (*models.RefreshToken, error) {
	args := m.Called(userID, expiresAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) DeleteByToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteByUserID(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

func createTestRequest(method, url string, body interface{}) *http.Request {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func assertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(t, expectedStatus, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func parseJSONResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), target)
	assert.NoError(t, err)
}

func SetupTestDB(t *testing.T) *repository.Database {
	return nil
}

func CleanupTestDB(t *testing.T, db *repository.Database) {
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
