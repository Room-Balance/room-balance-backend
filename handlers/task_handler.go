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
	var tasks []models.Task
	if err := db.DB.Find(&tasks).Error; err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
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
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := db.DB.Create(&task).Error; err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
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
