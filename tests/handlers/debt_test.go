package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/fazendapro/FazendaPro-api/tests/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDebtRouter(mockRepo *services.MockDebtRepository) (*chi.Mux, *services.MockDebtRepository) {
	debtService := service.NewDebtService(mockRepo)
	debtHandler := handlers.NewDebtHandler(debtService)
	r := chi.NewRouter()
	r.Post("/debts", debtHandler.CreateDebt)
	r.Get("/debts", debtHandler.GetDebts)
	r.Delete("/debts/{id}", debtHandler.DeleteDebt)
	r.Get("/debts/total-by-person", debtHandler.GetTotalByPerson)
	return r, mockRepo
}

func TestDebtHandler_CreateDebt_Success(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	debtData := map[string]interface{}{
		"person": "João Silva",
		"value":  1500.50,
	}

	jsonData, _ := json.Marshal(debtData)
	req, _ := http.NewRequest("POST", "/debts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("Create", mock.AnythingOfType("*models.Debt")).Return(nil).Run(func(args mock.Arguments) {
		debt := args.Get(0).(*models.Debt)
		debt.ID = 1
		debt.CreatedAt = time.Now()
		debt.UpdatedAt = time.Now()
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response handlers.DebtResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "João Silva", response.Person)
	assert.Equal(t, 1500.50, response.Value)
	mockRepo.AssertExpectations(t)
}

func TestDebtHandler_CreateDebt_InvalidJSON(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	req, _ := http.NewRequest("POST", "/debts", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDebtHandler_CreateDebt_ServiceError(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	debtData := map[string]interface{}{
		"person": "",
		"value":  1500.50,
	}

	jsonData, _ := json.Marshal(debtData)
	req, _ := http.NewRequest("POST", "/debts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDebtHandler_GetDebts_Success(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	expectedDebts := []models.Debt{
		{ID: 1, Person: "João Silva", Value: 1500.50, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Person: "Maria Santos", Value: 2000.00, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	req, _ := http.NewRequest("GET", "/debts?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindAllWithPagination", 1, 10, (*int)(nil), (*int)(nil)).Return(expectedDebts, int64(2), nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response handlers.DebtListResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response.Debts, 2)
	assert.Equal(t, int64(2), response.Total)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.Limit)
	mockRepo.AssertExpectations(t)
}

func TestDebtHandler_GetDebts_WithFilters(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	expectedDebts := []models.Debt{
		{ID: 1, Person: "João Silva", Value: 1500.50, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	year := 2024
	month := 1
	req, _ := http.NewRequest("GET", "/debts?page=1&limit=10&year=2024&month=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindAllWithPagination", 1, 10, &year, &month).Return(expectedDebts, int64(1), nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestDebtHandler_GetDebts_ServiceError(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/debts", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindAllWithPagination", 1, 10, (*int)(nil), (*int)(nil)).Return(nil, int64(0), errors.New("erro ao buscar"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestDebtHandler_DeleteDebt_Success(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	req, _ := http.NewRequest("DELETE", "/debts/1", nil)
	w := httptest.NewRecorder()

	existingDebt := &models.Debt{ID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existingDebt, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Dívida deletada com sucesso", response["message"])
	mockRepo.AssertExpectations(t)
}

func TestDebtHandler_DeleteDebt_MissingID(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)
	debtHandler := handlers.NewDebtHandler(debtService)

	req, _ := http.NewRequest("DELETE", "/debts/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	debtHandler.DeleteDebt(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDebtHandler_DeleteDebt_InvalidID(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)
	debtHandler := handlers.NewDebtHandler(debtService)

	req, _ := http.NewRequest("DELETE", "/debts/invalid", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	debtHandler.DeleteDebt(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDebtHandler_DeleteDebt_ServiceError(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	req, _ := http.NewRequest("DELETE", "/debts/1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("dívida não encontrada"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestDebtHandler_GetTotalByPerson_Success(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	expectedTotals := []repository.PersonTotal{
		{Person: "João Silva", Total: 1500.50},
		{Person: "Maria Santos", Total: 2000.00},
	}

	req, _ := http.NewRequest("GET", "/debts/total-by-person?year=2024&month=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("GetTotalByPersonInMonth", 2024, 1).Return(expectedTotals, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(2024), response["year"])
	assert.Equal(t, float64(1), response["month"])
	mockRepo.AssertExpectations(t)
}

func TestDebtHandler_GetTotalByPerson_MissingYear(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/debts/total-by-person?month=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDebtHandler_GetTotalByPerson_MissingMonth(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/debts/total-by-person?year=2024", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDebtHandler_GetTotalByPerson_InvalidYear(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/debts/total-by-person?year=invalid&month=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDebtHandler_GetTotalByPerson_InvalidMonth(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/debts/total-by-person?year=2024&month=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDebtHandler_GetTotalByPerson_ServiceError(t *testing.T) {
	mockRepo := new(services.MockDebtRepository)
	router, _ := setupDebtRouter(mockRepo)

	req, _ := http.NewRequest("GET", "/debts/total-by-person?year=2024&month=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("GetTotalByPersonInMonth", 2024, 1).Return(nil, errors.New("erro ao calcular"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}
