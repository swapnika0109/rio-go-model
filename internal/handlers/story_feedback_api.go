package handlers

import (
	"encoding/json"
	"net/http"
	"rio-go-model/internal/model"
	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"
	"log"
	// "strings"
)

type StoryFeedbackHandler struct {
	storyFeedbackDB *database.StoryDatabase
	logger *log.Logger
}

func NewStoryFeedbackHandler(storyFeedbackDB *database.StoryDatabase) *StoryFeedbackHandler {
	return &StoryFeedbackHandler{
		storyFeedbackDB: storyFeedbackDB,
		logger: log.New(log.Writer(), "[Story Feedback Service] ", log.LstdFlags|log.Lshortfile),
	}
}

// StoryFeedbackHandler is the handler for the story feedback API
// @Summary Create story feedback
// @Description Create story feedback
// @Tags story-feedback
// @Accept json
// @Produce json
// @Param storyFeedback body model.StoryFeedback true "StoryFeedback request"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Story feedback created successfully"
// @Failure 401 {object} util.HttpError "Unauthorized"
// @Failure 500 {object} util.HttpError "Internal Server Error"
// @Router /story-feedback [post]
func (h *StoryFeedbackHandler) HandleStoryFeedback(w http.ResponseWriter, r *http.Request) {
	_, email, err := util.VerifyAuth(r)
	if err != nil {
		h.logger.Printf("WARNING: Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var storyFeedback model.StoryFeedback
	if err := json.NewDecoder(r.Body).Decode(&storyFeedback); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	storyFeedback.Email = email
	_, err = h.storyFeedbackDB.CreateStoryFeedback(r.Context(), &storyFeedback)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Story feedback created successfully"})
	w.WriteHeader(http.StatusOK)
	return
}