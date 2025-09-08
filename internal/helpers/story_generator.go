package helpers

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
	"log"

	"rio-go-model/configs"
	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"

	"github.com/google/uuid"
)

// StoryGenerationHelper orchestrates the complete story generation process
type StoryGenerationHelper struct {
	settings         *configs.Settings
	logger           *util.CustomLogger
	storyCreator     *StoryCreator
	imageCreator     *ImageCreator
	audioGenerator   *AudioGenerator
	dynamicPrompting *DynamicPrompting
	storyDatabase    *database.StoryDatabase
	storageService   *database.StorageService
	httpClient       *HTTPClient
}

// HTTPClient represents an HTTP client with connection pooling
type HTTPClient struct {
	// Implementation would include httpx.AsyncClient equivalent
	// For now, we'll use standard http.Client
	// You can enhance this with connection pooling libraries
}

// StoryGenerationResponse represents the response structure for story generation
type StoryGenerationResponse struct {
	StoryID   string `json:"story_id"`
	Title     string `json:"title"`
	StoryText string `json:"story_text"`
	ImageURL  string `json:"image_url"`
	AudioURL  string `json:"audio_url"`
	AudioType string `json:"audio_type"`
	Theme     string `json:"theme"`
}

// UserProfile represents user profile data
type UserProfile struct {
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Country      string   `json:"country"`
	City         string   `json:"city"`
	Religions    []string `json:"religions"`
	Preferences  []string `json:"preferences"`
	ProcessingStatus string `json:"processing_status"`
}

// MetadataRequest represents metadata upload request
type MetadataRequest struct {
	Country     string   `json:"country"`
	City        string   `json:"city"`
	Religions   []string `json:"religions"`
	Preferences []string `json:"preferences"`
}

// NewStoryGenerationHelper creates a new story generation helper
func NewStoryGenerationHelper(
	storyDB *database.StoryDatabase,
	storageService *database.StorageService,
) *StoryGenerationHelper {
	settings := configs.LoadSettings()
	logger := util.GetLogger("story.generator", settings)

	return &StoryGenerationHelper{
		settings:         settings,
		logger:           logger,
		storyCreator:     NewStoryCreator(),
		imageCreator:     NewImageCreator(),
		audioGenerator:   NewAudioGenerator(),
		dynamicPrompting: NewDynamicPrompting(),
		storyDatabase:    storyDB,
		storageService:   storageService,
		httpClient:       &HTTPClient{}, // Initialize with proper client
	}
}

// GenerateImage generates an image from a prompt using AI
func (sgh *StoryGenerationHelper) GenerateImage(prompt string) (string, error){
	sgh.logger.Infof("Generating image for prompt: %s", prompt[:min(len(prompt), 50)])

	// Add kid-friendly modifiers to the prompt
	kidFriendlyPrompt := fmt.Sprintf(
		"kid-friendly, child-safe, colorful, cute, playful, %s, suitable for children, cartoon style, soft colors, friendly characters",
		prompt,
	)
	imgResp, err := sgh.imageCreator.CreateImage(kidFriendlyPrompt)
	if err != nil {
		return "", err

	}
	return imgResp.Base64, nil
}

// StoryHelper generates a complete story with image and audio
func (sgh *StoryGenerationHelper) StoryHelper(ctx context.Context, theme, topic string, version int, kwargs map[string]interface{}) (*StoryGenerationResponse, error) {
	sgh.logger.Infof("Generating story for theme: %s, topic: %s, version: %d", theme, topic, version)

	// Generate story using StoryCreator
	var storyResponse *StoryGenerationResponse

	if version == 1 {
		response, err := sgh.storyCreator.CreateStory(theme, topic, version, kwargs)
		if err != nil {
			return nil, fmt.Errorf("failed to generate story: %v", err)
		}
		if response.Error != "" {
			return nil, fmt.Errorf("story generation error: %s", response.Error)
		}
		if response.Story == "" {
			return nil, fmt.Errorf("no story text generated")
		}
		storyResponse = &StoryGenerationResponse{StoryText: response.Story}
	} else {
		// Version 2 with dynamic parameters
		response, err := sgh.storyCreator.CreateStory(theme, topic, version, kwargs)
		if err != nil {
			return nil, fmt.Errorf("failed to generate story: %v", err)
		}
		if response.Error != "" {
			return nil, fmt.Errorf("story generation error: %s", response.Error)
		}
		if response.Story == "" {
			return nil, fmt.Errorf("no story text generated")
		}
		storyResponse = &StoryGenerationResponse{StoryText: response.Story}
	}

	// Generate image and audio in parallel using worker pools
	imageResultChan := make(chan struct {
		data string
		err  error
	}, 1)
	audioResultChan := make(chan struct {
		data string
		err  error
	}, 1)

	// Start image generation worker
	go func() {
		imageData, err := sgh.GenerateImage(topic)
		imageResultChan <- struct {
			data string
			err  error
		}{imageData, err}
	}()

	// Start audio generation worker
	go func() {
		audioData, err := sgh.audioGenerator.GenerateAudio(storyResponse.StoryText)
		audioResultChan <- struct {
			data string
			err  error
		}{audioData, err}
	}()

	// Wait for both operations to complete with timeout
	ctx, cancel := context.WithTimeout(ctx, sgh.settings.HuggingFaceTimeout)
	defer cancel()

	var imageData, audioData string
	var imageErr, audioErr error

	// Collect results
	for i := 0; i < 2; i++ {
		select {
		case imageResult := <-imageResultChan:
			imageData = imageResult.data
			imageErr = imageResult.err
		case audioResult := <-audioResultChan:
			audioData = audioResult.data
			audioErr = audioResult.err
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for image/audio generation")
		}
	}

	if imageErr != nil {
		sgh.logger.Errorf("Image generation error: %v", imageErr)
		return nil, fmt.Errorf("image generation failed: %v", imageErr)
	}

	if audioErr != nil {
		sgh.logger.Errorf("Audio generation error: %v", audioErr)
		return nil, fmt.Errorf("audio generation failed: %v", audioErr)
	}

	// Generate unique story ID
	storyID := uuid.New().String()

	// Upload image and audio to storage in parallel
	imageUploadChan := make(chan struct {
		url string
		err error
	}, 1)
	audioUploadChan := make(chan struct {
		url string
		err error
	}, 1)

	// Start image upload worker
	go func() {
		// Decode base64 image data
		imageBytes, err := base64.StdEncoding.DecodeString(imageData)

		if err != nil {
			imageUploadChan <- struct {
				url string
				err error
			}{"", fmt.Errorf("failed to decode image: %v", err)}
			return
		}

		url, err := sgh.storageService.UploadFile(imageBytes, "images", "png")
		imageUploadChan <- struct {
			url string
			err error
		}{url, err}
	}()

	// Start audio upload worker
	go func() {
		// Decode base64 audio data
		audioBytes, err := base64.StdEncoding.DecodeString(audioData)
		if err != nil {
			audioUploadChan <- struct {
				url string
				err error
			}{"", fmt.Errorf("failed to decode audio: %v", err)}
			return
		}

		url, err := sgh.storageService.UploadFile(audioBytes, "audio", "wav")
		audioUploadChan <- struct {
			url string
			err error
		}{url, err}
	}()

	// Wait for uploads to complete
	var imageURL, audioURL string
	var uploadErr error

	for i := 0; i < 2; i++ {
		select {
		case imageResult := <-imageUploadChan:
			if imageResult.err != nil {
				uploadErr = imageResult.err
			} else {
				imageURL = imageResult.url
			}
		case audioResult := <-audioUploadChan:
			if audioResult.err != nil {
				uploadErr = audioResult.err
			} else {
				audioURL = audioResult.url
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for file uploads")
		}
	}

	if uploadErr != nil {
		sgh.logger.Errorf("Upload error: %v", uploadErr)
		return nil, fmt.Errorf("file upload failed: %v", uploadErr)
	}

	// Generate signed URLs for frontend access
	imageSignedURL, err := sgh.storageService.GenerateSignedURL(imageURL, 24*time.Hour)
	if err != nil {
		sgh.logger.Errorf("Failed to generate signed URL for image: %v", err)
		return nil, fmt.Errorf("failed to generate image URL: %v", err)
	}

	audioSignedURL, err := sgh.storageService.GenerateSignedURL(audioURL, 24*time.Hour)
	if err != nil {
		sgh.logger.Errorf("Failed to generate signed URL for audio: %v", err)
		return nil, fmt.Errorf("failed to generate audio URL: %v", err)
	}

	// Prepare response data
	responseData := &StoryGenerationResponse{
		StoryID:   storyID,
		Title:     topic,
		StoryText: storyResponse.StoryText,
		ImageURL:  imageSignedURL,
		AudioURL:  audioSignedURL,
		AudioType: "wav",
		Theme:     theme,
	}

	// Save to database (non-blocking)
	go func() {
		dbData := map[string]interface{}{
			"story_id":   storyID,
			"title":      topic,
			"story_text": storyResponse.StoryText,
			"image_url":  imageURL,
			"audio_url":  audioURL,
			"audio_type": "wav",
			"theme":      theme,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		if version == 1 {
			_, err = sgh.storyDatabase.CreateStory(ctx, dbData)
		} else {
			_, err = sgh.storyDatabase.CreateStoryV2(ctx, dbData)
		}

		if err != nil {
			sgh.logger.Errorf("Database save error: %v", err)
		} else {
			sgh.logger.Infof("Story saved to database with ID: %s", storyID)
		}
	}()

	log.Println("Story response: %v", storyResponse)
	return responseData, nil
}

// UploadMetadata handles metadata upload and triggers background processing
func (sgh *StoryGenerationHelper) UploadMetadata(ctx context.Context, token, username, email string, metadata *MetadataRequest) error {
	sgh.logger.Infof("Uploading metadata for user: %s", email)


	//Check if user profile exists
	userProfile, err := sgh.storyDatabase.GetUserProfile(ctx, username, email)
	if err != nil || userProfile == nil {
		// Create user profile
		profileData := map[string]interface{}{
			"username":          username,
			"email":             email,
			"country":           metadata.Country,
			"city":              metadata.City,
			"religions":         metadata.Religions,
			"preferences":       metadata.Preferences,
			"processing_status": "in_progress",
		}

		_, err := sgh.storyDatabase.CreateUserProfile(ctx, profileData)
		if err != nil {
			return fmt.Errorf("failed to create user profile: %v", err)
		}
	}

	// Start background processing (non-blocking)
	go sgh.runBackgroundTasks(email, metadata)

	return nil
}

// IsUserProfileExists checks if a user profile exists
func (sgh *StoryGenerationHelper) IsUserProfileExists(ctx context.Context, username, email string) (bool, error) {
	userProfile, err := sgh.storyDatabase.GetUserProfile(ctx, username, email)
	if err != nil {
		return false, err
	}
	return userProfile != nil, nil
}

// runBackgroundTasks processes metadata for all themes in parallel
func (sgh *StoryGenerationHelper) runBackgroundTasks(email string, metadata *MetadataRequest) {
	sgh.logger.Infof("Starting background tasks for user: %s", email)

	ctx := context.Background()
	
	// Create a wait group to track all background tasks
	var wg sync.WaitGroup
	wg.Add(3)

	//Process theme 1
	go func() {
		defer wg.Done()
		if err := sgh.getDynamicPromptingTheme1(ctx, metadata.Country, metadata.City, metadata.Preferences); err != nil {
			sgh.logger.Errorf("Theme 1 processing error: %v", err)
		}
	}()

	Process theme 2
	go func() {
		defer wg.Done()
		if err := sgh.getDynamicPromptingTheme2(ctx, metadata.Country, metadata.Religions, metadata.Preferences); err != nil {
			sgh.logger.Errorf("Theme 2 processing error: %v", err)
		}
	}()

	// Process theme 3
	go func() {
		defer wg.Done()
		if err := sgh.getDynamicPromptingTheme3(ctx, metadata.Preferences); err != nil {
			sgh.logger.Errorf("Theme 3 processing error: %v", err)
		}
	}()

	// Wait for all tasks to complete
	wg.Wait()

	// Update user profile status
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := sgh.storyDatabase.UpdateUserProfile(ctx, email, map[string]interface{}{
		"processing_status": "completed",
	})
	if err != nil {
		sgh.logger.Errorf("Failed to update user profile status: %v", err)
		// Try to set failed status
		sgh.storyDatabase.UpdateUserProfile(ctx, email, map[string]interface{}{
			"processing_status": "failed",
		})
	} else {
		sgh.logger.Infof("Background tasks completed for user: %s", email)
	}
}

// getDynamicPromptingTheme1 processes theme 1 with parallel story generation
func (sgh *StoryGenerationHelper) getDynamicPromptingTheme1(ctx context.Context, country, city string, preferences []string) error {
	sgh.logger.Infof("Starting theme 1 processing for country %s and city %s", country, city)

	// Check if topics already exist
	existing, err := sgh.storyDatabase.ReadMDTopics1(ctx, country, city, preferences)
	if err != nil {
		return fmt.Errorf("error checking existing topics: %v", err)
	}
	if existing != nil {
		sgh.logger.Infof("Topics already exist for theme 1")
		return nil
	}

	// Generate prompt
	prompt, err := sgh.dynamicPrompting.GetPlanetProtectorsStories(country, city, preferences)
	if err != nil {
		return fmt.Errorf("failed to generate prompt: %v", err)
	}

	// Create topics
	topicsResponse, err := sgh.storyCreator.CreateTopics(prompt)
	if err != nil {
		return fmt.Errorf("failed to create topics: %v", err)
	}
	if topicsResponse.Error != "" {
		return fmt.Errorf("topics creation error: %s", topicsResponse.Error)
	}

	topics := topicsResponse.Title
	sgh.logger.Infof("Generated %d topics for theme 1", len(topics))

	// Save topics to database
	_, err = sgh.storyDatabase.CreateMDTopics1(ctx, country, city, preferences, topics)
	if err != nil {
		return fmt.Errorf("failed to save topics: %v", err)
	}

	// Generate stories for first few topics in parallel
	storiesPerTheme := sgh.settings.StoriesPerTheme
	if len(topics) < storiesPerTheme {
		storiesPerTheme = len(topics)
	}

	var wg sync.WaitGroup
	for i := 0; i < storiesPerTheme; i++ {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()
			kwargs := map[string]interface{}{
				"country":     country,
				"city":        city,
				"preferences": preferences,
			}
			_, err := sgh.StoryHelper(ctx, "1", topic, 2, kwargs)
			if err != nil {
				sgh.logger.Errorf("Failed to generate story for topic %s: %v", topic, err)
			}
		}(topics[i])
	}

	wg.Wait()
	sgh.logger.Infof("Completed theme 1 processing")
	return nil
}

// getDynamicPromptingTheme2 processes theme 2 with parallel story generation
func (sgh *StoryGenerationHelper) getDynamicPromptingTheme2(ctx context.Context, country string, religions, preferences []string) error {
	sgh.logger.Infof("Starting theme 2 processing for country %s and religions %v", country, religions)

	// Check if topics already exist
	existing, err := sgh.storyDatabase.ReadMDTopics2(ctx, country, religions, preferences)
	if err != nil {
		return fmt.Errorf("error checking existing topics: %v", err)
	}
	if existing != nil {
		sgh.logger.Infof("Topics already exist for theme 2")
		return nil
	}

	// Generate prompt
	prompt, err := sgh.dynamicPrompting.GetMindfulStories(country, religions[0], preferences) // Use first religion
	if err != nil {
		return fmt.Errorf("failed to generate prompt: %v", err)
	}

	// Create topics
	topicsResponse, err := sgh.storyCreator.CreateTopics(prompt)
	if err != nil {
		return fmt.Errorf("failed to create topics: %v", err)
	}
	if topicsResponse.Error != "" {
		return fmt.Errorf("topics creation error: %s", topicsResponse.Error)
	}

	topics := topicsResponse.Title
	sgh.logger.Infof("Generated %d topics for theme 2", len(topics))

	// Save topics to database
	_, err = sgh.storyDatabase.CreateMDTopics2(ctx, country, religions, preferences, topics)
	if err != nil {
		return fmt.Errorf("failed to save topics: %v", err)
	}

	// Generate stories for first few topics in parallel
	storiesPerTheme := sgh.settings.StoriesPerTheme
	if len(topics) < storiesPerTheme {
		storiesPerTheme = len(topics)
	}

	var wg sync.WaitGroup
	for i := 0; i < storiesPerTheme; i++ {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()
			kwargs := map[string]interface{}{
				"country":     country,
				"religions":   religions,
				"preferences": preferences,
			}
			_, err := sgh.StoryHelper(ctx, "2", topic, 2, kwargs)
			if err != nil {
				sgh.logger.Errorf("Failed to generate story for topic %s: %v", topic, err)
			}
		}(topics[i])
	}

	wg.Wait()
	sgh.logger.Infof("Completed theme 2 processing")
	return nil
}

// getDynamicPromptingTheme3 processes theme 3 with parallel story generation
func (sgh *StoryGenerationHelper) getDynamicPromptingTheme3(ctx context.Context, preferences []string) error {
	sgh.logger.Infof("Starting theme 3 processing for preferences %v", preferences)

	// Check if topics already exist
	existing, err := sgh.storyDatabase.ReadMDTopics3(ctx, preferences)
	if err != nil {
		return fmt.Errorf("error checking existing topics: %v", err)
	}
	if existing != nil {
		sgh.logger.Infof("Topics already exist for theme 3")
		return nil
	}

	// Generate prompt
	prompt, err := sgh.dynamicPrompting.GetChillStories(preferences)
	if err != nil {
		return fmt.Errorf("failed to generate prompt: %v", err)
	}

	// Create topics
	topicsResponse, err := sgh.storyCreator.CreateTopics(prompt)
	if err != nil {
		return fmt.Errorf("failed to create topics: %v", err)
	}
	if topicsResponse.Error != "" {
		return fmt.Errorf("topics creation error: %s", topicsResponse.Error)
	}

	topics := topicsResponse.Title
	sgh.logger.Infof("Generated %d topics for theme 3", len(topics))

	// Save topics to database
	_, err = sgh.storyDatabase.CreateMDTopics3(ctx, preferences, topics)
	if err != nil {
		return fmt.Errorf("failed to save topics: %v", err)
	}

	// Generate stories for first few topics in parallel
	storiesPerTheme := sgh.settings.StoriesPerTheme
	if len(topics) < storiesPerTheme {
		storiesPerTheme = len(topics)
	}

	var wg sync.WaitGroup
	for i := 0; i < storiesPerTheme; i++ {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()
			kwargs := map[string]interface{}{
				"preferences": preferences,
			}
			_, err := sgh.StoryHelper(ctx, "3", topic, 2, kwargs)
			if err != nil {
				sgh.logger.Errorf("Failed to generate story for topic %s: %v", topic, err)
			}
		}(topics[i])
	}

	wg.Wait()
	sgh.logger.Infof("Completed theme 3 processing")
	return nil
}


