# Rio Logger - Enhanced Logging Utility

A comprehensive logging utility designed specifically for the Rio Go Model application, providing detailed logging with class names, method names, timestamps, and full stack traces for errors.

## Features

- **Class and Method Tracking**: Automatically captures class name and method name
- **Line Number Support**: Optional line number display for precise debugging
- **Timestamp Support**: Configurable date and time formatting
- **Error Stack Traces**: Full stack traces for error and fatal level logs
- **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **File and Console Logging**: Configurable output destinations
- **Specialized Logging Methods**: Pre-built methods for common operations
- **Performance Monitoring**: Built-in memory and goroutine tracking
- **Thread-Safe**: Safe for concurrent use

## Quick Start

### Basic Usage

```go
package main

import (
    "rio-go-model/internal/util"
    "errors"
)

func main() {
    // Create a logger for your class
    logger := util.CreateLogger("MyClass")
    defer logger.Close()

    // Basic logging
    logger.Info("MyMethod", "Application started")
    logger.Debug("MyMethod", "Debug information")
    logger.Warn("MyMethod", "Warning message")
    
    // Error logging with stack trace
    err := errors.New("something went wrong")
    logger.Error("MyMethod", "Operation failed", err)
}
```

### Class-Based Usage

```go
type StoryGenerator struct {
    logger *util.RioLogger
}

func NewStoryGenerator() *StoryGenerator {
    return &StoryGenerator{
        logger: util.CreateLogger("StoryGenerator"),
    }
}

func (sg *StoryGenerator) GenerateStory(theme string) error {
    sg.logger.Info("GenerateStory", "Starting story generation")
    
    if theme == "" {
        err := errors.New("theme cannot be empty")
        sg.logger.Error("GenerateStory", "Invalid theme provided", err)
        return err
    }
    
    sg.logger.Infof("GenerateStory", "Generating story for theme: %s", theme)
    
    // ... story generation logic ...
    
    sg.logger.Info("GenerateStory", "Story generation completed successfully")
    return nil
}

func (sg *StoryGenerator) Close() {
    sg.logger.Close()
}
```

## Log Levels

The logger supports five log levels in order of severity:

1. **RIO_DEBUG**: Detailed information for debugging
2. **RIO_INFO**: General information about program execution
3. **RIO_WARN**: Warning messages for potentially harmful situations
4. **RIO_ERROR**: Error messages for error conditions
5. **RIO_FATAL**: Fatal errors that cause the program to exit

### Setting Log Levels

```go
logger := util.CreateLogger("MyClass")

// Set to only show warnings and above
logger.SetLevel(util.RIO_WARN)

// Check if a level is enabled
if logger.IsDebugEnabled() {
    logger.Debug("MyMethod", "This will only log if DEBUG is enabled")
}
```

## Specialized Logging Methods

### Story Generation Logging

```go
logger.LogStoryGeneration("CreateStory", "PlanetProtector", "Water Conservation", "English", true, duration, 250)
```

### API Call Logging

```go
logger.LogAPI("CreateStory", "Gemini", "GenerateContent", true, duration)
```

### Database Operation Logging

```go
logger.LogDatabase("CreateStory", "INSERT", "stories", true, duration)
```

### Performance Monitoring

```go
// Log performance metrics
metrics := map[string]interface{}{
    "memory_usage": "45MB",
    "cpu_usage":    "12%",
    "response_time": "150ms",
}
logger.LogPerformance("CreateStory", "Story Generation", metrics)

// Log memory usage
logger.LogMemoryUsage("CreateStory")

// Log goroutine count
logger.LogGoroutineCount("CreateStory")
```

### User Action Logging

```go
logger.LogUserAction("CreateStory", "user123", "story_request", "Requested PlanetProtector story")
```

### Security Event Logging

```go
logger.LogSecurity("CreateStory", "unauthorized_access", "User attempted to access restricted resource")
```

## Logger Configurations

### Default Logger

```go
logger := util.CreateLogger("MyClass")
// Uses default configuration: INFO level, file logging, timestamps enabled
```

### Development Logger

```go
logger := util.CreateDevLogger("MyClass")
// DEBUG level, console only, all timestamps enabled
```

### Production Logger

```go
logger := util.CreateProdLogger("MyClass")
// INFO level, file logging to logs/rio-app-prod.log
```

### Test Logger

```go
logger := util.CreateTestLogger("MyClass")
// WARN level, console only, minimal timestamps
```

### Custom Configuration

```go
config := &util.LoggerConfig{
    Level:      util.RIO_ERROR, // Only log errors and fatal
    ShowCaller: true,
    ShowLine:   true,           // Enable line number display
    ShowTime:   true,
    ShowDate:   false,
    LogToFile:  true,
    FilePath:   "logs/errors-only.log",
}

logger := util.CreateLoggerWithConfig("MyClass", config)
```

### Line Number Feature

The logger can optionally display line numbers for precise debugging:

```go
// Logger with line numbers (default)
logger := util.CreateLogger("MyClass")
logger.Info("MyMethod", "This will show: [MyClass::MyMethod:42]")

// Logger without line numbers
config := &util.LoggerConfig{
    Level:      util.RIO_INFO,
    ShowCaller: true,
    ShowLine:   false, // Disable line numbers
    ShowTime:   true,
    ShowDate:   true,
    LogToFile:  false,
}
logger := util.CreateLoggerWithConfig("MyClass", config)
logger.Info("MyMethod", "This will show: [MyClass::MyMethod]")
```

## Log Output Format

The logger produces structured output with the following format:

```
[2024-01-15 14:30:25] [INFO] [StoryGenerator::CreateStory:42] Starting story generation
[2024-01-15 14:30:25] [ERROR] [StoryGenerator::CreateStory:58] API call failed | Error: connection timeout | StackTrace: goroutine 1 [running]:...
```

### Format Components

- **Timestamp**: `[2024-01-15 14:30:25]` (configurable)
- **Log Level**: `[INFO]`, `[ERROR]`, etc.
- **Class::Method:Line**: `[StoryGenerator::CreateStory:42]` (line number optional)
- **Message**: The actual log message
- **Error Details**: `| Error: connection timeout` (for errors)
- **Stack Trace**: `| StackTrace: ...` (for errors and fatal)

## File Logging

By default, logs are written to `logs/rio-app.log`. The logger automatically creates the logs directory if it doesn't exist.

### Log File Locations

- **Default**: `logs/rio-app.log`
- **Production**: `logs/rio-app-prod.log`
- **Custom**: Specify in `LoggerConfig.FilePath`

## Error Handling

The logger includes comprehensive error handling:

- **File Creation Errors**: Falls back to console logging
- **Permission Errors**: Logs to stderr and continues
- **Fatal Errors**: Logs to both file and stderr, then exits

## Performance Considerations

- **Minimal Overhead**: Logger checks log level before formatting
- **Efficient String Building**: Uses string builder for complex messages
- **Configurable Output**: Disable file logging for high-performance scenarios
- **Memory Monitoring**: Built-in memory usage tracking

## Best Practices

1. **Always Close Loggers**: Use `defer logger.Close()` to ensure proper cleanup
2. **Use Appropriate Levels**: Don't use DEBUG in production
3. **Include Context**: Provide meaningful method names and messages
4. **Handle Errors**: Always log errors with context
5. **Monitor Performance**: Use built-in performance logging methods
6. **Use Specialized Methods**: Leverage pre-built logging methods for common operations

## Migration from Standard Logger

Replace standard Go logging:

```go
// Old way
log.Printf("Creating story for theme: %s", theme)

// New way
logger.Infof("CreateStory", "Creating story for theme: %s", theme)
```

## Thread Safety

The RioLogger is thread-safe and can be used concurrently across multiple goroutines. Each logger instance maintains its own configuration and file handles.

## Examples

See `logger_example.go` for comprehensive usage examples including:
- Basic logging operations
- Class-based logging patterns
- Different log level configurations
- Custom logger configurations
- Environment-specific logger setups

## Dependencies

- `log`: Standard Go logging package
- `os`: File operations
- `runtime`: Stack trace and memory information
- `time`: Timestamp formatting
- `strings`: String manipulation
- `path/filepath`: File path operations
