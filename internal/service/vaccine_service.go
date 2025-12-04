package service

import (
	"errors"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type VaccineService struct {
	repository repository.VaccineRepositoryInterface
}

func NewVaccineService(repository repository.VaccineRepositoryInterface) *VaccineService {
	return &VaccineService{repository: repository}
}

func (s *VaccineService) CreateVaccine(vaccine *models.Vaccine) error {
	if vaccine.FarmID == 0 {
		return errors.New("ID da fazenda é obrigatório")
	}

	if vaccine.Name == "" {
		return errors.New("nome da vacina é obrigatório")
	}

	now := time.Now()
	vaccine.CreatedAt = now
	vaccine.UpdatedAt = now

	return s.repository.Create(vaccine)
}

func (s *VaccineService) GetVaccineByID(id uint) (*models.Vaccine, error) {
	if id == 0 {
		return nil, errors.New("ID é obrigatório")
	}
	return s.repository.FindByID(id)
}

func (s *VaccineService) GetVaccinesByFarmID(farmID uint) ([]models.Vaccine, error) {
	if farmID == 0 {
		return nil, errors.New("ID da fazenda é obrigatório")
	}
	return s.repository.FindByFarmID(farmID)
}

func (s *VaccineService) UpdateVaccine(vaccine *models.Vaccine) error {
	if vaccine.ID == 0 {
		return errors.New("ID da vacina é obrigatório para atualização")
	}

	if vaccine.Name == "" {
		return errors.New("nome da vacina é obrigatório")
	}

	existingVaccine, err := s.repository.FindByID(vaccine.ID)
	if err != nil {
		return err
	}

	if existingVaccine == nil {
		return errors.New("vacina não encontrada")
	}

	vaccine.UpdatedAt = time.Now()

	return s.repository.Update(vaccine)
}

func (s *VaccineService) DeleteVaccine(id uint) error {
	if id == 0 {
		return errors.New("ID é obrigatório")
	}

	existingVaccine, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if existingVaccine == nil {
		return errors.New("vacina não encontrada")
	}

	return s.repository.Delete(id)
}

