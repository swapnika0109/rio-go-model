package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"rio-go-model/internal/model"
	"rio-go-model/internal/services/emails"
	"rio-go-model/internal/util"
)

type Email struct {
	EmailRequest *model.EmailRequest
	logger       *log.Logger
}

// @Tags Email
// @Summary Sends an email to the user
// @Description Sends an email to the user
// @Accept json
// @Param emailRequest body model.EmailRequest true "Email request"
// @Success 200 {string} string "Email sent successfully"
// @Failure 400 {string} string "Invalid JSON body"
// @Failure 500 {string} string "Failed to send email"
// @Router /email [post]
// NewEmail handles POST /email and expects JSON body: {"name":"...","email":"...","message":"..."}
func (e *Email) NewEmail(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	e.logger = log.New(log.Writer(), "[Email Service] ", log.LstdFlags|log.Lshortfile)
	_, _, _, err := util.VerifyAuth(r)
	if err != nil {
		e.logger.Printf("WARNING: Invalid token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	var emailRequest model.EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&emailRequest); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if emailRequest.Name == "" || emailRequest.Email == "" || emailRequest.Message == "" {
		http.Error(w, "name, email and message are required", http.StatusBadRequest)
		return
	}

	e.EmailRequest = &emailRequest
	emailService := emails.NewEmailService(emailRequest.Name, emailRequest.Email, emailRequest.Message)
	if err := emailService.SendEmail(); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","message":"Email sent successfully"}`))
}
