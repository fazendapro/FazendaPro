package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUserRequest representa a requisição de criação de usuário
// @Description Dados necessários para criar um novo usuário
type CreateUserRequest struct {
	User   models.User   `json:"user"`   // Dados do usuário
	Person models.Person `json:"person"` // Dados da pessoa
}

// GetUser obtém um usuário por email
// @Summary      Obter usuário por email
// @Description  Retorna os dados de um usuário específico pelo email
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        email query string true "Email do usuário"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/users [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	user, err := h.service.GetUserByEmail(email)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		SendErrorResponse(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	SendSuccessResponse(w, user, "Usuário encontrado com sucesso", http.StatusOK)
}

// CreateUser cria um novo usuário
// @Summary      Criar usuário
// @Description  Cria um novo usuário no sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateUserRequest true "Dados do novo usuário"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(&req.User, &req.Person); err != nil {
		SendErrorResponse(w, "Erro ao criar usuário: "+err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"id":        req.User.ID,
		"person_id": req.User.PersonID,
		"farm_id":   req.User.FarmID,
	}
	SendSuccessResponse(w, data, "Usuário criado com sucesso", http.StatusCreated)
}

func (h *UserHandler) GetUserWithPerson(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		SendErrorResponse(w, "ID do usuário é obrigatório", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserWithPerson(1)
	if err != nil {
		SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		SendErrorResponse(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	SendSuccessResponse(w, user, "Usuário encontrado com sucesso", http.StatusOK)
}

func (h *UserHandler) UpdatePersonData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		SendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var personData models.Person
	if err := json.NewDecoder(r.Body).Decode(&personData); err != nil {
		SendErrorResponse(w, "Erro ao decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	userID := uint(1)

	if err := h.service.UpdatePersonData(userID, &personData); err != nil {
		SendErrorResponse(w, "Erro ao atualizar dados da pessoa: "+err.Error(), http.StatusBadRequest)
		return
	}

	SendSuccessResponse(w, nil, "Dados da pessoa atualizados com sucesso", http.StatusOK)
}
