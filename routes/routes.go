package routes

import (
	"github.com/Room-Balance/room-balance-backend.git/handlers"
	"github.com/Room-Balance/room-balance-backend.git/middlewares"
	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	// Initialize the router
	router := mux.NewRouter()

	// Public Routes (e.g., user signup)
	router.HandleFunc("/api/users", handlers.CreateUser).Methods("POST")

	// Protected Routes
	authRouter := router.PathPrefix("/api").Subrouter()
	authRouter.Use(middlewares.AuthMiddleware)

	// User Routes
	authRouter.HandleFunc("/users/me", handlers.GetUser).Methods("GET")
	authRouter.HandleFunc("/users/data", handlers.GetUserData).Methods("GET")

	// House Routes
	authRouter.HandleFunc("/houses/me", handlers.GetHouse).Methods("GET")
	authRouter.HandleFunc("/houses/me/rent", handlers.UpdateHouseRent).Methods("PUT")
	authRouter.HandleFunc("/houses/me/add-user", handlers.AddUserToHouse).Methods("POST")
	authRouter.HandleFunc("/houses/me/rent", handlers.UpdateRentPayments).Methods("POST")
	authRouter.HandleFunc("/houses/me/expenses", handlers.UpdateExpensePayments).Methods("POST")
	authRouter.HandleFunc("/houses", handlers.CreateHouse).Methods("POST")

	// Task Routes
	authRouter.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	authRouter.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")

	// Expense Routes
	authRouter.HandleFunc("/expenses", handlers.GetExpenses).Methods("GET")
	authRouter.HandleFunc("/expenses", handlers.CreateExpense).Methods("POST")

	// Event Routes
	authRouter.HandleFunc("/events", handlers.GetEvents).Methods("GET")
	authRouter.HandleFunc("/events", handlers.CreateEvent).Methods("POST")

	return router
}
