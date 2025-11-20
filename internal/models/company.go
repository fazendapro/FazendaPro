package models

import (
	"time"
)

type Company struct {
	ID          uint   `gorm:"primaryKey"`
	CompanyName string `gorm:"not null"`
	Location    string `gorm:"not null"`
	FarmCNPJ    string `gorm:"not null;uniqueIndex"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
