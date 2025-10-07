package handlers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"rio-go-model/internal/services/database"
)

type PubSubHandler struct {
	db *database.StoryDatabase
}

func NewPubSubHandler(db *database.StoryDatabase) *PubSubHandler {
	return &PubSubHandler{db: db}
}

type PubSubPushMessage struct {
	Message      PubSubMessage `json:"message"`
	Subscription string        `json:"subscription"`
}

type PubSubMessage struct {
	Data        string            `json:"data"`       // base64-encoded
	Attributes  map[string]string `json:"attributes"` // optional
	MessageID   string            `json:"messageId"`
	PublishTime string            `json:"publishTime"`
}

// @Summary PubSub Push Handler
// @Description Receives GCP Pub/Sub push messages and acknowledges with 200
// @Tags Triggers
// @Accept json
// @Produce json
// @Param request body PubSubPushMessage true "PubSub Push Message"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /triggers/gemini/pubsub [post]
func (h *PubSubHandler) PubSubPushGeminiHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("PubSubPushGeminiHandler called..")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var push PubSubPushMessage
	if err := json.Unmarshal(body, &push); err != nil {
		http.Error(w, "invalid pubsub push payload", http.StatusBadRequest)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(push.Message.Data)
	if err != nil {
		// Fallback for manual testing: treat data as plain text if not valid base64
		log.Printf("warning: base64 decode failed, using raw data: %v", err)
		decoded = []byte(push.Message.Data)
	}

	log.Printf("PubSub messageId=%s attrs=%v data=%s", push.Message.MessageID, push.Message.Attributes, string(decoded))
	_, err = h.db.CreateAPITrigger(r.Context(), "gemini")
	if err != nil {
		http.Error(w, "failed to create api trigger", http.StatusInternalServerError)
		return
	}
	// Return 200 quickly; retries happen if non-2xx is returned
	w.WriteHeader(http.StatusOK)
}

// @Summary PubSub Push Handler
// @Description Receives GCP Pub/Sub push messages and acknowledges with 200
// @Tags Triggers
// @Accept json
// @Produce json
// @Param request body PubSubPushMessage true "PubSub Push Message"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /triggers/audio/pubsub [post]
func (h *PubSubHandler) PubSubPushAudioHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("PubSubPushAudioHandler called..")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var push PubSubPushMessage
	if err := json.Unmarshal(body, &push); err != nil {
		http.Error(w, "invalid pubsub push payload", http.StatusBadRequest)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(push.Message.Data)
	if err != nil {
		// Fallback for manual testing: treat data as plain text if not valid base64
		log.Printf("warning: base64 decode failed, using raw data: %v", err)
		decoded = []byte(push.Message.Data)
	}

	log.Printf("PubSub messageId=%s attrs=%v data=%s", push.Message.MessageID, push.Message.Attributes, string(decoded))
	_, err = h.db.CreateAPITrigger(r.Context(), "audio")
	if err != nil {
		http.Error(w, "failed to create api trigger", http.StatusInternalServerError)
		return
	}
	// Return 200 quickly; retries happen if non-2xx is returned
	w.WriteHeader(http.StatusOK)
}
