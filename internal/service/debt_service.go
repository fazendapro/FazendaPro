package service

import (
	"errors"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
)

type DebtService struct {
	repository repository.DebtRepositoryInterface
}

func NewDebtService(repository repository.DebtRepositoryInterface) *DebtService {
	return &DebtService{repository: repository}
}

func (s *DebtService) CreateDebt(debt *models.Debt) error {
	if debt.Person == "" {
		return errors.New("nome da pessoa é obrigatório")
	}

	if debt.Value <= 0 {
		return errors.New("valor deve ser maior que zero")
	}

	now := time.Now()
	debt.CreatedAt = now
	debt.UpdatedAt = now

	return s.repository.Create(debt)
}

func (s *DebtService) GetDebtByID(id uint) (*models.Debt, error) {
	if id == 0 {
		return nil, errors.New("ID é obrigatório")
	}

	return s.repository.FindByID(id)
}

func (s *DebtService) GetDebtsWithPagination(page, limit int, year, month *int) ([]models.Debt, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	return s.repository.FindAllWithPagination(page, limit, year, month)
}

func (s *DebtService) DeleteDebt(id uint) error {
	if id == 0 {
		return errors.New("ID é obrigatório")
	}

	_, err := s.repository.FindByID(id)
	if err != nil {
		return errors.New("dívida não encontrada")
	}

	return s.repository.Delete(id)
}

func (s *DebtService) GetTotalByPersonInMonth(year, month int) ([]repository.PersonTotal, error) {
	if year < 2000 || year > 3000 {
		return nil, errors.New("ano deve estar entre 2000 e 3000")
	}

	if month < 1 || month > 12 {
		return nil, errors.New("mês deve estar entre 1 e 12")
	}

	return s.repository.GetTotalByPersonInMonth(year, month)
}
