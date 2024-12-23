package logging

import (
	"encoding/json"
	"fmt"
	"os"
)

var isDebugEnabled bool

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

// LogJSONInline logs a JSON object in a single line if debug mode is enabled
func LogJSONInline(prefix string, v interface{}) {
	if !isDebugEnabled {
		return
	}
	data, err := json.Marshal(v)
	if err != nil {
		LogDebug("%s: failed to marshal JSON: %v", prefix, err)
		return
	}
	LogDebug("%s: %s", prefix, string(data))
}

// LogJSON logs a JSON object with indentation if debug mode is enabled
func LogJSON(prefix string, v interface{}) {
	if !isDebugEnabled {
		return
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		LogDebug("%s: failed to marshal JSON: %v", prefix, err)
		return
	}
	LogDebug("%s:\n%s", prefix, string(data))
}
