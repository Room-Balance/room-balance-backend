package models

import "time"

type House struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`

	// New fields
	Rent            float64 `gorm:"not null;default:0" json:"rent"`
	TotalExpenses   float64 `gorm:"not null;default:0" json:"total_expenses"`
	RentPayments    string  `gorm:"type:json" json:"rent_payments"`    // JSON field for rent payments per user
	ExpensePayments string  `gorm:"type:json" json:"expense_payments"` // JSON field for expense payments per user

	// Relationships with cascade delete
	Tasks    []Task    `gorm:"constraint:OnDelete:CASCADE;" json:"tasks"`
	Expenses []Expense `gorm:"constraint:OnDelete:CASCADE;" json:"expenses"`
	Events   []Event   `gorm:"constraint:OnDelete:CASCADE;" json:"events"`
}
