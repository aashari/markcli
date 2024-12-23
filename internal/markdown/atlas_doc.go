package markdown

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// AtlasDocument represents the root structure of an atlas_doc_format document
type AtlasDocument struct {
	Type    string         `json:"type"`
	Content []AtlasContent `json:"content"`
	Version int            `json:"version"`
}

// AtlasContent represents a content node in the atlas document
type AtlasContent struct {
	Type    string          `json:"type"`
	Content []AtlasContent  `json:"content,omitempty"`
	Text    string          `json:"text,omitempty"`
	Marks   []AtlasMark     `json:"marks,omitempty"`
	Attrs   AtlasAttributes `json:"attrs,omitempty"`
}

// AtlasMark represents text formatting in the atlas document
type AtlasMark struct {
	Type  string `json:"type"`
	Attrs struct {
		URL   string `json:"url,omitempty"`
		Color string `json:"color,omitempty"`
	} `json:"attrs,omitempty"`
}

// AtlasAttrs represents the attributes of an Atlas content node
type AtlasAttrs struct {
	Level              int                    `json:"level,omitempty"`
	Text               string                 `json:"text,omitempty"`
	Title              string                 `json:"title,omitempty"`
	ExtensionType      string                 `json:"extensionType,omitempty"`
	ExtensionKey       string                 `json:"extensionKey,omitempty"`
	Parameters         json.RawMessage        `json:"parameters,omitempty"`
	URL                string                 `json:"url,omitempty"`
	ReferencePageID    string                 `json:"referencePageId,omitempty"`
	ReferenceStatus    string                 `json:"referenceStatus,omitempty"`
	ReferencePageTitle string                 `json:"referencePageTitle,omitempty"`
	Color              string                 `json:"color,omitempty"`
	PanelType          string                 `json:"panelType,omitempty"`
	State              string                 `json:"state,omitempty"`
	Language           string                 `json:"language,omitempty"`
	Timestamp          string                 `json:"timestamp,omitempty"`
	MediaWidth         int                    `json:"mediaWidth,omitempty"`
	MediaHeight        int                    `json:"mediaHeight,omitempty"`
	Other              map[string]interface{} `json:"-"`
}

// Update the existing type definitions
type MacroValue struct {
	Value interface{} `json:"value,omitempty"`
}

type MacroParams struct {
	Root       *MacroValue            `json:"root,omitempty"`
	Spaces     *MacroValue            `json:"spaces,omitempty"`
	StartDepth *MacroValue            `json:"startDepth,omitempty"`
	Excerpt    *MacroValue            `json:"excerpt,omitempty"`
	Other      map[string]interface{} `json:"-"`
}

type MacroMetadata struct {
	MacroId       *MacroValue            `json:"macroId,omitempty"`
	SchemaVersion *MacroValue            `json:"schemaVersion,omitempty"`
	Title         *MacroValue            `json:"title,omitempty"`
	Other         map[string]interface{} `json:"-"`
}

type Parameters struct {
	MacroParams   *MacroParams           `json:"macroParams,omitempty"`
	MacroMetadata *MacroMetadata         `json:"macroMetadata,omitempty"`
	Other         map[string]interface{} `json:"-"`
}

type AtlasAttributes struct {
	Level              int                    `json:"level,omitempty"`
	Text               string                 `json:"text,omitempty"`
	Title              string                 `json:"title,omitempty"`
	ExtensionType      string                 `json:"extensionType,omitempty"`
	ExtensionKey       string                 `json:"extensionKey,omitempty"`
	Parameters         json.RawMessage        `json:"parameters,omitempty"`
	URL                string                 `json:"url,omitempty"`
	ReferencePageID    string                 `json:"referencePageId,omitempty"`
	ReferenceStatus    string                 `json:"referenceStatus,omitempty"`
	ReferencePageTitle string                 `json:"referencePageTitle,omitempty"`
	Color              string                 `json:"color,omitempty"`
	PanelType          string                 `json:"panelType,omitempty"`
	State              string                 `json:"state,omitempty"`
	Language           string                 `json:"language,omitempty"`
	Timestamp          string                 `json:"timestamp,omitempty"`
	MediaWidth         int                    `json:"mediaWidth,omitempty"`
	MediaHeight        int                    `json:"mediaHeight,omitempty"`
	Other              map[string]interface{} `json:"-"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (a *AtlasAttributes) UnmarshalJSON(data []byte) error {
	// First unmarshal into a map to get all fields
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Check for panelType in the raw data
	if panelType, ok := raw["panelType"].(string); ok {
		a.PanelType = panelType
	}

	// Create a temporary struct to unmarshal the known fields
	type Alias AtlasAttributes
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	// Unmarshal into the temporary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Store any unhandled fields in Other
	a.Other = make(map[string]interface{})
	for k, v := range raw {
		switch k {
		case "level", "text", "title", "extensionType", "extensionKey", "parameters",
			"url", "referencePageId", "referenceStatus", "referencePageTitle", "color",
			"panelType", "state", "language", "timestamp", "mediaWidth", "mediaHeight":
			continue
		default:
			a.Other[k] = v
		}
	}

	return nil
}

// ConvertAtlasToMarkdown converts atlas_doc_format JSON string to Markdown
func ConvertAtlasToMarkdown(jsonStr string) (string, error) {
	var doc struct {
		Type    string         `json:"type"`
		Content []AtlasContent `json:"content"`
		Version int            `json:"version"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &doc); err != nil {
		fmt.Printf("Error unmarshaling atlas document: %v\n", err)
		return "", fmt.Errorf("failed to parse atlas document: %v", err)
	}

	var result strings.Builder
	for _, content := range doc.Content {
		text, err := convertContent(&content)
		if err != nil {
			fmt.Printf("Error converting content: %v\n", err)
			return "", err
		}
		result.WriteString(text)
	}

	return result.String(), nil
}

// convertContent converts a single content node to Markdown
func convertContent(content *AtlasContent) (string, error) {
	switch content.Type {
	case "paragraph":
		return convertParagraph(content)
	case "text":
		return convertText(content)
	case "table":
		return convertTable(content)
	case "tableRow":
		return convertTableRow(content)
	case "tableHeader", "tableCell":
		return convertTableCell(content)
	case "bulletList":
		return convertBulletList(content)
	case "orderedList":
		return convertOrderedList(content)
	case "listItem":
		return convertListItem(content)
	case "inlineCard":
		return convertInlineCard(content)
	case "status":
		return convertStatus(content)
	case "heading":
		return convertHeading(content)
	case "emoji":
		return convertEmoji(content)
	case "panel":
		return convertPanel(content)
	case "taskList":
		return convertTaskList(content)
	case "taskItem":
		return convertTaskItem(content)
	case "rule":
		return "---\n\n", nil
	case "codeBlock":
		return convertCodeBlock(content)
	case "hardBreak":
		return "\n", nil
	case "bodiedExtension":
		return convertBodiedExtension(content)
	case "extension":
		return convertExtension(content)
	case "date":
		return convertDate(content)
	case "mention":
		return convertMention(content)
	case "placeholder":
		return convertPlaceholder(content)
	case "layoutSection":
		return convertLayoutSection(content)
	case "layoutColumn":
		return convertLayoutColumn(content)
	case "mediaSingle":
		return convertMediaSingle(content)
	case "media":
		return convertMedia(content)
	case "nestedExpand":
		return convertNestedExpand(content)
	case "expand":
		return convertExpand(content)
	case "inlineExtension":
		return convertInlineExtension(content)
	default:
		return "", fmt.Errorf("unsupported content type: %s", content.Type)
	}
}

func convertParagraph(content *AtlasContent) (string, error) {
	var para strings.Builder
	for i, child := range content.Content {
		text, err := convertContent(&child)
		if err != nil {
			return "", err
		}
		para.WriteString(text)
		// Add space between inline elements if needed
		if i < len(content.Content)-1 && text != "" && !strings.HasSuffix(text, " ") {
			para.WriteString(" ")
		}
	}
	return para.String() + "\n", nil
}

func convertText(content *AtlasContent) (string, error) {
	text := content.Text

	// Clean up highlight markers
	text = strings.ReplaceAll(text, "@@@hl@@@", "**")
	text = strings.ReplaceAll(text, "@@@endhl@@@", "**")

	// Apply marks in reverse order to handle nested formatting
	if len(content.Marks) > 0 {
		for i := len(content.Marks) - 1; i >= 0; i-- {
			mark := content.Marks[i]
			switch mark.Type {
			case "strong":
				text = "**" + text + "**"
			case "em":
				text = "_" + text + "_"
			case "code":
				text = "`" + text + "`"
			case "link":
				// If the text is already formatted (e.g., with emphasis), keep the formatting
				text = "[" + text + "](" + mark.Attrs.URL + ")"
			case "textColor":
				// Add color indicator based on common colors
				switch mark.Attrs.Color {
				case "#ff0000", "#ff5630", "#de350b": // Red variants
					text = "🔴 " + text
				case "#00ff00", "#36b37e", "#00875a": // Green variants
					text = "🟢 " + text
				case "#ffff00", "#ff991f", "#ff8b00": // Yellow/Orange variants
					text = "⚠️ " + text
				case "#0000ff", "#0052cc", "#0747a6": // Blue variants
					text = "🔵 " + text
				default:
					// For other colors, add a note about the color
					text = text + " _(in " + mark.Attrs.Color + ")_"
				}
			default:
				// Log unknown mark types for future reference
				fmt.Printf("Unknown mark type: %s\n", mark.Type)
			}
		}
	}
	return text, nil
}

func convertTable(content *AtlasContent) (string, error) {
	var table strings.Builder

	// Process rows
	for i, row := range content.Content {
		text, err := convertContent(&row)
		if err != nil {
			return "", err
		}
		table.WriteString(text)

		// Add header separator after first row
		if i == 0 {
			// Count columns by looking at first row
			if len(row.Content) > 0 {
				table.WriteString("|")
				for range row.Content {
					table.WriteString("---|")
				}
				table.WriteString("\n")
			}
		}
	}

	return "\n" + table.String() + "\n", nil
}

func convertTableRow(content *AtlasContent) (string, error) {
	var row strings.Builder
	row.WriteString("|")

	for _, cell := range content.Content {
		text, err := convertContent(&cell)
		if err != nil {
			return "", err
		}
		row.WriteString(text)
		row.WriteString("|")
	}
	row.WriteString("\n")

	return row.String(), nil
}

func convertTableCell(content *AtlasContent) (string, error) {
	var cell strings.Builder
	for _, child := range content.Content {
		text, err := convertContent(&child)
		if err != nil {
			return "", err
		}
		cell.WriteString(strings.TrimSpace(text))
	}
	return " " + strings.ReplaceAll(cell.String(), "\n", " ") + " ", nil
}

func convertBulletList(content *AtlasContent) (string, error) {
	var list strings.Builder
	list.WriteString("\n")

	for _, item := range content.Content {
		text, err := convertContent(&item)
		if err != nil {
			return "", err
		}
		// Ensure proper indentation for nested lists
		lines := strings.Split(strings.TrimSpace(text), "\n")
		for i, line := range lines {
			if i > 0 && !strings.HasPrefix(line, "* ") && !strings.HasPrefix(line, "- ") {
				list.WriteString("  ")
			}
			list.WriteString(line)
			if i < len(lines)-1 {
				list.WriteString("\n")
			}
		}
	}

	list.WriteString("\n")
	return list.String(), nil
}

func convertOrderedList(content *AtlasContent) (string, error) {
	var list strings.Builder
	list.WriteString("\n")

	for i, item := range content.Content {
		text, err := convertContent(&item)
		if err != nil {
			return "", err
		}
		// Remove any bullet points that might be present
		text = strings.TrimPrefix(strings.TrimSpace(text), "* ")
		text = strings.TrimPrefix(text, "- ")

		// Handle multi-line list items
		lines := strings.Split(text, "\n")
		for j, line := range lines {
			if j == 0 {
				list.WriteString(fmt.Sprintf("%d. %s", i+1, line))
			} else {
				list.WriteString(fmt.Sprintf("   %s", line))
			}
			if j < len(lines)-1 {
				list.WriteString("\n")
			}
		}
		list.WriteString("\n")
	}

	return list.String(), nil
}

func convertListItem(content *AtlasContent) (string, error) {
	var item strings.Builder
	item.WriteString("* ")

	for i, child := range content.Content {
		text, err := convertContent(&child)
		if err != nil {
			return "", err
		}
		// Handle multi-line list items
		lines := strings.Split(strings.TrimSpace(text), "\n")
		for j, line := range lines {
			if j > 0 {
				item.WriteString("  ")
			}
			item.WriteString(line)
			if j < len(lines)-1 {
				item.WriteString("\n")
			}
		}
		if i < len(content.Content)-1 {
			item.WriteString(" ")
		}
	}

	return item.String() + "\n", nil
}

func convertInlineCard(content *AtlasContent) (string, error) {
	if content.Attrs.URL != "" {
		return fmt.Sprintf("[%s](%s)", content.Attrs.URL, content.Attrs.URL), nil
	}
	if content.Attrs.ReferencePageID != "" {
		if content.Attrs.ReferenceStatus == "deleted" || content.Attrs.ReferenceStatus == "trashed" {
			return fmt.Sprintf("[%s] _(referenced page no longer exists)_", content.Attrs.ReferencePageTitle), nil
		}
		return fmt.Sprintf("[%s](pages/%s)", content.Attrs.ReferencePageTitle, content.Attrs.ReferencePageID), nil
	}
	return "", nil
}

func convertStatus(content *AtlasContent) (string, error) {
	// Get status text from attributes
	text := content.Attrs.Text
	if text == "" {
		return "", nil
	}

	// Format based on color if available
	switch content.Attrs.Color {
	case "red":
		return fmt.Sprintf("❌ %s", text), nil
	case "green":
		return fmt.Sprintf("✅ %s", text), nil
	case "yellow":
		return fmt.Sprintf("⚠️ %s", text), nil
	case "blue":
		return fmt.Sprintf("ℹ️ %s", text), nil
	case "neutral":
		return fmt.Sprintf("⚪ %s", text), nil
	default:
		return text, nil
	}
}

func convertHeading(content *AtlasContent) (string, error) {
	level := content.Attrs.Level
	if level < 1 || level > 6 {
		level = 1
	}

	var heading strings.Builder
	heading.WriteString("\n")
	for i := 0; i < level; i++ {
		heading.WriteString("#")
	}
	heading.WriteString(" ")

	for _, child := range content.Content {
		text, err := convertContent(&child)
		if err != nil {
			return "", err
		}
		heading.WriteString(text)
	}

	heading.WriteString("\n\n")
	return heading.String(), nil
}

func convertEmoji(content *AtlasContent) (string, error) {
	return content.Attrs.Text, nil
}

func convertPanel(content *AtlasContent) (string, error) {
	var panel strings.Builder

	// // Debug logging
	// fmt.Printf("Panel Type: %s\n", content.Attrs.PanelType)
	// fmt.Printf("Panel Content: %+v\n", content.Content)

	// Add panel type indicator with blockquote
	switch content.Attrs.PanelType {
	case "info":
		panel.WriteString("\n> ℹ️  **Info**\n")
	case "note":
		panel.WriteString("\n> 📝 **Note**\n")
	case "warning":
		panel.WriteString("\n> ⚠️  **Warning**\n")
	case "error":
		panel.WriteString("\n> ❌ **Error**\n")
	case "success":
		panel.WriteString("\n> ✅ **Success**\n")
	default:
		panel.WriteString(fmt.Sprintf("\n> ℹ️  **%s**\n", strings.Title(strings.ToLower(content.Attrs.PanelType))))
	}

	// Convert panel content with blockquote formatting
	for i, child := range content.Content {
		text, err := convertContent(&child)
		if err != nil {
			return "", err
		}
		// // Debug logging
		// fmt.Printf("Panel Content %d: %s\n", i, text)

		// Clean up extra newlines and add blockquote prefix
		lines := strings.Split(strings.TrimSpace(text), "\n")
		for j, line := range lines {
			if j > 0 || i > 0 {
				panel.WriteString("> ")
			}
			panel.WriteString(line)
			if j < len(lines)-1 {
				panel.WriteString("\n")
			}
		}
		// Add newline between paragraphs
		if i < len(content.Content)-1 {
			panel.WriteString("\n>\n")
		}
	}

	panel.WriteString("\n")
	return panel.String(), nil
}

func convertTaskList(content *AtlasContent) (string, error) {
	var list strings.Builder
	list.WriteString("\n")

	for _, item := range content.Content {
		text, err := convertContent(&item)
		if err != nil {
			return "", err
		}
		list.WriteString(text)
	}

	return list.String() + "\n", nil
}

func convertTaskItem(content *AtlasContent) (string, error) {
	var item strings.Builder

	// Add checkbox
	if content.Attrs.State == "DONE" {
		item.WriteString("- [x] ")
	} else {
		item.WriteString("- [ ] ")
	}

	for _, child := range content.Content {
		text, err := convertContent(&child)
		if err != nil {
			return "", err
		}
		item.WriteString(strings.TrimSpace(text))
	}

	return item.String() + "\n", nil
}

func convertCodeBlock(content *AtlasContent) (string, error) {
	var code strings.Builder
	code.WriteString("\n```")

	// Add language if specified
	if content.Attrs.Language != "" {
		code.WriteString(content.Attrs.Language)
	}
	code.WriteString("\n")

	// Process content and handle newlines
	for i, child := range content.Content {
		text, err := convertContent(&child)
		if err != nil {
			return "", err
		}
		// Trim trailing whitespace but preserve indentation
		text = strings.TrimRightFunc(text, unicode.IsSpace)
		code.WriteString(text)

		// Add newline between content items except for the last one
		if i < len(content.Content)-1 {
			code.WriteString("\n")
		}
	}

	// Ensure single newline before closing fence
	if !strings.HasSuffix(code.String(), "\n") {
		code.WriteString("\n")
	}
	code.WriteString("```\n\n")
	return code.String(), nil
}

func convertBodiedExtension(content *AtlasContent) (string, error) {
	// Extract extension parameters
	extensionKey := content.Attrs.ExtensionKey

	// Handle different extension types
	switch extensionKey {
	case "details":
		// Convert the content inside the extension
		var result strings.Builder
		for _, child := range content.Content {
			text, err := convertContent(&child)
			if err != nil {
				return "", err
			}
			result.WriteString(text)
		}
		return result.String(), nil
	default:
		// For unknown extension types, just convert the content
		var result strings.Builder
		for _, child := range content.Content {
			text, err := convertContent(&child)
			if err != nil {
				return "", err
			}
			result.WriteString(text)
		}
		return result.String(), nil
	}
}

func convertExtension(content *AtlasContent) (string, error) {
	// Extract extension key
	extensionKey := content.Attrs.ExtensionKey

	// Handle different extension types
	switch extensionKey {
	case "children":
		return "\n_This page displays a list of child pages in Confluence._\n\n", nil
	default:
		// For unknown extension types, return empty string
		return "", nil
	}
}

func convertDate(content *AtlasContent) (string, error) {
	// Convert timestamp to readable date
	timestamp := content.Attrs.Timestamp
	if timestamp == "" {
		return "", nil
	}

	// Parse timestamp (milliseconds since epoch)
	ms, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return "", err
	}

	// Convert to time
	t := time.Unix(ms/1000, 0)
	return t.Format("2006-01-02"), nil
}

func convertMention(content *AtlasContent) (string, error) {
	// Return the display text of the mention
	return fmt.Sprintf("@%s", content.Attrs.Text), nil
}

func convertPlaceholder(content *AtlasContent) (string, error) {
	// Return the placeholder text in italics
	return fmt.Sprintf("_%s_", content.Attrs.Text), nil
}

func convertLayoutSection(content *AtlasContent) (string, error) {
	var section strings.Builder

	// Process each column in the section
	for _, column := range content.Content {
		text, err := convertContent(&column)
		if err != nil {
			return "", err
		}
		section.WriteString(text)
	}

	return section.String(), nil
}

func convertLayoutColumn(content *AtlasContent) (string, error) {
	var column strings.Builder

	// Process the content within the column
	for _, child := range content.Content {
		text, err := convertContent(&child)
		if err != nil {
			return "", err
		}
		column.WriteString(text)
	}

	return column.String(), nil
}

func convertMediaSingle(content *AtlasContent) (string, error) {
	var result strings.Builder

	// Process each media item in the mediaSingle container
	for _, media := range content.Content {
		text, err := convertContent(&media)
		if err != nil {
			return "", err
		}
		result.WriteString(text)
	}

	return result.String(), nil
}

func convertMedia(content *AtlasContent) (string, error) {
	// For now, just add a placeholder note about the image
	// In the future, we could potentially download and embed the image
	return fmt.Sprintf("\n_[Image: %dx%d]_\n\n", content.Attrs.MediaWidth, content.Attrs.MediaHeight), nil
}

func convertNestedExpand(content *AtlasContent) (string, error) {
	// Nested expands are handled the same way as regular expands
	return convertExpand(content)
}

func convertExpand(content *AtlasContent) (string, error) {
	var expand strings.Builder

	// Add title if available
	if content.Attrs.Title != "" {
		expand.WriteString(fmt.Sprintf("\n<details>\n<summary>%s</summary>\n\n", content.Attrs.Title))
	} else {
		expand.WriteString("\n<details>\n<summary>Details</summary>\n\n")
	}

	// Convert nested content
	for _, child := range content.Content {
		text, err := convertContent(&child)
		if err != nil {
			return "", err
		}
		expand.WriteString(text)
	}

	expand.WriteString("</details>\n\n")
	return expand.String(), nil
}

func convertInlineExtension(content *AtlasContent) (string, error) {
	// Extract extension parameters
	extensionKey := content.Attrs.ExtensionKey
	extensionType := content.Attrs.ExtensionType

	// Handle different extension types
	switch extensionKey {
	case "pagetree":
		// For pagetree macro, add a note about the tree structure
		if content.Attrs.Parameters != nil {
			var params struct {
				MacroParams struct {
					Root struct {
						Value string `json:"value"`
					} `json:"root"`
				} `json:"macroParams"`
			}
			if err := json.Unmarshal(content.Attrs.Parameters, &params); err == nil {
				if params.MacroParams.Root.Value != "" {
					return fmt.Sprintf("\n_[Page tree structure - showing child pages under %s]_\n\n",
						params.MacroParams.Root.Value), nil
				}
			}
		}
		return "\n_[Page tree structure - showing child pages]_\n\n", nil
	case "jira":
		// For JIRA macros
		return "\n_[JIRA issue/filter]_\n\n", nil
	case "confluence-content":
		// For content includes
		return "\n_[Included Confluence content]_\n\n", nil
	case "drawio":
		// For draw.io diagrams
		return "\n_[draw.io diagram]_\n\n", nil
	case "plantuml":
		// For PlantUML diagrams
		return "\n_[PlantUML diagram]_\n\n", nil
	default:
		// For unknown extension types, add a note about the macro and its type
		if extensionType != "" {
			return fmt.Sprintf("\n_[Confluence macro: %s (%s)]_\n\n", extensionKey, extensionType), nil
		}
		return fmt.Sprintf("\n_[Confluence macro: %s]_\n\n", extensionKey), nil
	}
}

// JiraAtlasDocument represents the root structure of a Jira atlas_doc_format document
type JiraAtlasDocument struct {
	Type    string             `json:"type"`
	Content []JiraAtlasContent `json:"content"`
	Version int                `json:"version"`
}

// JiraAtlasContent represents a content node in the Jira atlas document
type JiraAtlasContent struct {
	Type    string              `json:"type"`
	Content []JiraAtlasContent  `json:"content,omitempty"`
	Text    string              `json:"text,omitempty"`
	Marks   []JiraAtlasMark     `json:"marks,omitempty"`
	Attrs   JiraAtlasAttributes `json:"attrs,omitempty"`
}

// JiraAtlasMark represents text formatting in the Jira atlas document
type JiraAtlasMark struct {
	Type  string `json:"type"`
	Attrs struct {
		URL   string `json:"url,omitempty"`
		Color string `json:"color,omitempty"`
	} `json:"attrs,omitempty"`
}

// JiraAtlasAttributes represents the attributes of a Jira Atlas content node
type JiraAtlasAttributes struct {
	Level              int                    `json:"level,omitempty"`
	Text               string                 `json:"text,omitempty"`
	Title              string                 `json:"title,omitempty"`
	ExtensionType      string                 `json:"extensionType,omitempty"`
	ExtensionKey       string                 `json:"extensionKey,omitempty"`
	Parameters         json.RawMessage        `json:"parameters,omitempty"`
	URL                string                 `json:"url,omitempty"`
	ReferencePageID    string                 `json:"referencePageId,omitempty"`
	ReferenceStatus    string                 `json:"referenceStatus,omitempty"`
	ReferencePageTitle string                 `json:"referencePageTitle,omitempty"`
	Color              string                 `json:"color,omitempty"`
	PanelType          string                 `json:"panelType,omitempty"`
	State              string                 `json:"state,omitempty"`
	Language           string                 `json:"language,omitempty"`
	Timestamp          string                 `json:"timestamp,omitempty"`
	MediaWidth         int                    `json:"mediaWidth,omitempty"`
	MediaHeight        int                    `json:"mediaHeight,omitempty"`
	Other              map[string]interface{} `json:"-"`
}

// ConvertJiraAtlasToMarkdown converts Jira atlas_doc_format JSON string to Markdown
func ConvertJiraAtlasToMarkdown(jsonStr string) (string, error) {
	var doc JiraAtlasDocument
	if err := json.Unmarshal([]byte(jsonStr), &doc); err != nil {
		fmt.Printf("Error unmarshaling Jira atlas document: %v\n", err)
		return "", fmt.Errorf("failed to parse Jira atlas document: %v", err)
	}

	var result strings.Builder
	for _, content := range doc.Content {
		text, err := convertJiraContent(&content)
		if err != nil {
			fmt.Printf("Error converting Jira content: %v\n", err)
			return "", err
		}
		result.WriteString(text)
	}

	return result.String(), nil
}

// convertJiraContent converts a single Jira content node to Markdown
func convertJiraContent(content *JiraAtlasContent) (string, error) {
	switch content.Type {
	case "paragraph":
		return convertJiraParagraph(content)
	case "text":
		return convertJiraText(content)
	case "table":
		return convertJiraTable(content)
	case "tableRow":
		return convertJiraTableRow(content)
	case "tableHeader", "tableCell":
		return convertJiraTableCell(content)
	case "bulletList":
		return convertJiraBulletList(content)
	case "orderedList":
		return convertJiraOrderedList(content)
	case "listItem":
		return convertJiraListItem(content)
	case "inlineCard":
		return convertJiraInlineCard(content)
	case "status":
		return convertJiraStatus(content)
	case "heading":
		return convertJiraHeading(content)
	case "emoji":
		return convertJiraEmoji(content)
	case "panel":
		return convertJiraPanel(content)
	case "taskList":
		return convertJiraTaskList(content)
	case "taskItem":
		return convertJiraTaskItem(content)
	case "rule":
		return "---\n\n", nil
	case "codeBlock":
		return convertJiraCodeBlock(content)
	case "hardBreak":
		return "\n", nil
	case "bodiedExtension":
		return convertJiraBodiedExtension(content)
	case "extension":
		return convertJiraExtension(content)
	case "date":
		return convertJiraDate(content)
	case "mention":
		return convertJiraMention(content)
	case "placeholder":
		return convertJiraPlaceholder(content)
	case "layoutSection":
		return convertJiraLayoutSection(content)
	case "layoutColumn":
		return convertJiraLayoutColumn(content)
	case "mediaSingle":
		return convertJiraMediaSingle(content)
	case "media":
		return convertJiraMedia(content)
	case "nestedExpand":
		return convertJiraNestedExpand(content)
	case "expand":
		return convertJiraExpand(content)
	case "inlineExtension":
		return convertJiraInlineExtension(content)
	default:
		return "", fmt.Errorf("unsupported content type: %s", content.Type)
	}
}

func convertJiraParagraph(content *JiraAtlasContent) (string, error) {
	var para strings.Builder
	for _, child := range content.Content {
		text, err := convertJiraContent(&child)
		if err != nil {
			return "", err
		}
		para.WriteString(text)
		if len(content.Content) > 1 && !strings.HasSuffix(text, " ") {
			para.WriteString(" ")
		}
	}
	return para.String() + "\n", nil
}

func convertJiraText(content *JiraAtlasContent) (string, error) {
	text := content.Text

	// Clean up highlight markers
	text = strings.ReplaceAll(text, "@@@hl@@@", "**")
	text = strings.ReplaceAll(text, "@@@endhl@@@", "**")

	if len(content.Marks) > 0 {
		for i := len(content.Marks) - 1; i >= 0; i-- {
			mark := content.Marks[i]
			switch mark.Type {
			case "strong":
				text = "**" + text + "**"
			case "em":
				text = "_" + text + "_"
			case "code":
				text = "`" + text + "`"
			case "link":
				text = "[" + text + "](" + mark.Attrs.URL + ")"
			case "textColor":
				// Add color indicator based on common colors
				switch mark.Attrs.Color {
				case "#ff0000", "#ff5630", "#de350b": // Red variants
					text = "🔴 " + text
				case "#00ff00", "#36b37e", "#00875a": // Green variants
					text = "🟢 " + text
				case "#ffff00", "#ff991f", "#ff8b00": // Yellow/Orange variants
					text = "⚠️ " + text
				case "#0000ff", "#0052cc", "#0747a6": // Blue variants
					text = "🔵 " + text
				default:
					// For other colors, add a note about the color
					text = text + " _(in " + mark.Attrs.Color + ")_"
				}
			default:
				fmt.Printf("Unknown mark type: %s\n", mark.Type)
			}
		}
	}
	return text, nil
}

func convertJiraTable(content *JiraAtlasContent) (string, error) {
	var table strings.Builder

	// Process rows
	for i, row := range content.Content {
		text, err := convertJiraContent(&row)
		if err != nil {
			return "", err
		}
		table.WriteString(text)

		// Add header separator after first row
		if i == 0 {
			if len(row.Content) > 0 {
				table.WriteString("|")
				for range row.Content {
					table.WriteString("---|")
				}
				table.WriteString("\n")
			}
		}
	}

	return "\n" + table.String() + "\n", nil
}

func convertJiraTableRow(content *JiraAtlasContent) (string, error) {
	var row strings.Builder
	row.WriteString("|")

	for _, cell := range content.Content {
		text, err := convertJiraContent(&cell)
		if err != nil {
			return "", err
		}
		row.WriteString(text)
		row.WriteString("|")
	}
	row.WriteString("\n")

	return row.String(), nil
}

func convertJiraTableCell(content *JiraAtlasContent) (string, error) {
	var cell strings.Builder
	for _, child := range content.Content {
		text, err := convertJiraContent(&child)
		if err != nil {
			return "", err
		}
		cell.WriteString(strings.TrimSpace(text))
	}
	return " " + strings.ReplaceAll(cell.String(), "\n", " ") + " ", nil
}

func convertJiraBulletList(content *JiraAtlasContent) (string, error) {
	var list strings.Builder
	list.WriteString("\n")

	for _, item := range content.Content {
		text, err := convertJiraContent(&item)
		if err != nil {
			return "", err
		}
		lines := strings.Split(strings.TrimSpace(text), "\n")
		for i, line := range lines {
			if i > 0 && !strings.HasPrefix(line, "* ") && !strings.HasPrefix(line, "- ") {
				list.WriteString("  ")
			}
			list.WriteString(line)
			if i < len(lines)-1 {
				list.WriteString("\n")
			}
		}
	}

	list.WriteString("\n")
	return list.String(), nil
}

func convertJiraOrderedList(content *JiraAtlasContent) (string, error) {
	var list strings.Builder
	list.WriteString("\n")

	for i, item := range content.Content {
		text, err := convertJiraContent(&item)
		if err != nil {
			return "", err
		}
		text = strings.TrimPrefix(strings.TrimSpace(text), "* ")
		text = strings.TrimPrefix(text, "- ")

		lines := strings.Split(text, "\n")
		for j, line := range lines {
			if j == 0 {
				list.WriteString(fmt.Sprintf("%d. %s", i+1, line))
			} else {
				list.WriteString(fmt.Sprintf("   %s", line))
			}
			if j < len(lines)-1 {
				list.WriteString("\n")
			}
		}
		list.WriteString("\n")
	}

	return list.String() + "\n", nil
}

func convertJiraListItem(content *JiraAtlasContent) (string, error) {
	var item strings.Builder
	item.WriteString("* ")

	for i, child := range content.Content {
		text, err := convertJiraContent(&child)
		if err != nil {
			return "", err
		}
		lines := strings.Split(strings.TrimSpace(text), "\n")
		for j, line := range lines {
			if j > 0 {
				item.WriteString("  ")
			}
			item.WriteString(line)
			if j < len(lines)-1 {
				item.WriteString("\n")
			}
		}
		if i < len(content.Content)-1 {
			item.WriteString(" ")
		}
	}

	return item.String() + "\n", nil
}

func convertJiraInlineCard(content *JiraAtlasContent) (string, error) {
	if content.Attrs.URL != "" {
		return fmt.Sprintf("[%s](%s)", content.Attrs.URL, content.Attrs.URL), nil
	}
	if content.Attrs.ReferencePageID != "" {
		if content.Attrs.ReferenceStatus == "deleted" || content.Attrs.ReferenceStatus == "trashed" {
			return fmt.Sprintf("[%s] _(referenced page no longer exists)_", content.Attrs.ReferencePageTitle), nil
		}
		return fmt.Sprintf("[%s](pages/%s)", content.Attrs.ReferencePageTitle, content.Attrs.ReferencePageID), nil
	}
	return "", nil
}

func convertJiraStatus(content *JiraAtlasContent) (string, error) {
	text := content.Attrs.Text
	if text == "" {
		return "", nil
	}

	switch content.Attrs.Color {
	case "red":
		return fmt.Sprintf("❌ %s", text), nil
	case "green":
		return fmt.Sprintf("✅ %s", text), nil
	case "yellow":
		return fmt.Sprintf("⚠️ %s", text), nil
	case "blue":
		return fmt.Sprintf("ℹ️ %s", text), nil
	case "neutral":
		return fmt.Sprintf("⚪ %s", text), nil
	default:
		return text, nil
	}
}

func convertJiraHeading(content *JiraAtlasContent) (string, error) {
	level := content.Attrs.Level
	if level < 1 || level > 6 {
		level = 1
	}

	var heading strings.Builder
	heading.WriteString("\n")
	for i := 0; i < level; i++ {
		heading.WriteString("#")
	}
	heading.WriteString(" ")

	for _, child := range content.Content {
		text, err := convertJiraContent(&child)
		if err != nil {
			return "", err
		}
		heading.WriteString(text)
	}

	heading.WriteString("\n\n")
	return heading.String(), nil
}

func convertJiraEmoji(content *JiraAtlasContent) (string, error) {
	return content.Attrs.Text, nil
}

func convertJiraCodeBlock(content *JiraAtlasContent) (string, error) {
	var code strings.Builder
	code.WriteString("\n```")

	if content.Attrs.Language != "" {
		code.WriteString(content.Attrs.Language)
	}
	code.WriteString("\n")

	for i, child := range content.Content {
		text, err := convertJiraContent(&child)
		if err != nil {
			return "", err
		}
		text = strings.TrimRightFunc(text, unicode.IsSpace)
		code.WriteString(text)

		if i < len(content.Content)-1 {
			code.WriteString("\n")
		}
	}

	if !strings.HasSuffix(code.String(), "\n") {
		code.WriteString("\n")
	}
	code.WriteString("```\n\n")
	return code.String(), nil
}

func convertJiraTaskList(content *JiraAtlasContent) (string, error) {
	var list strings.Builder
	list.WriteString("\n")

	for _, item := range content.Content {
		text, err := convertJiraContent(&item)
		if err != nil {
			return "", err
		}
		list.WriteString(text)
	}

	return list.String() + "\n", nil
}

func convertJiraTaskItem(content *JiraAtlasContent) (string, error) {
	var item strings.Builder

	if content.Attrs.State == "DONE" {
		item.WriteString("- [x] ")
	} else {
		item.WriteString("- [ ] ")
	}

	for _, child := range content.Content {
		text, err := convertJiraContent(&child)
		if err != nil {
			return "", err
		}
		item.WriteString(strings.TrimSpace(text))
	}

	return item.String() + "\n", nil
}

func convertJiraPanel(content *JiraAtlasContent) (string, error) {
	var panel strings.Builder

	switch strings.ToLower(content.Attrs.PanelType) {
	case "info":
		panel.WriteString("\n> ℹ️  **Info**\n")
	case "note":
		panel.WriteString("\n> 📝 **Note**\n")
	case "warning":
		panel.WriteString("\n> ⚠️  **Warning**\n")
	case "error":
		panel.WriteString("\n> ❌ **Error**\n")
	case "success":
		panel.WriteString("\n> ✅ **Success**\n")
	default:
		if content.Attrs.PanelType != "" {
			panel.WriteString(fmt.Sprintf("\n> ℹ️  **%s**\n", strings.Title(strings.ToLower(content.Attrs.PanelType))))
		} else {
			panel.WriteString("\n> ℹ️  **Note**\n")
		}
	}

	for _, child := range content.Content {
		text, err := convertJiraContent(&child)
		if err != nil {
			return "", err
		}
		text = strings.TrimSpace(text)
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			if line != "" {
				panel.WriteString("> " + line + "\n")
			}
		}
	}

	panel.WriteString("\n")
	return panel.String(), nil
}

func convertJiraBodiedExtension(content *JiraAtlasContent) (string, error) {
	return fmt.Sprintf("[%s extension]", content.Attrs.ExtensionKey), nil
}

func convertJiraExtension(content *JiraAtlasContent) (string, error) {
	return fmt.Sprintf("[%s extension]", content.Attrs.ExtensionKey), nil
}

func convertJiraDate(content *JiraAtlasContent) (string, error) {
	if content.Attrs.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, content.Attrs.Timestamp); err == nil {
			return t.Format("Jan 02, 2006"), nil
		}
	}
	return content.Text, nil
}

func convertJiraMention(content *JiraAtlasContent) (string, error) {
	return content.Text, nil
}

func convertJiraPlaceholder(content *JiraAtlasContent) (string, error) {
	return content.Text, nil
}

func convertJiraLayoutSection(content *JiraAtlasContent) (string, error) {
	var section strings.Builder
	for _, child := range content.Content {
		text, err := convertJiraContent(&child)
		if err != nil {
			return "", err
		}
		section.WriteString(text)
	}
	return section.String(), nil
}

func convertJiraLayoutColumn(content *JiraAtlasContent) (string, error) {
	var column strings.Builder
	for _, child := range content.Content {
		text, err := convertJiraContent(&child)
		if err != nil {
			return "", err
		}
		column.WriteString(text)
	}
	return column.String(), nil
}

func convertJiraMediaSingle(content *JiraAtlasContent) (string, error) {
	return fmt.Sprintf("\n_[Image: %dx%d]_\n\n", content.Attrs.MediaWidth, content.Attrs.MediaHeight), nil
}

func convertJiraMedia(content *JiraAtlasContent) (string, error) {
	return fmt.Sprintf("\n_[Image: %dx%d]_\n\n", content.Attrs.MediaWidth, content.Attrs.MediaHeight), nil
}

func convertJiraNestedExpand(content *JiraAtlasContent) (string, error) {
	return convertJiraExpand(content)
}

func convertJiraExpand(content *JiraAtlasContent) (string, error) {
	var expand strings.Builder

	if content.Attrs.Title != "" {
		expand.WriteString(fmt.Sprintf("\n<details>\n<summary>%s</summary>\n\n", content.Attrs.Title))
	} else {
		expand.WriteString("\n<details>\n<summary>Details</summary>\n\n")
	}

	for _, child := range content.Content {
		text, err := convertJiraContent(&child)
		if err != nil {
			return "", err
		}
		expand.WriteString(text)
	}

	expand.WriteString("</details>\n\n")
	return expand.String(), nil
}

func convertJiraInlineExtension(content *JiraAtlasContent) (string, error) {
	return fmt.Sprintf("[%s extension]", content.Attrs.ExtensionKey), nil
}
