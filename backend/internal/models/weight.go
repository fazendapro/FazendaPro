package models

import (
	"time"
)

type Weight struct {
	ID           uint      `gorm:"primaryKey"`
	AnimalID     uint      `gorm:"not null"`
	Animal       Animal    `gorm:"foreignKey:AnimalID"`
	Date         time.Time `gorm:"not null"`
	AnimalWeight float64   `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
