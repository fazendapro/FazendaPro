package repository

import (
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"gorm.io/gorm"
)

type MilkCollectionRepository struct {
	db *gorm.DB
}

func NewMilkCollectionRepository(db *gorm.DB) *MilkCollectionRepository {
	return &MilkCollectionRepository{db: db}
}

func (r *MilkCollectionRepository) Create(milkCollection *models.MilkCollection) error {
	return r.db.Create(milkCollection).Error
}

func (r *MilkCollectionRepository) FindByID(id uint) (*models.MilkCollection, error) {
	var milkCollection models.MilkCollection
	err := r.db.Preload("Animal").First(&milkCollection, id).Error
	if err != nil {
		return nil, err
	}
	return &milkCollection, nil
}

func (r *MilkCollectionRepository) FindByFarmID(farmID uint) ([]models.MilkCollection, error) {
	var milkCollections []models.MilkCollection
	err := r.db.Preload("Animal", "farm_id = ?", farmID).
		Joins("JOIN animals ON milk_collections.animal_id = animals.id").
		Where("animals.farm_id = ?", farmID).
		Order("milk_collections.date DESC").
		Find(&milkCollections).Error
	return milkCollections, err
}

func (r *MilkCollectionRepository) FindByFarmIDWithDateRange(farmID uint, startDate, endDate *time.Time) ([]models.MilkCollection, error) {
	var milkCollections []models.MilkCollection
	query := r.db.Preload("Animal", "farm_id = ?", farmID).
		Joins("JOIN animals ON milk_collections.animal_id = animals.id").
		Where("animals.farm_id = ?", farmID)

	if startDate != nil {
		query = query.Where("milk_collections.date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("milk_collections.date <= ?", *endDate)
	}

	err := query.Order("milk_collections.date DESC").Find(&milkCollections).Error
	return milkCollections, err
}

func (r *MilkCollectionRepository) FindByAnimalID(animalID uint) ([]models.MilkCollection, error) {
	var milkCollections []models.MilkCollection
	err := r.db.Preload("Animal").
		Where("animal_id = ?", animalID).
		Order("date DESC").
		Find(&milkCollections).Error
	return milkCollections, err
}

func (r *MilkCollectionRepository) Update(milkCollection *models.MilkCollection) error {
	return r.db.Save(milkCollection).Error
}

func (r *MilkCollectionRepository) Delete(id uint) error {
	return r.db.Delete(&models.MilkCollection{}, id).Error
}
