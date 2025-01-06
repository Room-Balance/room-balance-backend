package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/Room-Balance/room-balance-backend.git/db"
	"github.com/Room-Balance/room-balance-backend.git/db/models"
	"github.com/gorilla/mux"
)

// Get all tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
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

	// Fetch tasks belonging to the house
	var tasks []models.Task
	if err := db.DB.Where("house_id = ?", house.ID).Find(&tasks).Error; err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Get a specific task
func GetTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var task models.Task
	if err := db.DB.First(&task, id).Error; err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)
}

// Create a task
func CreateTask(w http.ResponseWriter, r *http.Request) {
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

	// Decode the task data
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Assign the task to the house
	task.HouseID = house.ID
	if err := db.DB.Create(&task).Error; err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Update a task
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var task models.Task
	if err := db.DB.First(&task, id).Error; err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	var updatedData models.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	task.Description = updatedData.Description
	task.Type = updatedData.Type
	task.Status = updatedData.Status
	task.DueDate = updatedData.DueDate

	if err := db.DB.Save(&task).Error; err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// Delete a task
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if err := db.DB.Delete(&models.Task{}, id).Error; err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
