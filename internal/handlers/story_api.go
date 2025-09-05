package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"rio-go-model/internal/helpers"
	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"
)

// StoryTopics represents the story topics handler
type Story struct {
	storyGenerator *helpers.StoryGenerationHelper
	storageServiceOn bool
	storyDB *database.StoryDatabase
	storageService *database.StorageService
}

// NewStory creates a new story topics handler
func NewStory(storyDB *database.StoryDatabase, 
	storageService *database.StorageService) *Story {
	return &Story{
		storyGenerator: nil,
		storyDB: storyDB,
		storageService: storageService,
	}
}

// CreateStoryRequest represents the request body for creating stories
type CreateStoryRequest struct {
	Country     string                 `json:"country"`
	City        string                 `json:"city"`
	Religions   []string               `json:"religions"`
	Preferences []string               `json:"preferences"`
	
}

// MetadataUploadRequest represents metadata upload request
type MetadataUploadRequest struct {
	UserID      string                 `json:"user_id"`
	Country     string                 `json:"country"`
	City        string                 `json:"city"`
	Religions   []string               `json:"religions"`
	Preferences []string               `json:"preferences"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
}



// StoryResponse represents the API response structure
type StoryResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// GetStoryTopics handles GET request for story topics
// func (h *StoryTopics) GetStoryTopics(w http.ResponseWriter, r *http.Request) {
// 	// Verify authentication
// 	username, email, err := h.verifyAuth(r)
// 	if err != nil {
// 		h.sendErrorResponse(w, http.StatusUnauthorized, err.Error())
// 		return
// 	}

// 	var req CreateStoryRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	if req.Country == "" || req.City == "" {
// 		h.sendErrorResponse(w, http.StatusBadRequest, "Country and city are required")
// 		return
// 	}

// 	// Create context with timeout
// 	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
// 	defer cancel()

// 	helpers.NewStoryGenerationHelper(ctx, req).GenerateStoryTopics()

// 	// Get database service from context (you'll need to inject this)
// 	// For now, we'll return a mock response
// 	response := StoryResponse{
// 		Message: "Story topics fetched successfully",
// 		Data: map[string]interface{}{
// 			"username":    username,
// 			"email":       email,
// 			"country":     country,
// 			"city":        city,
// 			"religions":   religions,
// 			"preferences": preferences,
// 			"topics":      []string{"adventure", "mystery", "romance", "fantasy"},
// 		},
// 	}

// 	h.sendJSONResponse(w, http.StatusOK, response)
// }


// CreateStory handles POST request for creating a new story
// @Summary Create a new story
// @Description Create a new story with the provided details
// @Tags stories
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param story body CreateStoryRequest true "Story creation request"
// @Success 201 {object} StoryResponse
// @Failure 400 {object} StoryResponse
// @Failure 401 {object} StoryResponse
// @Failure 500 {object} StoryResponse
// @Router /story [post]
// @Security BearerAuth
func (h *Story) CreateStory(w http.ResponseWriter, r *http.Request) {

		// Verify authentication
	username, email, err := h.verifyAuth(r)
	if err != nil {
		log.Printf("❌ DEBUG: Authentication failed: %v", err)
		h.sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	log.Printf("✅ DEBUG: Authentication successful - Username: %s, Email: %s", username, email)

	// Parse request body
	var req CreateStoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ DEBUG: Failed to parse request body: %v", err)
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	log.Printf("✅ DEBUG: Request body parsed successfully: %+v", req)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if h.storyDB == nil {
		log.Println("🔧 Initializing database service...")
		h.storyDB = database.NewStoryDatabase()
		if err := h.storyDB.Init(ctx); err != nil {
			log.Printf("❌ Failed to initialize database: %v", err)
			http.Error(w, "Database initialization failed", http.StatusInternalServerError)
			return
		}
		log.Println("✅ Database service initialized successfully")
	}
	if h.storageService == nil {

		// Initialize storage service
		log.Println("🔧 Initializing storage service...")
		h.storageService = database.NewStorageService("kutty_bucket")
		if err := h.storageService.Init(ctx); err != nil {
			log.Printf("❌ Failed to initialize storage service: %v", err)
			http.Error(w, "Storage initialization failed", http.StatusInternalServerError)
			return
		}
		log.Println("✅ Storage service initialized successfully")
	}
	h.storyGenerator = helpers.NewStoryGenerationHelper(h.storyDB, h.storageService)
	log.Println("✅ All services initialized successfully - ready for future requests")
	
	err = h.storyGenerator.UploadMetadata(ctx, "", username, email, &helpers.MetadataRequest{
		Country:     req.Country,
		City:        req.City,
		Religions:   req.Religions,
		Preferences: req.Preferences,
	})

	if err != nil {
		log.Printf("❌ DEBUG: UploadMetadata failed: %v", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("✅ DEBUG: UploadMetadata completed successfully")

	// service from context (you'll need to inject this)
	// For now, we'll return a mock response
	response := StoryResponse{
		Message: "Story created successfully",
	}

	h.sendJSONResponse(w, http.StatusCreated, response)
}

// Helper methods

// verifyAuth verifies the authentication token
func (h *Story) verifyAuth(r *http.Request) (string, string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", "", fmt.Errorf("Authorization header is required")
	}

	// Remove "Bearer " prefix if present
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		authHeader = authHeader[7:]
	}

	username, email, err := util.VerifyToken(authHeader)
	if err != nil {
		return "", "", fmt.Errorf("Invalid token: %v", err)
	}

	return username, email, nil
}

// sendJSONResponse sends a JSON response
func (h *Story) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends an error response
func (h *Story) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := StoryResponse{
		Error: message,
	}
	h.sendJSONResponse(w, statusCode, response)
}
