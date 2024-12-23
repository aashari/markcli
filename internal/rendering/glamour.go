package rendering

import (
	"github.com/charmbracelet/glamour"
)

// GlamourRenderer wraps the Glamour renderer with our custom configuration
type GlamourRenderer struct {
	renderer *glamour.TermRenderer
}

// NewGlamourRenderer creates a new Glamour renderer with default settings
func NewGlamourRenderer() (*GlamourRenderer, error) {
	renderer, err := glamour.NewTermRenderer(
		// Use the dark theme as base
		glamour.WithEnvironmentConfig(),
		// Override specific settings
		glamour.WithWordWrap(0), // Disable word wrapping
		glamour.WithEmoji(),     // Enable emoji support
		glamour.WithBaseURL(""), // Disable base URL to show full links
	)
	if err != nil {
		return nil, err
	}

	return &GlamourRenderer{
		renderer: renderer,
	}, nil
}

// Render takes markdown input and returns the rendered output
func (r *GlamourRenderer) Render(markdown string) (string, error) {
	return r.renderer.Render(markdown)
}
