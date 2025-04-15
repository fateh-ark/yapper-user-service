package logger

import "time"

type logMessage struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Source    string                 `json:"source"`
	Component string                 `json:"component,omitempty"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Error     map[string]interface{} `json:"error,omitempty"`
	// Build     map[string]string      `json:"build,omitempty"`
}

type LogData struct {
	Timestamp *time.Time
	Level     LogLevel
	Component string
	Message   string
	Context   *map[string]interface{}
	Error     *map[string]interface{}
}

type LogLevel string

const (
	DebugLogLevel LogLevel = "DEBUG"
	InfoLogLevel  LogLevel = "INFO"
	WarnLogLevel  LogLevel = "WARN"
	ErrorLogLevel LogLevel = "ERROR"
	FatalLogLevel LogLevel = "FATAL"
)
