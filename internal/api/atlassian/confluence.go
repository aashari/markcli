package atlassian

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"markcli/internal/logging"
	"markcli/internal/types/atlassian"
)

// AtlassianConfluenceSearchPages searches for pages in Confluence
func (c *Client) AtlassianConfluenceSearchPages(opts atlassian.AtlassianConfluenceSearchOptions) (*atlassian.AtlassianConfluenceSearchResponse, error) {
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
	var result atlassian.AtlassianConfluenceSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// AtlassianConfluenceListSpaces returns a list of Confluence spaces
func (c *Client) AtlassianConfluenceListSpaces(includeAll bool) ([]atlassian.AtlassianConfluenceSpace, error) {
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
		Results []atlassian.AtlassianConfluenceSpace `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Results, nil
}

// AtlassianConfluenceGetPage gets a specific page by ID from Confluence API v2
func (c *Client) AtlassianConfluenceGetPage(pageID string) (*atlassian.AtlassianConfluencePageDetails, error) {
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

	// Parse response
	var result atlassian.AtlassianConfluencePageDetails
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// AtlassianConfluenceGetPageFooterComments retrieves footer comments for a specific page
func (c *Client) AtlassianConfluenceGetPageFooterComments(pageID string) (*atlassian.AtlassianConfluenceFooterCommentsResponse, error) {
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
		return &atlassian.AtlassianConfluenceFooterCommentsResponse{
			Results: []atlassian.AtlassianConfluenceFooterComment{},
		}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var result atlassian.AtlassianConfluenceFooterCommentsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
