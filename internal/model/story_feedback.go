package model

type StoryFeedback struct {
	Like bool   `json:"like"`
	StoryId string `json:"storyId"`
	Email     string `json:"email"`
}

func NewStoryFeedback(like bool, storyId string, email string) *StoryFeedback {
	return &StoryFeedback{
		Like: like,
		StoryId: storyId,
		Email: email,
	}
}
