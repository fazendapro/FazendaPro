package models

import (
	"time"
)

type Reproduction struct {
	ID            uint      `gorm:"primaryKey"`
	AnimalID      uint      `gorm:"not null"`
	Animal        Animal    `gorm:"foreignKey:AnimalID"`
	Type          string    `gorm:"not null"` // Natural, Artificial Insemination, etc.
	PregnancyDate time.Time `gorm:"not null"`
	Situation     int       `gorm:"not null"` // 0 = Pregnant, 1 = Not Pregnant, 2 = Unknown
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
