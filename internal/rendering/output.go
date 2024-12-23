package rendering

import (
	"fmt"
)

var defaultRenderer *GlamourRenderer

func init() {
	var err error
	defaultRenderer, err = NewGlamourRenderer()
	if err != nil {
		// If we can't create the renderer, we'll fall back to plain output
		fmt.Printf("Warning: Could not initialize Glamour renderer: %v\n", err)
	}
}

// RenderMarkdown takes markdown input and returns the rendered output.
// If the Glamour renderer is not available, it returns the plain markdown.
func RenderMarkdown(markdown string) string {
	if defaultRenderer == nil {
		return markdown
	}

	rendered, err := defaultRenderer.Render(markdown)
	if err != nil {
		// If rendering fails, fall back to plain markdown
		return markdown
	}

	return rendered
}

// PrintMarkdown renders and prints markdown to stdout
func PrintMarkdown(markdown string) {
	fmt.Print(RenderMarkdown(markdown))
}
