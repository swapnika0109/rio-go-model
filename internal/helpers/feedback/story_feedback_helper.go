package feedback

import (
	"context"
	"log"
	"rio-go-model/internal/model"
	"rio-go-model/internal/services/database"
)

func CreateStoryFeedback(ctx context.Context, db *database.StoryDatabase, storyId string, like bool, email string) (string, error) {
	log.Println("Creating story feedback for storyId: ", storyId, " and email: ", email)
	storyFeedback := model.NewStoryFeedback(like, storyId, email)
	if _, err := db.CreateStoryFeedback(ctx, storyFeedback); err != nil {
		log.Println("Error creating story feedback: ", err)
		return "", err
	}
	return "Story feedback created successfully", nil
}