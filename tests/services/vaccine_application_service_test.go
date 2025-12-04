package services

import (
	"errors"
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestVaccineApplicationService_CreateApplication_Success(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: applicationDate,
		BatchNumber:     "LOTE123",
		Veterinarian:    "Dr. João Silva",
		Observations:    "Aplicação realizada com sucesso",
	}

	mockRepo.On("Create", vaccineApplication).Return(nil)

	err := vaccineApplicationService.CreateApplication(vaccineApplication)

	assert.NoError(t, err)
	assert.NotZero(t, vaccineApplication.CreatedAt)
	assert.NotZero(t, vaccineApplication.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_CreateApplication_EmptyAnimalID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		AnimalID:        0,
		VaccineID:       1,
		ApplicationDate: applicationDate,
	}

	err := vaccineApplicationService.CreateApplication(vaccineApplication)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID do animal é obrigatório")
}

func TestVaccineApplicationService_CreateApplication_EmptyVaccineID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		AnimalID:        1,
		VaccineID:       0,
		ApplicationDate: applicationDate,
	}

	err := vaccineApplicationService.CreateApplication(vaccineApplication)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID da vacina é obrigatório")
}

func TestVaccineApplicationService_CreateApplication_EmptyDate(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	vaccineApplication := &models.VaccineApplication{
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: time.Time{},
	}

	err := vaccineApplicationService.CreateApplication(vaccineApplication)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data de aplicação é obrigatória")
}

func TestVaccineApplicationService_GetApplicationByID_Success(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	expectedApplication := &models.VaccineApplication{
		ID:              1,
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: time.Now(),
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedApplication, nil)

	application, err := vaccineApplicationService.GetApplicationByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, application)
	assert.Equal(t, uint(1), application.ID)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationByID_InvalidID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	application, err := vaccineApplicationService.GetApplicationByID(0)

	assert.Error(t, err)
	assert.Nil(t, application)
	assert.Contains(t, err.Error(), "ID é obrigatório")
}

func TestVaccineApplicationService_GetApplicationsByFarmID_Success(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	expectedApplications := []models.VaccineApplication{
		{ID: 1, AnimalID: 1, VaccineID: 1, ApplicationDate: time.Now()},
		{ID: 2, AnimalID: 2, VaccineID: 1, ApplicationDate: time.Now()},
	}

	mockRepo.On("FindByFarmID", uint(1)).Return(expectedApplications, nil)

	applications, err := vaccineApplicationService.GetApplicationsByFarmID(1)

	assert.NoError(t, err)
	assert.Len(t, applications, 2)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationsByFarmID_InvalidFarmID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applications, err := vaccineApplicationService.GetApplicationsByFarmID(0)

	assert.Error(t, err)
	assert.Nil(t, applications)
	assert.Contains(t, err.Error(), "ID da fazenda é obrigatório")
}

func TestVaccineApplicationService_GetApplicationsByFarmIDWithDateRange_Success(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	expectedApplications := []models.VaccineApplication{
		{ID: 1, AnimalID: 1, VaccineID: 1, ApplicationDate: time.Now()},
	}

	mockRepo.On("FindByFarmIDWithDateRange", uint(1), &startDate, &endDate).Return(expectedApplications, nil)

	applications, err := vaccineApplicationService.GetApplicationsByFarmIDWithDateRange(1, &startDate, &endDate)

	assert.NoError(t, err)
	assert.Len(t, applications, 1)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationsByFarmIDWithDateRange_InvalidFarmID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	applications, err := vaccineApplicationService.GetApplicationsByFarmIDWithDateRange(0, &startDate, &endDate)

	assert.Error(t, err)
	assert.Nil(t, applications)
	assert.Contains(t, err.Error(), "ID da fazenda é obrigatório")
}

func TestVaccineApplicationService_GetApplicationsByFarmIDWithDateRange_RepositoryError(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	mockRepo.On("FindByFarmIDWithDateRange", uint(1), &startDate, &endDate).Return(nil, errors.New("database error"))

	applications, err := vaccineApplicationService.GetApplicationsByFarmIDWithDateRange(1, &startDate, &endDate)

	assert.Error(t, err)
	assert.Nil(t, applications)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationsByAnimalID_Success(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	expectedApplications := []models.VaccineApplication{
		{ID: 1, AnimalID: 1, VaccineID: 1, ApplicationDate: time.Now()},
	}

	mockRepo.On("FindByAnimalID", uint(1)).Return(expectedApplications, nil)

	applications, err := vaccineApplicationService.GetApplicationsByAnimalID(1)

	assert.NoError(t, err)
	assert.Len(t, applications, 1)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationsByAnimalID_InvalidAnimalID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applications, err := vaccineApplicationService.GetApplicationsByAnimalID(0)

	assert.Error(t, err)
	assert.Nil(t, applications)
	assert.Contains(t, err.Error(), "ID do animal é obrigatório")
}

func TestVaccineApplicationService_GetApplicationsByAnimalID_RepositoryError(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	mockRepo.On("FindByAnimalID", uint(1)).Return(nil, errors.New("database error"))

	applications, err := vaccineApplicationService.GetApplicationsByAnimalID(1)

	assert.Error(t, err)
	assert.Nil(t, applications)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationsByVaccineID_Success(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	expectedApplications := []models.VaccineApplication{
		{ID: 1, AnimalID: 1, VaccineID: 1, ApplicationDate: time.Now()},
	}

	mockRepo.On("FindByVaccineID", uint(1)).Return(expectedApplications, nil)

	applications, err := vaccineApplicationService.GetApplicationsByVaccineID(1)

	assert.NoError(t, err)
	assert.Len(t, applications, 1)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationsByVaccineID_InvalidVaccineID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applications, err := vaccineApplicationService.GetApplicationsByVaccineID(0)

	assert.Error(t, err)
	assert.Nil(t, applications)
	assert.Contains(t, err.Error(), "ID da vacina é obrigatório")
}

func TestVaccineApplicationService_GetApplicationsByVaccineID_RepositoryError(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	mockRepo.On("FindByVaccineID", uint(1)).Return(nil, errors.New("database error"))

	applications, err := vaccineApplicationService.GetApplicationsByVaccineID(1)

	assert.Error(t, err)
	assert.Nil(t, applications)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationsByFarmIDWithPagination_Success(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	expectedApplications := []models.VaccineApplication{
		{ID: 1, AnimalID: 1, VaccineID: 1, ApplicationDate: time.Now()},
		{ID: 2, AnimalID: 2, VaccineID: 1, ApplicationDate: time.Now()},
	}

	mockRepo.On("FindByFarmIDWithPagination", uint(1), 1, 10).Return(expectedApplications, int64(2), nil)

	applications, total, err := vaccineApplicationService.GetApplicationsByFarmIDWithPagination(1, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, applications, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationsByFarmIDWithPagination_InvalidFarmID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applications, total, err := vaccineApplicationService.GetApplicationsByFarmIDWithPagination(0, 1, 10)

	assert.Error(t, err)
	assert.Nil(t, applications)
	assert.Equal(t, int64(0), total)
	assert.Contains(t, err.Error(), "ID da fazenda é obrigatório")
}

func TestVaccineApplicationService_GetApplicationsByFarmIDWithPagination_InvalidPage(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	expectedApplications := []models.VaccineApplication{
		{ID: 1, AnimalID: 1, VaccineID: 1, ApplicationDate: time.Now()},
	}

	mockRepo.On("FindByFarmIDWithPagination", uint(1), 1, 10).Return(expectedApplications, int64(1), nil)

	applications, total, err := vaccineApplicationService.GetApplicationsByFarmIDWithPagination(1, 0, 10)

	assert.NoError(t, err)
	assert.Len(t, applications, 1)
	assert.Equal(t, int64(1), total)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationsByFarmIDWithDateRangePaginated_Success(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	startDate := time.Now().AddDate(0, 0, -30)
	endDate := time.Now()

	expectedApplications := []models.VaccineApplication{
		{ID: 1, AnimalID: 1, VaccineID: 1, ApplicationDate: time.Now()},
	}

	mockRepo.On("FindByFarmIDWithDateRangePaginated", uint(1), &startDate, &endDate, 1, 10).Return(expectedApplications, int64(1), nil)

	applications, total, err := vaccineApplicationService.GetApplicationsByFarmIDWithDateRangePaginated(1, &startDate, &endDate, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, applications, 1)
	assert.Equal(t, int64(1), total)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_UpdateApplication_Success(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	existingApplication := &models.VaccineApplication{
		ID:              1,
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: time.Now(),
	}

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		ID:              1,
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: applicationDate,
		BatchNumber:     "LOTE456",
	}

	mockRepo.On("FindByID", uint(1)).Return(existingApplication, nil)
	mockRepo.On("Update", vaccineApplication).Return(nil)

	err := vaccineApplicationService.UpdateApplication(vaccineApplication)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_UpdateApplication_InvalidID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		ID:              0,
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: applicationDate,
	}

	err := vaccineApplicationService.UpdateApplication(vaccineApplication)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID da aplicação é obrigatório")
}

func TestVaccineApplicationService_UpdateApplication_NotFound(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		ID:              1,
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: applicationDate,
	}

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	err := vaccineApplicationService.UpdateApplication(vaccineApplication)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "aplicação de vacina não encontrada")
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_UpdateApplication_EmptyAnimalID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		ID:              1,
		AnimalID:        0,
		VaccineID:       1,
		ApplicationDate: applicationDate,
	}

	err := vaccineApplicationService.UpdateApplication(vaccineApplication)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID do animal é obrigatório")
}

func TestVaccineApplicationService_UpdateApplication_EmptyVaccineID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		ID:              1,
		AnimalID:        1,
		VaccineID:       0,
		ApplicationDate: applicationDate,
	}

	err := vaccineApplicationService.UpdateApplication(vaccineApplication)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID da vacina é obrigatório")
}

func TestVaccineApplicationService_UpdateApplication_EmptyDate(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	vaccineApplication := &models.VaccineApplication{
		ID:              1,
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: time.Time{},
	}

	err := vaccineApplicationService.UpdateApplication(vaccineApplication)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data de aplicação é obrigatória")
}

func TestVaccineApplicationService_UpdateApplication_RepositoryFindError(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		ID:              1,
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: applicationDate,
	}

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	err := vaccineApplicationService.UpdateApplication(vaccineApplication)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_UpdateApplication_RepositoryUpdateError(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	existingApplication := &models.VaccineApplication{
		ID:              1,
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: time.Now(),
	}

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		ID:              1,
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: applicationDate,
	}

	mockRepo.On("FindByID", uint(1)).Return(existingApplication, nil)
	mockRepo.On("Update", vaccineApplication).Return(errors.New("database error"))

	err := vaccineApplicationService.UpdateApplication(vaccineApplication)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_DeleteApplication_Success(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	existingApplication := &models.VaccineApplication{ID: 1}

	mockRepo.On("FindByID", uint(1)).Return(existingApplication, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := vaccineApplicationService.DeleteApplication(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_DeleteApplication_InvalidID(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	err := vaccineApplicationService.DeleteApplication(0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID é obrigatório")
}

func TestVaccineApplicationService_DeleteApplication_NotFound(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, nil)

	err := vaccineApplicationService.DeleteApplication(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "aplicação de vacina não encontrada")
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_CreateApplication_RepositoryError(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	applicationDate := time.Now()
	vaccineApplication := &models.VaccineApplication{
		AnimalID:        1,
		VaccineID:       1,
		ApplicationDate: applicationDate,
	}

	mockRepo.On("Create", vaccineApplication).Return(errors.New("database error"))

	err := vaccineApplicationService.CreateApplication(vaccineApplication)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationByID_RepositoryError(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	application, err := vaccineApplicationService.GetApplicationByID(1)

	assert.Error(t, err)
	assert.Nil(t, application)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_GetApplicationsByFarmID_RepositoryError(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	mockRepo.On("FindByFarmID", uint(1)).Return(nil, errors.New("database error"))

	applications, err := vaccineApplicationService.GetApplicationsByFarmID(1)

	assert.Error(t, err)
	assert.Nil(t, applications)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_DeleteApplication_RepositoryFindError(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	err := vaccineApplicationService.DeleteApplication(1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVaccineApplicationService_DeleteApplication_RepositoryDeleteError(t *testing.T) {
	mockRepo := new(MockVaccineApplicationRepository)
	vaccineApplicationService := service.NewVaccineApplicationService(mockRepo)

	existingApplication := &models.VaccineApplication{ID: 1}

	mockRepo.On("FindByID", uint(1)).Return(existingApplication, nil)
	mockRepo.On("Delete", uint(1)).Return(errors.New("database error"))

	err := vaccineApplicationService.DeleteApplication(1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
