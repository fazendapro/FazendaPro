package models

type User struct {
	ID       uint    `gorm:"primaryKey"`
	PersonID *uint   `gorm:"unique"` // Nullable for existing records
	Person   *Person `gorm:"foreignKey:PersonID"`
	FarmID   uint    `gorm:"not null"`
	Farm     Farm    `gorm:"foreignKey:FarmID"`
}
