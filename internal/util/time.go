package util

import (
	"time"
)

// ParseDate attempts to parse a date using various formats
func ParseDate(dateString string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05.999-0700",
		"2006-01-02T15:04:05.999+0700",
		time.RFC3339,
		"2006-01-02T15:04:05Z",
	}

	var t time.Time
	var err error

	for _, layout := range layouts {
		t, err = time.Parse(layout, dateString)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, err
}
