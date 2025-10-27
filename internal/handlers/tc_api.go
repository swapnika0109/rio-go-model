package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rio-go-model/internal/model"
	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"
	"time"
	// "strings"
)

type TcHandler struct {
	tcDB   *database.StoryDatabase
	logger *log.Logger
}

func NewTcHandler(tcDB *database.StoryDatabase) *TcHandler {
	return &TcHandler{
		tcDB:   tcDB,
		logger: log.New(log.Writer(), "[TC Service] ", log.LstdFlags|log.Lshortfile),
	}
}

// TcHandler is the handler for the TC API
// @Summary Create TC
// @Description Create TC
// @Tags tc
// @Accept json
// @Produce json
// @Param tc body model.Tc true "Tc request"
// @Security BearerAuth
// @Success 200 {object} map[string]string "TC created successfully"
// @Failure 401 {object} util.HttpError "Unauthorized"
// @Failure 500 {object} util.HttpError "Internal Server Error"
// @Router /tc [post]
func (h *TcHandler) HandleTc(w http.ResponseWriter, r *http.Request) {
	_, email, tokenVersion, err := util.VerifyAuth(r)
	if err != nil {
		h.logger.Printf("WARNING: Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userTokenVersion, err := h.tcDB.GetTokenVersion(ctx, email)
	if err != nil {
		h.logger.Printf("ERROR: Failed to get token version: %v", err)
		http.Error(w, "Failed to get token version", http.StatusInternalServerError)
		return
	}
	err = util.VerifyUserTokenVersion(tokenVersion, userTokenVersion)
	if err != nil {
		h.logger.Printf("❌ DEBUG: Token version mismatch: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	var tc model.Tc
	if err := json.NewDecoder(r.Body).Decode(&tc); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	tc.Email = email

	_, err = h.tcDB.CreateTc(r.Context(), &tc)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "TC created successfully"})
}
