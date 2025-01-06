package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Room-Balance/room-balance-backend.git/db"
	"github.com/Room-Balance/room-balance-backend.git/db/models"
	"github.com/gorilla/mux"
)

// Get all houses
func GetHouses(w http.ResponseWriter, r *http.Request) {
	var houses []models.House
	if err := db.DB.Find(&houses).Error; err != nil {
		http.Error(w, "Failed to fetch houses", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(houses)
}

// Get a specific house
func GetHouse(w http.ResponseWriter, r *http.Request) {
	// Extract Firebase UID from the request context
	firebaseUID, ok := r.Context().Value("firebase_uid").(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "Unauthorized: Firebase UID missing", http.StatusUnauthorized)
		return
	}

	// Fetch the user's house
	var house models.House
	if err := db.DB.Joins("JOIN user_houses ON user_houses.house_id = houses.id").
		Where("user_houses.firebase_uid = ?", firebaseUID).First(&house).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(house)
}

// Create a house
func CreateHouse(w http.ResponseWriter, r *http.Request) {
	// Extract Firebase UID from the request context
	firebaseUID, ok := r.Context().Value("firebase_uid").(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "Unauthorized: Firebase UID missing", http.StatusUnauthorized)
		return
	}

	// Check if the user already belongs to a house
	var existingUserHouse models.UserHouse
	if err := db.DB.Where("firebase_uid = ?", firebaseUID).First(&existingUserHouse).Error; err == nil {
		http.Error(w, "User already belongs to a house", http.StatusConflict)
		return
	}

	// Decode the house data from the request
	var request struct {
		Name string  `json:"name"`
		Rent float64 `json:"rent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate input
	if request.Name == "" {
		http.Error(w, "House name is required", http.StatusBadRequest)
		return
	}
	if request.Rent <= 0 {
		http.Error(w, "Rent must be a positive value", http.StatusBadRequest)
		return
	}

	// Initialize the new house with default payments
	house := models.House{
		Name:            request.Name,
		Rent:            request.Rent,
		TotalExpenses:   0,    // Initialize total expenses to zero
		RentPayments:    "{}", // Initialize as an empty JSON object
		ExpensePayments: "{}", // Initialize as an empty JSON object
	}

	// Create the house in the database
	if err := db.DB.Create(&house).Error; err != nil {
		http.Error(w, "Failed to create house", http.StatusInternalServerError)
		return
	}

	// Link the house to the user
	userHouse := models.UserHouse{
		FirebaseUID: firebaseUID,
		HouseID:     house.ID,
		JoinedAt:    time.Now(),
	}
	if err := db.DB.Create(&userHouse).Error; err != nil {
		http.Error(w, "Failed to link user to house", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(house)
}

// Update a house
func UpdateHouse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var house models.House
	if err := db.DB.First(&house, id).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	var updatedData models.House
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	house.Name = updatedData.Name

	if err := db.DB.Save(&house).Error; err != nil {
		http.Error(w, "Failed to update house", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(house)
}

func UpdateHouseRent(w http.ResponseWriter, r *http.Request) {
	// Extract Firebase UID from the context
	firebaseUID, ok := r.Context().Value("firebase_uid").(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "Unauthorized: Firebase UID missing", http.StatusUnauthorized)
		return
	}

	// Fetch the user's house
	var house models.House
	if err := db.DB.Joins("JOIN user_houses ON user_houses.house_id = houses.id").
		Where("user_houses.firebase_uid = ?", firebaseUID).First(&house).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	// Parse the request payload
	var request struct {
		NewRent float64 `json:"new_rent"` // New rent amount
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate the new rent amount
	if request.NewRent <= 0 {
		http.Error(w, "Rent must be a positive value", http.StatusBadRequest)
		return
	}

	// Update the house's rent
	house.Rent = request.NewRent

	// Save the updated house to the database
	if err := db.DB.Save(&house).Error; err != nil {
		http.Error(w, "Failed to update house rent", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(house)
}

func UpdateRentPayments(w http.ResponseWriter, r *http.Request) {
	// Extract Firebase UID from the request context for authorization
	_, ok := r.Context().Value("firebase_uid").(string)
	if !ok {
		http.Error(w, "Unauthorized: Firebase UID missing", http.StatusUnauthorized)
		return
	}

	// Fetch the user's house
	var house models.House
	if err := db.DB.Joins("JOIN user_houses ON user_houses.house_id = houses.id").
		First(&house).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	// Parse the request payload
	var request struct {
		TargetUID string  `json:"target_uid"` // Firebase UID of the target user
		Amount    float64 `json:"amount"`     // Payment amount
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate the target user
	if request.TargetUID == "" {
		http.Error(w, "Target user UID is required", http.StatusBadRequest)
		return
	}

	// Parse the current rent payments JSON
	var rentPayments map[string]float64
	if house.RentPayments == "" {
		rentPayments = make(map[string]float64)
	} else if err := json.Unmarshal([]byte(house.RentPayments), &rentPayments); err != nil {
		http.Error(w, "Failed to parse rent payments", http.StatusInternalServerError)
		return
	}

	// Update the target user's rent payment
	rentPayments[request.TargetUID] += request.Amount

	// Save the updated rent payments to the database
	updatedPayments, err := json.Marshal(rentPayments)
	if err != nil {
		http.Error(w, "Failed to serialize rent payments", http.StatusInternalServerError)
		return
	}
	house.RentPayments = string(updatedPayments)

	if err := db.DB.Save(&house).Error; err != nil {
		http.Error(w, "Failed to update rent payments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(house)
}

func UpdateExpensePayments(w http.ResponseWriter, r *http.Request) {
	// Extract Firebase UID from the request context for authorization
	_, ok := r.Context().Value("firebase_uid").(string)
	if !ok {
		http.Error(w, "Unauthorized: Firebase UID missing", http.StatusUnauthorized)
		return
	}

	// Fetch the user's house
	var house models.House
	if err := db.DB.Joins("JOIN user_houses ON user_houses.house_id = houses.id").
		First(&house).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	// Parse the request payload
	var request struct {
		TargetUID string  `json:"target_uid"` // Firebase UID of the target user
		Amount    float64 `json:"amount"`     // Payment amount
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate the target user
	if request.TargetUID == "" {
		http.Error(w, "Target user UID is required", http.StatusBadRequest)
		return
	}

	// Parse the current expense payments JSON
	var expensePayments map[string]float64
	if house.ExpensePayments == "" {
		expensePayments = make(map[string]float64)
	} else if err := json.Unmarshal([]byte(house.ExpensePayments), &expensePayments); err != nil {
		http.Error(w, "Failed to parse expense payments", http.StatusInternalServerError)
		return
	}

	// Update the target user's expense payment
	expensePayments[request.TargetUID] += request.Amount

	// Update the total expenses for the house
	house.TotalExpenses += request.Amount

	// Save the updated expense payments and total expenses to the database
	updatedPayments, err := json.Marshal(expensePayments)
	if err != nil {
		http.Error(w, "Failed to serialize expense payments", http.StatusInternalServerError)
		return
	}
	house.ExpensePayments = string(updatedPayments)

	if err := db.DB.Save(&house).Error; err != nil {
		http.Error(w, "Failed to update expense payments or total expenses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(house)
}

// Delete a house
func DeleteHouse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if err := db.DB.Delete(&models.House{}, id).Error; err != nil {
		http.Error(w, "Failed to delete house", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
