package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"rio-go-model/configs"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// CustomLogger represents a custom logger with rotation capabilities
type CustomLogger struct {
	logger   *log.Logger
	file     *os.File
	settings *configs.Settings
	level    LogLevel
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Logger    string
	Message   string
}

// NewCustomLogger creates a new custom logger instance
func NewCustomLogger(name string, settings *configs.Settings) *CustomLogger {
	level := parseLogLevel(settings.LogLevel)
	
	// Create logs directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Warning: Could not create logs directory: %v", err)
	}

	// Create or open log file
	logPath := filepath.Join(logDir, settings.LogFile)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Warning: Could not open log file, using stdout: %v", err)
		file = nil
	}

	// Create multi-writer for both console and file
	var writer io.Writer
	if file != nil {
		writer = io.MultiWriter(os.Stdout, file)
	} else {
		writer = os.Stdout
	}

	// Create logger with custom format
	logger := log.New(writer, fmt.Sprintf("[%s] ", name), 0)

	return &CustomLogger{
		logger:   logger,
		file:     file,
		settings: settings,
		level:    level,
	}
}

// parseLogLevel parses string log level to LogLevel enum
func parseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

// levelToString converts LogLevel to string
func levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "INFO"
	}
}

// shouldLog checks if message should be logged based on level
func (cl *CustomLogger) shouldLog(level LogLevel) bool {
	return level >= cl.level
}

// formatMessage formats the log message according to settings
func (cl *CustomLogger) formatMessage(level LogLevel, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := levelToString(level)
	
	// Use a simplified format that's more Go-friendly
	return fmt.Sprintf("%s - %s - %s", timestamp, levelStr, message)
}

// Debug logs a debug message
func (cl *CustomLogger) Debug(message string) {
	if cl.shouldLog(DEBUG) {
		formatted := cl.formatMessage(DEBUG, message)
		cl.logger.Println(formatted)
	}
}

// Debugf logs a formatted debug message
func (cl *CustomLogger) Debugf(format string, args ...interface{}) {
	if cl.shouldLog(DEBUG) {
		message := fmt.Sprintf(format, args...)
		cl.Debug(message)
	}
}

// Info logs an info message
func (cl *CustomLogger) Info(message string) {
	if cl.shouldLog(INFO) {
		formatted := cl.formatMessage(INFO, message)
		cl.logger.Println(formatted)
	}
}

// Infof logs a formatted info message
func (cl *CustomLogger) Infof(format string, args ...interface{}) {
	if cl.shouldLog(INFO) {
		message := fmt.Sprintf(format, args...)
		cl.Info(message)
	}
}

// Warn logs a warning message
func (cl *CustomLogger) Warn(message string) {
	if cl.shouldLog(WARN) {
		formatted := cl.formatMessage(WARN, message)
		cl.logger.Println(formatted)
	}
}

// Warnf logs a formatted warning message
func (cl *CustomLogger) Warnf(format string, args ...interface{}) {
	if cl.shouldLog(WARN) {
		message := fmt.Sprintf(format, args...)
		cl.Warn(message)
	}
}

// Error logs an error message
func (cl *CustomLogger) Error(message string) {
	if cl.shouldLog(ERROR) {
		formatted := cl.formatMessage(ERROR, message)
		cl.logger.Println(formatted)
	}
}

// Errorf logs a formatted error message
func (cl *CustomLogger) Errorf(format string, args ...interface{}) {
	if cl.shouldLog(ERROR) {
		message := fmt.Sprintf(format, args...)
		cl.Error(message)
	}
}

// Close closes the logger and any open files
func (cl *CustomLogger) Close() error {
	if cl.file != nil {
		return cl.file.Close()
	}
	return nil
}

// Global logger instance
var globalLogger *CustomLogger

// SetupLogging sets up global logging configuration
func SetupLogging(settings *configs.Settings) {
	globalLogger = NewCustomLogger("app", settings)
	
	// Log the logging setup
	globalLogger.Infof("Logging configured - Level: %s, File: %s", 
		settings.LogLevel, filepath.Join("logs", settings.LogFile))
}

// GetLogger returns a logger instance for the specified name
func GetLogger(name string, settings *configs.Settings) *CustomLogger {
	return NewCustomLogger(name, settings)
}

// Global logging functions that use the global logger

// Debug logs a debug message using the global logger
func Debug(message string) {
	if globalLogger != nil {
		globalLogger.Debug(message)
	}
}

// Debugf logs a formatted debug message using the global logger
func Debugf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debugf(format, args...)
	}
}

// Info logs an info message using the global logger
func Info(message string) {
	if globalLogger != nil {
		globalLogger.Info(message)
	}
}

// Infof logs a formatted info message using the global logger
func Infof(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Infof(format, args...)
	}
}

// Warn logs a warning message using the global logger
func Warn(message string) {
	if globalLogger != nil {
		globalLogger.Warn(message)
	}
}

// Warnf logs a formatted warning message using the global logger
func Warnf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warnf(format, args...)
	}
}

// Error logs an error message using the global logger
func Error(message string) {
	if globalLogger != nil {
		globalLogger.Error(message)
	}
}

// Errorf logs a formatted error message using the global logger
func Errorf(format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Errorf(format, args...)
	}
}

// CloseLogging closes the global logger
func CloseLogging() error {
	if globalLogger != nil {
		return globalLogger.Close()
	}
	return nil
}
