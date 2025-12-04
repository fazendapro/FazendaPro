package service

import (
	"errors"
	"log"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type WeightService struct {
	repository repository.WeightRepositoryInterface
}

func NewWeightService(repository repository.WeightRepositoryInterface) *WeightService {
	return &WeightService{repository: repository}
}

func (s *WeightService) CreateOrUpdateWeight(weight *models.Weight) error {
	log.Println("Creating or updating weight record", weight)

	if weight.AnimalID == 0 {
		return errors.New("ID do animal é obrigatório")
	}

	if weight.AnimalWeight <= 0 {
		return errors.New("peso do animal deve ser maior que zero")
	}

	existingWeight, err := s.repository.FindByAnimalID(weight.AnimalID)
	if err != nil {
		return err
	}

	now := time.Now()

	if existingWeight != nil {
		existingWeight.AnimalWeight = weight.AnimalWeight
		existingWeight.Date = weight.Date
		existingWeight.UpdatedAt = now
		return s.repository.Update(existingWeight)
	}

	weight.CreatedAt = now
	weight.UpdatedAt = now
	return s.repository.Create(weight)
}

func (s *WeightService) GetWeightByID(id uint) (*models.Weight, error) {
	return s.repository.FindByID(id)
}

func (s *WeightService) GetWeightByAnimalID(animalID uint) (*models.Weight, error) {
	return s.repository.FindByAnimalID(animalID)
}

func (s *WeightService) GetWeightsByFarmID(farmID uint) ([]models.Weight, error) {
	return s.repository.FindByFarmID(farmID)
}

func (s *WeightService) UpdateWeight(weight *models.Weight) error {
	if weight.ID == 0 {
		return errors.New("ID do registro de peso é obrigatório para atualização")
	}

	existingWeight, err := s.repository.FindByID(weight.ID)
	if err != nil {
		return err
	}

	if existingWeight == nil {
		return errors.New("registro de peso não encontrado")
	}

	if weight.AnimalWeight <= 0 {
		return errors.New("peso do animal deve ser maior que zero")
	}

	weight.UpdatedAt = time.Now()
	return s.repository.Update(weight)
}

func (s *WeightService) DeleteWeight(id uint) error {
	existingWeight, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if existingWeight == nil {
		return errors.New("registro de peso não encontrado")
	}

	return s.repository.Delete(id)
}
