package models

import (
	"time"
)

type Farm struct {
	ID        uint    `gorm:"primaryKey"`
	CompanyID uint    `gorm:"not null"`
	Company   Company `gorm:"foreignKey:CompanyID"`
	Logo      string
	Users     []User    `gorm:"foreignKey:FarmID"`
	Animals   []Animal  `gorm:"foreignKey:FarmID"`
	Expenses  []Expense `gorm:"foreignKey:FarmID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (f *Farm) ChangeFarm(userID uint, newFarmID uint) error {
	return nil
}
