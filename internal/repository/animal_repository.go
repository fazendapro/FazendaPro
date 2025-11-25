package repository

import (
	"fmt"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"gorm.io/gorm"
)

type AnimalRepository struct {
	db *Database
}

func NewAnimalRepository(db *Database) AnimalRepositoryInterface {
	return &AnimalRepository{db: db}
}

func (r *AnimalRepository) Create(animal *models.Animal) error {
	if err := r.db.DB.Create(animal).Error; err != nil {
		return fmt.Errorf("erro ao criar animal: %w", err)
	}
	return nil
}

func (r *AnimalRepository) FindByID(id uint) (*models.Animal, error) {
	var animal models.Animal
	if err := r.db.DB.Preload("Father").Preload("Mother").Where(SQLWhereID, id).First(&animal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar animal: %w", err)
	}
	return &animal, nil
}

func (r *AnimalRepository) FindByFarmID(farmID uint) ([]models.Animal, error) {
	var animals []models.Animal
	if err := r.db.DB.Where(SQLWhereFarmID, farmID).Find(&animals).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar animais da fazenda: %w", err)
	}
	return animals, nil
}

func (r *AnimalRepository) FindByEarTagNumber(farmID uint, earTagNumber int) (*models.Animal, error) {
	var animal models.Animal
	if err := r.db.DB.Where(SQLWhereFarmIDAndEarTag, farmID, earTagNumber).First(&animal).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar animal pela brinca: %w", err)
	}
	return &animal, nil
}

func (r *AnimalRepository) Update(animal *models.Animal) error {
	if err := r.db.DB.Save(animal).Error; err != nil {
		return fmt.Errorf("erro ao atualizar animal: %w", err)
	}
	return nil
}

func (r *AnimalRepository) Delete(id uint) error {
	if err := r.db.DB.Delete(&models.Animal{}, id).Error; err != nil {
		return fmt.Errorf("erro ao deletar animal: %w", err)
	}
	return nil
}

func (r *AnimalRepository) FindByFarmIDAndSex(farmID uint, sex int) ([]models.Animal, error) {
	var animals []models.Animal
	if err := r.db.DB.Where(SQLWhereFarmIDAndSex, farmID, sex).Find(&animals).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar animais por sexo: %w", err)
	}
	return animals, nil
}

func (r *AnimalRepository) CountBySex(farmID uint, sex int) (int64, error) {
	var count int64
	if err := r.db.DB.Model(&models.Animal{}).Where(SQLWhereFarmIDAndSex, farmID, sex).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("erro ao contar animais por sexo: %w", err)
	}
	return count, nil
}
