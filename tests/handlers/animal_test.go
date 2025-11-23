package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/fazendapro/FazendaPro-api/internal/cache"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/fazendapro/FazendaPro-api/tests/services"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAnimalRouter(mockRepo *services.MockAnimalRepository, mockCache *MockCache) (*chi.Mux, *services.MockAnimalRepository, *MockCache) {
	animalService := service.NewAnimalService(mockRepo, mockCache)
	animalHandler := handlers.NewAnimalHandler(animalService)
	r := chi.NewRouter()
	r.Post("/animals", animalHandler.CreateAnimal)
	r.Get("/animals", animalHandler.GetAnimal)
	r.Get("/animals/farm", animalHandler.GetAnimalsByFarm)
	r.Put("/animals", animalHandler.UpdateAnimal)
	r.Delete("/animals", animalHandler.DeleteAnimal)
	r.Get("/animals/sex", animalHandler.GetAnimalsBySex)
	r.Post("/animals/photo", animalHandler.UploadAnimalPhoto)
	return r, mockRepo, mockCache
}

func TestAnimalHandler_CreateAnimal_Success(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	birthDate := "2020-01-15"
	animalData := map[string]interface{}{
		"farm_id":                 1,
		"ear_tag_number_local":    123,
		"ear_tag_number_register": 456,
		"animal_name":             "Boi Teste",
		"sex":                     1,
		"breed":                   "Nelore",
		"type":                    "Bovino",
		"birth_date":              birthDate,
		"animal_type":             1,
		"status":                  0,
		"purpose":                 0,
		"current_batch":           1,
	}

	jsonData, _ := json.Marshal(animalData)
	req, _ := http.NewRequest("POST", "/animals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	mockRepo.On("FindByEarTagNumber", uint(1), 123).Return((*models.Animal)(nil), nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.Animal")).Return(nil).Run(func(args mock.Arguments) {
		animal := args.Get(0).(*models.Animal)
		animal.ID = 1
	})
	mockCache.On("Delete", "animals:farm:1").Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, float64(1), response["data"].(map[string]interface{})["id"])
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_CreateAnimal_InvalidMethod(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	animalService := service.NewAnimalService(mockRepo, new(MockCache))
	animalHandler := handlers.NewAnimalHandler(animalService)

	req, _ := http.NewRequest("GET", "/animals", nil)
	w := httptest.NewRecorder()

	animalHandler.CreateAnimal(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestAnimalHandler_CreateAnimal_InvalidJSON(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("POST", "/animals", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_CreateAnimal_ServiceError(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	animalData := map[string]interface{}{
		"farm_id":              1,
		"ear_tag_number_local": 123,
		"animal_name":          "Boi Teste",
		"breed":                "Nelore",
		"type":                 "Bovino",
	}

	jsonData, _ := json.Marshal(animalData)
	req, _ := http.NewRequest("POST", "/animals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("FindByEarTagNumber", uint(1), 123).Return((*models.Animal)(nil), errors.New("erro ao buscar"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_GetAnimal_Success(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	birthDate := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedAnimal := &models.Animal{
		ID:                1,
		FarmID:            1,
		EarTagNumberLocal: 123,
		AnimalName:        "Boi Teste",
		Sex:               1,
		Breed:             "Nelore",
		Type:              "Bovino",
		BirthDate:         &birthDate,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	req, _ := http.NewRequest("GET", "/animals?id=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(expectedAnimal, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_GetAnimal_MissingID(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("GET", "/animals", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_GetAnimal_InvalidID(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("GET", "/animals?id=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_GetAnimal_NotFound(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("GET", "/animals?id=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_GetAnimal_ServiceError(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("GET", "/animals?id=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("erro ao buscar animal"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_GetAnimalsByFarm_Success(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	expectedAnimals := []models.Animal{
		{ID: 1, FarmID: 1, AnimalName: "Animal 1"},
		{ID: 2, FarmID: 1, AnimalName: "Animal 2"},
	}

	req, _ := http.NewRequest("GET", "/animals/farm?farmId=1", nil)
	w := httptest.NewRecorder()

	mockCache.On("Get", "animals:farm:1", mock.Anything).Return(cache.ErrCacheMiss)
	mockRepo.On("FindByFarmID", uint(1)).Return(expectedAnimals, nil)
	mockCache.On("Set", "animals:farm:1", expectedAnimals, int32(300)).Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_GetAnimalsByFarm_MissingFarmID(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("GET", "/animals/farm", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_GetAnimalsByFarm_InvalidFarmID(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("GET", "/animals/farm?farmId=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_UpdateAnimal_Success(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	animalData := map[string]interface{}{
		"id":                   1,
		"farm_id":              1,
		"ear_tag_number_local": 123,
		"animal_name":          "Animal Atualizado",
		"breed":                "Nelore",
		"type":                 "Bovino",
	}

	jsonData, _ := json.Marshal(animalData)
	req, _ := http.NewRequest("PUT", "/animals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	updatedAnimal := &models.Animal{
		ID:                1,
		FarmID:            1,
		EarTagNumberLocal: 123,
		AnimalName:        "Animal Atualizado",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	existingAnimal := &models.Animal{ID: 1, FarmID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(nil)
	mockRepo.On("FindByID", uint(1)).Return(updatedAnimal, nil).Once()
	mockCache.On("Delete", "animals:farm:1").Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_UpdateAnimal_InvalidMethod(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	animalService := service.NewAnimalService(mockRepo, new(MockCache))
	animalHandler := handlers.NewAnimalHandler(animalService)

	req, _ := http.NewRequest("GET", "/animals", nil)
	w := httptest.NewRecorder()

	animalHandler.UpdateAnimal(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestAnimalHandler_DeleteAnimal_Success(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("DELETE", "/animals?id=1", nil)
	w := httptest.NewRecorder()

	existingAnimal := &models.Animal{ID: 1, FarmID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)
	mockCache.On("Delete", "animals:farm:1").Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_DeleteAnimal_MissingID(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("DELETE", "/animals", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_DeleteAnimal_InvalidID(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("DELETE", "/animals?id=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_DeleteAnimal_ServiceError(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("DELETE", "/animals?id=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("animal not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_UpdateAnimal_ServiceError(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	animalData := map[string]interface{}{
		"id":                   1,
		"farm_id":              1,
		"ear_tag_number_local": 123,
		"animal_name":          "Animal Atualizado",
	}

	jsonData, _ := json.Marshal(animalData)
	req, _ := http.NewRequest("PUT", "/animals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	existingAnimal := &models.Animal{ID: 1, FarmID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(errors.New("erro ao atualizar"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_UpdateAnimal_GetAnimalError(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	animalData := map[string]interface{}{
		"id":                   1,
		"farm_id":              1,
		"ear_tag_number_local": 123,
		"animal_name":          "Animal Atualizado",
	}

	jsonData, _ := json.Marshal(animalData)
	req, _ := http.NewRequest("PUT", "/animals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	existingAnimal := &models.Animal{ID: 1, FarmID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(nil)
	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("erro ao buscar"))
	mockCache.On("Delete", "animals:farm:1").Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_GetAnimalsByFarm_ServiceError(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("GET", "/animals/farm?farmId=1", nil)
	w := httptest.NewRecorder()

	mockCache.On("Get", "animals:farm:1", mock.Anything).Return(cache.ErrCacheMiss)
	mockRepo.On("FindByFarmID", uint(1)).Return(nil, errors.New("erro ao buscar"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_GetAnimalsBySex_ServiceError(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("GET", "/animals/sex?farmId=1&sex=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByFarmIDAndSex", uint(1), 1).Return(nil, errors.New("erro ao buscar"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_GetAnimalsBySex_Success(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	expectedAnimals := []models.Animal{
		{ID: 1, FarmID: 1, Sex: 1, AnimalName: "Macho 1"},
		{ID: 2, FarmID: 1, Sex: 1, AnimalName: "Macho 2"},
	}

	req, _ := http.NewRequest("GET", "/animals/sex?farmId=1&sex=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByFarmIDAndSex", uint(1), 1).Return(expectedAnimals, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_GetAnimalsBySex_MissingFarmID(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("GET", "/animals/sex?sex=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_GetAnimalsBySex_MissingSex(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("GET", "/animals/sex?farmId=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_GetAnimalsBySex_InvalidSex(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	req, _ := http.NewRequest("GET", "/animals/sex?farmId=1&sex=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_UploadAnimalPhoto_Success(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.WriteField("animal_id", "1")
	part, _ := writer.CreateFormFile("photo", "test.jpg")
	part.Write([]byte("fake image content"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/animals/photo", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	existingAnimal := &models.Animal{
		ID:         1,
		FarmID:     1,
		AnimalName: "Test Animal",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	updatedAnimal := &models.Animal{
		ID:         1,
		FarmID:     1,
		AnimalName: "Test Animal",
		Photo:      "data:image/jpeg;base64,ZmFrZSBpbWFnZSBjb250ZW50",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil).Once()
	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(nil)
	mockRepo.On("FindByID", uint(1)).Return(updatedAnimal, nil).Once()
	mockCache.On("Delete", "animals:farm:1").Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_UploadAnimalPhoto_InvalidMethod(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	animalService := service.NewAnimalService(mockRepo, new(MockCache))
	animalHandler := handlers.NewAnimalHandler(animalService)

	req, _ := http.NewRequest("GET", "/animals/photo", nil)
	w := httptest.NewRecorder()

	animalHandler.UploadAnimalPhoto(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestAnimalHandler_UploadAnimalPhoto_MissingAnimalID(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.CreateFormFile("photo", "test.jpg")
	writer.Close()

	req, _ := http.NewRequest("POST", "/animals/photo", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_UploadAnimalPhoto_InvalidAnimalID(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.WriteField("animal_id", "invalid")
	writer.Close()

	req, _ := http.NewRequest("POST", "/animals/photo", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnimalHandler_UploadAnimalPhoto_AnimalNotFound(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.WriteField("animal_id", "1")
	fileWriter, _ := writer.CreateFormFile("photo", "test.jpg")
	fileWriter.Write([]byte("fake image data"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/animals/photo", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("animal not found"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_UploadAnimalPhoto_UpdateError(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.WriteField("animal_id", "1")
	fileWriter, _ := writer.CreateFormFile("photo", "test.jpg")
	fileWriter.Write([]byte("fake image data"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/animals/photo", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	existingAnimal := &models.Animal{ID: 1, FarmID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil).Once()
	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(errors.New("erro ao atualizar"))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_UploadAnimalPhoto_GetUpdatedAnimalError(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.WriteField("animal_id", "1")
	fileWriter, _ := writer.CreateFormFile("photo", "test.jpg")
	fileWriter.Write([]byte("fake image data"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/animals/photo", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	existingAnimal := &models.Animal{ID: 1, FarmID: 1}
	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil).Once()
	mockRepo.On("FindByID", uint(1)).Return(existingAnimal, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*models.Animal")).Return(nil)
	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("erro ao buscar"))
	mockCache.On("Delete", "animals:farm:1").Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_CreateAnimal_WithFatherAndMother(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	fatherID := uint(10)
	motherID := uint(20)
	animalData := map[string]interface{}{
		"farm_id":                 1,
		"ear_tag_number_local":    123,
		"ear_tag_number_register": 456,
		"animal_name":             "Filhote Teste",
		"sex":                     0,
		"breed":                   "Holandesa",
		"type":                    "Bovino",
		"father_id":               fatherID,
		"mother_id":               motherID,
		"animal_type":             0,
		"purpose":                 1,
	}

	jsonData, _ := json.Marshal(animalData)
	req, _ := http.NewRequest("POST", "/animals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("FindByEarTagNumber", uint(1), 123).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.Animal")).Return(nil).Run(func(args mock.Arguments) {
		animal := args.Get(0).(*models.Animal)
		animal.ID = 1
	})
	mockCache.On("Delete", "animals:farm:1").Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_CreateAnimal_InvalidBirthDate(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	animalData := map[string]interface{}{
		"farm_id":                 1,
		"ear_tag_number_local":    123,
		"ear_tag_number_register": 456,
		"animal_name":             "Animal Teste",
		"sex":                     0,
		"breed":                   "Holandesa",
		"type":                    "Bovino",
		"birth_date":              "invalid-date",
		"animal_type":             0,
		"purpose":                 1,
	}

	jsonData, _ := json.Marshal(animalData)
	req, _ := http.NewRequest("POST", "/animals", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockRepo.On("FindByEarTagNumber", uint(1), 123).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.Animal")).Return(nil).Run(func(args mock.Arguments) {
		animal := args.Get(0).(*models.Animal)
		animal.ID = 1
	})
	mockCache.On("Delete", "animals:farm:1").Return(nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestAnimalHandler_GetAnimal_WithFatherAndMother(t *testing.T) {
	mockRepo := new(services.MockAnimalRepository)
	mockCache := new(MockCache)
	router, _, _ := setupAnimalRouter(mockRepo, mockCache)

	fatherID := uint(10)
	motherID := uint(20)
	animal := &models.Animal{
		ID:                1,
		FarmID:            1,
		AnimalName:        "Filhote",
		EarTagNumberLocal: 123,
		Father: &models.Animal{
			ID:                fatherID,
			AnimalName:        "Pai",
			EarTagNumberLocal: 10,
		},
		Mother: &models.Animal{
			ID:                motherID,
			AnimalName:        "Mãe",
			EarTagNumberLocal: 20,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req, _ := http.NewRequest("GET", "/animals?id=1", nil)
	w := httptest.NewRecorder()

	mockRepo.On("FindByID", uint(1)).Return(animal, nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["father"])
	assert.NotNil(t, data["mother"])
	mockRepo.AssertExpectations(t)
}
