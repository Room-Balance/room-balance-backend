package models

import "time"

type Task struct {
	ID               uint       `gorm:"primaryKey"`
	HouseID          uint       `gorm:"not null"`
	AssignedToUserID uint       `gorm:"not null"`
	Description      string     `gorm:"type:text;not null"`
	Type             string     `gorm:"type:task_type;not null"`
	Status           string     `gorm:"type:task_status;default:'pending'"`
	DueDate          *time.Time `gorm:"default:null"` // Nullable for repetitive tasks
	Frequency        *string    `gorm:"type:task_frequency"`
	StartDate        *time.Time `gorm:"default:null"` // Only for repetitive tasks
	EndDate          *time.Time `gorm:"default:null"` // Optional end date for repetitive tasks

	// Relationships
	House    House `gorm:"foreignKey:HouseID;constraint:OnDelete:CASCADE"`
	Assigned User  `gorm:"foreignKey:AssignedToUserID;constraint:OnDelete:SET NULL"`
}
