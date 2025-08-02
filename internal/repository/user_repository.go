package repository

import (
	"fmt"

	"github.com/fazendapro/FazendaPro-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *Database
}

func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	if err := r.db.DB.Create(user).Error; err != nil {
		return fmt.Errorf("erro ao criar usuário: %w", err)
	}
	return nil
}
