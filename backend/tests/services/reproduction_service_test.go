package services

import (
	"errors"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReproductionService_GetReproductionsByPhase(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200)

	mockReproductions := []models.Reproduction{
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                1,
				FarmID:            1,
				AnimalName:        "Tata Salt",
				EarTagNumberLocal: 123,
				Photo:             "src/assets/images/mocked/cows/tata.png",
			},
		},
		{
			ID:            2,
			AnimalID:      2,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                2,
				FarmID:            1,
				AnimalName:        "Lays",
				EarTagNumberLocal: 124,
				Photo:             "src/assets/images/mocked/cows/lays.png",
			},
		},
	}

	mockRepo.On("FindByPhase", models.PhasePrenhas).Return(mockReproductions, nil)

	result, err := reproductionService.GetReproductionsByPhase(models.PhasePrenhas)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, models.PhasePrenhas, result[0].CurrentPhase)
	assert.Equal(t, models.PhasePrenhas, result[1].CurrentPhase)
	assert.Equal(t, "Tata Salt", result[0].Animal.AnimalName)
	assert.Equal(t, "Lays", result[1].Animal.AnimalName)

	mockRepo.AssertExpectations(t)
}

func TestReproductionService_GetReproductionsByPhase_Error(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	mockRepo.On("FindByPhase", models.PhasePrenhas).Return(nil, errors.New("database error"))

	result, err := reproductionService.GetReproductionsByPhase(models.PhasePrenhas)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")

	mockRepo.AssertExpectations(t)
}

func TestReproductionService_GetReproductionsByPhase_EmptyResults(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	mockRepo.On("FindByPhase", models.PhasePrenhas).Return([]models.Reproduction{}, nil)

	result, err := reproductionService.GetReproductionsByPhase(models.PhasePrenhas)

	assert.NoError(t, err)
	assert.Len(t, result, 0)

	mockRepo.AssertExpectations(t)
}

func TestReproductionService_GetReproductionsByPhase_DifferentPhases(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	mockReproductions := []models.Reproduction{
		{
			ID:           1,
			AnimalID:     1,
			CurrentPhase: models.PhaseLactacao,
			Animal: models.Animal{
				ID:                1,
				FarmID:            1,
				AnimalName:        "Tata Salt",
				EarTagNumberLocal: 123,
			},
		},
	}

	mockRepo.On("FindByPhase", models.PhaseLactacao).Return(mockReproductions, nil)

	result, err := reproductionService.GetReproductionsByPhase(models.PhaseLactacao)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, models.PhaseLactacao, result[0].CurrentPhase)

	mockRepo.AssertExpectations(t)
}

func TestReproductionService_GetReproductionsByPhase_WithAnimalData(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200)

	mockReproductions := []models.Reproduction{
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                   1,
				FarmID:               1,
				AnimalName:           "Tata Salt",
				EarTagNumberLocal:    123,
				EarTagNumberRegister: 456,
				Photo:                "src/assets/images/mocked/cows/tata.png",
				Sex:                  1,
				Breed:                "Holandesa",
				Type:                 "vaca",
			},
		},
	}

	mockRepo.On("FindByPhase", models.PhasePrenhas).Return(mockReproductions, nil)

	result, err := reproductionService.GetReproductionsByPhase(models.PhasePrenhas)

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	reproduction := result[0]
	assert.Equal(t, uint(1), reproduction.AnimalID)
	assert.Equal(t, models.PhasePrenhas, reproduction.CurrentPhase)
	assert.NotNil(t, reproduction.PregnancyDate)
	assert.Equal(t, pregnancyDate, *reproduction.PregnancyDate)

	animal := reproduction.Animal
	assert.Equal(t, uint(1), animal.ID)
	assert.Equal(t, uint(1), animal.FarmID)
	assert.Equal(t, "Tata Salt", animal.AnimalName)
	assert.Equal(t, 123, animal.EarTagNumberLocal)
	assert.Equal(t, 456, animal.EarTagNumberRegister)
	assert.Equal(t, "src/assets/images/mocked/cows/tata.png", animal.Photo)
	assert.Equal(t, 1, animal.Sex)
	assert.Equal(t, "Holandesa", animal.Breed)
	assert.Equal(t, "vaca", animal.Type)

	mockRepo.AssertExpectations(t)
}

func TestReproductionService_GetReproductionsByPhase_WithMultipleFarms(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -200)

	mockReproductions := []models.Reproduction{
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                1,
				FarmID:            1,
				AnimalName:        "Tata Salt",
				EarTagNumberLocal: 123,
			},
		},
		{
			ID:            2,
			AnimalID:      2,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                2,
				FarmID:            2,
				AnimalName:        "Lays",
				EarTagNumberLocal: 124,
			},
		},
		{
			ID:            3,
			AnimalID:      3,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: &pregnancyDate,
			Animal: models.Animal{
				ID:                3,
				FarmID:            1,
				AnimalName:        "Matilda",
				EarTagNumberLocal: 125,
			},
		},
	}

	mockRepo.On("FindByPhase", models.PhasePrenhas).Return(mockReproductions, nil)

	result, err := reproductionService.GetReproductionsByPhase(models.PhasePrenhas)

	assert.NoError(t, err)
	assert.Len(t, result, 3)

	for _, reproduction := range result {
		assert.Equal(t, models.PhasePrenhas, reproduction.CurrentPhase)
		assert.NotNil(t, reproduction.PregnancyDate)
	}

	farm1Count := 0
	farm2Count := 0
	for _, reproduction := range result {
		if reproduction.Animal.FarmID == 1 {
			farm1Count++
		} else if reproduction.Animal.FarmID == 2 {
			farm2Count++
		}
	}
	assert.Equal(t, 2, farm1Count)
	assert.Equal(t, 1, farm2Count)

	mockRepo.AssertExpectations(t)
}

func TestReproductionService_GetReproductionsByPhase_WithNullPregnancyDate(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	mockReproductions := []models.Reproduction{
		{
			ID:            1,
			AnimalID:      1,
			CurrentPhase:  models.PhasePrenhas,
			PregnancyDate: nil,
			Animal: models.Animal{
				ID:                1,
				FarmID:            1,
				AnimalName:        "Tata Salt",
				EarTagNumberLocal: 123,
			},
		},
	}

	mockRepo.On("FindByPhase", models.PhasePrenhas).Return(mockReproductions, nil)

	result, err := reproductionService.GetReproductionsByPhase(models.PhasePrenhas)

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	reproduction := result[0]
	assert.Equal(t, models.PhasePrenhas, reproduction.CurrentPhase)
	assert.Nil(t, reproduction.PregnancyDate)

	mockRepo.AssertExpectations(t)
}

func TestReproductionService_CreateReproduction_Success(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	reproduction := &models.Reproduction{
		AnimalID:     1,
		CurrentPhase: models.PhaseVazias,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.Reproduction")).Return(nil)

	err := reproductionService.CreateReproduction(reproduction)

	assert.NoError(t, err)
	assert.Equal(t, models.PhaseVazias, reproduction.CurrentPhase)
	assert.NotZero(t, reproduction.CreatedAt)
	assert.NotZero(t, reproduction.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_CreateReproduction_ZeroAnimalID(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	reproduction := &models.Reproduction{
		AnimalID: 0,
	}

	err := reproductionService.CreateReproduction(reproduction)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID do animal é obrigatório")
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_CreateReproduction_AnimalAlreadyExists(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	existingReproduction := &models.Reproduction{
		ID:       1,
		AnimalID: 1,
	}

	reproduction := &models.Reproduction{
		AnimalID: 1,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(existingReproduction, nil)

	err := reproductionService.CreateReproduction(reproduction)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "já existe um registro de reprodução para este animal")
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_CreateReproduction_DefaultPhase(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	reproduction := &models.Reproduction{
		AnimalID:     1,
		CurrentPhase: 0,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.Reproduction")).Return(nil)

	err := reproductionService.CreateReproduction(reproduction)

	assert.NoError(t, err)
	assert.Equal(t, models.PhaseVazias, reproduction.CurrentPhase)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_GetReproductionByID_Success(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	expectedReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhasePrenhas,
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedReproduction, nil)

	result, err := reproductionService.GetReproductionByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedReproduction, result)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_GetReproductionByID_Error(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	result, err := reproductionService.GetReproductionByID(1)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_GetReproductionByAnimalID_Success(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	expectedReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhaseLactacao,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(expectedReproduction, nil)

	result, err := reproductionService.GetReproductionByAnimalID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedReproduction, result)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_GetReproductionsByFarmID_Success(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	expectedReproductions := []models.Reproduction{
		{
			ID:           1,
			AnimalID:     1,
			CurrentPhase: models.PhasePrenhas,
		},
		{
			ID:           2,
			AnimalID:     2,
			CurrentPhase: models.PhaseLactacao,
		},
	}

	mockRepo.On("FindByFarmID", uint(1)).Return(expectedReproductions, nil)

	result, err := reproductionService.GetReproductionsByFarmID(1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproduction_Success(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	existingReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhaseVazias,
	}

	updatedReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhasePrenhas,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(nil)

	err := reproductionService.UpdateReproduction(updatedReproduction)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproduction_ZeroID(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	reproduction := &models.Reproduction{
		ID: 0,
	}

	err := reproductionService.UpdateReproduction(reproduction)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID do registro de reprodução é obrigatório")
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproduction_NotFound(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	reproduction := &models.Reproduction{
		ID: 1,
	}

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	err := reproductionService.UpdateReproduction(reproduction)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "registro de reprodução não encontrado")
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproductionPhase_PhasePrenhas(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -100)
	inseminationDate := now.AddDate(0, 0, -110)

	existingReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhaseVazias,
	}

	additionalData := map[string]interface{}{
		"pregnancy_date":          pregnancyDate,
		"insemination_date":       inseminationDate,
		"insemination_type":       "IA",
		"veterinary_confirmation": true,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(nil).Run(func(args mock.Arguments) {
		reproduction := args.Get(0).(*models.Reproduction)
		assert.Equal(t, models.PhasePrenhas, reproduction.CurrentPhase)
		assert.NotNil(t, reproduction.PregnancyDate)
		assert.NotNil(t, reproduction.ExpectedBirthDate)
		assert.NotNil(t, reproduction.InseminationDate)
		assert.Equal(t, "IA", reproduction.InseminationType)
		assert.True(t, reproduction.VeterinaryConfirmation)
	})

	err := reproductionService.UpdateReproductionPhase(1, models.PhasePrenhas, additionalData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproductionPhase_PhaseLactacao(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	now := time.Now()
	lactationStartDate := now.AddDate(0, 0, -10)
	actualBirthDate := now.AddDate(0, 0, -10)

	existingReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhasePrenhas,
	}

	additionalData := map[string]interface{}{
		"lactation_start_date": lactationStartDate,
		"actual_birth_date":    actualBirthDate,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(nil).Run(func(args mock.Arguments) {
		reproduction := args.Get(0).(*models.Reproduction)
		assert.Equal(t, models.PhaseLactacao, reproduction.CurrentPhase)
		assert.NotNil(t, reproduction.LactationStartDate)
		assert.NotNil(t, reproduction.ActualBirthDate)
	})

	err := reproductionService.UpdateReproductionPhase(1, models.PhaseLactacao, additionalData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproductionPhase_PhaseLactacao_DefaultDate(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	existingReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhasePrenhas,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(nil).Run(func(args mock.Arguments) {
		reproduction := args.Get(0).(*models.Reproduction)
		assert.Equal(t, models.PhaseLactacao, reproduction.CurrentPhase)
		assert.NotNil(t, reproduction.LactationStartDate)
	})

	err := reproductionService.UpdateReproductionPhase(1, models.PhaseLactacao, map[string]interface{}{})

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproductionPhase_PhaseSecando(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	now := time.Now()
	dryPeriodStartDate := now.AddDate(0, 0, -5)
	lactationEndDate := now.AddDate(0, 0, -5)

	existingReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhaseLactacao,
	}

	additionalData := map[string]interface{}{
		"dry_period_start_date": dryPeriodStartDate,
		"lactation_end_date":     lactationEndDate,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(nil).Run(func(args mock.Arguments) {
		reproduction := args.Get(0).(*models.Reproduction)
		assert.Equal(t, models.PhaseSecando, reproduction.CurrentPhase)
		assert.NotNil(t, reproduction.DryPeriodStartDate)
		assert.NotNil(t, reproduction.LactationEndDate)
	})

	err := reproductionService.UpdateReproductionPhase(1, models.PhaseSecando, additionalData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproductionPhase_PhaseSecando_DefaultDate(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	existingReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhaseLactacao,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(nil).Run(func(args mock.Arguments) {
		reproduction := args.Get(0).(*models.Reproduction)
		assert.Equal(t, models.PhaseSecando, reproduction.CurrentPhase)
		assert.NotNil(t, reproduction.DryPeriodStartDate)
	})

	err := reproductionService.UpdateReproductionPhase(1, models.PhaseSecando, map[string]interface{}{})

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproductionPhase_PhaseVazias(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	now := time.Now()
	pregnancyDate := now.AddDate(0, 0, -100)

	existingReproduction := &models.Reproduction{
		ID:            1,
		AnimalID:      1,
		CurrentPhase:  models.PhaseSecando,
		PregnancyDate: &pregnancyDate,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(nil).Run(func(args mock.Arguments) {
		reproduction := args.Get(0).(*models.Reproduction)
		assert.Equal(t, models.PhaseVazias, reproduction.CurrentPhase)
		assert.Nil(t, reproduction.PregnancyDate)
		assert.Nil(t, reproduction.ExpectedBirthDate)
		assert.Nil(t, reproduction.ActualBirthDate)
		assert.Nil(t, reproduction.LactationStartDate)
		assert.Nil(t, reproduction.LactationEndDate)
		assert.Nil(t, reproduction.DryPeriodStartDate)
	})

	err := reproductionService.UpdateReproductionPhase(1, models.PhaseVazias, map[string]interface{}{})

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproductionPhase_WithObservations(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	existingReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhaseVazias,
	}

	additionalData := map[string]interface{}{
		"observations": "Observação de teste",
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(nil).Run(func(args mock.Arguments) {
		reproduction := args.Get(0).(*models.Reproduction)
		assert.Equal(t, "Observação de teste", reproduction.Observations)
	})

	err := reproductionService.UpdateReproductionPhase(1, models.PhasePrenhas, additionalData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproductionPhase_NotFound(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, nil)

	err := reproductionService.UpdateReproductionPhase(1, models.PhasePrenhas, map[string]interface{}{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "registro de reprodução não encontrado")
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_DeleteReproduction_Success(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	existingReproduction := &models.Reproduction{
		ID:       1,
		AnimalID: 1,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := reproductionService.DeleteReproduction(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_DeleteReproduction_NotFound(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	err := reproductionService.DeleteReproduction(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "registro de reprodução não encontrado")
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_DeleteReproduction_Error(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	err := reproductionService.DeleteReproduction(1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproduction_RepositoryError(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	reproduction := &models.Reproduction{
		ID: 1,
	}

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	err := reproductionService.UpdateReproduction(reproduction)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproduction_UpdateError(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	existingReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhaseVazias,
	}

	updatedReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhasePrenhas,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(errors.New("update error"))

	err := reproductionService.UpdateReproduction(updatedReproduction)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproductionPhase_RepositoryError(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, errors.New("database error"))

	err := reproductionService.UpdateReproductionPhase(1, models.PhasePrenhas, map[string]interface{}{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	mockRepo.AssertExpectations(t)
}

func TestReproductionService_UpdateReproductionPhase_UpdateError(t *testing.T) {
	mockRepo := &MockReproductionRepository{}
	reproductionService := service.NewReproductionService(mockRepo)

	existingReproduction := &models.Reproduction{
		ID:           1,
		AnimalID:     1,
		CurrentPhase: models.PhaseVazias,
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(existingReproduction, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Reproduction")).Return(errors.New("update error"))

	err := reproductionService.UpdateReproductionPhase(1, models.PhasePrenhas, map[string]interface{}{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
	mockRepo.AssertExpectations(t)
}
