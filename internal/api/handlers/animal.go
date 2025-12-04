package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
)

type AnimalHandler struct {
	service       *service.AnimalService
	weightService *service.WeightService
}

func NewAnimalHandler(service *service.AnimalService) *AnimalHandler {
	return &AnimalHandler{service: service}
}

func NewAnimalHandlerWithWeight(animalService *service.AnimalService, weightService *service.WeightService) *AnimalHandler {
	return &AnimalHandler{
		service:       animalService,
		weightService: weightService,
	}
}

// AnimalData representa os dados básicos de um animal
// @Description Dados principais de um animal na fazenda
type AnimalData struct {
	ID                   uint     `json:"id" example:"1"`                                       // ID único do animal
	FarmID               uint     `json:"farm_id" example:"1"`                                  // ID da fazenda
	EarTagNumberLocal    int      `json:"ear_tag_number_local" example:"123"`                   // Número da brinco local
	EarTagNumberRegister int      `json:"ear_tag_number_register" example:"456"`                // Número de registro
	AnimalName           string   `json:"animal_name" example:"Branquinha"`                     // Nome do animal
	Sex                  int      `json:"sex" example:"1"`                                      // Sexo (1=Fêmea, 2=Macho)
	Breed                string   `json:"breed" example:"Holandesa"`                            // Raça
	Type                 string   `json:"type" example:"Bovino"`                                // Tipo de animal
	BirthDate            string   `json:"birth_date,omitempty" example:"2020-01-15"`            // Data de nascimento (ISO format)
	Photo                string   `json:"photo,omitempty" example:"data:image/jpeg;base64,..."` // Foto em base64
	FatherID             *uint    `json:"father_id,omitempty" example:"10"`                     // ID do pai (opcional)
	MotherID             *uint    `json:"mother_id,omitempty" example:"20"`                     // ID da mãe (opcional)
	Confinement          bool     `json:"confinement" example:"false"`                          // Se está em confinamento
	AnimalType           int      `json:"animal_type" example:"1"`                              // Tipo de animal (enum)
	Status               int      `json:"status" example:"1"`                                   // Status do animal (enum)
	Fertilization        bool     `json:"fertilization" example:"true"`                         // Se está fertilizada
	Castrated            bool     `json:"castrated" example:"false"`                            // Se está castrado
	Purpose              int      `json:"purpose" example:"1"`                                  // Propósito do animal (enum)
	CurrentBatch         int      `json:"current_batch" example:"1"`                            // Lote atual
	Weight               *float64 `json:"weight,omitempty" example:"450.5"`                     // Peso atual do animal (kg)
}

// CreateAnimalRequest representa a requisição de criação de animal
// @Description Dados necessários para criar um novo animal
type CreateAnimalRequest struct {
	AnimalData // Dados do animal
}

// AnimalResponse representa a resposta com dados completos do animal
// @Description Resposta com dados do animal incluindo informações dos pais
type AnimalResponse struct {
	AnimalData               // Dados básicos do animal
	Father     *AnimalParent `json:"father,omitempty"`                         // Informações do pai (se disponível)
	Mother     *AnimalParent `json:"mother,omitempty"`                         // Informações da mãe (se disponível)
	CreatedAt  string        `json:"createdAt" example:"2024-01-15T10:30:00Z"` // Data de criação
	UpdatedAt  string        `json:"updatedAt" example:"2024-01-15T10:30:00Z"` // Data de atualização
}

// AnimalParent representa informações básicas de um animal pai/mãe
// @Description Informações resumidas do animal pai ou mãe
type AnimalParent struct {
	ID                uint   `json:"id" example:"10"`                    // ID do animal pai/mãe
	AnimalName        string `json:"animal_name" example:"Touro Bravo"`  // Nome do animal
	EarTagNumberLocal int    `json:"ear_tag_number_local" example:"100"` // Número da brinco
}

func animalDataToModel(data AnimalData) models.Animal {
	var birthDate *time.Time
	if data.BirthDate != "" {
		if parsedDate, err := time.Parse(DateFormatISO, data.BirthDate); err == nil {
			birthDate = &parsedDate
		} else {
			fmt.Printf("Erro ao fazer parse da data: %v\n", err)
		}
	}

	animal := models.Animal{
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

	return animal
}

func modelToAnimalResponse(animal *models.Animal) AnimalResponse {
	var birthDate string
	if animal.BirthDate != nil {
		birthDate = animal.BirthDate.Format(DateFormatISO)
	}

	var father *AnimalParent
	if animal.Father != nil {
		father = &AnimalParent{
			ID:                animal.Father.ID,
			AnimalName:        animal.Father.AnimalName,
			EarTagNumberLocal: animal.Father.EarTagNumberLocal,
		}
	}

	var mother *AnimalParent
	if animal.Mother != nil {
		mother = &AnimalParent{
			ID:                animal.Mother.ID,
			AnimalName:        animal.Mother.AnimalName,
			EarTagNumberLocal: animal.Mother.EarTagNumberLocal,
		}
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
			Weight:               nil, // Será preenchido no handler se necessário
		},
		Father:    father,
		Mother:    mother,
		CreatedAt: animal.CreatedAt.Format(DateFormatDateTime),
		UpdatedAt: animal.UpdatedAt.Format(DateFormatDateTime),
	}
}

// CreateAnimal cria um novo animal
// @Summary      Criar animal
// @Description  Cria um novo animal na fazenda
// @Tags         animals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateAnimalRequest true "Dados do animal"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/animals [post]
func (h *AnimalHandler) CreateAnimal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req CreateAnimalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
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

// GetAnimal obtém um animal por ID
// @Summary      Obter animal por ID
// @Description  Retorna os dados de um animal específico pelo ID
// @Tags         animals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id query int true "ID do animal"
// @Success      200  {object}  AnimalResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/animals [get]
func (h *AnimalHandler) GetAnimal(w http.ResponseWriter, r *http.Request) {
	animalID := r.URL.Query().Get("id")
	if animalID == "" {
		SendErrorResponse(w, ErrAnimalIDRequired, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(animalID, 10, 32)
	if err != nil {
		SendErrorResponse(w, ErrInvalidAnimalID, http.StatusBadRequest)
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

// GetAnimalsByFarm obtém todos os animais de uma fazenda
// @Summary      Obter animais por fazenda
// @Description  Retorna lista de todos os animais de uma fazenda específica
// @Tags         animals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        farmId query int true "ID da fazenda"
// @Success      200  {array}   AnimalResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/animals/farm [get]
func (h *AnimalHandler) GetAnimalsByFarm(w http.ResponseWriter, r *http.Request) {
	farmID := r.URL.Query().Get("farmId")
	if farmID == "" {
		fmt.Printf("Erro: farmId não fornecido\n")

		SendErrorResponse(w, "ID da fazenda é obrigatório", http.StatusBadRequest)
		return
	}

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
		response := modelToAnimalResponse(&animal)

		if h.weightService != nil {
			weight, err := h.weightService.GetWeightByAnimalID(animal.ID)
			if err != nil {
				fmt.Printf("Erro ao buscar peso para animal ID %d: %v\n", animal.ID, err)
			} else if weight != nil {
				response.Weight = &weight.AnimalWeight
				fmt.Printf("Peso encontrado para animal ID %d: %.2f kg\n", animal.ID, weight.AnimalWeight)
			} else {
				fmt.Printf("Peso não encontrado para animal ID %d (weight é nil)\n", animal.ID)
			}
		} else {
			fmt.Printf("WeightService não disponível no AnimalHandler para animal ID %d\n", animal.ID)
		}

		fmt.Printf("Response antes de adicionar: Weight = %v\n", response.Weight)
		responses = append(responses, response)
	}

	fmt.Printf("FarmID: %d, Animais encontrados: %d\n", id, len(animals))

	SendSuccessResponse(w, responses, fmt.Sprintf("Animais encontrados com sucesso (%d animais)", len(animals)), http.StatusOK)
}

// UpdateAnimal atualiza os dados de um animal
// @Summary      Atualizar animal
// @Description  Atualiza os dados de um animal existente
// @Tags         animals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateAnimalRequest true "Dados atualizados do animal"
// @Success      200  {object}  AnimalResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/animals [put]
func (h *AnimalHandler) UpdateAnimal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req CreateAnimalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	animal := animalDataToModel(req.AnimalData)

	if err := h.service.UpdateAnimal(&animal); err != nil {
		SendErrorResponse(w, "Erro ao atualizar animal: "+err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.service.GetAnimalByID(animal.ID)
	if err != nil || updated == nil {
		SendSuccessResponse(w, nil, "Animal atualizado com sucesso", http.StatusOK)
		return
	}

	response := modelToAnimalResponse(updated)
	SendSuccessResponse(w, response, "Animal atualizado com sucesso", http.StatusOK)
}

// DeleteAnimal remove um animal
// @Summary      Deletar animal
// @Description  Remove um animal do sistema
// @Tags         animals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id query int true "ID do animal"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/animals [delete]
func (h *AnimalHandler) DeleteAnimal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	animalID := r.URL.Query().Get("id")
	if animalID == "" {
		SendErrorResponse(w, ErrAnimalIDRequired, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(animalID, 10, 32)
	if err != nil {
		SendErrorResponse(w, ErrInvalidAnimalID, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteAnimal(uint(id)); err != nil {
		SendErrorResponse(w, "Erro ao deletar animal: "+err.Error(), http.StatusBadRequest)
		return
	}

	SendSuccessResponse(w, nil, "Animal deletado com sucesso", http.StatusOK)
}

// GetAnimalsBySex obtém animais filtrados por sexo
// @Summary      Obter animais por sexo
// @Description  Retorna lista de animais de uma fazenda filtrados por sexo
// @Tags         animals
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        farmId query int true "ID da fazenda"
// @Param        sex query int true "Sexo (1=Fêmea, 2=Macho)"
// @Success      200  {array}   AnimalResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/animals/sex [get]
func (h *AnimalHandler) GetAnimalsBySex(w http.ResponseWriter, r *http.Request) {
	farmID := r.URL.Query().Get("farmId")
	sex := r.URL.Query().Get("sex")

	if farmID == "" {
		SendErrorResponse(w, "ID da fazenda é obrigatório", http.StatusBadRequest)
		return
	}

	if sex == "" {
		SendErrorResponse(w, "Sexo é obrigatório", http.StatusBadRequest)
		return
	}

	farmIDUint, err := strconv.ParseUint(farmID, 10, 32)
	if err != nil {
		SendErrorResponse(w, "ID da fazenda inválido", http.StatusBadRequest)
		return
	}

	sexInt, err := strconv.Atoi(sex)
	if err != nil {
		SendErrorResponse(w, "Sexo inválido", http.StatusBadRequest)
		return
	}

	animals, err := h.service.GetAnimalsByFarmIDAndSex(uint(farmIDUint), sexInt)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responses []AnimalResponse
	for _, animal := range animals {
		responses = append(responses, modelToAnimalResponse(&animal))
	}

	SendSuccessResponse(w, responses, fmt.Sprintf("Animais encontrados com sucesso (%d animais)", len(animals)), http.StatusOK)
}

// UploadAnimalPhoto faz upload de foto de um animal
// @Summary      Upload foto do animal
// @Description  Faz upload de uma foto para um animal específico
// @Tags         animals
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        animal_id formData int true "ID do animal"
// @Param        photo formData file true "Arquivo de imagem"
// @Success      200  {object}  AnimalResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/animals/photo [post]
func (h *AnimalHandler) UploadAnimalPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		SendErrorResponse(w, "Erro ao fazer parse do formulário: "+err.Error(), http.StatusBadRequest)
		return
	}

	animalID := r.FormValue("animal_id")
	if animalID == "" {
		SendErrorResponse(w, ErrAnimalIDRequired, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(animalID, 10, 32)
	if err != nil {
		SendErrorResponse(w, ErrInvalidAnimalID, http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("photo")
	if err != nil {
		SendErrorResponse(w, "Erro ao obter arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			SendErrorResponse(w, "Erro ao ler arquivo: "+err.Error(), http.StatusBadRequest)
			return
		}
		if n == 0 {
			break
		}
		fileBytes = append(fileBytes, buffer[:n]...)
	}

	photoBase64 := fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(fileBytes))

	animal, err := h.service.GetAnimalByID(uint(id))
	if err != nil || animal == nil {
		SendErrorResponse(w, "Animal não encontrado", http.StatusNotFound)
		return
	}

	animal.Photo = photoBase64
	if err := h.service.UpdateAnimal(animal); err != nil {
		SendErrorResponse(w, "Erro ao atualizar foto do animal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	updatedAnimal, err := h.service.GetAnimalByID(uint(id))
	if err != nil || updatedAnimal == nil {
		SendErrorResponse(w, "Erro ao buscar animal atualizado", http.StatusInternalServerError)
		return
	}

	response := modelToAnimalResponse(updatedAnimal)
	SendSuccessResponse(w, response, "Foto do animal atualizada com sucesso", http.StatusOK)
}
