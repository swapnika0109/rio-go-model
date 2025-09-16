package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"log"
	"net/http"
	"time"

	"rio-go-model/internal/services/database"
	"rio-go-model/internal/util"
)

// AuthHandler handles authentication-related API requests.
type AuthHandler struct {
	db *database.StoryDatabase
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(db *database.StoryDatabase) *AuthHandler {
	return &AuthHandler{db: db}
}

// GoogleLoginRequest represents the expected JSON body for the Google login endpoint.
type GoogleLoginRequest struct {
	AccessToken string `json:"access_token"`
}

// RefreshTokenRequest represents the expected JSON body for the token refresh endpoint.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh"`
}

// GoogleLogin handles the server-side flow for Google OAuth2 authentication.
// @Summary      Google Login
// @Description  Accepts a Google access token, validates it, and returns a local JWT pair (access and refresh tokens). If the user doesn't exist, a new user profile is created.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        google_login_request body GoogleLoginRequest true "Google Access Token"
// @Success      200 {object} util.TokenPair
// @Failure      400 {object} util.HttpError "Invalid request body or missing access_token"
// @Failure      401 {object} util.HttpError "Invalid Google token"
// @Failure      500 {object} util.HttpError "Internal server error"
// @Router       /auth/google [post]
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	var req GoogleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.AccessToken == "" {
		http.Error(w, "access_token is required", http.StatusBadRequest)
		return
	}

	// 1. Validate the Google access token by calling Google's tokeninfo endpoint.
	googleUser, googleEmail, err := util.ValidateGoogleToken(req.AccessToken)
	if err != nil {
		log.Printf("ERROR: Invalid Google token: %v", err)
		http.Error(w, fmt.Sprintf("Invalid Google token: %v", err), http.StatusUnauthorized)
		return
	}

	// 2. Check if a user with this email exists in our database.
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// userProfile, err := h.db.GetUserProfileByEmail(ctx, googleUser.Email)
	// if err != nil {
	// 	log.Printf("ERROR: Failed to query user profile: %v", err)
	// 	http.Error(w, "Internal server error during user lookup", http.StatusInternalServerError)
	// 	return
	// }

	username := googleUser // Use Google's given name as the username

	// 3. If user does not exist, create a new profile.
	// if userProfile == nil {
	// 	log.Printf("User with email %s not found. Creating new profile.", googleUser.Email)
	// 	profileData := map[string]interface{}{
	// 		"username":          username,
	// 		"email":             googleUser.Email,
	// 		"processing_status": "not_started",
	// 		// Add other default fields as needed
	// 	}
	// 	if _, err := h.db.CreateUserProfile(ctx, profileData); err != nil {
	// 		log.Printf("ERROR: Failed to create user profile: %v", err)
	// 		http.Error(w, "Internal server error during user creation", http.StatusInternalServerError)
	// 		return
	// 	}
	// } else {
	// 	username = userProfile["username"].(string) // Use existing username if profile exists
	// }

	// 4. Get or create user profile and get token version
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	userProfile, err := h.db.GetUserProfileByEmail(ctx, googleEmail)
	if err != nil {
		log.Printf("ERROR: Failed to query user profile: %v", err)
		http.Error(w, "Internal server error during user lookup", http.StatusInternalServerError)
		return
	}

	var tokenVersion int64 = 0
	if userProfile != nil {
		if version, exists := userProfile["token_version"]; exists {
			if v, ok := version.(int64); ok {
				tokenVersion = v
			}
		}
	}

	// 5. Issue our own JWT access and refresh tokens.
	tokenPair, err := util.GenerateTokens(username, googleEmail, tokenVersion)
	if err != nil {
		log.Printf("ERROR: Failed to generate local JWTs: %v", err)
		http.Error(w, "Internal server error during token generation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
}

// RefreshToken validates a refresh token and issues a new access token.
// @Summary      Refresh Access Token
// @Description  Accepts a refresh token and returns a new, short-lived access token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        refresh_token_request body RefreshTokenRequest true "Refresh Token"
// @Success      200 {object} util.TokenPair "{\"access\":\"new_access_token\", \"refresh\":\"original_refresh_token\"}"
// @Failure      400 {object} util.HttpError "Invalid request body or missing refresh token"
// @Failure      401 {object} util.HttpError "Invalid refresh token"
// @Failure      500 {object} util.HttpError "Failed to generate new access token"
// @Router       /auth/token/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "refresh token is required", http.StatusBadRequest)
		return
	}

	// 1. Validate the refresh token. We can expand validateToken to check for "refresh" type.
	_, _, err := util.VerifyToken(req.RefreshToken) // VerifyToken internally calls validateToken
	if err != nil {
		log.Printf("ERROR: Invalid refresh token provided: %v", err)
		http.Error(w, fmt.Sprintf("Invalid refresh token: %v", err), http.StatusUnauthorized)
		return
	}
	
	// (Optional check for token type if you added it to validateToken)
	
	// 2. Generate a new access token.
	newAccessToken, err := util.GenerateAccessTokenFromRefresh(req.RefreshToken)
	if err != nil {
		log.Printf("ERROR: Failed to generate new access token: %v", err)
		http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access":  newAccessToken,
		"refresh": req.RefreshToken,
	})
}

// LogoutRequest represents the expected JSON body for the logout endpoint.
type LogoutRequest struct {
	Refresh string `json:"refresh"`
}

// Logout handles user logout by invalidating the refresh token.
// @Summary      User Logout
// @Description  Logs out a user by invalidating their refresh token. This is a client-side operation that clears tokens from storage.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        logout_request body LogoutRequest true "Refresh Token to Invalidate"
// @Success      200 {object} map[string]string "Logout successful"
// @Failure      400 {object} util.HttpError "Invalid request body or missing refresh token"
// @Failure      500 {object} util.HttpError "Internal server error"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Refresh == "" {
		http.Error(w, "refresh token is required", http.StatusBadRequest)
		return
	}

	// Validate the refresh token to get user email
	_, email, err := util.VerifyToken(req.Refresh)
	if err != nil {
		log.Printf("WARNING: Invalid refresh token during logout: %v", err)
		// Don't return error for invalid token during logout - just log it
	} else {
		// Increment token version to invalidate all existing tokens for this user
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		if err := h.db.IncrementTokenVersion(ctx, email); err != nil {
			log.Printf("WARNING: Failed to increment token version: %v", err)
		} else {
			log.Printf("Token version incremented for user: %s", email)
		}
	}
	
	log.Printf("User logged out successfully")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}
