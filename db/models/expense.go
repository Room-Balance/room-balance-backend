package models

import (
	"time"

	"gorm.io/datatypes"
)

type Expense struct {
	ID          uint           `gorm:"primaryKey"`
	HouseID     uint           `gorm:"not null"`
	PayerID     uint           `gorm:"not null"`
	Amount      float64        `gorm:"not null"`
	Description string         `gorm:"type:text;not null"`
	Date        time.Time      `gorm:"autoCreateTime"`
	SplitAmong  datatypes.JSON `gorm:"type:json;not null"` // JSON field for splitting details

	// Relationships
	House House `gorm:"foreignKey:HouseID;constraint:OnDelete:CASCADE"`
	Payer User  `gorm:"foreignKey:PayerID;constraint:OnDelete:CASCADE"`
}
