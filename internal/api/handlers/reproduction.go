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

func parseOptionalDate(dateStr *string) *time.Time {
	if dateStr == nil || *dateStr == "" {
		return nil
	}
	parsedDate, err := time.Parse(DateFormatISO, *dateStr)
	if err != nil {
		return nil
	}
	return &parsedDate
}

func reproductionDataToModel(data ReproductionData) models.Reproduction {
	reproduction := models.Reproduction{
		ID:                     data.ID,
		AnimalID:               data.AnimalID,
		CurrentPhase:           models.ReproductionPhase(data.CurrentPhase),
		InseminationType:       data.InseminationType,
		VeterinaryConfirmation: data.VeterinaryConfirmation,
		Observations:           data.Observations,
		InseminationDate:       parseOptionalDate(data.InseminationDate),
		PregnancyDate:          parseOptionalDate(data.PregnancyDate),
		ExpectedBirthDate:      parseOptionalDate(data.ExpectedBirthDate),
		ActualBirthDate:        parseOptionalDate(data.ActualBirthDate),
		LactationStartDate:     parseOptionalDate(data.LactationStartDate),
		LactationEndDate:       parseOptionalDate(data.LactationEndDate),
		DryPeriodStartDate:     parseOptionalDate(data.DryPeriodStartDate),
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
		CreatedAt: reproduction.CreatedAt.Format(DateFormatDateTime),
		UpdatedAt: reproduction.UpdatedAt.Format(DateFormatDateTime),
	}

	if reproduction.InseminationDate != nil {
		dateStr := reproduction.InseminationDate.Format(DateFormatISO)
		response.InseminationDate = &dateStr
	}

	if reproduction.PregnancyDate != nil {
		dateStr := reproduction.PregnancyDate.Format(DateFormatISO)
		response.PregnancyDate = &dateStr
	}

	if reproduction.ExpectedBirthDate != nil {
		dateStr := reproduction.ExpectedBirthDate.Format(DateFormatISO)
		response.ExpectedBirthDate = &dateStr
	}

	if reproduction.ActualBirthDate != nil {
		dateStr := reproduction.ActualBirthDate.Format(DateFormatISO)
		response.ActualBirthDate = &dateStr
	}

	if reproduction.LactationStartDate != nil {
		dateStr := reproduction.LactationStartDate.Format(DateFormatISO)
		response.LactationStartDate = &dateStr
	}

	if reproduction.LactationEndDate != nil {
		dateStr := reproduction.LactationEndDate.Format(DateFormatISO)
		response.LactationEndDate = &dateStr
	}

	if reproduction.DryPeriodStartDate != nil {
		dateStr := reproduction.DryPeriodStartDate.Format(DateFormatISO)
		response.DryPeriodStartDate = &dateStr
	}

	return response
}

func (h *ReproductionHandler) CreateReproduction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req CreateReproductionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
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
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req CreateReproductionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
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
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req UpdateReproductionPhaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

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
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
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

type NextToCalveResponse struct {
	ID                uint   `json:"id"`
	AnimalName        string `json:"animal_name"`
	EarTagNumberLocal int    `json:"ear_tag_number_local"`
	Photo             string `json:"photo"`
	PregnancyDate     string `json:"pregnancy_date"`
	ExpectedBirthDate string `json:"expected_birth_date"`
	DaysUntilBirth    int    `json:"days_until_birth"`
	Status            string `json:"status"`
}

func calculateBirthStatus(daysUntilBirth int) string {
	if daysUntilBirth <= 30 {
		return "Alto"
	}
	if daysUntilBirth <= 60 {
		return "Médio"
	}
	return "Baixo"
}

func sortByDaysUntilBirth(responses []NextToCalveResponse) {
	for i := 0; i < len(responses); i++ {
		for j := i + 1; j < len(responses); j++ {
			if responses[i].DaysUntilBirth > responses[j].DaysUntilBirth {
				responses[i], responses[j] = responses[j], responses[i]
			}
		}
	}
}

func buildNextToCalveResponses(reproductions []models.Reproduction, farmID uint, now time.Time) []NextToCalveResponse {
	var responses []NextToCalveResponse

	for _, reproduction := range reproductions {
		if reproduction.Animal.FarmID != farmID || reproduction.PregnancyDate == nil {
			continue
		}

		expectedBirth := reproduction.PregnancyDate.AddDate(0, 0, 283)
		daysUntilBirth := int(expectedBirth.Sub(now).Hours() / 24)

		response := NextToCalveResponse{
			ID:                reproduction.Animal.ID,
			AnimalName:        reproduction.Animal.AnimalName,
			EarTagNumberLocal: reproduction.Animal.EarTagNumberLocal,
			Photo:             reproduction.Animal.Photo,
			PregnancyDate:     reproduction.PregnancyDate.Format(DateFormatISO),
			ExpectedBirthDate: expectedBirth.Format(DateFormatISO),
			DaysUntilBirth:    daysUntilBirth,
			Status:            calculateBirthStatus(daysUntilBirth),
		}

		responses = append(responses, response)
	}

	return responses
}

func (h *ReproductionHandler) GetNextToCalve(w http.ResponseWriter, r *http.Request) {
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

	reproductions, err := h.service.GetReproductionsByPhase(models.PhasePrenhas)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := buildNextToCalveResponses(reproductions, uint(id), time.Now())
	sortByDaysUntilBirth(responses)

	SendSuccessResponse(w, responses, fmt.Sprintf("Próximas vacas a parir encontradas com sucesso (%d registros)", len(responses)), http.StatusOK)
}
