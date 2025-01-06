package models

import "time"

type Event struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	HouseID      uint      `gorm:"not null;constraint:OnDelete:CASCADE;" json:"house_id"`
	CreatedByUID string    `gorm:"not null" json:"created_by_uid"` // Firebase UID of the creator
	Name         string    `gorm:"not null" json:"name"`
	StartTime    time.Time `gorm:"not null" json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Description  string    `json:"description"`
}
