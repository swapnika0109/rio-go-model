package services

import (
	"context"
	"log"
	// "time"

	"rio-go-model/configs"
	"rio-go-model/internal/services/database"
)

// AppService manages the application lifecycle
type AppService struct {
	config    *configs.Config
	firestore *database.StoryDatabase
	storage   *database.StorageService
}

// NewAppService creates a new app service
func NewAppService(config *configs.Config) *AppService {
	return &AppService{
		config: config,
	}
}

// Initialize initializes all services
func (a *AppService) Initialize(ctx context.Context) error {
	log.Println("🚀 Initializing application services...")

	// Initialize Firestore
	if err := a.initializeFirestore(ctx); err != nil {
		return err
	}

	// Initialize Google Cloud Storage
	if err := a.initializeStorage(ctx); err != nil {
		return err
	}

	// Add more service initializations here
	// Example: Redis, external APIs, etc.

	log.Println("✅ All services initialized successfully")
	return nil
}

// initializeFirestore initializes the Firestore connection
func (a *AppService) initializeFirestore(ctx context.Context) error {
	log.Println("🔥 Initializing Firestore...")

	// Create StoryDatabase client
	firestoreClient := database.NewStoryDatabase()

	// Initialize connection
	if err := firestoreClient.Init(ctx); err != nil {
		return err
	}

	// Health check
	if err := firestoreClient.HealthCheck(ctx); err != nil {
		return err
	}

	a.firestore = firestoreClient
	log.Printf("✅ Firestore connected to project: %s", a.config.Firestore.ProjectID)
	return nil
}

// initializeStorage initializes the Google Cloud Storage connection
func (a *AppService) initializeStorage(ctx context.Context) error {
	log.Println("☁️ Initializing Google Cloud Storage...")

	// Create Storage client with bucket name from config
	bucketName := a.config.Storage.BucketName
	if bucketName == "" {
		bucketName = "kutty_bucket" // Default bucket name
	}

	storageClient := database.NewStorageService(bucketName)

	// Initialize connection
	if err := storageClient.Init(ctx); err != nil {
		return err
	}

	// Health check
	if err := storageClient.HealthCheck(ctx); err != nil {
		return err
	}

	a.storage = storageClient
	log.Printf("✅ Google Cloud Storage connected to bucket: %s", bucketName)
	return nil
}

// GetFirestore returns the StoryDatabase client
func (a *AppService) GetFirestore() *database.StoryDatabase {
	return a.firestore
}

// GetStorage returns the Storage client
func (a *AppService) GetStorage() *database.StorageService {
	return a.storage
}

// Shutdown gracefully shuts down all services
func (a *AppService) Shutdown(ctx context.Context) error {
	log.Println("🛑 Shutting down application services...")

	// Close Firestore connection
	if a.firestore != nil {
		if err := a.firestore.Close(); err != nil {
			log.Printf("⚠️ Error closing Firestore: %v", err)
		} else {
			log.Println("✅ Firestore connection closed")
		}
	}

	// Close Storage connection
	if a.storage != nil {
		if err := a.storage.Close(); err != nil {
			log.Printf("⚠️ Error closing Storage: %v", err)
		} else {
			log.Println("✅ Storage connection closed")
		}
	}

	// Add more cleanup here
	// Example: Close Redis, external API connections, etc.

	log.Println("✅ All services shut down successfully")
	return nil
}

// HealthCheck checks the health of all services
func (a *AppService) HealthCheck(ctx context.Context) error {
	// Check Firestore health
	if a.firestore != nil {
		if err := a.firestore.HealthCheck(ctx); err != nil {
			return err
		}
	}

	// Check Storage health
	if a.storage != nil {
		if err := a.storage.HealthCheck(ctx); err != nil {
			return err
		}
	}

	// Add more health checks here
	// Example: Redis, external APIs, etc.

	return nil
}
