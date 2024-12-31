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
	var expenses []models.Expense
	if err := db.DB.Find(&expenses).Error; err != nil {
		http.Error(w, "Failed to fetch expenses", http.StatusInternalServerError)
		return
	}
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
	var expense models.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := db.DB.Create(&expense).Error; err != nil {
		http.Error(w, "Failed to create expense", http.StatusInternalServerError)
		return
	}
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
