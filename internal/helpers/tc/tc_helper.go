package tc

import (
	"context"
	"log"
	"rio-go-model/internal/model"
	"rio-go-model/internal/services/database"
)

func CreateTc(ctx context.Context, db *database.StoryDatabase, accepted bool, email string) (string, error) {
	log.Println("Creating TC for email: ", email)
	tc := model.NewTc(accepted, email)
	if _, err := db.CreateTc(ctx, tc); err != nil {
		log.Println("Error creating TC: ", err)
		return "", err
	}
	return "TC created successfully", nil
}