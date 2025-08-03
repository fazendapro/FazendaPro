package models

import (
	"time"
)

type Animal struct {
	ID            uint    `gorm:"primaryKey"`
	FarmID        uint    `gorm:"not null"`
	Farm          Farm    `gorm:"foreignKey:FarmID"`
	EarTagNumber  int     `gorm:"not null;uniqueIndex"`
	AnimalName    string  `gorm:"not null"`
	Sex           int     `gorm:"not null"` // 0 = Female, 1 = Male
	Breed         string  `gorm:"not null"`
	Type          string  `gorm:"not null"`
	Age           int     `gorm:"not null"`
	Photo         string  // URL or path to animal photo
	FatherID      *uint   // FK para Animal (pai)
	Father        *Animal `gorm:"foreignKey:FatherID"`
	MotherID      *uint   // FK para Animal (mãe)
	Mother        *Animal `gorm:"foreignKey:MotherID"`
	Confinement   bool    `gorm:"default:false"`
	AnimalType    int     `gorm:"not null"` // 0 = Cattle, 1 = Buffalo, etc.
	Status        string  `gorm:"not null"`
	Fertilization string
	Castrated     bool   `gorm:"default:false"`
	Purpose       string `gorm:"not null"` // Meat, Milk, Breeding
	CurrentBatch  int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
