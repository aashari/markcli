package atlassian

import (
	"encoding/json"
	"time"
)

// AtlassianLinks represents common link attributes in Atlassian responses
type AtlassianLinks struct {
	Base    string `json:"base,omitempty"`
	Context string `json:"context,omitempty"`
	Next    string `json:"next,omitempty"`
	Self    string `json:"self,omitempty"`
	WebUI   string `json:"webui,omitempty"`
}

// AtlassianDocument represents the root structure of an Atlassian Document Format (ADF) document
type AtlassianDocument struct {
	Type    string             `json:"type"`
	Content []AtlassianContent `json:"content"`
	Version int                `json:"version"`
}

// AtlassianContent represents a content node in the Atlassian document.
type AtlassianContent struct {
	Type    string              `json:"type"`
	Content []AtlassianContent  `json:"content,omitempty"`
	Text    string              `json:"text,omitempty"`
	Marks   []AtlassianMark     `json:"marks,omitempty"`
	Attrs   AtlassianAttributes `json:"attrs,omitempty"`
}

// AtlassianMark represents text formatting in the Atlassian document.
type AtlassianMark struct {
	Type  string `json:"type"`
	Attrs struct {
		URL   string `json:"url,omitempty"`
		Color string `json:"color,omitempty"`
	} `json:"attrs,omitempty"`
}

// AtlassianAttributes represents the attributes of an Atlassian content node.
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

// AtlassianListResponse represents a generic list response with pagination.
type AtlassianListResponse struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
}

// AtlassianSearchOptions represents generic search options
type AtlassianSearchOptions struct {
	Query   string `json:"query"`
	StartAt int    `json:"startAt"`
	Limit   int    `json:"limit"`
}

// AtlassianAuthor represents an author in Atlassian
type AtlassianAuthor struct {
	AccountID   string `json:"accountId"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName"`
}

// AtlassianVersion represents a generic version
type AtlassianVersion struct {
	Number    int             `json:"number"`
	Message   string          `json:"message,omitempty"`
	CreatedAt time.Time       `json:"createdAt,omitempty"`
	Author    AtlassianAuthor `json:"author,omitempty"`
}
