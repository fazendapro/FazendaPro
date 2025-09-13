package service

import (
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type MilkCollectionService struct {
	repository repository.MilkCollectionRepositoryInterface
}

func NewMilkCollectionService(repository repository.MilkCollectionRepositoryInterface) *MilkCollectionService {
	return &MilkCollectionService{repository: repository}
}

func (s *MilkCollectionService) CreateMilkCollection(milkCollection *models.MilkCollection) error {
	return s.repository.Create(milkCollection)
}

func (s *MilkCollectionService) GetMilkCollectionByID(id uint) (*models.MilkCollection, error) {
	return s.repository.FindByID(id)
}

func (s *MilkCollectionService) GetMilkCollectionsByFarmID(farmID uint) ([]models.MilkCollection, error) {
	return s.repository.FindByFarmID(farmID)
}

func (s *MilkCollectionService) GetMilkCollectionsByFarmIDWithDateRange(farmID uint, startDate, endDate *time.Time) ([]models.MilkCollection, error) {
	return s.repository.FindByFarmIDWithDateRange(farmID, startDate, endDate)
}

func (s *MilkCollectionService) GetMilkCollectionsByAnimalID(animalID uint) ([]models.MilkCollection, error) {
	return s.repository.FindByAnimalID(animalID)
}

func (s *MilkCollectionService) UpdateMilkCollection(milkCollection *models.MilkCollection) error {
	return s.repository.Update(milkCollection)
}

func (s *MilkCollectionService) DeleteMilkCollection(id uint) error {
	return s.repository.Delete(id)
}
