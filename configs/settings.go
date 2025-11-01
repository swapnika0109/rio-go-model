package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Global settings instance
var GlobalSettings *Settings

// Settings represents the application settings
type Settings struct {
	// API Keys and Authentication
	HuggingFaceToken string
	SecretKey        string
	EmailHost        string
	EmailPort        int
	EmailAppPassword string
	EmailSender      string
	EmailTo          string

	// Story Generation Settings
	DefaultStoryToGenerate int
	StoriesPerTheme        int

	// File Upload Settings
	DataUploadMaxMemorySize int // in bytes
	FileUploadMaxMemorySize int // in bytes

	// Logging Configuration
	LogLevel       string
	LogFormat      string
	LogFile        string
	LogMaxSize     int
	LogBackupCount int
	LogEncoding    string

	// Performance and Concurrency Settings
	MaxWorkers       int
	MaxStoryWorkers  int
	MaxUploadWorkers int

	// Connection Pooling Settings
	MaxKeepaliveConnections  int
	MaxTotalConnections      int
	MaxRequestsPerConnection int
	ConnectionTimeout        time.Duration
	ReadTimeout              time.Duration
	KeepaliveTimeout         time.Duration

	// API-specific Timeouts
	HuggingFaceTimeout time.Duration
	TogetherAITimeout  time.Duration
	TTSTimeout         time.Duration

	// Story Configuration
	StoryConfig StoryConfig

	// Story Themes
	StoryThemes    map[string]string
	ChirpVoices    []string
	StandardVoices []string
	WaveNetVoices  []string
	DefaultVoice   string
}

// StoryConfig represents story-specific configuration
type StoryConfig struct {
	OllamaAPIURL string
	DefaultModel string
	Prompts      map[string]string
	MaxTokens    int
}

// PromptConfig represents a prompt configuration with system and prompt messages
type PromptConfig struct {
	System string `json:"system"`
	Prompt string `json:"prompt"`
}

// NewSettings creates a new Settings instance with default values
func NewSettings() *Settings {

	return &Settings{
		// Default values
		HuggingFaceToken:        getEnvString("HUGGINGFACE_TOKEN", ""),
		SecretKey:               getEnvString("SECRET_KEY", "********"),
		EmailHost:               getEnvString("EMAIL_HOST", "smtp.gmail.com"),
		EmailPort:               getEnvInt("EMAIL_PORT", 587),
		EmailAppPassword:        getEnvString("EMAIL_APP_PASSWORD", ""),
		EmailSender:             getEnvString("EMAIL_SENDER", "rio.oly.pluto@gmail.com"),
		EmailTo:                 getEnvString("EMAIL_TO", "rio.oly.pluto@gmail.com"),
		DefaultVoice:            getEnvString("DEFAULT_CHIRP_VOICE", "Standard"),
		DefaultStoryToGenerate:  getEnvInt("DEFAULT_STORY_TO_GENERATE", 0),
		StoriesPerTheme:         getEnvInt("STORIES_PER_THEME", 0),
		DataUploadMaxMemorySize: getEnvInt("DATA_UPLOAD_MAX_MEMORY_SIZE", 5242880), // 5MB
		FileUploadMaxMemorySize: getEnvInt("FILE_UPLOAD_MAX_MEMORY_SIZE", 5242880), // 5MB

		// Logging
		LogLevel:       getEnvString("LOG_LEVEL", "INFO"),
		LogFormat:      getEnvString("LOG_FORMAT", "2006-01-02 15:04:05 - %s - %s - %s"),
		LogFile:        getEnvString("LOG_FILE", "app.log"),
		LogMaxSize:     getEnvInt("LOG_MAX_SIZE", 10*1024*1024), // 10MB
		LogBackupCount: getEnvInt("LOG_BACKUP_COUNT", 5),
		LogEncoding:    getEnvString("LOG_ENCODING", "utf-8"),

		// Performance
		MaxWorkers:       getEnvInt("MAX_WORKERS", 4),
		MaxStoryWorkers:  getEnvInt("MAX_STORY_WORKERS", 2),
		MaxUploadWorkers: getEnvInt("MAX_UPLOAD_WORKERS", 2),

		// Connection Pooling
		MaxKeepaliveConnections:  getEnvInt("MAX_KEEPALIVE_CONNECTIONS", 20),
		MaxTotalConnections:      getEnvInt("MAX_TOTAL_CONNECTIONS", 100),
		MaxRequestsPerConnection: getEnvInt("MAX_REQUESTS_PER_CONNECTION", 1000),
		ConnectionTimeout:        time.Duration(getEnvInt("CONNECTION_TIMEOUT", 10)) * time.Second,
		ReadTimeout:              time.Duration(getEnvInt("READ_TIMEOUT", 60)) * time.Second,
		KeepaliveTimeout:         time.Duration(getEnvInt("KEEPALIVE_TIMEOUT", 30)) * time.Second,

		// API Timeouts
		HuggingFaceTimeout: time.Duration(getEnvInt("HUGGINGFACE_TIMEOUT", 120)) * time.Second,
		TogetherAITimeout:  time.Duration(getEnvInt("TOGETHER_AI_TIMEOUT", 180)) * time.Second,
		TTSTimeout:         time.Duration(getEnvInt("TTS_TIMEOUT", 90)) * time.Second,

		// Initialize complex configurations
		StoryConfig:    initStoryConfig(),
		StoryThemes:    initStoryThemes(),
		ChirpVoices:    initChirpVoices(),
		StandardVoices: initStandardVoices(),
		WaveNetVoices:  initWaveNetVoices(),
	}
}

// LoadSettings loads settings from environment variables
func LoadSettings() *Settings {
	settings := NewSettings()

	// Log initialization
	log.Println("Settings loaded")
	// log.Printf("Settings loaded - HuggingFace Token: %s, Log Level: %s",
	// 	maskToken(settings.HuggingFaceToken), settings.LogLevel)

	return settings
}

// Helper functions

// getEnvString gets environment variable with fallback (renamed to avoid conflict)
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func maskToken(token string) string {
	if len(token) > 10 {
		return token[:6] + "..." + token[len(token)-4:]
	}
	return "***"
}

// Configuration initializers

func initStoryConfig() StoryConfig {
	return StoryConfig{
		OllamaAPIURL: getEnvString("OLLAMA_API_URL", "http://localhost:11434/api/generate"),
		DefaultModel: getEnvString("DEFAULT_MODEL", "llama2"),
		Prompts: map[string]string{
			"story":     "Write a creative story about {topic}",
			"poem":      "Write a poem about {topic}",
			"adventure": "Write an adventure story about {topic} with a surprising twist",
		},
		MaxTokens: getEnvInt("MAX_TOKENS", 500),
	}
}

func initStoryThemes() map[string]string {
	return map[string]string{
		"1": "Eco-conscious, Clean planet, Nature, Wildlife, Environment",
		"2": "Mindful, Meditation, Yoga, Wellness, Self-care, tradition",
		"3": "Slow living, Minimalism, Minimalist lifestyle, Minimalist design, simple living",
	}
}

func initChirpVoices() []string {
	return []string{
		"-Chirp3-HD-Gacrux",
		"-Chirp3-HD-Achernar",
		"-Chirp3-HD-Callirrhoe",
		"-Chirp3-HD-Despina",
		"-Chirp3-HD-Iapetus",
		"-Chirp3-HD-Leda",
		"-Chirp3-HD-Zephyr",
		"-Chirp3-HD-Schedar",
		"-Chirp3-HD-Sadaltager",
		"-Chirp3-HD-Rasalgethi",
		"-Chirp3-HD-Umbriel",
		"-Chirp3-HD-Pulcherrima",
		"-Chirp3-HD-Charon",
		"-Chirp3-HD-Zubenelgenubi",
	}
}

func initStandardVoices() []string {
	return []string{
		"-Standard-A",
		"-Standard-B",
		"-Standard-C",
		"-Standard-D",
	}
}

func initWaveNetVoices() []string {
	return []string{
		"-WaveNet-A",
		"-WaveNet-B",
		"-WaveNet-C",
		"-WaveNet-D",
		"-WaveNet-E",
		"-WaveNet-F",
		"-WaveNet-G",
		"-WaveNet-H",
		"-WaveNet-I",
		"-WaveNet-J",
	}
}

// GetStoryTheme returns the story theme for a given theme key
func (s *Settings) GetStoryTheme(theme string) (string, bool) {
	themeText, exists := s.StoryThemes[theme]
	return themeText, exists
}

// FormatPrompt formats a prompt template with the given parameters
func (s *Settings) FormatPrompt(template string, params map[string]string) string {
	result := template
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		result = fmt.Sprintf(strings.ReplaceAll(result, placeholder, value))
	}
	return result
}

// Validate validates the settings configuration
func (s *Settings) Validate() error {
	if s.HuggingFaceToken == "" {
		return fmt.Errorf("HUGGINGFACE_TOKEN is required")
	}

	if s.DefaultStoryToGenerate <= 0 {
		return fmt.Errorf("DEFAULT_STORY_TO_GENERATE must be positive")
	}

	if s.MaxWorkers <= 0 {
		return fmt.Errorf("MAX_WORKERS must be positive")
	}

	return nil
}

// InitializeSettings loads settings once at startup
func InitializeSettings() {
	GlobalSettings = LoadSettings()
	log.Println("✅ Global settings loaded successfully")
}

// GetSettings returns the global settings instance
func GetSettings() *Settings {
	if GlobalSettings == nil {
		log.Println("⚠️  Settings not initialized, loading now...")
		GlobalSettings = LoadSettings()
	}
	return GlobalSettings
}
