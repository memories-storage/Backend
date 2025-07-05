package handlers

import (
	"Backend/internal/db"
	"Backend/internal/middleware"
	"Backend/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
)

type SuccessfulFile struct {
	ID     string `json:"id"`
	Link   string `json:"link"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type FailedFile struct {
	Name string `json:"name"`
}

type UploadResponse struct {
	Successful []SuccessfulFile `json:"successful"`
	Failed     []FailedFile     `json:"failed"`
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

	// Parse multipart form (allowing up to 32 MB in memory)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		AddFilerespondWithError(w, http.StatusBadRequest, "Error parsing form")
		return
	}

	// Get userId and deviceInfo from form data
	userId := r.FormValue("userId")
	deviceInfo := r.FormValue("deviceInfo")
	if userId == "" || deviceInfo == "" {
		AddFilerespondWithError(w, http.StatusBadRequest, "Missing userId or deviceInfo")
		return
	}

	// Get all uploaded files (files[])
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		AddFilerespondWithError(w, http.StatusBadRequest, "No files uploaded!")
		return
	}

	var successfulFiles []SuccessfulFile
	var failedFiles []FailedFile

	// Helper function to add failed file
	addFailedFile := func(filename string) {
		failedFiles = append(failedFiles, FailedFile{Name: filename})
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			addFailedFile(fileHeader.Filename)
			continue
		}

		// Save to temp location
		tmpDir := os.TempDir()
		tmpFile := filepath.Join(tmpDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename))
		out, err := os.Create(tmpFile)
		if err != nil {
			file.Close()
			addFailedFile(fileHeader.Filename)
			continue
		}

		_, err = io.Copy(out, file)
		out.Close()
		file.Close()
		if err != nil {
			os.Remove(tmpFile)
			addFailedFile(fileHeader.Filename)
			continue
		}

		// Upload to Cloudinary
		imageUrl, err := utils.UploadToCloudinary(tmpFile)
		os.Remove(tmpFile)
		if err != nil {
			addFailedFile(fileHeader.Filename)
			continue
		}

		// Insert into DB
		var id string
		err = db.DB.QueryRow(
			"INSERT INTO images (user_id, image_url, device_info) VALUES ($1, $2, $3) RETURNING id",
			userId, imageUrl, deviceInfo,
		).Scan(&id)
		if err != nil {
			addFailedFile(fileHeader.Filename)
			continue
		}

		// Add to successful files
		successfulFile := SuccessfulFile{
			ID:     id,
			UserID: userId,
			Link:   imageUrl,
			Name:   fileHeader.Filename,
		}
		successfulFiles = append(successfulFiles, successfulFile)
	}

	// Create structured response
	response := UploadResponse{
		Successful: successfulFiles,
		Failed:     failedFiles,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func AddFilerespondWithError(w http.ResponseWriter, code int, message ...string) {
	w.WriteHeader(code)
	msg := "error"
	if len(message) > 0 {
		msg = message[0]
		if len(message) > 1 {
			msg += ": " + message[1]
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
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
