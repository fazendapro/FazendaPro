package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/go-chi/chi/v5"
)

type SaleChiHandler struct {
	service service.SaleService
}

func NewSaleChiHandler(service service.SaleService) *SaleChiHandler {
	return &SaleChiHandler{service: service}
}

type CreateSaleRequest struct {
	AnimalID  uint    `json:"animal_id"`
	FarmID    *uint   `json:"farm_id,omitempty"`
	BuyerName string  `json:"buyer_name"`
	Price     float64 `json:"price"`
	SaleDate  string  `json:"sale_date"`
	Notes     string  `json:"notes"`
}

type UpdateSaleRequest struct {
	BuyerName string  `json:"buyer_name"`
	Price     float64 `json:"price"`
	SaleDate  string  `json:"sale_date"`
	Notes     string  `json:"notes"`
}

type SaleResponse struct {
	ID        uint           `json:"id"`
	AnimalID  uint           `json:"animal_id"`
	FarmID    uint           `json:"farm_id"`
	BuyerName string         `json:"buyer_name"`
	Price     float64        `json:"price"`
	SaleDate  time.Time      `json:"sale_date"`
	Notes     string         `json:"notes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Animal    *models.Animal `json:"animal,omitempty"`
}

func (h *SaleChiHandler) CreateSale(w http.ResponseWriter, r *http.Request) {
	var req CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	saleDate, err := time.Parse(DateFormatISO, req.SaleDate)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	farmIDFromContext, ok := r.Context().Value("farm_id").(uint)
	if !ok || farmIDFromContext == 0 {
		SendErrorResponse(w, ErrFarmIDNotFound, http.StatusUnauthorized)
		return
	}

	farmID := farmIDFromContext

	if req.FarmID != nil && *req.FarmID != farmIDFromContext {
		SendErrorResponse(w, "farm_id do body não corresponde ao farm_id do token. Use o farm_id do token ou não envie farm_id no body", http.StatusForbidden)
		return
	}

	sale := &models.Sale{
		AnimalID:  req.AnimalID,
		FarmID:    farmID,
		BuyerName: req.BuyerName,
		Price:     req.Price,
		SaleDate:  saleDate,
		Notes:     req.Notes,
	}

	err = h.service.CreateSale(r.Context(), sale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := SaleResponse{
		ID:        sale.ID,
		AnimalID:  sale.AnimalID,
		FarmID:    sale.FarmID,
		BuyerName: sale.BuyerName,
		Price:     sale.Price,
		SaleDate:  sale.SaleDate,
		Notes:     sale.Notes,
		CreatedAt: sale.CreatedAt,
		UpdatedAt: sale.UpdatedAt,
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *SaleChiHandler) GetSalesByFarm(w http.ResponseWriter, r *http.Request) {
	farmIDStr := r.URL.Query().Get("farmId")
	var farmID uint

	if farmIDStr != "" {
		id, parseErr := strconv.ParseUint(farmIDStr, 10, 32)
		if parseErr != nil {
			SendErrorResponse(w, ErrInvalidFarmID, http.StatusBadRequest)
			return
		}
		farmID = uint(id)
	} else {
		farmIDFromContext, ok := r.Context().Value("farm_id").(uint)
		if !ok {
			SendErrorResponse(w, ErrFarmIDNotFound, http.StatusUnauthorized)
			return
		}
		farmID = farmIDFromContext
	}

	sales, err := h.service.GetSalesByFarmID(r.Context(), farmID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]SaleResponse, len(sales))
	for i, sale := range sales {
		responses[i] = SaleResponse{
			ID:        sale.ID,
			AnimalID:  sale.AnimalID,
			FarmID:    sale.FarmID,
			BuyerName: sale.BuyerName,
			Price:     sale.Price,
			SaleDate:  sale.SaleDate,
			Notes:     sale.Notes,
			CreatedAt: sale.CreatedAt,
			UpdatedAt: sale.UpdatedAt,
			Animal:    &sale.Animal,
		}
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(responses)
}

func (h *SaleChiHandler) GetSalesHistory(w http.ResponseWriter, r *http.Request) {
	farmIDStr := r.URL.Query().Get("farmId")
	var farmID uint

	if farmIDStr != "" {
		id, parseErr := strconv.ParseUint(farmIDStr, 10, 32)
		if parseErr != nil {
			SendErrorResponse(w, ErrInvalidFarmID, http.StatusBadRequest)
			return
		}
		farmID = uint(id)
	} else {
		farmIDFromContext, ok := r.Context().Value("farm_id").(uint)
		if !ok {
			SendErrorResponse(w, ErrFarmIDNotFound, http.StatusUnauthorized)
			return
		}
		farmID = farmIDFromContext
	}

	sales, err := h.service.GetSalesHistory(r.Context(), farmID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]SaleResponse, len(sales))
	for i, sale := range sales {
		responses[i] = SaleResponse{
			ID:        sale.ID,
			AnimalID:  sale.AnimalID,
			FarmID:    sale.FarmID,
			BuyerName: sale.BuyerName,
			Price:     sale.Price,
			SaleDate:  sale.SaleDate,
			Notes:     sale.Notes,
			CreatedAt: sale.CreatedAt,
			UpdatedAt: sale.UpdatedAt,
			Animal:    &sale.Animal,
		}
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(responses)
}

func (h *SaleChiHandler) GetSalesByDateRange(w http.ResponseWriter, r *http.Request) {
	farmIDStr := r.URL.Query().Get("farmId")
	var farmID uint

	if farmIDStr != "" {
		id, parseErr := strconv.ParseUint(farmIDStr, 10, 32)
		if parseErr != nil {
			SendErrorResponse(w, ErrInvalidFarmID, http.StatusBadRequest)
			return
		}
		farmID = uint(id)
	} else {
		farmIDFromContext, ok := r.Context().Value("farm_id").(uint)
		if !ok {
			SendErrorResponse(w, ErrFarmIDNotFound, http.StatusUnauthorized)
			return
		}
		farmID = farmIDFromContext
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "Both start_date and end_date are required", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Invalid end date format", http.StatusBadRequest)
		return
	}

	sales, err := h.service.GetSalesByDateRange(r.Context(), farmID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]SaleResponse, len(sales))
	for i, sale := range sales {
		responses[i] = SaleResponse{
			ID:        sale.ID,
			AnimalID:  sale.AnimalID,
			FarmID:    sale.FarmID,
			BuyerName: sale.BuyerName,
			Price:     sale.Price,
			SaleDate:  sale.SaleDate,
			Notes:     sale.Notes,
			CreatedAt: sale.CreatedAt,
			UpdatedAt: sale.UpdatedAt,
			Animal:    &sale.Animal,
		}
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(responses)
}

func (h *SaleChiHandler) GetSaleByID(w http.ResponseWriter, r *http.Request) {
	farmID, ok := r.Context().Value("farm_id").(uint)
	if !ok {
		SendErrorResponse(w, ErrFarmIDNotFound, http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, ErrInvalidSaleID, http.StatusBadRequest)
		return
	}

	sale, err := h.service.GetSaleByID(r.Context(), uint(id), farmID)
	if err != nil {
		SendErrorResponse(w, ErrSaleNotFound, http.StatusForbidden)
		return
	}

	response := SaleResponse{
		ID:        sale.ID,
		AnimalID:  sale.AnimalID,
		FarmID:    sale.FarmID,
		BuyerName: sale.BuyerName,
		Price:     sale.Price,
		SaleDate:  sale.SaleDate,
		Notes:     sale.Notes,
		CreatedAt: sale.CreatedAt,
		UpdatedAt: sale.UpdatedAt,
		Animal:    &sale.Animal,
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(response)
}

func (h *SaleChiHandler) UpdateSale(w http.ResponseWriter, r *http.Request) {
	farmID, ok := r.Context().Value("farm_id").(uint)
	if !ok {
		SendErrorResponse(w, ErrFarmIDNotFound, http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, ErrInvalidSaleID, http.StatusBadRequest)
		return
	}

	var req UpdateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	saleDate, err := time.Parse(DateFormatISO, req.SaleDate)
	if err != nil {
		SendErrorResponse(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	sale := &models.Sale{
		ID:        uint(id),
		FarmID:    farmID,
		BuyerName: req.BuyerName,
		Price:     req.Price,
		SaleDate:  saleDate,
		Notes:     req.Notes,
	}

	err = h.service.UpdateSale(r.Context(), sale, farmID)
	if err != nil {
		if err.Error() == service.ErrSaleNotFoundOrNotBelongsToFarm {
			SendErrorResponse(w, ErrSaleNotBelongsToFarm, http.StatusForbidden)
		} else {
			SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	updatedSale, err := h.service.GetSaleByID(r.Context(), uint(id), farmID)
	if err != nil {
		SendErrorResponse(w, ErrSaleNotFound, http.StatusInternalServerError)
		return
	}

	response := SaleResponse{
		ID:        updatedSale.ID,
		AnimalID:  updatedSale.AnimalID,
		FarmID:    updatedSale.FarmID,
		BuyerName: updatedSale.BuyerName,
		Price:     updatedSale.Price,
		SaleDate:  updatedSale.SaleDate,
		Notes:     updatedSale.Notes,
		CreatedAt: updatedSale.CreatedAt,
		UpdatedAt: updatedSale.UpdatedAt,
		Animal:    &updatedSale.Animal,
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(response)
}

func (h *SaleChiHandler) DeleteSale(w http.ResponseWriter, r *http.Request) {
	farmID, ok := r.Context().Value("farm_id").(uint)
	if !ok {
		SendErrorResponse(w, ErrFarmIDNotFound, http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, ErrInvalidSaleID, http.StatusBadRequest)
		return
	}

	err = h.service.DeleteSale(r.Context(), uint(id), farmID)
	if err != nil {
		if err.Error() == service.ErrSaleNotFoundOrNotBelongsToFarm {
			SendErrorResponse(w, ErrSaleNotBelongsToFarm, http.StatusForbidden)
		} else {
			SendErrorResponse(w, err.Error(), http.StatusNotFound)
		}
		return
	}

	SendSuccessResponse(w, nil, "Sale deleted successfully", http.StatusOK)
}

func (h *SaleChiHandler) GetSalesByAnimal(w http.ResponseWriter, r *http.Request) {
	farmID, ok := r.Context().Value("farm_id").(uint)
	if !ok {
		SendErrorResponse(w, ErrFarmIDNotFound, http.StatusUnauthorized)
		return
	}

	animalIDStr := chi.URLParam(r, "animal_id")
	animalID, err := strconv.ParseUint(animalIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "Invalid animal ID", http.StatusBadRequest)
		return
	}

	sales, err := h.service.GetSalesByAnimalID(r.Context(), uint(animalID), farmID)
	if err != nil {
		if err.Error() == "animal does not belong to the specified farm" {
			SendErrorResponse(w, ErrAnimalNotBelongsToFarm, http.StatusForbidden)
		} else {
			SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	responses := make([]SaleResponse, len(sales))
	for i, sale := range sales {
		responses[i] = SaleResponse{
			ID:        sale.ID,
			AnimalID:  sale.AnimalID,
			FarmID:    sale.FarmID,
			BuyerName: sale.BuyerName,
			Price:     sale.Price,
			SaleDate:  sale.SaleDate,
			Notes:     sale.Notes,
			CreatedAt: sale.CreatedAt,
			UpdatedAt: sale.UpdatedAt,
			Animal:    &sale.Animal,
		}
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(responses)
}

func (h *SaleChiHandler) GetMonthlySalesStats(w http.ResponseWriter, r *http.Request) {
	farmIDStr := r.URL.Query().Get("farmId")
	var farmID uint

	if farmIDStr != "" {
		id, parseErr := strconv.ParseUint(farmIDStr, 10, 32)
		if parseErr != nil {
			SendErrorResponse(w, ErrInvalidFarmID, http.StatusBadRequest)
			return
		}
		farmID = uint(id)
	} else {
		farmIDFromContext, ok := r.Context().Value("farm_id").(uint)
		if !ok {
			SendErrorResponse(w, ErrFarmIDNotFound, http.StatusUnauthorized)
			return
		}
		farmID = farmIDFromContext
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	count, err := h.service.GetMonthlySalesCount(r.Context(), farmID, startOfMonth, endOfMonth)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"count": count,
	}

	SendSuccessResponse(w, data, "Estatísticas mensais de vendas recuperadas com sucesso", http.StatusOK)
}

func (h *SaleChiHandler) GetMonthlySalesAndPurchases(w http.ResponseWriter, r *http.Request) {
	farmIDStr := r.URL.Query().Get("farmId")
	var farmID uint

	if farmIDStr != "" {
		id, parseErr := strconv.ParseUint(farmIDStr, 10, 32)
		if parseErr != nil {
			SendErrorResponse(w, ErrInvalidFarmID, http.StatusBadRequest)
			return
		}
		farmID = uint(id)
	} else {
		farmIDFromContext, ok := r.Context().Value("farm_id").(uint)
		if !ok {
			SendErrorResponse(w, ErrFarmIDNotFound, http.StatusUnauthorized)
			return
		}
		farmID = farmIDFromContext
	}

	monthsStr := r.URL.Query().Get("months")
	months := 12
	if monthsStr != "" {
		parsed, err := strconv.Atoi(monthsStr)
		if err != nil || parsed <= 0 || parsed > 24 {
			SendErrorResponse(w, ErrInvalidMonthsParam, http.StatusBadRequest)
			return
		}
		months = parsed
	}

	salesData, err := h.service.GetMonthlySalesData(r.Context(), farmID, months)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	purchasesData := make([]map[string]interface{}, len(salesData))
	for i := range purchasesData {
		purchasesData[i] = map[string]interface{}{
			"month": salesData[i].Month,
			"year":  salesData[i].Year,
			"total": 0.0,
		}
	}

	data := map[string]interface{}{
		"sales":     salesData,
		"purchases": purchasesData,
	}

	SendSuccessResponse(w, data, "Dados mensais de vendas e compras recuperados com sucesso", http.StatusOK)
}

func (h *SaleChiHandler) GetOverviewStats(w http.ResponseWriter, r *http.Request) {
	farmIDStr := r.URL.Query().Get("farmId")
	var farmID uint

	if farmIDStr != "" {
		id, parseErr := strconv.ParseUint(farmIDStr, 10, 32)
		if parseErr != nil {
			SendErrorResponse(w, ErrInvalidFarmID, http.StatusBadRequest)
			return
		}
		farmID = uint(id)
	} else {
		farmIDFromContext, ok := r.Context().Value("farm_id").(uint)
		if !ok {
			SendErrorResponse(w, ErrFarmIDNotFound, http.StatusUnauthorized)
			return
		}
		farmID = farmIDFromContext
	}

	stats, err := h.service.GetOverviewStats(r.Context(), farmID)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"males_count":   stats.MalesCount,
		"females_count": stats.FemalesCount,
		"total_sold":    stats.TotalSold,
		"total_revenue": stats.TotalRevenue,
	}

	SendSuccessResponse(w, data, "Estatísticas gerais recuperadas com sucesso", http.StatusOK)
}
