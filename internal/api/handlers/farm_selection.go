package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

type FarmSelectionHandler struct {
	service   *service.UserService
	jwtSecret string
}

func NewFarmSelectionHandler(service *service.UserService, jwtSecret string) *FarmSelectionHandler {
	return &FarmSelectionHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

type GetUserFarmsResponse struct {
	Success        bool        `json:"success"`
	Message        string      `json:"message"`
	Farms          interface{} `json:"farms"`
	AutoSelect     bool        `json:"auto_select"`
	SelectedFarmID *uint       `json:"selected_farm_id,omitempty"`
}

type SelectFarmRequest struct {
	FarmID uint `json:"farm_id"`
}

type SelectFarmResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	FarmID  uint   `json:"farm_id"`
}

func (h *FarmSelectionHandler) GetUserFarms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.extractUserIDFromToken(r)
	if err != nil {
		SendErrorResponse(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	farms, err := h.service.GetUserFarms(userID)
	if err != nil {
		SendErrorResponse(w, "Erro ao buscar fazendas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	shouldAutoSelect, err := h.service.ShouldAutoSelectFarm(userID)
	if err != nil {
		SendErrorResponse(w, "Erro ao verificar auto-seleção: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := GetUserFarmsResponse{
		Success:    true,
		Message:    "Fazendas recuperadas com sucesso",
		Farms:      farms,
		AutoSelect: shouldAutoSelect,
	}

	if shouldAutoSelect && len(farms) > 0 {
		response.SelectedFarmID = &farms[0].ID
	}

	fmt.Printf("DEBUG: Retornando %d fazendas, auto_select: %t\n", len(farms), shouldAutoSelect)
	for i, farm := range farms {
		fmt.Printf("DEBUG: Farm %d - ID: %d, CompanyID: %d, Logo: %s\n", i, farm.ID, farm.CompanyID, farm.Logo)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *FarmSelectionHandler) SelectFarm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.extractUserIDFromToken(r)
	if err != nil {
		SendErrorResponse(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	var req SelectFarmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	farm, err := h.service.GetUserFarmByID(userID, req.FarmID)
	if err != nil {
		SendErrorResponse(w, "Fazenda não encontrada ou não pertence ao usuário", http.StatusBadRequest)
		return
	}

	response := SelectFarmResponse{
		Success: true,
		Message: "Fazenda selecionada com sucesso",
		FarmID:  farm.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *FarmSelectionHandler) extractUserIDFromToken(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, jwt.ErrTokenMalformed
	}

	tokenString := authHeader[7:]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, jwt.ErrTokenMalformed
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, jwt.ErrTokenMalformed
	}

	return uint(userIDFloat), nil
}
