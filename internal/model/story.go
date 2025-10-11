package model

// StoryResponse represents the response for story generation
type StoryResponse struct {
	Story string `json:"story,omitempty"`
	Error string `json:"error,omitempty"`
}

type PromptEngineConfig struct {
	System string `json:"system"`
	Prompt string `json:"prompt"`
}

type TopicResponse struct {
	Title []string `json:"title,omitempty"`
	Error string   `json:"error,omitempty"`
}
