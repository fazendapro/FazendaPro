package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	PersonID uint   `gorm:"not null;unique"`
	Person   Person `gorm:"foreignKey:PersonID"`
	FarmID   uint   `gorm:"not null"`
	Farm     Farm   `gorm:"foreignKey:FarmID"`
}
