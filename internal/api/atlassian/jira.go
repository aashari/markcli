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

// AtlassianJiraListProjects returns a list of Jira projects
func (c *Client) AtlassianJiraListProjects() ([]atlassian.AtlassianJiraProject, error) {
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
		var errorResp atlassian.AtlassianJiraError
		if err := json.Unmarshal(body, &errorResp); err == nil {
			errorResp.StatusCode = resp.StatusCode
			logging.LogDebug("Jira API Error - Status Code: %d, Error: %+v", resp.StatusCode, errorResp)
			return nil, &errorResp
		}
		logging.LogDebug("Jira API Error - Status Code: %d, Body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var projects []atlassian.AtlassianJiraProject
	if err := json.Unmarshal(body, &projects); err != nil {
		logging.LogDebug("Failed to decode response: %s", string(body))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return projects, nil
}

// AtlassianJiraSearchIssues searches for issues in Jira using JQL
func (c *Client) AtlassianJiraSearchIssues(opts atlassian.AtlassianJiraSearchOptions) (*atlassian.AtlassianJiraSearchResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("jql", opts.Query)
	params.Add("startAt", fmt.Sprintf("%d", opts.StartAt))
	params.Add("maxResults", fmt.Sprintf("%d", opts.Limit))
	params.Add("fields", "summary,status,priority,project,assignee,reporter,description,created,updated,duedate,resolution,issuetype")

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
		var errorResp atlassian.AtlassianJiraError
		if err := json.Unmarshal(body, &errorResp); err == nil {
			errorResp.StatusCode = resp.StatusCode
			logging.LogDebug("Jira API Error - Status Code: %d, Error: %+v", resp.StatusCode, errorResp)
			return nil, &errorResp
		}
		logging.LogDebug("Jira API Error - Status Code: %d, Body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result atlassian.AtlassianJiraSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		logging.LogDebug("Failed to decode response: %s", string(body))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// AtlassianJiraGetIssue gets a specific issue by ID from Jira API v3
func (c *Client) AtlassianJiraGetIssue(issueID string) (*atlassian.AtlassianJiraIssue, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("fields", "summary,status,priority,project,assignee,reporter,description,created,updated,duedate,resolution,issuetype")

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
		var errorResp atlassian.AtlassianJiraError
		if err := json.Unmarshal(body, &errorResp); err == nil {
			errorResp.StatusCode = resp.StatusCode
			logging.LogDebug("Jira API Error - Status Code: %d, Error: %+v", resp.StatusCode, errorResp)
			return nil, &errorResp
		}
		logging.LogDebug("Jira API Error - Status Code: %d, Body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var issue atlassian.AtlassianJiraIssue
	if err := json.Unmarshal(body, &issue); err != nil {
		logging.LogDebug("Failed to decode response: %s", string(body))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &issue, nil
}

// AtlassianJiraGetIssueComments gets comments for a specific issue
func (c *Client) AtlassianJiraGetIssueComments(issueID string) (*atlassian.AtlassianJiraCommentsResponse, error) {
	// Create request
	req, err := c.newRequest("GET", fmt.Sprintf("/rest/api/3/issue/%s/comment", issueID), nil)
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
		var errorResp atlassian.AtlassianJiraError
		if err := json.Unmarshal(body, &errorResp); err == nil {
			errorResp.StatusCode = resp.StatusCode
			logging.LogDebug("Jira API Error - Status Code: %d, Error: %+v", resp.StatusCode, errorResp)
			return nil, &errorResp
		}
		logging.LogDebug("Jira API Error - Status Code: %d, Body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result atlassian.AtlassianJiraCommentsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		logging.LogDebug("Failed to decode response: %s", string(body))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
