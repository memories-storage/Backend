package handlers

import (
	"Backend/internal/db"
	"Backend/internal/middleware"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Folder struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateFolderRequest struct {
	Name string `json:"name"`
}

type CreateFolderResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// POST /folders
func CreateFolderHandler(w http.ResponseWriter, r *http.Request) {
	// Get authenticated userId from JWT context
	userId, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userId == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid user or not logged in")
		return
	}

	// Parse request body
	var req CreateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate folder name
	if req.Name == "" || len(req.Name) > 255 {
		respondWithError(w, http.StatusBadRequest, "Folder name is required and must be less than 255 characters")
		return
	}

	// Generate folder ID
	folderID := uuid.New().String()

	// Insert folder into database
	_, err := db.DB.Exec(
		"INSERT INTO folders (id, user_id, name) VALUES ($1, $2, $3)",
		folderID, userId, req.Name,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create folder")
		return
	}

	// Return success response
	response := CreateFolderResponse{
		ID:   folderID,
		Name: req.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GET /folders
func GetFoldersHandler(w http.ResponseWriter, r *http.Request) {
	// Get authenticated userId from JWT context
	userId, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userId == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid user or not logged in")
		return
	}

	rows, err := db.DB.Query("SELECT id, user_id, name, created_at, updated_at FROM folders WHERE user_id = $1 ORDER BY created_at DESC", userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database query error")
		return
	}
	defer rows.Close()

	var folders []Folder
	for rows.Next() {
		var folder Folder
		if err := rows.Scan(&folder.ID, &folder.UserID, &folder.Name, &folder.CreatedAt, &folder.UpdatedAt); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error scanning folders")
			return
		}
		folders = append(folders, folder)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(folders)
}

// DELETE /folders/{id}
func DeleteFolderHandler(w http.ResponseWriter, r *http.Request) {
	// Get authenticated userId from JWT context
	userId, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userId == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid user or not logged in")
		return
	}

	folderID := chi.URLParam(r, "id")
	if folderID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing folder id parameter")
		return
	}

	// Check if the folder exists and belongs to the user
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM folders WHERE id = $1 AND user_id = $2)", folderID, userId).Scan(&exists)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if !exists {
		respondWithError(w, http.StatusNotFound, "Folder not found or access denied")
		return
	}

	// Delete the folder from the database
	_, err = db.DB.Exec("DELETE FROM folders WHERE id = $1 AND user_id = $2", folderID, userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete folder")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Folder deleted successfully",
	})
} 