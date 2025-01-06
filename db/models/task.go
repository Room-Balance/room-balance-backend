package models

import "time"

type Task struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	HouseID          uint      `gorm:"not null;constraint:OnDelete:CASCADE;" json:"house_id"`
	AssignedToUserID string    `gorm:"not null" json:"assigned_to_user_id"` // Firebase UID
	Description      string    `gorm:"not null" json:"description"`
	Type             string    `gorm:"type:task_type;not null" json:"type"`
	Status           string    `gorm:"type:task_status;default:'pending'" json:"status"`
	DueDate          time.Time `json:"due_date"`
}
