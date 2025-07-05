package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(id string) (string, error) {
	content := fmt.Sprintf("https://weddingnew.netlify.app/upload/id=%s",id)
	savePath := fmt.Sprintf("%s.png", id)

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	err := qrcode.WriteFile(content, qrcode.Medium, 256, savePath)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	url, err := UploadToCloudinary(savePath)
	if err != nil{
		return "", fmt.Errorf("failed to upload files at cloudinary %v",err)
	}

	// Optional: delete the local file to save disk space
	_ = os.Remove(savePath)

	return url, nil
}
