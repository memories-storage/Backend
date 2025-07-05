package utils

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"os"
)

func UploadToCloudinary(localFilePath string) (string, error) {
	// Load environment variables
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	
	// Validate environment variables
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return "", fmt.Errorf("missing Cloudinary environment variables")
	}

	// Initialize Cloudinary
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	// Try method 1: File reader upload
	url, err := uploadWithFileReader(cld, localFilePath)
	if err == nil && url != "" {
		return url, nil
	}
	
	// Try method 2: File path upload
	url, err = uploadWithFilePath(cld, localFilePath)
	if err == nil && url != "" {
		return url, nil
	}
	
	return "", fmt.Errorf("all upload methods failed")
}

func uploadWithFileReader(cld *cloudinary.Cloudinary, localFilePath string) (string, error) {
	// Open the file for reading
	file, err := os.Open(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	
	ctx := context.Background()
	useFilename := true
	uniqueFilename := true
	uploadParams := uploader.UploadParams{
		ResourceType: "auto",
		UseFilename:  &useFilename,
		UniqueFilename: &uniqueFilename,
	}
	
	resp, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", err
	}
	
	if resp.SecureURL == "" {
		return "", fmt.Errorf("SecureURL is empty")
	}
	
	return resp.SecureURL, nil
}

func uploadWithFilePath(cld *cloudinary.Cloudinary, localFilePath string) (string, error) {
	ctx := context.Background()
	useFilename := true
	uniqueFilename := true
	uploadParams := uploader.UploadParams{
		ResourceType: "auto",
		UseFilename:  &useFilename,
		UniqueFilename: &uniqueFilename,
	}
	
	resp, err := cld.Upload.Upload(ctx, localFilePath, uploadParams)
	if err != nil {
		return "", err
	}
	
	if resp.SecureURL == "" {
		return "", fmt.Errorf("SecureURL is empty")
	}
	
	return resp.SecureURL, nil
}
