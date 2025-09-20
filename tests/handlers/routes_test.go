package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	handlers.SendErrorResponse(w, "Test error", http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Test error", response["message"])
}

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"test": "data"}
	handlers.SendSuccessResponse(w, data, "Success message", http.StatusOK)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Success message", response["message"])
	assert.Equal(t, data["test"], response["data"].(map[string]interface{})["test"])
}

func TestLoginRequest(t *testing.T) {
	req := handlers.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
}

func TestRegisterRequest(t *testing.T) {
	user := models.User{
		ID:     1,
		FarmID: 1,
	}
	person := models.Person{
		Email:     "test@example.com",
		FirstName: "João",
		LastName:  "Silva",
	}

	req := handlers.RegisterRequest{
		User:   user,
		Person: person,
	}

	assert.Equal(t, user, req.User)
	assert.Equal(t, person, req.Person)
}

func TestRefreshTokenRequest(t *testing.T) {
	req := handlers.RefreshTokenRequest{
		RefreshToken: "refresh-token-123",
	}

	assert.Equal(t, "refresh-token-123", req.RefreshToken)
}

func TestCreateUserRequest(t *testing.T) {
	user := models.User{
		ID:     1,
		FarmID: 1,
	}
	person := models.Person{
		Email:     "test@example.com",
		FirstName: "João",
		LastName:  "Silva",
	}

	req := handlers.CreateUserRequest{
		User:   user,
		Person: person,
	}

	assert.Equal(t, user, req.User)
	assert.Equal(t, person, req.Person)
}

func TestAnimalData(t *testing.T) {
	animalData := handlers.AnimalData{
		ID:                   1,
		FarmID:               1,
		EarTagNumberLocal:    123,
		EarTagNumberRegister: 456,
		AnimalName:           "Bella",
		Sex:                  1,
		Breed:                "Holandesa",
		Type:                 "Vaca",
		BirthDate:            "2020-01-15",
		Confinement:          false,
		AnimalType:           1,
		Status:               1,
		Fertilization:        false,
		Castrated:            false,
		Purpose:              1,
		CurrentBatch:         1,
	}

	assert.Equal(t, uint(1), animalData.ID)
	assert.Equal(t, uint(1), animalData.FarmID)
	assert.Equal(t, 123, animalData.EarTagNumberLocal)
	assert.Equal(t, "Bella", animalData.AnimalName)
}

func TestCreateAnimalRequest(t *testing.T) {
	animalData := handlers.AnimalData{
		ID:                1,
		FarmID:            1,
		EarTagNumberLocal: 123,
		AnimalName:        "Bella",
		Sex:               1,
		Breed:             "Holandesa",
		Type:              "Vaca",
	}

	req := handlers.CreateAnimalRequest{
		AnimalData: animalData,
	}

	assert.Equal(t, animalData, req.AnimalData)
}

func TestMilkCollectionData(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	milkData := handlers.MilkCollectionData{
		ID:       1,
		AnimalID: 1,
		Liters:   25.5,
		Date:     date,
	}

	assert.Equal(t, uint(1), milkData.ID)
	assert.Equal(t, uint(1), milkData.AnimalID)
	assert.Equal(t, 25.5, milkData.Liters)
	assert.Equal(t, date, milkData.Date)
}

func TestCreateMilkCollectionRequest(t *testing.T) {
	req := handlers.CreateMilkCollectionRequest{
		AnimalID: 1,
		Liters:   25.5,
		Date:     "2024-01-15",
	}

	assert.Equal(t, uint(1), req.AnimalID)
	assert.Equal(t, 25.5, req.Liters)
	assert.Equal(t, "2024-01-15", req.Date)
}

func TestReproductionData(t *testing.T) {
	reproductionData := handlers.ReproductionData{
		ID:                     1,
		AnimalID:               1,
		AnimalName:             "Bella",
		EarTag:                 123,
		CurrentPhase:           1,
		InseminationType:       "IA",
		VeterinaryConfirmation: true,
		Observations:           "Primeira inseminação",
	}

	assert.Equal(t, uint(1), reproductionData.ID)
	assert.Equal(t, uint(1), reproductionData.AnimalID)
	assert.Equal(t, "Bella", reproductionData.AnimalName)
	assert.Equal(t, 1, reproductionData.CurrentPhase)
}

func TestCreateReproductionRequest(t *testing.T) {
	reproductionData := handlers.ReproductionData{
		ID:                     1,
		AnimalID:               1,
		CurrentPhase:           1,
		InseminationType:       "IA",
		VeterinaryConfirmation: true,
	}

	req := handlers.CreateReproductionRequest{
		ReproductionData: reproductionData,
	}

	assert.Equal(t, reproductionData, req.ReproductionData)
}

func TestUpdateReproductionPhaseRequest(t *testing.T) {
	additionalData := map[string]interface{}{
		"pregnancy_date": "2024-02-15",
	}

	req := handlers.UpdateReproductionPhaseRequest{
		AnimalID:       1,
		NewPhase:       2,
		AdditionalData: additionalData,
	}

	assert.Equal(t, uint(1), req.AnimalID)
	assert.Equal(t, 2, req.NewPhase)
	assert.Equal(t, additionalData, req.AdditionalData)
}
