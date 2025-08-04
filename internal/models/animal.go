package models

import (
	"time"
)

type Animal struct {
	ID                   uint   `gorm:"primaryKey"`
	FarmID               uint   `gorm:"not null"`
	Farm                 Farm   `gorm:"foreignKey:FarmID"`
	EarTagNumberLocal    int    `gorm:"not null"` // Número da brinca local
	EarTagNumberRegister int    // Número da brinca registrado
	AnimalName           string `gorm:"not null"`
	Sex                  int    `gorm:"not null"` // 0 = Female, 1 = Male
	Breed                string `gorm:"not null"`
	Type                 string `gorm:"not null"`
	BirthDate            *time.Time
	Photo                string  // URL or path to animal photo
	FatherID             *uint   // FK para Animal (pai)
	Father               *Animal `gorm:"foreignKey:FatherID"`
	MotherID             *uint   // FK para Animal (mãe)
	Mother               *Animal `gorm:"foreignKey:MotherID"`
	Confinement          bool    `gorm:"default:false"`
	AnimalType           int     `gorm:"not null"`  // 0 = Cattle, 1 = Buffalo, etc.
	Status               int     `gorm:"default:0"` // 0 = Active, 1 = Inactive, 2 = Sold, 3 = Dead
	Fertilization        bool    `gorm:"default:false"`
	Castrated            bool    `gorm:"default:false"`
	Purpose              int     `gorm:"default:0"` // 0 = Meat, 1 = Milk, 2 = Breeding
	CurrentBatch         int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
