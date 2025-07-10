package handlers

import (
	"Backend/internal/db"
	"Backend/internal/middleware"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserResponse struct {
	UserID     string    `json:"userId"`
	Email      string    `json:"email"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	QRCodeLink string    `json:"qrCodeLink"`
	CreatedAt  time.Time `json:"createdAt"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	var email, firstName, lastName, qrCodeLink string
	var createdAt time.Time

	err := db.DB.QueryRow(`
		SELECT first_name, last_name, email, qr_code_link, created_at 
		FROM users 
		WHERE id = $1
	`, userID).Scan(&firstName, &lastName, &email, &qrCodeLink, &createdAt)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User not found")
		return
	}

	response := UserResponse{
		UserID:     userID,
		Email:      email,
		FirstName:  firstName,
		LastName:   lastName,
		QRCodeLink: qrCodeLink,
		CreatedAt:  createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if req.FirstName == "" && req.LastName == "" {
		respondWithError(w, http.StatusBadRequest, "At least one of firstName or lastName is required")
		return
	}

	// Update user info in DB
	_, err := db.DB.Exec(`
		UPDATE users SET "first_name" = $1, "last_name" = $2 WHERE id = $3
	`, req.FirstName, req.LastName, userID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user profile")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "User profile updated successfully",
	})
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from JWT context
	userId, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userId == "" {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized or invalid user")
		return
	}

	// Delete the user from the users table
	_, err := db.DB.Exec("DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		fmt.Printf("%s",err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	// You may want to clean up related data (images, etc.) as well

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeleteUserResponse{
		Message: "User deleted successfully",
	})
}
