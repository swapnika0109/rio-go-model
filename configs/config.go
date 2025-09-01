package configs

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Firestore FirestoreConfig
	Storage   StorageConfig
	Server    ServerConfig
}

// FirestoreConfig holds Firestore-specific configuration
type FirestoreConfig struct {
	ProjectID string
	CredentialsPath string
}

// StorageConfig holds Google Cloud Storage configuration
type StorageConfig struct {
	BucketName string
	CredentialsPath string
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port string
	Host string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Firestore: FirestoreConfig{
			ProjectID: getEnv("GOOGLE_CLOUD_PROJECT", "your-project-id"),
			CredentialsPath: getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
		},
		Storage: StorageConfig{
			BucketName: getEnv("STORAGE_BUCKET_NAME", "kutty_bucket"),
			CredentialsPath: getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
		},
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
