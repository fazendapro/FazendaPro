package repository

import (
	"context"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"

	"gorm.io/gorm"
)

type SaleRepository interface {
	Create(ctx context.Context, sale *models.Sale) error
	GetByID(ctx context.Context, id uint) (*models.Sale, error)
	GetByFarmID(ctx context.Context, farmID uint) ([]*models.Sale, error)
	GetByAnimalID(ctx context.Context, animalID uint) ([]*models.Sale, error)
	GetByDateRange(ctx context.Context, farmID uint, startDate, endDate time.Time) ([]*models.Sale, error)
	Update(ctx context.Context, sale *models.Sale) error
	Delete(ctx context.Context, id uint) error
}

type saleRepository struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) SaleRepository {
	return &saleRepository{db: db}
}

func (r *saleRepository) Create(ctx context.Context, sale *models.Sale) error {
	return r.db.WithContext(ctx).Create(sale).Error
}

func (r *saleRepository) GetByID(ctx context.Context, id uint) (*models.Sale, error) {
	var sale models.Sale
	err := r.db.WithContext(ctx).Preload("Animal").Preload("Farm").First(&sale, id).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *saleRepository) GetByFarmID(ctx context.Context, farmID uint) ([]*models.Sale, error) {
	var sales []*models.Sale
	err := r.db.WithContext(ctx).Preload("Animal").Where("farm_id = ?", farmID).Order("sale_date DESC").Find(&sales).Error
	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *saleRepository) GetByAnimalID(ctx context.Context, animalID uint) ([]*models.Sale, error) {
	var sales []*models.Sale
	err := r.db.WithContext(ctx).Preload("Animal").Where("animal_id = ?", animalID).Order("sale_date DESC").Find(&sales).Error
	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *saleRepository) GetByDateRange(ctx context.Context, farmID uint, startDate, endDate time.Time) ([]*models.Sale, error) {
	var sales []*models.Sale
	err := r.db.WithContext(ctx).Preload("Animal").Where("farm_id = ? AND sale_date BETWEEN ? AND ?", farmID, startDate, endDate).Order("sale_date DESC").Find(&sales).Error
	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *saleRepository) Update(ctx context.Context, sale *models.Sale) error {
	return r.db.WithContext(ctx).Save(sale).Error
}

func (r *saleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Sale{}, id).Error
}
