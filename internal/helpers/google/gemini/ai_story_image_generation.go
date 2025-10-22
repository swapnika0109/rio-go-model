package gemini

import (
	"context"
	"fmt"
	"log"
	"os"
	"rio-go-model/internal/model"
	"time"

	"google.golang.org/genai"
)

type GeminiImageGenerationHelper struct {
	logger    *log.Logger
	apiKey    string
	modelName string
	client    *genai.Client
}

func NewGeminiImageGenerationHelper() *GeminiImageGenerationHelper {
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
	// storyCharacters := tokens.NewStoryCharacters()
	helper := &GeminiImageGenerationHelper{
		logger:    log.New(os.Stdout, "GeminiStoryGenerationHelper: ", log.LstdFlags),
		apiKey:    apiKey,
		modelName: "gemini-2.5-flash-image",
		client:    client,
	}
	log.Printf("✅ Gemini HTTP helper ready (publishers/google), model=%s", helper.modelName)
	return helper
}

func (s *GeminiImageGenerationHelper) CreateTopicsImage(prompt string) (imageBytes []byte, err error) {
	s.logger.Printf("Creating image for prompt")
	imageResponse := s.GenerateImage(prompt, "gemini-2.5-flash-image")
	if imageResponse.Error != "" {
		return nil, fmt.Errorf("failed to generate image: %v", imageResponse.Error)
	}
	return imageResponse.Image, nil
}

func (s *GeminiImageGenerationHelper) GenerateImage(prompt string, modelName string) *model.ImageResponse {
	s.logger.Printf("Creating image for prompt")

	if s.client == nil {
		return &model.ImageResponse{
			Error: "Gemini AI client not initialized - check GOOGLE_CLOUD_PROJECT",
		}
	}

	// Generate the story using Vertex AI
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{"Image"},
		ImageConfig: &genai.ImageConfig{
			AspectRatio: "1:1",
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
		return &model.ImageResponse{
			Error: "failed to generate story: %v",
		}
	}

	// Extract the generated text
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return &model.ImageResponse{
			Error: "No content generated",
		}
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			s.logger.Println(part.Text)
		} else if part.InlineData != nil {
			imageBytes := part.InlineData.Data
			s.logger.Printf("Image generated successfully")
			return &model.ImageResponse{
				Image: imageBytes,
			}
		}
	}

	return &model.ImageResponse{
		Error: "No image generated",
	}
}
