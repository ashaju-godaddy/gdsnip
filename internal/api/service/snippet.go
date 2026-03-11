package service

import (
	"fmt"

	"github.com/ashaju-godaddy/gdsnip/internal/api/repository"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/ashaju-godaddy/gdsnip/pkg/template"
	"github.com/ashaju-godaddy/gdsnip/pkg/validator"
	"github.com/google/uuid"
)

// SnippetService handles snippet business logic
type SnippetService struct {
	snippetRepo *repository.SnippetRepo
	userRepo    *repository.UserRepo
	teamRepo    *repository.TeamRepo
}

// NewSnippetService creates a new snippet service
func NewSnippetService(snippetRepo *repository.SnippetRepo, userRepo *repository.UserRepo, teamRepo *repository.TeamRepo) *SnippetService {
	return &SnippetService{
		snippetRepo: snippetRepo,
		userRepo:    userRepo,
		teamRepo:    teamRepo,
	}
}

// CreateRequest represents snippet creation request
type CreateRequest struct {
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Content     string            `json:"content"`
	Tags        []string          `json:"tags"`
	Visibility  string            `json:"visibility"`
	Metadata    map[string]string `json:"metadata"`
	TeamSlug    string            `json:"team_slug,omitempty"` // Optional: push to team namespace
}

// PullRequest represents snippet pull request
type PullRequest struct {
	Variables map[string]string `json:"variables"`
}

// PullResponse represents snippet pull response
type PullResponse struct {
	Content   string              `json:"content"`
	Snippet   *models.Snippet     `json:"snippet"`
	Variables []template.Variable `json:"variables"`
	Warnings  []string            `json:"warnings,omitempty"`
}

// Create creates a new snippet
func (s *SnippetService) Create(userID string, req CreateRequest) (*models.Snippet, error) {
	// Validate name
	if req.Name == "" {
		return nil, models.NewValidationError("name is required", nil)
	}

	// Validate content
	if req.Content == "" {
		return nil, models.NewValidationError("content is required", nil)
	}

	// Auto-generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = validator.GenerateSlug(req.Name)
	} else {
		// Validate provided slug
		if err := validator.ValidateSlug(slug); err != nil {
			return nil, models.NewValidationError(err.Error(), nil)
		}
	}

	// Extract variables from content
	extractedVarNames := template.ExtractVariables(req.Content)

	// Convert extracted variable names to Variable objects
	var extractedVars []template.Variable
	for _, name := range extractedVarNames {
		v := template.Variable{
			Name:     name,
			Required: true, // Default to required
		}
		// Apply metadata if provided
		if desc, ok := req.Metadata[name]; ok {
			v.Description = desc
		}
		extractedVars = append(extractedVars, v)
	}

	// Determine owner type, owner ID, namespace, and default visibility
	var ownerType, ownerID, namespace string
	visibility := req.Visibility

	if req.TeamSlug != "" {
		// Team snippet
		team, err := s.teamRepo.GetBySlug(req.TeamSlug)
		if err != nil {
			return nil, err
		}

		// Check permission to create snippets in this team
		role, err := s.teamRepo.GetUserRole(team.ID, userID)
		if err != nil {
			return nil, models.NewForbiddenError("you are not a member of this team")
		}
		if !role.HasPermission(models.PermCreateSnippet) {
			return nil, models.NewForbiddenError("you do not have permission to create snippets in this team")
		}

		ownerType = "team"
		ownerID = team.ID
		namespace = team.Slug

		// Team snippets default to "team" visibility
		if visibility == "" {
			visibility = "team"
		}

		// Validate visibility for team snippets
		if visibility != "public" && visibility != "private" && visibility != "team" {
			return nil, models.NewValidationError("visibility must be 'public', 'private', or 'team'", nil)
		}
	} else {
		// User snippet
		user, err := s.userRepo.GetByID(userID)
		if err != nil {
			return nil, err
		}
		ownerType = "user"
		ownerID = userID
		namespace = user.Username

		// Default visibility to private for user snippets
		if visibility == "" {
			visibility = "private"
		}

		// Validate visibility for user snippets
		if visibility != "public" && visibility != "private" {
			return nil, models.NewValidationError("visibility must be 'public' or 'private'", nil)
		}
	}

	// Create snippet
	snippet := &models.Snippet{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		Content:     req.Content,
		Variables:   extractedVars,
		Tags:        req.Tags,
		Visibility:  visibility,
		OwnerType:   ownerType,
		OwnerID:     ownerID,
		Namespace:   namespace,
		CreatedBy:   userID,
		Version:     1,
		PullCount:   0,
	}

	if err := s.snippetRepo.Create(snippet); err != nil {
		return nil, err
	}

	return snippet, nil
}

// GetByPath retrieves a snippet by namespace and slug
func (s *SnippetService) GetByPath(namespace, slug string) (*models.Snippet, error) {
	return s.snippetRepo.GetByNamespaceSlug(namespace, slug)
}

// Pull renders a snippet with provided variables
func (s *SnippetService) Pull(namespace, slug, userID string, variables map[string]string) (*PullResponse, error) {
	// Get snippet
	snippet, err := s.snippetRepo.GetByNamespaceSlug(namespace, slug)
	if err != nil {
		return nil, err
	}

	// Check visibility
	switch snippet.Visibility {
	case "private":
		// Private snippets can only be accessed by owner
		if snippet.OwnerID != userID {
			return nil, models.NewNotFoundError("snippet not found")
		}
	case "team":
		// Team-visibility snippets can be accessed by team members
		if snippet.OwnerType == "team" {
			isMember, err := s.teamRepo.IsMember(snippet.OwnerID, userID)
			if err != nil {
				return nil, models.NewInternalError("failed to check team membership")
			}
			if !isMember {
				return nil, models.NewNotFoundError("snippet not found")
			}
		} else {
			// User-owned snippet with team visibility (shouldn't happen, but handle it)
			if snippet.OwnerID != userID {
				return nil, models.NewNotFoundError("snippet not found")
			}
		}
	case "public":
		// Public snippets are accessible to everyone
	default:
		// Unknown visibility, deny access
		return nil, models.NewNotFoundError("snippet not found")
	}

	// Validate provided variables
	if err := template.Validate(snippet.Variables, variables); err != nil {
		// Get missing variables for detailed error response
		missing := template.GetMissingVariables(snippet.Variables, variables)
		return nil, models.NewMissingVariablesError(err.Error(), map[string]interface{}{
			"missing": missing,
		})
	}

	// Render template
	renderResult, err := template.Render(snippet.Content, snippet.Variables, variables)
	if err != nil {
		return nil, models.NewInternalError(fmt.Sprintf("failed to render template: %v", err))
	}

	// Increment pull count in background
	go func() {
		_ = s.snippetRepo.IncrementPullCount(snippet.ID)
	}()

	return &PullResponse{
		Content:   renderResult.Content,
		Snippet:   snippet,
		Variables: snippet.Variables,
		Warnings:  renderResult.Warnings,
	}, nil
}

// Search searches for public snippets
func (s *SnippetService) Search(query string, tags []string, limit, offset int) ([]models.Snippet, int, error) {
	// Apply defaults
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.snippetRepo.Search(query, tags, limit, offset)
}

// ListMine lists snippets owned by the current user
func (s *SnippetService) ListMine(userID string, limit, offset int) ([]models.Snippet, int, error) {
	// Apply defaults
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.snippetRepo.ListByOwner("user", userID, limit, offset)
}

// ListTeamSnippets lists snippets owned by a team
func (s *SnippetService) ListTeamSnippets(userID, teamSlug string, limit, offset int) ([]models.Snippet, int, error) {
	// Get team
	team, err := s.teamRepo.GetBySlug(teamSlug)
	if err != nil {
		return nil, 0, err
	}

	// Check membership
	_, err = s.teamRepo.GetUserRole(team.ID, userID)
	if err != nil {
		return nil, 0, models.NewForbiddenError("you are not a member of this team")
	}

	// Apply defaults
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.snippetRepo.ListByOwner("team", team.ID, limit, offset)
}
