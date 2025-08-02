package service

import (
	"errors"
	"strings"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	if !strings.Contains(email, "@") {
		return nil, errors.New("email inválido")
	}

	return s.repository.FindByEmail(email)
}
