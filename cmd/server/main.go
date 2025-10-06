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
	"time"

	"rio-go-model/docs"
	"rio-go-model/internal/handlers"
	"rio-go-model/internal/services/database"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/rs/cors"
	"rio-go-model/internal/util"
)

// Global variables for services and handler, initialized in background
var storyDB *database.StoryDatabase
var storageService *database.StorageService
var storyTopicsHandler *handlers.Story
var authHandler *handlers.AuthHandler
var emailHandler *handlers.Email
var tcHandler *handlers.TcHandler
var storyFeedbackHandler *handlers.StoryFeedbackHandler
// var servicesReady bool // No longer needed

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Println("🔧 Initializing services at startup...")

	// Database
	storyDB = database.NewStoryDatabase()
	if err := storyDB.Init(ctx); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	log.Println("✅ Database service initialized successfully")

	// Storage
	storageService = database.NewStorageService("kutty_bucket")
	if err := storageService.Init(ctx); err != nil {
		log.Fatalf("❌ Failed to initialize storage service: %v", err)
	}
	log.Println("✅ Storage service initialized successfully")

	// Handlers
	storyTopicsHandler = handlers.NewStory(storyDB, storageService)
	authHandler = handlers.NewAuthHandler(storyDB)
	emailHandler = &handlers.Email{}
	tcHandler = handlers.NewTcHandler(storyDB)
	storyFeedbackHandler = handlers.NewStoryFeedbackHandler(storyDB)
	log.Println("✅ All services initialized successfully!")
}


func main() {
	defer util.RecoverPanic()
	// Load environment variables
	if err := godotenv.Load(); err != nil {
        log.Println("ℹ️  No .env file found, using system environment variables")
    }
    
	// Get port from environment variable, default to 8080 for local development
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get host from environment variable or use default
    host := os.Getenv("HOST")
    if host == "" {
        host = "localhost"
    }

    // Read allowed origins from env (supports dev default)
    frontendOrigins := os.Getenv("FRONTEND_ORIGINS")
    var allowedOrigins []string
    if frontendOrigins != "" {
        // Split by comma for multiple origins
        allowedOrigins = []string{frontendOrigins}
    } else {
        // Default origins for development and production
        allowedOrigins = []string{
            "http://localhost:8081",           // Local development
            "https://riokutty.web.app",        // Firebase production
            "https://riokutty.firebaseapp.com", // Firebase alternative domain
        }
    }

	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "Story API"
	docs.SwaggerInfo.Description = "A comprehensive API for generating and managing stories with AI"
	docs.SwaggerInfo.Version = "1.0"
	if host == "localhost" {
		docs.SwaggerInfo.Host = host + ":" + port
		docs.SwaggerInfo.Schemes = []string{"http"}
	} else {
		docs.SwaggerInfo.Host = host
		docs.SwaggerInfo.Schemes = []string{"https"}
	}
	 // This will be overridden in production
	docs.SwaggerInfo.BasePath = "/api/v1"
	// docs.SwaggerInfo.Schemes = []string{"https", "http"}


	log.Printf("🚀 Server will start on port: %s", port)

	// REMOVED: Background initialization goroutine.
	// Initialization will now happen lazily in the handler.

	// Directly initialize the handler with nil services.
	// The services will be created on the first request.
	// storyTopicsHandler = handlers.NewStory(nil, nil)


	// Create router
	r := mux.NewRouter()

	// Add health check endpoint immediately
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// REMOVED: Readiness check is no longer needed with lazy initialization.
	// The server is ready to accept traffic as soon as it starts.

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
					<h3><span class="method post">POST</span> <span class="url">/api/v1/story</span> <span class="auth">🔒 Auth Required</span></h3>
					<p><strong>Description:</strong> Create a new story</p>
					<p><strong>Headers:</strong> Authorization: Bearer &lt;token&gt;</p>
					<p><strong>Body:</strong> JSON with story data</p>
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

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.AfterScript("const url = new URL(window.location.href); url.port = ''; url.pathname = '/swagger/doc.json'; document.querySelector('.swagger-ui .topbar a').href = url.href;"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Register API routes with panic recovery middleware
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(util.HTTPPanicRecoveryMiddleware)
	api.HandleFunc("/story", storyTopicsHandler.CreateStory).Methods("POST")
	api.HandleFunc("/stories", storyTopicsHandler.ListStories).Methods("GET")
	api.HandleFunc("/user-profile", storyTopicsHandler.UserProfile).Methods("GET")
	api.HandleFunc("/user-profile", storyTopicsHandler.UpdateUserProfile).Methods("PUT")
	api.HandleFunc("/user-profile", storyTopicsHandler.DeleteUserProfile).Methods("DELETE")
	api.HandleFunc("/email", emailHandler.NewEmail).Methods("POST")
	api.HandleFunc("/tc", tcHandler.HandleTc).Methods("POST")
	api.HandleFunc("/story-feedback", storyFeedbackHandler.HandleStoryFeedback).Methods("POST")
	api.HandleFunc("/triggers/pubsub", handlers.PubSubPushHandler).Methods("POST")

	// Add the new authentication routes
	authRouter := api.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/google", authHandler.GoogleLogin).Methods("POST")
	authRouter.HandleFunc("/google/", authHandler.GoogleLogin).Methods("POST")
	authRouter.HandleFunc("/token/refresh", authHandler.RefreshToken).Methods("POST")
	authRouter.HandleFunc("/token/refresh/", authHandler.RefreshToken).Methods("POST")
	authRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	authRouter.HandleFunc("/logout/", authHandler.Logout).Methods("POST")
	authRouter.HandleFunc("/w/logout/", authHandler.LogoutWeb).Methods("POST")
	authRouter.HandleFunc("/w/token/refresh/", authHandler.RefreshTokenWeb).Methods("POST")
	
    // Configure CORS (use rs/cors only)
    c := cors.New(cors.Options{
        AllowedOrigins:   allowedOrigins,
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })
    
    // Log CORS configuration for debugging
    log.Printf("🌐 CORS configured with allowed origins: %v", allowedOrigins)
    // Do not also use the custom CorsMiddleware when using rs/cors
	// r.Use(CorsMiddleware)
	// Create server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: c.Handler(r),
	}

	// Start server immediately in main thread
	log.Printf("🌐 Starting server on :%s", port)
	log.Printf("📚 Documentation available at: http://localhost:%s/docs", port)
	log.Println("🔒 All APIs require authentication")
	log.Println("✅ Server is ready to accept connections")

	// Start server in main thread - this will block until shutdown
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server error: %v", err)
	}
	// Forcing redeployment to fix clock skew issue.
}
