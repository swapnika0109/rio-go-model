package database

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	// "encoding/json"
	// "os"
	// "path/filepath"
	"strings"
	"time"
	// "io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	// "golang.org/x/oauth2/google"
)

// StorageService represents a Google Cloud Storage service
type StorageService struct {
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle
	googleAccessID string
	privateKey     []byte
}

type serviceAccountKey struct {
	Type        string `json:"type"`
	ProjectID   string `json:"project_id"`
	PrivateKeyID string `json:"private_key_id"`
	ClientEmail string `json:"client_email"`
	ClientID    string `json:"client_id"`
	AuthURI     string `json:"auth_uri"`
	TokenURI    string `json:"token_uri"`
	AuthCertURL string `json:"auth_provider_x509_cert_url"`
	ClientCertURL string `json:"client_x509_cert_url"`
	PrivateKey  string `json:"private_key"`
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

		// Read the service account file to get the credentials
		// data, err := ioutil.ReadFile(credPath)
		// if err != nil {
		// 	return fmt.Errorf("failed to read service account file: %v", err)
		// }
		
		// var key serviceAccountKey
		// if err := json.Unmarshal(data, &key); err != nil {
		// 	log.Printf("Error parsing credentials from JSON: %v", err)
		// 	return fmt.Errorf("error parsing credentials from JSON: %w", err)
		// }

		

		// s.googleAccessID = key.ClientEmail
		// s.privateKey = []byte(key.PrivateKey)	
		s.googleAccessID = "89695899419-compute@developer.gserviceaccount.com"
		s.privateKey = []byte(`-----BEGIN PRIVATE KEY-----
	MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDJhWAYzRXYSX2b
	e0vywPH1Kivp+g1QgwJkya46pTiR2xvXnz7ZQOhArdqTbTQlTER5OHk2gMqkP9q8
	bluI+ZL1d2Z2N6bSW3aqNewhXQ6DlrfcGthZCnNdTJ4m/T5z4HdElKRx/wVAtkaB
	cdbeHVkJFtX9rCHEuG1+PgzgDdDVajBJQPgbUORa+kWD2Ee1NJxrIX5y0eJ1WlFS
	N3KFhX/QT3Dg4GgmkMYPKlxibyF+62jCc99RmW/4RBcwngBxqqgEBH7bSQCBxVUj
	GuviGKDRl15J/PsbJo95zlh+C8DfOHRnXc9zHj5I5U7JfV3sWIJ/0RT89QVBc8Os
	PlFuN2dLAgMBAAECggEABIJLKh2QFDKlSiHxE6f7VSaJTs15/nScO2adBkUYlObh
	UZEhSl3euL0X5SQ68X5MiHNI3OOQXKH9yrG9kQnbuGvX3o6t4bU+ZQjnf7+9gzqx
	lXRGlMNxuymnRSztNF5WVm25qSZgGzZlqP3/lDIXtVXr1cA7hQ9CrLc1LL3OIOSJ
	hn0H6MjYU5fQiO4V7au5hZp7agFCKjHv/C5VnlZglZWDoKUVIEBDIYabhXugseK0
	r+Leq12HWOLwRvyJtGspffcxoJt14W8xtyQbMIzd+nuXHGMOYxA7s5dkpp7ejpV4
	DKi/JLmqSwTcnMHk8gs0M5XLUx02rc68cEP85EMv6QKBgQD19rWBU2CWI48JRaRh
	dKTw5FsH9FRKkL5G70ZqhpGwGK3XN49CL++uEISET0ozsvL5Rr+2/tzX308/AbmA
	yJa9BxSj2FXartPF7khYuGLdGvZNYtVfkYsq+RiHzMidP8Nf6eCy493nfseSUnrh
	RDzDB+foD9cqni/+AKFwuTmhEwKBgQDRvm0WVvUuI6eqJUrPzYa3HR43R9EK508w
	BeFq/oFJxhS9hY4g9uqeayG0+tV/k1YRcYz3SfGnKUPI5aotr0r5H9CcztzbGS0Z
	sInkMBV+dcEjP+TtjWXQXcoBwkxAGXM8PMcpoQQv0eeQYrdmKfn812bTXtaCf+pN
	inAR+4Of6QKBgQC5x5+u0ChLHw6h5T1U8wLGhOURLAYi0BM4hvB786rFBM8BmRCW
	4Jg7tHQzb6RPSmHl0P9rnDN4xk4X/Bh+YoQgwOFhJD0s0RJaFTIF+TeLZBsXtogO
	CGPCkKVrpUz+ITRUFNZIqH5qpULejXgNADqY4TbT/Gr74MHFK/rEptMViwKBgHsc
	z2NBf4CppQyV/yeid+SbztSb7vP7eduyV/I5mSH4hswHzLlEtcpvD8XvRcowbWCn
	yhqM6K855XPSeuV98v/v0L96HODuEi72FLpADx2/eLJ8Gp/lU1HO+3e02JT3W1CM
	TEr/HDoFd2qkxwnMsdPbi5ueG0NWWe1RyR6FB2mhAoGAdwQaO74tDGFGuojyzQeK
	HLtOXsppdNrnlM9+QTOV44Enm/dSDxVj62t3pBkiB+3IQVbA0kmj93eqg5PJ3oVM
	OSinHwFpOnydDCy8XdTXKofPozxkimEqpGwLkWhY/H1/H6ICWZyib0Ky77F5NtPR
	d3aatywNyxsAv8BroRox5Ak=
	-----END PRIVATE KEY-----`)
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
	if err != nil && err != fmt.Errorf("iterator done: %v", err) {
		log.Println("failed to test storage connection: %v", err)
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

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4, // Use V4 signing scheme.
		Method:  "GET",                   // The URL allows a GET request.
		Expires: time.Now().Add(15 * time.Minute),
		GoogleAccessID: s.googleAccessID,
		PrivateKey:     s.privateKey, // The URL will expire in 15 minutes.
	}

	url , err := s.bucket.SignedURL(blobPath, opts)

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


