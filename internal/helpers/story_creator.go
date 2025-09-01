package helpers

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
)

// StoryCreator represents a service for creating stories using AI models
type StoryCreator struct {
	logger    *log.Logger
	apiKey    string
	baseURL   string
	client    *http.Client
}

// AIRequest represents the request structure for AI model API calls
type AIRequest struct {
	Model       string        `json:"model"`
	Messages    []AIMessage   `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
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

// TopicResponse represents the response for topic generation
type TopicResponse struct {
	Title []string `json:"title,omitempty"`
	Error string   `json:"error,omitempty"`
}

// StoryResponse represents the response for story generation
type StoryResponse struct {
	Story string `json:"story,omitempty"`
	Error string `json:"error,omitempty"`
}

// NewStoryCreator creates a new StoryCreator instance
func NewStoryCreator() *StoryCreator {
	apiKey := os.Getenv("HUGGINGFACE_TOKEN")
	if apiKey == "" {
		log.Println("Warning: HUGGINGFACE_TOKEN not set")
	}

	return &StoryCreator{
		logger:  log.New(log.Writer(), "[story.views] ", log.LstdFlags),
		apiKey:  apiKey,
		baseURL: "https://api.together.xyz/v1", // Together AI endpoint
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// CreateTopics generates topics from a prompt using AI model
func (s *StoryCreator) CreateTopics(prompt string) (*TopicResponse, error) {
	s.logger.Printf("Creating topics from prompt")

	// Prepare the request
	request := AIRequest{
		Model: "mistralai/Mixtral-8x7B-Instruct-v0.1",
		Messages: []AIMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
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
		return &TopicResponse{
			Error: "No response from model",
		}, nil
	}

	topicsData := response.Choices[0].Message.Content
	if topicsData == "" {
		s.logger.Println("Warning: No topics data in response")
		return &TopicResponse{
			Error: "No response from model",
		}, nil
	}

	// Parse topics from response
	topics := s.parseTopics(topicsData)
	s.logger.Printf("Successfully generated %d topics", len(topics))

	return &TopicResponse{
		Title: topics,
	}, nil
}

// CreateStory generates a story based on theme, topic, and version
func (s *StoryCreator) CreateStory(theme, topic string, version int, kwargs map[string]interface{}) (*StoryResponse, error) {
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

	// Generate formatted prompt
	formattedPrompt, systemMessage, err := s.generateFormattedPrompt(theme, topic, version, kwargs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate formatted prompt: %v", err)
	}

	s.logger.Printf("Generated prompt: %s", formattedPrompt)

	// Prepare the request
	request := AIRequest{
		Model: "mistralai/Mixtral-8x7B-Instruct-v0.1",
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
		MaxTokens:  1000,
		TopP:       0.9,
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
		return &StoryResponse{
			Error: "No response from model",
		}, nil
	}

	story := strings.TrimSpace(response.Choices[0].Message.Content)
	s.logger.Println("Successfully generated story")

	return &StoryResponse{
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

// parseTopics parses topics from the AI response
func (s *StoryCreator) parseTopics(topicsData string) []string {
	topicsList := strings.Split(topicsData, "\n")
	topics := make([]string, 0)

	for _, topic := range topicsList {
		topic = strings.TrimSpace(topic)
		if topic == "" {
			continue
		}

		// Handle quoted topics
		if strings.Contains(topic, `"`) {
			parts := strings.Split(topic, `"`)
			if len(parts) > 1 {
				topics = append(topics, strings.TrimSpace(parts[1]))
			}
		} else {
			topics = append(topics, topic)
		}
	}

	return topics
}

// generateFormattedPrompt generates the formatted prompt based on version and parameters
func (s *StoryCreator) generateFormattedPrompt(theme, topic string, version int, kwargs map[string]interface{}) (string, string, error) {
	var formattedPrompt, systemMessage string

	if version == 2 {
		// Dynamic prompt for version 2
		promptTemplate, err := s.getDynamicPromptConfig(theme)
		if err != nil {
			return "", "", err
		}

		// Extract parameters
		country := s.getStringFromMap(kwargs, "country", "")
		city := s.getStringFromMap(kwargs, "city", "")
		religions := s.getStringSliceFromMap(kwargs, "religions")
		preferences := s.getStringSliceFromMap(kwargs, "preferences")

		religionsStr := strings.Join(religions, ", ")

		// Format the prompt
		formattedPrompt = fmt.Sprintf(promptTemplate.Prompt, topic, country, city, religionsStr, strings.Join(preferences, ", "))

		// Add preference-specific content
		for _, preference := range preferences {
			if prefContent := s.getPreferenceContent(strings.ToUpper(preference)); prefContent != "" {
				formattedPrompt += prefContent
			}
		}

		systemMessage = promptTemplate.System
	} else {
		// Standard prompt for version 1
		promptTemplate, err := s.getPromptConfig(theme)
		if err != nil {
			return "", "", err
		}

		formattedPrompt = fmt.Sprintf(promptTemplate.Prompt, topic)
		systemMessage = promptTemplate.System
	}

	return formattedPrompt, systemMessage, nil
}

// Configuration structures

// PromptTemplate represents a prompt template configuration
type PromptTemplate struct {
	Prompt string `json:"prompt"`
	System string `json:"system"`
}

// getDynamicPromptConfig gets the dynamic prompt configuration for a theme
func (s *StoryCreator) getDynamicPromptConfig(theme string) (*PromptTemplate, error) {
	// Load settings to get the actual prompt configuration
	settings := configs.LoadSettings()
	
	if config, exists := settings.GetDynamicPromptConfig(theme); exists {
		return &PromptTemplate{
			Prompt: config.Prompt,
			System: config.System,
		}, nil
	}
	
	// Fallback to default if theme not found
	return &PromptTemplate{
		Prompt: "Create a story about %s set in %s, %s. Include elements related to %s and incorporate these preferences: %s.",
		System: "You are a creative storyteller who creates engaging stories for children.",
	}, nil
}

// getPromptConfig gets the standard prompt configuration for a theme
func (s *StoryCreator) getPromptConfig(theme string) (*PromptTemplate, error) {
	// Load settings to get the actual prompt configuration
	settings := configs.LoadSettings()
	
	if config, exists := settings.GetPromptConfig(theme); exists {
		return &PromptTemplate{
			Prompt: config.Prompt,
			System: config.System,
		}, nil
	}
	
	// Fallback to default if theme not found
	return &PromptTemplate{
		Prompt: "Create a story about %s.",
		System: "You are a creative storyteller who creates engaging stories for children.",
	}, nil
}

// getPreferenceContent gets the content for a specific preference
func (s *StoryCreator) getPreferenceContent(preference string) string {
	// Load settings to get the actual preference configuration
	settings := configs.LoadSettings()
	
	if content, exists := settings.GetPreference(preference); exists {
		return content
	}

	// Fallback preferences if not found in settings
	fallbackPreferences := map[string]string{
		"NATURE":     " Focus on environmental themes and natural elements.",
		"ADVENTURE":  " Include exciting adventures and challenges.",
		"FRIENDSHIP": " Emphasize friendship and cooperation.",
		"LEARNING":   " Include educational elements and life lessons.",
	}

	if content, exists := fallbackPreferences[preference]; exists {
		return content
	}

	return ""
}

// Helper functions for map operations

// getStringFromMap safely extracts a string value from a map
func (s *StoryCreator) getStringFromMap(m map[string]interface{}, key, defaultValue string) string {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

// getStringSliceFromMap safely extracts a string slice from a map
func (s *StoryCreator) getStringSliceFromMap(m map[string]interface{}, key string) []string {
	if value, exists := m[key]; exists {
		if slice, ok := value.([]string); ok {
			return slice
		}
		// Handle interface{} slice
		if interfaceSlice, ok := value.([]interface{}); ok {
			result := make([]string, 0)
			for _, item := range interfaceSlice {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return []string{}
}
