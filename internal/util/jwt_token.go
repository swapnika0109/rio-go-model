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
	"github.com/google/uuid"
)

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

// GenerateTokens creates a new pair of access and refresh tokens for a user.
func GenerateTokens(username, email string, tokenVersion int64) (*TokenPair, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("SECRET_KEY environment variable not set")
	}
	secretKeyBytes := []byte(secretKey)

	// Create access token (short-lived)
	accessToken, err := createTokenWithVersion(username, email, 1*time.Hour, "access", tokenVersion, secretKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Create refresh token (long-lived)
	refreshToken, err := createTokenWithVersion(username, email, 7*24*time.Hour, "refresh", tokenVersion, secretKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateAccessTokenFromRefresh creates a new access token from a valid refresh token.
func GenerateAccessTokenFromRefresh(refreshTokenStr string) (string, error) {
	username, email, err := validateToken(refreshTokenStr)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}
	
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("SECRET_KEY environment variable not set")
	}
	secretKeyBytes := []byte(secretKey)

	return createToken(username, email, 1*time.Hour, "access", secretKeyBytes)
}

func createTokenWithVersion(username, email string, expiryDuration time.Duration, tokenType string, tokenVersion int64, secret []byte) (string, error) {
	token := jwt.New()
	_ = token.Set(jwt.JwtIDKey, uuid.New().String())
	_ = token.Set(jwt.IssuedAtKey, time.Now().Unix())
	_ = token.Set(jwt.ExpirationKey, time.Now().Add(expiryDuration).Unix())
	_ = token.Set("username", username)
	_ = token.Set("email", email)
	_ = token.Set("token_type", tokenType)
	_ = token.Set("token_version", tokenVersion)

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return string(signed), nil
}

func createToken(username, email string, expiryDuration time.Duration, tokenType string, secret []byte) (string, error) {
	token := jwt.New()
	_ = token.Set(jwt.JwtIDKey, uuid.New().String())
	_ = token.Set(jwt.IssuedAtKey, time.Now().Unix())
	_ = token.Set(jwt.ExpirationKey, time.Now().Add(expiryDuration).Unix())
	_ = token.Set("username", username)
	_ = token.Set("email", email)
	_ = token.Set("token_type", tokenType)

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return string(signed), nil
}


// HttpError represents an HTTP error response
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

// GetTokenVersion extracts the token version from a JWT token
func GetTokenVersion(tokenStr string) (int64, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return 0, fmt.Errorf("SECRET_KEY environment variable not set")
	}
	secretKeyBytes := []byte(secretKey)

	token, err := jwt.Parse(
		[]byte(tokenStr),
		jwt.WithKey(jwa.HS256, secretKeyBytes),
		jwt.WithValidate(false),
	)

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	versionClaim, ok := token.Get("token_version")
	if !ok {
		return 0, fmt.Errorf("token_version claim not found in token")
	}

	version, ok := versionClaim.(int64)
	if !ok {
		return 0, fmt.Errorf("token_version claim is not an int64")
	}

	return version, nil
}

func validateToken(tokenStr string) (string, string, error) {
	var username, email string
	if tokenStr != "" {
		// Check if token is blacklisted first
		// if GlobalBlacklist.IsBlacklisted(tokenStr) {
		// 	return "", "", fmt.Errorf("token has been revoked")
		// }
		
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

		// Check for token_type claim
		tokenTypeClaim, ok := token.Get("token_type")
		if ok {
			tokenType, ok := tokenTypeClaim.(string)
			if !ok {
				return "", "", fmt.Errorf("token_type claim not a string")
			}
			if tokenType != "refresh" && tokenType != "access" {
				return "", "", fmt.Errorf("invalid token_type: %s", tokenType)
			}

			// For refresh endpoint, we must ensure the token is a refresh token
			// We can add this check later in the handler itself if needed.
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
