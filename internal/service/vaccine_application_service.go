package service

import (
	"errors"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type VaccineApplicationService struct {
	repository repository.VaccineApplicationRepositoryInterface
}

func NewVaccineApplicationService(repository repository.VaccineApplicationRepositoryInterface) *VaccineApplicationService {
	return &VaccineApplicationService{repository: repository}
}

func (s *VaccineApplicationService) CreateApplication(vaccineApplication *models.VaccineApplication) error {
	if vaccineApplication.AnimalID == 0 {
		return errors.New("ID do animal é obrigatório")
	}

	if vaccineApplication.VaccineID == 0 {
		return errors.New("ID da vacina é obrigatório")
	}

	if vaccineApplication.ApplicationDate.IsZero() {
		return errors.New("data de aplicação é obrigatória")
	}

	now := time.Now()
	vaccineApplication.CreatedAt = now
	vaccineApplication.UpdatedAt = now

	return s.repository.Create(vaccineApplication)
}

func (s *VaccineApplicationService) GetApplicationByID(id uint) (*models.VaccineApplication, error) {
	if id == 0 {
		return nil, errors.New("ID é obrigatório")
	}
	return s.repository.FindByID(id)
}

func (s *VaccineApplicationService) GetApplicationsByFarmID(farmID uint) ([]models.VaccineApplication, error) {
	if farmID == 0 {
		return nil, errors.New("ID da fazenda é obrigatório")
	}
	return s.repository.FindByFarmID(farmID)
}

func (s *VaccineApplicationService) GetApplicationsByFarmIDWithDateRange(farmID uint, startDate, endDate *time.Time) ([]models.VaccineApplication, error) {
	if farmID == 0 {
		return nil, errors.New("ID da fazenda é obrigatório")
	}
	return s.repository.FindByFarmIDWithDateRange(farmID, startDate, endDate)
}

func (s *VaccineApplicationService) GetApplicationsByFarmIDWithPagination(farmID uint, page, limit int) ([]models.VaccineApplication, int64, error) {
	if farmID == 0 {
		return nil, 0, errors.New("ID da fazenda é obrigatório")
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return s.repository.FindByFarmIDWithPagination(farmID, page, limit)
}

func (s *VaccineApplicationService) GetApplicationsByFarmIDWithDateRangePaginated(farmID uint, startDate, endDate *time.Time, page, limit int) ([]models.VaccineApplication, int64, error) {
	if farmID == 0 {
		return nil, 0, errors.New("ID da fazenda é obrigatório")
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return s.repository.FindByFarmIDWithDateRangePaginated(farmID, startDate, endDate, page, limit)
}

func (s *VaccineApplicationService) GetApplicationsByAnimalID(animalID uint) ([]models.VaccineApplication, error) {
	if animalID == 0 {
		return nil, errors.New("ID do animal é obrigatório")
	}
	return s.repository.FindByAnimalID(animalID)
}

func (s *VaccineApplicationService) GetApplicationsByVaccineID(vaccineID uint) ([]models.VaccineApplication, error) {
	if vaccineID == 0 {
		return nil, errors.New("ID da vacina é obrigatório")
	}
	return s.repository.FindByVaccineID(vaccineID)
}

func (s *VaccineApplicationService) UpdateApplication(vaccineApplication *models.VaccineApplication) error {
	if vaccineApplication.ID == 0 {
		return errors.New("ID da aplicação é obrigatório para atualização")
	}

	if vaccineApplication.AnimalID == 0 {
		return errors.New("ID do animal é obrigatório")
	}

	if vaccineApplication.VaccineID == 0 {
		return errors.New("ID da vacina é obrigatório")
	}

	if vaccineApplication.ApplicationDate.IsZero() {
		return errors.New("data de aplicação é obrigatória")
	}

	existingApplication, err := s.repository.FindByID(vaccineApplication.ID)
	if err != nil {
		return err
	}

	if existingApplication == nil {
		return errors.New("aplicação de vacina não encontrada")
	}

	vaccineApplication.UpdatedAt = time.Now()

	return s.repository.Update(vaccineApplication)
}

func (s *VaccineApplicationService) DeleteApplication(id uint) error {
	if id == 0 {
		return errors.New("ID é obrigatório")
	}

	existingApplication, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if existingApplication == nil {
		return errors.New("aplicação de vacina não encontrada")
	}

	return s.repository.Delete(id)
}
