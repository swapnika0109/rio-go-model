package handlers

import (
	"encoding/json"
	"net/http"
	"rio-go-model/internal/model"
	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"
	"strings"
)

type StoryFeedbackHandler struct {
	storyFeedbackDB *database.StoryDatabase
}

func NewStoryFeedbackHandler(storyFeedbackDB *database.StoryDatabase) *StoryFeedbackHandler {
	return &StoryFeedbackHandler{
		storyFeedbackDB: storyFeedbackDB,
	}
}

// StoryFeedbackHandler is the handler for the story feedback API
// @Summary Create story feedback
// @Description Create story feedback
// @Tags story-feedback
// @Accept json
// @Produce json
// @Param storyId path string true "Story ID"
// @Param like body bool true "Like"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Story feedback created successfully"
// @Failure 401 {object} util.HttpError "Unauthorized"
// @Failure 500 {object} util.HttpError "Internal Server Error"
// @Router /story-feedback [post]
func (h *StoryFeedbackHandler) HandleStoryFeedback(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}


	token := strings.TrimPrefix(authHeader, "Bearer ")
	_, email, err := util.VerifyToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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