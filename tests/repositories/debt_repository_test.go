package repositories

import (
	"testing"
	"time"

	"github.com/fazendapro/FazendaPro-api/internal/models"
	"github.com/fazendapro/FazendaPro-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDebtTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Debt{})
	require.NoError(t, err)

	return db
}

func TestDebtRepository_Create(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	debt := &models.Debt{
		Person: "João Silva",
		Value:  1500.50,
	}

	err := repo.Create(debt)
	assert.NoError(t, err)
	assert.NotZero(t, debt.ID)
}

func TestDebtRepository_FindByID(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	debt := &models.Debt{
		Person: "João Silva",
		Value:  1500.50,
	}
	require.NoError(t, db.Create(debt).Error)

	found, err := repo.FindByID(debt.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, debt.ID, found.ID)
	assert.Equal(t, "João Silva", found.Person)
	assert.Equal(t, 1500.50, found.Value)
}

func TestDebtRepository_FindByID_NotFound(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	found, err := repo.FindByID(999)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestDebtRepository_FindAllWithPagination(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	for i := 0; i < 10; i++ {
		debt := &models.Debt{
			Person: "Pessoa " + string(rune('A'+i)),
			Value:  float64(100 * (i + 1)),
		}
		require.NoError(t, db.Create(debt).Error)
	}

	debts, total, err := repo.FindAllWithPagination(1, 5, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
	assert.Len(t, debts, 5)
}

func TestDebtRepository_FindAllWithPagination_SecondPage(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	for i := 0; i < 10; i++ {
		debt := &models.Debt{
			Person: "Pessoa " + string(rune('A'+i)),
			Value:  float64(100 * (i + 1)),
		}
		require.NoError(t, db.Create(debt).Error)
	}

	debts, total, err := repo.FindAllWithPagination(2, 5, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), total)
	assert.Len(t, debts, 5)
}

func TestDebtRepository_FindAllWithPagination_WithYear(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	now := time.Now()
	currentYear := now.Year()
	lastYear := currentYear - 1

	for i := 0; i < 5; i++ {
		debt := &models.Debt{
			Person:    "Pessoa Atual",
			Value:     100.0,
			CreatedAt: time.Date(currentYear, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC),
		}
		require.NoError(t, db.Create(debt).Error)
	}

	for i := 0; i < 3; i++ {
		debt := &models.Debt{
			Person:    "Pessoa Passado",
			Value:     100.0,
			CreatedAt: time.Date(lastYear, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC),
		}
		require.NoError(t, db.Create(debt).Error)
	}

	year := currentYear
	debts, total, err := repo.FindAllWithPagination(1, 10, &year, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Len(t, debts, 5)
}

func TestDebtRepository_FindAllWithPagination_WithYearAndMonth(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	for i := 0; i < 3; i++ {
		debt := &models.Debt{
			Person:    "Pessoa Mês Atual",
			Value:     100.0,
			CreatedAt: time.Date(currentYear, time.Month(currentMonth), i+1, 0, 0, 0, 0, time.UTC),
		}
		require.NoError(t, db.Create(debt).Error)
	}

	lastMonth := currentMonth - 1
	if lastMonth == 0 {
		lastMonth = 12
		currentYear--
	}
	for i := 0; i < 2; i++ {
		debt := &models.Debt{
			Person:    "Pessoa Mês Passado",
			Value:     100.0,
			CreatedAt: time.Date(currentYear, time.Month(lastMonth), i+1, 0, 0, 0, 0, time.UTC),
		}
		require.NoError(t, db.Create(debt).Error)
	}

	year := now.Year()
	month := int(now.Month())
	debts, total, err := repo.FindAllWithPagination(1, 10, &year, &month)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, debts, 3)
}

func TestDebtRepository_FindAllWithPagination_WithMonthDecember(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	year := 2023
	month := 12

	for i := 0; i < 3; i++ {
		debt := &models.Debt{
			Person:    "Pessoa Dezembro",
			Value:     100.0,
			CreatedAt: time.Date(year, time.December, i+1, 0, 0, 0, 0, time.UTC),
		}
		require.NoError(t, db.Create(debt).Error)
	}

	debts, total, err := repo.FindAllWithPagination(1, 10, &year, &month)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, debts, 3)
}

func TestDebtRepository_Delete(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	debt := &models.Debt{
		Person: "João Silva",
		Value:  1500.50,
	}
	require.NoError(t, db.Create(debt).Error)

	err := repo.Delete(debt.ID)
	assert.NoError(t, err)

	found, err := repo.FindByID(debt.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestDebtRepository_GetTotalByPersonInMonth(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	year := 2023
	month := 6

	for i := 0; i < 3; i++ {
		debt := &models.Debt{
			Person:    "João Silva",
			Value:     100.0,
			CreatedAt: time.Date(year, time.June, i+1, 0, 0, 0, 0, time.UTC),
		}
		require.NoError(t, db.Create(debt).Error)
	}

	for i := 0; i < 2; i++ {
		debt := &models.Debt{
			Person:    "Maria Santos",
			Value:     200.0,
			CreatedAt: time.Date(year, time.June, i+1, 0, 0, 0, 0, time.UTC),
		}
		require.NoError(t, db.Create(debt).Error)
	}

	debt := &models.Debt{
		Person:    "João Silva",
		Value:     500.0,
		CreatedAt: time.Date(year, time.July, 1, 0, 0, 0, 0, time.UTC),
	}
	require.NoError(t, db.Create(debt).Error)

	results, err := repo.GetTotalByPersonInMonth(year, month)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	assert.Equal(t, "Maria Santos", results[0].Person)
	assert.Equal(t, 400.0, results[0].Total)
	assert.Equal(t, "João Silva", results[1].Person)
	assert.Equal(t, 300.0, results[1].Total)
}

func TestDebtRepository_GetTotalByPersonInMonth_December(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	year := 2023
	month := 12

	debt := &models.Debt{
		Person:    "João Silva",
		Value:     100.0,
		CreatedAt: time.Date(year, time.December, 15, 0, 0, 0, 0, time.UTC),
	}
	require.NoError(t, db.Create(debt).Error)

	results, err := repo.GetTotalByPersonInMonth(year, month)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "João Silva", results[0].Person)
	assert.Equal(t, 100.0, results[0].Total)
}

func TestDebtRepository_GetTotalByPersonInMonth_Empty(t *testing.T) {
	db := setupDebtTestDB(t)
	repo := repository.NewDebtRepository(db)

	year := 2023
	month := 6

	results, err := repo.GetTotalByPersonInMonth(year, month)
	assert.NoError(t, err)
	assert.Len(t, results, 0)
}
