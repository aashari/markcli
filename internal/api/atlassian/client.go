package atlassian

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"markcli/internal/logging"
	"net/http"
	"net/url"
)

// Client represents an Atlassian API client
type Client struct {
	baseURL    string
	email      string
	token      string
	httpClient *http.Client
}

// NewClient creates a new Atlassian API client
func NewClient(baseURL, email, token string) *Client {
	return &Client{
		baseURL:    baseURL,
		email:      email,
		token:      token,
		httpClient: &http.Client{},
	}
}

// newRequest creates a new HTTP request with authentication and common headers
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	// Build full URL
	reqURL := fmt.Sprintf("%s%s", c.baseURL, path)

	// Create request body if provided
	var buf io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		buf = bytes.NewBuffer(jsonBody)
	}

	// Create request
	req, err := http.NewRequest(method, reqURL, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set common headers
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// SearchOptions represents options for searching pages
type SearchOptions struct {
	Query     string
	SpaceKey  string
	StartAt   int
	Limit     int
	SortBy    string
	SortOrder string
}

// SearchResponse represents a search response from the Confluence API
type SearchResponse struct {
	Results        []ContentResult `json:"results"`
	Start          int             `json:"start"`
	Limit          int             `json:"limit"`
	Size           int             `json:"size"`
	TotalSize      int             `json:"totalSize"`
	CQLQuery       string          `json:"cqlQuery"`
	SearchDuration int             `json:"searchDuration"`
	Links          struct {
		Base    string `json:"base"`
		Context string `json:"context"`
		Next    string `json:"next"`
		Self    string `json:"self"`
	} `json:"_links"`
}

// ContentResult represents a single search result
type ContentResult struct {
	Content struct {
		ID     string `json:"id"`
		Type   string `json:"type"`
		Status string `json:"status"`
		Title  string `json:"title"`
		Body   struct {
			AtlasDocFormat struct {
				Value string `json:"value"`
			} `json:"atlas_doc_format"`
		} `json:"body"`
		Links struct {
			WebUI  string `json:"webui"`
			Self   string `json:"self"`
			TinyUI string `json:"tinyui"`
		} `json:"_links"`
	} `json:"content"`
	Title                 string `json:"title"`
	Excerpt               string `json:"excerpt"`
	URL                   string `json:"url"`
	ResultGlobalContainer struct {
		Title      string `json:"title"`
		DisplayURL string `json:"displayUrl"`
	} `json:"resultGlobalContainer"`
	LastModified         string `json:"lastModified"`
	FriendlyLastModified string `json:"friendlyLastModified"`
}

// Body represents the content body of a page
type Body struct {
	View struct {
		Value          string `json:"value"`
		Representation string `json:"representation"`
	} `json:"view"`
}

// Version represents the version information of a page
type Version struct {
	Number    int    `json:"number"`
	When      string `json:"when"`
	MinorEdit bool   `json:"minorEdit"`
}

// Space represents a Confluence space
type Space struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

// Page represents a Confluence page
type Page struct {
	ID                   string  `json:"id"`
	Type                 string  `json:"type"`
	Status               string  `json:"status"`
	Title                string  `json:"title"`
	Body                 Body    `json:"body"`
	Space                Space   `json:"space"`
	Version              Version `json:"version"`
	URL                  string
	FriendlyLastModified string
}

// SearchPages searches for pages in Confluence
func (c *Client) SearchPages(opts SearchOptions) (*SearchResponse, error) {
	// Build CQL query
	cql := fmt.Sprintf("type=page AND text ~ \"%s\"", opts.Query)
	if opts.SpaceKey != "" {
		cql = fmt.Sprintf("%s AND space=\"%s\"", cql, opts.SpaceKey)
	}

	// Build query parameters
	params := url.Values{}
	params.Add("cql", cql)
	params.Add("start", fmt.Sprintf("%d", opts.StartAt))
	params.Add("limit", fmt.Sprintf("%d", opts.Limit))
	params.Add("expand", "content.space,content.version")

	// Create request
	req, err := c.newRequest("GET", "/wiki/rest/api/search?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for logging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log request and response for debugging
	logging.LogDebug("Request URL: %s", req.URL.String())
	logging.LogDebug("Response Status: %s", resp.Status)
	if logging.IsDebugEnabled() {
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err == nil {
			logging.LogJSONInline("Response Body", jsonData)
		} else {
			logging.LogDebug("Response Body: %s", string(body))
		}
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ListSpaces returns a list of Confluence spaces
func (c *Client) ListSpaces(includeAll bool) ([]Space, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("limit", "100") // Get maximum spaces per request
	if !includeAll {
		params.Add("type", "global")    // Only get global spaces
		params.Add("status", "current") // Only get active spaces
	}

	// Create request
	req, err := c.newRequest("GET", "/wiki/rest/api/space?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		Results []Space `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Results, nil
}

// GetPage gets a specific page by ID from Confluence API v2
func (c *Client) GetPage(pageID string) ([]byte, error) {
	endpoint := fmt.Sprintf("/wiki/api/v2/pages/%s", pageID)

	// Build query parameters
	params := url.Values{}
	params.Add("body-format", "atlas_doc_format")

	// Add query parameters to endpoint
	if len(params) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	// Create request
	req, err := c.newRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Log request for debugging
	logging.LogDebug("Request URL: %s", req.URL.String())
	logging.LogDebug("Request Headers: %v", req.Header)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log response for debugging
	logging.LogDebug("Response Status: %s", resp.Status)
	logging.LogDebug("Response Headers: %v", resp.Header)
	if logging.IsDebugEnabled() {
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err == nil {
			logging.LogJSONInline("Response Body", jsonData)
		} else {
			logging.LogDebug("Response Body: %s", string(body))
		}
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// makeRequest creates and sends an HTTP request
func (c *Client) makeRequest(method, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.email, c.token)

	logging.LogDebug("Request URL: %s", req.URL.String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	logging.LogDebug("Response Status: %s", resp.Status)
	return resp, nil
}

// FooterComment represents a Confluence footer comment
type FooterComment struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Title   string `json:"title"`
	Version struct {
		CreatedAt string `json:"createdAt"`
		Message   string `json:"message"`
		Number    int    `json:"number"`
		AuthorID  string `json:"authorId"`
	} `json:"version"`
	Body struct {
		Storage struct {
			Value          string `json:"value"`
			Representation string `json:"representation"`
		} `json:"storage"`
		AtlasDocFormat struct {
			Value          string `json:"value"`
			Representation string `json:"representation"`
		} `json:"atlas_doc_format"`
	} `json:"body"`
}

// FooterCommentsResponse represents the response from the footer comments API
type FooterCommentsResponse struct {
	Results []FooterComment `json:"results"`
	Links   struct {
		Next string `json:"next"`
		Base string `json:"base"`
	} `json:"_links"`
}

// GetPageFooterComments retrieves footer comments for a specific page
func (c *Client) GetPageFooterComments(pageID string) (*FooterCommentsResponse, error) {
	endpoint := fmt.Sprintf("/wiki/api/v2/pages/%s/footer-comments", pageID)

	// Build query parameters
	params := make(map[string]string)
	params["body-format"] = "atlas_doc_format"

	// Add query parameters to endpoint
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}
		endpoint = fmt.Sprintf("%s?%s", endpoint, values.Encode())
	}

	// Create request
	req, err := c.newRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Log request for debugging
	logging.LogDebug("Request URL: %s", req.URL.String())
	logging.LogDebug("Request Headers: %v", req.Header)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for logging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log response for debugging
	logging.LogDebug("Response Status: %s", resp.Status)
	logging.LogDebug("Request Headers: %v", req.Header)
	if logging.IsDebugEnabled() {
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err == nil {
			logging.LogJSONInline("Response Body", jsonData)
		} else {
			logging.LogDebug("Response Body: %s", string(body))
		}
	}

	// Check response status
	if resp.StatusCode == http.StatusNotFound {
		// Return empty response for 404 errors
		return &FooterCommentsResponse{
			Results: []FooterComment{},
		}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var result FooterCommentsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// APIError represents an error returned by the Atlassian API
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s (status code: %d)", e.Message, e.StatusCode)
}

// JiraProject represents a Jira project
type JiraProject struct {
	ID             string `json:"id"`
	Key            string `json:"key"`
	Name           string `json:"name"`
	ProjectTypeKey string `json:"projectTypeKey"`
	Style          string `json:"style"`
}

// JiraSearchOptions represents options for searching Jira issues
type JiraSearchOptions struct {
	Query   string
	StartAt int
	Limit   int
}

// JiraSearchResponse represents a search response from the Jira API
type JiraSearchResponse struct {
	Issues     []JiraIssue `json:"issues"`
	MaxResults int         `json:"maxResults"`
	StartAt    int         `json:"startAt"`
	Total      int         `json:"total"`
}

// JiraIssue represents a Jira issue
type JiraIssue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Self   string `json:"self"`
	Fields struct {
		Summary string `json:"summary"`
		Status  struct {
			Name string `json:"name"`
		} `json:"status"`
		Project struct {
			Key  string `json:"key"`
			Name string `json:"name"`
		} `json:"project"`
		Assignee *struct {
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		Description *struct {
			Type    string `json:"type"`
			Version int    `json:"version"`
			Content []struct {
				Type    string `json:"type"`
				Content []struct {
					Type  string `json:"type"`
					Text  string `json:"text,omitempty"`
					Marks []struct {
						Type  string `json:"type"`
						Attrs struct {
							URL string `json:"url,omitempty"`
						} `json:"attrs,omitempty"`
					} `json:"marks,omitempty"`
					Attrs map[string]interface{} `json:"attrs,omitempty"`
				} `json:"content,omitempty"`
			} `json:"content,omitempty"`
		} `json:"description"`
	} `json:"fields"`
}

// ListProjects returns a list of Jira projects
func (c *Client) ListProjects() ([]JiraProject, error) {
	// Create request
	req, err := c.newRequest("GET", "/rest/api/3/project", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for logging and error handling
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log request and response for debugging
	logging.LogDebug("Request URL: %s", req.URL.String())
	logging.LogDebug("Response Status: %s", resp.Status)
	if logging.IsDebugEnabled() {
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err == nil {
			logging.LogJSONInline("Response Body", jsonData)
		} else {
			logging.LogDebug("Response Body: %s", string(body))
		}
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Try to parse error message from response body
		var errorResp struct {
			ErrorMessages []string          `json:"errorMessages"`
			Errors        map[string]string `json:"errors"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			var message string
			if len(errorResp.ErrorMessages) > 0 {
				message = errorResp.ErrorMessages[0]
			} else if len(errorResp.Errors) > 0 {
				// Take the first error message from the map
				for _, v := range errorResp.Errors {
					message = v
					break
				}
			}
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Message:    message,
			}
		}
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    "unknown error",
		}
	}

	// Parse response
	var projects []JiraProject
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return projects, nil
}

// SearchIssues searches for issues in Jira using JQL
func (c *Client) SearchIssues(opts JiraSearchOptions) (*JiraSearchResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("jql", fmt.Sprintf("%s ORDER BY updated DESC", opts.Query))
	params.Add("startAt", fmt.Sprintf("%d", opts.StartAt))
	params.Add("maxResults", fmt.Sprintf("%d", opts.Limit))
	params.Add("fields", "summary,status,project,assignee,description")

	// Create request
	req, err := c.newRequest("GET", "/rest/api/3/search?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for logging and error handling
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log request and response for debugging
	logging.LogDebug("Request URL: %s", req.URL.String())
	logging.LogDebug("Response Status: %s", resp.Status)
	if logging.IsDebugEnabled() {
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err == nil {
			logging.LogJSONInline("Response Body", jsonData)
		} else {
			logging.LogDebug("Response Body: %s", string(body))
		}
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Try to parse error message from response body
		var errorResp struct {
			ErrorMessages []string          `json:"errorMessages"`
			Errors        map[string]string `json:"errors"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			var message string
			if len(errorResp.ErrorMessages) > 0 {
				message = errorResp.ErrorMessages[0]
			} else if len(errorResp.Errors) > 0 {
				// Take the first error message from the map
				for _, v := range errorResp.Errors {
					message = v
					break
				}
			}
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Message:    message,
			}
		}
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    "unknown error",
		}
	}

	// Parse response
	var result JiraSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// JiraIssueDetails represents a full Jira issue details response
type JiraIssueDetails struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Self   string `json:"self"`
	Fields struct {
		Summary string `json:"summary"`
		Status  struct {
			Name string `json:"name"`
		} `json:"status"`
		Project struct {
			Key  string `json:"key"`
			Name string `json:"name"`
		} `json:"project"`
		Assignee *struct {
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		Reporter *struct {
			DisplayName string `json:"displayName"`
		} `json:"reporter"`
		Description *struct {
			Type    string `json:"type"`
			Version int    `json:"version"`
			Content []struct {
				Type    string `json:"type"`
				Content []struct {
					Type  string `json:"type"`
					Text  string `json:"text,omitempty"`
					Marks []struct {
						Type  string `json:"type"`
						Attrs struct {
							URL string `json:"url,omitempty"`
						} `json:"attrs,omitempty"`
					} `json:"marks,omitempty"`
					Attrs map[string]interface{} `json:"attrs,omitempty"`
				} `json:"content,omitempty"`
			} `json:"content,omitempty"`
		} `json:"description"`
		Created string `json:"created"`
		Updated string `json:"updated"`
	} `json:"fields"`
}

// GetIssue gets a specific issue by ID from Jira API v3
func (c *Client) GetIssue(issueID string) (*JiraIssueDetails, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("fields", "summary,status,project,assignee,reporter,description,created,updated")

	// Create request
	req, err := c.newRequest("GET", fmt.Sprintf("/rest/api/3/issue/%s?%s", issueID, params.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for logging and error handling
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log request and response for debugging
	logging.LogDebug("Request URL: %s", req.URL.String())
	logging.LogDebug("Response Status: %s", resp.Status)
	if logging.IsDebugEnabled() {
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err == nil {
			logging.LogJSONInline("Response Body", jsonData)
		} else {
			logging.LogDebug("Response Body: %s", string(body))
		}
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Try to parse error message from response body
		var errorResp struct {
			ErrorMessages []string          `json:"errorMessages"`
			Errors        map[string]string `json:"errors"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			var message string
			if len(errorResp.ErrorMessages) > 0 {
				message = errorResp.ErrorMessages[0]
			} else if len(errorResp.Errors) > 0 {
				// Take the first error message from the map
				for _, v := range errorResp.Errors {
					message = v
					break
				}
			}
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Message:    message,
			}
		}
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    "unknown error",
		}
	}

	// Parse response
	var issue JiraIssueDetails
	if err := json.Unmarshal(body, &issue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &issue, nil
}
