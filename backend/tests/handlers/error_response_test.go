package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/api/handlers"
	"github.com/stretchr/testify/assert"
)

func TestSendErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()

	handlers.SendErrorResponse(w, "Test error message", http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response handlers.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Bad Request", response.Error)
	assert.Equal(t, "Test error message", response.Message)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestSendErrorResponse_NotFound(t *testing.T) {
	w := httptest.NewRecorder()

	handlers.SendErrorResponse(w, "Not found", http.StatusNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response handlers.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Not Found", response.Error)
	assert.Equal(t, "Not found", response.Message)
	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestSendErrorResponse_InternalServerError(t *testing.T) {
	w := httptest.NewRecorder()

	handlers.SendErrorResponse(w, "Internal error", http.StatusInternalServerError)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response handlers.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Internal Server Error", response.Error)
	assert.Equal(t, "Internal error", response.Message)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestSendSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]interface{}{
		"id":   1,
		"name": "Test",
	}

	handlers.SendSuccessResponse(w, data, "Success message", http.StatusOK)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Success message", response["message"])
	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.NotNil(t, response["data"])
}

func TestSendSuccessResponse_Created(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]interface{}{
		"id": 1,
	}

	handlers.SendSuccessResponse(w, data, "Created successfully", http.StatusCreated)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Created successfully", response["message"])
	assert.Equal(t, float64(http.StatusCreated), response["code"])
}

func TestSendSuccessResponse_WithNilData(t *testing.T) {
	w := httptest.NewRecorder()

	handlers.SendSuccessResponse(w, nil, "Success with nil data", http.StatusOK)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Success with nil data", response["message"])
	assert.Nil(t, response["data"])
}

func TestSendSuccessResponse_WithArrayData(t *testing.T) {
	w := httptest.NewRecorder()

	data := []string{"item1", "item2", "item3"}

	handlers.SendSuccessResponse(w, data, "Array data", http.StatusOK)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])
}

