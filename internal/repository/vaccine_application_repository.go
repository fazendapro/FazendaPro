package repository

import (
	"errors"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"gorm.io/gorm"
)

type VaccineApplicationRepository struct {
	db *gorm.DB
}

func NewVaccineApplicationRepository(db *gorm.DB) *VaccineApplicationRepository {
	return &VaccineApplicationRepository{db: db}
}

func (r *VaccineApplicationRepository) Create(vaccineApplication *models.VaccineApplication) error {
	return r.db.Create(vaccineApplication).Error
}

func (r *VaccineApplicationRepository) FindByID(id uint) (*models.VaccineApplication, error) {
	var vaccineApplication models.VaccineApplication
	err := r.db.Preload("Animal").Preload("Vaccine").
		Where(SQLWhereID, id).
		First(&vaccineApplication).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &vaccineApplication, nil
}

func (r *VaccineApplicationRepository) FindByFarmID(farmID uint) ([]models.VaccineApplication, error) {
	var vaccineApplications []models.VaccineApplication
	err := r.db.Preload("Animal", SQLWhereFarmID, farmID).
		Preload("Vaccine", SQLWhereFarmID, farmID).
		Joins("JOIN animals ON vaccine_applications.animal_id = animals.id").
		Where(SQLWhereAnimalsFarmID, farmID).
		Order("vaccine_applications.application_date DESC").
		Find(&vaccineApplications).Error
	return vaccineApplications, err
}

func (r *VaccineApplicationRepository) FindByFarmIDWithDateRange(farmID uint, startDate, endDate *time.Time) ([]models.VaccineApplication, error) {
	var vaccineApplications []models.VaccineApplication
	query := r.db.Preload("Animal", SQLWhereFarmID, farmID).
		Preload("Vaccine", SQLWhereFarmID, farmID).
		Joins("JOIN animals ON vaccine_applications.animal_id = animals.id").
		Where(SQLWhereAnimalsFarmID, farmID)

	if startDate != nil {
		query = query.Where("vaccine_applications.application_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("vaccine_applications.application_date <= ?", *endDate)
	}

	err := query.Order("vaccine_applications.application_date DESC").Find(&vaccineApplications).Error
	return vaccineApplications, err
}

func (r *VaccineApplicationRepository) FindByAnimalID(animalID uint) ([]models.VaccineApplication, error) {
	var vaccineApplications []models.VaccineApplication
	err := r.db.Preload("Animal").Preload("Vaccine").
		Where(SQLWhereAnimalID, animalID).
		Order("application_date DESC").
		Find(&vaccineApplications).Error
	return vaccineApplications, err
}

func (r *VaccineApplicationRepository) FindByVaccineID(vaccineID uint) ([]models.VaccineApplication, error) {
	var vaccineApplications []models.VaccineApplication
	err := r.db.Preload("Animal").Preload("Vaccine").
		Where("vaccine_id = ?", vaccineID).
		Order("application_date DESC").
		Find(&vaccineApplications).Error
	return vaccineApplications, err
}

func (r *VaccineApplicationRepository) Update(vaccineApplication *models.VaccineApplication) error {
	return r.db.Save(vaccineApplication).Error
}

func (r *VaccineApplicationRepository) Delete(id uint) error {
	return r.db.Delete(&models.VaccineApplication{}, id).Error
}

