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

type MilkCollectionHandler struct {
	service *service.MilkCollectionService
}

func NewMilkCollectionHandler(service *service.MilkCollectionService) *MilkCollectionHandler {
	return &MilkCollectionHandler{service: service}
}

type MilkCollectionData struct {
	ID        uint       `json:"id"`
	AnimalID  uint       `json:"animal_id"`
	Animal    AnimalData `json:"animal"`
	Liters    float64    `json:"liters"`
	Date      time.Time  `json:"date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateMilkCollectionRequest struct {
	AnimalID uint    `json:"animal_id" validate:"required"`
	Liters   float64 `json:"liters" validate:"required,min=0"`
	Date     string  `json:"date" validate:"required"`
}

type MilkCollectionResponse struct {
	Success bool               `json:"success"`
	Data    MilkCollectionData `json:"data,omitempty"`
	Message string             `json:"message,omitempty"`
}

type MilkCollectionsResponse struct {
	Success bool                 `json:"success"`
	Data    []MilkCollectionData `json:"data,omitempty"`
	Message string               `json:"message,omitempty"`
}

func (h *MilkCollectionHandler) CreateMilkCollection(w http.ResponseWriter, r *http.Request) {
	var req CreateMilkCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(DateFormatISO, req.Date)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	milkCollection := &models.MilkCollection{
		AnimalID: req.AnimalID,
		Liters:   req.Liters,
		Date:     date,
	}

	if err := h.service.CreateMilkCollection(milkCollection); err != nil {
		http.Error(w, "Failed to create milk collection", http.StatusInternalServerError)
		return
	}

	createdMilkCollection, err := h.service.GetMilkCollectionByID(milkCollection.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve created milk collection", http.StatusInternalServerError)
		return
	}

	response := MilkCollectionResponse{
		Success: true,
		Data:    h.mapToMilkCollectionData(createdMilkCollection),
		Message: "Milk collection created successfully",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *MilkCollectionHandler) UpdateMilkCollection(w http.ResponseWriter, r *http.Request) {
	milkCollectionIDStr := chi.URLParam(r, "id")
	milkCollectionID, err := strconv.ParseUint(milkCollectionIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid milk collection ID", http.StatusBadRequest)
		return
	}

	var req CreateMilkCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(DateFormatISO, req.Date)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	milkCollection := &models.MilkCollection{
		ID:       uint(milkCollectionID),
		AnimalID: req.AnimalID,
		Liters:   req.Liters,
		Date:     date,
	}

	if err := h.service.UpdateMilkCollection(milkCollection); err != nil {
		http.Error(w, "Failed to update milk collection", http.StatusInternalServerError)
		return
	}

	updatedMilkCollection, err := h.service.GetMilkCollectionByID(milkCollection.ID)
	if err != nil {
		http.Error(w, "Failed to retrieve updated milk collection", http.StatusInternalServerError)
		return
	}

	response := MilkCollectionResponse{
		Success: true,
		Data:    h.mapToMilkCollectionData(updatedMilkCollection),
		Message: "Milk collection updated successfully",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *MilkCollectionHandler) GetMilkCollectionsByFarmID(w http.ResponseWriter, r *http.Request) {
	farmIDStr := chi.URLParam(r, "farmId")
	farmID, err := strconv.ParseUint(farmIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid farm ID", http.StatusBadRequest)
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

	var milkCollections []models.MilkCollection
	if startDate != nil || endDate != nil {
		milkCollections, err = h.service.GetMilkCollectionsByFarmIDWithDateRange(uint(farmID), startDate, endDate)
	} else {
		milkCollections, err = h.service.GetMilkCollectionsByFarmID(uint(farmID))
	}

	if err != nil {
		http.Error(w, "Failed to retrieve milk collections", http.StatusInternalServerError)
		return
	}

	milkCollectionData := make([]MilkCollectionData, len(milkCollections))
	for i, mc := range milkCollections {
		milkCollectionData[i] = h.mapToMilkCollectionData(&mc)
	}

	response := MilkCollectionsResponse{
		Success: true,
		Data:    milkCollectionData,
		Message: "Milk collections retrieved successfully",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(response)
}

func (h *MilkCollectionHandler) GetMilkCollectionsByAnimalID(w http.ResponseWriter, r *http.Request) {
	animalIDStr := chi.URLParam(r, "animalId")
	animalID, err := strconv.ParseUint(animalIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid animal ID", http.StatusBadRequest)
		return
	}

	milkCollections, err := h.service.GetMilkCollectionsByAnimalID(uint(animalID))
	if err != nil {
		http.Error(w, "Failed to retrieve milk collections", http.StatusInternalServerError)
		return
	}

	milkCollectionData := make([]MilkCollectionData, len(milkCollections))
	for i, mc := range milkCollections {
		milkCollectionData[i] = h.mapToMilkCollectionData(&mc)
	}

	response := MilkCollectionsResponse{
		Success: true,
		Data:    milkCollectionData,
		Message: "Milk collections retrieved successfully",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	json.NewEncoder(w).Encode(response)
}

func (h *MilkCollectionHandler) mapToMilkCollectionData(mc *models.MilkCollection) MilkCollectionData {
	return MilkCollectionData{
		ID:       mc.ID,
		AnimalID: mc.AnimalID,
		Animal: AnimalData{
			ID:                   mc.Animal.ID,
			FarmID:               mc.Animal.FarmID,
			EarTagNumberLocal:    mc.Animal.EarTagNumberLocal,
			EarTagNumberRegister: mc.Animal.EarTagNumberRegister,
			AnimalName:           mc.Animal.AnimalName,
			Sex:                  mc.Animal.Sex,
			Breed:                mc.Animal.Breed,
			Type:                 mc.Animal.Type,
			BirthDate:            formatBirthDate(mc.Animal.BirthDate),
			Confinement:          mc.Animal.Confinement,
			AnimalType:           mc.Animal.AnimalType,
			Status:               mc.Animal.Status,
			Fertilization:        mc.Animal.Fertilization,
			Castrated:            mc.Animal.Castrated,
			Purpose:              mc.Animal.Purpose,
			CurrentBatch:         mc.Animal.CurrentBatch,
		},
		Liters:    mc.Liters,
		Date:      mc.Date,
		CreatedAt: mc.CreatedAt,
		UpdatedAt: mc.UpdatedAt,
	}
}

func formatBirthDate(birthDate *time.Time) string {
	if birthDate == nil {
		return ""
	}
	return birthDate.Format(DateFormatISO)
}

type TopMilkProducerResponse struct {
	ID                     uint    `json:"id"`
	AnimalName             string  `json:"animal_name"`
	EarTagNumberLocal      int     `json:"ear_tag_number_local"`
	Photo                  string  `json:"photo"`
	TotalProduction        float64 `json:"total_production"`
	AverageDailyProduction float64 `json:"average_daily_production"`
	FatContent             float64 `json:"fat_content"`
	LastCollectionDate     string  `json:"last_collection_date"`
	DaysInLactation        int     `json:"days_in_lactation"`
}

type animalStats struct {
	AnimalID           uint
	AnimalName         string
	EarTagNumberLocal  int
	Photo              string
	TotalProduction    float64
	CollectionCount    int
	FatContent         float64
	LastCollectionDate time.Time
	DaysInLactation    int
}

func calculateAnimalStats(milkCollections []models.MilkCollection) map[uint]*animalStats {
	stats := make(map[uint]*animalStats)

	for _, mc := range milkCollections {
		if existing, exists := stats[mc.AnimalID]; exists {
			existing.TotalProduction += mc.Liters
			existing.CollectionCount++
			if mc.Date.After(existing.LastCollectionDate) {
				existing.LastCollectionDate = mc.Date
			}
		} else {
			daysInLactation := int(time.Since(mc.Date).Hours()/24) + 60
			if daysInLactation < 0 {
				daysInLactation = 0
			}

			stats[mc.AnimalID] = &animalStats{
				AnimalID:           mc.AnimalID,
				AnimalName:         mc.Animal.AnimalName,
				EarTagNumberLocal:  mc.Animal.EarTagNumberLocal,
				Photo:              mc.Animal.Photo,
				TotalProduction:    mc.Liters,
				CollectionCount:    1,
				FatContent:         3.5,
				LastCollectionDate: mc.Date,
				DaysInLactation:    daysInLactation,
			}
		}
	}

	return stats
}

func buildTopMilkProducerResponses(stats map[uint]*animalStats) []TopMilkProducerResponse {
	responses := make([]TopMilkProducerResponse, 0, len(stats))

	for _, s := range stats {
		averageDailyProduction := s.TotalProduction / float64(s.CollectionCount)
		responses = append(responses, TopMilkProducerResponse{
			ID:                     s.AnimalID,
			AnimalName:             s.AnimalName,
			EarTagNumberLocal:      s.EarTagNumberLocal,
			Photo:                  s.Photo,
			TotalProduction:        s.TotalProduction,
			AverageDailyProduction: averageDailyProduction,
			FatContent:             s.FatContent,
			LastCollectionDate:     s.LastCollectionDate.Format(DateFormatISO),
			DaysInLactation:        s.DaysInLactation,
		})
	}

	return responses
}

func sortByTotalProduction(responses []TopMilkProducerResponse) {
	for i := 0; i < len(responses); i++ {
		for j := i + 1; j < len(responses); j++ {
			if responses[i].TotalProduction < responses[j].TotalProduction {
				responses[i], responses[j] = responses[j], responses[i]
			}
		}
	}
}

func parseTopMilkProducersParams(r *http.Request) (uint, int, int, error) {
	farmID := r.URL.Query().Get("farmId")
	if farmID == "" {
		return 0, 0, 0, fmt.Errorf("ID da fazenda é obrigatório")
	}

	id, err := strconv.ParseUint(farmID, 10, 32)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("ID da fazenda inválido")
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	periodDays := 30
	if periodStr := r.URL.Query().Get("periodDays"); periodStr != "" {
		if parsedPeriod, err := strconv.Atoi(periodStr); err == nil && parsedPeriod > 0 {
			periodDays = parsedPeriod
		}
	}

	return uint(id), limit, periodDays, nil
}

func (h *MilkCollectionHandler) GetTopMilkProducers(w http.ResponseWriter, r *http.Request) {
	farmID, limit, periodDays, err := parseTopMilkProducersParams(r)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate := time.Now().AddDate(0, 0, -periodDays)
	endDate := time.Now()
	milkCollections, err := h.service.GetMilkCollectionsByFarmIDWithDateRange(farmID, &startDate, &endDate)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats := calculateAnimalStats(milkCollections)
	responses := buildTopMilkProducerResponses(stats)
	sortByTotalProduction(responses)

	if len(responses) > limit {
		responses = responses[:limit]
	}

	SendSuccessResponse(w, responses, fmt.Sprintf("Maiores produtoras de leite encontradas com sucesso (%d registros)", len(responses)), http.StatusOK)
}
