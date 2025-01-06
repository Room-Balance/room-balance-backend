package models

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Email       string    `gorm:"unique;not null" json:"email"`
	FirebaseUID string    `gorm:"unique;not null" json:"firebase_uid"` // Used for user identification
	JoinedAt    time.Time `json:"joined_at"`
}
