package util

import (
	"errors"
	"time"
)

// ExampleUsage demonstrates how to use the RioLogger
func ExampleUsage() {
	// Create a logger for a specific class
	logger := CreateLogger("StoryGenerator")
	defer logger.Close()

	// Basic logging
	logger.Info("CreateStory", "Starting story generation")
	logger.Debug("CreateStory", "Debug information about story creation")
	logger.Warn("CreateStory", "Warning: Low memory detected")

	// Error logging with stack trace
	err := errors.New("failed to connect to AI service")
	logger.Error("CreateStory", "Failed to generate story", err)

	// Formatted logging
	theme := "PlanetProtector"
	topic := "Water Conservation"
	logger.Infof("CreateStory", "Generating story for theme: %s, topic: %s", theme, topic)

	// Error with formatted message
	logger.Errorf("CreateStory", "API call failed for %s: %v", err, "Gemini")

	// Specialized logging methods
	start := time.Now()
	time.Sleep(100 * time.Millisecond) // Simulate work
	duration := time.Since(start)

	// Log story generation
	logger.LogStoryGeneration("CreateStory", theme, topic, "English", true, duration, 250)

	// Log API call
	logger.LogAPI("CreateStory", "Gemini", "GenerateContent", true, duration)

	// Log database operation
	logger.LogDatabase("CreateStory", "INSERT", "stories", true, duration)

	// Log performance metrics
	metrics := map[string]interface{}{
		"memory_usage":  "45MB",
		"cpu_usage":     "12%",
		"response_time": "150ms",
	}
	logger.LogPerformance("CreateStory", "Story Generation", metrics)

	// Log user action
	logger.LogUserAction("CreateStory", "user123", "story_request", "Requested PlanetProtector story")

	// Log system event
	logger.LogSystemEvent("CreateStory", "service_start", "Story generation service started")

	// Log memory usage
	logger.LogMemoryUsage("CreateStory")

	// Log goroutine count
	logger.LogGoroutineCount("CreateStory")
}

// ExampleClass demonstrates logger usage in a class
type ExampleClass struct {
	logger *RioLogger
}

// NewExampleClass creates a new instance with logger
func NewExampleClass() *ExampleClass {
	return &ExampleClass{
		logger: CreateLogger("ExampleClass"),
	}
}

// ProcessData demonstrates method-level logging
func (ec *ExampleClass) ProcessData(data string) error {
	ec.logger.Info("ProcessData", "Starting data processing")

	if data == "" {
		err := errors.New("empty data provided")
		ec.logger.Error("ProcessData", "Invalid input data", err)
		return err
	}

	ec.logger.Debugf("ProcessData", "Processing data with length: %d", len(data))

	// Simulate processing
	time.Sleep(50 * time.Millisecond)

	ec.logger.Info("ProcessData", "Data processing completed successfully")
	return nil
}

// Close cleans up resources
func (ec *ExampleClass) Close() {
	ec.logger.Close()
}

// ExampleDifferentLogLevels demonstrates different log levels
func ExampleDifferentLogLevels() {
	logger := CreateLogger("LogLevelExample")
	defer logger.Close()

	// These will all be logged (DEBUG level)
	logger.SetLevel(RIO_DEBUG)
	logger.Debug("TestMethod", "This is a debug message")
	logger.Info("TestMethod", "This is an info message")
	logger.Warn("TestMethod", "This is a warning message")
	logger.Error("TestMethod", "This is an error message", errors.New("test error"))

	// Only WARN and above will be logged
	logger.SetLevel(RIO_WARN)
	logger.Debug("TestMethod", "This debug message will NOT be logged")
	logger.Info("TestMethod", "This info message will NOT be logged")
	logger.Warn("TestMethod", "This warning message WILL be logged")
	logger.Error("TestMethod", "This error message WILL be logged", errors.New("test error"))
}

// ExampleCustomConfig demonstrates custom logger configuration
func ExampleCustomConfig() {
	config := &LoggerConfig{
		Level:      RIO_ERROR, // Only log errors and fatal
		ShowCaller: true,
		ShowLine:   true,
		ShowTime:   true,
		ShowDate:   false,
		LogToFile:  true,
		FilePath:   "logs/errors-only.log",
	}

	logger := CreateLoggerWithConfig("CustomLogger", config)
	defer logger.Close()

	// Only this will be logged
	logger.Error("TestMethod", "This error will be logged", errors.New("custom error"))

	// These will not be logged due to level setting
	logger.Debug("TestMethod", "This will not be logged")
	logger.Info("TestMethod", "This will not be logged")
	logger.Warn("TestMethod", "This will not be logged")
}

// ExampleLineNumberFeature demonstrates line number logging
func ExampleLineNumberFeature() {
	// Logger with line numbers enabled
	loggerWithLines := CreateLogger("LineNumberExample")
	defer loggerWithLines.Close()

	// This will show line numbers
	loggerWithLines.Info("TestMethod", "This message will show line numbers")

	// Logger without line numbers
	config := &LoggerConfig{
		Level:      RIO_INFO,
		ShowCaller: true,
		ShowLine:   false, // Disable line numbers
		ShowTime:   true,
		ShowDate:   true,
		LogToFile:  false,
	}
	loggerWithoutLines := CreateLoggerWithConfig("NoLineExample", config)
	defer loggerWithoutLines.Close()

	// This will NOT show line numbers
	loggerWithoutLines.Info("TestMethod", "This message will NOT show line numbers")
}

// ExampleEnvironmentSpecificLoggers demonstrates different logger types
func ExampleEnvironmentSpecificLoggers() {
	// Development logger (debug level, console only)
	devLogger := CreateDevLogger("DevClass")
	devLogger.Debug("TestMethod", "Development debug message")
	devLogger.Close()

	// Production logger (info level, file logging)
	prodLogger := CreateProdLogger("ProdClass")
	prodLogger.Info("TestMethod", "Production info message")
	prodLogger.Close()

	// Test logger (warn level, console only)
	testLogger := CreateTestLogger("TestClass")
	testLogger.Warn("TestMethod", "Test warning message")
	testLogger.Close()
}
