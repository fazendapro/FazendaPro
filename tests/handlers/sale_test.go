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
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSaleService struct {
	mock.Mock
}

func (m *MockSaleService) CreateSale(ctx context.Context, sale *models.Sale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockSaleService) GetSaleByID(ctx context.Context, id uint) (*models.Sale, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Sale), args.Error(1)
}

func (m *MockSaleService) GetSalesByFarmID(ctx context.Context, farmID uint) ([]*models.Sale, error) {
	args := m.Called(ctx, farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleService) GetSalesByAnimalID(ctx context.Context, animalID uint) ([]*models.Sale, error) {
	args := m.Called(ctx, animalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleService) GetSalesByDateRange(ctx context.Context, farmID uint, startDate, endDate time.Time) ([]*models.Sale, error) {
	args := m.Called(ctx, farmID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleService) GetMonthlySalesCount(ctx context.Context, farmID uint, startDate, endDate time.Time) (int64, error) {
	args := m.Called(ctx, farmID, startDate, endDate)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSaleService) UpdateSale(ctx context.Context, sale *models.Sale) error {
	args := m.Called(ctx, sale)
	return args.Error(0)
}

func (m *MockSaleService) DeleteSale(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSaleService) GetSalesHistory(ctx context.Context, farmID uint) ([]*models.Sale, error) {
	args := m.Called(ctx, farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Sale), args.Error(1)
}

func (m *MockSaleService) GetOverviewStats(ctx context.Context, farmID uint) (*repository.OverviewStats, error) {
	args := m.Called(ctx, farmID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.OverviewStats), args.Error(1)
}

func setupSaleTestRouter() (*chi.Mux, *MockSaleService) {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	mockSaleService := new(MockSaleService)
	saleHandler := handlers.NewSaleChiHandler(mockSaleService)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "farm_id", uint(1))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	router.Post("/sales", saleHandler.CreateSale)
	router.Get("/sales/{id}", saleHandler.GetSaleByID)
	router.Get("/sales", saleHandler.GetSalesByFarm)
	router.Get("/animals/{animal_id}/sales", saleHandler.GetSalesByAnimal)
	router.Get("/sales/date-range", saleHandler.GetSalesByDateRange)
	router.Put("/sales/{id}", saleHandler.UpdateSale)
	router.Delete("/sales/{id}", saleHandler.DeleteSale)
	router.Get("/sales/history", saleHandler.GetSalesHistory)
	router.Get("/sales/monthly-stats", saleHandler.GetMonthlySalesStats)
	router.Get("/sales/overview", saleHandler.GetOverviewStats)

	return router, mockSaleService
}

func TestSaleHandler_CreateSale_Success(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	now := time.Now()

	mockService.On("CreateSale", mock.Anything, mock.AnythingOfType("*models.Sale")).Return(nil).Run(func(args mock.Arguments) {
		sale := args.Get(1).(*models.Sale)
		sale.ID = 1
		sale.CreatedAt = now
		sale.UpdatedAt = now
	})

	reqBody := map[string]interface{}{
		"animal_id":  1,
		"buyer_name": "João Silva",
		"price":      1500.50,
		"sale_date":  now.Format("2006-01-02"),
		"notes":      "Test sale",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/sales", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "João Silva", response["buyer_name"])
	assert.Equal(t, 1500.50, response["price"])

	mockService.AssertExpectations(t)
}

func TestSaleHandler_CreateSale_InvalidData(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	reqBody := map[string]interface{}{
		"animal_id": 1,
		"price":     1500.50,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/sales", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid request data")

	mockService.AssertExpectations(t)
}

func TestSaleHandler_CreateSale_InvalidDate(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	reqBody := map[string]interface{}{
		"animal_id":  1,
		"buyer_name": "João Silva",
		"price":      1500.50,
		"sale_date":  "invalid-date",
		"notes":      "Test sale",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/sales", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid date format")

	mockService.AssertExpectations(t)
}

func TestSaleHandler_CreateSale_ServiceError(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	mockService.On("CreateSale", mock.Anything, mock.AnythingOfType("*models.Sale")).Return(errors.New("animal not found"))

	now := time.Now()
	reqBody := map[string]interface{}{
		"animal_id":  1,
		"buyer_name": "João Silva",
		"price":      1500.50,
		"sale_date":  now.Format("2006-01-02"),
		"notes":      "Test sale",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/sales", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Failed to create sale")
	assert.Contains(t, response["message"], "animal not found")

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetSaleByID_Success(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	now := time.Now()
	expectedSale := &models.Sale{
		ID:        1,
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  now,
		Notes:     "Test sale",
		CreatedAt: now,
		UpdatedAt: now,
		Animal: models.Animal{
			ID:         1,
			AnimalName: "Test Animal",
		},
	}

	mockService.On("GetSaleByID", mock.Anything, uint(1)).Return(expectedSale, nil)

	req, _ := http.NewRequest("GET", "/sales/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "João Silva", response["buyer_name"])
	assert.Equal(t, 1500.50, response["price"])

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetSaleByID_NotFound(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	mockService.On("GetSaleByID", mock.Anything, uint(1)).Return(nil, errors.New("sale not found"))

	req, _ := http.NewRequest("GET", "/sales/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Sale not found")

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetSaleByID_InvalidID(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	req, _ := http.NewRequest("GET", "/sales/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid sale ID")

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetSalesByFarm_Success(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	now := time.Now()
	expectedSales := []*models.Sale{
		{
			ID:        1,
			AnimalID:  1,
			FarmID:    1,
			BuyerName: "João Silva",
			Price:     1500.50,
			SaleDate:  now,
			Notes:     "Test sale 1",
			CreatedAt: now,
			UpdatedAt: now,
			Animal: models.Animal{
				ID:         1,
				AnimalName: "Test Animal 1",
			},
		},
		{
			ID:        2,
			AnimalID:  2,
			FarmID:    1,
			BuyerName: "Maria Santos",
			Price:     2000.00,
			SaleDate:  now.Add(-24 * time.Hour),
			Notes:     "Test sale 2",
			CreatedAt: now.Add(-24 * time.Hour),
			UpdatedAt: now.Add(-24 * time.Hour),
			Animal: models.Animal{
				ID:         2,
				AnimalName: "Test Animal 2",
			},
		},
	}

	mockService.On("GetSalesByFarmID", mock.Anything, uint(1)).Return(expectedSales, nil)

	req, _ := http.NewRequest("GET", "/sales", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, float64(1), response[0]["id"])
	assert.Equal(t, "João Silva", response[0]["buyer_name"])
	assert.Equal(t, float64(2), response[1]["id"])
	assert.Equal(t, "Maria Santos", response[1]["buyer_name"])

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetSalesByAnimal_Success(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	now := time.Now()
	expectedSales := []*models.Sale{
		{
			ID:        1,
			AnimalID:  1,
			FarmID:    1,
			BuyerName: "João Silva",
			Price:     1500.50,
			SaleDate:  now,
			Notes:     "Test sale",
			CreatedAt: now,
			UpdatedAt: now,
			Animal: models.Animal{
				ID:         1,
				AnimalName: "Test Animal",
			},
		},
	}

	mockService.On("GetSalesByAnimalID", mock.Anything, uint(1)).Return(expectedSales, nil)

	req, _ := http.NewRequest("GET", "/animals/1/sales", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, float64(1), response[0]["id"])
	assert.Equal(t, "João Silva", response[0]["buyer_name"])

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetSalesByDateRange_Success(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	now := time.Now()
	startDate := now.Add(-7 * 24 * time.Hour)
	endDate := now

	expectedSales := []*models.Sale{
		{
			ID:        1,
			AnimalID:  1,
			FarmID:    1,
			BuyerName: "João Silva",
			Price:     1500.50,
			SaleDate:  now.Add(-3 * 24 * time.Hour),
			Notes:     "Test sale",
			CreatedAt: now.Add(-3 * 24 * time.Hour),
			UpdatedAt: now.Add(-3 * 24 * time.Hour),
			Animal: models.Animal{
				ID:         1,
				AnimalName: "Test Animal",
			},
		},
	}

	mockService.On("GetSalesByDateRange", mock.Anything, uint(1), startDate, endDate).Return(expectedSales, nil)

	req, _ := http.NewRequest("GET", "/sales/date-range?start_date=2024-01-01&end_date=2024-01-08", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, float64(1), response[0]["id"])
	assert.Equal(t, "João Silva", response[0]["buyer_name"])

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetSalesByDateRange_MissingParams(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	req, _ := http.NewRequest("GET", "/sales/date-range", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Missing date parameters")

	mockService.AssertExpectations(t)
}

func TestSaleHandler_UpdateSale_Success(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	now := time.Now()
	existingSale := &models.Sale{
		ID:        1,
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "Original Buyer",
		Price:     1000.00,
		SaleDate:  now.Add(-24 * time.Hour),
		Notes:     "Original notes",
		CreatedAt: now.Add(-24 * time.Hour),
		UpdatedAt: now.Add(-24 * time.Hour),
	}

	mockService.On("GetSaleByID", mock.Anything, uint(1)).Return(existingSale, nil)
	mockService.On("UpdateSale", mock.Anything, mock.AnythingOfType("*models.Sale")).Return(nil)

	reqBody := map[string]interface{}{
		"buyer_name": "Updated Buyer",
		"price":      2000.00,
		"sale_date":  now.Format("2006-01-02"),
		"notes":      "Updated notes",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("PUT", "/sales/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "Updated Buyer", response["buyer_name"])
	assert.Equal(t, 2000.00, response["price"])

	mockService.AssertExpectations(t)
}

func TestSaleHandler_DeleteSale_Success(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	mockService.On("DeleteSale", mock.Anything, uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/sales/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "Sale deleted successfully")

	mockService.AssertExpectations(t)
}

func TestSaleHandler_DeleteSale_NotFound(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	mockService.On("DeleteSale", mock.Anything, uint(1)).Return(errors.New("sale not found"))

	req, _ := http.NewRequest("DELETE", "/sales/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Failed to delete sale")

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetSalesHistory_Success(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	now := time.Now()
	expectedSales := []*models.Sale{
		{
			ID:        1,
			AnimalID:  1,
			FarmID:    1,
			BuyerName: "João Silva",
			Price:     1500.50,
			SaleDate:  now,
			Notes:     "Test sale",
			CreatedAt: now,
			UpdatedAt: now,
			Animal: models.Animal{
				ID:         1,
				AnimalName: "Test Animal",
			},
		},
	}

	mockService.On("GetSalesHistory", mock.Anything, uint(1)).Return(expectedSales, nil)

	req, _ := http.NewRequest("GET", "/sales/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, float64(1), response[0]["id"])
	assert.Equal(t, "João Silva", response[0]["buyer_name"])

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetMonthlySalesStats_Success(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
	expectedCount := int64(5)

	mockService.On("GetMonthlySalesCount", mock.Anything, uint(1), startOfMonth, endOfMonth).Return(expectedCount, nil)

	req, _ := http.NewRequest("GET", "/sales/monthly-stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.NotNil(t, response["data"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(5), data["count"])

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetMonthlySalesStats_ServiceError(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	mockService.On("GetMonthlySalesCount", mock.Anything, uint(1), startOfMonth, endOfMonth).Return(int64(0), errors.New("database error"))

	req, _ := http.NewRequest("GET", "/sales/monthly-stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetOverviewStats_Success(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	expectedStats := &repository.OverviewStats{
		MalesCount:   10,
		FemalesCount: 15,
		TotalSold:    5,
		TotalRevenue: 15000.50,
	}

	mockService.On("GetOverviewStats", mock.Anything, uint(1)).Return(expectedStats, nil)

	req, _ := http.NewRequest("GET", "/sales/overview", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.NotNil(t, response["data"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(10), data["males_count"])
	assert.Equal(t, float64(15), data["females_count"])
	assert.Equal(t, float64(5), data["total_sold"])
	assert.Equal(t, 15000.50, data["total_revenue"])

	mockService.AssertExpectations(t)
}

func TestSaleHandler_GetOverviewStats_ServiceError(t *testing.T) {
	router, mockService := setupSaleTestRouter()

	mockService.On("GetOverviewStats", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

	req, _ := http.NewRequest("GET", "/sales/overview", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}
