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

// MilkCollectionData representa os dados de uma coleta de leite
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

// CreateMilkCollection cria uma nova coleta de leite
func (h *MilkCollectionHandler) CreateMilkCollection(w http.ResponseWriter, r *http.Request) {
	var req CreateMilkCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse da data
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

	// Buscar a coleta criada com os dados do animal
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

// UpdateMilkCollection atualiza uma coleta de leite existente
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

	// Parse da data
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

	// Buscar a coleta atualizada com os dados do animal
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

// GetMilkCollectionsByFarmID obtém todas as coletas de leite de uma fazenda
func (h *MilkCollectionHandler) GetMilkCollectionsByFarmID(w http.ResponseWriter, r *http.Request) {
	farmIDStr := chi.URLParam(r, "farmId")
	farmID, err := strconv.ParseUint(farmIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid farm ID", http.StatusBadRequest)
		return
	}

	// Verificar se há filtros de data
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

// GetMilkCollectionsByAnimalID obtém todas as coletas de leite de um animal específico
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

// mapToMilkCollectionData converte um modelo para a estrutura de resposta
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
