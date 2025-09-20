package repository

import (
	"fmt"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *Database
}

func NewRefreshTokenRepository(db *Database) RefreshTokenRepositoryInterface {
	return &RefreshTokenRepository{db: db}
}

type RefreshTokenRepositoryInterface interface {
	Create(userID uint, expiresAt time.Time) (*models.RefreshToken, error)
	FindByToken(token string) (*models.RefreshToken, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID uint) error
	DeleteExpired() error
}

func (r *RefreshTokenRepository) Create(userID uint, expiresAt time.Time) (*models.RefreshToken, error) {
	token := &models.RefreshToken{
		Token:     uuid.New().String(),
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	if err := r.db.DB.Create(token).Error; err != nil {
		return nil, fmt.Errorf("error creating refresh token: %w", err)
	}

	return token, nil
}

func (r *RefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.DB.Preload("User").Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding refresh token: %w", err)
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepository) DeleteByToken(token string) error {
	if err := r.db.DB.Where("token = ?", token).Delete(&models.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("error deleting refresh token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteByUserID(userID uint) error {
	if err := r.db.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("error deleting user refresh tokens: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteExpired() error {
	if err := r.db.DB.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("error deleting expired refresh tokens: %w", err)
	}
	return nil
}
