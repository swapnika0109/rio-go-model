package handlers

import (
	"encoding/json"
	"net/http"
	"rio-go-model/internal/helpers/tc"
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
// @Param accepted body bool true "Accepted"
// @Success 200 {object} map[string]string "TC created successfully"
// @Failure 401 {object} util.HttpError "Unauthorized"
// @Failure 500 {object} util.HttpError "Internal Server Error"
// @Router /tc [get]
func (h *TcHandler) TcHandler(w http.ResponseWriter, r *http.Request) {
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
	
	accepted := r.URL.Query().Get("accepted")
	_, err = tc.CreateTc(r.Context(), h.tcDB, accepted == "true", email)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "TC created successfully"})
	w.WriteHeader(http.StatusOK)
	return
}