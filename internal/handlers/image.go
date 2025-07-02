package handlers

import (
	"Backend/internal/db"
	"Backend/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"Backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type ImageResponse struct {
	ID     string    `json:"id"`
	UserID string `json:"user_id"`
	URL    string `json:"url"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

type Image struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	URL    string `json:"url"`
}


// GET  /images
func GetImagesHandler(w http.ResponseWriter, r *http.Request) {
	// Get authenticated userId from JWT context
	userId, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userId == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid user or not logged in")
		return
	}

	rows, err := db.DB.Query("SELECT id, user_id, image_url FROM images WHERE user_id = $1", userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database query error")
		return
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.ID, &img.UserID, &img.URL); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error scanning images")
			return
		}
		images = append(images, img)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

// POST /images/{userId}
func AddImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get userId from URL param as string (for UUID)
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		respondWithError(w, http.StatusBadRequest, "Missing userId parameter")
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing form")
		return
	}

	// Get image file
	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Image/video file is required")
		return
	}
	defer file.Close()

	// Save file to temp location
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename))
	out, err := os.Create(tmpFile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create temp file")
		return
	}
	defer func() {
		out.Close()
		os.Remove(tmpFile) // cleanup
	}()
	_, err = io.Copy(out, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save temp file")
		return
	}

	// Upload to Cloudinary
	imageUrl, err := utils.UploadToCloudinary(tmpFile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to upload to Cloudinary", err.Error())
		return
	}

	// Insert into DB (store user_id as string/uuid in images table!)
	var id string
	err = db.DB.QueryRow(
		"INSERT INTO images (user_id, image_url) VALUES ($1, $2) RETURNING id",
		userId, imageUrl,
	).Scan(&id)
	if err != nil {
		fmt.Printf("%s", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to save image to database")
		return
	}

	resp := ImageResponse{
		ID:     id,
		UserID: userId,
		URL:    imageUrl,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DELETE /images/{id}
func DeleteImageHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing image id parameter")
		return
	}

	// Check if the image exists (optional but nice for UX)
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM images WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}
	if !exists {
		respondWithError(w, http.StatusNotFound, "Image not found")
		return
	}

	// Delete the image from the database
	_, err = db.DB.Exec("DELETE FROM images WHERE id = $1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete image")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeleteResponse{
		Message: "Image deleted successfully",
	})
}
