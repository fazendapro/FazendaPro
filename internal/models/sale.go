package models

import (
	"time"
)

type Sale struct {
	ID        uint      `gorm:"primaryKey"`
	AnimalID  uint      `gorm:"not null"`
	Animal    Animal    `gorm:"foreignKey:AnimalID"`
	FarmID    uint      `gorm:"not null"`
	Farm      Farm      `gorm:"foreignKey:FarmID"`
	BuyerName string    `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	SaleDate  time.Time `gorm:"not null"`
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
