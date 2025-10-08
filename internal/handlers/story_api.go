package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// "strings"
	"sync"
	"time"

	// "runtime/debug"

	"rio-go-model/internal/helpers"
	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"

	// "rio-go-model/configs"
	"rio-go-model/internal/model"
)

// StoryTopics represents the story topics handler
type Story struct {
	storyGenerator *helpers.StoryGenerationHelper
	storyDB        *database.StoryDatabase
	storageService *database.StorageService
	initMutex      sync.Mutex // Protects initialization
	isInitialized  bool       // Flag to check if services are initialized
	logger         *log.Logger
}

// NewStory creates a new story topics handler
func NewStory(storyDB *database.StoryDatabase,
	storageService *database.StorageService) *Story {
	return &Story{
		storyGenerator: nil,
		storyDB:        storyDB,
		storageService: storageService,
		logger:         log.New(log.Writer(), "[Story] ", log.LstdFlags|log.Lshortfile),
	}
}

// CreateStoryRequest represents the request body for creating stories
type CreateStoryRequest struct {
	Country     string   `json:"country"`
	City        string   `json:"city"`
	Religions   []string `json:"religions"`
	Preferences []string `json:"preferences"`
	Language    string   `json:"language"`
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

// ListStoriesRequest represents the expected query parameters for listing stories
type ListStoriesRequest struct {
	Theme string `json:"theme"`
	Limit int    `json:"limit"`
}

// StoryData represents a single story in the response
type StoryData struct {
	StoryID   string `json:"story_id"`
	Title     string `json:"title"`
	StoryText string `json:"story_text"`
	Image     string `json:"image"`
	Audio     string `json:"audio"`
	AudioType string `json:"audio_type"`
	Theme     string `json:"theme"`
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

// @Summary      Create a new story
// @Description  Creates a new story based on the user's profile and preferences.
// @Tags         Story
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        story body CreateStoryRequest true "Story creation request"
// @Success      201 {object} StoryResponse "Story created successfully"
// @Failure      401 {object} util.HttpError "Invalid or missing authorization token"
// @Failure      400 {object} util.HttpError "Invalid request body"
// @Failure      500 {object} util.HttpError "Internal server error"
// @Router       /story [post]
// CreateStory handles the creation of a new story
func (h *Story) CreateStory(w http.ResponseWriter, r *http.Request) {
	// Authentication
	// logger := h.logger
	// logger.Println("r.Header ", r.Header)
	var token string
	if !util.IsCookiesPresent(r) {
		// If cookies are not present, use the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("❌ DEBUG: Authorization header is required")
			h.sendErrorResponse(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		// Remove "Bearer " prefix if present
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			authHeader = authHeader[7:]
		}

		token = authHeader

	} else {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			h.sendErrorResponse(w, http.StatusUnauthorized, "Session token not found")
			return
		}
		token = cookie.Value
	}

	// log.Println("token ", token)
	// secretKey := configs.LoadSettings().SecretKey
	// log.Println("token secretKey ", strings.TrimSpace(secretKey))
	username, email, err := util.VerifyToken(token)
	if err != nil {
		log.Printf("❌ DEBUG: Invalid token: %v", err)
		h.sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}
	log.Printf("✅ DEBUG: Authentication successful - Username: %s, Email: %s", username, email)
	// Thread-safe lazy initialization
	// h.initMutex.Lock()
	// if !h.isInitialized {
	// 	// Use a new context for initialization; r.Context() might be cancelled
	// 	// if the client disconnects, but we want initialization to complete.
	// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Longer timeout for initialization
	// 	defer cancel()

	// 	log.Println("🔧 Lazily initializing services for the first time...")

	// 	// Initialize database service
	// 	storyDB := database.NewStoryDatabase()
	// 	if err := storyDB.Init(ctx); err != nil {
	// 		log.Printf("❌ Failed to initialize database: %v", err)
	// 		h.initMutex.Unlock() // Unlock on error
	// 		http.Error(w, "Database initialization failed", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	h.storyDB = storyDB
	// 	log.Println("✅ Database service initialized successfully")

	// 	// Initialize storage service
	// 	storageService := database.NewStorageService("kutty_bucket")
	// 	if err := storageService.Init(ctx); err != nil {
	// 		log.Printf("❌ Failed to initialize storage service: %v", err)
	// 		h.initMutex.Unlock() // Unlock on error
	// 		http.Error(w, "Storage initialization failed", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	h.storageService = storageService
	// 	log.Println("✅ Storage service initialized successfully")

	// 	// Create story generator with initialized services
	// 	h.storyGenerator = helpers.NewStoryGenerationHelper(h.storyDB, h.storageService)
	// 	h.isInitialized = true // Mark as initialized
	// 	log.Println("✅ All services initialized successfully!")
	// }
	// h.initMutex.Unlock()

	h.storyGenerator = helpers.NewStoryGenerationHelper(h.storyDB, h.storageService)
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

	err = h.storyGenerator.UploadMetadata(ctx, "", username, email, &helpers.MetadataRequest{
		Country:     req.Country,
		City:        req.City,
		Religions:   req.Religions,
		Preferences: req.Preferences,
		Language:    req.Language,
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

// ListStories handles listing generated stories for a user
// @Summary      List Generated Stories
// @Description  Lists stories for a user based on their profile and theme preference. Returns stories with signed URLs for images and audio.
// @Tags         Stories
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        theme query string false "Theme filter (1, 2, or 3)"
// @Param        limit query int false "Number of stories to return (default: 10)"
// @Success      200 {array} StoryData
// @Failure      401 {object} util.HttpError "Invalid or missing authorization token"
// @Failure      500 {object} util.HttpError "Internal server error"
// @Router       /stories [get]
func (h *Story) ListStories(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	logger.Println("Starting list_generated_stories request")

	defer util.RecoverPanicWithHandler(func(r interface{}) {
		http.Error(w, fmt.Sprintf("Error listing stories: %v", r), http.StatusInternalServerError)
	})

	username, email, err := util.VerifyAuth(r)
	if err != nil {
		logger.Printf("WARNING: Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	logger.Printf("INFO: Extracted username: %s, email: %s", username, email)

	// Get user profile data - exactly like Python: getUserProfile(username, email)
	var userProfile map[string]interface{}
	if username != "" && email != "" {
		userProfile, err = h.storyDB.GetUserProfileByEmail(r.Context(), email)
		if err != nil {
			logger.Printf("ERROR: Error getting user profile: %v", err)
		}
		logger.Printf("INFO: User profile data: %v", userProfile)
	}

	if userProfile == nil {
		logger.Println("WARNING: User profile not found")
		http.Error(w, "User profile not found", http.StatusUnauthorized)
		return
	}

	// Extract user data exactly like Python
	country := ""
	city := ""
	var preferences []string
	var religions []string

	if countryVal, ok := userProfile["country"].(string); ok {
		country = countryVal
	}
	if cityVal, ok := userProfile["city"].(string); ok {
		city = cityVal
	}
	if prefsVal, ok := userProfile["preferences"].([]interface{}); ok {
		for _, pref := range prefsVal {
			if prefStr, ok := pref.(string); ok {
				preferences = append(preferences, prefStr)
			}
		}
	}
	if relsVal, ok := userProfile["religions"].([]interface{}); ok {
		for _, rel := range relsVal {
			if relStr, ok := rel.(string); ok {
				religions = append(religions, relStr)
			}
		}
	}

	logger.Printf("INFO: User details - Country: %s, City: %s, Preferences: %v, Religions: %v", country, city, preferences, religions)

	// Get query parameters exactly like Python
	theme := r.URL.Query().Get("theme")
	limitStr := r.URL.Query().Get("limit")
	limit := 10

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err != nil {
			logger.Printf("WARNING: Invalid limit provided, using default: %d", limit)
		} else {
			limit = parsedLimit
		}
	}

	logger.Printf("INFO: Requested theme: %s, limit: %d", theme, limit)

	// Fetch theme data exactly like Python
	var themeData []map[string]interface{}
	switch theme {
	case "1":
		logger.Println("INFO: Fetching theme 1 data")
		themeData, err = h.storyDB.ReadMDTopics1(r.Context(), country, city, preferences)
		if err != nil {
			logger.Printf("ERROR: Error fetching theme 1 data: %v", err)
		}
		if len(themeData) == 0 {
			logger.Println("INFO: No theme data found, fetching theme 1 data directly")
			themeData, err = h.storyDB.InitialReadMDTopics1(r.Context())
			if err != nil {
				logger.Printf("ERROR: Error fetching theme 1 data directly: %v", err)
			}
		}
		logger.Printf("INFO: Theme 1 data: %v", themeData)
	case "2":
		logger.Println("INFO: Fetching theme 2 data")
		themeData, err = h.storyDB.ReadMDTopics2(r.Context(), country, religions, preferences)
		if err != nil {
			logger.Printf("ERROR: Error fetching theme 2 data: %v", err)
		}
		if len(themeData) == 0 {
			logger.Println("INFO: No theme data found, fetching theme 2 data directly")
			themeData, err = h.storyDB.InitialReadMDTopics2(r.Context())
			if err != nil {
				logger.Printf("ERROR: Error fetching theme 2 data directly: %v", err)
			}
		}
		logger.Printf("INFO: Theme 2 data: %v", themeData)
	case "3":
		logger.Println("INFO: Fetching theme 3 data")
		themeData, err = h.storyDB.ReadMDTopics3(r.Context(), preferences)
		if err != nil {
			logger.Printf("ERROR: Error fetching theme 3 data: %v", err)
		}
		if len(themeData) == 0 {
			logger.Println("INFO: No theme data found, fetching theme 3 data directly")
			themeData, err = h.storyDB.InitialReadMDTopics3(r.Context())
			if err != nil {
				logger.Printf("ERROR: Error fetching theme 3 data directly: %v", err)
			}
		}
		logger.Printf("INFO: Theme 3 data: %v", themeData)
	}

	var storiesData []StoryData
	var stories []map[string]interface{}

	// Exactly like Python: if theme_data is None or theme_data == []
	if len(themeData) == 0 {
		logger.Println("INFO: No theme data found, fetching stories directly")
		stories, err = h.storyDB.ListStories(r.Context(), limit, theme)
		if err != nil {
			logger.Printf("ERROR: Error fetching stories directly: %v", err)
		}
		logger.Printf("INFO: Direct stories fetch result: %v", len(stories))
	} else {
		// Exactly like Python: for theme_topic in theme_data
		for _, themeTopic := range themeData {
			if themeTopic != nil {
				id, ok := themeTopic["theme_id"]
				if ok {
					themeID, ok := id.(string)
					if !ok {
						logger.Printf("WARNING: theme_id is not a string: %v", id)
						continue
					}

					systemStories, err := h.storyDB.ListStoriesByThemeID(r.Context(), themeID, limit)
					if err != nil {
						logger.Printf("ERROR: Error fetching stories for theme_id %s: %v", themeID, err)
						continue
					}

					if len(systemStories) > 0 {
						logger.Printf("INFO: Found %d stories for theme_id: %s", len(systemStories), themeID)
						stories = append(stories, systemStories...)
					}
				}
			}
		}
	}

	// Exactly like Python: stories = [story for story in stories if story is not None]
	var filteredStories []map[string]interface{}
	for _, story := range stories {
		if story != nil {
			filteredStories = append(filteredStories, story)
		}
	}

	logger.Printf("INFO: Filtered stories count: %d", len(filteredStories))

	// Exactly like Python: if stories is None or stories == []
	if len(filteredStories) == 0 {
		logger.Println("INFO: No stories found from topics, fetching stories directly")
		directStories, err := h.storyDB.ListStories(r.Context(), limit, theme)
		if err != nil {
			logger.Printf("ERROR: Error in fallback stories fetch: %v", err)
		} else {
			filteredStories = directStories
		}
		logger.Printf("INFO: Fallback stories fetch result: %v", filteredStories)
	}

	// Process stories exactly like Python
	for _, story := range filteredStories {
		// Exactly like Python: image_blob_path = story.get('image_url', '').split('kutty_bucket/')[-1] if story.get('image_url') else None
		imageURL := ""
		audioURL := ""

		if imgVal, ok := story["image_url"].(string); ok {
			imageURL = imgVal
		}
		if audVal, ok := story["audio_url"].(string); ok {
			audioURL = audVal
		}

		var imageBlobPath = imageURL
		var audioBlobPath = audioURL

		// if imageURL != "" && strings.Contains(imageURL, "kutty_bucket/") {
		// 	parts := strings.Split(imageURL, "kutty_bucket/")
		// 	if len(parts) > 1 {
		// 		imageBlobPath = parts[1]
		// 	}
		// }

		// log.Println("imageBlobPath ", imageBlobPath)

		// if audioURL != "" && strings.Contains(audioURL, "kutty_bucket/") {
		// 	parts := strings.Split(audioURL, "kutty_bucket/")
		// 	if len(parts) > 1 {
		// 		audioBlobPath = parts[1]
		// 	}
		// }
		// log.Println("audioBlobPath ", audioBlobPath)

		logger.Printf("DEBUG: Processing story - Image path: %s, Audio path: %s", imageBlobPath, audioBlobPath)

		// Exactly like Python: image_signed_url = self.bucket.generate_signed_url(image_blob_path) if image_blob_path else None
		var imageSignedURL string
		var audioSignedURL string

		if imageBlobPath != "" {
			if signedURL, err := h.storageService.GenerateSignedURL(imageBlobPath, 3600); err == nil {
				imageSignedURL = signedURL
			}
		} else {
			imageSignedURL = imageURL
		}

		if audioBlobPath != "" {
			if signedURL, err := h.storageService.GenerateSignedURL(audioBlobPath, 3600); err == nil {
				audioSignedURL = signedURL
			}
		} else {
			audioSignedURL = audioURL
		}

		// Extract story data exactly like Python
		storyID := ""
		title := ""
		storyText := ""
		audioType := "audio/wav"
		storyTheme := ""

		if idVal, ok := story["id"].(string); ok {
			storyID = idVal
		} else if idVal, ok := story["story_id"].(string); ok {
			storyID = idVal
		}
		if titleVal, ok := story["title"].(string); ok {
			title = titleVal
		}
		if textVal, ok := story["story_text"].(string); ok {
			storyText = textVal
		}
		if audioTypeVal, ok := story["audio_type"].(string); ok {
			audioType = audioTypeVal
		}
		if themeVal, ok := story["theme"].(string); ok {
			storyTheme = themeVal
		}

		// Exactly like Python: stories_data.append({...})
		storiesData = append(storiesData, StoryData{
			StoryID:   storyID,
			Title:     title,
			StoryText: storyText,
			Image:     imageSignedURL,
			Audio:     audioSignedURL,
			AudioType: audioType,
			Theme:     storyTheme,
		})
	}

	logger.Printf("INFO: Returning %d stories", len(storiesData))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storiesData)
}

// UserProfile gets the user profile information for the authenticated user.
// @Summary      Get User Profile
// @Description  Gets the user profile information for the authenticated user.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]bool "User profile retrieved successfully"
// @Failure      401 {object} util.HttpError "Invalid or missing authorization token"
// @Failure      500 {object} util.HttpError "Internal server error"
// @Router       /user-profile [get]
func (h *Story) UserProfile(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	logger.Println("Starting user_profile request")

	username, email, err := util.VerifyAuth(r)
	if err != nil {
		logger.Printf("WARNING: Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	logger.Printf("INFO: Extracted username: %s, email: %s", username, email)
	user, err := h.storyDB.GetUserProfileByEmail(r.Context(), email)
	if err != nil {
		logger.Printf("ERROR: Error getting user profile: %v", err)
		http.Error(w, "Internal server error during user lookup", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if user == nil {
		logger.Println("WARNING: User profile not found")
		response := map[string]interface{}{"exists": false, "user": nil}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if username matches
	if userUsername, ok := user["username"].(string); !ok || userUsername != username {
		logger.Println("WARNING: User profile not found")
		response := map[string]interface{}{"exists": false, "user": (&model.UserProfile{}).FromMap(user)}
		json.NewEncoder(w).Encode(response)
		return
	}
	logger.Printf("INFO: User profile data: %v", user)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"exists": true, "user": (&model.UserProfile{}).FromMap(user)}
	json.NewEncoder(w).Encode(response)

}

// UpdateUserProfile updates the user profile for the authenticated user.
// @Summary      Update User Profile
// @Description  Updates the user profile for the authenticated user.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userProfile body model.UserProfile true "User profile request"
// @Success      200 {object} map[string]bool "User profile updated successfully"
// @Failure      401 {object} util.HttpError "Invalid or missing authorization token"
// @Failure      500 {object} util.HttpError "Internal server error"
// @Router       /user-profile [put]
func (h *Story) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	logger.Println("Starting update_user_profile request")

	_, email, err := util.VerifyAuth(r)
	if err != nil {
		logger.Printf("WARNING: Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var userProfileRequest model.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&userProfileRequest); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	if err := h.storyDB.UpdateUserProfileByEmail(r.Context(), email, userProfileRequest); err != nil {
		logger.Printf("ERROR: Error updating user profile: %v", err)
		http.Error(w, "Internal server error during user profile update", http.StatusInternalServerError)
		return
	}
	logger.Println("INFO: User profile updated successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "User profile updated successfully"})
}

// DeleteUserProfile deletes the user profile for the authenticated user.
// @Summary      Delete User Profile
// @Description  Deletes the user profile for the authenticated user.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]bool "User profile deleted successfully"
// @Failure      401 {object} util.HttpError "Invalid or missing authorization token"
// @Failure      500 {object} util.HttpError "Internal server error"
// @Router       /user-profile [delete]
func (h *Story) DeleteUserProfile(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	logger.Println("Starting delete_user_profile request")

	_, email, err := util.VerifyAuth(r)
	if err != nil {
		logger.Printf("WARNING: Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if err := h.storyDB.DeleteUserProfile(r.Context(), email); err != nil {
		logger.Printf("ERROR: Error deleting user profile: %v", err)
		http.Error(w, "Internal server error during user profile deletion", http.StatusInternalServerError)
		return
	}
	logger.Println("INFO: User profile deleted successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "User profile deleted successfully"})
}
