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

type AnimalHandler struct {
	service *service.AnimalService
}

func NewAnimalHandler(service *service.AnimalService) *AnimalHandler {
	return &AnimalHandler{service: service}
}

// AnimalData representa os dados base de um animal
type AnimalData struct {
	ID                   uint   `json:"id"`
	FarmID               uint   `json:"farmID"`
	EarTagNumberLocal    int    `json:"earringNumber"`
	EarTagNumberRegister int    `json:"earringNumberGlobal"`
	AnimalName           string `json:"animalName"`
	Sex                  int    `json:"sex"`
	Breed                string `json:"breed"`
	Type                 string `json:"type"`
	BirthDate            string `json:"birthDate,omitempty"`
	Photo                string `json:"photo,omitempty"`
	FatherID             *uint  `json:"fatherID,omitempty"`
	MotherID             *uint  `json:"motherID,omitempty"`
	Confinement          bool   `json:"confinement"`
	AnimalType           int    `json:"animalType"`
	Status               int    `json:"status"`
	Fertilization        bool   `json:"fertilization"`
	Castrated            bool   `json:"castrated"`
	Purpose              int    `json:"purpose"`
	CurrentBatch         int    `json:"currentBatch"`
}

type CreateAnimalRequest struct {
	AnimalData
}

type AnimalResponse struct {
	AnimalData
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func animalDataToModel(data AnimalData) models.Animal {
	var birthDate *time.Time
	if data.BirthDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", data.BirthDate); err == nil {
			birthDate = &parsedDate
		}
	}

	return models.Animal{
		ID:                   data.ID,
		FarmID:               data.FarmID,
		EarTagNumberLocal:    data.EarTagNumberLocal,
		EarTagNumberRegister: data.EarTagNumberRegister,
		AnimalName:           data.AnimalName,
		Sex:                  data.Sex,
		Breed:                data.Breed,
		Type:                 data.Type,
		BirthDate:            birthDate,
		Photo:                data.Photo,
		FatherID:             data.FatherID,
		MotherID:             data.MotherID,
		Confinement:          data.Confinement,
		AnimalType:           data.AnimalType,
		Status:               data.Status,
		Fertilization:        data.Fertilization,
		Castrated:            data.Castrated,
		Purpose:              data.Purpose,
		CurrentBatch:         data.CurrentBatch,
	}
}

func modelToAnimalResponse(animal *models.Animal) AnimalResponse {
	var birthDate string
	if animal.BirthDate != nil {
		birthDate = animal.BirthDate.Format("2006-01-02")
	}

	return AnimalResponse{
		AnimalData: AnimalData{
			ID:                   animal.ID,
			FarmID:               animal.FarmID,
			EarTagNumberLocal:    animal.EarTagNumberLocal,
			EarTagNumberRegister: animal.EarTagNumberRegister,
			AnimalName:           animal.AnimalName,
			Sex:                  animal.Sex,
			Breed:                animal.Breed,
			Type:                 animal.Type,
			BirthDate:            birthDate,
			Photo:                animal.Photo,
			FatherID:             animal.FatherID,
			MotherID:             animal.MotherID,
			Confinement:          animal.Confinement,
			AnimalType:           animal.AnimalType,
			Status:               animal.Status,
			Fertilization:        animal.Fertilization,
			Castrated:            animal.Castrated,
			Purpose:              animal.Purpose,
			CurrentBatch:         animal.CurrentBatch,
		},
		CreatedAt: animal.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: animal.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (h *AnimalHandler) CreateAnimal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req CreateAnimalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	animal := animalDataToModel(req.AnimalData)

	if err := h.service.CreateAnimal(&animal); err != nil {
		SendErrorResponse(w, "Erro ao criar animal: "+err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"id": animal.ID,
	}
	SendSuccessResponse(w, data, "Animal criado com sucesso", http.StatusCreated)
}

func (h *AnimalHandler) GetAnimal(w http.ResponseWriter, r *http.Request) {
	animalID := r.URL.Query().Get("id")
	if animalID == "" {
		SendErrorResponse(w, "ID do animal é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(animalID, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID do animal inválido", http.StatusBadRequest)
		return
	}

	animal, err := h.service.GetAnimalByID(uint(id))
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if animal == nil {
		SendErrorResponse(w, "Animal não encontrado", http.StatusNotFound)
		return
	}

	response := modelToAnimalResponse(animal)
	SendSuccessResponse(w, response, "Animal encontrado com sucesso", http.StatusOK)
}

func (h *AnimalHandler) GetAnimalsByFarm(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GetAnimalsByFarm chamado - URL: %s, Query: %s\n", r.URL.Path, r.URL.RawQuery)

	farmID := r.URL.Query().Get("farmId")
	if farmID == "" {
		fmt.Printf("Erro: farmId não fornecido\n")
		SendErrorResponse(w, "ID da fazenda é obrigatório", http.StatusBadRequest)
		return
	}

	fmt.Printf("FarmID recebido: %s\n", farmID)

	id, err := strconv.ParseUint(farmID, 10, 32)
	if err != nil {
		fmt.Printf("Erro ao converter farmId: %v\n", err)
		SendErrorResponse(w, "ID da fazenda inválido", http.StatusBadRequest)
		return
	}

	fmt.Printf("FarmID convertido: %d\n", id)

	animals, err := h.service.GetAnimalsByFarmID(uint(id))
	if err != nil {
		fmt.Printf("Erro ao buscar animais: %v\n", err)
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []AnimalResponse
	for _, animal := range animals {
		responses = append(responses, modelToAnimalResponse(&animal))
	}

	fmt.Printf("FarmID: %d, Animais encontrados: %d\n", id, len(animals))

	SendSuccessResponse(w, responses, fmt.Sprintf("Animais encontrados com sucesso (%d animais)", len(animals)), http.StatusOK)
}

func (h *AnimalHandler) UpdateAnimal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req CreateAnimalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	animal := animalDataToModel(req.AnimalData)

	if err := h.service.UpdateAnimal(&animal); err != nil {
		SendErrorResponse(w, "Erro ao atualizar animal: "+err.Error(), http.StatusBadRequest)
		return
	}

	SendSuccessResponse(w, nil, "Animal atualizado com sucesso", http.StatusOK)
}

func (h *AnimalHandler) DeleteAnimal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	animalID := r.URL.Query().Get("id")
	if animalID == "" {
		SendErrorResponse(w, "ID do animal é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(animalID, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID do animal inválido", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAnimal(uint(id)); err != nil {
		SendErrorResponse(w, "Erro ao deletar animal: "+err.Error(), http.StatusBadRequest)
		return
	}

	SendSuccessResponse(w, nil, "Animal deletado com sucesso", http.StatusOK)
}
