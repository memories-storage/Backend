package handlers

import (
	"Backend/internal/db"
	"Backend/internal/middleware"
	"encoding/json"
	"net/http"
)

type UserResponse struct {
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	QRCodeLink string `json:"qrCodeLink"`
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	var email, firstName, lastName, qrCodeLink string
	err := db.DB.QueryRow(`SELECT first_name, last_name, email, qr_code_link FROM users WHERE id = $1`, userID).
		Scan(&firstName, &lastName, &email, &qrCodeLink)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User not found")
		return
	}

	response := UserResponse{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		QRCodeLink: qrCodeLink,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
