package database

import (
	"context"
	"fmt"
	"log"
	"os"
	// "path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	// "rio-go-model/configs"
)

// StoryDatabase represents a Firestore database service for stories
type StoryDatabase struct {
	client            *firestore.Client
	collection        string
	collectionV2      string
	mdCollection1     string
	mdCollection2     string
	mdCollection3     string
	userProfiles      string
	appHelper         *AppHelper
	// configs           *configs.ServiceAccount
}

// AppHelper represents the helper utility for document ID generation
type AppHelper struct {
	// Add any helper methods you need
}

// NewStoryDatabase creates a new story database service
func NewStoryDatabase() *StoryDatabase {
	return &StoryDatabase{
		collection:    "riostories",
		collectionV2:  "riostories_v2",
		mdCollection1: "riostories_topics_metadata_1",
		mdCollection2: "riostories_topics_metadata_2",
		mdCollection3: "riostories_topics_metadata_3",
		userProfiles:  "user_profiles",
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

	iter := s.client.Collection(s.collection).Limit(1).Documents(ctx)
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

// CreateMDTopics1 creates metadata topics collection 1
func (s *StoryDatabase) CreateMDTopics1(ctx context.Context, country, city string, preferences, topics []string) (string, error) {
	data := map[string]interface{}{
		"country":     country,
		"city":        city,
		"preferences": preferences,
		"topics":      topics,
		"created_at":  firestore.ServerTimestamp,
		"updated_at":  firestore.ServerTimestamp,
	}

	docRef, _, err := s.client.Collection(s.mdCollection1).Add(ctx, data)
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

	doc, err := s.client.Collection(s.userProfiles).Doc(emailID).Get(ctx)
	if err != nil {
		log.Printf("❌ DEBUG: Error reading user profile: %v", err)
		return nil, fmt.Errorf("error reading user profile: %v", err)
	}

	data := doc.Data()
	log.Printf("✅ DEBUG: User profile data: %v", data)
	return data, nil
}

// CreateUserProfile creates a new user profile
func (s *StoryDatabase) CreateUserProfile(ctx context.Context, userData map[string]interface{}) (string, error) {
	email := userData["email"].(string)
	
	// Add timestamps
	userData["created_at"] = firestore.ServerTimestamp
	userData["updated_at"] = firestore.ServerTimestamp

	_, err := s.client.Collection(s.userProfiles).Doc(email).Set(ctx, userData)
	if err != nil {
		return "", fmt.Errorf("error creating user profile: %v", err)
	}

	return email, nil
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

// ReadMDTopics1 reads metadata topics collection 1
func (s *StoryDatabase) ReadMDTopics1(ctx context.Context, country, city string, preferences []string) (*firestore.DocumentSnapshot, error) {
	log.Printf("Reading metadata topics 1 for country: %s, city: %s, preferences: %v", country, city, preferences)

	query := s.client.Collection(s.mdCollection1).
		Where("country", "==", country).
		Where("city", "==", city).
		Where("preferences", "array-contains-any", preferences)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	if len(docs) == 0 {
		log.Println("No metadata topics found")
		return nil, nil
	}

	log.Printf("Found metadata topics: %v", docs[0].Data())
	return docs[0], nil
}

// CreateMDTopics2 creates metadata topics collection 2
func (s *StoryDatabase) CreateMDTopics2(ctx context.Context, country string, religions, preferences, topics []string) (string, error) {
	data := map[string]interface{}{
		"country":     country,
		"religions":   religions,
		"preferences": preferences,
		"topics":      topics,
		"created_at":  firestore.ServerTimestamp,
		"updated_at":  firestore.ServerTimestamp,
	}

	docRef, _, err := s.client.Collection(s.mdCollection2).Add(ctx, data)
	if err != nil {
		return "", fmt.Errorf("error creating metadata topics 2: %v", err)
	}

	return docRef.ID, nil
}

// ReadMDTopics2 reads metadata topics collection 2
func (s *StoryDatabase) ReadMDTopics2(ctx context.Context, country string, religions, preferences []string) (*firestore.DocumentSnapshot, error) {
	// First filter by country and religions using array-contains-any
	initialQuery := s.client.Collection(s.mdCollection2).
		Where("country", "==", country).
		Where("religions", "array-contains-any", religions)

	docs, err := initialQuery.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error executing initial query: %v", err)
	}

	// Then filter results in memory for preferences match
	for _, doc := range docs {
		docPreferences := doc.Data()["preferences"].([]interface{})
		// Check if any preference matches
		for _, pref := range preferences {
			for _, docPref := range docPreferences {
				if pref == docPref {
					return doc, nil
				}
			}
		}
	}

	return nil, nil
}

// CreateMDTopics3 creates metadata topics collection 3
func (s *StoryDatabase) CreateMDTopics3(ctx context.Context, preferences, topics []string) (string, error) {
	data := map[string]interface{}{
		"preferences": preferences,
		"topics":      topics,
		"created_at":  firestore.ServerTimestamp,
		"updated_at":  firestore.ServerTimestamp,
	}

	docRef, _, err := s.client.Collection(s.mdCollection3).Add(ctx, data)
	if err != nil {
		return "", fmt.Errorf("error creating metadata topics 3: %v", err)
	}

	return docRef.ID, nil
}

// ReadMDTopics3 reads metadata topics collection 3
func (s *StoryDatabase) ReadMDTopics3(ctx context.Context, preferences []string) (*firestore.DocumentSnapshot, error) {
	query := s.client.Collection(s.mdCollection3).
		Where("preferences", "array-contains-any", preferences)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	if len(docs) > 0 {
		return docs[0], nil
	}

	return nil, nil
}

// CreateStory creates a new story
func (s *StoryDatabase) CreateStory(ctx context.Context, storyData map[string]interface{}) (string, error) {
	// Add timestamps
	storyData["created_at"] = firestore.ServerTimestamp
	storyData["updated_at"] = firestore.ServerTimestamp

	docRef, _, err := s.client.Collection(s.collection).Add(ctx, storyData)
	if err != nil {
		return "", fmt.Errorf("error creating story: %v", err)
	}

	return docRef.ID, nil
}

// CreateStoryV2 creates a new story with custom document ID
func (s *StoryDatabase) CreateStoryV2(ctx context.Context, storyData map[string]interface{}) (string, error) {
	title := storyData["title"].(string)
	theme := storyData["theme"].(string)
	
	docID := s.appHelper.GetDocID(title, theme)
	
	// Add timestamps
	storyData["created_at"] = firestore.ServerTimestamp
	storyData["updated_at"] = firestore.ServerTimestamp

	_, err := s.client.Collection(s.collectionV2).Doc(docID).Set(ctx, storyData)
	if err != nil {
		return "", fmt.Errorf("error creating story v2: %v", err)
	}

	return docID, nil
}

// GetStory retrieves a story by ID
func (s *StoryDatabase) GetStory(ctx context.Context, storyID string) (map[string]interface{}, error) {
	doc, err := s.client.Collection(s.collection).Doc(storyID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting story: %v", err)
	}

	if !doc.Exists() {
		return nil, nil
	}

	return doc.Data(), nil
}

// GetStoryV2 retrieves a story v2 by ID
func (s *StoryDatabase) GetStoryV2(ctx context.Context, storyID string) (map[string]interface{}, error) {
	doc, err := s.client.Collection(s.collectionV2).Doc(storyID).Get(ctx)
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

	doc, err := s.client.Collection(s.collectionV2).Doc(docID).Get(ctx)
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

// ListStories lists stories with filtering
func (s *StoryDatabase) ListStories(ctx context.Context, limit int, theme string) ([]map[string]interface{}, error) {
	log.Printf("Listing stories with limit: %d, theme: %s", limit, theme)
	
	var docs []*firestore.DocumentSnapshot
	var err error
	
	if theme != "" {
		docs, err = s.client.Collection(s.collection).Where("theme", "==", theme).Limit(limit).Documents(ctx).GetAll()
	} else {
		docs, err = s.client.Collection(s.collection).Limit(limit).Documents(ctx).GetAll()
	}
	
	if err != nil {
		log.Printf("Error listing stories: %v", err)
		return nil, fmt.Errorf("error listing stories: %v", err)
	}

	result := make([]map[string]interface{}, 0)
	for _, doc := range docs {
		data := doc.Data()
		data["id"] = doc.Ref.ID
		result = append(result, data)
	}

	log.Printf("Found %d stories", len(result))
	return result, nil
}

// UpdateStory updates an existing story
func (s *StoryDatabase) UpdateStory(ctx context.Context, storyID string, storyData map[string]interface{}) error {
	// Convert map to firestore updates
	var updates []firestore.Update
	for key, value := range storyData {
		updates = append(updates, firestore.Update{Path: key, Value: value})
	}
	updates = append(updates, firestore.Update{Path: "updated_at", Value: firestore.ServerTimestamp})
	
	_, err := s.client.Collection(s.collection).Doc(storyID).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("error updating story: %v", err)
	}
	
	return nil
}

// DeleteStory deletes a story
func (s *StoryDatabase) DeleteStory(ctx context.Context, storyID string) error {
	_, err := s.client.Collection(s.collection).Doc(storyID).Delete(ctx)
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