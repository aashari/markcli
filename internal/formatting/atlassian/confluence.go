package formatting

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"markcli/internal/logging"
	"markcli/internal/types/atlassian"
	"markcli/internal/util"
)

// AtlassianConfluenceSpaceTableFormatter formats Confluence spaces as a markdown table
type AtlassianConfluenceSpaceTableFormatter struct {
	spaces []atlassian.AtlassianConfluenceSpace
}

// AtlassianConfluenceSearchResultsFormatter formats search results into markdown
type AtlassianConfluenceSearchResultsFormatter struct {
	results []atlassian.AtlassianConfluenceContentResult
}

// AtlassianConfluencePageDetailsFormatter formats a single Confluence page's details
type AtlassianConfluencePageDetailsFormatter struct {
	page atlassian.AtlassianConfluencePageDetails
}

// AtlassianConfluenceCreateSpaceTableFormatter creates a new SpaceTableFormatter
func AtlassianConfluenceCreateSpaceTableFormatter(spaces []atlassian.AtlassianConfluenceSpace) *AtlassianConfluenceSpaceTableFormatter {
	return &AtlassianConfluenceSpaceTableFormatter{
		spaces: spaces,
	}
}

// AtlassianConfluenceCreateSearchResultsFormatter creates a new search results formatter
func AtlassianConfluenceCreateSearchResultsFormatter(results []atlassian.AtlassianConfluenceContentResult) *AtlassianConfluenceSearchResultsFormatter {
	return &AtlassianConfluenceSearchResultsFormatter{
		results: results,
	}
}

// AtlassianConfluenceCreatePageDetailsFormatter creates a new PageDetailsFormatter
func AtlassianConfluenceCreatePageDetailsFormatter(page atlassian.AtlassianConfluencePageDetails) *AtlassianConfluencePageDetailsFormatter {
	return &AtlassianConfluencePageDetailsFormatter{
		page: page,
	}
}

// AtlassianConfluenceFormatSpacesAsMarkdown returns a raw markdown table
func (f *AtlassianConfluenceSpaceTableFormatter) AtlassianConfluenceFormatSpacesAsMarkdown() string {
	if len(f.spaces) == 0 {
		return "No spaces found."
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Write header
	fmt.Fprintf(w, "| Key | Name | Type | Status |\n")
	fmt.Fprintf(w, "|-----|------|------|--------|\n")

	// Write rows
	for _, space := range f.spaces {
		fmt.Fprintf(w, "| %s | %s | %s | %s |\n",
			space.Key,
			space.Name,
			strings.Title(strings.ToLower(space.Type)),
			strings.Title(strings.ToLower(space.Status)),
		)
	}

	w.Flush()
	return fmt.Sprintf("Found %d spaces:\n\n%s", len(f.spaces), buf.String())
}

// AtlassianConfluenceFormatSearchResultsAsMarkdown returns raw markdown formatted search results
func (f *AtlassianConfluenceSearchResultsFormatter) AtlassianConfluenceFormatSearchResultsAsMarkdown() string {
	if len(f.results) == 0 {
		return "No results found."
	}

	var md strings.Builder

	for i, result := range f.results {
		if i > 0 {
			md.WriteString("\n---\n\n")
		}

		// Title and metadata
		md.WriteString(fmt.Sprintf("Title: %s\n", AtlassianConfluenceCleanTitle(result.Title)))
		md.WriteString(fmt.Sprintf("Space: %s\n", result.ResultGlobalContainer.Title))
		md.WriteString(fmt.Sprintf("Status: %s\n", result.Content.Status))
		md.WriteString(fmt.Sprintf("Last Modified: %s\n", result.FriendlyLastModified))
		md.WriteString(fmt.Sprintf("URL: %s\n", AtlassianConfluenceFormatURL(result.URL)))
		md.WriteString("\n")

		// Content
		var description string
		if result.Content.Body.AtlasDocFormat.Value != "" {
			content, err := atlassian.AtlassianDocumentConvertJSONToMarkdown(result.Content.Body.AtlasDocFormat.Value)
			if err != nil {
				description = AtlassianConfluenceCleanContent(result.Excerpt)
			} else {
				description = AtlassianConfluenceCleanContent(content)
				if len(description) < 100 && result.Excerpt != "" {
					excerptDesc := AtlassianConfluenceCleanContent(result.Excerpt)
					if excerptDesc != "" {
						if description != "" {
							description += ". "
						}
						description += excerptDesc
					}
				}
			}
		} else {
			description = AtlassianConfluenceCleanContent(result.Excerpt)
		}

		if description == "" {
			description = "(No description available)"
		}

		description = util.TruncateText(description, 1000)
		md.WriteString(description)
		md.WriteString("\n")
	}

	return md.String()
}

// AtlassianConfluenceFormatPageDetailsAsMarkdown returns a raw markdown representation of the page details
func (f *AtlassianConfluencePageDetailsFormatter) AtlassianConfluenceFormatPageDetailsAsMarkdown() string {
	var output strings.Builder

	// Print page metadata
	output.WriteString(fmt.Sprintf("# %s\n\n", f.page.Title))
	output.WriteString("**Page Information**\n")
	output.WriteString(fmt.Sprintf("- **ID**: %s\n", f.page.ID))
	output.WriteString(fmt.Sprintf("- **Status**: %s\n", f.page.Status))
	output.WriteString(fmt.Sprintf("- **Version**: %d\n", f.page.Version.Number))
	if !f.page.Version.CreatedAt.IsZero() {
		output.WriteString(fmt.Sprintf("- **Last Modified**: %s\n", f.page.Version.CreatedAt.Format("Jan 02, 2006 15:04:05")))
	}
	if f.page.Version.Author.DisplayName != "" {
		output.WriteString(fmt.Sprintf("- **Author**: %s\n", f.page.Version.Author.DisplayName))
	}
	if f.page.SpaceId != "" {
		output.WriteString(fmt.Sprintf("- **Space ID**: %s\n", f.page.SpaceId))
	}
	if f.page.Links.WebUI != "" {
		output.WriteString(fmt.Sprintf("- **Web URL**: %s\n", f.page.Links.WebUI))
	}
	output.WriteString("\n---\n\n")

	// Convert body to markdown if available
	if f.page.Body.AtlasDocFormat.Value != "" {
		logging.LogDebug("Converting ADF content: %s", f.page.Body.AtlasDocFormat.Value)
		md, err := atlassian.AtlassianDocumentConvertJSONToMarkdown(f.page.Body.AtlasDocFormat.Value)
		if err != nil {
			output.WriteString(fmt.Sprintf("Error converting to markdown: %v\n", err))
			logging.LogDebug("ADF conversion error: %v", err)
		} else {
			output.WriteString(md)
			output.WriteString("\n")
		}
	} else {
		output.WriteString("*No content available*\n")
	}

	// Add comments if available
	if f.page.Comments != nil && len(f.page.Comments.Results) > 0 {
		output.WriteString("\n---\n\n## Comments\n\n")
		for _, comment := range f.page.Comments.Results {
			output.WriteString(fmt.Sprintf("### %s\n", comment.Title))
			if comment.Version.Author.DisplayName != "" {
				output.WriteString(fmt.Sprintf("**Author**: %s\n", comment.Version.Author.DisplayName))
			}
			if !comment.Version.CreatedAt.IsZero() {
				output.WriteString(fmt.Sprintf("**Last Modified**: %s\n\n", comment.Version.CreatedAt.Format("Jan 02, 2006 15:04:05")))
			}

			if comment.Body.AtlasDocFormat.Value != "" {
				logging.LogDebug("Converting comment ADF content: %s", comment.Body.AtlasDocFormat.Value)
				md, err := atlassian.AtlassianDocumentConvertJSONToMarkdown(comment.Body.AtlasDocFormat.Value)
				if err != nil {
					output.WriteString(fmt.Sprintf("Error converting comment to markdown: %v\n", err))
					logging.LogDebug("Comment ADF conversion error: %v", err)
				} else {
					output.WriteString(md)
					output.WriteString("\n")
				}
			}
			output.WriteString("\n---\n\n")
		}
	}

	return output.String()
}

// Helper functions

// AtlassianConfluenceCleanTitle replaces highlight markers with bold markers
func AtlassianConfluenceCleanTitle(title string) string {
	title = strings.ReplaceAll(title, "@@@hl@@@", "**")
	title = strings.ReplaceAll(title, "@@@endhl@@@", "**")
	return title
}

// AtlassianConfluenceFormatURL makes the URL more readable
func AtlassianConfluenceFormatURL(url string) string {
	url = strings.ReplaceAll(url, "+", " ")
	url = strings.TrimPrefix(url, "/spaces/")
	url = strings.TrimPrefix(url, "/wiki/spaces/")
	url = strings.TrimPrefix(url, "/pages/")
	return url
}

// AtlassianConfluenceCleanContent removes newlines, headings, and lists from the content
func AtlassianConfluenceCleanContent(content string) string {
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

// AtlassianConfluenceWrapText wraps text at the specified line length
func AtlassianConfluenceWrapText(text string, lineLength int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= lineLength {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return strings.Join(lines, "\n")
}
