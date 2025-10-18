package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
)

type DebtHandler struct {
	service *service.DebtService
}

func NewDebtHandler(service *service.DebtService) *DebtHandler {
	return &DebtHandler{service: service}
}

type CreateDebtRequest struct {
	Person string  `json:"person"`
	Value  float64 `json:"value"`
}

type DebtResponse struct {
	ID        uint    `json:"id"`
	Person    string  `json:"person"`
	Value     float64 `json:"value"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type DebtListResponse struct {
	Debts []DebtResponse `json:"debts"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

func (h *DebtHandler) CreateDebt(w http.ResponseWriter, r *http.Request) {
	var req CreateDebtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	debt := &models.Debt{
		Person: req.Person,
		Value:  req.Value,
	}

	if err := h.service.CreateDebt(debt); err != nil {
		http.Error(w, "Erro ao criar dívida: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := DebtResponse{
		ID:        debt.ID,
		Person:    debt.Person,
		Value:     debt.Value,
		CreatedAt: debt.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt: debt.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *DebtHandler) GetDebts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	// Default values
	page := 1
	limit := 10
	var year, month *int

	// Parse page
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Parse year
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = &y
		}
	}

	// Parse month
	if monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			month = &m
		}
	}

	debts, total, err := h.service.GetDebtsWithPagination(page, limit, year, month)
	if err != nil {
		http.Error(w, "Erro ao buscar dívidas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response format
	debtResponses := make([]DebtResponse, len(debts))
	for i, debt := range debts {
		debtResponses[i] = DebtResponse{
			ID:        debt.ID,
			Person:    debt.Person,
			Value:     debt.Value,
			CreatedAt: debt.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
			UpdatedAt: debt.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
		}
	}

	response := DebtListResponse{
		Debts: debtResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DebtHandler) DeleteDebt(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteDebt(uint(id)); err != nil {
		http.Error(w, "Erro ao deletar dívida: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dívida deletada com sucesso"})
}
