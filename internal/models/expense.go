package models

import (
	"time"
)

type Expense struct {
	ID          uint      `gorm:"primaryKey"`
	FarmID      uint      `gorm:"not null"`
	Farm        Farm      `gorm:"foreignKey:FarmID"`
	Description string    `gorm:"not null"`
	Amount      float64   `gorm:"not null"`
	Category    string    `gorm:"not null"` // Food, Veterinary, Equipment, etc.
	Date        time.Time `gorm:"not null"`
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
