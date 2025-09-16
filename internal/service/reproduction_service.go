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

	// Verificar se já existe um registro de reprodução para este animal
	existingReproduction, err := s.repository.FindByAnimalID(reproduction.AnimalID)
	if err != nil {
		return err
	}

	if existingReproduction != nil {
		return errors.New("já existe um registro de reprodução para este animal")
	}

	// Definir valores padrão
	if reproduction.CurrentPhase == 0 {
		reproduction.CurrentPhase = models.PhaseVazias // Vazias por padrão
	}

	// Definir timestamps
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

	// Verificar se o registro existe
	existingReproduction, err := s.repository.FindByID(reproduction.ID)
	if err != nil {
		return err
	}

	if existingReproduction == nil {
		return errors.New("registro de reprodução não encontrado")
	}

	// Atualizar timestamp
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

	// Atualizar a fase
	reproduction.CurrentPhase = newPhase
	now := time.Now()

	// Atualizar datas específicas baseadas na nova fase
	switch newPhase {
	case models.PhasePrenhas:
		if pregnancyDate, ok := additionalData["pregnancy_date"].(time.Time); ok {
			reproduction.PregnancyDate = &pregnancyDate
			// Calcular data prevista do parto (aproximadamente 280 dias)
			expectedBirth := pregnancyDate.AddDate(0, 0, 280)
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
		// Resetar algumas datas quando voltar para vazias
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
	// Verificar se o registro existe
	existingReproduction, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}

	if existingReproduction == nil {
		return errors.New("registro de reprodução não encontrado")
	}

	return s.repository.Delete(id)
}

