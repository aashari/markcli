package atlassian

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"
)

// AtlassianDocument represents the root structure of an Atlassian Document Format (ADF) document
type AtlassianDocument struct {
	Type    string             `json:"type"`
	Content []AtlassianContent `json:"content"`
	Version int                `json:"version"`
}

// AtlassianContent represents a content node in the Atlassian document.
// Supported content types include:
// - paragraph: A block of text
// - text: Plain text content
// - table, tableRow, tableHeader, tableCell: Table structures
// - bulletList, orderedList, listItem: List structures
// - inlineCard: Links to other content
// - status: Status indicators
// - heading: Section headers
// - emoji: Emoji characters
// - panel: Info, note, warning panels
// - taskList, taskItem: Checkable task items
// - rule: Horizontal rules
// - codeBlock: Code snippets
// - hardBreak: Line breaks
// - bodiedExtension, extension: Confluence macros
// - date: Date values
// - mention: User mentions
// - placeholder: Placeholder text
// - layoutSection, layoutColumn: Layout containers
// - mediaSingle, media: Images and media
// - nestedExpand, expand: Collapsible sections
// - inlineExtension: Inline macros
type AtlassianContent struct {
	Type    string              `json:"type"`
	Content []AtlassianContent  `json:"content,omitempty"`
	Text    string              `json:"text,omitempty"`
	Marks   []AtlassianMark     `json:"marks,omitempty"`
	Attrs   AtlassianAttributes `json:"attrs,omitempty"`
}

// AtlassianMark represents text formatting in the Atlassian document.
// Supported mark types include:
// - strong: Bold text
// - em: Italic text
// - code: Inline code
// - link: Hyperlinks
// - textColor: Colored text
type AtlassianMark struct {
	Type  string `json:"type"`
	Attrs struct {
		URL   string `json:"url,omitempty"`
		Color string `json:"color,omitempty"`
	} `json:"attrs,omitempty"`
}

// AtlassianAttributes represents the attributes of an Atlassian content node.
// Common attributes include:
// - level: Used for headings (1-6)
// - text: Plain text content
// - title: Title text for expandable sections
// - extensionType, extensionKey: Used for macros and extensions
// - parameters: Raw JSON parameters for macros
// - url: Used for links and cards
// - referencePageId, referenceStatus, referencePageTitle: Used for page references
// - color: Used for status and text colors
// - panelType: Type of panel (info, note, warning, error, success)
// - state: Used for task items (DONE or not)
// - language: Used for code blocks
// - timestamp: Used for date values
// - mediaWidth, mediaHeight: Used for images and media
// Additional attributes are stored in the Other field
type AtlassianAttributes struct {
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
func (a *AtlassianAttributes) UnmarshalJSON(data []byte) error {
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
	type Alias AtlassianAttributes
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

// AtlassianDocumentConvertToMarkdown converts an AtlassianDocument to Markdown format
func (doc *AtlassianDocument) AtlassianDocumentConvertToMarkdown() (string, error) {
	if doc.Content == nil {
		return "", nil
	}

	var result strings.Builder
	for _, content := range doc.Content {
		text, err := content.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", fmt.Errorf("failed to convert content: %w", err)
		}
		result.WriteString(text)
	}

	return result.String(), nil
}

// AtlassianDocumentConvertToMarkdown converts an AtlassianContent node to Markdown format
func (content *AtlassianContent) AtlassianDocumentConvertToMarkdown() (string, error) {
	switch content.Type {
	case "paragraph":
		return content.convertParagraph()
	case "text":
		return content.convertText()
	case "table":
		return content.convertTable()
	case "tableRow":
		return content.convertTableRow()
	case "tableHeader", "tableCell":
		return content.convertTableCell()
	case "bulletList":
		return content.convertBulletList()
	case "orderedList":
		return content.convertOrderedList()
	case "listItem":
		return content.convertListItem()
	case "inlineCard":
		return content.convertInlineCard()
	case "status":
		return content.convertStatus()
	case "heading":
		return content.convertHeading()
	case "emoji":
		return content.convertEmoji()
	case "panel":
		return content.convertPanel()
	case "taskList":
		return content.convertTaskList()
	case "taskItem":
		return content.convertTaskItem()
	case "rule":
		return "---\n\n", nil
	case "codeBlock":
		return content.convertCodeBlock()
	case "hardBreak":
		return "\n", nil
	case "bodiedExtension":
		return content.convertBodiedExtension()
	case "extension":
		return content.convertExtension()
	case "date":
		return content.convertDate()
	case "mention":
		return content.convertMention()
	case "placeholder":
		return content.convertPlaceholder()
	case "layoutSection":
		return content.convertLayoutSection()
	case "layoutColumn":
		return content.convertLayoutColumn()
	case "mediaSingle":
		return content.convertMediaSingle()
	case "media":
		return content.convertMedia()
	case "nestedExpand":
		return content.convertNestedExpand()
	case "expand":
		return content.convertExpand()
	case "inlineExtension":
		return content.convertInlineExtension()
	default:
		return "", fmt.Errorf("unsupported content type: %s", content.Type)
	}
}

// Helper methods for converting different content types to Markdown
func (content *AtlassianContent) convertParagraph() (string, error) {
	var para strings.Builder
	for i, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		para.WriteString(text)
		if i < len(content.Content)-1 && !strings.HasSuffix(text, " ") {
			para.WriteString(" ")
		}
	}
	return para.String() + "\n\n", nil
}

func (content *AtlassianContent) convertText() (string, error) {
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
					text = text + " _(in " + mark.Attrs.Color + ")_"
				}
			}
		}
	}
	return text, nil
}

func (content *AtlassianContent) convertTable() (string, error) {
	var table strings.Builder
	table.WriteString("\n")

	for i, row := range content.Content {
		text, err := row.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		table.WriteString(text)

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

	table.WriteString("\n")
	return table.String(), nil
}

func (content *AtlassianContent) convertTableRow() (string, error) {
	var row strings.Builder
	row.WriteString("|")

	for _, cell := range content.Content {
		text, err := cell.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		row.WriteString(text)
		row.WriteString("|")
	}
	row.WriteString("\n")

	return row.String(), nil
}

func (content *AtlassianContent) convertTableCell() (string, error) {
	var cell strings.Builder
	for _, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		cell.WriteString(strings.TrimSpace(text))
	}
	return " " + strings.ReplaceAll(cell.String(), "\n", " ") + " ", nil
}

func (content *AtlassianContent) convertBulletList() (string, error) {
	var list strings.Builder
	list.WriteString("\n")

	for _, item := range content.Content {
		text, err := item.AtlassianDocumentConvertToMarkdown()
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

func (content *AtlassianContent) convertOrderedList() (string, error) {
	var list strings.Builder
	list.WriteString("\n")

	for i, item := range content.Content {
		text, err := item.AtlassianDocumentConvertToMarkdown()
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

	return list.String(), nil
}

func (content *AtlassianContent) convertListItem() (string, error) {
	var item strings.Builder
	item.WriteString("* ")

	for i, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
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

func (content *AtlassianContent) convertInlineCard() (string, error) {
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

func (content *AtlassianContent) convertStatus() (string, error) {
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

func (content *AtlassianContent) convertHeading() (string, error) {
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
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		heading.WriteString(text)
	}

	heading.WriteString("\n\n")
	return heading.String(), nil
}

func (content *AtlassianContent) convertEmoji() (string, error) {
	return content.Attrs.Text, nil
}

func (content *AtlassianContent) convertCodeBlock() (string, error) {
	var code strings.Builder
	code.WriteString("\n```")

	if content.Attrs.Language != "" {
		code.WriteString(content.Attrs.Language)
	}
	code.WriteString("\n")

	for i, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
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

func (content *AtlassianContent) convertTaskList() (string, error) {
	var list strings.Builder
	list.WriteString("\n")

	for _, item := range content.Content {
		text, err := item.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		list.WriteString(text)
	}

	return list.String() + "\n", nil
}

func (content *AtlassianContent) convertTaskItem() (string, error) {
	var item strings.Builder

	if content.Attrs.State == "DONE" {
		item.WriteString("- [x] ")
	} else {
		item.WriteString("- [ ] ")
	}

	for _, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		item.WriteString(strings.TrimSpace(text))
	}

	return item.String() + "\n", nil
}

func (content *AtlassianContent) convertPanel() (string, error) {
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
		text, err := child.AtlassianDocumentConvertToMarkdown()
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

func (content *AtlassianContent) convertBodiedExtension() (string, error) {
	return fmt.Sprintf("[%s extension]", content.Attrs.ExtensionKey), nil
}

func (content *AtlassianContent) convertExtension() (string, error) {
	return fmt.Sprintf("[%s extension]", content.Attrs.ExtensionKey), nil
}

func (content *AtlassianContent) convertDate() (string, error) {
	if content.Attrs.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, content.Attrs.Timestamp); err == nil {
			return t.Format("Jan 02, 2006"), nil
		}
	}
	return content.Text, nil
}

func (content *AtlassianContent) convertMention() (string, error) {
	return content.Text, nil
}

func (content *AtlassianContent) convertPlaceholder() (string, error) {
	return content.Text, nil
}

func (content *AtlassianContent) convertLayoutSection() (string, error) {
	var section strings.Builder
	for _, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		section.WriteString(text)
	}
	return section.String(), nil
}

func (content *AtlassianContent) convertLayoutColumn() (string, error) {
	var column strings.Builder
	for _, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		column.WriteString(text)
	}
	return column.String(), nil
}

func (content *AtlassianContent) convertMediaSingle() (string, error) {
	return fmt.Sprintf("\n_[Image: %dx%d]_\n\n", content.Attrs.MediaWidth, content.Attrs.MediaHeight), nil
}

func (content *AtlassianContent) convertMedia() (string, error) {
	return fmt.Sprintf("\n_[Image: %dx%d]_\n\n", content.Attrs.MediaWidth, content.Attrs.MediaHeight), nil
}

func (content *AtlassianContent) convertNestedExpand() (string, error) {
	return content.convertExpand()
}

func (content *AtlassianContent) convertExpand() (string, error) {
	var expand strings.Builder

	if content.Attrs.Title != "" {
		expand.WriteString(fmt.Sprintf("\n<details>\n<summary>%s</summary>\n\n", content.Attrs.Title))
	} else {
		expand.WriteString("\n<details>\n<summary>Details</summary>\n\n")
	}

	for _, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		expand.WriteString(text)
	}

	expand.WriteString("</details>\n\n")
	return expand.String(), nil
}

func (content *AtlassianContent) convertInlineExtension() (string, error) {
	return fmt.Sprintf("[%s extension]", content.Attrs.ExtensionKey), nil
}

// ParseDocument parses a JSON string into an AtlassianDocument.
// This function is used to convert raw JSON content from Atlassian APIs into a structured document.
// Returns an error if the JSON is invalid or cannot be parsed into an AtlassianDocument.
func ParseDocument(jsonStr string) (*AtlassianDocument, error) {
	var doc AtlassianDocument
	if err := json.Unmarshal([]byte(jsonStr), &doc); err != nil {
		return nil, fmt.Errorf("failed to parse document: %v", err)
	}
	return &doc, nil
}

// AtlassianDocumentConvertJSONToMarkdown converts a JSON string containing an Atlassian document to Markdown.
// This is a convenience function that parses the JSON and calls AtlassianDocumentConvertToMarkdown.
func AtlassianDocumentConvertJSONToMarkdown(jsonStr string) (string, error) {
	var doc AtlassianDocument
	if err := json.Unmarshal([]byte(jsonStr), &doc); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}
	return doc.AtlassianDocumentConvertToMarkdown()
}
