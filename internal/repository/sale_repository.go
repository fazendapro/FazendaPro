package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"

	"gorm.io/gorm"
)

type MonthlySalesData struct {
	Month string  `json:"month"`
	Year  int     `json:"year"`
	Sales float64 `json:"sales"`
	Count int64   `json:"count"`
}

type OverviewStats struct {
	MalesCount   int64   `json:"males_count"`
	FemalesCount int64   `json:"females_count"`
	TotalSold    int64   `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

type SaleRepository interface {
	Create(ctx context.Context, sale *models.Sale) error
	GetByID(ctx context.Context, id uint) (*models.Sale, error)
	GetByFarmID(ctx context.Context, farmID uint) ([]*models.Sale, error)
	GetByAnimalID(ctx context.Context, animalID uint) ([]*models.Sale, error)
	GetByDateRange(ctx context.Context, farmID uint, startDate, endDate time.Time) ([]*models.Sale, error)
	GetMonthlySalesCount(ctx context.Context, farmID uint, startDate, endDate time.Time) (int64, error)
	GetMonthlySalesData(ctx context.Context, farmID uint, months int) ([]MonthlySalesData, error)
	GetOverviewStats(ctx context.Context, farmID uint) (*OverviewStats, error)
	Update(ctx context.Context, sale *models.Sale) error
	Delete(ctx context.Context, id uint) error
}

type saleRepository struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) SaleRepository {
	return &saleRepository{db: db}
}

func (r *saleRepository) Create(ctx context.Context, sale *models.Sale) error {
	return r.db.WithContext(ctx).Create(sale).Error
}

func (r *saleRepository) GetByID(ctx context.Context, id uint) (*models.Sale, error) {
	var sale models.Sale
	err := r.db.WithContext(ctx).Preload("Animal").Preload("Farm").First(&sale, id).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *saleRepository) GetByFarmID(ctx context.Context, farmID uint) ([]*models.Sale, error) {
	var sales []*models.Sale
	err := r.db.WithContext(ctx).Preload("Animal").Where("farm_id = ?", farmID).Order("sale_date DESC").Find(&sales).Error
	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *saleRepository) GetByAnimalID(ctx context.Context, animalID uint) ([]*models.Sale, error) {
	var sales []*models.Sale
	err := r.db.WithContext(ctx).Preload("Animal").Where("animal_id = ?", animalID).Order("sale_date DESC").Find(&sales).Error
	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *saleRepository) GetByDateRange(ctx context.Context, farmID uint, startDate, endDate time.Time) ([]*models.Sale, error) {
	var sales []*models.Sale
	err := r.db.WithContext(ctx).Preload("Animal").Where("farm_id = ? AND sale_date BETWEEN ? AND ?", farmID, startDate, endDate).Order("sale_date DESC").Find(&sales).Error
	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *saleRepository) GetMonthlySalesCount(ctx context.Context, farmID uint, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Sale{}).
		Where("farm_id = ? AND sale_date BETWEEN ? AND ?", farmID, startDate, endDate).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *saleRepository) GetMonthlySalesData(ctx context.Context, farmID uint, months int) ([]MonthlySalesData, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -months+1, 0)
	endDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, 0).Add(-time.Nanosecond)

	type Result struct {
		Year  int     `gorm:"column:year"`
		Month int     `gorm:"column:month"`
		Sales float64 `gorm:"column:sales"`
		Count int64   `gorm:"column:count"`
	}

	var results []Result
	err := r.db.WithContext(ctx).
		Table("sales").
		Select("EXTRACT(YEAR FROM sale_date)::int as year, EXTRACT(MONTH FROM sale_date)::int as month, COALESCE(SUM(price), 0) as sales, COUNT(*)::bigint as count").
		Where("farm_id = ? AND sale_date >= ? AND sale_date <= ?", farmID, startDate, endDate).
		Group("EXTRACT(YEAR FROM sale_date), EXTRACT(MONTH FROM sale_date)").
		Order("year ASC, month ASC").
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching monthly sales data: %w", err)
	}

	monthNames := []string{"Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"}
	monthlyData := make([]MonthlySalesData, 0, months)

	resultMap := make(map[string]Result)
	for _, result := range results {
		key := fmt.Sprintf("%d-%d", result.Year, result.Month)
		resultMap[key] = result
	}

	for i := 0; i < months; i++ {
		currentDate := startDate.AddDate(0, i, 0)
		year := currentDate.Year()
		month := int(currentDate.Month())
		key := fmt.Sprintf("%d-%d", year, month)

		var data MonthlySalesData
		if result, ok := resultMap[key]; ok {
			data = MonthlySalesData{
				Month: monthNames[month-1],
				Year:  year,
				Sales: result.Sales,
				Count: result.Count,
			}
		} else {
			data = MonthlySalesData{
				Month: monthNames[month-1],
				Year:  year,
				Sales: 0,
				Count: 0,
			}
		}
		monthlyData = append(monthlyData, data)
	}

	return monthlyData, nil
}

func (r *saleRepository) Update(ctx context.Context, sale *models.Sale) error {
	return r.db.WithContext(ctx).Save(sale).Error
}

func (r *saleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Sale{}, id).Error
}

func (r *saleRepository) GetOverviewStats(ctx context.Context, farmID uint) (*OverviewStats, error) {
	stats := &OverviewStats{}

	var malesCount int64
	err := r.db.WithContext(ctx).Model(&models.Animal{}).Where("farm_id = ? AND sex = ?", farmID, 1).Count(&malesCount).Error
	if err != nil {
		return nil, fmt.Errorf("error counting males: %w", err)
	}
	stats.MalesCount = malesCount

	var femalesCount int64
	err = r.db.WithContext(ctx).Model(&models.Animal{}).Where("farm_id = ? AND sex = ?", farmID, 0).Count(&femalesCount).Error
	if err != nil {
		return nil, fmt.Errorf("error counting females: %w", err)
	}
	stats.FemalesCount = femalesCount

	var totalSold int64
	err = r.db.WithContext(ctx).Model(&models.Sale{}).Where("farm_id = ?", farmID).Count(&totalSold).Error
	if err != nil {
		return nil, fmt.Errorf("error counting total sold: %w", err)
	}
	stats.TotalSold = totalSold

	var totalRevenue float64
	err = r.db.WithContext(ctx).Model(&models.Sale{}).
		Where("farm_id = ?", farmID).
		Select("COALESCE(SUM(price), 0)").
		Scan(&totalRevenue).Error
	if err != nil {
		return nil, fmt.Errorf("error calculating total revenue: %w", err)
	}
	stats.TotalRevenue = totalRevenue

	return stats, nil
}
