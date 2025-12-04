package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/go-chi/chi/v5"
)

type VaccineHandler struct {
	service *service.VaccineService
}

func NewVaccineHandler(service *service.VaccineService) *VaccineHandler {
	return &VaccineHandler{service: service}
}

type VaccineData struct {
	ID           uint   `json:"id"`
	FarmID       uint   `json:"farm_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Manufacturer string `json:"manufacturer"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateVaccineRequest representa a requisição de criação de vacina
// @Description Dados necessários para criar uma nova vacina
type CreateVaccineRequest struct {
	FarmID       uint   `json:"farm_id" validate:"required" example:"1"`        // ID da fazenda
	Name         string `json:"name" validate:"required" example:"Vacina Aftosa"` // Nome da vacina
	Description  string `json:"description" example:"Vacina contra febre aftosa"` // Descrição da vacina
	Manufacturer string `json:"manufacturer" example:"Fabricante XYZ"`        // Fabricante
}

// VaccineResponse representa a resposta de vacina
// @Description Resposta com dados de uma vacina
type VaccineResponse struct {
	Success bool        `json:"success" example:"true"` // Indica sucesso
	Data    VaccineData `json:"data,omitempty"`         // Dados da vacina
	Message string      `json:"message,omitempty"`      // Mensagem de resposta
}

// VaccinesResponse representa a resposta com múltiplas vacinas
// @Description Resposta com lista de vacinas
type VaccinesResponse struct {
	Success bool          `json:"success" example:"true"` // Indica sucesso
	Data    []VaccineData `json:"data,omitempty"`         // Lista de vacinas
	Message string        `json:"message,omitempty"`       // Mensagem de resposta
}

// CreateVaccine cria uma nova vacina
// @Summary      Criar vacina
// @Description  Registra um novo tipo de vacina na fazenda
// @Tags         vaccines
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateVaccineRequest true "Dados da vacina"
// @Success      201  {object}  VaccineResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/vaccines [post]
func (h *VaccineHandler) CreateVaccine(w http.ResponseWriter, r *http.Request) {
	var req CreateVaccineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	vaccine := &models.Vaccine{
		FarmID:       req.FarmID,
		Name:         req.Name,
		Description:  req.Description,
		Manufacturer: req.Manufacturer,
	}

	if err := h.service.CreateVaccine(vaccine); err != nil {
		SendErrorResponse(w, "Erro ao criar vacina: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdVaccine, err := h.service.GetVaccineByID(vaccine.ID)
	if err != nil {
		SendErrorResponse(w, "Erro ao recuperar vacina criada: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := VaccineResponse{
		Success: true,
		Data:    h.mapToVaccineData(createdVaccine),
		Message: "Vacina criada com sucesso",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetVaccinesByFarmID obtém lista de vacinas da fazenda
// @Summary      Obter vacinas da fazenda
// @Description  Retorna lista de vacinas cadastradas na fazenda
// @Tags         vaccines
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        farmId path int true "ID da fazenda"
// @Success      200  {object}  VaccinesResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/vaccines/farm/{farmId} [get]
func (h *VaccineHandler) GetVaccinesByFarmID(w http.ResponseWriter, r *http.Request) {
	farmIDStr := chi.URLParam(r, "farmId")
	farmID, err := strconv.ParseUint(farmIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da fazenda inválido", http.StatusBadRequest)
		return
	}

	vaccines, err := h.service.GetVaccinesByFarmID(uint(farmID))
	if err != nil {
		SendErrorResponse(w, "Erro ao buscar vacinas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	vaccineData := make([]VaccineData, len(vaccines))
	for i, v := range vaccines {
		vaccineData[i] = h.mapToVaccineData(&v)
	}

	response := VaccinesResponse{
		Success: true,
		Data:    vaccineData,
		Message: "Vacinas recuperadas com sucesso",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(response)
}

// GetVaccineByID obtém uma vacina por ID
// @Summary      Obter vacina por ID
// @Description  Retorna os dados de uma vacina específica pelo ID
// @Tags         vaccines
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID da vacina"
// @Success      200  {object}  VaccineResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/vaccines/{id} [get]
func (h *VaccineHandler) GetVaccineByID(w http.ResponseWriter, r *http.Request) {
	vaccineIDStr := chi.URLParam(r, "id")
	vaccineID, err := strconv.ParseUint(vaccineIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da vacina inválido", http.StatusBadRequest)
		return
	}

	vaccine, err := h.service.GetVaccineByID(uint(vaccineID))
	if err != nil {
		SendErrorResponse(w, "Erro ao buscar vacina: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if vaccine == nil {
		SendErrorResponse(w, "Vacina não encontrada", http.StatusNotFound)
		return
	}

	response := VaccineResponse{
		Success: true,
		Data:    h.mapToVaccineData(vaccine),
		Message: "Vacina encontrada com sucesso",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(response)
}

// UpdateVaccine atualiza uma vacina
// @Summary      Atualizar vacina
// @Description  Atualiza os dados de uma vacina existente
// @Tags         vaccines
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID da vacina"
// @Param        request body CreateVaccineRequest true "Dados atualizados"
// @Success      200  {object}  VaccineResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/vaccines/{id} [put]
func (h *VaccineHandler) UpdateVaccine(w http.ResponseWriter, r *http.Request) {
	vaccineIDStr := chi.URLParam(r, "id")
	vaccineID, err := strconv.ParseUint(vaccineIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da vacina inválido", http.StatusBadRequest)
		return
	}

	var req CreateVaccineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	vaccine := &models.Vaccine{
		ID:           uint(vaccineID),
		FarmID:       req.FarmID,
		Name:         req.Name,
		Description:  req.Description,
		Manufacturer: req.Manufacturer,
	}

	if err := h.service.UpdateVaccine(vaccine); err != nil {
		SendErrorResponse(w, "Erro ao atualizar vacina: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedVaccine, err := h.service.GetVaccineByID(vaccine.ID)
	if err != nil {
		SendErrorResponse(w, "Erro ao recuperar vacina atualizada: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := VaccineResponse{
		Success: true,
		Data:    h.mapToVaccineData(updatedVaccine),
		Message: "Vacina atualizada com sucesso",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteVaccine deleta uma vacina
// @Summary      Deletar vacina
// @Description  Remove uma vacina do sistema
// @Tags         vaccines
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID da vacina"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/vaccines/{id} [delete]
func (h *VaccineHandler) DeleteVaccine(w http.ResponseWriter, r *http.Request) {
	vaccineIDStr := chi.URLParam(r, "id")
	vaccineID, err := strconv.ParseUint(vaccineIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da vacina inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteVaccine(uint(vaccineID)); err != nil {
		SendErrorResponse(w, "Erro ao deletar vacina: "+err.Error(), http.StatusBadRequest)
		return
	}

	SendSuccessResponse(w, nil, "Vacina deletada com sucesso", http.StatusOK)
}

func (h *VaccineHandler) mapToVaccineData(v *models.Vaccine) VaccineData {
	return VaccineData{
		ID:           v.ID,
		FarmID:      v.FarmID,
		Name:         v.Name,
		Description:  v.Description,
		Manufacturer: v.Manufacturer,
		CreatedAt:    v.CreatedAt.Format(DateFormatISO8601),
		UpdatedAt:    v.UpdatedAt.Format(DateFormatISO8601),
	}
}

