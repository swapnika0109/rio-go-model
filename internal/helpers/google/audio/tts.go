package audio

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"google.golang.org/api/option"
)

type GoogleTTS struct {
	Client *texttospeech.Client
	Logger *log.Logger
}

type GoogleTTSRequest struct {
	Text         string
	LanguageName string
	LanguageCode string
	SSML         string
}

type GoogleTTSResponse struct {
	AudioContent []byte
	AudioFormat  string
	Error        string
}

func NewGoogleTTS() *GoogleTTS {
	ctx := context.Background()
	var client *texttospeech.Client
	var err error
	// Try to use service account file first
	credPath := "serviceAccount.json"
	if _, statErr := os.Stat(credPath); statErr == nil {
		log.Println("Using service account from file for texttospeech")
		client, err = texttospeech.NewClient(ctx, option.WithCredentialsFile(credPath))
		if err != nil {
			log.Fatalf("Failed to create texttospeech client: %v", err)
		}
	} else {
		log.Println("Using default credentials for texttospeech")
		// In Cloud Run, use the default service account
		client, err = texttospeech.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create texttospeech client: %v", err)
		}
	}
	return &GoogleTTS{
		Client: client,
		Logger: log.New(os.Stdout, "GoogleTTS: ", log.LstdFlags),
	}
}

func (g *GoogleTTS) GenerateAudioAdapter(text string) ([]byte, error) {
	request := GoogleTTSRequest{
		Text: text,
	}
	response := g.GenerateAudio(request)
	if response.Error != "" {
		return nil, fmt.Errorf("%s", response.Error)
	}
	return response.AudioContent, nil
}

func (g *GoogleTTS) GenerateAudio(request GoogleTTSRequest) GoogleTTSResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var input *texttospeechpb.SynthesisInput
	if request.Text != "" {
		input = &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: request.Text,
			},
		}
	} else if request.SSML != "" {
		input = &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Ssml{
				Ssml: request.SSML,
			},
		}
	} else {
		return GoogleTTSResponse{
			Error: "text or ssml is required",
		}
	}

	if request.LanguageCode == "" {
		request.LanguageCode = "en-US"
	}
	if request.LanguageName == "" {
		request.LanguageName = "en-US-Chirp3-HD-Achernar"
	}

	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: input,
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: request.LanguageCode,
			Name:         request.LanguageName,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
			// EnableTimePointing: true,
		},
	}

	response, err := g.Client.SynthesizeSpeech(ctx, req)
	if err != nil {
		g.Logger.Printf("failed to synthesize speech: %v", err)
		return GoogleTTSResponse{
			Error: "failed to synthesize speech",
		}
	}

	return GoogleTTSResponse{
		AudioContent: response.AudioContent,
		AudioFormat:  "mp3",
		Error:        "",
	}
}

func (g *GoogleTTS) ListVoices(languageCode string) ([]*texttospeechpb.Voice, error) {
	if languageCode == "" {
		g.Logger.Println("Warning: languageCode is required but not provided")
		languageCode = "en-US"
	}
	ctx := context.Background()
	voices, err := g.Client.ListVoices(ctx, &texttospeechpb.ListVoicesRequest{})
	if err != nil {
		g.Logger.Printf("failed to list voices: %v", err)
		return nil, fmt.Errorf("failed to list voices: %v", err)
	}
	return voices.Voices, nil
}
