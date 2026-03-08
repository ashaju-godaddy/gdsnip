package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/ashaju-godaddy/gdsnip/pkg/template"
)

// Snippet represents a code/config template snippet
type Snippet struct {
	ID           string              `json:"id" db:"id"`
	Name         string              `json:"name" db:"name"`
	Slug         string              `json:"slug" db:"slug"`
	Description  string              `json:"description" db:"description"`
	Content      string              `json:"content" db:"content"`
	Variables    []template.Variable `json:"variables" db:"variables"`
	Tags         []string            `json:"tags" db:"tags"`
	Visibility   string              `json:"visibility" db:"visibility"` // public, private, team
	OwnerType    string              `json:"owner_type" db:"owner_type"` // user, team
	OwnerID      string              `json:"owner_id" db:"owner_id"`
	Namespace    string              `json:"namespace" db:"namespace"` // username or team slug
	CreatedBy    string              `json:"created_by" db:"created_by"`
	Version      int                 `json:"version" db:"version"`
	PullCount    int                 `json:"pull_count" db:"pull_count"`
	SearchVector *string             `json:"-" db:"search_vector"` // PostgreSQL tsvector, not exposed in JSON
	CreatedAt    time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at" db:"updated_at"`
}

// FullPath returns the full path in namespace/slug format
func (s *Snippet) FullPath() string {
	return fmt.Sprintf("%s/%s", s.Namespace, s.Slug)
}

// PullCommand returns an example CLI command to pull this snippet.
// Only includes required variables that don't have defaults.
func (s *Snippet) PullCommand() string {
	cmd := fmt.Sprintf("gdsnip pull %s", s.FullPath())

	// Add required variables without defaults
	for _, v := range s.Variables {
		if v.Required && v.Default == "" {
			cmd += fmt.Sprintf(" --%s=<value>", v.Name)
		}
	}

	return cmd
}

// SnippetSummary is a lightweight version of Snippet for list/search results
type SnippetSummary struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Namespace   string              `json:"namespace"`
	Description string              `json:"description"`
	Tags        []string            `json:"tags"`
	Visibility  string              `json:"visibility"`
	Version     int                 `json:"version"`
	PullCount   int                 `json:"pull_count"`
	Variables   []template.Variable `json:"variables"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// Summary converts a Snippet to SnippetSummary (excludes content)
func (s *Snippet) Summary() SnippetSummary {
	return SnippetSummary{
		ID:          s.ID,
		Name:        s.Name,
		Slug:        s.Slug,
		Namespace:   s.Namespace,
		Description: s.Description,
		Tags:        s.Tags,
		Visibility:  s.Visibility,
		Version:     s.Version,
		PullCount:   s.PullCount,
		Variables:   s.Variables,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// CreateSnippetRequest represents the request body for creating a snippet
type CreateSnippetRequest struct {
	Name        string              `json:"name" validate:"required"`
	Slug        string              `json:"slug"`
	Description string              `json:"description"`
	Content     string              `json:"content" validate:"required"`
	Variables   []template.Variable `json:"variables"`
	Tags        []string            `json:"tags"`
	Visibility  string              `json:"visibility"` // public, private, team
	TeamSlug    string              `json:"team_slug"`  // For team snippets
}

// UpdateSnippetRequest represents the request body for updating a snippet
type UpdateSnippetRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Content     *string             `json:"content"`
	Variables   []template.Variable `json:"variables"`
	Tags        []string            `json:"tags"`
	Visibility  *string             `json:"visibility"`
}

// PullSnippetRequest represents the request body for pulling/rendering a snippet
type PullSnippetRequest struct {
	Variables map[string]string `json:"variables"`
}

// PullSnippetResponse represents the response for a pull request
type PullSnippetResponse struct {
	Snippet            string   `json:"snippet"`             // Full path (namespace/slug)
	Version            int      `json:"version"`             // Snippet version
	Rendered           string   `json:"rendered"`            // Rendered content
	FilenameSuggestion string   `json:"filename_suggestion"` // Suggested filename
	Warnings           []string `json:"warnings,omitempty"`  // Any warnings from rendering
}

// SearchSnippetsQuery represents search parameters
type SearchSnippetsQuery struct {
	Query      string   `query:"q"`
	Tags       []string `query:"tags"`
	Visibility string   `query:"visibility"` // public, private, all (for authenticated users)
	Limit      int      `query:"limit"`
	Offset     int      `query:"offset"`
}

// ListSnippetsResponse represents a paginated list of snippets
type ListSnippetsResponse struct {
	Snippets []SnippetSummary `json:"snippets"`
	Total    int              `json:"total"`
	Limit    int              `json:"limit"`
	Offset   int              `json:"offset"`
}

// GuessFilename attempts to suggest a filename based on snippet content and tags
func (s *Snippet) GuessFilename() string {
	// Check tags for common file types
	for _, tag := range s.Tags {
		tag = strings.ToLower(tag)
		switch tag {
		case "docker", "docker-compose":
			return "docker-compose.yml"
		case "dockerfile":
			return "Dockerfile"
		case "kubernetes", "k8s":
			return s.Slug + ".yaml"
		case "nginx":
			return "nginx.conf"
		case "terraform":
			return "main.tf"
		case "github-actions", "github":
			return ".github/workflows/" + s.Slug + ".yml"
		case "makefile":
			return "Makefile"
		case "go":
			return s.Slug + ".go"
		case "python":
			return s.Slug + ".py"
		case "javascript", "js":
			return s.Slug + ".js"
		case "typescript", "ts":
			return s.Slug + ".ts"
		}
	}

	// Check content for shebang or recognizable patterns
	lines := strings.Split(s.Content, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		if strings.HasPrefix(firstLine, "#!") {
			if strings.Contains(firstLine, "bash") || strings.Contains(firstLine, "sh") {
				return s.Slug + ".sh"
			}
			if strings.Contains(firstLine, "python") {
				return s.Slug + ".py"
			}
		}
		if strings.HasPrefix(firstLine, "version:") || strings.HasPrefix(firstLine, "services:") {
			return "docker-compose.yml"
		}
		if strings.HasPrefix(firstLine, "apiVersion:") {
			return s.Slug + ".yaml"
		}
	}

	// Default: use slug with no extension
	return s.Slug
}
