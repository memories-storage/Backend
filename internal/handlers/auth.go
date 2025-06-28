package handlers

import (
	"Backend/internal/db"
	"Backend/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
	"time"
	"database/sql"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type SignUpReponse struct {
	Message string `json:"message"`
	Email   string `json:"email"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Role    string `json:"role"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func respondWithError(w http.ResponseWriter, statusCode int, message string, details ...string) {
	resp := ErrorResponse{
		Error: message,
	}
	if len(details) > 0 {
		resp.Details = details[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")
	id := uuid.New()

	if email == "" || password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Check if user already exists
	var existingID string
	err = db.DB.QueryRow(`SELECT id FROM users WHERE email = $1`, email).Scan(&existingID)
	if err == nil {
		// User found
		respondWithError(w, http.StatusConflict, "User already exists with this email")
		return
	} else if err != sql.ErrNoRows {
		// Some other DB error
		respondWithError(w, http.StatusInternalServerError, "Database error while checking user")
		return
	}

	// Generate QR code
	QRCodeLink, err := utils.GenerateQRCode(id.String())
	if err != nil || QRCodeLink == "" {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate QR code link")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Insert into database
	_, err = db.DB.Exec(`
		INSERT INTO users (id, first_name, last_name, email, password, qr_code_link)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, id, firstName, lastName, email, string(hashedPassword), QRCodeLink)

	if err != nil {
		fmt.Println("Insert error:", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to save user in database")
		return
	}

	// Success response
	response := SignUpReponse{
		Message: "Signed up successfully",
		Email:   email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")

	if email == "" || password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Query user from DB
	var id, hashedPassword, role string
	err = db.DB.QueryRow(`SELECT id, password, role FROM users WHERE email = $1`, email).Scan(&id, &hashedPassword, &role)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	role = strings.Trim(role, `"`)
	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	fmt.Println(role)
	response := LoginResponse{
		Token:   tokenString,
		Message: "User LoggedIn successfully",
		Role:    role,
	}
	// Return token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
