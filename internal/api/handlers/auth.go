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

// LoginRequest representa a requisição de login
// @Description Dados necessários para autenticação
type LoginRequest struct {
	Email    string `json:"email" example:"usuario@example.com"` // Email do usuário
	Password string `json:"password" example:"senha123"`         // Senha do usuário
}

// RegisterRequest representa a requisição de registro
// @Description Dados necessários para criar um novo usuário
type RegisterRequest struct {
	User   models.User   `json:"user"`   // Dados do usuário
	Person models.Person `json:"person"` // Dados da pessoa
}

// RegisterResponse representa a resposta de registro
// @Description Resposta após criação de usuário com tokens de autenticação
type RegisterResponse struct {
	Success      bool   `json:"success" example:"true"`                                         // Indica se a operação foi bem-sucedida
	Message      string `json:"message" example:"Usuário criado e logado com sucesso"`          // Mensagem de resposta
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Token de acesso JWT
	RefreshToken string `json:"refresh_token" example:"abc123def456..."`                        // Token de refresh
	User         struct {
		ID    uint   `json:"id"`    // ID do usuário
		Email string `json:"email"` // Email do usuário
		Name  string `json:"name"`  // Nome completo do usuário
	} `json:"user"` // Dados do usuário criado
}

// LoginResponse representa a resposta de login
// @Description Resposta após autenticação bem-sucedida
type LoginResponse struct {
	Success      bool   `json:"success" example:"true"`                                         // Indica se a operação foi bem-sucedida
	Message      string `json:"message" example:"Login realizado com sucesso"`                  // Mensagem de resposta
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Token de acesso JWT
	RefreshToken string `json:"refresh_token" example:"abc123def456..."`                        // Token de refresh
	User         struct {
		ID    uint   `json:"id"`    // ID do usuário
		Email string `json:"email"` // Email do usuário
		Name  string `json:"name"`  // Nome completo do usuário
	} `json:"user"` // Dados do usuário autenticado
}

// Login autentica um usuário e retorna tokens JWT
// @Summary      Realizar login
// @Description  Autentica um usuário com email e senha, retornando tokens de acesso e refresh
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Credenciais de login"
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/auth/login [post]
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

// Register cria um novo usuário e retorna tokens JWT
// @Summary      Registrar novo usuário
// @Description  Cria um novo usuário no sistema e retorna tokens de autenticação
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Dados do novo usuário"
// @Success      201  {object}  RegisterResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/auth/register [post]
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

// RefreshTokenRequest representa a requisição de refresh token
// @Description Token de refresh para renovar o access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"abc123def456..."` // Token de refresh
}

// RefreshTokenResponse representa a resposta de refresh token
// @Description Novo access token gerado
type RefreshTokenResponse struct {
	Success     bool   `json:"success" example:"true"`                                         // Indica se a operação foi bem-sucedida
	Message     string `json:"message" example:"Token renovado com sucesso"`                   // Mensagem de resposta
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // Novo token de acesso JWT
}

// RefreshToken renova o access token usando o refresh token
// @Summary      Renovar token de acesso
// @Description  Gera um novo access token usando um refresh token válido
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "Refresh token"
// @Success      200  {object}  RefreshTokenResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/auth/refresh [post]
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

// Logout invalida o refresh token do usuário
// @Summary      Realizar logout
// @Description  Invalida o refresh token, efetivando o logout do usuário
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "Refresh token a ser invalidado"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/auth/logout [post]
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
