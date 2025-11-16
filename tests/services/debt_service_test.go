package services

import (
	"errors"
	"testing"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/fazendapro/FazendaPro-api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestDebtService_CreateDebt_Success(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	debt := &models.Debt{
		Person: "João Silva",
		Value:  1500.50,
	}

	mockRepo.On("Create", debt).Return(nil)

	err := debtService.CreateDebt(debt)

	assert.NoError(t, err)
	assert.NotZero(t, debt.CreatedAt)
	assert.NotZero(t, debt.UpdatedAt)
	mockRepo.AssertExpectations(t)
}

func TestDebtService_CreateDebt_EmptyPerson(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	debt := &models.Debt{
		Person: "",
		Value:  1500.50,
	}

	err := debtService.CreateDebt(debt)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nome da pessoa é obrigatório")
}

func TestDebtService_CreateDebt_InvalidValue(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	debt := &models.Debt{
		Person: "João Silva",
		Value:  0,
	}

	err := debtService.CreateDebt(debt)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "valor deve ser maior que zero")
}

func TestDebtService_CreateDebt_NegativeValue(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	debt := &models.Debt{
		Person: "João Silva",
		Value:  -100,
	}

	err := debtService.CreateDebt(debt)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "valor deve ser maior que zero")
}

func TestDebtService_GetDebtByID_Success(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	expectedDebt := &models.Debt{
		ID:     1,
		Person: "João Silva",
		Value:  1500.50,
	}

	mockRepo.On("FindByID", uint(1)).Return(expectedDebt, nil)

	debt, err := debtService.GetDebtByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, debt)
	assert.Equal(t, uint(1), debt.ID)
	mockRepo.AssertExpectations(t)
}

func TestDebtService_GetDebtByID_InvalidID(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	debt, err := debtService.GetDebtByID(0)

	assert.Error(t, err)
	assert.Nil(t, debt)
	assert.Contains(t, err.Error(), "ID é obrigatório")
}

func TestDebtService_GetDebtsWithPagination_Success(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	expectedDebts := []models.Debt{
		{ID: 1, Person: "João Silva", Value: 1500.50},
		{ID: 2, Person: "Maria Santos", Value: 2000.00},
	}

	mockRepo.On("FindAllWithPagination", 1, 10, (*int)(nil), (*int)(nil)).Return(expectedDebts, int64(2), nil)

	debts, total, err := debtService.GetDebtsWithPagination(1, 10, nil, nil)

	assert.NoError(t, err)
	assert.Len(t, debts, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestDebtService_GetDebtsWithPagination_DefaultValues(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	expectedDebts := []models.Debt{}

	mockRepo.On("FindAllWithPagination", 1, 10, (*int)(nil), (*int)(nil)).Return(expectedDebts, int64(0), nil)

	debts, total, err := debtService.GetDebtsWithPagination(0, 0, nil, nil)

	assert.NoError(t, err)
	assert.Len(t, debts, 0)
	assert.Equal(t, int64(0), total)
	mockRepo.AssertExpectations(t)
}

func TestDebtService_GetDebtsWithPagination_WithFilters(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	expectedDebts := []models.Debt{
		{ID: 1, Person: "João Silva", Value: 1500.50},
	}

	year := 2024
	month := 1
	mockRepo.On("FindAllWithPagination", 2, 20, &year, &month).Return(expectedDebts, int64(1), nil)

	debts, total, err := debtService.GetDebtsWithPagination(2, 20, &year, &month)

	assert.NoError(t, err)
	assert.Len(t, debts, 1)
	assert.Equal(t, int64(1), total)
	mockRepo.AssertExpectations(t)
}

func TestDebtService_DeleteDebt_Success(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	existingDebt := &models.Debt{ID: 1}

	mockRepo.On("FindByID", uint(1)).Return(existingDebt, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := debtService.DeleteDebt(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDebtService_DeleteDebt_InvalidID(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	err := debtService.DeleteDebt(0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ID é obrigatório")
}

func TestDebtService_DeleteDebt_NotFound(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

	err := debtService.DeleteDebt(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dívida não encontrada")
	mockRepo.AssertExpectations(t)
}

func TestDebtService_GetTotalByPersonInMonth_Success(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	expectedTotals := []repository.PersonTotal{
		{Person: "João Silva", Total: 1500.50},
		{Person: "Maria Santos", Total: 2000.00},
	}

	mockRepo.On("GetTotalByPersonInMonth", 2024, 1).Return(expectedTotals, nil)

	totals, err := debtService.GetTotalByPersonInMonth(2024, 1)

	assert.NoError(t, err)
	assert.Len(t, totals, 2)
	mockRepo.AssertExpectations(t)
}

func TestDebtService_GetTotalByPersonInMonth_InvalidYear(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	totals, err := debtService.GetTotalByPersonInMonth(1999, 1)

	assert.Error(t, err)
	assert.Nil(t, totals)
	assert.Contains(t, err.Error(), "ano deve estar entre 2000 e 3000")
}

func TestDebtService_GetTotalByPersonInMonth_InvalidMonth(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	totals, err := debtService.GetTotalByPersonInMonth(2024, 0)

	assert.Error(t, err)
	assert.Nil(t, totals)
	assert.Contains(t, err.Error(), "mês deve estar entre 1 e 12")
}

func TestDebtService_GetTotalByPersonInMonth_MonthTooHigh(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	totals, err := debtService.GetTotalByPersonInMonth(2024, 13)

	assert.Error(t, err)
	assert.Nil(t, totals)
	assert.Contains(t, err.Error(), "mês deve estar entre 1 e 12")
}

func TestDebtService_GetTotalByPersonInMonth_YearTooHigh(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	totals, err := debtService.GetTotalByPersonInMonth(3001, 1)

	assert.Error(t, err)
	assert.Nil(t, totals)
	assert.Contains(t, err.Error(), "ano deve estar entre 2000 e 3000")
}

func TestDebtService_GetTotalByPersonInMonth_RepositoryError(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	mockRepo.On("GetTotalByPersonInMonth", 2024, 1).Return(nil, errors.New("database error"))

	totals, err := debtService.GetTotalByPersonInMonth(2024, 1)

	assert.Error(t, err)
	assert.Nil(t, totals)
	mockRepo.AssertExpectations(t)
}

func TestDebtService_GetDebtsWithPagination_RepositoryError(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	mockRepo.On("FindAllWithPagination", 1, 10, (*int)(nil), (*int)(nil)).Return(nil, int64(0), errors.New("database error"))

	debts, total, err := debtService.GetDebtsWithPagination(1, 10, nil, nil)

	assert.Error(t, err)
	assert.Nil(t, debts)
	assert.Equal(t, int64(0), total)
	mockRepo.AssertExpectations(t)
}

func TestDebtService_GetDebtByID_RepositoryError(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))

	debt, err := debtService.GetDebtByID(1)

	assert.Error(t, err)
	assert.Nil(t, debt)
	mockRepo.AssertExpectations(t)
}

func TestDebtService_CreateDebt_RepositoryError(t *testing.T) {
	mockRepo := new(MockDebtRepository)
	debtService := service.NewDebtService(mockRepo)

	debt := &models.Debt{
		Person: "João Silva",
		Value:  1500.50,
	}

	mockRepo.On("Create", debt).Return(errors.New("database error"))

	err := debtService.CreateDebt(debt)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

