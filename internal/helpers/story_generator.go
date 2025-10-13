package helpers

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	// "log"

	"rio-go-model/configs"
	"rio-go-model/internal/helpers/google/audio"
	"rio-go-model/internal/helpers/google/gemini"
	"rio-go-model/internal/helpers/google/vertex"
	"rio-go-model/internal/helpers/huggingface"
	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"

	"rio-go-model/internal/model"

	"github.com/google/uuid"
)

// StoryGenerationHelper orchestrates the complete story generation process
type StoryGenerationHelper struct {
	settings               *configs.Settings
	logger                 *util.CustomLogger
	storyCreator           *huggingface.StoryCreator
	vertexAiStoryGenerator *vertex.VertexStoryGenerationHelper
	geminiStoryGenerator   *gemini.GeminiStoryGenerationHelper
	audioStoryGenerator    *audio.GoogleTTS
	imageCreator           *ImageCreator
	audioGenerator         *AudioGenerator
	dynamicPrompting       *DynamicPrompting
	storyDatabase          *database.StoryDatabase
	storageService         *database.StorageService
	httpClient             *HTTPClient
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

// A helper struct to associate a topic with its original key (preference or religion).
type topicWithKey struct {
	Key   string
	Topic string
}

// UserProfile represents user profile data
type UserProfile struct {
	Username         string   `json:"username"`
	Email            string   `json:"email"`
	Country          string   `json:"country"`
	City             string   `json:"city"`
	Religions        []string `json:"religions"`
	Preferences      []string `json:"preferences"`
	ProcessingStatus string   `json:"processing_status"`
}

// MetadataRequest represents metadata upload request
type MetadataRequest struct {
	Country     string   `json:"country"`
	City        string   `json:"city"`
	Religions   []string `json:"religions"`
	Preferences []string `json:"preferences"`
	Language    string   `json:"language"`
}

// NewStoryGenerationHelper creates a new story generation helper
func NewStoryGenerationHelper(
	storyDB *database.StoryDatabase,
	storageService *database.StorageService,
) *StoryGenerationHelper {
	settings := configs.GetSettings()
	logger := util.GetLogger("story.generator", settings)

	return &StoryGenerationHelper{
		settings:               settings,
		logger:                 logger,
		storyCreator:           huggingface.NewStoryCreator(),
		vertexAiStoryGenerator: vertex.NewVertexStoryGenerationHelper(),
		geminiStoryGenerator:   gemini.NewGeminiStoryGenerationHelper(),
		audioStoryGenerator:    audio.NewGoogleTTS(),
		imageCreator:           NewImageCreator(),
		audioGenerator:         NewAudioGenerator(),
		dynamicPrompting:       NewDynamicPrompting(),
		storyDatabase:          storyDB,
		storageService:         storageService,
		httpClient:             &HTTPClient{}, // Initialize with proper client
	}
}

// GenerateImage generates an image from a prompt using AI
func (sgh *StoryGenerationHelper) GenerateImage(prompt string) ([]byte, error) {
	sgh.logger.Infof("Generating image for prompt: %s", prompt[:min(len(prompt), 50)])

	// Add kid-friendly modifiers to the prompt
	kidFriendlyPrompt := fmt.Sprintf(
		"kid-friendly, child-safe, colorful, cute, playful, %s, suitable for children, cartoon style, soft colors, friendly characters",
		prompt,
	)
	imgResp, err := sgh.imageCreator.CreateImage(kidFriendlyPrompt)
	if err != nil {
		return nil, err

	}
	return imgResp.Data, nil
}

// StoryHelper generates a complete story with image and audio
func (sgh *StoryGenerationHelper) StoryHelper(ctx context.Context, theme, theme_id, topic string, kwargs map[string]interface{}) error {
	sgh.logger.Infof("Generating story for theme: %s, topic: %s, version: %d", theme, topic)

	// Generate story using StoryCreator
	var storyResponse *StoryGenerationResponse
	isSuspended, err := sgh.storyDatabase.SuspendGeminiAPI(ctx, "gemini")
	var response *model.StoryResponse
	if (err != nil || isSuspended) && kwargs["language"].(string) != "Telugu" {
		response, err = sgh.storyCreator.CreateStory(theme, topic, kwargs)
	} else {
		response, err = sgh.geminiStoryGenerator.CreateStory(theme, topic, kwargs)
	}
	// Version 2 with dynamic parameters
	if err != nil {
		return fmt.Errorf("failed to generate story: %v", err)
	}
	if response.Error != "" {
		return fmt.Errorf("story generation error: %s", response.Error)
	}
	if response.Story == "" {
		return fmt.Errorf("no story text generated")
	}
	storyResponse = &StoryGenerationResponse{StoryText: response.Story}
	// Generate image and audio in parallel using worker pools
	imageResultChan := make(chan struct {
		data []byte
		err  error
	}, 1)
	audioResultChan := make(chan struct {
		data []byte
		err  error
	}, 1)

	// Start image generation worker
	util.GoroutineWithRecovery(func() {
		imageData, err := sgh.GenerateImage(topic)
		imageResultChan <- struct {
			data []byte
			err  error
		}{imageData, err}
	})

	// Start audio generation worker
	util.GoroutineWithRecovery(func() {
		var audioData []byte
		language := kwargs["language"].(string)
		suspended, err := sgh.storyDatabase.SuspendAudioAPI(ctx, "audio")
		if suspended || err != nil {
			if err != nil {
				sgh.logger.Errorf("Failed to read audio api trigger: %v", err)
			} else {
				sgh.logger.Errorf("Google Audio API trigger is suspended; using fallback audio generator")
			}
			audioData, err = sgh.audioGenerator.GenerateAudio(storyResponse.StoryText)
		} else {
			sgh.logger.Infof("Using Google Audio API to generate story audio...")
			audioData, err = sgh.audioStoryGenerator.GenerateAudioAdapter(storyResponse.StoryText, language)
		}
		audioResultChan <- struct {
			data []byte
			err  error
		}{audioData, err}
	})

	// Wait for both operations to complete with timeout
	ctx, cancel := context.WithTimeout(ctx, sgh.settings.HuggingFaceTimeout)
	defer cancel()

	var imageData []byte
	var audioData []byte
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
			return fmt.Errorf("timeout waiting for image/audio generation")
		}
	}

	if imageErr != nil {
		sgh.logger.Errorf("Image generation error: %v", imageErr)
		return fmt.Errorf("image generation failed: %v", imageErr)
	}

	if audioErr != nil {
		sgh.logger.Errorf("Audio generation error: %v", audioErr)
		return fmt.Errorf("audio generation failed: %v", audioErr)
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
	util.GoroutineWithRecovery(func() {
		url, err := sgh.storageService.UploadFile(imageData, "images", "png")
		imageUploadChan <- struct {
			url string
			err error
		}{url, err}
	})

	// Start audio upload worker
	util.GoroutineWithRecovery(func() {
		url, err := sgh.storageService.UploadFile(audioData, "audio", "wav")
		audioUploadChan <- struct {
			url string
			err error
		}{url, err}
	})

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
			return fmt.Errorf("timeout waiting for file uploads")
		}
	}

	if uploadErr != nil {
		sgh.logger.Errorf("Upload error: %v", uploadErr)
		return fmt.Errorf("file upload failed: %v", uploadErr)
	}
	// Save to database (non-blocking)
	util.GoroutineWithRecovery(func() {
		dbData := map[string]interface{}{
			"story_id":   storyID,
			"title":      topic,
			"story_text": storyResponse.StoryText,
			"image_url":  imageURL,
			"audio_url":  audioURL,
			"audio_type": "wav",
			"theme":      theme,
			"language":   kwargs["language"].(string),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := sgh.storyDatabase.CreateStoryV2(ctx, theme_id, dbData)
		if err != nil {
			sgh.logger.Errorf("Database save error: %v", err)
		} else {
			sgh.logger.Infof("Story saved to database with ID: %s", storyID)
		}
	})

	// log.Println("Story response: %v", storyResponse)
	// return responseData, nil
	return nil
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
			"language":          metadata.Language,
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

// runBackgroundTasks processes metadata for all themes in parallel, with a global semaphore
// to limit the total number of concurrent story generations.
func (sgh *StoryGenerationHelper) runBackgroundTasks(email string, metadata *MetadataRequest) {
	sgh.logger.Infof("Starting background tasks for user: %s", email)

	ctx := context.Background()
	// Add panic recovery for the entire background task
	defer util.RecoverPanic()

	var wg sync.WaitGroup
	wg.Add(3)

	// This semaphore limits the total number of concurrent StoryHelper calls across all themes.
	// A value of 5 is a safe starting point for a 2-CPU instance.
	const maxConcurrentStories = 5
	semaphore := make(chan struct{}, maxConcurrentStories)

	// Process theme 1 in a goroutine
	util.GoroutineWithRecoveryAndHandler(func() {
		defer wg.Done()
		if err := sgh.getDynamicPromptingTheme1(ctx, metadata.Country, metadata.City, metadata.Preferences, metadata.Language, semaphore); err != nil {
			sgh.logger.Errorf("Theme 1 processing error: %v", err)
		}
	}, func(r interface{}) {
		wg.Done() // Ensure wg.Done() is called even on panic
	})

	// Process theme 2 in a goroutine
	util.GoroutineWithRecoveryAndHandler(func() {
		defer wg.Done()
		if err := sgh.getDynamicPromptingTheme2(ctx, metadata.Country, metadata.Religions, metadata.Preferences, metadata.Language, semaphore); err != nil {
			sgh.logger.Errorf("Theme 2 processing error: %v", err)
		}
	}, func(r interface{}) {
		wg.Done() // Ensure wg.Done() is called even on panic
	})

	// Process theme 3 in a goroutine
	util.GoroutineWithRecoveryAndHandler(func() {
		defer wg.Done()
		if err := sgh.getDynamicPromptingTheme3(ctx, metadata.Preferences, metadata.Language, semaphore); err != nil {
			sgh.logger.Errorf("Theme 3 processing error: %v", err)
		}
	}, func(r interface{}) {
		wg.Done() // Ensure wg.Done() is called even on panic
	})

	// Wait for all theme processing to complete
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

func (sgh *StoryGenerationHelper) TopicsGenerator(ctx context.Context, prompt string, language string) ([]string, error) {
	sgh.logger.Infof("Starting topics generator ")
	// Create topics
	isSuspended, err := sgh.storyDatabase.SuspendGeminiAPI(ctx, "gemini")
	var topicsResponse *model.TopicResponse
	if (err != nil || isSuspended) && language == "English" {
		topicsResponse, err = sgh.storyCreator.CreateTopics(prompt)
	} else {
		topicsResponse, err = sgh.geminiStoryGenerator.CreateTopics(prompt)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create topics: %v", err)
	}
	if topicsResponse.Error != "" {
		return nil, fmt.Errorf("topics creation error: %s", topicsResponse.Error)
	}

	topics := topicsResponse.Title
	return topics, nil
}

// getDynamicPromptingTheme1 processes theme 1 with parallel story generation controlled by a semaphore
func (sgh *StoryGenerationHelper) getDynamicPromptingTheme1(ctx context.Context, country, city string, preferences []string, language string, semaphore chan struct{}) error {
	sgh.logger.Infof("Starting theme 1 processing for country %s and city %s", country, city)
	// Check if topics already exist
	existing, err := sgh.storyDatabase.ReadMDTopics1(ctx, country, city, preferences, language)
	if err != nil {
		return fmt.Errorf("error checking existing topics: %v", err)
	}

	if existing != nil && len(existing) >= sgh.settings.DefaultStoryToGenerate {
		sgh.logger.Infof("Topics already exist for theme 1")
		return nil
	}
	theme1_id := uuid.New().String()
	var storiesPerPreference = int(math.Round(float64(sgh.settings.DefaultStoryToGenerate) / float64(len(preferences))))
	var allTopics []topicWithKey
	var concatTopics = make(map[string][]string)
	for _, preference := range preferences {
		// Generate prompt
		prompt, err := sgh.dynamicPrompting.GetPlanetProtectorsStories(country, city, preference, language, storiesPerPreference)
		if err != nil {
			return fmt.Errorf("failed to generate prompt: %v", err)
		}

		// Create topics
		topics, err := sgh.TopicsGenerator(ctx, prompt, language)
		if err != nil {
			return fmt.Errorf("failed to create topics: %v", err)
		}
		sgh.logger.Infof("Generated %d topics for theme 1", len(topics))

		// Save topics to database
		_, err = sgh.storyDatabase.CreateMDTopics1(ctx, theme1_id, country, city, preference, topics, language)
		if err != nil {
			return fmt.Errorf("failed to save topics: %v", err)
		}

		concatTopics[preference] = append(concatTopics[preference], topics...)

	}

	for key, topics := range concatTopics {
		for _, topic := range topics {
			allTopics = append(allTopics, topicWithKey{Key: key, Topic: topic})
		}
	}

	// Generate stories for first few topics in parallel
	storiesPerTheme := sgh.settings.StoriesPerTheme
	if len(allTopics) < storiesPerTheme {
		storiesPerTheme = len(allTopics)
	}

	var wg sync.WaitGroup

	for _, topic := range allTopics[:storiesPerTheme] {
		wg.Add(1)
		util.GoroutineWithRecoveryAndHandler(func() {
			semaphore <- struct{}{}        // Acquire a spot
			defer func() { <-semaphore }() // Release the spot
			defer wg.Done()
			kwargs := map[string]interface{}{
				"country":     country,
				"city":        city,
				"preferences": topic.Key,
				"language":    language,
			}

			err := sgh.StoryHelper(ctx, "1", theme1_id, topic.Topic, kwargs)
			if err != nil {
				sgh.logger.Errorf("Failed to generate story for topic %s: %v", topic, err)
			}
		}, func(r interface{}) {
			<-semaphore // Release semaphore on panic
			wg.Done()   // Ensure wg.Done() is called even on panic
		})
	}
	wg.Wait()

	sgh.logger.Infof("Completed theme 1 processing")
	return nil
}

// getDynamicPromptingTheme2 processes theme 2 with parallel story generation controlled by a semaphore
func (sgh *StoryGenerationHelper) getDynamicPromptingTheme2(ctx context.Context, country string, religions, preferences []string, language string, semaphore chan struct{}) error {
	sgh.logger.Infof("Starting theme 2 processing for country %s and religions %v", country, religions)

	// Check if topics already exist
	existing, err := sgh.storyDatabase.ReadMDTopics2(ctx, country, religions, preferences, language)
	if err != nil {
		return fmt.Errorf("error checking existing topics: %v", err)
	}
	if existing != nil && len(existing) >= sgh.settings.DefaultStoryToGenerate {
		sgh.logger.Infof("Topics already exist for theme 2")
		return nil
	}
	theme2_id := uuid.New().String()
	// CORRECT: Initialize the map using make()
	concatTopics := make(map[string][]string)
	storiesPerPreference := int(math.Round(float64(sgh.settings.DefaultStoryToGenerate) / float64(len(religions))))

	for _, religion := range religions {
		if strings.EqualFold(religion, "any") {
			continue
		}
		prompt, err := sgh.dynamicPrompting.GetMindfulStories(country, religion, preferences, language, storiesPerPreference)
		if err != nil {
			return fmt.Errorf("failed to generate prompt: %v", err)
		}

		// Create topics
		topics, err := sgh.TopicsGenerator(ctx, prompt, language)
		if err != nil {
			return fmt.Errorf("failed to create topics: %v", err)
		}
		sgh.logger.Infof("Generated %d topics for theme 2", len(topics))

		// Save topics to database
		_, err = sgh.storyDatabase.CreateMDTopics2(ctx, theme2_id, country, religion, language, preferences, topics)
		if err != nil {
			return fmt.Errorf("failed to save topics: %v", err)
		}
		// CORRECT: Use map assignment syntax
		concatTopics[religion] = append(concatTopics[religion], topics...)
	}

	// Generate stories for first few topics in parallel
	var allTopics []topicWithKey
	for key, topics := range concatTopics {
		for _, topic := range topics {
			allTopics = append(allTopics, topicWithKey{Key: key, Topic: topic})
		}
	}

	storiesPerTheme := sgh.settings.StoriesPerTheme
	if len(allTopics) < storiesPerTheme {
		storiesPerTheme = len(allTopics)
	}

	var wg sync.WaitGroup
	for _, item := range allTopics[:storiesPerTheme] {
		wg.Add(1)
		util.GoroutineWithRecoveryAndHandler(func() {
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			defer wg.Done()

			kwargs := map[string]interface{}{
				"country":     country,
				"religions":   item.Key,
				"preferences": preferences,
				"language":    language,
			}

			err := sgh.StoryHelper(ctx, "2", theme2_id, item.Topic, kwargs)
			if err != nil {
				sgh.logger.Errorf("Failed to generate story for topic %s: %v", item.Topic, err)
			}
		}, func(r interface{}) {
			<-semaphore // Release semaphore on panic
			wg.Done()   // Ensure wg.Done() is called even on panic
		})
	}
	wg.Wait()

	sgh.logger.Infof("Completed theme 2 processing")
	return nil
}

// getDynamicPromptingTheme3 processes theme 3 with parallel story generation controlled by a semaphore
func (sgh *StoryGenerationHelper) getDynamicPromptingTheme3(ctx context.Context, preferences []string, language string, semaphore chan struct{}) error {
	sgh.logger.Infof("Starting theme 3 processing for preferences %v", preferences)

	// Check if topics already exist
	existing, err := sgh.storyDatabase.ReadMDTopics3(ctx, preferences, language)
	if err != nil {
		return fmt.Errorf("error checking existing topics: %v", err)
	}
	if existing != nil && len(existing) >= sgh.settings.DefaultStoryToGenerate {
		sgh.logger.Infof("Topics already exist for theme 3")
		return nil
	}

	theme3_id := uuid.New().String()

	var allTopics []topicWithKey
	var concatTopics = make(map[string][]string)
	var storiesPerPreference = int(math.Round(float64(sgh.settings.DefaultStoryToGenerate) / float64(len(preferences))))
	// log.Println("storiesPerPreference", storiesPerPreference)
	for _, preference := range preferences {
		// Generate prompt
		prompt, err := sgh.dynamicPrompting.GetChillStories(preference, language, storiesPerPreference)
		if err != nil {
			return fmt.Errorf("failed to generate prompt: %v", err)
		}

		// log.Println("prompt .. ", prompt)
		topics, err := sgh.TopicsGenerator(ctx, prompt, language)
		if err != nil {
			return fmt.Errorf("failed to create topics: %v", err)
		}
		sgh.logger.Infof("Generated %d topics for theme 3", len(topics))

		// log.Println("length of topics", len(topics))
		// Save topics to database
		_, err = sgh.storyDatabase.CreateMDTopics3(ctx, theme3_id, preference, language, topics)
		if err != nil {
			return fmt.Errorf("failed to save topics: %v", err)
		}
		concatTopics[preference] = append(concatTopics[preference], topics...)
	}

	// log.Println("length of concatTopics", len(concatTopics))
	for key, topics := range concatTopics {
		for _, topic := range topics {
			allTopics = append(allTopics, topicWithKey{Key: key, Topic: topic})
		}
	}
	// log.Println("length of allTopics", len(allTopics))

	// Generate stories for first few topics in parallel
	storiesPerTheme := sgh.settings.StoriesPerTheme
	if len(allTopics) < storiesPerTheme {
		storiesPerTheme = len(allTopics)
	}

	var wg sync.WaitGroup

	for _, topic := range allTopics[:storiesPerTheme] {
		wg.Add(1)
		util.GoroutineWithRecoveryAndHandler(func() {
			semaphore <- struct{}{}        // Acquire a spot
			defer func() { <-semaphore }() // Release the spot
			defer wg.Done()

			kwargs := map[string]interface{}{
				"preferences": topic.Key,
				"language":    language,
			}
			err := sgh.StoryHelper(ctx, "3", theme3_id, topic.Topic, kwargs)
			if err != nil {
				sgh.logger.Errorf("Failed to generate story for topic %s: %v", topic, err)
			}
		}, func(r interface{}) {
			<-semaphore // Release semaphore on panic
			wg.Done()   // Ensure wg.Done() is called even on panic
		})
	}
	wg.Wait()

	sgh.logger.Infof("Completed theme 3 processing")
	return nil
}
