package formatting

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"markcli/internal/types/atlassian"
	"markcli/internal/util"
)

// AtlassianJiraProjectTableFormatter formats Jira projects as a markdown table
type AtlassianJiraProjectTableFormatter struct {
	projects []atlassian.AtlassianJiraProject
	sortBy   string // Can be "key", "name", "type", or "style"
}

// AtlassianJiraSearchResultsFormatter formats Jira issue search results as Markdown
type AtlassianJiraSearchResultsFormatter struct {
	issues []atlassian.AtlassianJiraIssue
}

// AtlassianJiraIssueDetailsFormatter formats a single Jira issue's details
type AtlassianJiraIssueDetailsFormatter struct {
	issue    atlassian.AtlassianJiraIssue
	comments *atlassian.AtlassianJiraCommentsResponse
}

// AtlassianJiraCreateProjectTableFormatter creates a new project table formatter
func AtlassianJiraCreateProjectTableFormatter(projects []atlassian.AtlassianJiraProject, sortBy string) *AtlassianJiraProjectTableFormatter {
	if sortBy == "" {
		sortBy = "key"
	}
	return &AtlassianJiraProjectTableFormatter{
		projects: projects,
		sortBy:   strings.ToLower(sortBy),
	}
}

// AtlassianJiraCreateSearchResultsFormatter creates a new search results formatter
func AtlassianJiraCreateSearchResultsFormatter(issues []atlassian.AtlassianJiraIssue) *AtlassianJiraSearchResultsFormatter {
	return &AtlassianJiraSearchResultsFormatter{
		issues: issues,
	}
}

// AtlassianJiraCreateIssueDetailsFormatter creates a new issue details formatter
func AtlassianJiraCreateIssueDetailsFormatter(issue atlassian.AtlassianJiraIssue) *AtlassianJiraIssueDetailsFormatter {
	return &AtlassianJiraIssueDetailsFormatter{
		issue: issue,
	}
}

// WithComments adds comments to the formatter
func (f *AtlassianJiraIssueDetailsFormatter) WithComments(comments *atlassian.AtlassianJiraCommentsResponse) *AtlassianJiraIssueDetailsFormatter {
	f.comments = comments
	return f
}

// AtlassianJiraFormatProjectsAsMarkdown returns a markdown table of projects
func (f *AtlassianJiraProjectTableFormatter) AtlassianJiraFormatProjectsAsMarkdown() string {
	if len(f.projects) == 0 {
		return "No projects found."
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Write header
	fmt.Fprintf(w, "| Key | Name | Type | Style |\n")
	fmt.Fprintf(w, "|-----|------|------|-------|\n")

	// Write rows
	for _, project := range f.projects {
		name := util.TruncateText(project.Name, 50)
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

// AtlassianJiraFormatSearchResultsAsMarkdown returns search results in markdown format
func (f *AtlassianJiraSearchResultsFormatter) AtlassianJiraFormatSearchResultsAsMarkdown() string {
	if len(f.issues) == 0 {
		return "No issues found."
	}

	var output strings.Builder

	for i, issue := range f.issues {
		if i > 0 {
			output.WriteString("\n---\n\n")
		}

		// Title and key
		output.WriteString(fmt.Sprintf("Title: %s\n", issue.Fields.Summary))
		output.WriteString(fmt.Sprintf("Key: %s\n", issue.Key))
		output.WriteString(fmt.Sprintf("Project: %s\n", issue.Fields.Project.Name))
		output.WriteString(fmt.Sprintf("Status: %s\n", issue.Fields.Status.Name))
		output.WriteString(fmt.Sprintf("Priority: %s\n", issue.Fields.Priority.Name))

		// Add assignee if available
		if issue.Fields.Assignee != nil {
			output.WriteString(fmt.Sprintf("Assignee: %s\n", issue.Fields.Assignee.DisplayName))
		}

		// Add last modified date
		if issue.Fields.Updated != "" {
			t, err := util.ParseDate(issue.Fields.Updated)
			if err == nil {
				output.WriteString(fmt.Sprintf("Last Modified: %s\n", t.Format("Jan 02, 2006")))
			} else {
				output.WriteString(fmt.Sprintf("Last Modified: %s\n", issue.Fields.Updated))
			}
		}

		// Add URL by constructing from Self link
		if issue.Self != "" {
			webURL := strings.Replace(issue.Self, "/rest/api/3/issue/", "/browse/", 1)
			output.WriteString(fmt.Sprintf("URL: %s\n", webURL))
		}

		// Add description if available (truncated for consistency with Confluence excerpts)
		if issue.Fields.Description != nil && len(issue.Fields.Description.Content) > 0 {
			output.WriteString("\n")
			doc := &atlassian.AtlassianDocument{
				Type:    "doc",
				Content: issue.Fields.Description.Content,
				Version: issue.Fields.Description.Version,
			}
			if desc, err := doc.AtlassianDocumentConvertToMarkdown(); err == nil {
				desc = util.TruncateText(desc, 300)
				output.WriteString(desc)
			}
			output.WriteString("\n")
		}
	}

	return output.String()
}

// AtlassianJiraFormatIssueDetailsAsMarkdown returns issue details in markdown format
func (f *AtlassianJiraIssueDetailsFormatter) AtlassianJiraFormatIssueDetailsAsMarkdown() string {
	issue := f.issue
	var output strings.Builder

	// Title
	output.WriteString(fmt.Sprintf("# %s\n\n", issue.Fields.Summary))

	// Metadata
	output.WriteString("**Issue Information**\n")
	output.WriteString(fmt.Sprintf("- **Key**: %s\n", issue.Key))
	output.WriteString(fmt.Sprintf("- **Type**: %s\n", issue.Fields.IssueType.Name))
	output.WriteString(fmt.Sprintf("- **Status**: %s\n", issue.Fields.Status.Name))
	output.WriteString(fmt.Sprintf("- **Priority**: %s\n", issue.Fields.Priority.Name))
	output.WriteString(fmt.Sprintf("- **Project**: %s (%s)\n", issue.Fields.Project.Name, issue.Fields.Project.Key))

	if issue.Fields.Assignee != nil {
		output.WriteString(fmt.Sprintf("- **Assignee**: %s\n", issue.Fields.Assignee.DisplayName))
	} else {
		output.WriteString("- **Assignee**: Unassigned\n")
	}

	if issue.Fields.Reporter != nil {
		output.WriteString(fmt.Sprintf("- **Reporter**: %s\n", issue.Fields.Reporter.DisplayName))
	}

	if issue.Fields.Created != "" {
		t, err := util.ParseDate(issue.Fields.Created)
		if err == nil {
			output.WriteString(fmt.Sprintf("- **Created**: %s\n", t.Format("Jan 02, 2006 15:04:05")))
		} else {
			output.WriteString(fmt.Sprintf("- **Created**: %s\n", issue.Fields.Created))
		}
	}

	if issue.Fields.Updated != "" {
		t, err := util.ParseDate(issue.Fields.Updated)
		if err == nil {
			output.WriteString(fmt.Sprintf("- **Last Modified**: %s\n", t.Format("Jan 02, 2006 15:04:05")))
		} else {
			output.WriteString(fmt.Sprintf("- **Last Modified**: %s\n", issue.Fields.Updated))
		}
	}

	if issue.Fields.Resolution.Name != "" {
		output.WriteString(fmt.Sprintf("- **Resolution**: %s\n", issue.Fields.Resolution.Name))
	}

	// Add URL if available
	if issue.Self != "" {
		webURL := strings.Replace(issue.Self, "/rest/api/3/issue/", "/browse/", 1)
		output.WriteString(fmt.Sprintf("- **Web URL**: %s\n", webURL))
	}

	output.WriteString("\n---\n\n")

	// Description
	if issue.Fields.Description != nil && len(issue.Fields.Description.Content) > 0 {
		doc := &atlassian.AtlassianDocument{
			Type:    "doc",
			Content: issue.Fields.Description.Content,
			Version: issue.Fields.Description.Version,
		}
		if desc, err := doc.AtlassianDocumentConvertToMarkdown(); err == nil {
			output.WriteString("## Description\n\n")
			output.WriteString(desc)
			output.WriteString("\n\n")
		}
	}

	// Comments
	if f.comments != nil && len(f.comments.Comments) > 0 {
		output.WriteString("## Comments\n\n")
		for _, comment := range f.comments.Comments {
			if comment.Author != nil {
				output.WriteString(fmt.Sprintf("**%s** ", comment.Author.DisplayName))
			}
			if comment.Created != "" {
				t, err := util.ParseDate(comment.Created)
				if err == nil {
					output.WriteString(fmt.Sprintf("on %s", t.Format("Jan 02, 2006 15:04:05")))
				} else {
					output.WriteString(fmt.Sprintf("on %s", comment.Created))
				}
			}
			output.WriteString("\n\n")

			if comment.Body != nil {
				doc := &atlassian.AtlassianDocument{
					Type:    "doc",
					Content: comment.Body.Content,
					Version: comment.Body.Version,
				}
				if body, err := doc.AtlassianDocumentConvertToMarkdown(); err == nil {
					output.WriteString(body)
					output.WriteString("\n\n")
				}
			}

			output.WriteString("---\n\n")
		}
	} else {
		output.WriteString("## Comments\n\nNo comments found\n")
	}

	return output.String()
}
