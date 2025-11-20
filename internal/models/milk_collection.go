package models

import (
	"time"
)

type MilkCollection struct {
	ID        uint      `gorm:"primaryKey"`
	AnimalID  uint      `gorm:"not null"`
	Animal    Animal    `gorm:"foreignKey:AnimalID"`
	Liters    float64   `gorm:"not null"`
	Date      time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
