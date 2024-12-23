package atlassian

import (
	"strings"
)

// truncateText truncates text to the specified length, adding ellipsis if needed
func truncateText(text string, maxLength int) string {
	text = strings.TrimSpace(text)
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength-3] + "..."
}

// wrapText wraps text at word boundaries to the specified line length
func wrapText(text string, lineLength int) string {
	var wrapped strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for i, word := range words {
		wordLen := len(word)
		if lineLen+wordLen+1 > lineLength && lineLen > 0 {
			wrapped.WriteString("\n")
			lineLen = 0
		} else if i > 0 {
			wrapped.WriteString(" ")
			lineLen++
		}
		wrapped.WriteString(word)
		lineLen += wordLen
	}

	return wrapped.String()
}

// cleanTitle replaces highlight markers with bold markers
func cleanTitle(title string) string {
	title = strings.ReplaceAll(title, "@@@hl@@@", "**")
	title = strings.ReplaceAll(title, "@@@endhl@@@", "**")
	return title
}

// formatURL makes the URL more readable
func formatURL(url string) string {
	url = strings.ReplaceAll(url, "+", " ")
	url = strings.TrimPrefix(url, "/spaces/")
	url = strings.TrimPrefix(url, "/wiki/spaces/")
	url = strings.TrimPrefix(url, "/pages/")
	return url
}

// cleanContent removes newlines, headings, and lists from the content
func cleanContent(content string) string {
	content = strings.ReplaceAll(content, "@@@hl@@@", "**")
	content = strings.ReplaceAll(content, "@@@endhl@@@", "**")

	lines := strings.Split(content, "\n")
	var cleaned []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "---" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "-") {
			continue
		}
		cleaned = append(cleaned, line)
	}

	result := strings.Join(cleaned, " ")
	result = strings.ReplaceAll(result, "  ", " ")
	result = strings.ReplaceAll(result, "..", ".")
	result = strings.ReplaceAll(result, ". .", ".")
	return strings.TrimSpace(result)
}
