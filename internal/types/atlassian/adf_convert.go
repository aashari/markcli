package atlassian

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

	// Apply marks
	if content.Marks != nil {
		for _, mark := range content.Marks {
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
				// Skip color formatting in markdown
			}
		}
	}

	return text, nil
}

func (content *AtlassianContent) convertTable() (string, error) {
	if len(content.Content) == 0 {
		return "", nil
	}

	var table strings.Builder
	var headerRow bool

	// Process each row
	for i, row := range content.Content {
		if row.Type != "tableRow" {
			continue
		}

		text, err := row.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		table.WriteString(text)

		// Add header separator after first row
		if i == 0 && !headerRow {
			headerRow = true
			if len(row.Content) > 0 {
				table.WriteString("|")
				for range row.Content {
					table.WriteString(" --- |")
				}
				table.WriteString("\n")
			}
		}
	}

	return table.String() + "\n", nil
}

func (content *AtlassianContent) convertTableRow() (string, error) {
	var row strings.Builder
	row.WriteString("|")

	for _, cell := range content.Content {
		text, err := cell.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		// Clean up newlines in cells
		text = strings.ReplaceAll(text, "\n", " ")
		text = strings.TrimSpace(text)
		row.WriteString(" " + text + " |")
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
		cell.WriteString(text)
	}
	return cell.String(), nil
}

func (content *AtlassianContent) convertBulletList() (string, error) {
	var list strings.Builder
	for _, item := range content.Content {
		if item.Type != "listItem" {
			continue
		}

		text, err := item.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}

		// Add bullet point to each line
		lines := strings.Split(strings.TrimSpace(text), "\n")
		for i, line := range lines {
			if i == 0 {
				list.WriteString("* " + line + "\n")
			} else {
				list.WriteString("  " + line + "\n")
			}
		}
	}

	return list.String() + "\n", nil
}

func (content *AtlassianContent) convertOrderedList() (string, error) {
	var list strings.Builder
	for i, item := range content.Content {
		if item.Type != "listItem" {
			continue
		}

		text, err := item.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}

		// Add number to each line
		lines := strings.Split(strings.TrimSpace(text), "\n")
		for j, line := range lines {
			if j == 0 {
				list.WriteString(fmt.Sprintf("%d. %s\n", i+1, line))
			} else {
				list.WriteString("   " + line + "\n")
			}
		}
	}

	return list.String() + "\n", nil
}

func (content *AtlassianContent) convertListItem() (string, error) {
	var item strings.Builder
	for _, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		item.WriteString(text)
	}
	return item.String(), nil
}

func (content *AtlassianContent) convertInlineCard() (string, error) {
	if content.Attrs.URL != "" {
		return fmt.Sprintf("[%s](%s)", content.Attrs.URL, content.Attrs.URL), nil
	}
	return "", nil
}

func (content *AtlassianContent) convertStatus() (string, error) {
	text := content.Attrs.Text
	if text == "" {
		return "", nil
	}

	color := content.Attrs.Color
	if color != "" {
		return fmt.Sprintf("[%s]", text), nil
	}
	return text, nil
}

func (content *AtlassianContent) convertHeading() (string, error) {
	level := content.Attrs.Level
	if level < 1 {
		level = 1
	} else if level > 6 {
		level = 6
	}

	var heading strings.Builder
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

	return heading.String() + "\n\n", nil
}

func (content *AtlassianContent) convertEmoji() (string, error) {
	return content.Text, nil
}

func (content *AtlassianContent) convertCodeBlock() (string, error) {
	var code strings.Builder
	language := content.Attrs.Language
	if language == "" {
		language = "text"
	}

	code.WriteString("```" + language + "\n")

	for _, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		code.WriteString(text)
	}

	// Ensure the code block ends with a newline
	if !strings.HasSuffix(code.String(), "\n") {
		code.WriteString("\n")
	}

	code.WriteString("```\n\n")
	return code.String(), nil
}

func (content *AtlassianContent) convertTaskList() (string, error) {
	var list strings.Builder
	for _, item := range content.Content {
		if item.Type != "taskItem" {
			continue
		}

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
	state := content.Attrs.State

	// Add checkbox
	if state == "DONE" {
		item.WriteString("- [x] ")
	} else {
		item.WriteString("- [ ] ")
	}

	for _, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		item.WriteString(text)
	}
	item.WriteString("\n")

	return item.String(), nil
}

func (content *AtlassianContent) convertPanel() (string, error) {
	var panel strings.Builder
	panelType := content.Attrs.PanelType

	// Start panel with a horizontal rule and panel type
	panel.WriteString("---\n")
	if panelType != "" {
		panel.WriteString(fmt.Sprintf("**%s**\n\n", strings.ToUpper(panelType)))
	}

	// Convert panel content
	for _, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		panel.WriteString(text)
	}

	// End panel with a horizontal rule
	panel.WriteString("---\n\n")
	return panel.String(), nil
}

func (content *AtlassianContent) convertBodiedExtension() (string, error) {
	// For now, just convert the content
	var ext strings.Builder
	for _, child := range content.Content {
		text, err := child.AtlassianDocumentConvertToMarkdown()
		if err != nil {
			return "", err
		}
		ext.WriteString(text)
	}
	return ext.String(), nil
}

func (content *AtlassianContent) convertExtension() (string, error) {
	// For now, just return empty string
	return "", nil
}

func (content *AtlassianContent) convertDate() (string, error) {
	if content.Attrs.Timestamp != "" {
		return content.Attrs.Timestamp, nil
	}
	return "", nil
}

func (content *AtlassianContent) convertMention() (string, error) {
	if content.Attrs.Text != "" {
		return "@" + content.Attrs.Text, nil
	}
	return "", nil
}

func (content *AtlassianContent) convertPlaceholder() (string, error) {
	if content.Attrs.Text != "" {
		return fmt.Sprintf("_%s_", content.Attrs.Text), nil
	}
	return "", nil
}

func (content *AtlassianContent) convertLayoutSection() (string, error) {
	var section strings.Builder
	for _, column := range content.Content {
		text, err := column.AtlassianDocumentConvertToMarkdown()
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
	return content.convertMedia()
}

func (content *AtlassianContent) convertMedia() (string, error) {
	if content.Attrs.URL != "" {
		return fmt.Sprintf("![Image](%s)\n\n", content.Attrs.URL), nil
	}
	return "", nil
}

func (content *AtlassianContent) convertNestedExpand() (string, error) {
	return content.convertExpand()
}

func (content *AtlassianContent) convertExpand() (string, error) {
	var expand strings.Builder
	title := content.Attrs.Title
	if title == "" {
		title = "Details"
	}

	expand.WriteString(fmt.Sprintf("<details>\n<summary>%s</summary>\n\n", title))

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
	// For now, just return empty string
	return "", nil
}

// ParseDocument parses a JSON string into an AtlassianDocument
func ParseDocument(jsonStr string) (*AtlassianDocument, error) {
	var doc AtlassianDocument
	if err := json.Unmarshal([]byte(jsonStr), &doc); err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}
	return &doc, nil
}

// AtlassianDocumentConvertJSONToMarkdown converts a JSON string representing an Atlassian document to Markdown
func AtlassianDocumentConvertJSONToMarkdown(jsonStr string) (string, error) {
	doc, err := ParseDocument(jsonStr)
	if err != nil {
		return "", err
	}
	return doc.AtlassianDocumentConvertToMarkdown()
}
