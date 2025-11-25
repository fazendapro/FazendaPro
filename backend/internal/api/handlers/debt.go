package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/go-chi/chi/v5"
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
		http.Error(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
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
		CreatedAt: debt.CreatedAt.Format(DateFormatISO8601),
		UpdatedAt: debt.UpdatedAt.Format(DateFormatISO8601),
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

type queryParams struct {
	page  int
	limit int
	year  *int
	month *int
}

func parseQueryParams(r *http.Request) queryParams {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	params := queryParams{
		page:  1,
		limit: 10,
	}

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			params.page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			params.limit = l
		}
	}

	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			params.year = &y
		}
	}

	if monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			params.month = &m
		}
	}

	return params
}

func (h *DebtHandler) GetDebts(w http.ResponseWriter, r *http.Request) {
	params := parseQueryParams(r)

	debts, total, err := h.service.GetDebtsWithPagination(params.page, params.limit, params.year, params.month)
	if err != nil {
		http.Error(w, "Erro ao buscar dívidas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	debtResponses := make([]DebtResponse, len(debts))
	for i, debt := range debts {
		debtResponses[i] = DebtResponse{
			ID:        debt.ID,
			Person:    debt.Person,
			Value:     debt.Value,
			CreatedAt: debt.CreatedAt.Format(DateFormatISO8601),
			UpdatedAt: debt.UpdatedAt.Format(DateFormatISO8601),
		}
	}

	response := DebtListResponse{
		Debts: debtResponses,
		Total: total,
		Page:  params.page,
		Limit: params.limit,
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(response)
}

func (h *DebtHandler) DeleteDebt(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
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

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dívida deletada com sucesso"})
}

func (h *DebtHandler) GetTotalByPerson(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	if yearStr == "" || monthStr == "" {
		http.Error(w, "Parâmetros 'year' e 'month' são obrigatórios", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Ano deve ser um número válido", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		http.Error(w, "Mês deve ser um número válido", http.StatusBadRequest)
		return
	}

	totals, err := h.service.GetTotalByPersonInMonth(year, month)
	if err != nil {
		http.Error(w, "Erro ao calcular total por pessoa: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"year":   year,
		"month":  month,
		"totals": totals,
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(response)
}
