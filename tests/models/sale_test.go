package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSale_TableName(t *testing.T) {
	sale := Sale{}
	assert.Equal(t, "sales", sale.TableName())
}

func TestSale_ValidSale(t *testing.T) {
	now := time.Now()
	sale := Sale{
		AnimalID:  1,
		FarmID:    1,
		BuyerName: "João Silva",
		Price:     1500.50,
		SaleDate:  now,
		Notes:     "Venda para engorda",
	}

	assert.Equal(t, uint(1), sale.AnimalID)
	assert.Equal(t, uint(1), sale.FarmID)
	assert.Equal(t, "João Silva", sale.BuyerName)
	assert.Equal(t, 1500.50, sale.Price)
	assert.Equal(t, now, sale.SaleDate)
	assert.Equal(t, "Venda para engorda", sale.Notes)
}

func TestSale_EmptyFields(t *testing.T) {
	sale := Sale{}

	assert.Equal(t, uint(0), sale.AnimalID)
	assert.Equal(t, uint(0), sale.FarmID)
	assert.Equal(t, "", sale.BuyerName)
	assert.Equal(t, 0.0, sale.Price)
	assert.Equal(t, time.Time{}, sale.SaleDate)
	assert.Equal(t, "", sale.Notes)
}
