package helpers

import (
	"fmt"
	"log"
	"strings"

	"rio-go-model/configs"
)

// DynamicPrompting represents a service for generating dynamic story prompts
type DynamicPrompting struct {
	logger *log.Logger
}

// NewDynamicPrompting creates a new DynamicPrompting instance
func NewDynamicPrompting() *DynamicPrompting {
	return &DynamicPrompting{
		logger: log.New(log.Writer(), "[story.views] ", log.LstdFlags),
	}
}

// GetPlanetProtectorsStories generates planet protector story prompts
func (d *DynamicPrompting) GetPlanetProtectorsStories(country, city string, preferences []string) (string, error) {
	d.logger.Printf("Generating planet protector stories for country: %s, city: %s", country, city)
	
	// Convert preferences slice to string
	preferencesStr := strings.Join(preferences, ", ")
	
	superPrompt := fmt.Sprintf(
		"Generate a minimum of %d topics based on PLANET, ENVIRONMENT, ANIMALS, PLACES, PEOPLE, and other things that are related to the theme %s. "+
			"The topics should be easy, catchy and interesting in a way that kids can understand."+
			"The topics should be illustrate a story that kids can understand."+
			"The topics should also take them to different world and to illustrate what is environment/nature. %s."+
			"Each topic should be creative, entertainment-driven, engaging, fantasy-based, and align with the provided preferences: %s. "+
			"Each topic should be have exactly two parts title and description."+
			"title should be a short and catchy title that kids can understand."+
			"description should also be short and concise."+
			"seperate title and description with a colon. and maintain only one colon in the whole string."+
			"Return the topics as a list of strings it should be in title:description format.",
		GetDefaultStoryCount(),
		GetStoryTheme("1"),
		preferencesStr,
	)
	
	// d.logger.Printf("Generated prompt: %s", superPrompt)
	return superPrompt, nil
}

// GetMindfulStories generates mindful story prompts
func (d *DynamicPrompting) GetMindfulStories(country, religion string, preferences []string) (string, error) {
	d.logger.Printf("Generating mindful stories for country: %s, religion: %s", country, religion)
	
	superPrompt := fmt.Sprintf(
		"Create %d story topics that TEACH %s VALUES through SIMPLE STORIES. "+
			"Extract the topics from real %s scriptures/books/history and turn easy topic that kids can understand. "+
			"Each topic should be have exactly two parts title and description."+
			"title should be a short and catchy title that kids can understand."+
			"description should also be short and concise."+
			"seperate title and description with a colon. and maintain only one colon in the whole string."+
			"Return the topics as a list of strings it should be in title:description format.",
		GetDefaultStoryCount(),
		religion,
		religion,
	)
	
	// d.logger.Printf("Generated prompt: %s", superPrompt)
	return superPrompt, nil
}

// GetChillStories generates chill story prompts
func (d *DynamicPrompting) GetChillStories(preferences []string) (string, error) {
	d.logger.Printf("Generating chill stories for preferences: %v", preferences)
	
	// Convert preferences slice to string
	preferencesStr := strings.Join(preferences, ", ")
	
	superPrompt := fmt.Sprintf(
		"Create %d story topics that TEACH Simple/Slow Living VALUES through SIMPLE STORIES. "+
			"The topics should illustrate in the way of preferences: %s. "+
			"Each topic should be have exactly two parts title and description."+
			"title should be a short and catchy title that kids can understand."+
			"description should also be short and concise."+
			"seperate title and description with a colon. and maintain only one colon in the whole string."+
			"Return the topics as a list of strings it should be in title:description format.",
		GetDefaultStoryCount(),
		preferencesStr,
	)
	
	// d.logger.Printf("Generated prompt: %s", superPrompt)
	return superPrompt, nil
}

// Helper functions to get configuration values
// These would typically come from your config package

// GetDefaultStoryCount returns the default number of stories to generate
func GetDefaultStoryCount() int {
	settings := configs.LoadSettings()
	return settings.DefaultStoryToGenerate
}

// GetStoryTheme returns the story theme by key
func GetStoryTheme(key string) string {
	settings := configs.LoadSettings()
	
	if theme, exists := settings.GetStoryTheme(key); exists {
		return theme
	}
	
	return "DEFAULT_THEME"
}

// GetSettings returns configuration settings
func GetSettings() map[string]interface{} {
	settings := configs.LoadSettings()
	
	return map[string]interface{}{
		"DEFAULT_STORY_TO_GENERATE": settings.DefaultStoryToGenerate,
		"STORIES_PER_THEME":        settings.StoriesPerTheme,
		"STORY_THEMES":             settings.StoryThemes,
		"MAX_WORKERS":              settings.MaxWorkers,
		"MAX_STORY_WORKERS":        settings.MaxStoryWorkers,
		"HUGGINGFACE_TIMEOUT":      settings.HuggingFaceTimeout.Seconds(),
		"TTS_TIMEOUT":             settings.TTSTimeout.Seconds(),
	}
}
