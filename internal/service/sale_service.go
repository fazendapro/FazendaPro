package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/cache"
	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type SaleService interface {
	CreateSale(ctx context.Context, sale *models.Sale) error
	GetSaleByID(ctx context.Context, id uint, farmID uint) (*models.Sale, error)
	GetSalesByFarmID(ctx context.Context, farmID uint) ([]*models.Sale, error)
	GetSalesByAnimalID(ctx context.Context, animalID uint, farmID uint) ([]*models.Sale, error)
	GetSalesByDateRange(ctx context.Context, farmID uint, startDate, endDate time.Time) ([]*models.Sale, error)
	GetMonthlySalesCount(ctx context.Context, farmID uint, startDate, endDate time.Time) (int64, error)
	GetMonthlySalesData(ctx context.Context, farmID uint, months int) ([]repository.MonthlySalesData, error)
	GetOverviewStats(ctx context.Context, farmID uint) (*repository.OverviewStats, error)
	UpdateSale(ctx context.Context, sale *models.Sale, farmID uint) error
	DeleteSale(ctx context.Context, id uint, farmID uint) error
	GetSalesHistory(ctx context.Context, farmID uint) ([]*models.Sale, error)
}

type saleService struct {
	saleRepo   repository.SaleRepository
	animalRepo repository.AnimalRepositoryInterface
	cache      cache.CacheInterface
}

func NewSaleService(saleRepo repository.SaleRepository, animalRepo repository.AnimalRepositoryInterface, cacheClient cache.CacheInterface) SaleService {
	return &saleService{
		saleRepo:   saleRepo,
		animalRepo: animalRepo,
		cache:      cacheClient,
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
		log.Printf("Erro ao buscar animal ID %d: %v", sale.AnimalID, err)
		return errors.New(ErrAnimalNotFound)
	}
	if animal == nil {
		log.Printf("Animal ID %d não encontrado", sale.AnimalID)
		return errors.New(ErrAnimalNotFound)
	}
	log.Printf("Animal ID %d encontrado - FarmID do animal: %d, FarmID da venda: %d", sale.AnimalID, animal.FarmID, sale.FarmID)
	if animal.FarmID != sale.FarmID {
		log.Printf("Erro: Animal ID %d pertence à fazenda %d, mas a venda está sendo criada para a fazenda %d", sale.AnimalID, animal.FarmID, sale.FarmID)
		return fmt.Errorf("animal ID %d pertence à fazenda %d, mas a venda está sendo criada para a fazenda %d", sale.AnimalID, animal.FarmID, sale.FarmID)
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
		s.saleRepo.Delete(ctx, sale.ID, sale.FarmID)
		return errors.New("failed to update animal status")
	}

	s.invalidateDashboardCache(sale.FarmID)

	return nil
}

func (s *saleService) GetSaleByID(ctx context.Context, id uint, farmID uint) (*models.Sale, error) {
	return s.saleRepo.GetByID(ctx, id, farmID)
}

func (s *saleService) GetSalesByFarmID(ctx context.Context, farmID uint) ([]*models.Sale, error) {
	return s.saleRepo.GetByFarmID(ctx, farmID)
}

func (s *saleService) GetSalesByAnimalID(ctx context.Context, animalID uint, farmID uint) ([]*models.Sale, error) {
	animal, err := s.animalRepo.FindByID(animalID)
	if err != nil {
		return nil, errors.New(ErrAnimalNotFound)
	}
	if animal == nil {
		return nil, errors.New(ErrAnimalNotFound)
	}
	if animal.FarmID != farmID {
		return nil, errors.New("animal does not belong to the specified farm")
	}

	return s.saleRepo.GetByAnimalID(ctx, animalID, farmID)
}

func (s *saleService) GetSalesByDateRange(ctx context.Context, farmID uint, startDate, endDate time.Time) ([]*models.Sale, error) {
	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}
	return s.saleRepo.GetByDateRange(ctx, farmID, startDate, endDate)
}

func (s *saleService) GetMonthlySalesCount(ctx context.Context, farmID uint, startDate, endDate time.Time) (int64, error) {
	if startDate.After(endDate) {
		return 0, errors.New("start date cannot be after end date")
	}
	return s.saleRepo.GetMonthlySalesCount(ctx, farmID, startDate, endDate)
}

func (s *saleService) GetMonthlySalesData(ctx context.Context, farmID uint, months int) ([]repository.MonthlySalesData, error) {
	if months <= 0 {
		months = 12
	}
	if months > 24 {
		months = 24
	}

	cacheKey := fmt.Sprintf("dashboard:monthly:%d:%d", farmID, months)
	var cachedData []repository.MonthlySalesData

	err := s.cache.Get(cacheKey, &cachedData)
	if err == nil {
		log.Printf("Cache HIT para dados mensais da fazenda %d (meses: %d)", farmID, months)
		return cachedData, nil
	}

	log.Printf("Cache MISS para dados mensais da fazenda %d (meses: %d)", farmID, months)
	data, err := s.saleRepo.GetMonthlySalesData(ctx, farmID, months)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(cacheKey, data, 900); err != nil {
		log.Printf("Erro ao salvar no cache (não crítico): %v", err)
	}

	return data, nil
}

func (s *saleService) GetOverviewStats(ctx context.Context, farmID uint) (*repository.OverviewStats, error) {
	cacheKey := fmt.Sprintf("dashboard:overview:%d", farmID)
	var cachedStats repository.OverviewStats

	err := s.cache.Get(cacheKey, &cachedStats)
	if err == nil {
		log.Printf("Cache HIT para estatísticas gerais da fazenda %d", farmID)
		return &cachedStats, nil
	}

	log.Printf("Cache MISS para estatísticas gerais da fazenda %d", farmID)
	stats, err := s.saleRepo.GetOverviewStats(ctx, farmID)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(cacheKey, stats, 600); err != nil {
		log.Printf("Erro ao salvar no cache (não crítico): %v", err)
	}

	return stats, nil
}

func (s *saleService) UpdateSale(ctx context.Context, sale *models.Sale, farmID uint) error {
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

	existingSale, err := s.saleRepo.GetByID(ctx, sale.ID, farmID)
	if err != nil {
		return errors.New(ErrSaleNotFoundOrNotBelongsToFarm)
	}
	if existingSale == nil {
		return errors.New(ErrSaleNotFoundOrNotBelongsToFarm)
	}

	sale.FarmID = farmID
	sale.AnimalID = existingSale.AnimalID

	err = s.saleRepo.Update(ctx, sale)
	if err != nil {
		return err
	}

	s.invalidateDashboardCache(farmID)

	return nil
}

func (s *saleService) DeleteSale(ctx context.Context, id uint, farmID uint) error {
	sale, err := s.saleRepo.GetByID(ctx, id, farmID)
	if err != nil {
		return errors.New(ErrSaleNotFoundOrNotBelongsToFarm)
	}
	if sale == nil {
		return errors.New(ErrSaleNotFoundOrNotBelongsToFarm)
	}

	err = s.saleRepo.Delete(ctx, id, farmID)
	if err != nil {
		return err
	}

	animal, err := s.animalRepo.FindByID(sale.AnimalID)
	if err == nil && animal != nil {
		animal.Status = models.AnimalStatusActive
		s.animalRepo.Update(animal)
	}

	s.invalidateDashboardCache(farmID)

	return nil
}

func (s *saleService) GetSalesHistory(ctx context.Context, farmID uint) ([]*models.Sale, error) {
	return s.saleRepo.GetByFarmID(ctx, farmID)
}

func (s *saleService) invalidateDashboardCache(farmID uint) {
	overviewKey := fmt.Sprintf("dashboard:overview:%d", farmID)
	if err := s.cache.Delete(overviewKey); err != nil {
		log.Printf(ErrInvalidateCache, err)
	}

	for months := 6; months <= 24; months += 6 {
		monthlyKey := fmt.Sprintf("dashboard:monthly:%d:%d", farmID, months)
		if err := s.cache.Delete(monthlyKey); err != nil {
			log.Printf(ErrInvalidateCache, err)
		}
	}
}
