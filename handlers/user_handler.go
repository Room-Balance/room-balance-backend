package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"strconv"

	"github.com/Room-Balance/room-balance-backend.git/db"
	"github.com/Room-Balance/room-balance-backend.git/db/models"
	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract Firebase UID from the request context
	firebaseUID, ok := r.Context().Value("firebase_uid").(string)
	if !ok || firebaseUID == "" {
		http.Error(w, "Unauthorized: Firebase UID missing", http.StatusUnauthorized)
		return
	}

	// Fetch user information using Firebase UID
	var user models.User
	if err := db.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with the user data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if Firebase UID is provided
	if user.FirebaseUID == "" {
		http.Error(w, "Firebase UID is required", http.StatusBadRequest)
		return
	}

	// Check for existing user with the same Firebase UID
	var existingUser models.User
	if err := db.DB.Where("firebase_uid = ?", user.FirebaseUID).First(&existingUser).Error; err == nil {
		http.Error(w, "User with this Firebase UID already exists", http.StatusConflict)
		return
	}

	// Create the new user
	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updatedData models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user.Name = updatedData.Name
	user.Email = updatedData.Email

	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Get aggregated user-related data
func GetUserData(w http.ResponseWriter, r *http.Request) {
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

	// Fetch tasks, expenses, and events associated with the house
	var tasks []models.Task
	var expenses []models.Expense
	var events []models.Event

	db.DB.Where("house_id = ?", house.ID).Find(&tasks)
	db.DB.Where("house_id = ?", house.ID).Find(&expenses)
	db.DB.Where("house_id = ?", house.ID).Find(&events)

	// Fetch users belonging to the house
	var users []models.User
	db.DB.Joins("JOIN user_houses ON user_houses.firebase_uid = users.firebase_uid").
		Where("user_houses.house_id = ?", house.ID).Find(&users)

	// Aggregate the response
	response := map[string]interface{}{
		"house":    house,
		"tasks":    tasks,
		"expenses": expenses,
		"events":   events,
		"users":    users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AddUserToHouse(w http.ResponseWriter, r *http.Request) {
	// Extract Firebase UID of the requesting user from the context
	requestingUID, ok := r.Context().Value("firebase_uid").(string)
	if !ok || requestingUID == "" {
		http.Error(w, "Unauthorized: Firebase UID missing", http.StatusUnauthorized)
		return
	}

	// Fetch the house associated with the requesting user
	var house models.House
	if err := db.DB.Joins("JOIN user_houses ON user_houses.house_id = houses.id").
		Where("user_houses.firebase_uid = ?", requestingUID).First(&house).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}

	// Parse the payload for the new user's Firebase UID
	var request struct {
		TargetUID string `json:"target_uid"` // Firebase UID of the user to add
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the target UID is provided
	if request.TargetUID == "" {
		http.Error(w, "Target UID is required", http.StatusBadRequest)
		return
	}

	// Check if the target user is already part of the house
	var existingUserHouse models.UserHouse
	if err := db.DB.Where("firebase_uid = ? AND house_id = ?", request.TargetUID, house.ID).
		First(&existingUserHouse).Error; err == nil {
		http.Error(w, "User is already part of this house", http.StatusConflict)
		return
	}

	// Link the target user to the house
	newUserHouse := models.UserHouse{
		FirebaseUID: request.TargetUID,
		HouseID:     house.ID,
		JoinedAt:    time.Now(),
	}
	if err := db.DB.Create(&newUserHouse).Error; err != nil {
		http.Error(w, "Failed to add user to house", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUserHouse)
}
