package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var isDebugEnabled bool

// List of patterns to identify sensitive data
var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`"token":\s*"[^"]*"`),
	regexp.MustCompile(`"password":\s*"[^"]*"`),
	regexp.MustCompile(`"api_key":\s*"[^"]*"`),
	regexp.MustCompile(`"secret":\s*"[^"]*"`),
	regexp.MustCompile(`"Authorization":\s*"[^"]*"`),
	regexp.MustCompile(`"auth":\s*"[^"]*"`),
}

// EnableDebug enables debug logging
func EnableDebug() {
	isDebugEnabled = true
}

// IsDebugEnabled returns whether debug logging is enabled
func IsDebugEnabled() bool {
	return isDebugEnabled
}

// LogDebug logs a debug message if debug mode is enabled
func LogDebug(format string, args ...interface{}) {
	if !isDebugEnabled {
		return
	}
	fmt.Fprintf(os.Stderr, "DEBUG: "+format+"\n", args...)
}

// redactSensitiveData replaces sensitive data with "[REDACTED]"
func redactSensitiveData(data string) string {
	for _, pattern := range sensitivePatterns {
		data = pattern.ReplaceAllStringFunc(data, func(match string) string {
			parts := strings.SplitN(match, ":", 2)
			if len(parts) != 2 {
				return match
			}
			return fmt.Sprintf("%s: \"[REDACTED]\"", parts[0])
		})
	}
	return data
}

// LogJSONInline logs a JSON object in a single line if debug mode is enabled
func LogJSONInline(prefix string, v interface{}) {
	if !isDebugEnabled {
		return
	}
	data, err := json.Marshal(v)
	if err == nil {
		redactedData := redactSensitiveData(string(data))
		LogDebug("%s: %s", prefix, redactedData)
	} else {
		LogDebug("%s: failed to marshal JSON: %v", prefix, err)
	}
}

// LogJSON logs a JSON object with indentation if debug mode is enabled
func LogJSON(prefix string, v interface{}) {
	if !isDebugEnabled {
		return
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		redactedData := redactSensitiveData(string(data))
		LogDebug("%s:\n%s", prefix, redactedData)
	} else {
		LogDebug("%s: failed to marshal JSON: %v", prefix, err)
	}
}
