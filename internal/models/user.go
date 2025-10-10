package models

type User struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	PersonID  *uint      `gorm:"not null" json:"person_id"`
	Person    *Person    `gorm:"foreignKey:PersonID" json:"person"`
	FarmID    uint       `gorm:"not null" json:"farm_id"`
	Farm      Farm       `gorm:"foreignKey:FarmID" json:"farm"`
	UserFarms []UserFarm `gorm:"foreignKey:UserID" json:"user_farms,omitempty"`
}
