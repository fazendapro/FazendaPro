package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/go-chi/chi/v5"
)

type WeightHandler struct {
	service *service.WeightService
}

func NewWeightHandler(service *service.WeightService) *WeightHandler {
	return &WeightHandler{service: service}
}

type WeightData struct {
	ID           uint    `json:"id"`
	AnimalID     uint    `json:"animal_id"`
	AnimalName   string  `json:"animal_name,omitempty"`
	EarTag       int     `json:"ear_tag,omitempty"`
	Date         string  `json:"date"`
	AnimalWeight float64 `json:"animal_weight"`
}

// CreateOrUpdateWeightRequest representa a requisição de criação ou atualização de peso
// @Description Dados necessários para criar ou atualizar um registro de peso
type CreateOrUpdateWeightRequest struct {
	WeightData // Dados do peso
}

// UpdateWeightRequest representa a requisição de atualização de peso
// @Description Dados para atualizar um registro de peso existente
type UpdateWeightRequest struct {
	WeightData // Dados do peso
}

// WeightResponse representa a resposta de peso
// @Description Resposta com dados completos de um peso
type WeightResponse struct {
	WeightData        // Dados do peso
	CreatedAt  string `json:"created_at" example:"2024-01-15T10:30:00Z"` // Data de criação
	UpdatedAt  string `json:"updated_at" example:"2024-01-15T10:30:00Z"` // Data de atualização
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse(DateFormatISO, dateStr)
}

func weightDataToModel(data WeightData) (models.Weight, error) {
	date, err := parseDate(data.Date)
	if err != nil {
		return models.Weight{}, fmt.Errorf("data inválida: %v", err)
	}

	weight := models.Weight{
		ID:           data.ID,
		AnimalID:     data.AnimalID,
		Date:         date,
		AnimalWeight: data.AnimalWeight,
	}

	return weight, nil
}

func modelToWeightResponse(weight *models.Weight) WeightResponse {
	response := WeightResponse{
		WeightData: WeightData{
			ID:           weight.ID,
			AnimalID:     weight.AnimalID,
			AnimalName:   weight.Animal.AnimalName,
			EarTag:       weight.Animal.EarTagNumberLocal,
			Date:         weight.Date.Format(DateFormatISO),
			AnimalWeight: weight.AnimalWeight,
		},
		CreatedAt: weight.CreatedAt.Format(DateFormatDateTime),
		UpdatedAt: weight.UpdatedAt.Format(DateFormatDateTime),
	}

	return response
}

// CreateOrUpdateWeight cria ou atualiza um registro de peso
// @Summary      Criar ou atualizar peso
// @Description  Cria um novo registro de peso ou atualiza o existente para um animal
// @Tags         weights
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateOrUpdateWeightRequest true "Dados do peso"
// @Success      201  {object}  WeightResponse
// @Success      200  {object}  WeightResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/weights [post]
func (h *WeightHandler) CreateOrUpdateWeight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req CreateOrUpdateWeightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	weight, err := weightDataToModel(req.WeightData)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingWeight, _ := h.service.GetWeightByAnimalID(weight.AnimalID)
	isUpdate := existingWeight != nil

	if err := h.service.CreateOrUpdateWeight(&weight); err != nil {
		SendErrorResponse(w, "Erro ao criar ou atualizar registro de peso: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedWeight, err := h.service.GetWeightByAnimalID(weight.AnimalID)
	if err != nil || updatedWeight == nil {
		SendErrorResponse(w, "Erro ao recuperar registro de peso", http.StatusInternalServerError)
		return
	}

	response := modelToWeightResponse(updatedWeight)
	message := "Registro de peso criado com sucesso"
	if isUpdate {
		message = "Registro de peso atualizado com sucesso"
		SendSuccessResponse(w, response, message, http.StatusOK)
	} else {
		SendSuccessResponse(w, response, message, http.StatusCreated)
	}
}

// GetWeightByAnimal obtém o peso de um animal
// @Summary      Obter peso por animal
// @Description  Retorna o peso mais recente de um animal específico
// @Tags         weights
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        animalId path int true "ID do animal"
// @Success      200  {object}  WeightResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/weights/animal/{animalId} [get]
func (h *WeightHandler) GetWeightByAnimal(w http.ResponseWriter, r *http.Request) {
	animalIDStr := chi.URLParam(r, "animalId")
	if animalIDStr == "" {
		SendErrorResponse(w, ErrAnimalIDRequired, http.StatusBadRequest)
		return
	}

	animalID, err := strconv.ParseUint(animalIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, ErrInvalidAnimalID, http.StatusBadRequest)
		return
	}

	weight, err := h.service.GetWeightByAnimalID(uint(animalID))
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if weight == nil {
		SendErrorResponse(w, "Peso não encontrado para este animal", http.StatusNotFound)
		return
	}

	response := modelToWeightResponse(weight)
	SendSuccessResponse(w, response, "Peso encontrado com sucesso", http.StatusOK)
}

// GetWeightsByFarm obtém todos os pesos de uma fazenda
// @Summary      Obter pesos por fazenda
// @Description  Retorna todos os registros de peso de uma fazenda
// @Tags         weights
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        farmId path int true "ID da fazenda"
// @Success      200  {object}  []WeightResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/weights/farm/{farmId} [get]
func (h *WeightHandler) GetWeightsByFarm(w http.ResponseWriter, r *http.Request) {
	farmIDStr := chi.URLParam(r, "farmId")
	if farmIDStr == "" {
		SendErrorResponse(w, "ID da fazenda é obrigatório", http.StatusBadRequest)
		return
	}

	farmID, err := strconv.ParseUint(farmIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, ErrInvalidFarmID, http.StatusBadRequest)
		return
	}

	weights, err := h.service.GetWeightsByFarmID(uint(farmID))
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []WeightResponse
	for _, weight := range weights {
		responses = append(responses, modelToWeightResponse(&weight))
	}

	SendSuccessResponse(w, responses, fmt.Sprintf("Registros de peso encontrados com sucesso (%d registros)", len(weights)), http.StatusOK)
}

// UpdateWeight atualiza um registro de peso
// @Summary      Atualizar peso
// @Description  Atualiza um registro de peso existente
// @Tags         weights
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateWeightRequest true "Dados do peso"
// @Success      200  {object}  WeightResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/weights [put]
func (h *WeightHandler) UpdateWeight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req UpdateWeightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	weight, err := weightDataToModel(req.WeightData)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateWeight(&weight); err != nil {
		SendErrorResponse(w, "Erro ao atualizar registro de peso: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedWeight, err := h.service.GetWeightByID(weight.ID)
	if err != nil || updatedWeight == nil {
		SendErrorResponse(w, "Erro ao recuperar registro de peso atualizado", http.StatusInternalServerError)
		return
	}

	response := modelToWeightResponse(updatedWeight)
	SendSuccessResponse(w, response, "Registro de peso atualizado com sucesso", http.StatusOK)
}
