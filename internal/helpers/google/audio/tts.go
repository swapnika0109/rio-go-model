package audio

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"rio-go-model/internal/util"

	"rio-go-model/internal/util/tokens"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"google.golang.org/api/option"
)

type GoogleTTS struct {
	Client          *texttospeech.Client
	Logger          *log.Logger
	storyCharacters *tokens.StoryCharacters
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
	log.Println("Initializing Google TTS client...")
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
	log.Println("Google TTS client initialized successfully")
	storyCharacters := tokens.NewStoryCharacters()
	return &GoogleTTS{
		Client:          client,
		storyCharacters: storyCharacters,
		Logger:          log.New(os.Stdout, "GoogleTTS: ", log.LstdFlags),
	}
}

func (g *GoogleTTS) GenerateAudioAdapter(text string, language string, theme string, voice string) ([]byte, int32, error) {
	g.Logger.Printf("GenerateAudioAdapter called - Language: %s, Text length: %d, Theme: %s", language, len(text), theme)
	var ssml string
	var totalTokens int32
	languageCode := util.LanguageMapper(language, theme)
	voiceList := util.GetVoiceList(voice)
	no, err := util.RandomFromLength(len(voiceList))
	if err != nil {
		g.Logger.Printf("Failed to get random voice number: %v", err)
		return nil, totalTokens, fmt.Errorf("failed to get random voice number: %v", err)
	}
	languageName := util.GetVoice(languageCode, no, voiceList)
	log.Println("Language name: %s", languageName)
	// voices, err := g.ListVoices(languageCode)
	// if err != nil {
	// 	g.Logger.Printf("Failed to list voices: %v", err)
	// 	return nil, totalTokens, fmt.Errorf("failed to list voices: %v", err)
	// }
	// g.Logger.Printf("Voices: %v", voices)
	g.Logger.Printf("Mapped language code: %s, Voice name: %s", languageCode, languageName)

	// if language == "Telugu" && voice == tts.Standard.String() {
	// 	g.Logger.Printf("Processing Telugu text with SSML...")
	// 	teluguSSMLBuilder := util.NewTeluguSSMLBuilder()
	// 	ssml = teluguSSMLBuilder.BuildTeluguSSML(text)
	// 	g.Logger.Printf("Generated Telugu SSML length: %d bytes", len(ssml))
	// 	g.Logger.Printf("Updated voice name for Telugu: %s", languageName)
	// 	// Check if SSML exceeds 5000 byte limit
	// 	if len(ssml) > 5000 {
	// 		g.Logger.Printf("SSML exceeds 5000 byte limit (%d bytes), splitting into chunks...", len(ssml))
	// 		audio, err := g.generateAudioInChunks(text, language, languageCode, languageName)
	// 		return audio, totalTokens, err
	// 	}
	// }

	// Check if normal text exceeds 5000 byte limit (for non-SSML languages)
	if len(ssml) == 0 && len(text) > 5000 {
		g.Logger.Printf("Text exceeds 5000 byte limit (%d bytes), splitting into chunks...", len(text))
		audio, err := g.generateAudioInChunksNormal(text, language, languageCode, languageName)
		return audio, totalTokens, err
	}

	request := GoogleTTSRequest{
		Text:         text,
		SSML:         ssml,
		LanguageCode: languageCode,
		LanguageName: languageName,
	}
	// totalTokens is set from input char counts below
	// Log audio character counts (approx billing units)
	if g.storyCharacters != nil {
		totalTokens = int32(g.storyCharacters.CountAudioChars(text, ssml))
	}

	g.Logger.Printf("Calling GenerateAudio with request...")
	response := g.GenerateAudio(request)
	if response.Error != "" {
		g.Logger.Printf("GenerateAudio returned error: %s", response.Error)
		return nil, totalTokens, fmt.Errorf("%s", response.Error)
	}
	g.Logger.Printf("GenerateAudio succeeded, audio content length: %d", len(response.AudioContent))
	return response.AudioContent, totalTokens, nil
}

// generateAudioInChunksNormal splits long text into smaller chunks without SSML and combines the audio
func (g *GoogleTTS) generateAudioInChunksNormal(text, language, languageCode, languageName string) ([]byte, error) {
	g.Logger.Printf("Splitting text into chunks for processing (normal text)...")

	// Split text into sentences for better chunking
	sentences := g.splitIntoSentences(text)
	var audioChunks [][]byte

	g.Logger.Printf("Original text length: %d characters, split into %d sentences", len(text), len(sentences))

	// Process each sentence as a separate chunk
	// Fail-fast: If ANY chunk fails, stop immediately and return error (story will be bypassed)
	for i, sentence := range sentences {
		if len(sentence) == 0 {
			continue
		}

		g.Logger.Printf("Processing chunk %d/%d: %d characters", i+1, len(sentences), len(sentence))

		// If chunk is still too long, split further
		if len(sentence) > 5000 {
			g.Logger.Printf("Chunk %d still too long (%d bytes), splitting further...", i+1, len(sentence))
			subChunks := g.splitTextIntoSmallerChunks(sentence, 4000)
			for j, subChunk := range subChunks {
				audioData, err := g.generateSingleChunk("", subChunk, languageCode, languageName)
				if err != nil {
					g.Logger.Printf("❌ Failed to generate audio for sub-chunk %d of chunk %d: %v", j+1, i+1, err)
					g.Logger.Printf("⏹️  Stopping audio generation - story will be bypassed")
					return nil, fmt.Errorf("audio generation failed at chunk %d, sub-chunk %d: %v", i+1, j+1, err)
				}
				audioChunks = append(audioChunks, audioData)
			}
		} else {
			// Generate audio for this chunk
			audioData, err := g.generateSingleChunk("", sentence, languageCode, languageName)
			if err != nil {
				g.Logger.Printf("❌ Failed to generate audio for chunk %d: %v", i+1, err)
				g.Logger.Printf("⏹️  Stopping audio generation - story will be bypassed")
				return nil, fmt.Errorf("audio generation failed at chunk %d: %v", i+1, err)
			}
			audioChunks = append(audioChunks, audioData)
		}
	}

	// Validation: Check if all chunks were processed
	if len(audioChunks) == 0 {
		return nil, fmt.Errorf("no audio chunks generated")
	}

	g.Logger.Printf("✅ All %d chunks processed successfully", len(audioChunks))

	// Combine all audio chunks
	g.Logger.Printf("Combining %d audio chunks...", len(audioChunks))
	combinedAudio := g.combineAudioChunks(audioChunks)

	// Final validation: Check combined audio size is reasonable
	if len(combinedAudio) == 0 {
		return nil, fmt.Errorf("combined audio is empty")
	}
	g.Logger.Printf("✅ Audio generation complete: %d bytes", len(combinedAudio))

	return combinedAudio, nil
}

// generateAudioInChunks splits long text into smaller chunks and combines the audio
func (g *GoogleTTS) generateAudioInChunks(text, language, languageCode, languageName string) ([]byte, error) {
	g.Logger.Printf("Splitting text into chunks for processing...")

	// Split text into sentences for better chunking
	sentences := g.splitIntoSentences(text)
	var audioChunks [][]byte

	// Process each sentence as a separate chunk
	// Fail-fast: If ANY chunk fails, stop immediately and return error (story will be bypassed)
	for i, sentence := range sentences {
		if len(sentence) == 0 {
			continue
		}

		g.Logger.Printf("Processing chunk %d/%d: %d characters", i+1, len(sentences), len(sentence))

		// Generate SSML for this chunk
		teluguSSMLBuilder := util.NewTeluguSSMLBuilder()
		chunkSSML := teluguSSMLBuilder.BuildTeluguSSML(sentence)

		// If chunk SSML is still too long, split further
		if len(chunkSSML) > 5000 {
			g.Logger.Printf("Chunk %d SSML still too long (%d bytes), splitting further...", i+1, len(chunkSSML))
			subChunks := g.splitTextIntoSmallerChunks(sentence, 4000) // Leave room for SSML markup
			for j, subChunk := range subChunks {
				subSSML := teluguSSMLBuilder.BuildTeluguSSML(subChunk)
				audioData, err := g.generateSingleChunk(subSSML, "", languageCode, languageName)
				if err != nil {
					g.Logger.Printf("❌ Failed to generate audio for sub-chunk %d of chunk %d: %v", j+1, i+1, err)
					g.Logger.Printf("⏹️  Stopping audio generation - story will be bypassed")
					return nil, fmt.Errorf("audio generation failed at chunk %d, sub-chunk %d: %v", i+1, j+1, err)
				}
				audioChunks = append(audioChunks, audioData)
			}
		} else {
			audioData, err := g.generateSingleChunk(chunkSSML, "", languageCode, languageName)
			if err != nil {
				g.Logger.Printf("❌ Failed to generate audio for chunk %d: %v", i+1, err)
				g.Logger.Printf("⏹️  Stopping audio generation - story will be bypassed")
				return nil, fmt.Errorf("audio generation failed at chunk %d: %v", i+1, err)
			}
			audioChunks = append(audioChunks, audioData)
		}
	}

	if len(audioChunks) == 0 {
		return nil, fmt.Errorf("no audio chunks were generated successfully")
	}

	g.Logger.Printf("Successfully generated %d audio chunks, combining...", len(audioChunks))
	return g.combineAudioChunks(audioChunks), nil
}

// generateSingleChunk generates audio for a single SSML chunk
func (g *GoogleTTS) generateSingleChunk(ssml, text, languageCode, languageName string) ([]byte, error) {
	request := GoogleTTSRequest{
		SSML:         ssml,
		Text:         text,
		LanguageCode: languageCode,
		LanguageName: languageName,
	}
	response := g.GenerateAudio(request)
	if response.Error != "" {
		return nil, fmt.Errorf("failed to generate audio: %s", response.Error)
	}
	return response.AudioContent, nil
}

// splitIntoSentences splits text into sentences
func (g *GoogleTTS) splitIntoSentences(text string) []string {
	// Simple sentence splitting by periods, exclamation marks, and question marks
	// This is a basic implementation - you might want to use a more sophisticated NLP library
	re := regexp.MustCompile(`[.!?]+`)
	sentences := re.Split(text, -1)
	var result []string
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) > 0 {
			result = append(result, sentence)
		}
	}
	return result
}

// splitTextIntoSmallerChunks splits text into smaller chunks by character count
func (g *GoogleTTS) splitTextIntoSmallerChunks(text string, maxChars int) []string {
	var chunks []string
	for i := 0; i < len(text); i += maxChars {
		end := i + maxChars
		if end > len(text) {
			end = len(text)
		}
		chunk := text[i:end]
		chunks = append(chunks, chunk)
	}
	return chunks
}

// combineAudioChunks combines multiple audio chunks into one
func (g *GoogleTTS) combineAudioChunks(chunks [][]byte) []byte {
	// For MP3 files, we need to concatenate the raw audio data
	// This is a simple concatenation - for production, you might want to use a proper audio library
	var combined []byte
	for _, chunk := range chunks {
		combined = append(combined, chunk...)
	}
	return combined
}

func (g *GoogleTTS) GenerateAudio(request GoogleTTSRequest) GoogleTTSResponse {
	g.Logger.Printf("GenerateAudio called with SSML: %d, Text: %d, LanguageCode: %s", len(request.SSML), len(request.Text), request.LanguageCode)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	var input *texttospeechpb.SynthesisInput
	if request.SSML != "" {
		input = &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Ssml{
				Ssml: request.SSML,
			},
		}
	} else if request.Text != "" {
		input = &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: request.Text,
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

	g.Logger.Printf("Calling Google TTS API...")
	response, err := g.Client.SynthesizeSpeech(ctx, req)
	if err != nil {
		g.Logger.Printf("=== TTS Request Debug ===")
		g.Logger.Printf("Language Code: %s", request.LanguageCode)
		g.Logger.Printf("Language Name: %s", request.LanguageName)
		g.Logger.Printf("Text Length: %d", len(request.Text))
		g.Logger.Printf("SSML Length: %d", len(request.SSML))
		if len(request.SSML) > 0 {
			g.Logger.Printf("SSML Content: %s", request.SSML)
		}
		if len(request.Text) > 0 {
			g.Logger.Printf("Text Content: %s", request.Text)
		}
		g.Logger.Printf("========================")
		g.Logger.Printf("failed to synthesize speech: %v", err)
		return GoogleTTSResponse{
			Error: "failed to synthesize speech",
		}
	}
	g.Logger.Printf("Google TTS API call successful, audio content length: %d", len(response.AudioContent))

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
