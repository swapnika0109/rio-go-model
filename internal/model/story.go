package model

// StoryResponse represents the response for story generation
type StoryResponse struct {
	Story       string `json:"story,omitempty"`
	TotalTokens int32  `json:"total_tokens,omitempty"`
	Error       string `json:"error,omitempty"`
}

type PromptEngineConfig struct {
	System string `json:"system"`
	Prompt string `json:"prompt"`
}

type TopicResponse struct {
	Title       []string `json:"title,omitempty"`
	TotalTokens int32    `json:"total_tokens,omitempty"`
	Error       string   `json:"error,omitempty"`
}
