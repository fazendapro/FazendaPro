package repository

import (
	"errors"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"gorm.io/gorm"
)

type WeightRepository struct {
	db *gorm.DB
}

func NewWeightRepository(db *gorm.DB) *WeightRepository {
	return &WeightRepository{db: db}
}

func (r *WeightRepository) Create(weight *models.Weight) error {
	return r.db.Create(weight).Error
}

func (r *WeightRepository) FindByID(id uint) (*models.Weight, error) {
	var weight models.Weight
	err := r.db.Preload("Animal").First(&weight, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &weight, nil
}

func (r *WeightRepository) FindByAnimalID(animalID uint) (*models.Weight, error) {
	var weight models.Weight
	err := r.db.Preload("Animal").
		Where(SQLWhereAnimalID, animalID).
		Order("date DESC").
		First(&weight).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &weight, nil
}

func (r *WeightRepository) FindByFarmID(farmID uint) ([]models.Weight, error) {
	var weights []models.Weight
	err := r.db.Preload("Animal").
		Joins("JOIN animals ON weights.animal_id = animals.id").
		Where(SQLWhereAnimalsFarmID, farmID).
		Order("weights.date DESC").
		Find(&weights).Error
	return weights, err
}

func (r *WeightRepository) Update(weight *models.Weight) error {
	return r.db.Save(weight).Error
}

func (r *WeightRepository) Delete(id uint) error {
	return r.db.Delete(&models.Weight{}, id).Error
}
