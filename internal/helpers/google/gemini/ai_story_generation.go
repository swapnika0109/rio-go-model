package gemini

import (
	"context"
	"fmt"
	"log"
	"os"
	"rio-go-model/internal/model"
	"rio-go-model/internal/util"
	"strings"
	"time"

	"rio-go-model/configs"

	"google.golang.org/genai"
)

type GeminiStoryGenerationHelper struct {
	logger    *log.Logger
	apiKey    string
	modelName string
	client    *genai.Client
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GeminiStoryGenerationRequest struct {
	Model       string      `json:"model"`
	Contents    []AIMessage `json:"messages"`
	Temperature float32     `json:"temperature,omitempty"`
	MaxTokens   int32       `json:"max_tokens,omitempty"`
	TopP        float32     `json:"top_p,omitempty"`
	TopK        float32     `json:"top_k,omitempty"`
}

func NewGeminiStoryGenerationHelper() *GeminiStoryGenerationHelper {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("Warning: GEMINI_API_KEY not set")
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
		APIKey:  apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	helper := &GeminiStoryGenerationHelper{
		logger:    log.New(os.Stdout, "GeminiStoryGenerationHelper: ", log.LstdFlags),
		apiKey:    apiKey,
		modelName: "gemini-2.5-flash-lite",
		client:    client,
	}
	log.Printf("✅ Gemini HTTP helper ready (publishers/google), model=%s", helper.modelName)
	return helper
}

func (s *GeminiStoryGenerationHelper) CreateTopics(prompt string) (*model.TopicResponse, error) {
	s.logger.Printf("Creating topics for prompt")
	topicsResponse, err := s.GenerateText(prompt, "gemini-2.0-flash-lite")
	if err != nil {
		return nil, fmt.Errorf("failed to generate topics: %v", err)
	}

	topicsData := topicsResponse.Story

	// Parse topics from response
	topics := util.ParseTopics(topicsData)
	s.logger.Printf("Successfully generated %d topics", len(topics))
	settings := configs.GetSettings()

	if len(topics) != settings.DefaultStoryToGenerate {
		s.logger.Printf("Warning: Generated %d topics, expected %d", len(topics), settings.DefaultStoryToGenerate)
	}

	return &model.TopicResponse{
		Title: topics,
	}, nil
}

func (s *GeminiStoryGenerationHelper) CreateStory(theme, topic string, kwargs map[string]interface{}) (*model.StoryResponse, error) {
	s.logger.Printf("Creating story with gemini for theme: %s, topic: %s", theme, topic)

	if theme == "" {
		s.logger.Println("Warning: Theme is required but not provided")
		return &model.StoryResponse{
			Error: "Theme is required",
		}, nil
	}

	if topic == "" {
		s.logger.Println("Warning: No topic was selected")
		return &model.StoryResponse{
			Error: "No topic was selected",
		}, nil
	}

	// Generate formatted prompt
	formattedPrompt, systemMessage, err := util.GenerateFormattedPrompt(theme, topic, kwargs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate formatted prompt: %v", err)
	}

	// Create the complete prompt
	fullPrompt := fmt.Sprintf("%s\n\n%s", systemMessage, formattedPrompt)
	return s.GenerateText(fullPrompt, s.modelName)
}

func (s *GeminiStoryGenerationHelper) GenerateText(prompt string, modelName string) (*model.StoryResponse, error) {
	s.logger.Printf("Creating story for prompt")

	if s.client == nil {
		return &model.StoryResponse{
			Error: "Gemini AI client not initialized - check GOOGLE_CLOUD_PROJECT",
		}, nil
	}

	// Set safety settings for child-friendly content
	safetySettings := []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: "BLOCK_ONLY_HIGH",
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: "BLOCK_ONLY_HIGH",
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: "BLOCK_ONLY_HIGH",
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: "BLOCK_ONLY_HIGH",
		},
	}

	// Generate the story using Vertex AI
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	config := &genai.GenerateContentConfig{
		SafetySettings:  safetySettings,
		Temperature:     genai.Ptr(float32(0.8)),
		TopP:            genai.Ptr(float32(0.95)),
		TopK:            genai.Ptr(float32(40.0)),
		MaxOutputTokens: int32(2048),
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: genai.Ptr(int32(0)),
		},
	}

	resp, err := s.client.Models.GenerateContent(
		ctx,
		modelName,
		genai.Text(prompt),
		config,
	)
	if err != nil {
		s.logger.Printf("Error generating story: %v", err)
		return nil, fmt.Errorf("failed to generate story: %v", err)
	}

	// Extract the generated text
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return &model.StoryResponse{
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
		return &model.StoryResponse{
			Error: "Empty story generated",
		}, nil
	}

	s.logger.Printf("Successfully generated story with %d characters", len(storyText))

	return &model.StoryResponse{
		Story: storyText,
	}, nil
}
