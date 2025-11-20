package models

import "time"

type UserFarm struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	FarmID    uint      `gorm:"not null" json:"farm_id"`
	Farm      Farm      `gorm:"foreignKey:FarmID" json:"farm"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
