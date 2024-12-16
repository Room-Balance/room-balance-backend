package models

import "time"

type UserHouse struct {
	UserID   uint      `gorm:"primaryKey;autoIncrement:false"`
	HouseID  uint      `gorm:"primaryKey;autoIncrement:false"`
	JoinedAt time.Time `gorm:"autoCreateTime"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	House House `gorm:"foreignKey:HouseID;constraint:OnDelete:CASCADE"`
}
