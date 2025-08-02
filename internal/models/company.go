package models

import (
	"time"
)

type Company struct {
	ID        uint   `gorm:"primaryKey"`
	Location  string `gorm:"not null"`
	FarmCNPJ  string `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
