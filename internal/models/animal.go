package models

import (
	"time"
)

type Animal struct {
	ID                   uint `gorm:"primaryKey"`
	FarmID               uint `gorm:"not null"`
	Farm                 Farm `gorm:"foreignKey:FarmID"`
	EarTagNumberLocal    int  `gorm:"not null"`
	EarTagNumberRegister int
	AnimalName           string `gorm:"not null"`
	Sex                  int    `gorm:"not null"`
	Breed                string `gorm:"not null"`
	Type                 string `gorm:"not null"`
	BirthDate            *time.Time
	Photo                string
	FatherID             *uint
	Father               *Animal `gorm:"foreignKey:FatherID"`
	MotherID             *uint
	Mother               *Animal `gorm:"foreignKey:MotherID"`
	Confinement          bool    `gorm:"default:false"`
	AnimalType           int     `gorm:"not null"`
	Status               int     `gorm:"default:0"`
	Fertilization        bool    `gorm:"default:false"`
	Castrated            bool    `gorm:"default:false"`
	Purpose              int     `gorm:"default:0"`
	CurrentBatch         int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
