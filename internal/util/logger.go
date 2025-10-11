package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// RioLogLevel represents the severity level of a log message
type RioLogLevel int

const (
	RIO_DEBUG RioLogLevel = iota
	RIO_INFO
	RIO_WARN
	RIO_ERROR
	RIO_FATAL
)

// LoggerConfig holds configuration for the logger
type LoggerConfig struct {
	Level      RioLogLevel
	ShowCaller bool
	ShowLine   bool
	ShowTime   bool
	ShowDate   bool
	LogToFile  bool
	FilePath   string
}

// RioLogger is a custom logger with enhanced features
type RioLogger struct {
	config *LoggerConfig
	logger *log.Logger
	file   *os.File
}

// DefaultLoggerConfig returns a default configuration
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:      RIO_INFO,
		ShowCaller: true,
		ShowLine:   true,
		ShowTime:   true,
		ShowDate:   true,
		LogToFile:  true,
		FilePath:   "logs/rio-app.log",
	}
}

// NewRioLogger creates a new RioLogger instance
func NewRioLogger(className string, config *LoggerConfig) *RioLogger {
	if config == nil {
		config = DefaultLoggerConfig()
	}

	// Ensure log directory exists
	if config.LogToFile {
		dir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create log directory: %v", err)
			config.LogToFile = false
		}
	}

	var file *os.File
	var logger *log.Logger

	if config.LogToFile {
		var err error
		file, err = os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Failed to open log file: %v", err)
			config.LogToFile = false
		} else {
			logger = log.New(file, "", 0)
		}
	}

	if !config.LogToFile || logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}

	return &RioLogger{
		config: config,
		logger: logger,
		file:   file,
	}
}

// Close closes the logger and any open files
func (rl *RioLogger) Close() {
	if rl.file != nil {
		rl.file.Close()
	}
}

// getCallerInfo returns information about the calling function
func (rl *RioLogger) getCallerInfo() (className, methodName string, lineNumber int) {
	pc, file, line, ok := runtime.Caller(3) // Skip 3 levels to get the actual caller
	if !ok {
		return "Unknown", "Unknown", 0
	}

	// Extract class name from file path
	fileName := filepath.Base(file)
	className = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// Extract method name
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		fullName := fn.Name()
		parts := strings.Split(fullName, ".")
		if len(parts) > 0 {
			methodName = parts[len(parts)-1]
		}
	}

	return className, methodName, line
}

// formatMessage formats the log message with all required information
func (rl *RioLogger) formatMessage(level RioLogLevel, methodName string, message string, err error) string {
	var parts []string

	// Add timestamp
	if rl.config.ShowDate || rl.config.ShowTime {
		now := time.Now()
		if rl.config.ShowDate && rl.config.ShowTime {
			parts = append(parts, fmt.Sprintf("[%s]", now.Format("2006-01-02 15:04:05")))
		} else if rl.config.ShowDate {
			parts = append(parts, fmt.Sprintf("[%s]", now.Format("2006-01-02")))
		} else if rl.config.ShowTime {
			parts = append(parts, fmt.Sprintf("[%s]", now.Format("15:04:05")))
		}
	}

	// Add log level
	levelStr := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[level]
	parts = append(parts, fmt.Sprintf("[%s]", levelStr))

	// Add class and method info
	if rl.config.ShowCaller {
		className, _, lineNumber := rl.getCallerInfo()
		if rl.config.ShowLine {
			parts = append(parts, fmt.Sprintf("[%s::%s:%d]", className, methodName, lineNumber))
		} else {
			parts = append(parts, fmt.Sprintf("[%s::%s]", className, methodName))
		}
	}

	// Add message
	parts = append(parts, message)

	// Add error details if present
	if err != nil {
		parts = append(parts, fmt.Sprintf("| Error: %v", err))

		// Add full stack trace for errors
		if level == RIO_ERROR || level == RIO_FATAL {
			stackTrace := rl.getStackTrace()
			parts = append(parts, fmt.Sprintf("| StackTrace: %s", stackTrace))
		}
	}

	return strings.Join(parts, " ")
}

// getStackTrace returns a formatted stack trace
func (rl *RioLogger) getStackTrace() string {
	buf := make([]byte, 1024*4)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// log writes the formatted message to the logger
func (rl *RioLogger) log(level RioLogLevel, methodName string, message string, err error) {
	if level < rl.config.Level {
		return
	}

	formattedMessage := rl.formatMessage(level, methodName, message, err)
	rl.logger.Println(formattedMessage)

	// For fatal errors, also log to stderr and exit
	if level == RIO_FATAL {
		fmt.Fprintf(os.Stderr, "FATAL ERROR: %s\n", formattedMessage)
		os.Exit(1)
	}
}

// Debug logs a debug message
func (rl *RioLogger) Debug(methodName, message string) {
	rl.log(RIO_DEBUG, methodName, message, nil)
}

// Info logs an info message
func (rl *RioLogger) Info(methodName, message string) {
	rl.log(RIO_INFO, methodName, message, nil)
}

// Warn logs a warning message
func (rl *RioLogger) Warn(methodName, message string) {
	rl.log(RIO_WARN, methodName, message, nil)
}

// Error logs an error message with optional error details
func (rl *RioLogger) Error(methodName, message string, err error) {
	rl.log(RIO_ERROR, methodName, message, err)
}

// Fatal logs a fatal error message and exits the program
func (rl *RioLogger) Fatal(methodName, message string, err error) {
	rl.log(RIO_FATAL, methodName, message, err)
}

// Debugf logs a formatted debug message
func (rl *RioLogger) Debugf(methodName, format string, args ...interface{}) {
	rl.Debug(methodName, fmt.Sprintf(format, args...))
}

// Infof logs a formatted info message
func (rl *RioLogger) Infof(methodName, format string, args ...interface{}) {
	rl.Info(methodName, fmt.Sprintf(format, args...))
}

// Warnf logs a formatted warning message
func (rl *RioLogger) Warnf(methodName, format string, args ...interface{}) {
	rl.Warn(methodName, fmt.Sprintf(format, args...))
}

// Errorf logs a formatted error message with optional error details
func (rl *RioLogger) Errorf(methodName, format string, err error, args ...interface{}) {
	rl.Error(methodName, fmt.Sprintf(format, args...), err)
}

// Fatalf logs a formatted fatal error message and exits the program
func (rl *RioLogger) Fatalf(methodName, format string, err error, args ...interface{}) {
	rl.Fatal(methodName, fmt.Sprintf(format, args...), err)
}

// LogRequest logs HTTP request details
func (rl *RioLogger) LogRequest(methodName, method, url string, statusCode int, duration time.Duration) {
	rl.Infof(methodName, "HTTP %s %s | Status: %d | Duration: %v", method, url, statusCode, duration)
}

// LogAPI logs API call details
func (rl *RioLogger) LogAPI(methodName, apiName, operation string, success bool, duration time.Duration) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	rl.Infof(methodName, "API %s | Operation: %s | Status: %s | Duration: %v", apiName, operation, status, duration)
}

// LogDatabase logs database operation details
func (rl *RioLogger) LogDatabase(methodName, operation, table string, success bool, duration time.Duration) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	rl.Infof(methodName, "DB %s | Table: %s | Status: %s | Duration: %v", operation, table, status, duration)
}

// LogStoryGeneration logs story generation details
func (rl *RioLogger) LogStoryGeneration(methodName, theme, topic, language string, success bool, duration time.Duration, wordCount int) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	rl.Infof(methodName, "STORY_GEN | Theme: %s | Topic: %s | Language: %s | Status: %s | Duration: %v | Words: %d",
		theme, topic, language, status, duration, wordCount)
}

// LogTopicGeneration logs topic generation details
func (rl *RioLogger) LogTopicGeneration(methodName, theme, language string, success bool, duration time.Duration, topicCount int) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	rl.Infof(methodName, "TOPIC_GEN | Theme: %s | Language: %s | Status: %s | Duration: %v | Count: %d",
		theme, language, status, duration, topicCount)
}

// LogAudioGeneration logs audio generation details
func (rl *RioLogger) LogAudioGeneration(methodName, text string, success bool, duration time.Duration, audioLength time.Duration) {
	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	rl.Infof(methodName, "AUDIO_GEN | Text Length: %d | Status: %s | Duration: %v | Audio Length: %v",
		len(text), status, duration, audioLength)
}

// LogPerformance logs performance metrics
func (rl *RioLogger) LogPerformance(methodName, operation string, metrics map[string]interface{}) {
	var parts []string
	for key, value := range metrics {
		parts = append(parts, fmt.Sprintf("%s: %v", key, value))
	}
	rl.Infof(methodName, "PERFORMANCE | %s | %s", operation, strings.Join(parts, " | "))
}

// LogSecurity logs security-related events
func (rl *RioLogger) LogSecurity(methodName, event, details string) {
	rl.Warnf(methodName, "SECURITY | Event: %s | Details: %s", event, details)
}

// LogUserAction logs user actions
func (rl *RioLogger) LogUserAction(methodName, userID, action, details string) {
	rl.Infof(methodName, "USER_ACTION | User: %s | Action: %s | Details: %s", userID, action, details)
}

// LogSystemEvent logs system events
func (rl *RioLogger) LogSystemEvent(methodName, event, details string) {
	rl.Infof(methodName, "SYSTEM | Event: %s | Details: %s", event, details)
}

// LogConfigChange logs configuration changes
func (rl *RioLogger) LogConfigChange(methodName, configKey, oldValue, newValue string) {
	rl.Infof(methodName, "CONFIG_CHANGE | Key: %s | Old: %s | New: %s", configKey, oldValue, newValue)
}

// LogMemoryUsage logs memory usage information
func (rl *RioLogger) LogMemoryUsage(methodName string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	rl.Debugf(methodName, "MEMORY | Alloc: %d KB | TotalAlloc: %d KB | Sys: %d KB | NumGC: %d",
		m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC)
}

// LogGoroutineCount logs the current number of goroutines
func (rl *RioLogger) LogGoroutineCount(methodName string) {
	count := runtime.NumGoroutine()
	rl.Debugf(methodName, "GOROUTINES | Count: %d", count)
}

// SetLevel changes the log level
func (rl *RioLogger) SetLevel(level RioLogLevel) {
	rl.config.Level = level
}

// GetLevel returns the current log level
func (rl *RioLogger) GetLevel() RioLogLevel {
	return rl.config.Level
}

// IsDebugEnabled returns true if debug logging is enabled
func (rl *RioLogger) IsDebugEnabled() bool {
	return rl.config.Level <= RIO_DEBUG
}

// IsInfoEnabled returns true if info logging is enabled
func (rl *RioLogger) IsInfoEnabled() bool {
	return rl.config.Level <= RIO_INFO
}

// IsWarnEnabled returns true if warning logging is enabled
func (rl *RioLogger) IsWarnEnabled() bool {
	return rl.config.Level <= RIO_WARN
}

// IsErrorEnabled returns true if error logging is enabled
func (rl *RioLogger) IsErrorEnabled() bool {
	return rl.config.Level <= RIO_ERROR
}

// Helper function to create a logger for a specific class
func CreateLogger(className string) *RioLogger {
	return NewRioLogger(className, DefaultLoggerConfig())
}

// Helper function to create a logger with custom config
func CreateLoggerWithConfig(className string, config *LoggerConfig) *RioLogger {
	return NewRioLogger(className, config)
}

// Helper function to create a development logger (debug level, no file)
func CreateDevLogger(className string) *RioLogger {
	config := &LoggerConfig{
		Level:      RIO_DEBUG,
		ShowCaller: true,
		ShowLine:   true,
		ShowTime:   true,
		ShowDate:   true,
		LogToFile:  false,
	}
	return NewRioLogger(className, config)
}

// Helper function to create a production logger (info level, file logging)
func CreateProdLogger(className string) *RioLogger {
	config := &LoggerConfig{
		Level:      RIO_INFO,
		ShowCaller: true,
		ShowLine:   true,
		ShowTime:   true,
		ShowDate:   true,
		LogToFile:  true,
		FilePath:   "logs/rio-app-prod.log",
	}
	return NewRioLogger(className, config)
}

// Helper function to create a test logger (warn level, no file)
func CreateTestLogger(className string) *RioLogger {
	config := &LoggerConfig{
		Level:      RIO_WARN,
		ShowCaller: true,
		ShowLine:   true,
		ShowTime:   false,
		ShowDate:   false,
		LogToFile:  false,
	}
	return NewRioLogger(className, config)
}
