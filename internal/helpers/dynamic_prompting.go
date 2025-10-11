package helpers

import (
	"log"

	// "strings"

	"rio-go-model/configs"
	"rio-go-model/configs/english"
	"rio-go-model/configs/telugu"
	"rio-go-model/internal/util"
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
func (d *DynamicPrompting) GetPlanetProtectorsStories(country, city string, preference string, language string, storiesPerPreference int) (string, error) {
	d.logger.Printf("Generating planet protector stories for country: %s, city: %s", country, city)

	planetProtectors := english.ThemesSettings()
	var promptText string

	for i := 0; i < storiesPerPreference; i++ {
		topicNumber, err := util.RandomFrom(planetProtectors.PlanetProtectorTopicsList)
		if err != nil {
			topicNumber = 0
		}
		topic := planetProtectors.PlanetProtectorTopicsList[topicNumber]
		if i < storiesPerPreference-1 {
			promptText += topic + ", "
		} else {
			promptText += topic
		}
	}

	var superPrompt string
	switch language {
	case "English":
		superPrompt = english.SuperPlanetProtectorPrompt(promptText, preference, storiesPerPreference)
	case "Telugu":
		superPrompt = telugu.SuperPlanetProtectorPrompt(promptText, preference, storiesPerPreference)
	}

	d.logger.Printf("Generated prompt: %s", superPrompt)
	return superPrompt, nil
}

// GetMindfulStories generates mindful story prompts
func (d *DynamicPrompting) GetMindfulStories(country, religion string, preferences []string, language string, storiesPerPreference int) (string, error) {
	d.logger.Printf("Generating mindful stories for country: %s, religion: %s", country, religion)

	mindfulStoriesSettings := english.ThemesSettings()
	var promptText string

	for i := 0; i < storiesPerPreference; i++ {
		topicNumber, err := util.RandomFrom(mindfulStoriesSettings.MindfulStoriesList[religion])
		if err != nil {
			topicNumber = 0
		}
		topic := mindfulStoriesSettings.MindfulStoriesList[religion][topicNumber]
		if i < storiesPerPreference-1 {
			promptText += topic + ", "
		} else {
			promptText += topic
		}
	}

	var superPrompt string

	switch language {
	case "English":
		superPrompt = english.SuperMindfulStoriesPrompt(promptText, religion, storiesPerPreference)
	case "Telugu":
		superPrompt = telugu.SuperMindfulStoriesPrompt(promptText, religion, storiesPerPreference)
	}
	return superPrompt, nil
}

// GetChillStories generates chill story prompts
func (d *DynamicPrompting) GetChillStories(preference string, language string, storiesPerPreference int) (string, error) {
	d.logger.Printf("Generating chill stories for preferences: %v", preference)
	d.logger.Printf("storiesPerPreference: %d", storiesPerPreference)

	chillStoriesSettings := english.ThemesSettings()
	var promptText string

	for i := 0; i < storiesPerPreference; i++ {
		topicNumber, err := util.RandomFrom(chillStoriesSettings.ChillStoriesList)
		if err != nil {
			topicNumber = 0
		}
		topic := chillStoriesSettings.ChillStoriesList[topicNumber]
		if i < storiesPerPreference-1 {
			promptText += topic + ", "
		} else {
			promptText += topic
		}
	}
	var superPrompt string
	switch language {
	case "English":
		superPrompt = english.SuperChillStoriesPrompt(promptText, preference, storiesPerPreference)
	case "Telugu":
		superPrompt = telugu.SuperChillStoriesPrompt(promptText, preference, storiesPerPreference)
	}
	return superPrompt, nil
}

// Helper functions to get configuration values
// These would typically come from your config package

// GetDefaultStoryCount returns the default number of stories to generate
func GetDefaultStoryCount() int {
	settings := configs.GetSettings()
	return settings.DefaultStoryToGenerate
}

// GetStoryTheme returns the story theme by key
func GetStoryTheme(key string) string {
	settings := configs.GetSettings()

	if theme, exists := settings.GetStoryTheme(key); exists {
		return theme
	}

	return "DEFAULT_THEME"
}

// GetSettings returns configuration settings
func GetSettings() map[string]interface{} {
	settings := configs.GetSettings()

	return map[string]interface{}{
		"DEFAULT_STORY_TO_GENERATE": settings.DefaultStoryToGenerate,
		"STORIES_PER_THEME":         settings.StoriesPerTheme,
		"STORY_THEMES":              settings.StoryThemes,
		"MAX_WORKERS":               settings.MaxWorkers,
		"MAX_STORY_WORKERS":         settings.MaxStoryWorkers,
		"HUGGINGFACE_TIMEOUT":       settings.HuggingFaceTimeout.Seconds(),
		"TTS_TIMEOUT":               settings.TTSTimeout.Seconds(),
	}
}
