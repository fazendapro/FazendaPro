package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
)

type FarmHandler struct {
	service *service.FarmService
}

func NewFarmHandler(service *service.FarmService) *FarmHandler {
	return &FarmHandler{service: service}
}

// UpdateFarmRequest representa a requisição de atualização de fazenda
// @Description Dados para atualizar uma fazenda
type UpdateFarmRequest struct {
	Logo     string `json:"logo" example:"data:image/png;base64,..."` // Logo da fazenda em base64
	Language string `json:"language" example:"pt"`                    // Idioma da fazenda (pt, en, es)
}

// GetFarm obtém uma fazenda por ID
// @Summary      Obter fazenda por ID
// @Description  Retorna os dados de uma fazenda específica pelo ID
// @Tags         farm
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id query int true "ID da fazenda"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/farm [get]
func (h *FarmHandler) GetFarm(w http.ResponseWriter, r *http.Request) {
	farmIDStr := r.URL.Query().Get("id")
	if farmIDStr == "" {
		SendErrorResponse(w, "ID da fazenda é obrigatório", http.StatusBadRequest)
		return
	}

	farmID, err := strconv.ParseUint(farmIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da fazenda inválido", http.StatusBadRequest)
		return
	}

	farm, err := h.service.GetFarmByID(uint(farmID))
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if farm == nil {
		SendErrorResponse(w, "Fazenda não encontrada", http.StatusNotFound)
		return
	}

	if err := h.service.LoadCompanyData(farm); err != nil {
		SendErrorResponse(w, "Erro ao carregar dados da empresa: "+err.Error(), http.StatusInternalServerError)
		return
	}

	SendSuccessResponse(w, farm, "Fazenda encontrada com sucesso", http.StatusOK)
}

// UpdateFarm atualiza os dados de uma fazenda
// @Summary      Atualizar fazenda
// @Description  Atualiza os dados de uma fazenda existente
// @Tags         farm
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id query int true "ID da fazenda"
// @Param        request body UpdateFarmRequest true "Dados atualizados da fazenda"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/farm [put]
func (h *FarmHandler) UpdateFarm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	farmIDStr := r.URL.Query().Get("id")
	if farmIDStr == "" {
		SendErrorResponse(w, "ID da fazenda é obrigatório", http.StatusBadRequest)
		return
	}

	farmID, err := strconv.ParseUint(farmIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da fazenda inválido", http.StatusBadRequest)
		return
	}

	var req UpdateFarmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	farm := &models.Farm{
		ID:       uint(farmID),
		Logo:     req.Logo,
		Language: req.Language,
	}

	if err := h.service.UpdateFarm(farm); err != nil {
		SendErrorResponse(w, "Erro ao atualizar fazenda: "+err.Error(), http.StatusBadRequest)
		return
	}

	SendSuccessResponse(w, farm, "Fazenda atualizada com sucesso", http.StatusOK)
}
