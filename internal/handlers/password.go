package handlers

import (
	"Backend/internal/db"
	"Backend/internal/middleware"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		respondWithError(w, http.StatusBadRequest, "Both current and new passwords are required")
		return
	}

	// Fetch current hashed password
	var hashedPassword string
	err := db.DB.QueryRow(`SELECT password FROM users WHERE id = $1`, userID).Scan(&hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User not found")
		return
	}

	// Compare current password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.CurrentPassword))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	// Hash new password
	newHashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash new password")
		return
	}

	// Update DB
	_, err = db.DB.Exec(`UPDATE users SET password = $1 WHERE id = $2`, string(newHashed), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Password changed successfully",
	})
}
