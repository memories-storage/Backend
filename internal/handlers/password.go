package handlers

import (
	"Backend/internal/db"
	"Backend/internal/middleware"
	"Backend/internal/utils"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword"`
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

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	// Default response message
	message := "If the email exists, a password reset link has been sent."
	response := map[string]string{
		"message":   message,
		"resetLink": "", // Default empty
	}

	// Try to find user
	var userID string
	err := db.DB.QueryRow(`SELECT id FROM users WHERE email = $1`, email).Scan(&userID)
	// Send email if user found
	if err == nil {
		token, tokenErr := utils.GenerateResetToken(userID)
		if tokenErr == nil {
			resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
			_ = utils.SendResetEmail(email, resetLink) // send reset email
			response["resetLink"] = resetLink          // for testing, can be removed later
		}
	}

	respondWithJSON(w, http.StatusOK, response)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from URL param
	token := r.URL.Query().Get("token")
	if token == "" {
		respondWithError(w, http.StatusBadRequest, "Reset token is required in the query parameters")
		return
	}

	// Parse JSON body
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if req.NewPassword == "" {
		respondWithError(w, http.StatusBadRequest, "New password is required")
		return
	}

	// Parse and verify token
	userID, err := utils.ParseResetToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired reset token")
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update password in the DB
	_, err = db.DB.Exec(`UPDATE users SET password = $1 WHERE id = $2`, string(hashedPassword), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Password has been reset successfully",
	})
}
