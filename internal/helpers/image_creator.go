package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
	"time"

	// "rio-go-model/configs"
)

// StoryCreator represents a service for creating stories using AI models
type ImageCreator struct {
	logger    *log.Logger
	apiKey    string
	baseURL   string
	client    *http.Client
}

// AIRequest represents the request structure for AI model API calls
type AIImageRequest struct {
	Prompt    string `json:"prompt"`
	ResponseFormat string `json:"response_format"`
	Model string `json:"model"`
}



// TopicResponse represents the response for topic generation
type ImageResponse struct {
	Base64 string `json:"base64,omitempty"`
	Error string   `json:"error,omitempty"`
}


// NewStoryCreator creates a new StoryCreator instance
func NewImageCreator() *ImageCreator {
	apiKey := os.Getenv("HUGGINGFACE_TOKEN")
	if apiKey == "" {
		log.Println("Warning: HUGGINGFACE_TOKEN not set")
	}

	return &ImageCreator{
		logger:  log.New(log.Writer(), "[story.views] ", log.LstdFlags),
		apiKey:  apiKey,
		// baseURL: "https://api.together.xyz/v1", // Together AI endpoint
		baseURL: "https://router.huggingface.co/together/v1", // Fal AI endpoint
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// CreateTopics generates topics from a prompt using AI model
func (s *ImageCreator) CreateImage(prompt string) (*ImageResponse, error) {
	s.logger.Printf("Creating topics from prompt")

	// Prepare the request
	request := AIImageRequest{
		Model: "black-forest-labs/FLUX.1-dev",
		Prompt: prompt,
		ResponseFormat: "base64",
	}

	// Make API call
	response, err := s.makeAIRequest("/images/generations", request)
	if err != nil {
		s.logger.Printf("Error generating topics: %v", err)
		return nil, fmt.Errorf("failed to generate topics: %v", err)
	}

	// Parse response
	if response == nil {
		s.logger.Println("Warning: No completions or choices in response")
		return &ImageResponse{
			Error: "No response from model",
		}, nil
	}


	return response, nil
	
}

// makeAIRequest makes a request to the AI model API
func (s *ImageCreator) makeAIRequest(endpoint string, request AIImageRequest) (*ImageResponse, error) {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	

	
	// Try to parse as JSON first
	var jsonResponse struct {
		Data []struct {
			Base64 string `json:"b64_json"`
		} `json:"data"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	
	if err := json.Unmarshal(body, &jsonResponse); err == nil {
		// It's a JSON response
		if jsonResponse.Error != nil {
			return &ImageResponse{
				Error: jsonResponse.Error.Message,
			}, nil
		}
		
		if len(jsonResponse.Data) > 0 && jsonResponse.Data[0].Base64 != "" {
			return &ImageResponse{
				Base64: jsonResponse.Data[0].Base64,
				Error: "",
			}, nil
		}
	}
	
	// If JSON parsing failed, try treating as raw base64
	base64Image := string(body)
	
	// Validate it's actually base64 by trying to decode it
	_, err = base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return &ImageResponse{
			Error: fmt.Sprintf("Invalid base64 data: %v", err),
		}, nil
	}
	
	return &ImageResponse{
		Base64: base64Image,
		Error: "",
	}, nil

}

