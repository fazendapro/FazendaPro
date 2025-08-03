package models

import (
	"time"
)

type Farm struct {
	ID        uint      `gorm:"primaryKey"`
	CompanyID uint      `gorm:"not null"`
	Company   Company   `gorm:"foreignKey:CompanyID"`
	Users     []User    `gorm:"foreignKey:FarmID"`
	Animals   []Animal  `gorm:"foreignKey:FarmID"`
	Expenses  []Expense `gorm:"foreignKey:FarmID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ChangeFarm allows a user to change farms
func (f *Farm) ChangeFarm(userID uint, newFarmID uint) error {
	// Implementation of farm change logic
	// Here you can add validations and business rules
	return nil
}
