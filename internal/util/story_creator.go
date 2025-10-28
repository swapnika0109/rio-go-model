package util

import (
	"log"
	"regexp"
	"strings"

	"rio-go-model/configs/english"
	"rio-go-model/configs/telugu"
)

// AIMessage represents a message in the AI conversation
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StoryResponse represents the response for story generation
type StoryResponse struct {
	Story string `json:"story,omitempty"`
	Error string `json:"error,omitempty"`
}

// generateFormattedPrompt generates the formatted prompt based on version and parameters
func GenerateFormattedPrompt(theme, topic string, kwargs map[string]interface{}) (string, string, error) {
	var formattedPrompt, systemMessage string
	var promptTemplate *PromptTemplate
	// Extract parameters
	country := getStringFromMap(kwargs, "country", "")
	city := getStringFromMap(kwargs, "city", "")
	religion := getStringFromMap(kwargs, "religions", "")
	preference := getStringFromMap(kwargs, "preferences", "")
	log.Printf("Generated preferences: %v", preference)

	switch theme {
	case "1":
		if kwargs["language"] == "Telugu" {
			cfg := telugu.PlanetProtectorPromptConfig(topic, country, city)
			promptTemplate = &PromptTemplate{
				Prompt: cfg.Prompt,
				System: cfg.System,
			}
			formattedPrompt = promptTemplate.Prompt
			systemMessage = promptTemplate.System
		} else {
			cfg := english.PlanetProtectorPromptConfig(topic, country, city)
			promptTemplate = &PromptTemplate{
				Prompt: cfg.Prompt,
				System: cfg.System,
			}
			formattedPrompt = promptTemplate.Prompt
			systemMessage = promptTemplate.System
		}

	case "2":
		if kwargs["language"] == "Telugu" {
			cfg := telugu.MindfulStoriesPromptConfig(topic, religion)
			promptTemplate = &PromptTemplate{
				Prompt: cfg.Prompt,
				System: cfg.System,
			}
		} else {
			cfg := english.MindfulStoriesPromptConfig(topic, religion)
			promptTemplate = &PromptTemplate{
				Prompt: cfg.Prompt,
				System: cfg.System,
			}
		}
		formattedPrompt = promptTemplate.Prompt
		systemMessage = promptTemplate.System

	case "3":
		if kwargs["language"] == "Telugu" {
			cfg := telugu.ChillStoriesPromptConfig(topic)
			promptTemplate = &PromptTemplate{
				Prompt: cfg.Prompt,
				System: cfg.System,
			}
		} else {
			cfg := english.ChillStoriesPromptConfig(topic)
			promptTemplate = &PromptTemplate{
				Prompt: cfg.Prompt,
				System: cfg.System,
			}
		}
		formattedPrompt = promptTemplate.Prompt
		systemMessage = promptTemplate.System
	}

	// Add preference-specific content
	preference = strings.ToUpper(preference)
	if kwargs["language"] == "Telugu" {
		prefContent := telugu.Preferences()[strings.ToUpper(preference)]
		if prefContent != "" {
			formattedPrompt += prefContent
		}
	} else {
		log.Printf("Generated preference: %s", preference)
		prefContent := english.Preferences()[strings.ToUpper(preference)]
		log.Printf("Generated preference content: %s", prefContent)
		if prefContent != "" {
			formattedPrompt += prefContent
			log.Printf("Generated preference content: %s", formattedPrompt)
		}
	}

	log.Printf("Generated prompt 3: %s", formattedPrompt)
	return formattedPrompt, systemMessage, nil
}

// Configuration structures

// PromptTemplate represents a prompt template configuration
type PromptTemplate struct {
	Prompt string `json:"prompt"`
	System string `json:"system"`
}

// Helper functions for map operations

// getStringFromMap safely extracts a string value from a map
func getStringFromMap(m map[string]interface{}, key, defaultValue string) string {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

// getStringSliceFromMap safely extracts a string slice from a map
func getStringSliceFromMap(m map[string]interface{}, key string) []string {
	if value, exists := m[key]; exists {
		if slice, ok := value.([]string); ok {
			return slice
		}
		// Handle interface{} slice
		if interfaceSlice, ok := value.([]interface{}); ok {
			result := make([]string, 0)
			for _, item := range interfaceSlice {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return []string{}
}

// parseTopics parses topics from the AI response
func ParseTopics(topicsData string, storiesPerPreference int) []string {
	topicsDataLngth := strings.Split(topicsData, "[")
	if len(topicsDataLngth) > 1 {
		topicsData = topicsDataLngth[1]
	} else {
		topicsData = topicsDataLngth[0]
	}
	topicsData = strings.Split(topicsData, "]")[0]
	topicsList := strings.Split(topicsData, ";")
	topics := make([]string, 0)

	for _, topic := range topicsList {
		topic = strings.TrimSpace(topic)
		if topic == "" || topic == "[" || topic == "]" {
			continue
		}

		if strings.Contains(topic, "]") {
			topic = strings.ReplaceAll(topic, "]", "")
			topic = strings.TrimSpace(topic)
		}

		if strings.Contains(topic, "[") {
			topic = strings.ReplaceAll(topic, "[", "")
			topic = strings.TrimSpace(topic)
		}

		topicValidationStr := strings.Split(topic, ":")
		if len(topicValidationStr) > 2 {
			continue
		}
		if len(topicValidationStr) == 2 {
			topic = topicValidationStr[1]
		}

		// Skip introductory text
		// if strings.Contains(topic, "ఖచ్చితంగా") || strings.Contains(topic, "టాపిక్స్") ||
		// 	strings.Contains(topic, "ఇక్కడ ఉన్నాయి") || strings.Contains(topic, "అందించండి") {
		// 	continue
		// }

		// Handle markdown formatted topics (e.g., "**1. Title: Description**")
		if strings.Contains(topic, "**") {
			// Remove markdown formatting
			cleanTopic := strings.ReplaceAll(topic, "**", "")

			// Remove numbering (e.g., "1. ", "2. ", etc.)
			re := regexp.MustCompile(`^\d+\.\s*`)
			cleanTopic = re.ReplaceAllString(cleanTopic, "")

			// Append as a single topic string (no title:description required)
			cleanTopic = strings.TrimSpace(cleanTopic)
			if len(cleanTopic) > 0 {
				topics = append(topics, cleanTopic)
			}
			continue
		}

		// Handle quoted topics
		if strings.Contains(topic, `"`) {
			parts := strings.Split(topic, `"`)
			if len(parts) > 1 {
				finalTopic := strings.TrimSpace(parts[1])
				if len(finalTopic) > 0 {
					topics = append(topics, finalTopic)
				}
			}
		} else {
			// Handle numbered topics without markdown (e.g., "1. Some Topic Text")
			re := regexp.MustCompile(`^\d+\.\s*`)
			cleanTopic := re.ReplaceAllString(topic, "")
			cleanTopic = strings.TrimSpace(cleanTopic)

			if len(cleanTopic) > 0 {
				topics = append(topics, cleanTopic)
			}
		}
	}
	if len(topics) > storiesPerPreference {
		topics = topics[:storiesPerPreference]
	}
	return topics
}
