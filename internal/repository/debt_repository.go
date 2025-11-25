package repository

import (
	"fmt"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"gorm.io/gorm"
)

type PersonTotal struct {
	Person string  `json:"person"`
	Total  float64 `json:"total"`
}

type DebtRepository struct {
	db *gorm.DB
}

func NewDebtRepository(db *gorm.DB) DebtRepositoryInterface {
	return &DebtRepository{db: db}
}

func (r *DebtRepository) Create(debt *models.Debt) error {
	return r.db.Create(debt).Error
}

func (r *DebtRepository) FindByID(id uint) (*models.Debt, error) {
	var debt models.Debt
	err := r.db.First(&debt, id).Error
	if err != nil {
		return nil, err
	}
	return &debt, nil
}

func (r *DebtRepository) FindAllWithPagination(page, limit int, year, month *int) ([]models.Debt, int64, error) {
	var debts []models.Debt
	var total int64

	query := r.db.Model(&models.Debt{})

	if year != nil {
		startOfYear := time.Date(*year, 1, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := time.Date(*year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		query = query.Where(SQLWhereCreatedAtRange, startOfYear, endOfYear)
	}

	if month != nil && year != nil {
		startOfMonth := time.Date(*year, time.Month(*month), 1, 0, 0, 0, 0, time.UTC)
		var endOfMonth time.Time
		if *month == 12 {
			endOfMonth = time.Date(*year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		} else {
			endOfMonth = time.Date(*year, time.Month(*month+1), 1, 0, 0, 0, 0, time.UTC)
		}
		query = query.Where(SQLWhereCreatedAtRange, startOfMonth, endOfMonth)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf(ErrCountingDebts, err)
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&debts).Error; err != nil {
		return nil, 0, fmt.Errorf(ErrFindingDebts, err)
	}

	return debts, total, nil
}

func (r *DebtRepository) Delete(id uint) error {
	return r.db.Delete(&models.Debt{}, id).Error
}

func (r *DebtRepository) GetTotalByPersonInMonth(year, month int) ([]PersonTotal, error) {
	var results []PersonTotal

	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	var endOfMonth time.Time
	if month == 12 {
		endOfMonth = time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	} else {
		endOfMonth = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
	}

	err := r.db.Model(&models.Debt{}).
		Select("person, SUM(value) as total").
		Where(SQLWhereCreatedAtRange, startOfMonth, endOfMonth).
		Group("person").
		Order("total DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf(ErrCalculatingTotal, err)
	}

	return results, nil
}
