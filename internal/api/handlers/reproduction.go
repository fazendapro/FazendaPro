package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
)

type ReproductionHandler struct {
	service *service.ReproductionService
}

func NewReproductionHandler(service *service.ReproductionService) *ReproductionHandler {
	return &ReproductionHandler{service: service}
}

// ReproductionData representa os dados de reprodução
type ReproductionData struct {
	ID                     uint    `json:"id"`
	AnimalID               uint    `json:"animal_id"`
	AnimalName             string  `json:"animal_name,omitempty"`
	EarTag                 int     `json:"ear_tag,omitempty"`
	CurrentPhase           int     `json:"current_phase"`
	InseminationDate       *string `json:"insemination_date,omitempty"`
	InseminationType       string  `json:"insemination_type,omitempty"`
	PregnancyDate          *string `json:"pregnancy_date,omitempty"`
	ExpectedBirthDate      *string `json:"expected_birth_date,omitempty"`
	ActualBirthDate        *string `json:"actual_birth_date,omitempty"`
	LactationStartDate     *string `json:"lactation_start_date,omitempty"`
	LactationEndDate       *string `json:"lactation_end_date,omitempty"`
	DryPeriodStartDate     *string `json:"dry_period_start_date,omitempty"`
	VeterinaryConfirmation bool    `json:"veterinary_confirmation"`
	Observations           string  `json:"observations,omitempty"`
}

type CreateReproductionRequest struct {
	ReproductionData
}

type UpdateReproductionPhaseRequest struct {
	AnimalID       uint                   `json:"animal_id"`
	NewPhase       int                    `json:"new_phase"`
	AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
}

type ReproductionResponse struct {
	ReproductionData
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func reproductionDataToModel(data ReproductionData) models.Reproduction {
	reproduction := models.Reproduction{
		ID:                     data.ID,
		AnimalID:               data.AnimalID,
		CurrentPhase:           models.ReproductionPhase(data.CurrentPhase),
		InseminationType:       data.InseminationType,
		VeterinaryConfirmation: data.VeterinaryConfirmation,
		Observations:           data.Observations,
	}

	// Parse dates
	if data.InseminationDate != nil && *data.InseminationDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", *data.InseminationDate); err == nil {
			reproduction.InseminationDate = &parsedDate
		}
	}

	if data.PregnancyDate != nil && *data.PregnancyDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", *data.PregnancyDate); err == nil {
			reproduction.PregnancyDate = &parsedDate
		}
	}

	if data.ExpectedBirthDate != nil && *data.ExpectedBirthDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", *data.ExpectedBirthDate); err == nil {
			reproduction.ExpectedBirthDate = &parsedDate
		}
	}

	if data.ActualBirthDate != nil && *data.ActualBirthDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", *data.ActualBirthDate); err == nil {
			reproduction.ActualBirthDate = &parsedDate
		}
	}

	if data.LactationStartDate != nil && *data.LactationStartDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", *data.LactationStartDate); err == nil {
			reproduction.LactationStartDate = &parsedDate
		}
	}

	if data.LactationEndDate != nil && *data.LactationEndDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", *data.LactationEndDate); err == nil {
			reproduction.LactationEndDate = &parsedDate
		}
	}

	if data.DryPeriodStartDate != nil && *data.DryPeriodStartDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", *data.DryPeriodStartDate); err == nil {
			reproduction.DryPeriodStartDate = &parsedDate
		}
	}

	return reproduction
}

func modelToReproductionResponse(reproduction *models.Reproduction) ReproductionResponse {
	response := ReproductionResponse{
		ReproductionData: ReproductionData{
			ID:                     reproduction.ID,
			AnimalID:               reproduction.AnimalID,
			AnimalName:             reproduction.Animal.AnimalName,
			EarTag:                 reproduction.Animal.EarTagNumberLocal,
			CurrentPhase:           int(reproduction.CurrentPhase),
			InseminationType:       reproduction.InseminationType,
			VeterinaryConfirmation: reproduction.VeterinaryConfirmation,
			Observations:           reproduction.Observations,
		},
		CreatedAt: reproduction.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: reproduction.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Format dates
	if reproduction.InseminationDate != nil {
		dateStr := reproduction.InseminationDate.Format("2006-01-02")
		response.InseminationDate = &dateStr
	}

	if reproduction.PregnancyDate != nil {
		dateStr := reproduction.PregnancyDate.Format("2006-01-02")
		response.PregnancyDate = &dateStr
	}

	if reproduction.ExpectedBirthDate != nil {
		dateStr := reproduction.ExpectedBirthDate.Format("2006-01-02")
		response.ExpectedBirthDate = &dateStr
	}

	if reproduction.ActualBirthDate != nil {
		dateStr := reproduction.ActualBirthDate.Format("2006-01-02")
		response.ActualBirthDate = &dateStr
	}

	if reproduction.LactationStartDate != nil {
		dateStr := reproduction.LactationStartDate.Format("2006-01-02")
		response.LactationStartDate = &dateStr
	}

	if reproduction.LactationEndDate != nil {
		dateStr := reproduction.LactationEndDate.Format("2006-01-02")
		response.LactationEndDate = &dateStr
	}

	if reproduction.DryPeriodStartDate != nil {
		dateStr := reproduction.DryPeriodStartDate.Format("2006-01-02")
		response.DryPeriodStartDate = &dateStr
	}

	return response
}

func (h *ReproductionHandler) CreateReproduction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req CreateReproductionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	reproduction := reproductionDataToModel(req.ReproductionData)

	if err := h.service.CreateReproduction(&reproduction); err != nil {
		SendErrorResponse(w, "Erro ao criar registro de reprodução: "+err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"id": reproduction.ID,
	}
	SendSuccessResponse(w, data, "Registro de reprodução criado com sucesso", http.StatusCreated)
}

func (h *ReproductionHandler) GetReproduction(w http.ResponseWriter, r *http.Request) {
	reproductionID := r.URL.Query().Get("id")
	if reproductionID == "" {
		SendErrorResponse(w, "ID do registro de reprodução é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(reproductionID, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID do registro de reprodução inválido", http.StatusBadRequest)
		return
	}

	reproduction, err := h.service.GetReproductionByID(uint(id))
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if reproduction == nil {
		SendErrorResponse(w, "Registro de reprodução não encontrado", http.StatusNotFound)
		return
	}

	response := modelToReproductionResponse(reproduction)
	SendSuccessResponse(w, response, "Registro de reprodução encontrado com sucesso", http.StatusOK)
}

func (h *ReproductionHandler) GetReproductionByAnimal(w http.ResponseWriter, r *http.Request) {
	animalID := r.URL.Query().Get("animalId")
	if animalID == "" {
		SendErrorResponse(w, "ID do animal é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(animalID, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID do animal inválido", http.StatusBadRequest)
		return
	}

	reproduction, err := h.service.GetReproductionByAnimalID(uint(id))
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if reproduction == nil {
		SendErrorResponse(w, "Registro de reprodução não encontrado para este animal", http.StatusNotFound)
		return
	}

	response := modelToReproductionResponse(reproduction)
	SendSuccessResponse(w, response, "Registro de reprodução encontrado com sucesso", http.StatusOK)
}

func (h *ReproductionHandler) GetReproductionsByFarm(w http.ResponseWriter, r *http.Request) {
	farmID := r.URL.Query().Get("farmId")
	if farmID == "" {
		SendErrorResponse(w, "ID da fazenda é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(farmID, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da fazenda inválido", http.StatusBadRequest)
		return
	}

	reproductions, err := h.service.GetReproductionsByFarmID(uint(id))
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []ReproductionResponse
	for _, reproduction := range reproductions {
		responses = append(responses, modelToReproductionResponse(&reproduction))
	}

	SendSuccessResponse(w, responses, fmt.Sprintf("Registros de reprodução encontrados com sucesso (%d registros)", len(reproductions)), http.StatusOK)
}

func (h *ReproductionHandler) GetReproductionsByPhase(w http.ResponseWriter, r *http.Request) {
	phaseStr := r.URL.Query().Get("phase")
	if phaseStr == "" {
		SendErrorResponse(w, "Fase é obrigatória", http.StatusBadRequest)
		return
	}

	phase, err := strconv.Atoi(phaseStr)
	if err != nil {
		SendErrorResponse(w, "Fase inválida", http.StatusBadRequest)
		return
	}

	reproductions, err := h.service.GetReproductionsByPhase(models.ReproductionPhase(phase))
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []ReproductionResponse
	for _, reproduction := range reproductions {
		responses = append(responses, modelToReproductionResponse(&reproduction))
	}

	SendSuccessResponse(w, responses, fmt.Sprintf("Registros de reprodução encontrados com sucesso (%d registros)", len(reproductions)), http.StatusOK)
}

func (h *ReproductionHandler) UpdateReproduction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req CreateReproductionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	reproduction := reproductionDataToModel(req.ReproductionData)

	if err := h.service.UpdateReproduction(&reproduction); err != nil {
		SendErrorResponse(w, "Erro ao atualizar registro de reprodução: "+err.Error(), http.StatusBadRequest)
		return
	}

	SendSuccessResponse(w, nil, "Registro de reprodução atualizado com sucesso", http.StatusOK)
}

func (h *ReproductionHandler) UpdateReproductionPhase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateReproductionPhaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Converter additionalData para o formato correto
	additionalData := make(map[string]interface{})
	for key, value := range req.AdditionalData {
		additionalData[key] = value
	}

	if err := h.service.UpdateReproductionPhase(req.AnimalID, models.ReproductionPhase(req.NewPhase), additionalData); err != nil {
		SendErrorResponse(w, "Erro ao atualizar fase de reprodução: "+err.Error(), http.StatusBadRequest)
		return
	}

	SendSuccessResponse(w, nil, "Fase de reprodução atualizada com sucesso", http.StatusOK)
}

func (h *ReproductionHandler) DeleteReproduction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	reproductionID := r.URL.Query().Get("id")
	if reproductionID == "" {
		SendErrorResponse(w, "ID do registro de reprodução é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(reproductionID, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID do registro de reprodução inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteReproduction(uint(id)); err != nil {
		SendErrorResponse(w, "Erro ao deletar registro de reprodução: "+err.Error(), http.StatusBadRequest)
		return
	}

	SendSuccessResponse(w, nil, "Registro de reprodução deletado com sucesso", http.StatusOK)
}
