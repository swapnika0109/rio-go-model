package storymonitor

import (
	"context"
	"fmt"
	"rio-go-model/configs"
	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"

	"cloud.google.com/go/firestore"
)

type AdminRepo struct {
	db      *database.StoryDatabase
	storage *database.StorageService
	logger  *util.CustomLogger
}

func NewAdminRepo(db *database.StoryDatabase, storage *database.StorageService) *AdminRepo {
	return &AdminRepo{db: db, storage: storage, logger: util.GetLogger("admin_repo", configs.GetSettings())}
}

func (a *AdminRepo) DeleteStoryByID(ctx context.Context, storyID string) error {
	a.logger.Infof("Deleting story by ID: %s", storyID)
	doc, err := a.db.GetClient().Collection(a.db.CollectionV2).Doc(storyID).Get(ctx)
	if err != nil {
		return fmt.Errorf("error getting story by ID: %v", err)
	}
	if !doc.Exists() {
		return fmt.Errorf("story not found")
	}
	storyData := doc.Data()
	audioUrl := storyData["audio_url"].(string)
	imageUrl := storyData["image_url"].(string)
	if audioUrl != "" {
		a.logger.Infof("Deleting audio file: %s", audioUrl)
		a.storage.DeleteFile(audioUrl)
	}
	if imageUrl != "" {
		a.logger.Infof("Deleting image file: %s", imageUrl)
		a.storage.DeleteFile(imageUrl)
	}
	a.db.DeleteStory(ctx, storyID)
	theme_id := storyData["theme_id"].(string)
	theme := storyData["theme"].(string)
	var docs []*firestore.DocumentSnapshot
	switch theme {
	case "1":
		docs, err = a.db.GetClient().Collection(a.db.MdCollection1).Where("theme_id", "==", theme_id).Documents(ctx).GetAll()
		if err != nil {
			return fmt.Errorf("error getting theme 1 by ID: %v", err)
		}

	case "2":
		docs, err = a.db.GetClient().Collection(a.db.MdCollection2).Where("theme_id", "==", theme_id).Documents(ctx).GetAll()
		if err != nil {
			return fmt.Errorf("error getting theme 2 by ID: %v", err)
		}
	case "3":
		docs, err = a.db.GetClient().Collection(a.db.MdCollection3).Where("theme_id", "==", theme_id).Documents(ctx).GetAll()
		if err != nil {
			return fmt.Errorf("error getting theme 3 by ID: %v", err)
		}
	}
	for _, doc := range docs {
		docData := doc.Data()
		topics := docData["topics"].([]interface{})
		var finalTopics []string
		for _, topic := range topics {
			topicStr := topic.(string)
			if topicStr == storyData["title"].(string) {
				continue
			} else {
				finalTopics = append(finalTopics, topicStr)
				a.logger.Infof("Deleting topic: %s", topicStr)
			}
		}
		docData["topics"] = finalTopics
		a.logger.Infof("Updating metadata topics: %v", docData["topics"])
		switch theme {
		case "1":
			_, err = a.db.GetClient().Collection(a.db.MdCollection1).Doc(doc.Ref.ID).Set(ctx, docData)
			if err != nil {
				return fmt.Errorf("error updating metadata topics 1: %v", err)
			}
		case "2":
			_, err = a.db.GetClient().Collection(a.db.MdCollection2).Doc(doc.Ref.ID).Set(ctx, docData)
			if err != nil {
				return fmt.Errorf("error updating metadata topics 2: %v", err)
			}
		case "3":
			_, err = a.db.GetClient().Collection(a.db.MdCollection3).Doc(doc.Ref.ID).Set(ctx, docData)
			if err != nil {
				return fmt.Errorf("error updating metadata topics 3: %v", err)
			}
		}
	}
	return nil
}

// func (a *AdminRepo) MigrateAllStories(ctx context.Context) error {
// 	a.logger.Infof("Migrating all stories")
// 	docs, err := a.db.GetClient().Collection(a.db.CollectionV2).Documents(ctx).GetAll()
// 	if err != nil {
// 		return fmt.Errorf("error getting all stories: %v", err)
// 	}
// 	for _, doc := range docs {
// 		log.Printf("Migrating story: %s", doc.Ref.ID)
// 		docData := doc.Data()
// 		docData["story_type"] = "Premium"
// 		_, err = a.db.GetClient().Collection(a.db.CollectionV2).Doc(doc.Ref.ID).Set(ctx, docData)
// 		if err != nil {
// 			return fmt.Errorf("error updating story: %v", err)
// 		}
// 	}

// 	return nil
// }
