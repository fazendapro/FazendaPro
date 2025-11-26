package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	service          *service.UserService
	refreshTokenRepo repository.RefreshTokenRepositoryInterface
	jwtSecret        string
}

func NewAuthHandler(service *service.UserService, refreshTokenRepo repository.RefreshTokenRepositoryInterface, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		service:          service,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	User   models.User   `json:"user"`
	Person models.Person `json:"person"`
}

type RegisterResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
}

type LoginResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByEmail(req.Email)
	if err != nil {
		SendErrorResponse(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	if user == nil {
		SendErrorResponse(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	valid, err := h.service.ValidatePasswordByEmail(req.Email, req.Password)
	if err != nil {
		SendErrorResponse(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}
	if !valid {
		SendErrorResponse(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.generateJWT(user)
	if err != nil {
		SendErrorResponse(w, ErrGenerateToken, http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.refreshTokenRepo.Create(user.ID, time.Now().Add(time.Hour*24*7))
	if err != nil {
		SendErrorResponse(w, "Erro ao gerar refresh token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Success:      true,
		Message:      "Login realizado com sucesso",
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		User: struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name"`
		}{
			ID:    user.ID,
			Email: user.Person.Email,
			Name:  user.Person.FirstName + " " + user.Person.LastName,
		},
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(&req.User, &req.Person); err != nil {
		SendErrorResponse(w, "Erro ao criar usuário: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByEmail(req.Person.Email)
	if err != nil {
		SendErrorResponse(w, "Erro ao buscar usuário criado", http.StatusInternalServerError)
		return
	}

	if user == nil {
		SendErrorResponse(w, "Erro ao buscar usuário criado", http.StatusInternalServerError)
		return
	}

	accessToken, err := h.generateJWT(user)
	if err != nil {
		SendErrorResponse(w, ErrGenerateToken, http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.refreshTokenRepo.Create(user.ID, time.Now().Add(time.Hour*24*7))
	if err != nil {
		SendErrorResponse(w, "Erro ao gerar refresh token", http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{
		Success:      true,
		Message:      "Usuário criado e logado com sucesso",
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		User: struct {
			ID    uint   `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name"`
		}{
			ID:    user.ID,
			Email: user.Person.Email,
			Name:  user.Person.FirstName + " " + user.Person.LastName,
		},
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) generateJWT(user *models.User) (string, error) {
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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	refreshToken, err := h.refreshTokenRepo.FindByToken(req.RefreshToken)
	if err != nil {
		SendErrorResponse(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	if refreshToken == nil {
		SendErrorResponse(w, "Refresh token inválido ou expirado", http.StatusUnauthorized)
		return
	}

	accessToken, err := h.generateJWT(&refreshToken.User)
	if err != nil {
		SendErrorResponse(w, ErrGenerateToken, http.StatusInternalServerError)
		return
	}

	response := RefreshTokenResponse{
		Success:     true,
		Message:     "Token renovado com sucesso",
		AccessToken: accessToken,
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, ErrDecodeJSON+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.refreshTokenRepo.DeleteByToken(req.RefreshToken); err != nil {
		SendErrorResponse(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Logout realizado com sucesso",
	}

	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
