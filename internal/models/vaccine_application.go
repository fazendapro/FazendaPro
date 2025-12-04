package models

import (
	"time"
)

type VaccineApplication struct {
	ID              uint      `gorm:"primaryKey"`
	AnimalID        uint      `gorm:"not null"`
	Animal          Animal    `gorm:"foreignKey:AnimalID"`
	VaccineID       uint      `gorm:"not null"`
	Vaccine         Vaccine   `gorm:"foreignKey:VaccineID"`
	ApplicationDate time.Time `gorm:"not null"`
	BatchNumber     string
	Veterinarian    string
	Observations    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

