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

	log.Printf("PubSub messageId=%s attrs=%v data=%s", push.Message.MessageID, push.Message.Attributes, decoded)
	cost, budget, ok := extractBudgetNumbers(decoded)
	if !ok {
		log.Printf("warning: failed to extract budget numbers")
		w.WriteHeader(http.StatusOK)
		return
	}
	// If this is a Cloud Billing budget message, ignore until 90% threshold
	if shouldIgnoreByBudget(cost, budget) {
		w.WriteHeader(http.StatusOK)
		return
	}
	_, err = h.db.CreateAPITrigger(r.Context(), "gemini", budget, cost)
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

	log.Printf("PubSub messageId=%s attrs=%v data=%s", push.Message.MessageID, push.Message.Attributes, decoded)
	cost, budget, ok := extractBudgetNumbers(decoded)
	if !ok {
		log.Printf("warning: failed to extract budget numbers")
		w.WriteHeader(http.StatusOK)
		return
	}
	// If this is a Cloud Billing budget message, ignore until 90% threshold
	if shouldIgnoreByBudget(cost, budget) {
		w.WriteHeader(http.StatusOK)
		return
	}
	_, err = h.db.CreateAPITrigger(r.Context(), "audio", budget, cost)
	if err != nil {
		http.Error(w, "failed to create api trigger", http.StatusInternalServerError)
		return
	}
	// Return 200 quickly; retries happen if non-2xx is returned
	w.WriteHeader(http.StatusOK)
}

// shouldIgnoreByBudget inspects a (possibly noisy) billing JSON text and returns true
// if costAmount < 0.9 * budgetAmount. It tolerates extra timestamps by using regex.
func shouldIgnoreByBudget(cost float64, budget float64) bool {
	// Guard against negative or zero budget
	if budget <= 0 {
		return false
	}
	return cost < 0.9*budget
}

type PubSubData struct {
	DisplayName  string  `json:"budgetDisplayName"`
	Threshold    float64 `json:"alertThresholdExceeded"`
	CostAmount   float64 `json:"costAmount"`
	BudgetAmount float64 `json:"budgetAmount"`
	Currency     string  `json:"currencyCode"`
}

// extractBudgetNumbers parses JSON data to extract costAmount and budgetAmount
func extractBudgetNumbers(s []byte) (cost float64, budget float64, ok bool) {
	var data PubSubData
	if err := json.Unmarshal(s, &data); err != nil {
		log.Printf("DEBUG: JSON unmarshal failed: %v", err)
		return 0, 0, false
	}

	log.Printf("DEBUG: Parsed data - costAmount: %f, budgetAmount: %f", data.CostAmount, data.BudgetAmount)
	return data.CostAmount, data.BudgetAmount, true
}
