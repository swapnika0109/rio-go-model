package helpers

import (
	"context"
	"fmt"
	"rio-go-model/configs"
	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"
)

type StoryAudioCrud struct {
	db             *database.StoryDatabase
	storageService *database.StorageService
	storyGenerator *StoryGenerationHelper
	logger         *util.CustomLogger
}

func NewStoryAudioCrud(db *database.StoryDatabase, storageService *database.StorageService, storyGenerator *StoryGenerationHelper) *StoryAudioCrud {
	logger := util.GetLogger("story_audio_crud", configs.GetSettings())
	return &StoryAudioCrud{db: db, storageService: storageService, storyGenerator: storyGenerator, logger: logger}

}

func (s *StoryAudioCrud) ResetAudioByThemeID(ctx context.Context, themeID string) error {
	s.logger.Infof("Resetting audio by theme id: %s", themeID)
	stories, err := s.db.GetStoryByThemeID(ctx, themeID)
	if err != nil {
		return fmt.Errorf("error getting stories by theme id: %v", err)
	}

	for _, story := range stories {
		s.storageService.DeleteFile(story["audio_url"].(string))
		var audioData []byte
		language := story["language"].(string)
		theme := story["theme"].(string)
		suspended, err := s.db.SuspendAudioAPI(ctx, "audio")
		if language != "Telugu" && (suspended || err != nil) {
			if err != nil {
				s.logger.Errorf("Failed to read audio api trigger: %v", err)
			} else {
				s.logger.Errorf("Google Audio API trigger is suspended; using fallback audio generator")
			}
			audioData, err = s.storyGenerator.audioGenerator.GenerateAudio(story["story_text"].(string))
		} else {
			s.logger.Infof("Using Google Audio API to generate story audio...")
			s.logger.Infof("Language: %s, Story text length: %d", language, len(story["story_text"].(string)))

			if s.storyGenerator == nil {
				s.logger.Errorf("storyGenerator is nil!")
				continue
			}
			if s.storyGenerator.audioStoryGenerator == nil {
				s.logger.Errorf("audioStoryGenerator is nil!")
				continue
			}

			audioData, _, err = s.storyGenerator.audioStoryGenerator.GenerateAudioAdapter(story["story_text"].(string), language, theme)
			if err != nil {
				s.logger.Errorf("GenerateAudioAdapter failed: %v", err)
			} else {
				s.logger.Infof("GenerateAudioAdapter succeeded, audio data length: %d", len(audioData))
			}
		}
		if err != nil {
			s.logger.Errorf("Failed to generate audio file: %v", err)
			continue
		}
		url, err := s.storageService.UploadFile(audioData, "audio", "wav")
		if err != nil {
			s.logger.Errorf("Failed to upload audio file: %v", err)
			continue
		}
		story["audio_url"] = url
		s.db.UpdateStory(ctx, story["story_id"].(string), story)

	}

	return nil
}
