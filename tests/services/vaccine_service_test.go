package services

import (
	"errors"
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestVaccineService_CreateVaccine_Success(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	vaccine := &models.Vaccine{
		FarmID:       1,
		Name:         "Vacina Aftosa",
		Description:  "Vacina contra febre aftosa",
		Manufacturer: "Fabricante XYZ",
	}

	mockRepo.On("Create", vaccine).Return(nil)

	err := vaccineService.CreateVaccine(vaccine)

	assert.NoError(t, err)
	assert.NotZero(t, vaccine.CreatedAt)
	assert.NotZero(t, vaccine.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_CreateVaccine_EmptyFarmID(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	vaccine := &models.Vaccine{
		FarmID: 0,
		Name:   "Vacina Aftosa",
	}

	err := vaccineService.CreateVaccine(vaccine)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID da fazenda é obrigatório")
}

func TestVaccineService_CreateVaccine_EmptyName(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	vaccine := &models.Vaccine{
		FarmID: 1,
		Name:   "",
	}

	err := vaccineService.CreateVaccine(vaccine)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nome da vacina é obrigatório")
}

func TestVaccineService_GetVaccineByID_Success(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	expectedVaccine := &models.Vaccine{
		ID:     1,
		FarmID: 1,
		Name:   "Vacina Aftosa",
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedVaccine, nil)

	vaccine, err := vaccineService.GetVaccineByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, vaccine)
	assert.Equal(t, uint(1), vaccine.ID)
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_GetVaccineByID_InvalidID(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	vaccine, err := vaccineService.GetVaccineByID(0)

	assert.Error(t, err)
	assert.Nil(t, vaccine)
	assert.Contains(t, err.Error(), "ID é obrigatório")
}

func TestVaccineService_GetVaccinesByFarmID_Success(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	expectedVaccines := []models.Vaccine{
		{ID: 1, FarmID: 1, Name: "Vacina Aftosa"},
		{ID: 2, FarmID: 1, Name: "Vacina Brucelose"},
	}

	mockRepo.On("FindByFarmID", uint(1)).Return(expectedVaccines, nil)

	vaccines, err := vaccineService.GetVaccinesByFarmID(1)

	assert.NoError(t, err)
	assert.Len(t, vaccines, 2)
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_GetVaccinesByFarmID_InvalidFarmID(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	vaccines, err := vaccineService.GetVaccinesByFarmID(0)

	assert.Error(t, err)
	assert.Nil(t, vaccines)
	assert.Contains(t, err.Error(), "ID da fazenda é obrigatório")
}

func TestVaccineService_UpdateVaccine_Success(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	existingVaccine := &models.Vaccine{
		ID:     1,
		FarmID: 1,
		Name:   "Vacina Aftosa",
	}

	vaccine := &models.Vaccine{
		ID:     1,
		FarmID: 1,
		Name:   "Vacina Aftosa Atualizada",
	}

	mockRepo.On("FindByID", uint(1)).Return(existingVaccine, nil)
	mockRepo.On("Update", vaccine).Return(nil)

	err := vaccineService.UpdateVaccine(vaccine)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_UpdateVaccine_InvalidID(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	vaccine := &models.Vaccine{
		ID:     0,
		FarmID: 1,
		Name:   "Vacina Aftosa",
	}

	err := vaccineService.UpdateVaccine(vaccine)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID da vacina é obrigatório")
}

func TestVaccineService_UpdateVaccine_EmptyName(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	vaccine := &models.Vaccine{
		ID:     1,
		FarmID: 1,
		Name:   "",
	}

	err := vaccineService.UpdateVaccine(vaccine)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nome da vacina é obrigatório")
}

func TestVaccineService_UpdateVaccine_NotFound(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	vaccine := &models.Vaccine{
		ID:     1,
		FarmID: 1,
		Name:   "Vacina Aftosa",
	}

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	err := vaccineService.UpdateVaccine(vaccine)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vacina não encontrada")
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_UpdateVaccine_RepositoryFindError(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	vaccine := &models.Vaccine{
		ID:     1,
		FarmID: 1,
		Name:   "Vacina Aftosa",
	}

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	err := vaccineService.UpdateVaccine(vaccine)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_UpdateVaccine_RepositoryUpdateError(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	existingVaccine := &models.Vaccine{
		ID:     1,
		FarmID: 1,
		Name:   "Vacina Aftosa",
	}

	vaccine := &models.Vaccine{
		ID:     1,
		FarmID: 1,
		Name:   "Vacina Aftosa Atualizada",
	}

	mockRepo.On("FindByID", uint(1)).Return(existingVaccine, nil)
	mockRepo.On("Update", vaccine).Return(errors.New("database error"))

	err := vaccineService.UpdateVaccine(vaccine)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_DeleteVaccine_Success(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	existingVaccine := &models.Vaccine{ID: 1}

	mockRepo.On("FindByID", uint(1)).Return(existingVaccine, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := vaccineService.DeleteVaccine(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_DeleteVaccine_InvalidID(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	err := vaccineService.DeleteVaccine(0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID é obrigatório")
}

func TestVaccineService_DeleteVaccine_NotFound(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	err := vaccineService.DeleteVaccine(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vacina não encontrada")
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_DeleteVaccine_RepositoryFindError(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	err := vaccineService.DeleteVaccine(1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_DeleteVaccine_RepositoryDeleteError(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	existingVaccine := &models.Vaccine{ID: 1}

	mockRepo.On("FindByID", uint(1)).Return(existingVaccine, nil)
	mockRepo.On("Delete", uint(1)).Return(errors.New("database error"))

	err := vaccineService.DeleteVaccine(1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_CreateVaccine_RepositoryError(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	vaccine := &models.Vaccine{
		FarmID: 1,
		Name:   "Vacina Aftosa",
	}

	mockRepo.On("Create", vaccine).Return(errors.New("database error"))

	err := vaccineService.CreateVaccine(vaccine)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_GetVaccineByID_RepositoryError(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	vaccine, err := vaccineService.GetVaccineByID(1)

	assert.Error(t, err)
	assert.Nil(t, vaccine)
	mockRepo.AssertExpectations(t)
}

func TestVaccineService_GetVaccinesByFarmID_RepositoryError(t *testing.T) {
	mockRepo := new(MockVaccineRepository)
	vaccineService := service.NewVaccineService(mockRepo)

	mockRepo.On("FindByFarmID", uint(1)).Return(nil, errors.New("database error"))

	vaccines, err := vaccineService.GetVaccinesByFarmID(1)

	assert.Error(t, err)
	assert.Nil(t, vaccines)
	mockRepo.AssertExpectations(t)
}
