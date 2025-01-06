package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/Room-Balance/room-balance-backend.git/db"
	"github.com/Room-Balance/room-balance-backend.git/db/models"
	"google.golang.org/api/option"
)

var firebaseAuth *auth.Client

func InitFirebase() {
	opt := option.WithCredentialsFile("firebase-service-account.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Auth: %v", err)

	}

	firebaseAuth = client
	log.Println("Firebase initialized successfully!")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token, err := firebaseAuth.VerifyIDToken(context.Background(), tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Sync user with database
		email, _ := token.Claims["email"].(string)
		err = SyncUserWithDatabase(token.UID, email)
		if err != nil {
			http.Error(w, "Failed to sync user", http.StatusInternalServerError)
			return
		}

		// Add Firebase UID to context
		ctx := context.WithValue(r.Context(), "firebase_uid", token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SyncUserWithDatabase(firebaseUID, email string) error {
	var user models.User

	// Check if the user already exists
	if err := db.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err == nil {
		return nil // User already exists
	}

	// Create a new user if not found
	user = models.User{
		FirebaseUID: firebaseUID,
		Email:       email,
		JoinedAt:    time.Now(),
	}
	if err := db.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
