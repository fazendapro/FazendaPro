package repository

import (
	"errors"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"gorm.io/gorm"
)

type ReproductionRepository struct {
	db *gorm.DB
}

func NewReproductionRepository(db *gorm.DB) *ReproductionRepository {
	return &ReproductionRepository{db: db}
}

func (r *ReproductionRepository) Create(reproduction *models.Reproduction) error {
	return r.db.Create(reproduction).Error
}

func (r *ReproductionRepository) FindByID(id uint) (*models.Reproduction, error) {
	var reproduction models.Reproduction
	err := r.db.Preload("Animal").First(&reproduction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reproduction, nil
}

func (r *ReproductionRepository) FindByAnimalID(animalID uint) (*models.Reproduction, error) {
	var reproduction models.Reproduction
	err := r.db.Preload("Animal").Where(SQLWhereAnimalID, animalID).First(&reproduction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reproduction, nil
}

func (r *ReproductionRepository) FindByFarmID(farmID uint) ([]models.Reproduction, error) {
	var reproductions []models.Reproduction
	err := r.db.Preload("Animal").
		Joins("JOIN animals ON reproductions.animal_id = animals.id").
		Where(SQLWhereAnimalsFarmID, farmID).
		Find(&reproductions).Error
	return reproductions, err
}

func (r *ReproductionRepository) FindByPhase(phase models.ReproductionPhase) ([]models.Reproduction, error) {
	var reproductions []models.Reproduction
	err := r.db.Preload("Animal").Where("current_phase = ?", phase).Find(&reproductions).Error
	return reproductions, err
}

func (r *ReproductionRepository) Update(reproduction *models.Reproduction) error {
	return r.db.Save(reproduction).Error
}

func (r *ReproductionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Reproduction{}, id).Error
}
