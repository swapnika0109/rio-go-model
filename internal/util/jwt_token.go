package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"rio-go-model/configs"
)

type HttpError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Status, e.Message)
}


func VerifyToken(token string) (string, string, error) {
	username, email, err := validateToken(token)
	if err != nil {
		return "", "", err
	}

	return username, email, nil
}

func validateToken(token string) (string, string, error) {
	var username, email string
	if token != "" {
		secretKey := configs.LoadSettings().SecretKey
		token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			return "", "", fmt.Errorf("failed to parse JWT: %v", err)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			email = claims["email"].(string)
			username = claims["username"].(string)
		}
		
	}
	 
	return username, email, nil
}

func ValidateGoogleToken(token string) (string, string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %v", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+token)

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", &HttpError{
			Status:  401,
			Message: "Invalid Google access token",
		}
	}

	// Parse response
	userInfo := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return "", "", &HttpError{
			Status:  401,
			Message: fmt.Sprintf("Failed to parse response: %v", err),
		}
	}

	// Extract email and name
	email, ok := userInfo["email"].(string)
	if !ok || email == "" {
		return "", "", &HttpError{
			Status:  401,
			Message: "No email in Google user info",
		}
	}

	name, _ := userInfo["name"].(string)

	// Generate username
	var username string
	if name != "" {
		username = name
	} else {
		emailParts := strings.Split(email, "@")
		if len(emailParts) > 0 {
			username = emailParts[0]
		}
	}

	return username, email, nil
}
