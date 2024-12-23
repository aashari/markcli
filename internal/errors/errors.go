package errors

import (
	"errors"
	"fmt"
)

// Common errors
var (
	ErrNoAtlassianConfig = errors.New("no Atlassian configuration found")
	ErrInvalidInput      = errors.New("invalid input")
)

// WrapError wraps an error with additional context
func WrapError(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
