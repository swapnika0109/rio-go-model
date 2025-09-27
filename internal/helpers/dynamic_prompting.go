package helpers

import (
	"fmt"
	"log"
	// "strings"

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
func (d *DynamicPrompting) GetPlanetProtectorsStories(country, city string, preference string, storiesPerPreference int) (string, error) {
	d.logger.Printf("Generating planet protector stories for country: %s, city: %s", country, city)
	
	planetProtectors := configs.ThemesSettings()
	var promptText string

	for i := 0; i < storiesPerPreference; i++ {
		topicNumber, err := configs.RandomFrom(planetProtectors.PlanetProtectorTopicsList)
		if err != nil {
			topicNumber = 0
		}
		topic := planetProtectors.PlanetProtectorTopicsList[topicNumber]
		if i < storiesPerPreference - 1 {
			promptText += topic + ", "
		} else {
			promptText += topic
		}	
	}

	
	superPrompt := fmt.Sprintf(
		
			"Generate one topic for each topic in the list "+ promptText + " and other things that are related to the theme %s. "+
			"The topics should be very easy, catchy and interesting in a way that toddlers can understand."+
			"The topics should be illustrate a story that kids can understand."+
			"The topics should also take them to different world and to illustrate the topic in a very creative way."+
			"Each topic should be creative, entertainment-driven, engaging, fantasy-based, and align with the provided preferences: %s. "+
			"Each topic should be have exactly two parts title and description."+
			"title should be a short and catchy title that kids can understand."+
			"description should also be short and concise."+
			"seperate title and description with a colon. and maintain only one colon in the whole string."+
			"Return the topics as a list of strings and it should be in title:description format."+
			"Always validate the length of the topics should be alwys %d.",
		GetStoryTheme("1"),	
		preference,
	)
	
	d.logger.Printf("Generated prompt: %s", superPrompt)
	return superPrompt, nil
}

// GetMindfulStories generates mindful story prompts
func (d *DynamicPrompting) GetMindfulStories(country, religion string, preferences []string, storiesPerPreference int) (string, error) {
	d.logger.Printf("Generating mindful stories for country: %s, religion: %s", country, religion)

	mindfulStoriesSettings := configs.ThemesSettings()
	var promptText string

	for i := 0; i < storiesPerPreference; i++ {
		topicNumber, err := configs.RandomFrom(mindfulStoriesSettings.MindfulStoriesList[religion])
		if err != nil {
			topicNumber = 0
		}
		topic := mindfulStoriesSettings.MindfulStoriesList[religion][topicNumber]
		if i < storiesPerPreference - 1 {
			promptText += topic + ", "
		} else {
			promptText += topic
		}	
	}
	
	superPrompt := fmt.Sprintf(
		"Create one topic for each topic in the list : "+ promptText + ". that TEACH %s VALUES through SIMPLE STORIES. "+
			"The topic should be based on a real/existing topic that kids can understand. "+
			"Each topic should be have exactly two parts title and description."+
			"title should be a short and catchy title that kids can understand."+
			"description should also be short and concise."+
			"seperate title and description with a colon. and maintain only one colon in the whole string."+
			"Return the topics as a list of strings and it should be in title:description format.",
			"Always validate the length of the topics should be alwys %d.",
			"Dont add any direct book or scripture name in the title or description.",
		religion,
		storiesPerPreference,
	)
	
	d.logger.Printf("Generated prompt: %s", superPrompt)
	return superPrompt, nil
}

// GetChillStories generates chill story prompts
func (d *DynamicPrompting) GetChillStories(preference string, storiesPerPreference int) (string, error) {
	d.logger.Printf("Generating chill stories for preferences: %v", preference)
	d.logger.Printf("storiesPerPreference: %s", storiesPerPreference)

	chillStoriesSettings := configs.ThemesSettings()
	var promptText string

	for i := 0; i < storiesPerPreference; i++ {
		topicNumber, err := configs.RandomFrom(chillStoriesSettings.ChillStoriesList)
		if err != nil {
			topicNumber = 0
		}
		topic := chillStoriesSettings.ChillStoriesList[topicNumber]
		if i < storiesPerPreference - 1 {
			promptText += topic + ", "
		} else {
			promptText += topic
		}	
	}
	
	superPrompt := fmt.Sprintf(
		"Create one topic for each topic in the list "+ promptText + ". that TEACH VALUES and Courage "+
			"The topics should illustrate a journey of %s. "+
			"Each topic should be have exactly two parts title and description."+
			"title should be a short and catchy title that kids can understand."+
			"description should also be short and concise."+
			"seperate title and description with a colon. and maintain only one colon in the whole string."+
			"Return the topics as a list of strings and it should be in title:description format."+
			"Always validate the length of the topics should be alwys %d.",
		preference,
		storiesPerPreference,
	)
	
	d.logger.Printf("Generated prompt: %s", superPrompt)
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
