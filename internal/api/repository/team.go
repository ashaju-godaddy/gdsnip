package repository

import (
	"database/sql"
	"fmt"

	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// TeamRepo handles team database operations
type TeamRepo struct {
	db *sqlx.DB
}

// NewTeamRepo creates a new team repository
func NewTeamRepo(db *sqlx.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

// Create inserts a new team into the database
func (r *TeamRepo) Create(team *models.Team) error {
	query := `
		INSERT INTO teams (id, name, slug, description, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`
	_, err := r.db.Exec(query, team.ID, team.Name, team.Slug, team.Description, team.CreatedBy)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return models.NewDuplicateTeamSlugError(team.Slug)
			}
		}
		return fmt.Errorf("failed to create team: %w", err)
	}
	return nil
}

// GetByID retrieves a team by ID
func (r *TeamRepo) GetByID(id string) (*models.Team, error) {
	var team models.Team
	query := `
		SELECT id, name, slug, description, created_by, created_at, updated_at
		FROM teams
		WHERE id = $1
	`
	err := r.db.Get(&team, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewTeamNotFoundError(id)
		}
		return nil, fmt.Errorf("failed to get team by ID: %w", err)
	}
	return &team, nil
}

// GetBySlug retrieves a team by slug
func (r *TeamRepo) GetBySlug(slug string) (*models.Team, error) {
	var team models.Team
	query := `
		SELECT id, name, slug, description, created_by, created_at, updated_at
		FROM teams
		WHERE slug = $1
	`
	err := r.db.Get(&team, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewTeamNotFoundError(slug)
		}
		return nil, fmt.Errorf("failed to get team by slug: %w", err)
	}
	return &team, nil
}

// GetBySlugWithUserRole retrieves a team by slug with the user's role
func (r *TeamRepo) GetBySlugWithUserRole(slug, userID string) (*models.Team, error) {
	var team models.Team
	query := `
		SELECT
			t.id, t.name, t.slug, t.description, t.created_by, t.created_at, t.updated_at,
			COALESCE(tm.role, '') as role,
			(SELECT COUNT(*) FROM team_members WHERE team_id = t.id) as member_count
		FROM teams t
		LEFT JOIN team_members tm ON t.id = tm.team_id AND tm.user_id = $2
		WHERE t.slug = $1
	`
	err := r.db.Get(&team, query, slug, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewTeamNotFoundError(slug)
		}
		return nil, fmt.Errorf("failed to get team with role: %w", err)
	}
	return &team, nil
}

// Update updates a team's information
func (r *TeamRepo) Update(team *models.Team) error {
	query := `
		UPDATE teams
		SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.db.Exec(query, team.ID, team.Name, team.Description)
	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.NewTeamNotFoundError(team.ID)
	}
	return nil
}

// Delete removes a team from the database
func (r *TeamRepo) Delete(id string) error {
	query := `DELETE FROM teams WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.NewTeamNotFoundError(id)
	}
	return nil
}

// ListByUser retrieves all teams a user is a member of
func (r *TeamRepo) ListByUser(userID string, limit, offset int) ([]models.Team, int, error) {
	var teams []models.Team
	var totalCount int

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1
	`
	err := r.db.Get(&totalCount, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count teams: %w", err)
	}

	// Get teams with role
	query := `
		SELECT
			t.id, t.name, t.slug, t.description, t.created_by, t.created_at, t.updated_at,
			tm.role,
			(SELECT COUNT(*) FROM team_members WHERE team_id = t.id) as member_count
		FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1
		ORDER BY t.name ASC
		LIMIT $2 OFFSET $3
	`
	err = r.db.Select(&teams, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list teams: %w", err)
	}

	return teams, totalCount, nil
}

// AddMember adds a user to a team with a role
func (r *TeamRepo) AddMember(member *models.TeamMember) error {
	query := `
		INSERT INTO team_members (id, team_id, user_id, role, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := r.db.Exec(query, member.ID, member.TeamID, member.UserID, member.Role)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return models.NewAlreadyMemberError(member.UserID)
			}
		}
		return fmt.Errorf("failed to add member: %w", err)
	}
	return nil
}

// RemoveMember removes a user from a team
func (r *TeamRepo) RemoveMember(teamID, userID string) error {
	query := `DELETE FROM team_members WHERE team_id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, teamID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.NewMemberNotFoundError(userID)
	}
	return nil
}

// UpdateMemberRole updates a member's role within a team
func (r *TeamRepo) UpdateMemberRole(teamID, userID string, role models.TeamRole) error {
	query := `
		UPDATE team_members
		SET role = $3
		WHERE team_id = $1 AND user_id = $2
	`
	result, err := r.db.Exec(query, teamID, userID, role)
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return models.NewMemberNotFoundError(userID)
	}
	return nil
}

// GetMember retrieves a specific team membership
func (r *TeamRepo) GetMember(teamID, userID string) (*models.TeamMember, error) {
	var member models.TeamMember
	query := `
		SELECT tm.id, tm.team_id, tm.user_id, tm.role, tm.created_at,
		       u.username, u.email
		FROM team_members tm
		INNER JOIN users u ON tm.user_id = u.id
		WHERE tm.team_id = $1 AND tm.user_id = $2
	`
	err := r.db.Get(&member, query, teamID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewMemberNotFoundError(userID)
		}
		return nil, fmt.Errorf("failed to get member: %w", err)
	}
	return &member, nil
}

// GetMemberByUsername retrieves a team membership by username
func (r *TeamRepo) GetMemberByUsername(teamID, username string) (*models.TeamMember, error) {
	var member models.TeamMember
	query := `
		SELECT tm.id, tm.team_id, tm.user_id, tm.role, tm.created_at,
		       u.username, u.email
		FROM team_members tm
		INNER JOIN users u ON tm.user_id = u.id
		WHERE tm.team_id = $1 AND u.username = $2
	`
	err := r.db.Get(&member, query, teamID, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewMemberNotFoundError(username)
		}
		return nil, fmt.Errorf("failed to get member by username: %w", err)
	}
	return &member, nil
}

// ListMembers retrieves all members of a team
func (r *TeamRepo) ListMembers(teamID string) ([]models.TeamMember, error) {
	var members []models.TeamMember
	query := `
		SELECT tm.id, tm.team_id, tm.user_id, tm.role, tm.created_at,
		       u.username, u.email
		FROM team_members tm
		INNER JOIN users u ON tm.user_id = u.id
		WHERE tm.team_id = $1
		ORDER BY
			CASE tm.role
				WHEN 'owner' THEN 1
				WHEN 'admin' THEN 2
				WHEN 'member' THEN 3
				WHEN 'viewer' THEN 4
			END,
			u.username ASC
	`
	err := r.db.Select(&members, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list members: %w", err)
	}
	return members, nil
}

// GetUserRole retrieves a user's role in a team
func (r *TeamRepo) GetUserRole(teamID, userID string) (models.TeamRole, error) {
	var role models.TeamRole
	query := `SELECT role FROM team_members WHERE team_id = $1 AND user_id = $2`
	err := r.db.Get(&role, query, teamID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", models.NewMemberNotFoundError(userID)
		}
		return "", fmt.Errorf("failed to get user role: %w", err)
	}
	return role, nil
}

// IsMember checks if a user is a member of a team
func (r *TeamRepo) IsMember(teamID, userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM team_members WHERE team_id = $1 AND user_id = $2`
	err := r.db.Get(&count, query, teamID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check membership: %w", err)
	}
	return count > 0, nil
}
