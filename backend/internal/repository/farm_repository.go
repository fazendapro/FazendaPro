package repository

import (
	"github.com/fazendapro/FazendaPro-api/internal/models"
)

type FarmRepository struct {
	db *Database
}

func NewFarmRepository(db *Database) FarmRepositoryInterface {
	return &FarmRepository{db: db}
}

func (r *FarmRepository) FindByID(id uint) (*models.Farm, error) {
	var farm models.Farm
	err := r.db.DB.First(&farm, id).Error
	if err != nil {
		return nil, err
	}
	return &farm, nil
}

func (r *FarmRepository) Update(farm *models.Farm) error {
	return r.db.DB.Model(farm).Update("logo", farm.Logo).Error
}

func (r *FarmRepository) LoadCompanyData(farm *models.Farm) error {
	return r.db.DB.Preload("Company").First(farm, farm.ID).Error
}
