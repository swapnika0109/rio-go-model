package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"rio-go-model/internal/helpers/google/audio"
	"rio-go-model/internal/util"
)

func main() {
	// Telugu test text
	teluguText := `ఒక ఉదయం, చిట్టి అనే చిన్నారి, తన ఫ్రెండ్స్ తో కలిసి, ఒక పెద్ద అడవిలోకి వెళ్ళింది! Wow, ఎంత అందంగా ఉందో ఆ అడవి! రకరకాల చెట్లు, రంగురంగుల పువ్వులు, ఎగురుతున్న పక్షులు... చూడడానికి చాలా బాగున్నాయి!`
	// debugVothulu()
	fmt.Println("🎤 Telugu TTS Test Starting...")
	fmt.Printf("📝 Input text length: %d characters\n", len(teluguText))
	fmt.Println("📝 Input text:")
	fmt.Println(teluguText)
	fmt.Println(strings.Repeat("=", 80))

	// Step 1: Convert Telugu text to SSML
	fmt.Println("🔄 Step 1: Converting Telugu text to SSML...")
	teluguSSMLBuilder := util.NewTeluguSSMLBuilder()
	ssml := teluguSSMLBuilder.BuildTeluguSSML(teluguText)

	fmt.Printf("✅ SSML generated successfully!\n")
	fmt.Printf("📊 SSML length: %d bytes\n", len(ssml))
	fmt.Println("📄 Generated SSML:")
	fmt.Println(ssml)
	fmt.Println(strings.Repeat("=", 80))

	// Step 2: Initialize Google TTS
	fmt.Println("🔄 Step 2: Initializing Google TTS...")
	tts := audio.NewGoogleTTS()
	fmt.Println("✅ Google TTS initialized successfully!")

	// Step 3: Generate audio
	fmt.Println("🔄 Step 3: Generating audio from SSML...")
	startTime := time.Now()

	audioData, err := tts.GenerateAudioAdapter(teluguText, "Telugu")
	if err != nil {
		log.Fatalf("❌ Failed to generate audio: %v", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("✅ Audio generated successfully in %v!\n", duration)
	fmt.Printf("📊 Audio data length: %d bytes\n", len(audioData))

	// Step 4: Save audio file
	fmt.Println("🔄 Step 4: Saving audio file...")

	// Create output directory if it doesn't exist
	outputDir := "test_audio_output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("❌ Failed to create output directory: %v", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("telugu_test_%s.mp3", timestamp)
	filepath := filepath.Join(outputDir, filename)

	// Write audio data to file
	if err := ioutil.WriteFile(filepath, audioData, 0644); err != nil {
		log.Fatalf("❌ Failed to save audio file: %v", err)
	}

	fmt.Printf("✅ Audio file saved successfully!\n")
	fmt.Printf("📁 File path: %s\n", filepath)
	fmt.Printf("📊 File size: %d bytes\n", len(audioData))

	// Step 5: Display summary
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎉 Test completed successfully!")
	fmt.Printf("📝 Input text: %d characters\n", len(teluguText))
	fmt.Printf("📄 SSML: %d bytes\n", len(ssml))
	fmt.Printf("🎵 Audio: %d bytes\n", len(audioData))
	fmt.Printf("⏱️  Generation time: %v\n", duration)
	fmt.Printf("💾 Output file: %s\n", filepath)
	fmt.Println(strings.Repeat("=", 80))
}
