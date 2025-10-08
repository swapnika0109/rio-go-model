package model

// StoryResponse represents the response for story generation
type StoryResponse struct {
	Story string `json:"story,omitempty"`
	Error string `json:"error,omitempty"`
}
