package models

import "time"

type Event struct {
	ID              uint      `gorm:"primaryKey"`
	HouseID         uint      `gorm:"not null"`
	CreatedByUserID uint      `gorm:"not null"`
	Name            string    `gorm:"type:varchar(100);not null"`
	StartTime       time.Time `gorm:"not null"`
	EndTime         time.Time `gorm:"not null"`
	Description     string    `gorm:"type:text"`

	// Relationships
	House   House `gorm:"foreignKey:HouseID;constraint:OnDelete:CASCADE"`
	Creator User  `gorm:"foreignKey:CreatedByUserID;constraint:OnDelete:CASCADE"`
}
