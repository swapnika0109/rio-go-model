package database

import (
	"context"
	"fmt"
	"log"
	"os"

	// "path/filepath"
	"strings"
	"time"

	"rio-go-model/internal/model"
	"rio-go-model/internal/util"

	"rio-go-model/internal/services/tts"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	// "rio-go-model/configs"
)

// getUTCTimestamp returns the current time in UTC
func getUTCTimestamp() time.Time {
	return time.Now().UTC()
}

// getNextMonthResetTime returns the 2nd day of the next month at 00:00:00 UTC
func getNextMonthResetTime() time.Time {
	now := time.Now().UTC()

	// Get the first day of next month
	nextMonth := now.AddDate(0, 1, 0) // Add 1 month
	firstDayOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Add 1 day to get the 2nd day of next month
	secondDayOfNextMonth := firstDayOfNextMonth.AddDate(0, 0, 1)

	return secondDayOfNextMonth
}

type APITriggerOptions struct {
	DisplayName  string  `json:"budgetDisplayName"`
	Threshold    float64 `json:"alertThresholdExceeded"`
	CostAmount   float64 `json:"costAmount"`
	BudgetAmount float64 `json:"budgetAmount"`
	Currency     string  `json:"currencyCode"`
}

// safeStringSlice converts interface{} to []string safely
// moved to util.SafeStringSlice

// StoryDatabase represents a Firestore database service for stories
type StoryDatabase struct {
	client        *firestore.Client
	CollectionV2  string
	MdCollection1 string
	MdCollection2 string
	MdCollection3 string
	userProfiles  string
	tcDocuments   string
	storyFeedback string
	apiTrigger    string
	appHelper     *AppHelper
	// configs           *configs.ServiceAccount
}

// AppHelper represents the helper utility for document ID generation
type AppHelper struct {
	// Add any helper methods you need
}

// NewStoryDatabase creates a new story database service
func NewStoryDatabase() *StoryDatabase {
	return &StoryDatabase{
		CollectionV2:  "riostories_v2",
		MdCollection1: "riostories_topics_metadata_1",
		MdCollection2: "riostories_topics_metadata_2",
		MdCollection3: "riostories_topics_metadata_3",
		userProfiles:  "user_profiles",
		tcDocuments:   "tc_documents",
		storyFeedback: "story_feedback",
		apiTrigger:    "api_trigger",
		appHelper:     &AppHelper{},
	}
}

// Init initializes the Firestore connection
func (s *StoryDatabase) Init(ctx context.Context) error {
	log.Println("Initializing StoryDatabase")

	var client *firestore.Client
	var err error

	// Try to use service account file first
	// credPath := filepath.Join(filepath.Dir(".././"), "configs", "serviceAccount.json")
	credPath := "serviceAccount.json"
	_, err = os.Stat(credPath)
	if err == nil {
		log.Println("Using service account from file")
		client, err = firestore.NewClient(ctx, "riokutty", option.WithCredentialsFile(credPath))
	} else {
		log.Println("Using default credentials")
		// In Cloud Run, use the default service account
		client, err = firestore.NewClient(ctx, "riokutty")
	}

	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %v", err)
	}

	s.client = client

	// Test connection
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	iter := s.client.Collection(s.CollectionV2).Limit(1).Documents(ctx)
	_, err = iter.Next()
	if err != nil && err != iterator.Done {
		return fmt.Errorf("failed to test Firestore connection: %v", err)
	}

	log.Println("Successfully connected to Firestore")
	return nil
}

// Close closes the Firestore client
func (s *StoryDatabase) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// GetClient returns the Firestore client (for external packages that need direct access)
func (s *StoryDatabase) GetClient() *firestore.Client {
	return s.client
}

// Create Trigger Document on a API
// CreateAPITrigger creates an API trigger with optional parameters
// Usage examples:
//
//	CreateAPITrigger(ctx, "gemini")  // Basic usage
//	CreateAPITrigger(ctx, "gemini", 100.0, 50.0)  // With budget and cost
//	CreateAPITrigger(ctx, "gemini", 100.0, 50.0, "production")  // With all parameters
func (s *StoryDatabase) CreateAPITrigger(ctx context.Context, api_model string, optionalParams ...interface{}) (string, error) {
	log.Printf("Creating API Trigger for %s", api_model)
	ud, err := s.client.Collection(s.apiTrigger).Doc(api_model).Get(ctx)
	var userData map[string]interface{}
	if err != nil {
		userData = map[string]interface{}{
			"suspend":    true,
			"created_at": getUTCTimestamp(),
			"updated_at": getUTCTimestamp(),
			"reset_at":   getNextMonthResetTime(),
			"api_model":  api_model,
		}
	} else {
		userData = ud.Data()
	}

	// Process optional parameters
	// Parameter order: budgetAmount, costAmount, tag
	for i, param := range optionalParams {
		switch i {
		case 0: // budgetAmount
			if budgetAmount, ok := param.(float64); ok {
				userData["budgetAmount"] = budgetAmount
			}
		case 1: // costAmount
			if costAmount, ok := param.(float64); ok {
				userData["costAmount"] = costAmount
			}
		case 2: // tag
			if api_model == "audio" {
				if tag, ok := param.(string); ok {
					userData["tag"] = tag
				}
			}

		}
	}

	_, err = s.client.Collection(s.apiTrigger).Doc(api_model).Set(ctx, userData)
	if err != nil {
		return "", fmt.Errorf("error creating api model: %v", err)
	}

	return "Document written successfully", nil
}

// CreateAPITriggerWithOptions creates an API trigger using a struct for optional parameters
// This is a more type-safe alternative to the variadic approach
// Usage examples:
//
//	CreateAPITriggerWithOptions(ctx, "audio", nil)  // Basic usage
//	CreateAPITriggerWithOptions(ctx, "audio", &APITriggerOptions{BudgetAmount: &budget, CostAmount: &cost})  // With budget and cost
//	CreateAPITriggerWithOptions(ctx, "audio", &APITriggerOptions{BudgetAmount: &budget, CostAmount: &cost, Tag: &tag})  // With all parameters
func (s *StoryDatabase) CreateAPIAudioTrigger(ctx context.Context, api_model string, options *APITriggerOptions) (string, error) {
	log.Printf("Creating API Trigger for %s with options", api_model)
	// Initialize userData with required fields
	userData := map[string]interface{}{
		"created_at":   getUTCTimestamp(),
		"updated_at":   getUTCTimestamp(),
		"reset_at":     getNextMonthResetTime(),
		"api_model":    api_model,
		"tag":          tts.Chirp3HD.String(),
		"budgetAmount": options.BudgetAmount,
		"costAmount":   options.CostAmount,
		"threshold":    options.Threshold,
		"displayName":  options.DisplayName,
		"currency":     options.Currency,
	}

	if options.Threshold > 0.25 {
		userData["tag"] = tts.Standard.String()
	}

	_, err := s.client.Collection(s.apiTrigger).Doc(api_model).Set(ctx, userData)
	if err != nil {
		return "", fmt.Errorf("error creating api model: %v", err)
	}

	return "Document written successfully", nil
}

// Create Trigger Document on a user
func (s *StoryDatabase) SuspendAudioAPI(ctx context.Context, api_model string) (bool, string, error) {
	log.Printf("Reading API Trigger for %s", api_model)
	data, err := s.client.Collection(s.apiTrigger).Doc(api_model).Get(ctx)
	if err != nil {
		return false, "", fmt.Errorf("error reading api model: %v", err)
	}
	curr_time := getUTCTimestamp()
	reset_at := data.Data()["reset_at"].(time.Time)

	// Handle both int64 and float64 types for costAmount
	threshold, ok := data.Data()["threshold"].(float64)
	if !ok {
		if intVal, intOk := data.Data()["threshold"].(int64); intOk {
			threshold = float64(intVal)
		} else {
			threshold = 0
		}
	}
	tag, ok := data.Data()["tag"].(string)
	if !ok {
		if intVal, intOk := data.Data()["tag"].(int64); intOk {
			tag = string(intVal)
		} else {
			tag = ""
		}
	}
	if curr_time.After(reset_at) {
		s.CreateAPITrigger(ctx, "audio", tts.Chirp3HD.String(), &APITriggerOptions{BudgetAmount: 0, CostAmount: 0, Threshold: 0, DisplayName: "", Currency: ""})
	}
	if threshold >= 0.88 {
		return true, tag, nil
	}
	return false, tag, nil
}

func (s *StoryDatabase) SuspendGeminiAPI(ctx context.Context, api_model string) (bool, error) {
	log.Printf("Reading API Trigger for %s", api_model)
	data, err := s.client.Collection(s.apiTrigger).Doc(api_model).Get(ctx)
	if err != nil {
		return false, fmt.Errorf("error reading api model: %v", err)
	}
	curr_time := getUTCTimestamp()
	reset_at := data.Data()["reset_at"].(time.Time)

	// Handle both int64 and float64 types for costAmount
	costAmount, ok := data.Data()["costAmount"].(float64)
	if !ok {
		if intVal, intOk := data.Data()["costAmount"].(int64); intOk {
			costAmount = float64(intVal)
		} else {
			costAmount = 0
		}
	}

	// Handle both int64 and float64 types for budgetAmount
	budgetAmount, ok := data.Data()["budgetAmount"].(float64)
	if !ok {
		if intVal, intOk := data.Data()["budgetAmount"].(int64); intOk {
			budgetAmount = float64(intVal)
		} else {
			budgetAmount = 0
		}
	}

	thresholdPercCal := (costAmount / budgetAmount) * 100
	//If the threshold percentage is greater than 90% and the current time is before the reset time, return true
	if thresholdPercCal >= 60 && curr_time.Before(reset_at) {
		return true, nil
	}
	return false, nil
}

// Create TC Document on a user
func (s *StoryDatabase) CreateTc(ctx context.Context, data *model.Tc) (string, error) {

	userData := map[string]interface{}{
		"accepted":   data.Accepted,
		"created_at": getUTCTimestamp(),
		"updated_at": getUTCTimestamp(),
		"email":      data.Email,
	}
	_, err := s.client.Collection(s.tcDocuments).Doc(data.Email).Set(ctx, userData)
	if err != nil {
		return "", fmt.Errorf("error creating user profile: %v", err)
	}

	return "Document written successfully", nil
}

// Create TC Document on a user
func (s *StoryDatabase) CreateStoryFeedback(ctx context.Context, data *model.StoryFeedback) (string, error) {

	userData := map[string]interface{}{
		"like":       data.Like,
		"created_at": getUTCTimestamp(),
		"storyId":    data.StoryId,
		"email":      data.Email,
	}
	_, err := s.client.Collection(s.storyFeedback).Doc(data.StoryId).Set(ctx, userData)
	if err != nil {
		return "", fmt.Errorf("error creating user profile: %v", err)
	}

	return "Document written successfully", nil
}

// CreateMDTopics1 creates metadata topics collection 1
func (s *StoryDatabase) CreateMDTopics1(ctx context.Context, theme_id, country, city string, preference string, topics []string, language string) (string, error) {
	data := map[string]interface{}{
		"theme_id":   theme_id,
		"country":    country,
		"city":       city,
		"preference": preference,
		"topics":     topics,
		"language":   language,
		"created_at": getUTCTimestamp(),
		"updated_at": getUTCTimestamp(),
	}

	docRef, _, err := s.client.Collection(s.MdCollection1).Add(ctx, data)
	if err != nil {
		return "", fmt.Errorf("error creating metadata topics 1: %v", err)
	}

	return docRef.ID, nil
}

// GetUserProfile retrieves a user profile by email
func (s *StoryDatabase) GetUserProfile(ctx context.Context, username, emailID string) (map[string]interface{}, error) {

	if s.client == nil {
		log.Printf("❌ DEBUG: Firestore client is nil!")
		return nil, fmt.Errorf("Firestore client not initialized")
	}

	doc, err := s.client.Collection("user_profiles").Doc(emailID).Get(ctx)
	if err != nil {
		log.Printf("❌ DEBUG: Error reading user profile: %v", err)
		return nil, fmt.Errorf("error reading user profile: %v", err)
	}

	data := doc.Data()
	log.Printf("✅ DEBUG: User profile data: %v", data)
	return data, nil
}

// GetUserProfileByEmail retrieves a user profile from Firestore by email.
// Returns nil if the user is not found.
func (db *StoryDatabase) GetUserProfileByEmail(ctx context.Context, email string) (map[string]interface{}, error) {
	iter := db.client.Collection(db.userProfiles).Where("email", "==", email).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil // User not found, which is not an error in this case
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user profile by email: %w", err)
	}

	return doc.Data(), nil
}

// IncrementTokenVersion increments the token version for a user (for logout)
func (db *StoryDatabase) IncrementTokenVersion(ctx context.Context, email string) error {
	iter := db.client.Collection(db.userProfiles).Where("email", "==", email).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return fmt.Errorf("user not found")
	}
	if err != nil {
		return fmt.Errorf("failed to query user profile: %w", err)
	}

	// Get current token version or default to 0
	var currentVersion int64 = 0
	if version, exists := doc.Data()["token_version"]; exists {
		if v, ok := version.(int64); ok {
			currentVersion = v
		}
	}

	// Increment token version
	_, err = doc.Ref.Update(ctx, []firestore.Update{
		{Path: "token_version", Value: currentVersion + 1},
		{Path: "last_logout", Value: time.Now()},
	})

	if err != nil {
		return fmt.Errorf("failed to update token version: %w", err)
	}

	return nil
}

// GetTokenVersion retrieves the current token version for a user
func (db *StoryDatabase) GetTokenVersion(ctx context.Context, email string) (int64, error) {
	iter := db.client.Collection(db.userProfiles).Where("email", "==", email).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return 0, fmt.Errorf("user not found")
	}
	if err != nil {
		return 0, fmt.Errorf("failed to query user profile: %w", err)
	}

	// Get current token version or default to 0
	var currentVersion int64 = 0
	if version, exists := doc.Data()["token_version"]; exists {
		if v, ok := version.(int64); ok {
			currentVersion = v
		}
	}

	return currentVersion, nil
}

// CreateUserProfile creates a new user profile
func (s *StoryDatabase) CreateUserProfile(ctx context.Context, userData map[string]interface{}) (string, error) {
	email := userData["email"].(string)

	// Add timestamps
	userData["created_at"] = getUTCTimestamp()
	userData["updated_at"] = getUTCTimestamp()

	_, err := s.client.Collection(s.userProfiles).Doc(email).Set(ctx, userData)
	if err != nil {
		return "", fmt.Errorf("error creating user profile: %v", err)
	}

	return email, nil
}

func (s *StoryDatabase) DeleteUserProfile(ctx context.Context, email string) error {
	_, err := s.client.Collection(s.userProfiles).Doc(email).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error deleting user profile: %v", err)
	}
	return nil
}

func (s *StoryDatabase) UpdateUserProfileByEmail(ctx context.Context, email string, data model.UserProfile) error {
	userProfile, err := s.GetUserProfileByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("error getting user profile: %v", err)
	}

	userProfile["country"] = data.Country
	userProfile["city"] = data.City
	userProfile["religions"] = data.Religions
	userProfile["preferences"] = data.Preferences
	userProfile["updated_at"] = getUTCTimestamp()

	_, err = s.client.Collection(s.userProfiles).Doc(email).Set(ctx, userProfile)
	if err != nil {
		return fmt.Errorf("error updating user profile: %v", err)
	}
	return nil
}

// UpdateUserProfile updates an existing user profile
func (s *StoryDatabase) UpdateUserProfile(ctx context.Context, email string, data map[string]interface{}) error {
	// Convert map to firestore updates
	var updates []firestore.Update
	for key, value := range data {
		updates = append(updates, firestore.Update{Path: key, Value: value})
	}

	_, err := s.client.Collection(s.userProfiles).Doc(email).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("error updating user profile: %v", err)
	}
	return nil
}

func (s *StoryDatabase) InitialReadMDTopics1(ctx context.Context) ([]map[string]interface{}, error) {
	log.Printf("Reading metadata for initial context ")
	userProfile, err := s.GetUserProfileByEmail(ctx, "rio.oly.pluto@gmail.com")
	if err != nil {
		return nil, fmt.Errorf("error getting user profile: %v", err)
	}
	log.Printf("User profile: %v", userProfile)
	language := userProfile["language"].(string)
	preferences := util.SafeStringSlice(userProfile["preferences"])

	var allDocs []*firestore.DocumentSnapshot
	for _, preference := range preferences {
		query := s.client.Collection(s.MdCollection1).
			Where("language", "==", language).
			Where("preference", "==", preference)

		iterationDocs, err := query.Documents(ctx).GetAll()
		if err != nil {
			return nil, fmt.Errorf("error executing query: %v", err)
		}

		if len(iterationDocs) == 0 {
			log.Printf("No metadata topics found for preference: %s", preference)
			continue
		}
		allDocs = append(allDocs, iterationDocs...)
	}

	if len(allDocs) == 0 {
		log.Println("No metadata topics found for any preference")
		return nil, nil
	}

	var results []map[string]interface{}
	for _, doc := range allDocs {
		results = append(results, doc.Data())
	}

	if len(results) > 0 {
		log.Printf("Found metadata topics: %v", results[0]["theme_id"])
	}
	return results, nil
}

// ReadMDTopics1 reads metadata topics collection 1
func (s *StoryDatabase) ReadMDTopics1(ctx context.Context, country, city string, preferences []string, language string) ([]map[string]interface{}, error) {
	log.Printf("Reading metadata topics 1 for country: %s, city: %s, preferences: %v", country, city, preferences)

	var allDocs []*firestore.DocumentSnapshot
	for _, preference := range preferences {
		var query firestore.Query
		var iterationDocs []*firestore.DocumentSnapshot
		var err error

		if country != "Any" {
			query = s.client.Collection(s.MdCollection1).
				Where("language", "==", language).
				Where("country", "==", country).
				Where("preference", "==", preference)
			iterationDocs, err = query.Documents(ctx).GetAll()
			if err != nil {
				return nil, fmt.Errorf("error executing query: %v", err)
			}
		}
		if country == "Any" || len(iterationDocs) == 0 {
			query = s.client.Collection(s.MdCollection1).
				Where("language", "==", language).
				Where("preference", "==", preference)
			iterationDocs, err = query.Documents(ctx).GetAll()
			if err != nil {
				return nil, fmt.Errorf("error executing query: %v", err)
			}
		}

		if len(iterationDocs) == 0 {
			log.Printf("No metadata topics found for preference: %s", preference)
			continue
		}
		allDocs = append(allDocs, iterationDocs...)
	}

	if len(allDocs) == 0 {
		log.Println("No metadata topics found for any preference... ")
		return nil, nil
	}

	var results []map[string]interface{}
	for _, doc := range allDocs {
		results = append(results, doc.Data())
	}

	if len(results) < 0 {
		return s.InitialReadMDTopics1(ctx)
	}
	return results, nil
}

// CreateMDTopics2 creates metadata topics collection 2
func (s *StoryDatabase) CreateMDTopics2(ctx context.Context, theme_id, country string, religion string, language string, preferences, topics []string) (string, error) {
	data := map[string]interface{}{
		"theme_id":    theme_id,
		"country":     country,
		"religion":    religion,
		"preferences": preferences,
		"topics":      topics,
		"language":    language,
		"created_at":  getUTCTimestamp(),
		"updated_at":  getUTCTimestamp(),
	}

	docRef, _, err := s.client.Collection(s.MdCollection2).Add(ctx, data)
	if err != nil {
		return "", fmt.Errorf("error creating metadata topics 2: %v", err)
	}

	return docRef.ID, nil
}

func (s *StoryDatabase) InitialReadMDTopics2(ctx context.Context) ([]map[string]interface{}, error) {
	log.Printf("Reading metadata for initial context ")
	userProfile, err := s.GetUserProfileByEmail(ctx, "rio.oly.pluto@gmail.com")
	if err != nil {
		return nil, fmt.Errorf("error getting user profile: %v", err)
	}
	log.Printf("User profile: %v", userProfile)
	religions := util.SafeStringSlice(userProfile["religions"])
	preferences := util.SafeStringSlice(userProfile["preferences"])
	language := userProfile["language"].(string)

	var allDocs []*firestore.DocumentSnapshot
	for _, religion := range religions {
		query := s.client.Collection(s.MdCollection2).
			Where("religion", "==", religion).
			Where("language", "==", language).
			Where("preferences", "array-contains-any", preferences)

		iterationDocs, err := query.Documents(ctx).GetAll()
		if err != nil {
			return nil, fmt.Errorf("error executing query: %v", err)
		}

		if len(iterationDocs) == 0 {
			log.Printf("No metadata topics found for religion: %s", religion)
			continue
		}
		allDocs = append(allDocs, iterationDocs...)
	}

	if len(allDocs) == 0 {
		log.Println("No metadata topics found for any preference")
		return nil, nil
	}

	var results []map[string]interface{}
	for _, doc := range allDocs {
		results = append(results, doc.Data())
	}

	if len(results) > 0 {
		log.Printf("Found metadata topics: %v", results[0]["theme_id"])
	}
	return results, nil
}

// ReadMDTopics2 reads metadata topics collection 2
func (s *StoryDatabase) ReadMDTopics2(ctx context.Context, country string, religions, preferences []string, language string) ([]map[string]interface{}, error) {
	// First filter by country and religions using array-contains-any
	var allDocs []*firestore.DocumentSnapshot
	for _, religion := range religions {
		var iterationDocs []*firestore.DocumentSnapshot
		var err error
		if country != "Any" {
			query := s.client.Collection(s.MdCollection2).
				Where("language", "==", language).
				Where("country", "==", country).
				Where("religion", "==", religion).
				Where("preferences", "array-contains-any", preferences)
			iterationDocs, err = query.Documents(ctx).GetAll()
			if err != nil {
				return nil, fmt.Errorf("error executing query: %v", err)
			}
		}
		if religion == "Any" || len(iterationDocs) == 0 {
			query := s.client.Collection(s.MdCollection2).
				Where("language", "==", language).
				Where("religion", "==", religion).
				Where("preferences", "array-contains-any", preferences)
			iterationDocs, err = query.Documents(ctx).GetAll()
			if err != nil {
				return nil, fmt.Errorf("error executing query: %v", err)
			}
		}
		if len(iterationDocs) == 0 {
			log.Printf("No metadata topics found for religion: %s", religion)
			continue
		}
		allDocs = append(allDocs, iterationDocs...)
	}

	if len(allDocs) == 0 {
		log.Println("No metadata topics found for any religion")
		return nil, nil
	}

	var results []map[string]interface{}
	for _, doc := range allDocs {
		results = append(results, doc.Data())
	}

	if len(results) > 0 {
		log.Printf("Found metadata topics: %v", results[0]["theme_id"])
	}
	return results, nil
}

// CreateMDTopics3 creates metadata topics collection 3
func (s *StoryDatabase) CreateMDTopics3(ctx context.Context, theme_id, preference string, language string, topics []string) (string, error) {
	data := map[string]interface{}{
		"theme_id":   theme_id,
		"preference": preference,
		"topics":     topics,
		"language":   language,
		"created_at": getUTCTimestamp(),
		"updated_at": getUTCTimestamp(),
	}

	docRef, _, err := s.client.Collection(s.MdCollection3).Add(ctx, data)
	if err != nil {
		return "", fmt.Errorf("error creating metadata topics 3: %v", err)
	}

	return docRef.ID, nil
}

func (s *StoryDatabase) InitialReadMDTopics3(ctx context.Context) ([]map[string]interface{}, error) {
	log.Printf("Reading metadata for initial context ")
	userProfile, err := s.GetUserProfileByEmail(ctx, "rio.oly.pluto@gmail.com")
	if err != nil {
		return nil, fmt.Errorf("error getting user profile: %v", err)
	}
	log.Printf("User profile: %v", userProfile)
	preferences := util.SafeStringSlice(userProfile["preferences"])

	var allDocs []*firestore.DocumentSnapshot
	for _, preference := range preferences {
		query := s.client.Collection(s.MdCollection3).
			Where("preference", "==", preference)

		iterationDocs, err := query.Documents(ctx).GetAll()
		if err != nil {
			return nil, fmt.Errorf("error executing query: %v", err)
		}

		if len(iterationDocs) == 0 {
			log.Printf("No metadata topics found for preference: %s", preference)
			continue
		}
		allDocs = append(allDocs, iterationDocs...)
	}

	if len(allDocs) == 0 {
		log.Println("No metadata topics found for any preference")
		return nil, nil
	}

	var results []map[string]interface{}
	for _, doc := range allDocs {
		results = append(results, doc.Data())
	}

	if len(results) > 0 {
		log.Printf("Found metadata topics: %v", results[0])
	}
	return results, nil
}

// ReadMDTopics3 reads metadata topics collection 3
func (s *StoryDatabase) ReadMDTopics3(ctx context.Context, preferences []string, language string) ([]map[string]interface{}, error) {
	var allDocs []*firestore.DocumentSnapshot
	for _, preference := range preferences {
		query := s.client.Collection(s.MdCollection3).
			Where("language", "==", language).
			Where("preference", "==", preference)

		iterationDocs, err := query.Documents(ctx).GetAll()
		if err != nil {
			return nil, fmt.Errorf("error executing query: %v", err)
		}
		if len(iterationDocs) == 0 {
			log.Printf("No metadata topics found for preference: %s", preference)
			continue
		}
		allDocs = append(allDocs, iterationDocs...)
	}

	if len(allDocs) == 0 {
		log.Println("No metadata topics found for any preference")
		return nil, nil
	}

	var results []map[string]interface{}
	for _, doc := range allDocs {
		results = append(results, doc.Data())
	}

	if len(results) > 0 {
		log.Printf("Found metadata topics: %v", results[0])
	}
	return results, nil
}

// CreateStoryV2 creates a new story with custom document ID
func (s *StoryDatabase) CreateStoryV2(ctx context.Context, theme_id string, storyData map[string]interface{}) (string, error) {
	title := storyData["title"].(string)
	theme := storyData["theme"].(string)

	docID := s.appHelper.GetDocID(title, theme)

	// Add timestamps
	storyData["created_at"] = getUTCTimestamp()
	storyData["updated_at"] = getUTCTimestamp()
	storyData["theme_id"] = theme_id

	_, err := s.client.Collection(s.CollectionV2).Doc(docID).Set(ctx, storyData)
	if err != nil {
		return "", fmt.Errorf("error creating story v2: %v", err)
	}

	return docID, nil
}

// GetStory retrieves a story by Theme ID
func (s *StoryDatabase) GetStoryByThemeID(ctx context.Context, themeID string) ([]map[string]interface{}, error) {
	query := s.client.Collection(s.CollectionV2).Where("theme_id", "==", themeID)
	log.Printf("Getting story by theme id: %s", themeID)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting story: %v", err)
	}

	result := make([]map[string]interface{}, 0)
	for _, doc := range docs {
		data := doc.Data()
		data["id"] = doc.Ref.ID
		result = append(result, data)
	}
	log.Printf("Found %d stories by theme id: %s", len(result), themeID)
	return result, nil
}

// GetStoryV2 retrieves a story v2 by ID
func (s *StoryDatabase) GetStoryV2(ctx context.Context, storyID string) (map[string]interface{}, error) {
	doc, err := s.client.Collection(s.CollectionV2).Doc(storyID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting story v2: %v", err)
	}

	if !doc.Exists() {
		return nil, nil
	}

	return doc.Data(), nil
}

// ListStoriesV2 lists stories v2 with filtering
func (s *StoryDatabase) ListStoriesV2(ctx context.Context, limit int, theme, title string) (map[string]interface{}, error) {
	log.Printf("Listing stories v2 with limit: %d, theme: %s, title: %s", limit, theme, title)

	docID := s.appHelper.GetDocID(title, theme)
	log.Printf("Generated doc_id: %s", docID)

	doc, err := s.client.Collection(s.CollectionV2).Doc(docID).Get(ctx)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	if !doc.Exists() {
		return nil, nil
	}

	data := doc.Data()
	log.Printf("Found story data: %v", data)
	return data, nil
}

// ListStoriesByThemeID lists stories v2 by theme_id
func (s *StoryDatabase) ListStoriesByThemeID(ctx context.Context, themeID string, limit int) ([]map[string]interface{}, error) {
	log.Printf("Listing stories v2 by theme_id: %s with limit: %d", themeID, limit)

	query := s.client.Collection(s.CollectionV2).Where("theme_id", "==", themeID)

	if limit > 0 {
		query = query.Limit(limit)
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Printf("Error listing stories by theme_id: %v", err)
		return nil, fmt.Errorf("error listing stories by theme_id: %v", err)
	}

	result := make([]map[string]interface{}, 0)
	for _, doc := range docs {
		data := doc.Data()
		data["id"] = doc.Ref.ID
		result = append(result, data)
	}

	log.Printf("Found %d stories with theme_id: %s", len(result), themeID)
	return result, nil
}

// UpdateStory updates an existing story
func (s *StoryDatabase) UpdateStory(ctx context.Context, storyID string, storyData map[string]interface{}) error {
	// Convert map to firestore updates
	var updates []firestore.Update
	for key, value := range storyData {
		updates = append(updates, firestore.Update{Path: key, Value: value})
	}
	updates = append(updates, firestore.Update{Path: "updated_at", Value: getUTCTimestamp()})

	_, err := s.client.Collection(s.CollectionV2).Doc(storyID).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("error updating story: %v", err)
	}

	return nil
}

// DeleteStory deletes a story
func (s *StoryDatabase) DeleteStory(ctx context.Context, storyID string) error {
	_, err := s.client.Collection(s.CollectionV2).Doc(storyID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error deleting story: %v", err)
	}

	return nil
}

// HealthCheck checks if Firestore is accessible
func (s *StoryDatabase) HealthCheck(ctx context.Context) error {
	if s.client == nil {
		return fmt.Errorf("Firestore client not initialized")
	}

	// Try to access a collection to check connectivity
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.client.Collection("_health").Doc("check").Get(ctx)
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		return fmt.Errorf("Firestore health check failed: %v", err)
	}

	return nil
}

// AppHelper methods

// GetDocID generates a document ID based on title and theme
func (a *AppHelper) GetDocID(title, theme string) string {
	// Simple implementation - you can enhance this as needed
	return fmt.Sprintf("%s_%s", title, theme)
}

func (s *StoryDatabase) UpdateAPITokens(ctx context.Context, api_model string, tokensToAdd int64) (string, error) {
	log.Printf("Updating API Trigger for %s is %d", api_model, tokensToAdd)
	ud, err := s.client.Collection(s.apiTrigger).Doc(api_model).Get(ctx)
	var userData map[string]interface{}
	if err != nil {
		log.Printf("Error updating API Trigger for %s: %v", api_model, err)
		userData = map[string]interface{}{
			"created_at":   getUTCTimestamp(),
			"updated_at":   getUTCTimestamp(),
			"reset_at":     getNextMonthResetTime(),
			"api_model":    api_model,
			"budgetAmount": 0,
			"costAmount":   0,
			"threshold":    0,
			"displayName":  "",
			"currency":     "",
			"tag":          "",
		}
	} else {
		log.Printf("Else: Updating API Trigger for %s is %d", api_model, ud.Data()["tokensUsed"])
		userData = ud.Data()
	}
	log.Printf(" token used: %v", userData["tokensUsed"])
	tokensUsed := userData["tokensUsed"].(int64)
	var allTokensUsed int64

	if tokensUsed > 0 {
		log.Printf("All tokens used for %d is %d", tokensUsed, allTokensUsed)
		allTokensUsed = tokensUsed + tokensToAdd
	} else {
		allTokensUsed = tokensToAdd
		log.Printf("Else: All tokens used for %d is %d", tokensUsed, allTokensUsed)
	}
	log.Printf("Updating API Trigger for %s is %d", api_model, allTokensUsed)
	userData["tokensUsed"] = allTokensUsed
	userData["updated_at"] = getUTCTimestamp()
	_, err = s.client.Collection(s.apiTrigger).Doc(api_model).Set(ctx, userData)
	if err != nil {
		return "", fmt.Errorf("error creating api model: %v", err)
	}
	log.Printf("Document Updated successfully for %s", api_model)
	return "Document Updated successfully", nil
}

// func (s *StoryDatabase) MigrateAllThemeId(ctx context.Context) ([]map[string]interface{}, error) {
// 	query := s.client.Collection(s.mdCollection3).Documents(ctx)
// 	for {
// 		theme1_id := uuid.New().String()
// 		docSnapshot, err := query.Next()
// 		if err != nil {
// 			if err.Error() == "iterator.Done" {
// 				break
// 			}
// 			return nil, fmt.Errorf("error iterating documents: %v", err)
// 		}
// 		oldTheme1_id := docSnapshot.Data()["theme_id"].(string)
// 		log.Printf("Migrating theme 1 id: %s, old theme 1 id: %s", theme1_id, oldTheme1_id)

// 		storyData := docSnapshot.Data()
// 		storyData["theme_id"] = theme1_id
// 		_, err = s.client.Collection(s.mdCollection3).Doc(docSnapshot.Ref.ID).Set(ctx, storyData)
// 		if err != nil {
// 			log.Printf("Error migrating story %s: %v", docSnapshot.Ref.ID, err)
// 			continue
// 		}

// 		log.Printf("Migrating story: %s", docSnapshot.Ref.ID)

// 		// Handle both []string and []interface{} types from Firestore
// 		var topics []string
// 		if topicsInterface, ok := storyData["topics"].([]string); ok {
// 			topics = topicsInterface
// 		} else if topicsArray, ok := storyData["topics"].([]interface{}); ok {
// 			for _, topic := range topicsArray {
// 				if topicStr, ok := topic.(string); ok {
// 					topics = append(topics, topicStr)
// 				}
// 			}
// 		} else {
// 			log.Printf("Warning: topics field is not []string or []interface{} for story %s", docSnapshot.Ref.ID)
// 			continue
// 		}

// 		for _, topic := range topics {
// 			log.Printf("Migrating topic: %s", topic)
// 			docID := topic + "_3"
// 			// docID := s.appHelper.GetDocID(topic, storyData["theme"].(string))
// 			log.Printf("Migrating topic: %s, docID: %s", topic, docID)

// 			// Check if topic document exists
// 			topicDocSnapshot, err := s.client.Collection(s.CollectionV2).Doc(docID).Get(ctx)
// 			if err == nil && topicDocSnapshot.Exists() {
// 				log.Printf("Topic document exists: %s", docID)
// 				// Document exists, update theme_id
// 				topicData := topicDocSnapshot.Data()
// 				topicData["theme_id"] = theme1_id
// 				_, err = s.client.Collection(s.CollectionV2).Doc(docID).Set(ctx, topicData)
// 				if err != nil {
// 					log.Printf("Error updating theme_id for topic %s: %v", docID, err)
// 				} else {
// 					log.Printf("Updated theme_id for existing topic: %s", docID)
// 				}
// 			}
// 		}
// 		log.Printf("Migrating story: %s completed", docSnapshot.Ref.ID)
// 	}
// 	return nil, nil
// }
