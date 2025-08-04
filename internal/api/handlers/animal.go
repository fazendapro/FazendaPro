package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
)

type AnimalHandler struct {
	service *service.AnimalService
}

func NewAnimalHandler(service *service.AnimalService) *AnimalHandler {
	return &AnimalHandler{service: service}
}

type CreateAnimalRequest struct {
	models.Animal
}

type AnimalResponse struct {
	ID                   uint   `json:"id"`
	FarmID               uint   `json:"farm_id"`
	EarTagNumberLocal    int    `json:"ear_tag_number_local"`
	EarTagNumberRegister int    `json:"ear_tag_number_register"`
	AnimalName           string `json:"animal_name"`
	Sex                  int    `json:"sex"`
	Breed                string `json:"breed"`
	Type                 string `json:"type"`
	BirthDate            string `json:"birth_date,omitempty"`
	Photo                string `json:"photo,omitempty"`
	FatherID             *uint  `json:"father_id,omitempty"`
	MotherID             *uint  `json:"mother_id,omitempty"`
	Confinement          bool   `json:"confinement"`
	AnimalType           int    `json:"animal_type"`
	Status               int    `json:"status"`
	Fertilization        bool   `json:"fertilization"`
	Castrated            bool   `json:"castrated"`
	Purpose              int    `json:"purpose"`
	CurrentBatch         int    `json:"current_batch"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
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

	if err := h.service.CreateAnimal(&req.Animal); err != nil {
		SendErrorResponse(w, "Erro ao criar animal: "+err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"id": req.Animal.ID,
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

	SendSuccessResponse(w, animal, "Animal encontrado com sucesso", http.StatusOK)
}

func (h *AnimalHandler) GetAnimalsByFarm(w http.ResponseWriter, r *http.Request) {
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

	animals, err := h.service.GetAnimalsByFarmID(uint(id))
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendSuccessResponse(w, animals, "Animais encontrados com sucesso", http.StatusOK)
}

func (h *AnimalHandler) UpdateAnimal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var animal models.Animal
	if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
		SendErrorResponse(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

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
