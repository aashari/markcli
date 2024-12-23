package markdown

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"markcli/internal/api/atlassian"
)

// SpaceTableFormatter formats Confluence spaces as a markdown table
type SpaceTableFormatter struct {
	spaces []atlassian.Space
}

// SearchResultsFormatter formats search results into markdown
type SearchResultsFormatter struct {
	results []atlassian.ContentResult
}

// PageDetailsResponse represents the response from the Confluence API v2 for a single page
type PageDetailsResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Version struct {
		Number    int       `json:"number"`
		Message   string    `json:"message"`
		CreatedAt time.Time `json:"createdAt"`
		Author    struct {
			AccountID   string `json:"accountId"`
			Email       string `json:"email"`
			DisplayName string `json:"displayName"`
		} `json:"author"`
	} `json:"version"`
	Body struct {
		AtlasDocFormat struct {
			Value string `json:"value"`
		} `json:"atlas_doc_format"`
	} `json:"body"`
	SpaceId string `json:"spaceId"`
	Links   struct {
		WebUI string `json:"webui"`
	} `json:"_links"`
	Comments *atlassian.FooterCommentsResponse `json:"-"`
}

// PageDetailsFormatter formats a single Confluence page's details
type PageDetailsFormatter struct {
	page PageDetailsResponse
}

// JiraProjectTableFormatter formats Jira projects as a markdown table
type JiraProjectTableFormatter struct {
	projects []atlassian.JiraProject
	sortBy   string // Can be "key", "name", "type", or "style"
}

// JiraSearchResultsFormatter formats Jira issue search results as Markdown.
type JiraSearchResultsFormatter struct {
	issues []atlassian.JiraIssue
}

// JiraIssueDetailsFormatter formats a single Jira issue's details
type JiraIssueDetailsFormatter struct {
	issue atlassian.JiraIssueDetails
}

// NewSpaceTableFormatter creates a new SpaceTableFormatter
func NewSpaceTableFormatter(spaces []atlassian.Space) *SpaceTableFormatter {
	return &SpaceTableFormatter{
		spaces: spaces,
	}
}

// NewSearchResultsFormatter creates a new search results formatter
func NewSearchResultsFormatter(results []atlassian.ContentResult) *SearchResultsFormatter {
	return &SearchResultsFormatter{
		results: results,
	}
}

// NewPageDetailsFormatter creates a new PageDetailsFormatter
func NewPageDetailsFormatter(page PageDetailsResponse) *PageDetailsFormatter {
	return &PageDetailsFormatter{
		page: page,
	}
}

// NewJiraProjectTableFormatter creates a new JiraProjectTableFormatter
func NewJiraProjectTableFormatter(projects []atlassian.JiraProject, sortBy string) *JiraProjectTableFormatter {
	// Default to sorting by key if no sort field is specified
	if sortBy == "" {
		sortBy = "key"
	}
	return &JiraProjectTableFormatter{
		projects: projects,
		sortBy:   strings.ToLower(sortBy),
	}
}

// NewJiraSearchResultsFormatter creates a new JiraSearchResultsFormatter
func NewJiraSearchResultsFormatter(issues []atlassian.JiraIssue) *JiraSearchResultsFormatter {
	return &JiraSearchResultsFormatter{
		issues: issues,
	}
}

// NewJiraIssueDetailsFormatter creates a new JiraIssueDetailsFormatter
func NewJiraIssueDetailsFormatter(issue atlassian.JiraIssueDetails) *JiraIssueDetailsFormatter {
	return &JiraIssueDetailsFormatter{
		issue: issue,
	}
}

// RawMarkdown returns a raw markdown table
func (f *SpaceTableFormatter) RawMarkdown() string {
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

// RawMarkdown returns raw markdown formatted search results
func (f *SearchResultsFormatter) RawMarkdown() string {
	if len(f.results) == 0 {
		return "No results found."
	}

	var md strings.Builder

	for i, result := range f.results {
		if i > 0 {
			md.WriteString("\n---\n\n")
		}

		// Title and metadata
		md.WriteString(fmt.Sprintf("Title: %s\n", cleanTitle(result.Title)))
		md.WriteString(fmt.Sprintf("Space: %s\n", result.ResultGlobalContainer.Title))
		md.WriteString(fmt.Sprintf("Status: %s\n", result.Content.Status))
		md.WriteString(fmt.Sprintf("Last Modified: %s\n", result.FriendlyLastModified))
		md.WriteString(fmt.Sprintf("URL: %s\n", formatURL(result.URL)))
		md.WriteString("\n")

		// Content
		var description string
		if result.Content.Body.AtlasDocFormat.Value != "" {
			content, err := ConvertAtlasToMarkdown(result.Content.Body.AtlasDocFormat.Value)
			if err != nil {
				description = cleanContent(result.Excerpt)
			} else {
				// Try to get more content from the atlas doc format
				description = cleanContent(content)
				if len(description) < 100 && result.Excerpt != "" {
					// If atlas doc content is too short, append the excerpt
					excerptDesc := cleanContent(result.Excerpt)
					if excerptDesc != "" {
						if description != "" {
							description += ". "
						}
						description += excerptDesc
					}
				}
			}
		} else {
			description = cleanContent(result.Excerpt)
		}

		// Skip if no meaningful content
		if description == "" {
			// Try to get content from the excerpt even if it's not in atlas doc format
			if result.Excerpt != "" {
				description = cleanContent(result.Excerpt)
			}
			if description == "" {
				description = "(No description available)"
			}
		}

		// Truncate to 1500 characters and wrap at 100 characters
		description = truncateText(description, 1500)
		// Wrap text at word boundaries
		words := strings.Fields(description)
		var wrapped strings.Builder
		lineLength := 0
		for j, word := range words {
			if lineLength+len(word)+1 > 100 {
				wrapped.WriteString("\n")
				lineLength = 0
			} else if j > 0 {
				wrapped.WriteString(" ")
				lineLength++
			}
			wrapped.WriteString(word)
			lineLength += len(word)
		}
		md.WriteString(wrapped.String())
		md.WriteString("\n")
	}

	return md.String()
}

// RawMarkdown returns a raw markdown representation of the page details
func (f *PageDetailsFormatter) RawMarkdown() string {
	var output strings.Builder

	// Print page metadata
	output.WriteString(fmt.Sprintf("# %s\n\n", f.page.Title))
	output.WriteString("**Page Information**\n")
	output.WriteString(fmt.Sprintf("- **ID**: %s\n", f.page.ID))
	output.WriteString(fmt.Sprintf("- **Status**: %s\n", f.page.Status))
	output.WriteString(fmt.Sprintf("- **Version**: %d\n", f.page.Version.Number))
	output.WriteString(fmt.Sprintf("- **Last Modified**: %s\n", f.page.Version.CreatedAt.Format("Jan 02, 2006 15:04:05 MST")))
	output.WriteString(fmt.Sprintf("- **Last Editor**: %s (%s)\n", f.page.Version.Author.DisplayName, f.page.Version.Author.Email))
	output.WriteString(fmt.Sprintf("- **Space ID**: %s\n", f.page.SpaceId))
	output.WriteString(fmt.Sprintf("- **Web URL**: %s\n", f.page.Links.WebUI))
	output.WriteString("\n---\n\n")

	// Convert atlas_doc_format to markdown
	md, err := ConvertAtlasToMarkdown(f.page.Body.AtlasDocFormat.Value)
	if err != nil {
		if strings.Contains(err.Error(), "unsupported content type") {
			// Handle unsupported content type by showing a message and the raw content
			output.WriteString(fmt.Sprintf("Note: Some content types in this page are not yet supported. Raw content shown below:\n\n%s\n",
				f.page.Body.AtlasDocFormat))
		} else {
			output.WriteString(fmt.Sprintf("Error converting to markdown: %v\n", err))
		}
	} else {
		output.WriteString(md)
	}

	// Add comments section if available
	if f.page.Comments != nil && len(f.page.Comments.Results) > 0 {
		output.WriteString("\n---\n\n## Comments\n\n")
		for _, comment := range f.page.Comments.Results {
			// Add comment metadata
			output.WriteString(fmt.Sprintf("### Comment by %s\n", comment.Version.AuthorID))
			output.WriteString(fmt.Sprintf("*Posted on %s*\n\n", comment.Version.CreatedAt))

			// Try atlas_doc_format first
			if comment.Body.AtlasDocFormat.Value != "" {
				commentMd, err := ConvertAtlasToMarkdown(comment.Body.AtlasDocFormat.Value)
				if err != nil {
					if strings.Contains(err.Error(), "unsupported content type") {
						output.WriteString("Note: Some content types in this comment are not supported. Raw content shown below:\n\n")
						output.WriteString(comment.Body.AtlasDocFormat.Value)
					} else {
						output.WriteString(fmt.Sprintf("Error converting comment to markdown: %v\n", err))
					}
				} else {
					output.WriteString(commentMd)
				}
			} else if comment.Body.Storage.Value != "" {
				// Fallback to storage format if atlas_doc_format is not available
				// Storage format is typically HTML, so we'll display it with a note
				output.WriteString("Note: This comment is in legacy format:\n\n")
				output.WriteString(comment.Body.Storage.Value)
			} else {
				output.WriteString("*[Empty comment]*\n")
			}
			output.WriteString("\n---\n\n")
		}
	}

	return output.String()
}

// RawMarkdown returns a raw markdown table
func (f *JiraProjectTableFormatter) RawMarkdown() string {
	if len(f.projects) == 0 {
		return "No projects found."
	}

	// Sort projects based on the specified field
	sort.Slice(f.projects, func(i, j int) bool {
		switch f.sortBy {
		case "name":
			return f.projects[i].Name < f.projects[j].Name
		case "type":
			return f.projects[i].ProjectTypeKey < f.projects[j].ProjectTypeKey
		case "style":
			return f.projects[i].Style < f.projects[j].Style
		default: // "key" is the default sort field
			return f.projects[i].Key < f.projects[j].Key
		}
	})

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Write header
	fmt.Fprintf(w, "| Key | Name | Type | Style |\n")
	fmt.Fprintf(w, "|-----|------|------|-------|\n")

	// Write rows
	for _, project := range f.projects {
		name := truncateText(project.Name, 50)
		fmt.Fprintf(w, "| %s | %s | %s | %s |\n",
			project.Key,
			name,
			strings.Title(strings.ToLower(project.ProjectTypeKey)),
			strings.Title(strings.ToLower(project.Style)),
		)
	}

	w.Flush()
	return fmt.Sprintf("Found %d projects:\n\n%s", len(f.projects), buf.String())
}

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
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, ">") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "#") {
			// Don't wrap code blocks, blockquotes, lists, or headings
			wrapped.WriteString(line)
		} else {
			words := strings.Fields(line)
			lineLen := 0

			for j, word := range words {
				wordLen := len(word)
				if lineLen+wordLen+1 > lineLength && lineLen > 0 {
					wrapped.WriteString("\n")
					lineLen = 0
				} else if j > 0 {
					wrapped.WriteString(" ")
					lineLen++
				}
				wrapped.WriteString(word)
				lineLen += wordLen
			}
		}

		if i < len(lines)-1 {
			wrapped.WriteString("\n")
		}
	}

	return wrapped.String()
}

// formatIssueURL formats a Jira API URL to a web UI URL
func formatIssueURL(apiURL string) string {
	// Parse the API URL
	u, err := url.Parse(apiURL)
	if err != nil {
		// If URL parsing fails, fall back to string replacement
		return strings.Replace(apiURL, "rest/api/3/issue", "browse", 1)
	}

	// Get the issue key from the path
	parts := strings.Split(u.Path, "/")
	var issueKey string
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			issueKey = parts[i]
			break
		}
	}

	// Create the web UI URL
	webUI := &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   path.Join("browse", issueKey),
	}

	return webUI.String()
}

// RawMarkdown implements the formatting logic for Jira search results
func (f *JiraSearchResultsFormatter) RawMarkdown() string {
	if len(f.issues) == 0 {
		return "No issues found."
	}

	var md strings.Builder

	for i, issue := range f.issues {
		if i > 0 {
			md.WriteString("\n---\n\n")
		}

		// Title and metadata
		summary := truncateText(issue.Fields.Summary, 100)
		md.WriteString(fmt.Sprintf("### [%s] %s\n", issue.Key, summary))

		projectInfo := fmt.Sprintf("%s (%s)", issue.Fields.Project.Name, issue.Fields.Project.Key)
		projectInfo = truncateText(projectInfo, 50)
		md.WriteString(fmt.Sprintf("**Project**: %s\n", projectInfo))

		md.WriteString(fmt.Sprintf("**Status**: %s\n", issue.Fields.Status.Name))

		if issue.Fields.Assignee != nil && issue.Fields.Assignee.DisplayName != "" {
			assignee := truncateText(issue.Fields.Assignee.DisplayName, 30)
			md.WriteString(fmt.Sprintf("**Assignee**: %s\n", assignee))
		} else {
			md.WriteString("**Assignee**: Unassigned\n")
		}

		// Format the URL
		webURL := formatIssueURL(issue.Self)
		md.WriteString(fmt.Sprintf("**URL**: %s\n", webURL))

		// Format and add the description
		if issue.Fields.Description != nil {
			// Convert description to JSON for atlas_doc parsing
			descJSON, err := json.Marshal(issue.Fields.Description)
			if err != nil {
				md.WriteString("\n**Description**: _(Error formatting description)_\n")
				continue
			}

			desc, err := ConvertAtlasToMarkdown(string(descJSON))
			if err != nil {
				md.WriteString("\n**Description**: _(Error parsing description)_\n")
				continue
			}

			// Clean, truncate and wrap the description
			desc = cleanContent(desc)
			if desc != "" {
				desc = truncateText(desc, 500)
				desc = wrapText(desc, 100)
				md.WriteString("\n**Description**:\n")
				md.WriteString(desc)
				md.WriteString("\n")
			}
		} else {
			md.WriteString("\n**Description**: _(No description provided)_\n")
		}
	}

	return md.String()
}

// cleanContent removes newlines, headings, and lists from the content
func cleanContent(content string) string {
	// First, replace highlight markers with bold markers
	content = strings.ReplaceAll(content, "@@@hl@@@", "**")
	content = strings.ReplaceAll(content, "@@@endhl@@@", "**")

	// Split into lines
	lines := strings.Split(content, "\n")
	var cleaned []string
	inCodeBlock := false
	inTable := false
	var tableHeaders []string
	var tableRows [][]string
	var currentRow []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Handle code blocks
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			cleaned = append(cleaned, line)
			continue
		}
		if inCodeBlock {
			cleaned = append(cleaned, line)
			continue
		}

		// Handle table rows
		if strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|") {
			// Skip table separator lines
			if strings.Contains(line, "---") {
				continue
			}

			// Extract cell content
			cells := strings.Split(line, "|")
			currentRow = nil
			for _, cell := range cells {
				cell = strings.TrimSpace(cell)
				if cell != "" {
					currentRow = append(currentRow, cell)
				}
			}

			if !inTable {
				// This is the header row
				inTable = true
				tableHeaders = currentRow
			} else {
				// This is a data row
				if len(currentRow) > 0 {
					tableRows = append(tableRows, currentRow)
				}
			}
			continue
		} else if inTable {
			// Table ended
			inTable = false
			if len(tableHeaders) > 0 && len(tableRows) > 0 {
				// Process table data
				var descriptions []string
				for _, row := range tableRows {
					if len(row) > 0 {
						// Get the first few columns that have meaningful content
						var rowDesc []string
						for i, cell := range row {
							if i >= len(tableHeaders) {
								break
							}
							if len(cell) > 10 && !strings.Contains(cell, "http") {
								rowDesc = append(rowDesc, fmt.Sprintf("%s: %s", tableHeaders[i], cell))
							}
						}
						if len(rowDesc) > 0 {
							descriptions = append(descriptions, strings.Join(rowDesc[:1], ", "))
						}
					}
				}
				if len(descriptions) > 0 {
					// Take only the first few meaningful entries
					if len(descriptions) > 3 {
						descriptions = descriptions[:3]
					}
					cleaned = append(cleaned, strings.Join(descriptions, ". "))
				}
			}
			tableHeaders = nil
			tableRows = nil
		}

		// Skip empty lines and special lines
		if line == "" || line == "---" || strings.Contains(line, "${") {
			continue
		}

		// Preserve panels
		if strings.HasPrefix(line, ">") {
			cleaned = append(cleaned, line)
			continue
		}

		// Preserve headings and lists if they have content
		if (strings.HasPrefix(line, "#") || strings.HasPrefix(line, "*") ||
			strings.HasPrefix(line, "-")) && len(line) > 2 {
			cleaned = append(cleaned, line)
			continue
		}

		// Remove HTML entities
		line = strings.ReplaceAll(line, "&#39;", "'")
		line = strings.ReplaceAll(line, "&quot;", "\"")
		line = strings.ReplaceAll(line, "&lt;", "<")
		line = strings.ReplaceAll(line, "&gt;", ">")
		line = strings.ReplaceAll(line, "&amp;", "&")

		// Skip if line is empty after cleaning or too short
		if strings.TrimSpace(line) == "" || len(strings.TrimSpace(line)) < 3 {
			continue
		}

		cleaned = append(cleaned, line)
	}

	// Join lines and clean up extra spaces
	result := strings.Join(cleaned, "\n")
	result = strings.ReplaceAll(result, "  ", " ")
	result = strings.ReplaceAll(result, "..", ".")
	result = strings.ReplaceAll(result, ". .", ".")
	result = strings.TrimSpace(result)

	return result
}

// formatURL makes the URL more readable
func formatURL(url string) string {
	// Remove URL encoding
	url = strings.ReplaceAll(url, "+", " ")
	// Remove common prefixes
	url = strings.TrimPrefix(url, "/spaces/")
	url = strings.TrimPrefix(url, "/wiki/spaces/")
	url = strings.TrimPrefix(url, "/pages/")
	return url
}

// cleanTitle replaces highlight markers with bold markers
func cleanTitle(title string) string {
	title = strings.ReplaceAll(title, "@@@hl@@@", "**")
	title = strings.ReplaceAll(title, "@@@endhl@@@", "**")
	return title
}

// RawMarkdown returns a raw markdown representation of the Jira issue details
func (f *JiraIssueDetailsFormatter) RawMarkdown() string {
	var output strings.Builder
	issue := f.issue

	// Print issue metadata
	output.WriteString(fmt.Sprintf("# %s\n\n", issue.Fields.Summary))
	output.WriteString("**Issue Information**\n")
	output.WriteString(fmt.Sprintf("- **ID**: %s\n", issue.ID))
	output.WriteString(fmt.Sprintf("- **Key**: %s\n", issue.Key))
	output.WriteString(fmt.Sprintf("- **Project**: %s (%s)\n", issue.Fields.Project.Name, issue.Fields.Project.Key))

	if issue.Fields.Assignee != nil && issue.Fields.Assignee.DisplayName != "" {
		output.WriteString(fmt.Sprintf("- **Assignee**: %s\n", issue.Fields.Assignee.DisplayName))
	} else {
		output.WriteString("- **Assignee**: Unassigned\n")
	}
	if issue.Fields.Reporter != nil && issue.Fields.Reporter.DisplayName != "" {
		output.WriteString(fmt.Sprintf("- **Reporter**: %s\n", issue.Fields.Reporter.DisplayName))
	}
	output.WriteString(fmt.Sprintf("- **Status**: %s\n", issue.Fields.Status.Name))

	// Format dates in a more readable format
	if created, err := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.Created); err == nil {
		output.WriteString(fmt.Sprintf("- **Created**: %s\n", created.Format("Jan 02, 2006 15:04:05 MST")))
	} else {
		output.WriteString(fmt.Sprintf("- **Created**: %s\n", issue.Fields.Created))
	}
	if updated, err := time.Parse("2006-01-02T15:04:05.999-0700", issue.Fields.Updated); err == nil {
		output.WriteString(fmt.Sprintf("- **Updated**: %s\n", updated.Format("Jan 02, 2006 15:04:05 MST")))
	} else {
		output.WriteString(fmt.Sprintf("- **Updated**: %s\n", issue.Fields.Updated))
	}

	output.WriteString(fmt.Sprintf("- **URL**: %s\n", formatIssueURL(issue.Self)))

	output.WriteString("\n---\n\n")

	// Format and add the description using atlas_doc
	if issue.Fields.Description != nil {
		descJSON, err := json.Marshal(issue.Fields.Description)
		if err != nil {
			output.WriteString(fmt.Sprintf("Error formatting description: %v\n", err))
		} else {
			desc, err := ConvertJiraAtlasToMarkdown(string(descJSON))
			if err != nil {
				if strings.Contains(err.Error(), "unsupported content type") {
					// Handle unsupported content type by showing a message and the raw content
					output.WriteString(fmt.Sprintf("Note: Some content types in this issue are not yet supported. Raw content shown below:\n\n%s\n",
						string(descJSON)))
				} else {
					output.WriteString(fmt.Sprintf("Error converting to markdown: %v\n", err))
				}
			} else {
				// Format the description
				if desc != "" {
					// Clean up extra newlines and format the text
					desc = strings.TrimSpace(desc)
					desc = strings.ReplaceAll(desc, "\n\n\n", "\n\n")
					// Wrap text at word boundaries for better readability
					desc = wrapText(desc, 100)
					output.WriteString(desc)
					output.WriteString("\n\n")
				} else {
					output.WriteString("_(No description provided)_\n\n")
				}
			}
		}
	} else {
		output.WriteString("_(No description provided)_\n\n")
	}

	return output.String()
}
