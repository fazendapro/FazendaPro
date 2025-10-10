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

	date, err := time.Parse("2006-01-02", req.Date)
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

	w.Header().Set("Content-Type", "application/json")
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

	date, err := time.Parse("2006-01-02", req.Date)
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

	fmt.Printf("DEBUG: Updating milk collection with ID: %d, AnimalID: %d, Liters: %.2f, Date: %s\n",
		milkCollection.ID, milkCollection.AnimalID, milkCollection.Liters, milkCollection.Date.Format("2006-01-02"))

	if err := h.service.UpdateMilkCollection(milkCollection); err != nil {
		fmt.Printf("DEBUG: Error updating milk collection: %v\n", err)
		http.Error(w, "Failed to update milk collection", http.StatusInternalServerError)
		return
	}

	fmt.Printf("DEBUG: Milk collection updated successfully\n")

	updatedMilkCollection, err := h.service.GetMilkCollectionByID(milkCollection.ID)
	if err != nil {
		fmt.Printf("DEBUG: Error retrieving updated milk collection: %v\n", err)
		http.Error(w, "Failed to retrieve updated milk collection", http.StatusInternalServerError)
		return
	}

	fmt.Printf("DEBUG: Retrieved updated milk collection: ID=%d, Liters=%.2f\n",
		updatedMilkCollection.ID, updatedMilkCollection.Liters)

	response := MilkCollectionResponse{
		Success: true,
		Data:    h.mapToMilkCollectionData(updatedMilkCollection),
		Message: "Milk collection updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
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
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}

	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
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

	w.Header().Set("Content-Type", "application/json")
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

	w.Header().Set("Content-Type", "application/json")
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
	return birthDate.Format("2006-01-02")
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

func (h *MilkCollectionHandler) GetTopMilkProducers(w http.ResponseWriter, r *http.Request) {
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

	startDate := time.Now().AddDate(0, 0, -periodDays)
	endDate := time.Now()
	milkCollections, err := h.service.GetMilkCollectionsByFarmIDWithDateRange(uint(id), &startDate, &endDate)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	animalStats := make(map[uint]*struct {
		AnimalID           uint
		AnimalName         string
		EarTagNumberLocal  int
		Photo              string
		TotalProduction    float64
		CollectionCount    int
		FatContent         float64
		LastCollectionDate time.Time
		DaysInLactation    int
	})

	for _, mc := range milkCollections {
		if stats, exists := animalStats[mc.AnimalID]; exists {
			stats.TotalProduction += mc.Liters
			stats.CollectionCount++
			if mc.Date.After(stats.LastCollectionDate) {
				stats.LastCollectionDate = mc.Date
			}
		} else {
			daysInLactation := int(time.Since(mc.Date).Hours()/24) + 60
			if daysInLactation < 0 {
				daysInLactation = 0
			}

			animalStats[mc.AnimalID] = &struct {
				AnimalID           uint
				AnimalName         string
				EarTagNumberLocal  int
				Photo              string
				TotalProduction    float64
				CollectionCount    int
				FatContent         float64
				LastCollectionDate time.Time
				DaysInLactation    int
			}{
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

	var responses []TopMilkProducerResponse
	for _, stats := range animalStats {
		averageDailyProduction := stats.TotalProduction / float64(stats.CollectionCount)

		response := TopMilkProducerResponse{
			ID:                     stats.AnimalID,
			AnimalName:             stats.AnimalName,
			EarTagNumberLocal:      stats.EarTagNumberLocal,
			Photo:                  stats.Photo,
			TotalProduction:        stats.TotalProduction,
			AverageDailyProduction: averageDailyProduction,
			FatContent:             stats.FatContent,
			LastCollectionDate:     stats.LastCollectionDate.Format("2006-01-02"),
			DaysInLactation:        stats.DaysInLactation,
		}
		responses = append(responses, response)
	}

	for i := 0; i < len(responses); i++ {
		for j := i + 1; j < len(responses); j++ {
			if responses[i].TotalProduction < responses[j].TotalProduction {
				responses[i], responses[j] = responses[j], responses[i]
			}
		}
	}

	if len(responses) > limit {
		responses = responses[:limit]
	}

	SendSuccessResponse(w, responses, fmt.Sprintf("Maiores produtoras de leite encontradas com sucesso (%d registros)", len(responses)), http.StatusOK)
}
