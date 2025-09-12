package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Settings represents the application settings
type Settings struct {
	// API Keys and Authentication
	HuggingFaceToken string
	SecretKey        string

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
	MaxWorkers        int
	MaxStoryWorkers   int
	MaxUploadWorkers  int

	// Connection Pooling Settings
	MaxKeepaliveConnections   int
	MaxTotalConnections       int
	MaxRequestsPerConnection  int
	ConnectionTimeout         time.Duration
	ReadTimeout              time.Duration
	KeepaliveTimeout         time.Duration

	// API-specific Timeouts
	HuggingFaceTimeout time.Duration
	TogetherAITimeout  time.Duration
	TTSTimeout         time.Duration

	// Story Configuration
	StoryConfig StoryConfig

	// Story Themes
	StoryThemes map[string]string

	// Prompts Configuration
	PromptsConfig map[string]PromptConfig

	// Preferences
	Preferences map[string]string

	// Dynamic Prompts Configuration
	DynamicPromptsConfig map[string]PromptConfig
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
		HuggingFaceToken:       getEnvString("HUGGINGFACE_TOKEN", ""),
		SecretKey:             getEnvString("SECRET_KEY", "********"),
		DefaultStoryToGenerate: getEnvInt("DEFAULT_STORY_TO_GENERATE", 10),
		StoriesPerTheme:       getEnvInt("STORIES_PER_THEME", 10),
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
		ReadTimeout:             time.Duration(getEnvInt("READ_TIMEOUT", 60)) * time.Second,
		KeepaliveTimeout:        time.Duration(getEnvInt("KEEPALIVE_TIMEOUT", 30)) * time.Second,

		// API Timeouts
		HuggingFaceTimeout: time.Duration(getEnvInt("HUGGINGFACE_TIMEOUT", 120)) * time.Second,
		TogetherAITimeout:  time.Duration(getEnvInt("TOGETHER_AI_TIMEOUT", 180)) * time.Second,
		TTSTimeout:         time.Duration(getEnvInt("TTS_TIMEOUT", 90)) * time.Second,

		// Initialize complex configurations
		StoryConfig:          initStoryConfig(),
		StoryThemes:          initStoryThemes(),
		PromptsConfig:        initPromptsConfig(),
		Preferences:          initPreferences(),
		DynamicPromptsConfig: initDynamicPromptsConfig(),
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

func initPreferences() map[string]string {
	return map[string]string{
		"FUN":       "The ENTIRE story must be funny. Characters MUST say funny things, do silly things, and create humorous situations throughout. Include jokes, wordplay, silly mistakes, and funny dialogue. Make kids laugh out loud! ",
		"EXCITED":   "The ENTIRE story must be exciting. Include high-energy moments, surprises, and thrilling discoveries that get kids excited. Include unexpected twists, exciting finds, and moments that make kids gasp with wonder.",
		"ADVENTURE": "The ENTIRE story must be adventurous. Take kids on a real journey with exciting discoveries, new places, challenges to overcome, and thrilling moments. Include obstacles, new locations, and exciting discoveries along the way. ",
		"KINDNESS":  "The ENTIRE story must focus on kindness. Show characters helping each other, sharing resources, and being kind in specific situations throughout the story.",
		"HAPPY":     "The ENTIRE story must be joyful. Include celebrations, achievements, and moments of pure joy throughout. Make kids feel good!",
		"CHILL":     "The ENTIRE story must be calm and peaceful. Include quiet moments, gentle activities, and peaceful scenes throughout. ",
	}
}

func initPromptsConfig() map[string]PromptConfig {
	return map[string]PromptConfig{
		"1": {
			System: "You are a creative storyteller who weaves magical tales that inspire children to think innovatively about real-world challenges",
			Prompt: `Create an enchanting story about {topic} that inspires children to think innovatively about environmental challenges. 
The story should spark imagination and creative problem-solving through: 
- Use ONLY simple, everyday words that 3-year-olds understand (like 'trees', 'flowers', 'animals', 'friends', 'helpers', not 'rainforest', 'ecosystem', 'warriors', 'enchantment') 
- Use made-up, fun place names instead of real cities (like 'Sunny Valley', 'Green Meadows', 'Happy Town') 
- Use simple character names like 'Tree Friends', 'Garden Helpers', 'Nature Friends' instead of complex terms 
- Use simple problems like 'dirty water', 'too much trash', 'sick trees' instead of complex terms 
IMPORTANT: Write this as an ACTUAL STORY with: 
- A real journey that takes kids from one place to another 
- Natural dialogue that sounds like real conversations 
- Specific scenes and moments, not summaries 
- Add interactivity & choices
- Expand educational elements
- Problem-solving, kindness moments
- Boost sensory description and inclusivity
- Explore more emotions & ending resolution
- Show, don't tell - let kids experience the story 
Write the story as a flowing narrative that takes kids on a journey, not as a book summary. Make every scene vivid and engaging. Use only words a 3-year-old would understand. NO complex terms!`,
		},
		"2": {
			System: "You are a wise storyteller who brings ancient wisdom to life through enchanting tales that children can understand and live by.",
			Prompt: `Create a magical story about {topic} that teaches the deep wisdom of ancient texts and spiritual teachings in a way children can understand and live by. 
The story should bring sacred wisdom to life through: 
- Simple, beautiful language that makes complex wisdom accessible to children 
- Characters who embody the teachings of their ancestors and gods 
- Stories from sacred texts and ancient wisdom, simplified for children 
- Gentle lessons about values, morals, and spiritual principles 
- Beautiful, peaceful settings that feel both ancient and magical 
- Natural dialogue that teaches through example and story 
- Heartwarming messages about kindness, courage, and inner strength 
- Simple breathing or mindfulness moments that connect to spiritual practices 
- Creative ways to apply ancient wisdom to modern life 
Write the story as a flowing narrative that takes kids on a journey, not as a book summary. Make every scene vivid and engaging. Ensure children can understand and implement the teachings in their daily lives.`,
		},
		"3": {
			System: "You are a joyful storyteller who celebrates simple living while inspiring children to help others creatively.",
			Prompt: `Create a delightful story about {topic} that celebrates simple living while inspiring children to help others. 
The story should encourage creative problem-solving through: 
- Cheerful, natural language that makes children smile 
- Lovable characters who find joy in simple pleasures 
- Cozy, warm settings that feel like home 
- Gentle humor and happy moments throughout 
- Clear, natural dialogue that children can easily follow 
- A heartwarming message about appreciating what we have 
- A satisfying ending that leaves children feeling content 
- Creative ways to solve problems with simple, practical solutions 
Write the story as a flowing narrative that takes kids on a journey, not as a book summary. Make every scene vivid and engaging. Make every word count.`,
		},
	}
}

func initDynamicPromptsConfig() map[string]PromptConfig {
	return map[string]PromptConfig{
		"1": {
			System: "You are a creative entertainment-driven and animated imaginative storyteller who weaves magical tales that inspire children to think innovatively about environmental themes. NEVER use complex terms like 'rainforest', 'ecosystem', 'warriors', or 'enchantment'. Write ONLY simple, engaging stories with natural dialogue.",
			Prompt: `Create a VERY ELABORATE and enchanting story about {topic} that can be easily understandable by people staying in {country} and {city}. 
CRITICAL REQUIREMENTS - FOLLOW THESE EXACTLY: 
- The story should include humour and should be understandable by kids of 3-5 years old
- Write this as a VERY ELABORATE STORY with rich details and have at least 8 to 10 illustrated/animated detailed scenes with vivid descriptions
- Add interactivity, challenges & choices to the story by having deep character development
- Make the story non-linear with more opportunities for kids to interact, pick solutions or answer questions
- Include character emotions, thoughts, and reactions throughout by having natural dialogue that sounds like real time conversations
- Use catchy names for kids to understand and imagine
- Whenever needed Add rich sensory details (sounds, smells, colors, textures, tastes) and support illustration/animation
- Add brief moments of character when they are having uncertain emotions
- Add more educational elements, STEM learning, nature learning, scientific learning etc.
- Whenever needed Add surprising twists and discoveries and illustration too
- Add clear descriptions of any new places or objects to make kids imagine like an animation movie
- Explore more emotions & ending resolution
- The story has to interact more deeply with characters/places but not with the user.
- Don't end the story abruptly, don't ask user to share ideas. and also don't repeat the story at the end.
IMPORTANT: Write ONLY the story. NO notes, NO explanations, NO meta-commentary. Just write the story as a flowing narrative that takes kids on a journey. Use only words a 3-year-old would understand. NO complex terms!`,
		},
		"2": {
			System: "You are a wise grandparent who brings ancient wisdom and history in the form of stories to the children in a way they can understand and live by.",
			Prompt: `Read the topic: {topic} and fill the real/existing story behind it.
Add more details, emotions, and interactions to the story.
Add more illustrations and animations and make the story more engaging and interactive and understandable for the kids of 3-9 years old
Illustrate the story with more educational elements, STEM learning, nature learning, scientific learning etc. these learnings should be part of story.
Kids should learn the story by understanding the science and moral in it.
Add more surprises and discoveries to the story.
Add more interactions with characters and places to the story.
- Whenever needed Add surprising twists and discoveries and illustration too
- The story has to interact more deeply with characters/places but not with the user.
- Don't mention about learnings in the end of the story. it should be part of story.
IMPORTANT: Write ONLY the story. NO notes, NO explanations, NO meta-commentary. Just write the story as a flowing narrative that takes kids on a journey. Ensure children can understand and implement the teachings in their daily lives.`,
		},
		"3": {
			System: "You are a creative, entertainment-driven and animated storyteller",
			Prompt: `Illustrate a story like disney animated movie about {topic}, by explaining how important it is to live the life slowly and observe what life is giving us with an example related to the story.
Also explain the science and moral behind the story by adding more details like multiple scenes with more interactions having beautiful emotions
Each scene should be very engaging and give surprising illustrations and animations
Use catchy and interesting names. For human characters please use easy or real human names for the kids. 
Add more surprises when needed.
Make the story more engaging and interactive and understandable for the kids of 3-9 years old
Kids should learn the story by understanding the science and moral in it.
Add more interactions with characters and places to the story.
- The story has to interact more deeply with characters/places but not with the user.
- Don't mention about learnings in the end of the story. it should be part of story.
- Don't end the story abruptly.
IMPORTANT: Write ONLY the story. NO notes, NO explanations, NO meta-commentary. Just write the story as a flowing narrative that takes kids on a journey. Ensure children can understand and implement the teachings in their daily lives.`,
		},
	}
}

// GetPromptConfig returns the prompt configuration for a given theme
func (s *Settings) GetPromptConfig(theme string) (PromptConfig, bool) {
	config, exists := s.PromptsConfig[theme]
	return config, exists
}

// GetDynamicPromptConfig returns the dynamic prompt configuration for a given theme
func (s *Settings) GetDynamicPromptConfig(theme string) (PromptConfig, bool) {
	config, exists := s.DynamicPromptsConfig[theme]
	return config, exists
}

// GetPreference returns the preference text for a given preference key
func (s *Settings) GetPreference(preference string) (string, bool) {
	pref, exists := s.Preferences[preference]
	return pref, exists
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
