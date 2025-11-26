package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
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
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	FarmID      uint   `json:"farm_id"`
	AccessToken string `json:"access_token"`
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

	w.Header().Set(HeaderContentType, ContentTypeJSON)
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

	user, err := h.service.GetUserWithPerson(userID)
	if err != nil {
		SendErrorResponse(w, "Erro ao buscar usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		SendErrorResponse(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	userForToken := *user
	userForToken.FarmID = farm.ID

	accessToken, err := h.generateJWT(&userForToken)
	if err != nil {
		SendErrorResponse(w, "Erro ao gerar token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := SelectFarmResponse{
		Success:     true,
		Message:     "Fazenda selecionada com sucesso",
		FarmID:      farm.ID,
		AccessToken: accessToken,
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
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

func (h *FarmSelectionHandler) generateJWT(user *models.User) (string, error) {
	if user.Person == nil {
		return "", fmt.Errorf("user person is nil")
	}

	claims := jwt.MapClaims{
		"sub":     user.ID,
		"email":   user.Person.Email,
		"farm_id": user.FarmID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
