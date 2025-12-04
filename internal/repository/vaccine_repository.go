package repository

import (
	"errors"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"gorm.io/gorm"
)

type VaccineRepository struct {
	db *gorm.DB
}

func NewVaccineRepository(db *gorm.DB) *VaccineRepository {
	return &VaccineRepository{db: db}
}

func (r *VaccineRepository) Create(vaccine *models.Vaccine) error {
	return r.db.Create(vaccine).Error
}

func (r *VaccineRepository) FindByID(id uint) (*models.Vaccine, error) {
	var vaccine models.Vaccine
	err := r.db.First(&vaccine, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &vaccine, nil
}

func (r *VaccineRepository) FindByFarmID(farmID uint) ([]models.Vaccine, error) {
	var vaccines []models.Vaccine
	err := r.db.Where(SQLWhereFarmID, farmID).
		Order("name ASC").
		Find(&vaccines).Error
	return vaccines, err
}

func (r *VaccineRepository) Update(vaccine *models.Vaccine) error {
	return r.db.Save(vaccine).Error
}

func (r *VaccineRepository) Delete(id uint) error {
	return r.db.Delete(&models.Vaccine{}, id).Error
}

