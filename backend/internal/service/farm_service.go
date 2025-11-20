package service

import (
	"errors"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type FarmService struct {
	repository repository.FarmRepositoryInterface
}

func NewFarmService(repository repository.FarmRepositoryInterface) *FarmService {
	return &FarmService{repository: repository}
}

func (s *FarmService) GetFarmByID(farmID uint) (*models.Farm, error) {
	if farmID == 0 {
		return nil, errors.New("farm ID is required")
	}

	return s.repository.FindByID(farmID)
}

func (s *FarmService) UpdateFarm(farm *models.Farm) error {
	if farm.ID == 0 {
		return errors.New("farm ID is required")
	}

	return s.repository.Update(farm)
}

func (s *FarmService) LoadCompanyData(farm *models.Farm) error {
	return s.repository.LoadCompanyData(farm)
}
