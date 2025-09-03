// @title Story API
// @version 1.0
// @description A comprehensive API for generating and managing stories with AI
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rio-go-model/docs"
	"rio-go-model/internal/handlers"
	"rio-go-model/internal/services/database"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load environment variables: %v", err)
	}



	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "Story API"
	docs.SwaggerInfo.Description = "A comprehensive API for generating and managing stories with AI"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Initialize services
	ctx := context.Background()
	
	// Initialize database service
	log.Println("🔧 Initializing database service...")
	storyDB := database.NewStoryDatabase()
	if err := storyDB.Init(ctx); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	log.Println("✅ Database service initialized successfully")
	
	// Initialize storage service
	log.Println("🔧 Initializing storage service...")
	storageService := database.NewStorageService("kutty_bucket")
	if err := storageService.Init(ctx); err != nil {
		log.Fatalf("❌ Failed to initialize storage service: %v", err)
	}
	log.Println("✅ Storage service initialized successfully")
	
	// Create router
	r := mux.NewRouter()

	// Create story topics handler with proper dependency injection
	storyTopicsHandler := handlers.NewStory(storyDB, storageService)

	// Documentation endpoint
	r.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Story API Documentation</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
				.container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
				.endpoint { background: #f8f9fa; padding: 20px; margin: 20px 0; border-radius: 8px; border-left: 4px solid #007bff; }
				.method { color: white; padding: 5px 12px; border-radius: 4px; font-weight: bold; font-size: 12px; }
				.get { background: #28a745; }
				.post { background: #007bff; }
				.put { background: #ffc107; color: #212529; }
				.delete { background: #dc3545; }
				h1 { color: #333; border-bottom: 2px solid #007bff; padding-bottom: 10px; }
				.url { font-family: monospace; background: #e9ecef; padding: 5px 10px; border-radius: 4px; }
				.auth { background: #ffc107; color: #212529; padding: 5px 10px; border-radius: 4px; font-size: 12px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>📚 Story API Documentation</h1>
				
				<div class="endpoint">
					<h3><span class="method get">GET</span> <span class="url">/api/v1/story-topics</span> <span class="auth">🔒 Auth Required</span></h3>
					<p><strong>Description:</strong> Get all story topics</p>
					<p><strong>Headers:</strong> Authorization: Bearer &lt;token&gt;</p>
					<p><strong>Response:</strong> List of all story topics</p>
				</div>
				
				<div class="endpoint">
					<h3><span class="method post">POST</span> <span class="url">/api/v1/story-topics</span> <span class="auth">🔒 Auth Required</span></h3>
					<p><strong>Description:</strong> Create a new story topic</p>
					<p><strong>Headers:</strong> Authorization: Bearer &lt;token&gt;</p>
					<p><strong>Body:</strong> JSON with country, city, religions, preferences</p>
				</div>
				
				<div class="endpoint">
					<h3><span class="method get">GET</span> <span class="url">/api/v1/story-topics/{id}</span> <span class="auth">🔒 Auth Required</span></h3>
					<p><strong>Description:</strong> Get a specific story topic by ID</p>
					<p><strong>Headers:</strong> Authorization: Bearer &lt;token&gt;</p>
					<p><strong>Parameters:</strong> id in URL path</p>
				</div>
				
				<div class="endpoint">
					<h3><span class="method put">PUT</span> <span class="url">/api/v1/story-topics/{id}</span> <span class="auth">🔒 Auth Required</span></h3>
					<p><strong>Description:</strong> Update a story topic</p>
					<p><strong>Headers:</strong> Authorization: Bearer &lt;token&gt;</p>
					<p><strong>Body:</strong> JSON with fields to update</p>
				</div>
				
				<div class="endpoint">
					<h3><span class="method delete">DELETE</span> <span class="url">/api/v1/story-topics/{id}</span> <span class="auth">🔒 Auth Required</span></h3>
					<p><strong>Description:</strong> Delete a story topic</p>
					<p><strong>Headers:</strong> Authorization: Bearer &lt;token&gt;</p>
				</div>
				
				<hr style="margin: 30px 0;">
				<p><strong>Base URL:</strong> <span class="url">http://localhost:8080</span></p>
				<p><strong>Authentication:</strong> All endpoints require Bearer token in Authorization header</p>
				<p><strong>Test your APIs:</strong> Use curl, Postman, or your browser!</p>
			</div>
		</body>
		</html>
		`
		
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}).Methods("GET")

	// API routes using your actual handlers
	api := r.PathPrefix("/api/v1").Subrouter()
	
	// Story Topics routes - using your actual API methods
	// api.HandleFunc("/story-topics", storyTopicsHandler.GetStoryTopics).Methods("POST")
	api.HandleFunc("/story", storyTopicsHandler.CreateStory).Methods("POST")
	// api.HandleFunc("/story-topics/{id}", storyTopicsHandler.GetStory).Methods("POST")
	// api.HandleFunc("/story-topics/{id}", storyTopicsHandler.UpdateStory).Methods("PUT")
	// api.HandleFunc("/story-topics/{id}", storyTopicsHandler.DeleteStory).Methods("DELETE")
	
	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Start server
	// Create server with graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		log.Println("📚 Documentation available at: http://localhost:8080/docs")
		log.Println("🔒 All APIs require authentication")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Close database connections
	if err := storyDB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	if err := storageService.Close(); err != nil {
		log.Printf("Error closing storage service: %v", err)
	}

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
