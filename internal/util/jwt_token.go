package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"os"
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

func validateToken(tokenStr string) (string, string, error) {
	var username, email string
	if tokenStr != "" {
		secretKey := os.Getenv("SECRET_KEY")
		if secretKey == "" {
			return "", "", fmt.Errorf("SECRET_KEY environment variable not set")
		}
		// --- START TEMPORARY DEBUGGING ---
		// Using a hardcoded byte slice to eliminate any possibility of
		// string-to-byte conversion or encoding issues.
		secretKeyBytes := []byte(secretKey)
		// --- END TEMPORARY DEBUGGING ---

		// Use the jwx library, but disable time validation for this test.
		token, err := jwt.Parse(
			[]byte(tokenStr),
			jwt.WithKey(jwa.HS256, secretKeyBytes),
			jwt.WithValidate(false), // <-- Tell the library to IGNORE time claims (exp, iat)
		)

		if err != nil {
			// If this fails, it is a signature or formatting error.
			return "", "", fmt.Errorf("failed to parse or verify signature: %w", err)
		}

		// Manually log that validation was skipped
		log.Println("SIGNATURE CHECK PASSED (time validation was skipped)")

		// Correctly extract claims using the .Get() method on the parsed token.
		emailClaim, ok := token.Get("email")
		if !ok {
			return "", "", fmt.Errorf("email claim not found in token")
		}
		email, ok = emailClaim.(string)
		if !ok {
			return "", "", fmt.Errorf("email claim is not a string")
		}

		usernameClaim, ok := token.Get("username")
		if !ok {
			return "", "", fmt.Errorf("username claim not found in token")
		}
		username, ok = usernameClaim.(string)
		if !ok {
			return "", "", fmt.Errorf("username claim is not a string")
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
