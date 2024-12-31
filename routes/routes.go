package routes

import (
	"github.com/Room-Balance/room-balance-backend.git/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	// Initialize the router
	router := mux.NewRouter()

	// User Routes
	router.HandleFunc("/api/users", handlers.GetUsers).Methods("GET")
	router.HandleFunc("/api/users", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id:[0-9]+}", handlers.GetUser).Methods("GET")
	router.HandleFunc("/api/users/{id:[0-9]+}", handlers.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id:[0-9]+}", handlers.DeleteUser).Methods("DELETE")

	// House Routes
	router.HandleFunc("/api/houses", handlers.GetHouses).Methods("GET")
	router.HandleFunc("/api/houses", handlers.CreateHouse).Methods("POST")
	router.HandleFunc("/api/houses/{id:[0-9]+}", handlers.GetHouse).Methods("GET")
	router.HandleFunc("/api/houses/{id:[0-9]+}", handlers.UpdateHouse).Methods("PUT")
	router.HandleFunc("/api/houses/{id:[0-9]+}", handlers.DeleteHouse).Methods("DELETE")

	// Task Routes
	router.HandleFunc("/api/tasks", handlers.GetTasks).Methods("GET")
	router.HandleFunc("/api/tasks", handlers.CreateTask).Methods("POST")
	router.HandleFunc("/api/tasks/{id:[0-9]+}", handlers.GetTask).Methods("GET")
	router.HandleFunc("/api/tasks/{id:[0-9]+}", handlers.UpdateTask).Methods("PUT")
	router.HandleFunc("/api/tasks/{id:[0-9]+}", handlers.DeleteTask).Methods("DELETE")

	// Expense Routes
	router.HandleFunc("/api/expenses", handlers.GetExpenses).Methods("GET")
	router.HandleFunc("/api/expenses", handlers.CreateExpense).Methods("POST")
	router.HandleFunc("/api/expenses/{id:[0-9]+}", handlers.GetExpense).Methods("GET")
	router.HandleFunc("/api/expenses/{id:[0-9]+}", handlers.UpdateExpense).Methods("PUT")
	router.HandleFunc("/api/expenses/{id:[0-9]+}", handlers.DeleteExpense).Methods("DELETE")

	// Event Routes
	router.HandleFunc("/api/events", handlers.GetEvents).Methods("GET")
	router.HandleFunc("/api/events", handlers.CreateEvent).Methods("POST")
	router.HandleFunc("/api/events/{id:[0-9]+}", handlers.GetEvent).Methods("GET")
	router.HandleFunc("/api/events/{id:[0-9]+}", handlers.UpdateEvent).Methods("PUT")
	router.HandleFunc("/api/events/{id:[0-9]+}", handlers.DeleteEvent).Methods("DELETE")

	return router
}
