package atlassian

// AtlassianJiraProject represents a Jira project
type AtlassianJiraProject struct {
	ID             string `json:"id"`
	Key            string `json:"key"`
	Name           string `json:"name"`
	ProjectTypeKey string `json:"projectTypeKey"`
	Style          string `json:"style"`
}

// AtlassianJiraSearchOptions represents options for searching Jira issues
type AtlassianJiraSearchOptions struct {
	Query   string
	StartAt int
	Limit   int
}

// AtlassianJiraSearchResponse represents a search response from the Jira API
type AtlassianJiraSearchResponse struct {
	Expand     string               `json:"expand"`
	StartAt    int                  `json:"startAt"`
	MaxResults int                  `json:"maxResults"`
	Total      int                  `json:"total"`
	Issues     []AtlassianJiraIssue `json:"issues"`
}

// AtlassianJiraIssue represents a Jira issue
type AtlassianJiraIssue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Self   string `json:"self"`
	Fields struct {
		Summary     string `json:"summary"`
		Description *struct {
			Type    string             `json:"type"`
			Version int                `json:"version"`
			Content []AtlassianContent `json:"content,omitempty"`
		} `json:"description"`
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
		Priority struct {
			Name string `json:"name"`
		} `json:"priority"`
		Created    string `json:"created"`
		Updated    string `json:"updated"`
		DueDate    string `json:"duedate"`
		Resolution struct {
			Name string `json:"name"`
		} `json:"resolution"`
		Project struct {
			Key  string `json:"key"`
			Name string `json:"name"`
		} `json:"project"`
		IssueType struct {
			Name string `json:"name"`
		} `json:"issuetype"`
		Assignee *struct {
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		Reporter *struct {
			DisplayName string `json:"displayName"`
		} `json:"reporter"`
	} `json:"fields"`
}

// AtlassianJiraError represents an error returned by the Jira API
type AtlassianJiraError struct {
	StatusCode    int
	Message       string
	ErrorMessages []string          `json:"errorMessages"`
	Errors        map[string]string `json:"errors"`
}

func (e *AtlassianJiraError) Error() string {
	if len(e.ErrorMessages) > 0 {
		return e.ErrorMessages[0]
	}
	if len(e.Errors) > 0 {
		for _, v := range e.Errors {
			return v
		}
	}
	return e.Message
}

// AtlassianJiraComment represents a comment on a Jira issue
type AtlassianJiraComment struct {
	ID     string `json:"id"`
	Author *struct {
		DisplayName string `json:"displayName"`
	} `json:"author"`
	Body *struct {
		Type    string             `json:"type"`
		Version int                `json:"version"`
		Content []AtlassianContent `json:"content"`
	} `json:"body"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	JSDPublic bool   `json:"jsdPublic"`
}

// AtlassianJiraCommentsResponse represents a response containing Jira comments
type AtlassianJiraCommentsResponse struct {
	Comments   []AtlassianJiraComment `json:"comments"`
	StartAt    int                    `json:"startAt"`
	MaxResults int                    `json:"maxResults"`
	Total      int                    `json:"total"`
}
