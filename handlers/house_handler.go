package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var house models.House
	if err := db.DB.First(&house, id).Error; err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(house)
}

// Create a house
func CreateHouse(w http.ResponseWriter, r *http.Request) {
	var house models.House
	if err := json.NewDecoder(r.Body).Decode(&house); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := db.DB.Create(&house).Error; err != nil {
		http.Error(w, "Failed to create house", http.StatusInternalServerError)
		return
	}
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
