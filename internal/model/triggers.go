package model

import "time"

type GeminiTriggerRequest struct {
	Email string `json:"email"`
	Trigger string `json:"trigger"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status string `json:"status"`
}

type AudioChirpTriggerResponse struct {
	Email string `json:"email"`
	Trigger string `json:"trigger"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status string `json:"status"`
}	