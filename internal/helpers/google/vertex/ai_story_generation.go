package vertex

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"rio-go-model/internal/util"
	"google.golang.org/genai"
	// "google.golang.org/api/option"
)

type VertexStoryGenerationHelper struct {
	client        *genai.Client
	logger        *log.Logger
	projectID     string
	location      string
	modelName     string
}

type StoryResponse struct {
	Story string `json:"story,omitempty"`
	Error string `json:"error,omitempty"`
}

func NewVertexStoryGenerationHelper() *VertexStoryGenerationHelper {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Println("Warning: GOOGLE_CLOUD_PROJECT not set")
	}

	location := os.Getenv("GOOGLE_CLOUD_LOCATION")
	if location == "" {
		location = "us-central1" // Default location
	}

	return &VertexStoryGenerationHelper{
		projectID:     projectID,
		location:      location,
		modelName:     "gemini-1.5-flash",
		logger:        log.New(os.Stdout, "VertexAI: ", log.LstdFlags),
	}
}

// Init initializes the Vertex AI client
func (s *VertexStoryGenerationHelper) Init(ctx context.Context) error {
	s.logger.Println("Initializing Vertex AI Helper")

	if s.projectID == "" {
		return fmt.Errorf("GOOGLE_CLOUD_PROJECT environment variable is required")
	}

	var client *genai.Client
	var err error

	// Try to use service account file first (same pattern as Firestore)
	credPath := "serviceAccount.json"
	_, err = os.Stat(credPath)
	if err == nil {
		log.Println("Using service account from file for vertex")
		client, err = genai.NewClient(ctx, &genai.ClientConfig{Backend: genai.BackendVertexAI,
		})
	} else {
		log.Println("Using deafult service account for Vertex AI")
		client, err = genai.NewClient(ctx, &genai.ClientConfig{Backend: genai.BackendVertexAI,
		})
	}

	if err != nil {
		return fmt.Errorf("failed to create Vertex AI client: %v", err)
	}

	s.client = client
	s.logger.Printf("Vertex AI client initialized for project: %s, location: %s", s.projectID, s.location)
	return nil
}

// CreateStory generates a story using Vertex AI
func (s *VertexStoryGenerationHelper) CreateStory(theme, topic string, version int, kwargs map[string]interface{}) (*StoryResponse, error) {
	s.logger.Printf("Creating story for theme: %s, topic: %s, version: %d", theme, topic, version)

	// Validate inputs
	if theme == "" {
		s.logger.Println("Warning: Theme is required but not provided")
		return &StoryResponse{
			Error: "Theme is required",
		}, nil
	}

	if topic == "" {
		s.logger.Println("Warning: No topic was selected")
		return &StoryResponse{
			Error: "No topic was selected",
		}, nil
	}

	if s.client == nil {
		return &StoryResponse{
			Error: "Vertex AI client not initialized - check GOOGLE_CLOUD_PROJECT",
		}, nil
	}

	// Generate formatted prompt
	formattedPrompt, systemMessage, err := util.GenerateFormattedPrompt(theme, topic, version, kwargs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate formatted prompt: %v", err)
	}
	
	// Set safety settings for child-friendly content
	// safetySettings := []*genai.SafetySetting{
	// 	{
	// 		Category:  genai.HarmCategoryHarassment,
	// 		Threshold: genai.BlockSome,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategoryHateSpeech,
	// 		Threshold: genai.BlockSome,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategorySexuallyExplicit,
	// 		Threshold: genai.BlockSome,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategoryDangerousContent,
	// 		Threshold: genai.BlockSome,
	// 	},
	// }

	// Create the complete prompt
	fullPrompt := fmt.Sprintf("%s\n\n%s", systemMessage, formattedPrompt)

	// Generate the story using Vertex AI
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	config := &genai.GenerateContentConfig{
		// SafetySettings: safetySettings,
		Temperature:     genai.Ptr(float32(0.8)),
		TopP:           genai.Ptr(float32(0.9)),
		TopK:           genai.Ptr(float32(40.0)),
		MaxOutputTokens: int32(2048),
	}

	resp, err := s.client.Models.GenerateContent(
		ctx,
		s.modelName,
		genai.Text(fullPrompt),
		config,
	)
	if err != nil {
		s.logger.Printf("Error generating story: %v", err)
		return nil, fmt.Errorf("failed to generate story: %v", err)
	}

	// Extract the generated text
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return &StoryResponse{
			Error: "No content generated",
		}, nil
	}

	storyText := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if part == nil {
			continue
		}
		// Prefer text if available
		if part.Text != "" {
			storyText += part.Text
			continue
		}
		// Fallback: stringify any non-text part
		storyText += fmt.Sprint(part)
	}

	storyText = strings.TrimSpace(storyText)
	if storyText == "" {
		return &StoryResponse{
			Error: "Empty story generated",
		}, nil
	}

	s.logger.Printf("Successfully generated story with %d characters", len(storyText))

	return &StoryResponse{
		Story: storyText,
	}, nil
}
