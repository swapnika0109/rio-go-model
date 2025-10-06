package handlers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type PubSubPushMessage struct {
	Message      PubSubMessage `json:"message"`
	Subscription string        `json:"subscription"`
}

type PubSubMessage struct {
	Data        string            `json:"data"`        // base64-encoded
	Attributes  map[string]string `json:"attributes"`  // optional
	MessageID   string            `json:"messageId"`
	PublishTime string            `json:"publishTime"`
}

//@Summary PubSub Push Handler
//@Description Receives GCP Pub/Sub push messages and acknowledges with 200
//@Tags Triggers
//@Accept json
//@Produce json
//@Param request body PubSubPushMessage true "PubSub Push Message"
//@Success 200 {string} string "OK"
//@Failure 400 {string} string "Bad Request"
//@Failure 500 {string} string "Internal Server Error"
//@Router /triggers/pubsub [post]
func PubSubPushHandler(w http.ResponseWriter, r *http.Request) {
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

	// Return 200 quickly; retries happen if non-2xx is returned
	w.WriteHeader(http.StatusOK)
}