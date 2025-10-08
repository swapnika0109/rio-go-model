package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"rio-go-model/internal/util"
	"strings"
	"time"
)

type GeminiStoryGenerationHelper struct {
	logger    *log.Logger
	apiKey    string
	modelName string
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

type StoryResponse struct {
	Story string `json:"story,omitempty"`
	Error string `json:"error,omitempty"`
}

func NewGeminiStoryGenerationHelper() *GeminiStoryGenerationHelper {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("Warning: GEMINI_API_KEY not set")
	}

	helper := &GeminiStoryGenerationHelper{
		logger:    log.New(os.Stdout, "GeminiStoryGenerationHelper: ", log.LstdFlags),
		apiKey:    apiKey,
		modelName: "gemini-2.5-flash-lite",
	}
	log.Printf("✅ Gemini HTTP helper ready (publishers/google), model=%s", helper.modelName)
	return helper
}

func (s *GeminiStoryGenerationHelper) CreateStory(theme, topic string, version int, kwargs map[string]interface{}) (*StoryResponse, error) {
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

	if s.apiKey == "" {
		return &StoryResponse{Error: "GEMINI_API_KEY not set"}, nil
	}

	// Generate formatted prompt
	formattedPrompt, systemMessage, err := util.GenerateFormattedPrompt(theme, topic, version, kwargs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate formatted prompt: %v", err)
	}
	// Build HTTP payload for Gemini API (publishers/google)
	type part struct {
		Text string `json:"text"`
	}
	type content struct {
		Role  string `json:"role"`
		Parts []part `json:"parts"`
	}
	type ThinkingConfig struct {
		ThinkingBudget int32 `json:"thinkingBudget,omitempty"`
	}
	type generationConfig struct {
		Temperature     float32        `json:"temperature,omitempty"`
		MaxOutputTokens int32          `json:"maxOutputTokens,omitempty"`
		TopP            float32        `json:"topP,omitempty"`
		TopK            float32        `json:"topK,omitempty"`
		ThinkingConfig  ThinkingConfig `json:"thinkingConfig,omitempty"`
	}
	type safetySetting struct {
		Category  string `json:"category"`
		Threshold string `json:"threshold"`
	}
	payload := struct {
		Contents         []content        `json:"contents"`
		GenerationConfig generationConfig `json:"generationConfig"`
		SafetySettings   []safetySetting  `json:"safetySettings,omitempty"`
	}{}

	fullPrompt := fmt.Sprintf("%s\n\n%s", systemMessage, formattedPrompt)
	// log.Printf("Check the fullPrompt: %v", fullPrompt)
	payload.Contents = []content{
		{
			Role:  "user",
			Parts: []part{{Text: fullPrompt}},
		},
	}
	payload.GenerationConfig = generationConfig{Temperature: 0.7, MaxOutputTokens: 3024, TopP: 0.9, TopK: 40.0, ThinkingConfig: ThinkingConfig{ThinkingBudget: -1}}
	payload.SafetySettings = []safetySetting{
		{Category: "HARM_CATEGORY_HATE_SPEECH", Threshold: "OFF"},
		{Category: "HARM_CATEGORY_DANGEROUS_CONTENT", Threshold: "OFF"},
		{Category: "HARM_CATEGORY_SEXUALLY_EXPLICIT", Threshold: "OFF"},
		{Category: "HARM_CATEGORY_HARASSMENT", Threshold: "OFF"},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		s.modelName, s.apiKey,
	)
	// s.logger.Printf("Check the url: %v", url)
	// s.logger.Printf("Making request to Gemini API")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini http error: %v", err)
	}
	defer httpResp.Body.Close()
	respBody, _ := io.ReadAll(httpResp.Body)
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		s.logger.Printf("Gemini API error status=%d body=%s", httpResp.StatusCode, string(respBody))
		return nil, fmt.Errorf("gemini api error: status=%d", httpResp.StatusCode)
	}

	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return &StoryResponse{Error: "No content generated"}, nil
	}

	storyText := ""
	for _, p := range parsed.Candidates[0].Content.Parts {
		storyText += p.Text
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
