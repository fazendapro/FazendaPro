package service

import (
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type MilkCollectionService struct {
	repository   repository.MilkCollectionRepositoryInterface
	batchService *BatchService
}

func NewMilkCollectionService(repository repository.MilkCollectionRepositoryInterface, batchService *BatchService) *MilkCollectionService {
	return &MilkCollectionService{
		repository:   repository,
		batchService: batchService,
	}
}

func (s *MilkCollectionService) CreateMilkCollection(milkCollection *models.MilkCollection) error {
	err := s.repository.Create(milkCollection)
	if err != nil {
		return err
	}

	err = s.batchService.UpdateAnimalBatch(milkCollection.AnimalID)
	if err != nil {
		return err
	}

	return nil
}

func (s *MilkCollectionService) GetMilkCollectionByID(id uint) (*models.MilkCollection, error) {
	return s.repository.FindByID(id)
}

func (s *MilkCollectionService) GetMilkCollectionsByFarmID(farmID uint) ([]models.MilkCollection, error) {
	return s.repository.FindByFarmID(farmID)
}

func (s *MilkCollectionService) GetMilkCollectionsByFarmIDWithPagination(farmID uint, page, limit int) ([]models.MilkCollection, int64, error) {
	return s.repository.FindByFarmIDWithPagination(farmID, page, limit)
}

func (s *MilkCollectionService) GetMilkCollectionsByFarmIDWithDateRange(farmID uint, startDate, endDate *time.Time) ([]models.MilkCollection, error) {
	return s.repository.FindByFarmIDWithDateRange(farmID, startDate, endDate)
}

func (s *MilkCollectionService) GetMilkCollectionsByFarmIDWithDateRangePaginated(farmID uint, startDate, endDate *time.Time, page, limit int) ([]models.MilkCollection, int64, error) {
	return s.repository.FindByFarmIDWithDateRangePaginated(farmID, startDate, endDate, page, limit)
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
