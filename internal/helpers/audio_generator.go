package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// AudioGenerator represents a service for generating audio from text
type AudioGenerator struct {
	logger  *log.Logger
	apiKey  string
	baseURL string
	client  *http.Client
}

// AudioRequest represents the request structure for audio generation
type AudioRequest struct {
	Text string `json:"text"`
}

// AudioResponse represents the response structure from audio generation
type AudioResponse struct {
	Audio string `json:"audio,omitempty"`
	Error string `json:"error,omitempty"`
}

// KokoroAudioResponse represents the response structure from Kokoro TTS
type KokoroAudioResponse struct {
	Audio []byte `json:"audio"`
}

// NewAudioGenerator creates a new AudioGenerator instance
func NewAudioGenerator() *AudioGenerator {
	apiKey := os.Getenv("HUGGINGFACE_TOKEN")
	if apiKey == "" {
		log.Println("Warning: HUGGINGFACE_TOKEN not set")
	}

	return &AudioGenerator{
		logger:  log.New(log.Writer(), "[audio.generator] ", log.LstdFlags),
		apiKey:  apiKey,
		baseURL: "https://router.huggingface.co/fal-ai/fal-ai/kokoro/american-english", // fal.ai endpoint
		client: &http.Client{
			Timeout: 120 * time.Second, // Longer timeout for audio generation
		},
	}
}

// GenerateAudio generates audio from text using fal.ai (Kokoro model)
func (a *AudioGenerator) GenerateAudio(prompt string) (string, error) {
	a.logger.Printf("Generating audio from prompt: %s", prompt[:min(len(prompt), 100)])

	// Clean the prompt
	tempPrompt := strings.ReplaceAll(prompt, "\n", "")
	
	// Check cache first (you can implement caching later)
	// cached := cacheHelper.GetCache(prompt)
	// if cached != nil {
	//     return base64.StdEncoding.EncodeToString(cached), nil
	// }

	// Generate audio using fal.ai API
	audio, err := a.generateAudioFalAI(tempPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate audio: %v", err)
	}

	// Cache the audio (you can implement caching later)
	// cacheHelper.SetCache(prompt, audio)

	// Convert to base64
	audioBase64 := base64.StdEncoding.EncodeToString(audio)
	a.logger.Printf("Successfully generated audio, base64 length: %d", len(audioBase64))
	
	return audioBase64, nil
}

// GenerateAudioOpenVoice generates audio using local OpenVoice TTS service
func (a *AudioGenerator) GenerateAudioOpenVoice(prompt string) (string, error) {
	a.logger.Printf("Generating audio using OpenVoice for prompt: %s", prompt[:min(len(prompt), 100)])

	// Clean the prompt
	tempPrompt := strings.ReplaceAll(prompt, "\n", "")
	
	// Check cache first (you can implement caching later)
	// cached := cacheHelper.GetCache(prompt)
	// if cached != nil {
	//     return base64.StdEncoding.EncodeToString(cached), nil
	// }

	// Make request to local OpenVoice service
	url := "http://localhost:8001/tts"
	request := AudioRequest{
		Text: tempPrompt,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request to OpenVoice: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenVoice service returned status: %d - %s", resp.StatusCode, resp.Status)
	}

	// Parse response
	var audioResponse AudioResponse
	if err := json.NewDecoder(resp.Body).Decode(&audioResponse); err != nil {
		return "", fmt.Errorf("failed to decode OpenVoice response: %v", err)
	}

	// Check for error in response
	if audioResponse.Error != "" {
		return "", fmt.Errorf("OpenVoice service error: %s", audioResponse.Error)
	}

	// Check if audio data is present
	if audioResponse.Audio == "" {
		return "", fmt.Errorf("no audio data in OpenVoice response")
	}

	// Cache the audio (you can implement caching later)
	// cacheHelper.SetCache(prompt, audioResponse.Audio)

	a.logger.Printf("Successfully generated audio using OpenVoice, base64 length: %d", len(audioResponse.Audio))
	return audioResponse.Audio, nil
}

// generateAudioFalAI generates audio using fal.ai API (Kokoro model)
func (a *AudioGenerator) generateAudioFalAI(text string) ([]byte, error) {
	// Prepare the request for fal.ai
	request := map[string]interface{}{
		"version": "latest",
		"text": text,
		// "input": map[string]interface{}{
		// 	"text": text,
		// },
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", a.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	// Make request
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to fal.ai: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		// Read error response body for debugging
		errorBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fal.ai API returned status: %d - %s. Error: %s", resp.StatusCode, resp.Status, string(errorBody))
	}

	// Read the audio data
	audioData, err := io.ReadAll(resp.Body)
	
	if err != nil {
		return nil, fmt.Errorf("failed to read audio response: %v", err)
	}

	return audioData, nil
}

// GenerateAudioWithModel generates audio using a specific model
func (a *AudioGenerator) GenerateAudioWithModel(prompt, model string) (string, error) {
	a.logger.Printf("Generating audio with model %s for prompt: %s", model, prompt[:min(len(prompt), 100)])

	// Clean the prompt
	tempPrompt := strings.ReplaceAll(prompt, "\n", "")
	
	// Choose generation method based on model
	switch model {
	case "kokoro", "hexgrad/Kokoro-82M":
		return a.GenerateAudio(tempPrompt)
	case "openvoice", "local-openvoice":
		return a.GenerateAudioOpenVoice(tempPrompt)
	default:
		return "", fmt.Errorf("unsupported model: %s", model)
	}
}



// checkOpenVoiceHealth checks if the local OpenVoice service is running
func (a *AudioGenerator) checkOpenVoiceHealth() string {
	url := "http://localhost:8001/health"
	
	// Create a client with shorter timeout for health check
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "unhealthy"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "healthy"
	}

	return "unhealthy"
}

// SaveAudioToFile saves base64 audio data to a file
func (a *AudioGenerator) SaveAudioToFile(audioBase64, filepath string) error {
	// Decode base64 audio
	audioData, err := base64.StdEncoding.DecodeString(audioBase64)
	if err != nil {
		return fmt.Errorf("failed to decode base64 audio: %v", err)
	}

	// Create directory if it doesn't exist
	dir := filepath[:strings.LastIndex(filepath, "/")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write audio data to file
	if err := os.WriteFile(filepath, audioData, 0644); err != nil {
		return fmt.Errorf("failed to write audio file: %v", err)
	}

	a.logger.Printf("Audio saved to file: %s", filepath)
	return nil
}

// GetAudioInfo returns information about the audio data
func (a *AudioGenerator) GetAudioInfo(audioBase64 string) map[string]interface{} {
	audioData, err := base64.StdEncoding.DecodeString(audioBase64)
	if err != nil {
		return map[string]interface{}{
			"error": "Invalid base64 audio data",
		}
	}

	return map[string]interface{}{
		"size_bytes":        len(audioData),
		"size_kb":           len(audioData) / 1024,
		"base64_length":     len(audioBase64),
		"format":            "unknown", // You can add format detection logic
		"estimated_duration": "unknown", // You can add duration calculation logic
	}
}

// Helper functions

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
