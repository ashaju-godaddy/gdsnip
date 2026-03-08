package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// SnippetRepo handles snippet database operations
type SnippetRepo struct {
	db *sqlx.DB
}

// NewSnippetRepo creates a new snippet repository
func NewSnippetRepo(db *sqlx.DB) *SnippetRepo {
	return &SnippetRepo{db: db}
}

// Create inserts a new snippet into the database
func (r *SnippetRepo) Create(snippet *models.Snippet) error {
	// Convert variables to JSON
	variablesJSON, err := json.Marshal(snippet.Variables)
	if err != nil {
		return fmt.Errorf("failed to marshal variables: %w", err)
	}

	query := `
		INSERT INTO snippets (
			id, name, slug, description, content, variables, tags,
			visibility, owner_type, owner_id, namespace, created_by,
			version, pull_count, created_at, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12,
			$13, $14, NOW(), NOW()
		)
	`

	_, err = r.db.Exec(query,
		snippet.ID, snippet.Name, snippet.Slug, snippet.Description,
		snippet.Content, variablesJSON, pq.Array(snippet.Tags),
		snippet.Visibility, snippet.OwnerType, snippet.OwnerID,
		snippet.Namespace, snippet.CreatedBy,
		snippet.Version, snippet.PullCount,
	)

	if err != nil {
		// Check for unique constraint violations
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return models.NewDuplicateSlugError(snippet.Namespace, snippet.Slug)
			}
		}
		return fmt.Errorf("failed to create snippet: %w", err)
	}

	return nil
}

// GetByNamespaceSlug retrieves a snippet by namespace and slug
func (r *SnippetRepo) GetByNamespaceSlug(namespace, slug string) (*models.Snippet, error) {
	var snippet models.Snippet
	var variablesJSON []byte

	query := `
		SELECT
			id, name, slug, description, content, variables, tags,
			visibility, owner_type, owner_id, namespace, created_by,
			version, pull_count, created_at, updated_at
		FROM snippets
		WHERE namespace = $1 AND slug = $2
	`

	err := r.db.QueryRow(query, namespace, slug).Scan(
		&snippet.ID, &snippet.Name, &snippet.Slug, &snippet.Description,
		&snippet.Content, &variablesJSON, pq.Array(&snippet.Tags),
		&snippet.Visibility, &snippet.OwnerType, &snippet.OwnerID,
		&snippet.Namespace, &snippet.CreatedBy,
		&snippet.Version, &snippet.PullCount, &snippet.CreatedAt, &snippet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewSnippetNotFoundError(namespace + "/" + slug)
		}
		return nil, fmt.Errorf("failed to get snippet: %w", err)
	}

	// Unmarshal variables JSON
	if err := json.Unmarshal(variablesJSON, &snippet.Variables); err != nil {
		return nil, fmt.Errorf("failed to unmarshal variables: %w", err)
	}

	return &snippet, nil
}

// Search searches for snippets using full-text search and optional tag filtering
func (r *SnippetRepo) Search(query string, tags []string, limit, offset int) ([]models.Snippet, int, error) {
	var snippets []models.Snippet
	var totalCount int

	// Build WHERE clause
	whereClauses := []string{"visibility = 'public'"}
	args := []interface{}{}
	argIdx := 1

	// Add full-text search if query provided
	if query != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("search_vector @@ plainto_tsquery('english', $%d)", argIdx))
		args = append(args, query)
		argIdx++
	}

	// Add tag filtering if tags provided
	if len(tags) > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("tags && $%d", argIdx))
		args = append(args, pq.Array(tags))
		argIdx++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM snippets WHERE %s", whereClause)
	err := r.db.Get(&totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count snippets: %w", err)
	}

	// Get snippets with pagination
	selectQuery := fmt.Sprintf(`
		SELECT
			id, name, slug, description, content, variables, tags,
			visibility, owner_type, owner_id, namespace, created_by,
			version, pull_count, created_at, updated_at
		FROM snippets
		WHERE %s
		ORDER BY pull_count DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search snippets: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var snippet models.Snippet
		var variablesJSON []byte

		err := rows.Scan(
			&snippet.ID, &snippet.Name, &snippet.Slug, &snippet.Description,
			&snippet.Content, &variablesJSON, pq.Array(&snippet.Tags),
			&snippet.Visibility, &snippet.OwnerType, &snippet.OwnerID,
			&snippet.Namespace, &snippet.CreatedBy,
			&snippet.Version, &snippet.PullCount, &snippet.CreatedAt, &snippet.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan snippet: %w", err)
		}

		// Unmarshal variables JSON
		if err := json.Unmarshal(variablesJSON, &snippet.Variables); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal variables: %w", err)
		}

		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating snippets: %w", err)
	}

	return snippets, totalCount, nil
}

// ListByOwner retrieves all snippets owned by a user
func (r *SnippetRepo) ListByOwner(ownerType, ownerID string, limit, offset int) ([]models.Snippet, int, error) {
	var snippets []models.Snippet
	var totalCount int

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM snippets
		WHERE owner_type = $1 AND owner_id = $2
	`
	err := r.db.Get(&totalCount, countQuery, ownerType, ownerID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user snippets: %w", err)
	}

	// Get snippets with pagination
	selectQuery := `
		SELECT
			id, name, slug, description, content, variables, tags,
			visibility, owner_type, owner_id, namespace, created_by,
			version, pull_count, created_at, updated_at
		FROM snippets
		WHERE owner_type = $1 AND owner_id = $2
		ORDER BY updated_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(selectQuery, ownerType, ownerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list user snippets: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var snippet models.Snippet
		var variablesJSON []byte

		err := rows.Scan(
			&snippet.ID, &snippet.Name, &snippet.Slug, &snippet.Description,
			&snippet.Content, &variablesJSON, pq.Array(&snippet.Tags),
			&snippet.Visibility, &snippet.OwnerType, &snippet.OwnerID,
			&snippet.Namespace, &snippet.CreatedBy,
			&snippet.Version, &snippet.PullCount, &snippet.CreatedAt, &snippet.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan snippet: %w", err)
		}

		// Unmarshal variables JSON
		if err := json.Unmarshal(variablesJSON, &snippet.Variables); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal variables: %w", err)
		}

		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating snippets: %w", err)
	}

	return snippets, totalCount, nil
}

// IncrementPullCount atomically increments the pull count for a snippet
func (r *SnippetRepo) IncrementPullCount(snippetID string) error {
	query := `
		UPDATE snippets
		SET pull_count = pull_count + 1
		WHERE id = $1
	`

	_, err := r.db.Exec(query, snippetID)
	if err != nil {
		return fmt.Errorf("failed to increment pull count: %w", err)
	}

	return nil
}
