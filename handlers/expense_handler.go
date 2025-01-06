package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/Room-Balance/room-balance-backend.git/db"
	"github.com/Room-Balance/room-balance-backend.git/db/models"
	"github.com/gorilla/mux"
)

// Get all expenses
func GetExpenses(w http.ResponseWriter, r *http.Request) {
	// Extract Firebase UID from the request context
	firebaseUID, ok := r.Context().Value("firebase_uid").(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "Unauthorized: Firebase UID missing", http.StatusUnauthorized)
		return
	}

	// Fetch the house associated with the user
	var house models.House
	if err := db.DB.Joins("JOIN user_houses ON user_houses.house_id = houses.id").
		Where("user_houses.firebase_uid = ?", firebaseUID).First(&house).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	// Fetch expenses belonging to the house
	var expenses []models.Expense
	if err := db.DB.Where("house_id = ?", house.ID).Find(&expenses).Error; err != nil {
		http.Error(w, "Failed to fetch expenses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// Get a specific expense
func GetExpense(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var expense models.Expense
	if err := db.DB.First(&expense, id).Error; err != nil {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(expense)
}

// Create an expense
func CreateExpense(w http.ResponseWriter, r *http.Request) {
	// Extract Firebase UID
	firebaseUID, ok := r.Context().Value("firebase_uid").(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "Unauthorized: Firebase UID missing", http.StatusUnauthorized)
		return
	}

	// Fetch the user’s house
	var house models.House
	if err := db.DB.Joins("JOIN user_houses ON user_houses.house_id = houses.id").
		Where("user_houses.firebase_uid = ?", firebaseUID).First(&house).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	// Decode the expense data
	var expense models.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Assign the expense to the house and payer
	expense.HouseID = house.ID
	expense.PayerUID = firebaseUID
	if err := db.DB.Create(&expense).Error; err != nil {
		http.Error(w, "Failed to create expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

// Update an expense
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var expense models.Expense
	if err := db.DB.First(&expense, id).Error; err != nil {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}

	var updatedData models.Expense
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	expense.Amount = updatedData.Amount
	expense.Description = updatedData.Description
	expense.Date = updatedData.Date

	if err := db.DB.Save(&expense).Error; err != nil {
		http.Error(w, "Failed to update expense", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(expense)
}

// Delete an expense
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if err := db.DB.Delete(&models.Expense{}, id).Error; err != nil {
		http.Error(w, "Failed to delete expense", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
