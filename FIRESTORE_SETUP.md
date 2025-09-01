# 🔥 Firestore Setup Guide

## Prerequisites

1. **Google Cloud Project** with Firestore enabled
2. **Service Account Key** with Firestore permissions
3. **Go 1.21+** installed

## Setup Steps

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Configure Environment Variables

Copy `configs/env.example` to `.env` and update:

```bash
# Google Cloud Configuration
GOOGLE_CLOUD_PROJECT=your-actual-project-id
GOOGLE_APPLICATION_CREDENTIALS=./configs/service-account-key.json

# Server Configuration
PORT=8080
HOST=localhost
```

### 3. Get Service Account Key

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Navigate to **IAM & Admin** > **Service Accounts**
3. Create a new service account or select existing one
4. Add **Firestore Admin** role
5. Create and download JSON key
6. Save as `configs/service-account-key.json`

### 4. Update Project ID

Replace `your-actual-project-id` in your `.env` file with your real Google Cloud project ID.

### 5. Run the Application

```bash
go run cmd/server/main.go
```

## Project Structure

```
rio-go-model/
├── cmd/server/main.go              # Main application entry point
├── configs/
│   ├── config.go                   # Configuration management
│   └── env.example                 # Environment variables template
├── internal/
│   ├── handlers/
│   │   └── story_api.go           # Your API handlers
│   └── services/
│       ├── app_service.go         # App lifecycle management
│       └── database/
│           └── firestore.go       # Firestore client
├── go.mod                          # Go dependencies
└── FIRESTORE_SETUP.md              # This file
```

## Features

✅ **Automatic Firestore Connection**  
✅ **Health Checks**  
✅ **Graceful Shutdown**  
✅ **Configuration Management**  
✅ **Error Handling**  
✅ **Context Management**  

## Testing

### Health Check
```bash
curl http://localhost:8080/health
```

### Documentation
```bash
open http://localhost:8080/docs
```

### Your APIs
```bash
# Get all story topics (requires auth)
curl -H "Authorization: Bearer your-token" \
  http://localhost:8080/api/v1/story-topics

# Create story topic (requires auth)
curl -X POST \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{"country":"USA","city":"New York","religions":["Christianity"],"preferences":["fiction"]}' \
  http://localhost:8080/api/v1/story-topics
```

## Troubleshooting

### Common Issues

1. **"Failed to create Firestore client"**
   - Check service account key path
   - Verify project ID is correct
   - Ensure service account has Firestore permissions

2. **"Permission denied"**
   - Add **Firestore Admin** role to service account
   - Check if Firestore is enabled in your project

3. **"Project not found"**
   - Verify project ID in environment variables
   - Check if you have access to the project

### Debug Mode

Set environment variable for verbose logging:
```bash
export GOOGLE_APPLICATION_CREDENTIALS=./configs/service-account-key.json
export GOOGLE_CLOUD_PROJECT=your-project-id
go run cmd/server/main.go
```

## Next Steps

1. **Add More Collections**: Extend Firestore operations
2. **Implement Caching**: Add Redis or in-memory cache
3. **Add Authentication**: Integrate with Firebase Auth
4. **Add Monitoring**: Cloud Logging and Metrics
5. **Deploy**: Deploy to Google Cloud Run or App Engine

## Support

- [Firestore Go Client Documentation](https://pkg.go.dev/cloud.google.com/go/firestore)
- [Google Cloud Console](https://console.cloud.google.com/)
- [Firestore Documentation](https://firebase.google.com/docs/firestore)

