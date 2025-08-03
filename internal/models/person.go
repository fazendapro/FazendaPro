package models

import (
	"time"
)

type Person struct {
	ID         uint   `gorm:"primaryKey"`
	FirstName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
	Email      string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	LastAccess *time.Time
	CPF        string `gorm:"unique;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
