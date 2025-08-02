package models

import (
	"time"
)

type Farm struct {
	ID        uint    `gorm:"primaryKey"`
	CompanyID uint    `gorm:"not null"`
	Company   Company `gorm:"foreignKey:CompanyID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
