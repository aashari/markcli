package atlassian

import (
	"fmt"
	"time"
)

// AtlassianConfluenceSearchOptions represents options for searching pages in Confluence.
// These options are used to filter and paginate search results.
type AtlassianConfluenceSearchOptions struct {
	Query     string `json:"query"`     // The search query text
	SpaceKey  string `json:"spaceKey"`  // Optional space key to limit search to
	StartAt   int    `json:"startAt"`   // Pagination start index (0-based)
	Limit     int    `json:"limit"`     // Maximum number of results to return
	SortBy    string `json:"sortBy"`    // Field to sort by (e.g., "modified", "created")
	SortOrder string `json:"sortOrder"` // Sort direction ("asc" or "desc")
}

// AtlassianConfluenceLinks represents common link attributes in Confluence responses
type AtlassianConfluenceLinks struct {
	Base    string `json:"base,omitempty"`
	Context string `json:"context,omitempty"`
	Next    string `json:"next,omitempty"`
	Self    string `json:"self,omitempty"`
	WebUI   string `json:"webui,omitempty"`
}

// AtlassianConfluenceSearchResponse represents a search response from the Confluence API.
// It contains the search results along with pagination information.
type AtlassianConfluenceSearchResponse struct {
	Results        []AtlassianConfluenceContentResult `json:"results"`        // List of search results
	Start          int                                `json:"start"`          // Starting index of this page
	Limit          int                                `json:"limit"`          // Maximum number of results per page
	Size           int                                `json:"size"`           // Actual number of results in this page
	TotalSize      int                                `json:"totalSize"`      // Total number of results across all pages
	CQLQuery       string                             `json:"cqlQuery"`       // The CQL query that was executed
	SearchDuration int                                `json:"searchDuration"` // Time taken to execute the search in ms
	Links          AtlassianConfluenceLinks           `json:"_links"`         // Navigation links for pagination
}

// AtlassianConfluenceBody represents the content body in Confluence.
// The body can be in different representations:
// - view: HTML representation
// - storage: Legacy storage format
// - atlas_doc_format: Modern ADF format
type AtlassianConfluenceBody struct {
	View struct {
		Value          string `json:"value"`
		Representation string `json:"representation,omitempty"`
	} `json:"view,omitempty"`
	Storage struct {
		Value          string `json:"value"`
		Representation string `json:"representation,omitempty"`
	} `json:"storage,omitempty"`
	AtlasDocFormat struct {
		Value          string `json:"value"`
		Representation string `json:"representation,omitempty"`
	} `json:"atlas_doc_format,omitempty"`
}

// AtlassianConfluenceContent represents the content structure in search results
type AtlassianConfluenceContent struct {
	ID     string                   `json:"id"`
	Type   string                   `json:"type"`
	Status string                   `json:"status"`
	Title  string                   `json:"title"`
	Body   AtlassianConfluenceBody  `json:"body,omitempty"`
	Links  AtlassianConfluenceLinks `json:"_links,omitempty"`
}

// AtlassianConfluenceResultContainer represents the global container in search results.
// This type provides context about where the result was found.
type AtlassianConfluenceResultContainer struct {
	Title      string `json:"title"`      // Title of the container (e.g., space name)
	DisplayURL string `json:"displayUrl"` // URL to view the container in the browser
}

// AtlassianConfluenceContentResult represents a single search result from Confluence.
// It contains both the content metadata and the actual content data.
type AtlassianConfluenceContentResult struct {
	Content               AtlassianConfluenceContent         `json:"content"`               // The actual content that was found
	Title                 string                             `json:"title"`                 // Title of the content
	Excerpt               string                             `json:"excerpt,omitempty"`     // Snippet of text showing the match
	URL                   string                             `json:"url,omitempty"`         // URL to view the content
	ResultGlobalContainer AtlassianConfluenceResultContainer `json:"resultGlobalContainer"` // Container information
	LastModified          string                             `json:"lastModified"`          // ISO timestamp of last modification
	FriendlyLastModified  string                             `json:"friendlyLastModified"`  // Human-readable last modified date
}

// AtlassianConfluenceVersion represents the version information of a Confluence page or comment.
// It tracks the version number, timestamp, and whether it was a minor edit.
type AtlassianConfluenceVersion struct {
	Number    int                       `json:"number"`              // Version number, increments with each edit
	Message   string                    `json:"message,omitempty"`   // Optional message describing the changes
	When      string                    `json:"when,omitempty"`      // ISO timestamp of when this version was created
	MinorEdit bool                      `json:"minorEdit,omitempty"` // Whether this was a minor edit
	CreatedAt time.Time                 `json:"createdAt,omitempty"` // Timestamp when this version was created
	Author    AtlassianConfluenceAuthor `json:"author,omitempty"`    // User who created this version
	AuthorID  string                    `json:"authorId,omitempty"`  // ID of the user who created this version
}

// AtlassianConfluenceSpace represents a Confluence space.
// A space is a container for pages and other content.
type AtlassianConfluenceSpace struct {
	Key    string `json:"key"`    // Unique identifier for the space
	Name   string `json:"name"`   // Display name of the space
	Type   string `json:"type"`   // Type of space (e.g., "global", "personal")
	Status string `json:"status"` // Status of the space (e.g., "current", "archived")
}

// AtlassianConfluencePage represents a Confluence page.
// It contains all the metadata and content of a page.
type AtlassianConfluencePage struct {
	ID                   string                     `json:"id"`
	Type                 string                     `json:"type"`
	Status               string                     `json:"status"`
	Title                string                     `json:"title"`
	Body                 AtlassianConfluenceBody    `json:"body,omitempty"`
	Space                AtlassianConfluenceSpace   `json:"space,omitempty"`
	Version              AtlassianConfluenceVersion `json:"version,omitempty"`
	URL                  string                     `json:"url,omitempty"`
	FriendlyLastModified string                     `json:"friendlyLastModified,omitempty"`
}

// AtlassianConfluenceAuthor represents an author in Confluence
type AtlassianConfluenceAuthor struct {
	AccountID   string `json:"accountId"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName"`
}

// AtlassianConfluenceFooterComment represents a Confluence footer comment.
// Comments can be in either storage format (legacy) or atlas_doc_format.
type AtlassianConfluenceFooterComment struct {
	ID      string                     `json:"id"`
	Status  string                     `json:"status"`
	Title   string                     `json:"title"`
	Version AtlassianConfluenceVersion `json:"version"`
	Body    AtlassianConfluenceBody    `json:"body"`
}

// AtlassianConfluenceFooterCommentsResponse represents the response from the footer comments API.
// It contains a list of comments and pagination links.
type AtlassianConfluenceFooterCommentsResponse struct {
	Results []AtlassianConfluenceFooterComment `json:"results"`
	Links   AtlassianConfluenceLinks           `json:"_links"`
}

// AtlassianConfluenceMacroValue represents a macro value in Confluence.
// This type is used to store dynamic values for macros, which can be of any type.
type AtlassianConfluenceMacroValue struct {
	Value interface{} `json:"value,omitempty"` // The actual value, can be any valid JSON type
}

// AtlassianConfluenceMacroParams represents macro parameters in Confluence.
// This type defines the common parameters used by Confluence macros.
// The Other field captures any additional parameters not explicitly defined.
type AtlassianConfluenceMacroParams struct {
	Root       *AtlassianConfluenceMacroValue `json:"root,omitempty"`       // Root parameter for tree-like macros
	Spaces     *AtlassianConfluenceMacroValue `json:"spaces,omitempty"`     // Space keys for space-scoped macros
	StartDepth *AtlassianConfluenceMacroValue `json:"startDepth,omitempty"` // Starting depth for tree-like macros
	Excerpt    *AtlassianConfluenceMacroValue `json:"excerpt,omitempty"`    // Excerpt settings for content macros
	Other      map[string]interface{}         `json:"-"`                    // Additional parameters not listed above
}

// AtlassianConfluenceMacroMetadata represents macro metadata in Confluence.
// This type contains metadata about the macro itself, such as its ID and version.
// The Other field captures any additional metadata not explicitly defined.
type AtlassianConfluenceMacroMetadata struct {
	MacroId       *AtlassianConfluenceMacroValue `json:"macroId,omitempty"`
	SchemaVersion *AtlassianConfluenceMacroValue `json:"schemaVersion,omitempty"`
	Title         *AtlassianConfluenceMacroValue `json:"title,omitempty"`
	Other         map[string]interface{}         `json:"-"`
}

// AtlassianConfluenceParameters represents parameters in Confluence.
// This type is a container for both macro parameters and metadata.
// The Other field captures any additional parameters not explicitly defined.
type AtlassianConfluenceParameters struct {
	MacroParams   *AtlassianConfluenceMacroParams   `json:"macroParams,omitempty"`
	MacroMetadata *AtlassianConfluenceMacroMetadata `json:"macroMetadata,omitempty"`
	Other         map[string]interface{}            `json:"-"`
}

// AtlassianConfluencePageDetails represents the response from the Confluence API v2 for a single page
type AtlassianConfluencePageDetails struct {
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
	Comments *AtlassianConfluenceFooterCommentsResponse `json:"comments,omitempty"`
}

// AtlassianConfluenceError represents an error response from the Confluence API
type AtlassianConfluenceError struct {
	StatusCode   int    `json:"-"`
	Message      string `json:"message"`
	Reason       string `json:"reason"`
	ErrorMessage string `json:"error,omitempty"`
	ErrorDetails string `json:"error_description,omitempty"`
	Context      string `json:"context,omitempty"`
	Status       int    `json:"status,omitempty"`
}

// Error implements the error interface
func (e *AtlassianConfluenceError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("Confluence API error: %s", e.Message)
	}
	if e.ErrorMessage != "" {
		return fmt.Sprintf("Confluence API error: %s - %s", e.ErrorMessage, e.ErrorDetails)
	}
	if e.Reason != "" {
		return fmt.Sprintf("Confluence API error: %s", e.Reason)
	}
	return fmt.Sprintf("Confluence API error: status code %d", e.StatusCode)
}
