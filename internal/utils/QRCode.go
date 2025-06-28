package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(id string) (string, error) {
	content := fmt.Sprintf("http://localhost:8000?id=%s",id)
	fileName := fmt.Sprintf("%s.png", id)
	savePath := filepath.Join("assets", "qrcodes", fileName)

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	err := qrcode.WriteFile(content, qrcode.Medium, 256, savePath)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	return savePath, nil
}
