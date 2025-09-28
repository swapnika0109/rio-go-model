package handlers

import (
	"encoding/json"
	"net/http"
	"rio-go-model/internal/model"
	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"
	"strings"
)

type TcHandler struct {
	tcDB          *database.StoryDatabase
}

func NewTcHandler(tcDB *database.StoryDatabase) *TcHandler {
	return &TcHandler{
		tcDB: tcDB,
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
	return
}