package models

import (
	"time"
)

type Person struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	FirstName  string     `gorm:"not null" json:"first_name"`
	LastName   string     `gorm:"not null" json:"last_name"`
	Email      string     `gorm:"unique;not null" json:"email"`
	Password   string     `gorm:"not null" json:"password"`
	LastAccess *time.Time `json:"last_access"`
	CPF        string     `gorm:"unique;not null" json:"cpf"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
