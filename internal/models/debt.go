package models

import (
	"time"
)

type Debt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Person    string    `gorm:"not null" json:"person"`
	Value     float64   `gorm:"not null" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
