# Rio Go Model API

A simple Go web service with two REST API endpoints built using Gorilla Mux.

## Features

- **Hello API**: Returns a greeting message with timestamp
- **Health API**: Returns service health status and version information
- CORS enabled for cross-origin requests
- Request logging middleware
- JSON response format

## API Endpoints

### 1. Hello API
- **URL**: `GET /api/v1/hello`
- **Description**: Returns a greeting message
- **Response**:
```json
{
  "message": "Hello from Rio Go Model API!",
  "timestamp": "2024-01-01T12:00:00Z",
  "status": "success"
}
```

### 2. Health API
- **URL**: `GET /api/v1/health`
- **Description**: Returns service health status
- **Response**:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "service": "rio-go-model",
  "version": "1.0.0"
}
```

## Prerequisites

- Go 1.21 or later
- Git

## Installation & Setup

1. **Clone the repository** (if applicable):
```bash
git clone <repository-url>
cd rio-go-model
```

2. **Install dependencies**:
```bash
go mod tidy
```

3. **Run the application**:
```bash
go run main.go
```

The server will start on port 8080.

## Testing the APIs

### Using curl:
```bash
# Test Hello API
curl http://localhost:8080/api/v1/hello

# Test Health API
curl http://localhost:8080/api/v1/health
```

### Using a web browser:
- Hello API: http://localhost:8080/api/v1/hello
- Health API: http://localhost:8080/api/v1/health

## Project Structure

```
rio-go-model/
├── go.mod          # Go module file with dependencies
├── main.go         # Main application with API handlers
└── README.md       # This file
```

## Dependencies

- `github.com/gorilla/mux`: HTTP router and URL matcher
- `github.com/gorilla/handlers`: HTTP middleware for Go

## Development

To add new endpoints, modify the `main.go` file and add new handler functions. The current structure makes it easy to extend with additional APIs.

## License

This project is open source and available under the MIT License.
