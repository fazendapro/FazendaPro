package service

import (
	"context"
	"errors"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type SaleService interface {
	CreateSale(ctx context.Context, sale *models.Sale) error
	GetSaleByID(ctx context.Context, id uint) (*models.Sale, error)
	GetSalesByFarmID(ctx context.Context, farmID uint) ([]*models.Sale, error)
	GetSalesByAnimalID(ctx context.Context, animalID uint) ([]*models.Sale, error)
	GetSalesByDateRange(ctx context.Context, farmID uint, startDate, endDate time.Time) ([]*models.Sale, error)
	UpdateSale(ctx context.Context, sale *models.Sale) error
	DeleteSale(ctx context.Context, id uint) error
	GetSalesHistory(ctx context.Context, farmID uint) ([]*models.Sale, error)
}

type saleService struct {
	saleRepo   repository.SaleRepository
	animalRepo repository.AnimalRepositoryInterface
}

func NewSaleService(saleRepo repository.SaleRepository, animalRepo repository.AnimalRepositoryInterface) SaleService {
	return &saleService{
		saleRepo:   saleRepo,
		animalRepo: animalRepo,
	}
}

func (s *saleService) CreateSale(ctx context.Context, sale *models.Sale) error {
	if sale.AnimalID == 0 {
		return errors.New("animal ID is required")
	}
	if sale.FarmID == 0 {
		return errors.New("farm ID is required")
	}
	if sale.BuyerName == "" {
		return errors.New("buyer name is required")
	}
	if sale.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if sale.SaleDate.IsZero() {
		return errors.New("sale date is required")
	}

	animal, err := s.animalRepo.FindByID(sale.AnimalID)
	if err != nil {
		return errors.New("animal not found")
	}
	if animal.FarmID != sale.FarmID {
		return errors.New("animal does not belong to the specified farm")
	}

	if animal.Status == models.AnimalStatusSold {
		return errors.New("animal is already sold")
	}

	err = s.saleRepo.Create(ctx, sale)
	if err != nil {
		return err
	}

	animal.Status = models.AnimalStatusSold
	err = s.animalRepo.Update(animal)
	if err != nil {
		s.saleRepo.Delete(ctx, sale.ID)
		return errors.New("failed to update animal status")
	}

	return nil
}

func (s *saleService) GetSaleByID(ctx context.Context, id uint) (*models.Sale, error) {
	return s.saleRepo.GetByID(ctx, id)
}

func (s *saleService) GetSalesByFarmID(ctx context.Context, farmID uint) ([]*models.Sale, error) {
	return s.saleRepo.GetByFarmID(ctx, farmID)
}

func (s *saleService) GetSalesByAnimalID(ctx context.Context, animalID uint) ([]*models.Sale, error) {
	return s.saleRepo.GetByAnimalID(ctx, animalID)
}

func (s *saleService) GetSalesByDateRange(ctx context.Context, farmID uint, startDate, endDate time.Time) ([]*models.Sale, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}
	return s.saleRepo.GetByDateRange(ctx, farmID, startDate, endDate)
}

func (s *saleService) UpdateSale(ctx context.Context, sale *models.Sale) error {
	if sale.ID == 0 {
		return errors.New("sale ID is required")
	}
	if sale.BuyerName == "" {
		return errors.New("buyer name is required")
	}
	if sale.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if sale.SaleDate.IsZero() {
		return errors.New("sale date is required")
	}

	return s.saleRepo.Update(ctx, sale)
}

func (s *saleService) DeleteSale(ctx context.Context, id uint) error {
	sale, err := s.saleRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("sale not found")
	}

	err = s.saleRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	animal, err := s.animalRepo.FindByID(sale.AnimalID)
	if err == nil {
		animal.Status = models.AnimalStatusActive
		s.animalRepo.Update(animal)
	}

	return nil
}

func (s *saleService) GetSalesHistory(ctx context.Context, farmID uint) ([]*models.Sale, error) {
	return s.saleRepo.GetByFarmID(ctx, farmID)
}
