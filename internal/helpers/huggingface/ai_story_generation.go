package huggingface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"rio-go-model/configs"
	"rio-go-model/internal/model"
	"rio-go-model/internal/util"
)

// StoryCreator represents a service for creating stories using AI models
type StoryCreator struct {
	logger  *log.Logger
	apiKey  string
	baseURL string
	client  *http.Client
}

// AIRequest represents the request structure for AI model API calls
type AIRequest struct {
	Model       string      `json:"model"`
	Messages    []AIMessage `json:"messages"`
	Temperature float64     `json:"temperature,omitempty"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
	TopP        float64     `json:"top_p,omitempty"`
	Stream      bool        `json:"stream,omitempty"`
}

// AIMessage represents a message in the AI conversation
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIResponse represents the response structure from AI model API
type AIResponse struct {
	Choices []AIChoice `json:"choices"`
}

// AIChoice represents a choice from the AI model
type AIChoice struct {
	Message AIMessage `json:"message"`
}

// NewStoryCreator creates a new StoryCreator instance
func NewStoryCreator() *StoryCreator {
	apiKey := os.Getenv("HUGGINGFACE_TOKEN")
	if apiKey == "" {
		log.Println("Warning: HUGGINGFACE_TOKEN not set")
	}

	return &StoryCreator{
		logger: log.New(log.Writer(), "[story.views] ", log.LstdFlags),
		apiKey: apiKey,
		// baseURL: "https://api.together.xyz/v1", // Together AI endpoint
		baseURL: "https://router.huggingface.co/v1", // Fal AI endpoint
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// CreateTopics generates topics from a prompt using AI model
func (s *StoryCreator) CreateTopics(prompt string) (*model.TopicResponse, error) {
	s.logger.Printf("Creating topics from prompt ... " + prompt)

	// Prepare the request
	request := AIRequest{
		Model: "openai/gpt-oss-120b:together",
		Messages: []AIMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	// Make API call
	response, err := s.makeAIRequest("/chat/completions", request)
	if err != nil {
		s.logger.Printf("Error generating topics: %v", err)
		return nil, fmt.Errorf("failed to generate topics: %v", err)
	}

	// Parse response
	if response == nil || len(response.Choices) == 0 {
		s.logger.Println("Warning: No completions or choices in response")
		return &model.TopicResponse{
			Error: "No response from model",
		}, nil
	}

	topicsData := response.Choices[0].Message.Content
	if topicsData == "" {
		s.logger.Println("Warning: No topics data in response")
		return &model.TopicResponse{
			Error: "No response from model",
		}, nil
	}

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

// CreateStory generates a story based on theme, topic, and version
func (s *StoryCreator) CreateStory(theme, topic string, kwargs map[string]interface{}) (*model.StoryResponse, error) {
	s.logger.Printf("Creating story with huggingface for theme: %s, topic: %s", theme, topic)

	// Validate inputs
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

	// s.logger.Printf("Generated prompt: %s", formattedPrompt)

	// Prepare the request
	request := AIRequest{
		Model: "Qwen/Qwen2.5-7B-Instruct:together",
		Messages: []AIMessage{
			{
				Role:    "system",
				Content: systemMessage,
			},
			{
				Role:    "user",
				Content: formattedPrompt,
			},
		},
		Temperature: 0.9,
		MaxTokens:   1000,
		TopP:        0.9,
		Stream:      false,
	}

	// Make API call
	response, err := s.makeAIRequest("/chat/completions", request)
	if err != nil {
		s.logger.Printf("Error generating story: %v", err)
		return nil, fmt.Errorf("failed to generate story: %v", err)
	}

	// Parse response
	if response == nil || len(response.Choices) == 0 {
		s.logger.Println("Warning: No completions or choices in response")
		return &model.StoryResponse{
			Error: "No response from model",
		}, nil
	}

	story := strings.TrimSpace(response.Choices[0].Message.Content)
	s.logger.Println("Successfully generated story")

	return &model.StoryResponse{
		Story: story,
	}, nil
}

// Helper methods

// makeAIRequest makes a request to the AI model API
func (s *StoryCreator) makeAIRequest(endpoint string, request AIRequest) (*AIResponse, error) {
	// Convert request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", s.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	// Make request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// Parse response
	var aiResponse AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &aiResponse, nil
}
