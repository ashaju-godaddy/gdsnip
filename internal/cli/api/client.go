package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ashaju-godaddy/gdsnip/internal/cli/config"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/ashaju-godaddy/gdsnip/pkg/template"
)

// Client is the API client for GDSNIP
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient() (*Client, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Load credentials (optional)
	token := ""
	creds, err := config.LoadCredentials()
	if err == nil {
		token = creds.Token
	}

	return &Client{
		baseURL: cfg.APIURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// User represents a user in responses
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Register registers a new user
func (c *Client) Register(email, username, password string) (*AuthResponse, error) {
	payload := map[string]string{
		"email":    email,
		"username": username,
		"password": password,
	}

	var response AuthResponse
	if err := c.doRequest("POST", "/auth/register", payload, &response, false); err != nil {
		return nil, err
	}

	return &response, nil
}

// Login authenticates a user
func (c *Client) Login(email, password string) (*AuthResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	var response AuthResponse
	if err := c.doRequest("POST", "/auth/login", payload, &response, false); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMe returns the current authenticated user
func (c *Client) GetMe() (*User, error) {
	var user User
	if err := c.doRequest("GET", "/auth/me", nil, &user, true); err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateSnippetRequest represents a snippet creation request
type CreateSnippetRequest struct {
	Name        string            `json:"name"`
	Slug        string            `json:"slug,omitempty"`
	Description string            `json:"description,omitempty"`
	Content     string            `json:"content"`
	Tags        []string          `json:"tags,omitempty"`
	Visibility  string            `json:"visibility,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// CreateSnippet creates a new snippet
func (c *Client) CreateSnippet(req *CreateSnippetRequest) (*models.Snippet, error) {
	var snippet models.Snippet
	if err := c.doRequest("POST", "/snippets", req, &snippet, true); err != nil {
		return nil, err
	}

	return &snippet, nil
}

// GetSnippet retrieves a snippet by namespace and slug
func (c *Client) GetSnippet(namespace, slug string) (*models.Snippet, error) {
	var snippet models.Snippet
	path := fmt.Sprintf("/snippets/%s/%s", namespace, slug)
	if err := c.doRequest("GET", path, nil, &snippet, false); err != nil {
		return nil, err
	}

	return &snippet, nil
}

// PullRequest represents a pull request
type PullRequest struct {
	Variables map[string]string `json:"variables"`
}

// PullResponse represents a pull response
type PullResponse struct {
	Content   string              `json:"content"`
	Snippet   *models.Snippet     `json:"snippet"`
	Variables []template.Variable `json:"variables"`
	Warnings  []string            `json:"warnings,omitempty"`
}

// PullSnippet pulls and renders a snippet
func (c *Client) PullSnippet(namespace, slug string, variables map[string]string) (*PullResponse, error) {
	path := fmt.Sprintf("/snippets/%s/%s/pull", namespace, slug)
	req := PullRequest{Variables: variables}

	var response PullResponse
	if err := c.doRequest("POST", path, req, &response, false); err != nil {
		return nil, err
	}

	return &response, nil
}

// SearchResponse represents search results
type SearchResponse struct {
	Snippets   []models.Snippet `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Search searches for public snippets
func (c *Client) Search(query string, tags []string, limit int) (*SearchResponse, error) {
	path := fmt.Sprintf("/snippets?q=%s&limit=%d", query, limit)
	for _, tag := range tags {
		path += fmt.Sprintf("&tags=%s", tag)
	}

	var response SearchResponse
	if err := c.doRequest("GET", path, nil, &response, false); err != nil {
		return nil, err
	}

	return &response, nil
}

// ListResponse represents list results
type ListResponse struct {
	Snippets   []models.Snippet `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

// ListMySnippets lists the current user's snippets
func (c *Client) ListMySnippets(limit, offset int) (*ListResponse, error) {
	path := fmt.Sprintf("/users/me/snippets?limit=%d&offset=%d", limit, offset)

	var response ListResponse
	if err := c.doRequest("GET", path, nil, &response, true); err != nil {
		return nil, err
	}

	return &response, nil
}

// doRequest performs an HTTP request
func (c *Client) doRequest(method, path string, payload interface{}, result interface{}, requireAuth bool) error {
	// Build full URL
	url := c.baseURL + path

	// Prepare request body
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	// Create request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add auth header if available or required
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	} else if requireAuth {
		return fmt.Errorf("not logged in: please run 'gdsnip auth login' first")
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode >= 400 {
		return parseError(respBody)
	}

	// Parse success response
	var envelope struct {
		Success    bool            `json:"success"`
		Data       json.RawMessage `json:"data"`
		Pagination *Pagination     `json:"pagination,omitempty"`
	}

	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Unmarshal data into result
	if result != nil {
		// Check if result is a struct with pagination
		switch v := result.(type) {
		case *SearchResponse:
			if err := json.Unmarshal(envelope.Data, &v.Snippets); err != nil {
				return fmt.Errorf("failed to parse response data: %w", err)
			}
			if envelope.Pagination != nil {
				v.Pagination = *envelope.Pagination
			}
		case *ListResponse:
			if err := json.Unmarshal(envelope.Data, &v.Snippets); err != nil {
				return fmt.Errorf("failed to parse response data: %w", err)
			}
			if envelope.Pagination != nil {
				v.Pagination = *envelope.Pagination
			}
		default:
			if err := json.Unmarshal(envelope.Data, result); err != nil {
				return fmt.Errorf("failed to parse response data: %w", err)
			}
		}
	}

	return nil
}

// parseError parses an error response
func parseError(body []byte) error {
	var envelope struct {
		Success bool          `json:"success"`
		Error   *models.APIError `json:"error"`
	}

	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("API error (failed to parse): %s", string(body))
	}

	if envelope.Error != nil {
		return envelope.Error
	}

	return fmt.Errorf("unknown API error")
}
