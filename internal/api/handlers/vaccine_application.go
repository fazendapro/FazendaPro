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

type VaccineApplicationHandler struct {
	service *service.VaccineApplicationService
}

func NewVaccineApplicationHandler(service *service.VaccineApplicationService) *VaccineApplicationHandler {
	return &VaccineApplicationHandler{service: service}
}

type VaccineApplicationData struct {
	ID              uint       `json:"id"`
	AnimalID        uint       `json:"animal_id"`
	Animal          AnimalData `json:"animal"`
	VaccineID       uint       `json:"vaccine_id"`
	Vaccine         VaccineData `json:"vaccine"`
	ApplicationDate string     `json:"application_date"`
	BatchNumber     string     `json:"batch_number"`
	Veterinarian    string     `json:"veterinarian"`
	Observations    string     `json:"observations"`
	CreatedAt       string     `json:"created_at"`
	UpdatedAt       string     `json:"updated_at"`
}

// CreateVaccineApplicationRequest representa a requisição de criação de aplicação de vacina
// @Description Dados necessários para criar uma nova aplicação de vacina
type CreateVaccineApplicationRequest struct {
	AnimalID        uint   `json:"animal_id" validate:"required" example:"1"`              // ID do animal
	VaccineID       uint   `json:"vaccine_id" validate:"required" example:"1"`               // ID da vacina
	ApplicationDate string `json:"application_date" validate:"required" example:"2024-01-15"` // Data da aplicação (YYYY-MM-DD)
	BatchNumber     string `json:"batch_number" example:"LOTE123"`                            // Número do lote
	Veterinarian    string `json:"veterinarian" example:"Dr. João Silva"`                    // Nome do veterinário
	Observations    string `json:"observations" example:"Aplicação realizada com sucesso"`     // Observações
}

// VaccineApplicationResponse representa a resposta de aplicação de vacina
// @Description Resposta com dados de uma aplicação de vacina
type VaccineApplicationResponse struct {
	Success bool                   `json:"success" example:"true"` // Indica sucesso
	Data    VaccineApplicationData `json:"data,omitempty"`       // Dados da aplicação
	Message string                 `json:"message,omitempty"`     // Mensagem de resposta
}

// VaccineApplicationsResponse representa a resposta com múltiplas aplicações
// @Description Resposta com lista de aplicações de vacinas
type VaccineApplicationsResponse struct {
	Success bool                     `json:"success" example:"true"` // Indica sucesso
	Data    []VaccineApplicationData `json:"data,omitempty"`        // Lista de aplicações
	Message string                   `json:"message,omitempty"`      // Mensagem de resposta
}

// CreateVaccineApplication cria uma nova aplicação de vacina
// @Summary      Criar aplicação de vacina
// @Description  Registra uma nova aplicação de vacina em um animal
// @Tags         vaccine-applications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateVaccineApplicationRequest true "Dados da aplicação"
// @Success      201  {object}  VaccineApplicationResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/vaccine-applications [post]
func (h *VaccineApplicationHandler) CreateVaccineApplication(w http.ResponseWriter, r *http.Request) {
	var req CreateVaccineApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse(DateFormatISO, req.ApplicationDate)
	if err != nil {
		SendErrorResponse(w, "Formato de data inválido. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	vaccineApplication := &models.VaccineApplication{
		AnimalID:        req.AnimalID,
		VaccineID:       req.VaccineID,
		ApplicationDate: date,
		BatchNumber:     req.BatchNumber,
		Veterinarian:    req.Veterinarian,
		Observations:    req.Observations,
	}

	if err := h.service.CreateApplication(vaccineApplication); err != nil {
		SendErrorResponse(w, "Erro ao criar aplicação de vacina: "+err.Error(), http.StatusBadRequest)
		return
	}

	createdApplication, err := h.service.GetApplicationByID(vaccineApplication.ID)
	if err != nil {
		SendErrorResponse(w, "Erro ao recuperar aplicação criada: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := VaccineApplicationResponse{
		Success: true,
		Data:    h.mapToVaccineApplicationData(createdApplication),
		Message: "Aplicação de vacina criada com sucesso",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetVaccineApplicationsByFarmID obtém lista de aplicações de vacinas da fazenda
// @Summary      Obter aplicações de vacinas da fazenda
// @Description  Retorna lista de aplicações de vacinas com filtros opcionais de data
// @Tags         vaccine-applications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        farmId path int true "ID da fazenda"
// @Param        start_date query string false "Data inicial (YYYY-MM-DD)"
// @Param        end_date query string false "Data final (YYYY-MM-DD)"
// @Success      200  {object}  VaccineApplicationsResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/vaccine-applications/farm/{farmId} [get]
func (h *VaccineApplicationHandler) GetVaccineApplicationsByFarmID(w http.ResponseWriter, r *http.Request) {
	farmIDStr := chi.URLParam(r, "farmId")
	farmID, err := strconv.ParseUint(farmIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da fazenda inválido", http.StatusBadRequest)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		if parsed, err := time.Parse(DateFormatISO, startDateStr); err == nil {
			startDate = &parsed
		}
	}

	if endDateStr != "" {
		if parsed, err := time.Parse(DateFormatISO, endDateStr); err == nil {
			endDate = &parsed
		}
	}

	var vaccineApplications []models.VaccineApplication
	if startDate != nil || endDate != nil {
		vaccineApplications, err = h.service.GetApplicationsByFarmIDWithDateRange(uint(farmID), startDate, endDate)
	} else {
		vaccineApplications, err = h.service.GetApplicationsByFarmID(uint(farmID))
	}

	if err != nil {
		SendErrorResponse(w, "Erro ao buscar aplicações de vacinas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	vaccineApplicationData := make([]VaccineApplicationData, len(vaccineApplications))
	for i, va := range vaccineApplications {
		vaccineApplicationData[i] = h.mapToVaccineApplicationData(&va)
	}

	response := VaccineApplicationsResponse{
		Success: true,
		Data:    vaccineApplicationData,
		Message: "Aplicações de vacinas recuperadas com sucesso",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(response)
}

// GetVaccineApplicationsByAnimalID obtém lista de aplicações de vacinas de um animal
// @Summary      Obter aplicações de vacinas de um animal
// @Description  Retorna lista de aplicações de vacinas de um animal específico
// @Tags         vaccine-applications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        animalId path int true "ID do animal"
// @Success      200  {object}  VaccineApplicationsResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/vaccine-applications/animal/{animalId} [get]
func (h *VaccineApplicationHandler) GetVaccineApplicationsByAnimalID(w http.ResponseWriter, r *http.Request) {
	animalIDStr := chi.URLParam(r, "animalId")
	animalID, err := strconv.ParseUint(animalIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID do animal inválido", http.StatusBadRequest)
		return
	}

	vaccineApplications, err := h.service.GetApplicationsByAnimalID(uint(animalID))
	if err != nil {
		SendErrorResponse(w, "Erro ao buscar aplicações de vacinas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	vaccineApplicationData := make([]VaccineApplicationData, len(vaccineApplications))
	for i, va := range vaccineApplications {
		vaccineApplicationData[i] = h.mapToVaccineApplicationData(&va)
	}

	response := VaccineApplicationsResponse{
		Success: true,
		Data:    vaccineApplicationData,
		Message: "Aplicações de vacinas recuperadas com sucesso",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(response)
}

// GetVaccineApplicationByID obtém uma aplicação de vacina por ID
// @Summary      Obter aplicação de vacina por ID
// @Description  Retorna os dados de uma aplicação de vacina específica pelo ID
// @Tags         vaccine-applications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID da aplicação"
// @Success      200  {object}  VaccineApplicationResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/vaccine-applications/{id} [get]
func (h *VaccineApplicationHandler) GetVaccineApplicationByID(w http.ResponseWriter, r *http.Request) {
	applicationIDStr := chi.URLParam(r, "id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da aplicação inválido", http.StatusBadRequest)
		return
	}

	vaccineApplication, err := h.service.GetApplicationByID(uint(applicationID))
	if err != nil {
		SendErrorResponse(w, "Erro ao buscar aplicação de vacina: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if vaccineApplication == nil {
		SendErrorResponse(w, "Aplicação de vacina não encontrada", http.StatusNotFound)
		return
	}

	response := VaccineApplicationResponse{
		Success: true,
		Data:    h.mapToVaccineApplicationData(vaccineApplication),
		Message: "Aplicação de vacina encontrada com sucesso",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(response)
}

// UpdateVaccineApplication atualiza uma aplicação de vacina
// @Summary      Atualizar aplicação de vacina
// @Description  Atualiza os dados de uma aplicação de vacina existente
// @Tags         vaccine-applications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID da aplicação"
// @Param        request body CreateVaccineApplicationRequest true "Dados atualizados"
// @Success      200  {object}  VaccineApplicationResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/vaccine-applications/{id} [put]
func (h *VaccineApplicationHandler) UpdateVaccineApplication(w http.ResponseWriter, r *http.Request) {
	applicationIDStr := chi.URLParam(r, "id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da aplicação inválido", http.StatusBadRequest)
		return
	}

	var req CreateVaccineApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse(DateFormatISO, req.ApplicationDate)
	if err != nil {
		SendErrorResponse(w, "Formato de data inválido. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	vaccineApplication := &models.VaccineApplication{
		ID:              uint(applicationID),
		AnimalID:        req.AnimalID,
		VaccineID:       req.VaccineID,
		ApplicationDate: date,
		BatchNumber:     req.BatchNumber,
		Veterinarian:    req.Veterinarian,
		Observations:    req.Observations,
	}

	if err := h.service.UpdateApplication(vaccineApplication); err != nil {
		SendErrorResponse(w, "Erro ao atualizar aplicação de vacina: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedApplication, err := h.service.GetApplicationByID(vaccineApplication.ID)
	if err != nil {
		SendErrorResponse(w, "Erro ao recuperar aplicação atualizada: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := VaccineApplicationResponse{
		Success: true,
		Data:    h.mapToVaccineApplicationData(updatedApplication),
		Message: "Aplicação de vacina atualizada com sucesso",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteVaccineApplication deleta uma aplicação de vacina
// @Summary      Deletar aplicação de vacina
// @Description  Remove uma aplicação de vacina do sistema
// @Tags         vaccine-applications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID da aplicação"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/vaccine-applications/{id} [delete]
func (h *VaccineApplicationHandler) DeleteVaccineApplication(w http.ResponseWriter, r *http.Request) {
	applicationIDStr := chi.URLParam(r, "id")
	applicationID, err := strconv.ParseUint(applicationIDStr, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da aplicação inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteApplication(uint(applicationID)); err != nil {
		SendErrorResponse(w, "Erro ao deletar aplicação de vacina: "+err.Error(), http.StatusBadRequest)
		return
	}

	SendSuccessResponse(w, nil, "Aplicação de vacina deletada com sucesso", http.StatusOK)
}

func (h *VaccineApplicationHandler) mapToVaccineApplicationData(va *models.VaccineApplication) VaccineApplicationData {
	vaccineData := VaccineData{
		ID:           va.Vaccine.ID,
		FarmID:      va.Vaccine.FarmID,
		Name:         va.Vaccine.Name,
		Description:  va.Vaccine.Description,
		Manufacturer: va.Vaccine.Manufacturer,
		CreatedAt:    va.Vaccine.CreatedAt.Format(DateFormatISO8601),
		UpdatedAt:    va.Vaccine.UpdatedAt.Format(DateFormatISO8601),
	}

	return VaccineApplicationData{
		ID:              va.ID,
		AnimalID:        va.AnimalID,
		Animal: AnimalData{
			ID:                   va.Animal.ID,
			FarmID:               va.Animal.FarmID,
			EarTagNumberLocal:    va.Animal.EarTagNumberLocal,
			EarTagNumberRegister: va.Animal.EarTagNumberRegister,
			AnimalName:           va.Animal.AnimalName,
			Sex:                  va.Animal.Sex,
			Breed:                va.Animal.Breed,
			Type:                 va.Animal.Type,
			BirthDate:            formatBirthDate(va.Animal.BirthDate),
			Confinement:          va.Animal.Confinement,
			AnimalType:           va.Animal.AnimalType,
			Status:               va.Animal.Status,
			Fertilization:        va.Animal.Fertilization,
			Castrated:            va.Animal.Castrated,
			Purpose:              va.Animal.Purpose,
			CurrentBatch:         va.Animal.CurrentBatch,
		},
		VaccineID:       va.VaccineID,
		Vaccine:         vaccineData,
		ApplicationDate: va.ApplicationDate.Format(DateFormatISO),
		BatchNumber:     va.BatchNumber,
		Veterinarian:    va.Veterinarian,
		Observations:    va.Observations,
		CreatedAt:       va.CreatedAt.Format(DateFormatISO8601),
		UpdatedAt:       va.UpdatedAt.Format(DateFormatISO8601),
	}
}

