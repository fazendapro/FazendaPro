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
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVaccineRepository struct {
	mock.Mock
}

func (m *MockVaccineRepository) Create(vaccine *models.Vaccine) error {
	args := m.Called(vaccine)
	return args.Error(0)
}

func (m *MockVaccineRepository) FindByID(id uint) (*models.Vaccine, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Vaccine), args.Error(1)
}

func (m *MockVaccineRepository) FindByFarmID(farmID uint) ([]models.Vaccine, error) {
	args := m.Called(farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Vaccine), args.Error(1)
}

func (m *MockVaccineRepository) Update(vaccine *models.Vaccine) error {
	args := m.Called(vaccine)
	return args.Error(0)
}

func (m *MockVaccineRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

var _ repository.VaccineRepositoryInterface = (*MockVaccineRepository)(nil)

func setupVaccineRouter(mockRepo *MockVaccineRepository) *chi.Mux {
	svc := service.NewVaccineService(mockRepo)
	handler := handlers.NewVaccineHandler(svc)
	router := chi.NewRouter()
	router.Post("/vaccines", handler.CreateVaccine)
	router.Get("/vaccines/farm/{farmId}", handler.GetVaccinesByFarmID)
	router.Get("/vaccines/{id}", handler.GetVaccineByID)
	router.Put("/vaccines/{id}", handler.UpdateVaccine)
	router.Delete("/vaccines/{id}", handler.DeleteVaccine)
	return router
}

func TestVaccineHandler_CreateVaccine_Success(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	router := setupVaccineRouter(mockRepo)

	reqBody := handlers.CreateVaccineRequest{
		FarmID:       1,
		Name:         "Vacina A",
		Description:  "Desc",
		Manufacturer: "Fab",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/vaccines", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	created := &models.Vaccine{
		ID:           10,
		FarmID:       1,
		Name:         "Vacina A",
		Description:  "Desc",
		Manufacturer: "Fab",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.
		On("Create", mock.AnythingOfType("*models.Vaccine")).
		Run(func(args mock.Arguments) {
			v := args.Get(0).(*models.Vaccine)
			v.ID = created.ID
		}).
		Return(nil)
	mockRepo.On("FindByID", uint(10)).Return(created, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp handlers.VaccineResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "Vacina A", resp.Data.Name)
	mockRepo.AssertExpectations(t)
}

func TestVaccineHandler_GetVaccinesByFarm_InvalidID(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	router := setupVaccineRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/vaccines/farm/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVaccineHandler_GetVaccineByID_NotFound(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	router := setupVaccineRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/vaccines/5", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(5)).Return(nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestVaccineHandler_UpdateVaccine_Success(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	router := setupVaccineRouter(mockRepo)

	reqBody := handlers.CreateVaccineRequest{
		FarmID:       1,
		Name:         "Atualizada",
		Description:  "Nova",
		Manufacturer: "Fab",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", "/vaccines/5", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	updated := &models.Vaccine{
		ID:           5,
		FarmID:       1,
		Name:         "Atualizada",
		Description:  "Nova",
		Manufacturer: "Fab",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("FindByID", uint(5)).Return(updated, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Vaccine")).Return(nil)
	mockRepo.On("FindByID", uint(5)).Return(updated, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp handlers.VaccineResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Atualizada", resp.Data.Name)
	mockRepo.AssertExpectations(t)
}

func TestVaccineHandler_DeleteVaccine_Success(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	router := setupVaccineRouter(mockRepo)

	req, _ := http.NewRequest("DELETE", "/vaccines/3", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(3)).Return(&models.Vaccine{}, nil)
	mockRepo.On("Delete", uint(3)).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}
