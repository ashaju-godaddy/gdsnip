package service

import (
	"fmt"
	"testing"

	"github.com/ashaju-godaddy/gdsnip/internal/api/repository"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnippetService_Create(t *testing.T) {
	userRepo, cleanup := setupTestDB(t)
	if cleanup == nil {
		return
	}
	defer cleanup()

	// Create snippet repo
	db, _ := repository.NewDB("postgres://gdsnip:gdsnip@localhost:5432/gdsnip?sslmode=disable")
	snippetRepo := repository.NewSnippetRepo(db)
	defer db.Close()

	snippetService := NewSnippetService(snippetRepo, userRepo)
	authService := NewAuthService(userRepo, "test-secret-key-that-is-very-long", 24*3600*1000000000)

	// Create a test user
	user, _, err := authService.Register("snippet@example.com", "snippetuser", "password123")
	require.NoError(t, err)

	t.Run("create snippet with auto-slug", func(t *testing.T) {
		req := CreateRequest{
			Name:        "Docker Compose PostgreSQL",
			Content:     "version: \"3.8\"\nservices:\n  postgres:\n    image: postgres:{{VERSION}}\n    environment:\n      POSTGRES_PASSWORD: {{DB_PASSWORD}}",
			Visibility:  "public",
			Tags:        []string{"docker", "postgres"},
			Description: "PostgreSQL Docker Compose template",
		}

		snippet, err := snippetService.Create(user.ID, req)
		require.NoError(t, err)
		require.NotNil(t, snippet)

		assert.Equal(t, "Docker Compose PostgreSQL", snippet.Name)
		assert.Equal(t, "docker-compose-postgresql", snippet.Slug)
		assert.Equal(t, "snippetuser", snippet.Namespace)
		assert.Equal(t, "public", snippet.Visibility)
		assert.Equal(t, 2, len(snippet.Variables))
		assert.Equal(t, 2, len(snippet.Tags))
		assert.Equal(t, 1, snippet.Version)
		assert.Equal(t, 0, snippet.PullCount)
	})

	t.Run("create snippet with custom slug", func(t *testing.T) {
		req := CreateRequest{
			Name:       "My Custom Template",
			Slug:       "custom-slug",
			Content:    "echo {{MESSAGE}}",
			Visibility: "private",
		}

		snippet, err := snippetService.Create(user.ID, req)
		require.NoError(t, err)

		assert.Equal(t, "custom-slug", snippet.Slug)
		assert.Equal(t, "private", snippet.Visibility)
	})

	t.Run("create snippet with metadata", func(t *testing.T) {
		req := CreateRequest{
			Name:    "Template with Metadata",
			Content: "{{VAR1}} and {{VAR2}}",
			Metadata: map[string]string{
				"VAR1": "First variable description",
				"VAR2": "Second variable description",
			},
		}

		snippet, err := snippetService.Create(user.ID, req)
		require.NoError(t, err)

		assert.Equal(t, 2, len(snippet.Variables))
		assert.Equal(t, "First variable description", snippet.Variables[0].Description)
		assert.Equal(t, "Second variable description", snippet.Variables[1].Description)
	})

	t.Run("validation errors", func(t *testing.T) {
		// Missing name
		_, err := snippetService.Create(user.ID, CreateRequest{Content: "test"})
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "VALIDATION_ERROR", apiErr.Code)

		// Missing content
		_, err = snippetService.Create(user.ID, CreateRequest{Name: "test"})
		require.Error(t, err)
		apiErr, ok = err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "VALIDATION_ERROR", apiErr.Code)

		// Invalid visibility
		_, err = snippetService.Create(user.ID, CreateRequest{
			Name:       "test",
			Content:    "test",
			Visibility: "invalid",
		})
		require.Error(t, err)
		apiErr, ok = err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "VALIDATION_ERROR", apiErr.Code)
	})

	t.Run("duplicate slug", func(t *testing.T) {
		req := CreateRequest{
			Name:    "Duplicate Test",
			Slug:    "duplicate-slug",
			Content: "test content",
		}

		// First creation should succeed
		_, err := snippetService.Create(user.ID, req)
		require.NoError(t, err)

		// Second creation with same slug should fail
		_, err = snippetService.Create(user.ID, req)
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "DUPLICATE_SLUG", apiErr.Code)
	})
}

func TestSnippetService_Pull(t *testing.T) {
	userRepo, cleanup := setupTestDB(t)
	if cleanup == nil {
		return
	}
	defer cleanup()

	db, _ := repository.NewDB("postgres://gdsnip:gdsnip@localhost:5432/gdsnip?sslmode=disable")
	snippetRepo := repository.NewSnippetRepo(db)
	defer db.Close()

	snippetService := NewSnippetService(snippetRepo, userRepo)
	authService := NewAuthService(userRepo, "test-secret-key-that-is-very-long", 24*3600*1000000000)

	// Create test user and snippet
	user, _, err := authService.Register("pull@example.com", "pulluser", "password123")
	require.NoError(t, err)

	req := CreateRequest{
		Name:    "Pull Test Template",
		Content: "Hello {{NAME}}, your port is {{PORT}}",
	}
	snippet, err := snippetService.Create(user.ID, req)
	require.NoError(t, err)

	t.Run("successful pull with all variables", func(t *testing.T) {
		result, err := snippetService.Pull(snippet.Namespace, snippet.Slug, user.ID, map[string]string{
			"NAME": "World",
			"PORT": "8080",
		})
		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, "Hello World, your port is 8080", result.Content)
		assert.Equal(t, snippet.ID, result.Snippet.ID)
		assert.Equal(t, 2, len(result.Variables))
	})

	t.Run("missing required variables", func(t *testing.T) {
		_, err := snippetService.Pull(snippet.Namespace, snippet.Slug, user.ID, map[string]string{
			"NAME": "World",
		})
		require.Error(t, err)

		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "MISSING_VARIABLES", apiErr.Code)

		// Check that missing variables are included in details
		details, ok := apiErr.Details.(map[string]interface{})
		require.True(t, ok)
		missing, ok := details["missing"]
		require.True(t, ok)
		assert.NotNil(t, missing)
	})

	t.Run("snippet not found", func(t *testing.T) {
		_, err := snippetService.Pull("nonexistent", "nonexistent", user.ID, map[string]string{})
		require.Error(t, err)

		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "SNIPPET_NOT_FOUND", apiErr.Code)
	})

	t.Run("private snippet visibility", func(t *testing.T) {
		// Create another user
		user2, _, err := authService.Register("pull2@example.com", "pulluser2", "password123")
		require.NoError(t, err)

		// Create private snippet
		privateReq := CreateRequest{
			Name:       "Private Template",
			Content:    "Private content {{VAR}}",
			Visibility: "private",
		}
		privateSnippet, err := snippetService.Create(user.ID, privateReq)
		require.NoError(t, err)

		// Owner can access
		_, err = snippetService.Pull(privateSnippet.Namespace, privateSnippet.Slug, user.ID, map[string]string{"VAR": "value"})
		require.NoError(t, err)

		// Other user cannot access
		_, err = snippetService.Pull(privateSnippet.Namespace, privateSnippet.Slug, user2.ID, map[string]string{"VAR": "value"})
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "NOT_FOUND", apiErr.Code)
	})
}

func TestSnippetService_Search(t *testing.T) {
	userRepo, cleanup := setupTestDB(t)
	if cleanup == nil {
		return
	}
	defer cleanup()

	db, _ := repository.NewDB("postgres://gdsnip:gdsnip@localhost:5432/gdsnip?sslmode=disable")
	snippetRepo := repository.NewSnippetRepo(db)
	defer db.Close()

	snippetService := NewSnippetService(snippetRepo, userRepo)
	authService := NewAuthService(userRepo, "test-secret-key-that-is-very-long", 24*3600*1000000000)

	// Create test user and snippets
	user, _, err := authService.Register("search@example.com", "searchuser", "password123")
	require.NoError(t, err)

	// Create public snippets with different tags
	snippets := []CreateRequest{
		{
			Name:       "Docker PostgreSQL",
			Content:    "docker postgres {{DB_PASSWORD}}",
			Visibility: "public",
			Tags:       []string{"docker", "postgres"},
		},
		{
			Name:       "Docker Redis",
			Content:    "docker redis {{REDIS_PASSWORD}}",
			Visibility: "public",
			Tags:       []string{"docker", "redis"},
		},
		{
			Name:       "Kubernetes Deployment",
			Content:    "kubernetes deployment {{IMAGE}}",
			Visibility: "public",
			Tags:       []string{"kubernetes"},
		},
		{
			Name:       "Private Template",
			Content:    "private content",
			Visibility: "private",
			Tags:       []string{"private"},
		},
	}

	for _, req := range snippets {
		_, err := snippetService.Create(user.ID, req)
		require.NoError(t, err)
	}

	t.Run("search by query", func(t *testing.T) {
		results, total, err := snippetService.Search("docker", nil, 10, 0)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, total, 2)
		assert.GreaterOrEqual(t, len(results), 2)
	})

	t.Run("search by tags", func(t *testing.T) {
		results, total, err := snippetService.Search("", []string{"postgres"}, 10, 0)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, total, 1)
		assert.GreaterOrEqual(t, len(results), 1)
	})

	t.Run("search with pagination", func(t *testing.T) {
		// Search all
		_, total, err := snippetService.Search("", nil, 100, 0)
		require.NoError(t, err)

		// Should only return public snippets (3 out of 4)
		assert.GreaterOrEqual(t, total, 3)

		// Test pagination
		page1, _, err := snippetService.Search("", nil, 2, 0)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(page1), 2)

		page2, _, err := snippetService.Search("", nil, 2, 2)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(page2), 2)
	})

	t.Run("limit enforcement", func(t *testing.T) {
		// Test default limit
		results, _, err := snippetService.Search("", nil, 0, 0)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(results), 20)

		// Test max limit
		results, _, err = snippetService.Search("", nil, 200, 0)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(results), 100)
	})
}

func TestSnippetService_ListMine(t *testing.T) {
	userRepo, cleanup := setupTestDB(t)
	if cleanup == nil {
		return
	}
	defer cleanup()

	db, _ := repository.NewDB("postgres://gdsnip:gdsnip@localhost:5432/gdsnip?sslmode=disable")
	snippetRepo := repository.NewSnippetRepo(db)
	defer db.Close()

	snippetService := NewSnippetService(snippetRepo, userRepo)
	authService := NewAuthService(userRepo, "test-secret-key-that-is-very-long", 24*3600*1000000000)

	// Create test users
	user1, _, err := authService.Register("list1@example.com", "listuser1", "password123")
	require.NoError(t, err)

	user2, _, err := authService.Register("list2@example.com", "listuser2", "password123")
	require.NoError(t, err)

	// Create snippets for both users
	for i := 0; i < 3; i++ {
		_, err := snippetService.Create(user1.ID, CreateRequest{
			Name:    fmt.Sprintf("User1 Snippet %d", i),
			Content: "content {{VAR}}",
		})
		require.NoError(t, err)
	}

	for i := 0; i < 2; i++ {
		_, err := snippetService.Create(user2.ID, CreateRequest{
			Name:    fmt.Sprintf("User2 Snippet %d", i),
			Content: "content {{VAR}}",
		})
		require.NoError(t, err)
	}

	t.Run("list my snippets", func(t *testing.T) {
		results, total, err := snippetService.ListMine(user1.ID, 10, 0)
		require.NoError(t, err)

		assert.Equal(t, 3, total)
		assert.Equal(t, 3, len(results))

		// All snippets should belong to user1
		for _, s := range results {
			assert.Equal(t, user1.ID, s.OwnerID)
		}
	})

	t.Run("different user sees different snippets", func(t *testing.T) {
		results, total, err := snippetService.ListMine(user2.ID, 10, 0)
		require.NoError(t, err)

		assert.Equal(t, 2, total)
		assert.Equal(t, 2, len(results))

		for _, s := range results {
			assert.Equal(t, user2.ID, s.OwnerID)
		}
	})

	t.Run("pagination", func(t *testing.T) {
		page1, _, err := snippetService.ListMine(user1.ID, 2, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, len(page1))

		page2, _, err := snippetService.ListMine(user1.ID, 2, 2)
		require.NoError(t, err)
		assert.Equal(t, 1, len(page2))
	})
}

func TestSnippetService_GetByPath(t *testing.T) {
	userRepo, cleanup := setupTestDB(t)
	if cleanup == nil {
		return
	}
	defer cleanup()

	db, _ := repository.NewDB("postgres://gdsnip:gdsnip@localhost:5432/gdsnip?sslmode=disable")
	snippetRepo := repository.NewSnippetRepo(db)
	defer db.Close()

	snippetService := NewSnippetService(snippetRepo, userRepo)
	authService := NewAuthService(userRepo, "test-secret-key-that-is-very-long", 24*3600*1000000000)

	user, _, err := authService.Register("getpath@example.com", "getpathuser", "password123")
	require.NoError(t, err)

	req := CreateRequest{
		Name:    "Get Path Test",
		Content: "content {{VAR}}",
	}
	created, err := snippetService.Create(user.ID, req)
	require.NoError(t, err)

	t.Run("get existing snippet", func(t *testing.T) {
		snippet, err := snippetService.GetByPath(created.Namespace, created.Slug)
		require.NoError(t, err)
		require.NotNil(t, snippet)

		assert.Equal(t, created.ID, snippet.ID)
		assert.Equal(t, created.Name, snippet.Name)
		assert.Equal(t, created.Slug, snippet.Slug)
	})

	t.Run("get non-existent snippet", func(t *testing.T) {
		_, err := snippetService.GetByPath("nonexistent", "nonexistent")
		require.Error(t, err)

		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "SNIPPET_NOT_FOUND", apiErr.Code)
	})
}
