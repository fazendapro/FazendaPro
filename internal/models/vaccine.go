package models

import (
	"time"
)

type Vaccine struct {
	ID          uint      `gorm:"primaryKey"`
	FarmID     uint      `gorm:"not null"`
	Farm       Farm      `gorm:"foreignKey:FarmID"`
	Name       string    `gorm:"not null"`
	Description string
	Manufacturer string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

