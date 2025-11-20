package service

import (
	"errors"
	"log"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type ReproductionService struct {
	repository repository.ReproductionRepositoryInterface
}

func NewReproductionService(repository repository.ReproductionRepositoryInterface) *ReproductionService {
	return &ReproductionService{repository: repository}
}

func (s *ReproductionService) CreateReproduction(reproduction *models.Reproduction) error {
	log.Println("Creating reproduction record", reproduction)

	if reproduction.AnimalID == 0 {
		return errors.New("ID do animal é obrigatório")
	}

	existingReproduction, err := s.repository.FindByAnimalID(reproduction.AnimalID)
	if err != nil {
		return err
	}

	if existingReproduction != nil {
		return errors.New("já existe um registro de reprodução para este animal")
	}

	if reproduction.CurrentPhase == 0 {
		reproduction.CurrentPhase = models.PhaseVazias
	}

	now := time.Now()
	reproduction.CreatedAt = now
	reproduction.UpdatedAt = now

	return s.repository.Create(reproduction)
}

func (s *ReproductionService) GetReproductionByID(id uint) (*models.Reproduction, error) {
	return s.repository.FindByID(id)
}

func (s *ReproductionService) GetReproductionByAnimalID(animalID uint) (*models.Reproduction, error) {
	return s.repository.FindByAnimalID(animalID)
}

func (s *ReproductionService) GetReproductionsByFarmID(farmID uint) ([]models.Reproduction, error) {
	return s.repository.FindByFarmID(farmID)
}

func (s *ReproductionService) GetReproductionsByPhase(phase models.ReproductionPhase) ([]models.Reproduction, error) {
	return s.repository.FindByPhase(phase)
}

func (s *ReproductionService) UpdateReproduction(reproduction *models.Reproduction) error {
	if reproduction.ID == 0 {
		return errors.New("ID do registro de reprodução é obrigatório para atualização")
	}

	existingReproduction, err := s.repository.FindByID(reproduction.ID)
	if err != nil {
		return err
	}

	if existingReproduction == nil {
		return errors.New("registro de reprodução não encontrado")
	}

	reproduction.UpdatedAt = time.Now()

	return s.repository.Update(reproduction)
}

func (s *ReproductionService) UpdateReproductionPhase(animalID uint, newPhase models.ReproductionPhase, additionalData map[string]interface{}) error {
	reproduction, err := s.repository.FindByAnimalID(animalID)
	if err != nil {
		return err
	}

	if reproduction == nil {
		return errors.New("registro de reprodução não encontrado para este animal")
	}

	reproduction.CurrentPhase = newPhase
	now := time.Now()

	switch newPhase {
	case models.PhasePrenhas:
		if pregnancyDate, ok := additionalData["pregnancy_date"].(time.Time); ok {
			reproduction.PregnancyDate = &pregnancyDate
			expectedBirth := pregnancyDate.AddDate(0, 0, 283)
			reproduction.ExpectedBirthDate = &expectedBirth
		}
		if inseminationDate, ok := additionalData["insemination_date"].(time.Time); ok {
			reproduction.InseminationDate = &inseminationDate
		}
		if inseminationType, ok := additionalData["insemination_type"].(string); ok {
			reproduction.InseminationType = inseminationType
		}
		if veterinaryConfirmation, ok := additionalData["veterinary_confirmation"].(bool); ok {
			reproduction.VeterinaryConfirmation = veterinaryConfirmation
		}

	case models.PhaseLactacao:
		if lactationStartDate, ok := additionalData["lactation_start_date"].(time.Time); ok {
			reproduction.LactationStartDate = &lactationStartDate
		} else {
			reproduction.LactationStartDate = &now
		}
		if actualBirthDate, ok := additionalData["actual_birth_date"].(time.Time); ok {
			reproduction.ActualBirthDate = &actualBirthDate
		}

	case models.PhaseSecando:
		if dryPeriodStartDate, ok := additionalData["dry_period_start_date"].(time.Time); ok {
			reproduction.DryPeriodStartDate = &dryPeriodStartDate
		} else {
			reproduction.DryPeriodStartDate = &now
		}
		if lactationEndDate, ok := additionalData["lactation_end_date"].(time.Time); ok {
			reproduction.LactationEndDate = &lactationEndDate
		}

	case models.PhaseVazias:
		reproduction.PregnancyDate = nil
		reproduction.ExpectedBirthDate = nil
		reproduction.ActualBirthDate = nil
		reproduction.LactationStartDate = nil
		reproduction.LactationEndDate = nil
		reproduction.DryPeriodStartDate = nil
	}

	if observations, ok := additionalData["observations"].(string); ok {
		reproduction.Observations = observations
	}

	reproduction.UpdatedAt = now

	return s.repository.Update(reproduction)
}

func (s *ReproductionService) DeleteReproduction(id uint) error {
	existingReproduction, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if existingReproduction == nil {
		return errors.New("registro de reprodução não encontrado")
	}

	return s.repository.Delete(id)
}
