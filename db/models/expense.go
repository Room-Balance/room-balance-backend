package models

import "time"

type Expense struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	HouseID     uint      `gorm:"not null;constraint:OnDelete:CASCADE;" json:"house_id"`
	PayerUID    string    `gorm:"not null" json:"payer_uid"` // Firebase UID
	Amount      float64   `gorm:"not null" json:"amount"`
	Date        time.Time `gorm:"not null" json:"date"`
	SplitAmong  string    `json:"split_among"` // JSON array of Firebase UIDs
	Description string    `json:"description"`
}
