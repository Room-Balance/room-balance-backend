package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/Room-Balance/room-balance-backend.git/db"
	"github.com/Room-Balance/room-balance-backend.git/db/models"
	"github.com/gorilla/mux"
)

// Get all events
func GetEvents(w http.ResponseWriter, r *http.Request) {
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

	// Fetch events belonging to the house
	var events []models.Event
	if err := db.DB.Where("house_id = ?", house.ID).Find(&events).Error; err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// Get a specific event
func GetEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var event models.Event
	if err := db.DB.First(&event, id).Error; err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(event)
}

// Create an event
func CreateEvent(w http.ResponseWriter, r *http.Request) {
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

	// Decode the event data
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Assign the event to the house and creator
	event.HouseID = house.ID
	event.CreatedByUID = firebaseUID
	if err := db.DB.Create(&event).Error; err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// Update an event
func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var event models.Event
	if err := db.DB.First(&event, id).Error; err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	var updatedData models.Event
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	event.Name = updatedData.Name
	event.StartTime = updatedData.StartTime
	event.EndTime = updatedData.EndTime
	event.Description = updatedData.Description

	if err := db.DB.Save(&event).Error; err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(event)
}

// Delete an event
func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	if err := db.DB.Delete(&models.Event{}, id).Error; err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
