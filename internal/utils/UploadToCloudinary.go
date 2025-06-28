package utils

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"os"
)

func UploadToCloudinary(localFilePath string) (string, error) {
	// Load from env
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	ctx := context.Background()
	resp, err := cld.Upload.Upload(ctx, localFilePath, uploader.UploadParams{})
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}

	return resp.SecureURL, nil
}
