package models

import "time"

// TeamRole represents a user's role within a team
type TeamRole string

const (
	RoleOwner  TeamRole = "owner"  // Full control, can delete team
	RoleAdmin  TeamRole = "admin"  // Manage members, manage all snippets
	RoleMember TeamRole = "member" // Create/edit own snippets, pull team snippets
	RoleViewer TeamRole = "viewer" // Pull team snippets only
)

// Permission represents a specific action within a team
type Permission string

const (
	PermViewTeam       Permission = "team:view"
	PermManageTeam     Permission = "team:manage"
	PermDeleteTeam     Permission = "team:delete"
	PermManageMembers  Permission = "team:members:manage"
	PermCreateSnippet  Permission = "snippet:create"
	PermEditAnySnippet Permission = "snippet:edit:any"
	PermEditOwnSnippet Permission = "snippet:edit:own"
	PermDeleteSnippet  Permission = "snippet:delete"
	PermPullSnippet    Permission = "snippet:pull"
)

// RolePermissions maps roles to their permissions
var RolePermissions = map[TeamRole][]Permission{
	RoleOwner: {
		PermViewTeam, PermManageTeam, PermDeleteTeam,
		PermManageMembers, PermCreateSnippet, PermEditAnySnippet,
		PermDeleteSnippet, PermPullSnippet,
	},
	RoleAdmin: {
		PermViewTeam, PermManageTeam, PermManageMembers,
		PermCreateSnippet, PermEditAnySnippet, PermDeleteSnippet,
		PermPullSnippet,
	},
	RoleMember: {
		PermViewTeam, PermCreateSnippet, PermEditOwnSnippet,
		PermPullSnippet,
	},
	RoleViewer: {
		PermViewTeam, PermPullSnippet,
	},
}

// HasPermission checks if the role has the specified permission
func (r TeamRole) HasPermission(p Permission) bool {
	perms, ok := RolePermissions[r]
	if !ok {
		return false
	}
	for _, perm := range perms {
		if perm == p {
			return true
		}
	}
	return false
}

// IsValid checks if the role is a valid team role
func (r TeamRole) IsValid() bool {
	switch r {
	case RoleOwner, RoleAdmin, RoleMember, RoleViewer:
		return true
	default:
		return false
	}
}

// IsAssignable checks if the role can be assigned via API (owner cannot be assigned)
func (r TeamRole) IsAssignable() bool {
	switch r {
	case RoleAdmin, RoleMember, RoleViewer:
		return true
	default:
		return false
	}
}

// Team represents a team/organization
type Team struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description string    `json:"description" db:"description"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Populated on fetch via JOIN or separate query
	MemberCount int      `json:"member_count,omitempty" db:"member_count"`
	Role        TeamRole `json:"role,omitempty" db:"role"` // Current user's role
}

// TeamMember represents a user's membership in a team
type TeamMember struct {
	ID        string    `json:"id" db:"id"`
	TeamID    string    `json:"team_id" db:"team_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Role      TeamRole  `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Populated via JOIN
	Username string `json:"username,omitempty" db:"username"`
	Email    string `json:"email,omitempty" db:"email"`
}

// Request/Response types

// CreateTeamRequest represents the request body for creating a team
type CreateTeamRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateTeamRequest represents the request body for updating a team
type UpdateTeamRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// AddMemberRequest represents the request body for adding a team member
type AddMemberRequest struct {
	Username string   `json:"username"`
	Role     TeamRole `json:"role"`
}

// UpdateMemberRoleRequest represents the request body for updating a member's role
type UpdateMemberRoleRequest struct {
	Role TeamRole `json:"role"`
}

// TeamMemberListResponse represents a list of team members
type TeamMemberListResponse struct {
	Members []TeamMember `json:"members"`
	Total   int          `json:"total"`
}
