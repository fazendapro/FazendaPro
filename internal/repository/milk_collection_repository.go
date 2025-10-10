package repository

import (
	"fmt"
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
	fmt.Printf("DEBUG: Repository FindByID - Looking for ID: %d\n", id)

	err := r.db.Preload("Animal").Where("id = ?", id).First(&milkCollection).Error
	if err != nil {
		fmt.Printf("DEBUG: Repository FindByID Error: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Repository FindByID - Found: ID=%d, Liters=%.2f\n",
		milkCollection.ID, milkCollection.Liters)
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
	fmt.Printf("DEBUG: Repository Update - ID: %d, AnimalID: %d, Liters: %.2f\n",
		milkCollection.ID, milkCollection.AnimalID, milkCollection.Liters)

	var existingMilkCollection models.MilkCollection
	err := r.db.Where("id = ?", milkCollection.ID).First(&existingMilkCollection).Error
	if err != nil {
		fmt.Printf("DEBUG: Record not found with ID %d: %v\n", milkCollection.ID, err)
		return err
	}

	fmt.Printf("DEBUG: Found existing record - ID: %d, Liters: %.2f\n",
		existingMilkCollection.ID, existingMilkCollection.Liters)

	result := r.db.Model(&models.MilkCollection{}).Where("id = ?", milkCollection.ID).Updates(map[string]interface{}{
		"animal_id": milkCollection.AnimalID,
		"liters":    milkCollection.Liters,
		"date":      milkCollection.Date,
	})

	if result.Error != nil {
		fmt.Printf("DEBUG: Repository Update Error: %v\n", result.Error)
		return result.Error
	}

	fmt.Printf("DEBUG: Repository Update - Rows affected: %d\n", result.RowsAffected)
	return nil
}

func (r *MilkCollectionRepository) Delete(id uint) error {
	return r.db.Delete(&models.MilkCollection{}, id).Error
}
