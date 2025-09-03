package database

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	// "path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// StorageService represents a Google Cloud Storage service
type StorageService struct {
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle
}

// NewStorageService creates a new storage service
func NewStorageService(bucketName string) *StorageService {
	return &StorageService{
		bucketName: bucketName,
	}
}

// Init initializes the storage service
func (s *StorageService) Init(ctx context.Context) error {
	var client *storage.Client
	var err error

	// Try to use service account file first
	credPath := "serviceAccount.json"
	if _, err := os.Stat(credPath); err == nil {
		log.Println("Using service account from file for storage")
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credPath))
	} else {
		log.Println("Using default credentials for storage")
		// In Cloud Run, use the default service account
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return fmt.Errorf("failed to create storage client: %v", err)
	}

	s.client = client
	s.bucket = client.Bucket(s.bucketName)

	// Test connection
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	it := s.bucket.Objects(ctx, nil)
	_, err = it.Next()
	if err != nil && err != iterator.Done {
		return fmt.Errorf("failed to test storage connection: %v", err)
	}

	log.Println("Successfully connected to Storage")
	return nil
}

// Close closes the storage client
func (s *StorageService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// GenerateSignedURL generates a signed URL for a blob
func (s *StorageService) GenerateSignedURL(blobPath string, expiration time.Duration) (string, error) {
	if s.bucket == nil {
		return "", fmt.Errorf("storage service not initialized")
	}

	// Generate signed URL
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expiration),
	}

	url, err := storage.SignedURL(s.bucketName, blobPath, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %v", err)
	}

	return url, nil
}

// UploadFile uploads a file to cloud storage
func (s *StorageService) UploadFile(fileData []byte, prefix, extension string) (string, error) {
	if s.bucket == nil {
		return "", fmt.Errorf("storage service not initialized")
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s/%s.%s", prefix, generateUUID(), extension)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create blob
	blob := s.bucket.Object(filename)

	// Upload file
	writer := blob.NewWriter(ctx)
	defer writer.Close()

	// Set content type based on extension
	contentType := getContentType(extension)
	writer.ContentType = contentType

	// Write data
	if _, err := writer.Write(fileData); err != nil {
		return "", fmt.Errorf("failed to write file data: %v", err)
	}

	log.Printf("Successfully uploaded file: %s", filename)
	return filename, nil
}

// UploadBase64 uploads base64 encoded data to cloud storage
func (s *StorageService) UploadBase64(base64Data, prefix, extension string) (string, error) {
	// Decode base64 data
	fileData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %v", err)
	}

	return s.UploadFile(fileData, prefix, extension)
}

// ReadImage reads an image file from cloud storage
func (s *StorageService) ReadImage(filename string) ([]byte, error) {
	return s.readFile("images", filename)
}

// ReadAudio reads an audio file from cloud storage
func (s *StorageService) ReadAudio(filename string) ([]byte, error) {
	return s.readFile("audio", filename)
}

// ReadFile reads any file from cloud storage
func (s *StorageService) ReadFile(folder, filename string) ([]byte, error) {
	return s.readFile(folder, filename)
}

// readFile is a helper function to read files from storage
func (s *StorageService) readFile(folder, filename string) ([]byte, error) {
	if s.bucket == nil {
		return nil, fmt.Errorf("storage service not initialized")
	}

	// Extract filename from path if it contains slashes
	parts := strings.Split(filename, "/")
	actualFilename := parts[len(parts)-1]

	// Construct full path
	fullPath := fmt.Sprintf("%s/%s", folder, actualFilename)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get blob
	blob := s.bucket.Object(fullPath)

	// Download file
	reader, err := blob.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader for %s: %v", fullPath, err)
	}
	defer reader.Close()

	// Read all data
	data := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			data = append(data, buffer[:n]...)
		}
		if err != nil {
			break
		}
	}

	log.Printf("Successfully read file: %s", fullPath)
	return data, nil
}

// DeleteFile deletes a file from cloud storage
func (s *StorageService) DeleteFile(filename string) error {
	if s.bucket == nil {
		return fmt.Errorf("storage service not initialized")
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete blob
	blob := s.bucket.Object(filename)
	if err := blob.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file %s: %v", filename, err)
	}

	log.Printf("Successfully deleted file: %s", filename)
	return nil
}

// ListFiles lists files in a folder
func (s *StorageService) ListFiles(folder string) ([]string, error) {
	if s.bucket == nil {
		return nil, fmt.Errorf("storage service not initialized")
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var files []string
	it := s.bucket.Objects(ctx, &storage.Query{
		Prefix:    folder + "/",
		Delimiter: "/",
	})

	for {
		obj, err := it.Next()
		if err != nil {
			break
		}
		if obj.Name != folder+"/" {
			files = append(files, obj.Name)
		}
	}

	return files, nil
}

// HealthCheck checks if storage service is accessible
func (s *StorageService) HealthCheck(ctx context.Context) error {
	if s.bucket == nil {
		return fmt.Errorf("storage service not initialized")
	}

	// Try to list objects to check connectivity
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	it := s.bucket.Objects(ctx, &storage.Query{})
	_, err := it.Next()
	if err != nil && err != iterator.Done {
		return fmt.Errorf("storage health check failed: %v", err)
	}

	return nil
}

// Helper functions

// generateUUID generates a UUID (you can use a proper UUID library if needed)
func generateUUID() string {
	// Simple UUID generation - in production, use github.com/google/uuid
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// getContentType returns the MIME type based on file extension
func getContentType(extension string) string {
	switch strings.ToLower(extension) {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	case "mp3":
		return "audio/mpeg"
	case "wav":
		return "audio/wav"
	case "ogg":
		return "audio/ogg"
	case "pdf":
		return "application/pdf"
	case "txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}
