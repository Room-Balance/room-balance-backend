package models

import "time"

type UserHouse struct {
	FirebaseUID string    `gorm:"not null;index" json:"firebase_uid"` // Use firebase_uid instead of user_id
	HouseID     uint      `gorm:"not null" json:"house_id"`
	JoinedAt    time.Time `json:"joined_at"`
}
