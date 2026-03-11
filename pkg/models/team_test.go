package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== TeamRole.HasPermission Tests ====================

func TestTeamRole_HasPermission_Owner(t *testing.T) {
	tests := []struct {
		permission Permission
		expected   bool
	}{
		{PermViewTeam, true},
		{PermManageTeam, true},
		{PermDeleteTeam, true},
		{PermManageMembers, true},
		{PermCreateSnippet, true},
		{PermEditAnySnippet, true},
		{PermEditOwnSnippet, false}, // Owner has EditAny, not EditOwn
		{PermDeleteSnippet, true},
		{PermPullSnippet, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.permission), func(t *testing.T) {
			result := RoleOwner.HasPermission(tt.permission)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTeamRole_HasPermission_Admin(t *testing.T) {
	tests := []struct {
		permission Permission
		expected   bool
	}{
		{PermViewTeam, true},
		{PermManageTeam, true},
		{PermDeleteTeam, false}, // Admin cannot delete team
		{PermManageMembers, true},
		{PermCreateSnippet, true},
		{PermEditAnySnippet, true},
		{PermDeleteSnippet, true},
		{PermPullSnippet, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.permission), func(t *testing.T) {
			result := RoleAdmin.HasPermission(tt.permission)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTeamRole_HasPermission_Member(t *testing.T) {
	tests := []struct {
		permission Permission
		expected   bool
	}{
		{PermViewTeam, true},
		{PermManageTeam, false},
		{PermDeleteTeam, false},
		{PermManageMembers, false},
		{PermCreateSnippet, true},
		{PermEditAnySnippet, false},
		{PermEditOwnSnippet, true},
		{PermDeleteSnippet, false},
		{PermPullSnippet, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.permission), func(t *testing.T) {
			result := RoleMember.HasPermission(tt.permission)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTeamRole_HasPermission_Viewer(t *testing.T) {
	tests := []struct {
		permission Permission
		expected   bool
	}{
		{PermViewTeam, true},
		{PermManageTeam, false},
		{PermDeleteTeam, false},
		{PermManageMembers, false},
		{PermCreateSnippet, false},
		{PermEditAnySnippet, false},
		{PermEditOwnSnippet, false},
		{PermDeleteSnippet, false},
		{PermPullSnippet, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.permission), func(t *testing.T) {
			result := RoleViewer.HasPermission(tt.permission)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTeamRole_HasPermission_InvalidRole(t *testing.T) {
	invalidRole := TeamRole("invalid")
	assert.False(t, invalidRole.HasPermission(PermViewTeam))
	assert.False(t, invalidRole.HasPermission(PermCreateSnippet))
}

// ==================== TeamRole.IsValid Tests ====================

func TestTeamRole_IsValid(t *testing.T) {
	tests := []struct {
		role     TeamRole
		expected bool
	}{
		{RoleOwner, true},
		{RoleAdmin, true},
		{RoleMember, true},
		{RoleViewer, true},
		{TeamRole("invalid"), false},
		{TeamRole(""), false},
		{TeamRole("OWNER"), false}, // Case sensitive
		{TeamRole("Owner"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			result := tt.role.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ==================== TeamRole.IsAssignable Tests ====================

func TestTeamRole_IsAssignable(t *testing.T) {
	tests := []struct {
		role     TeamRole
		expected bool
	}{
		{RoleOwner, false},  // Owner cannot be assigned via API
		{RoleAdmin, true},
		{RoleMember, true},
		{RoleViewer, true},
		{TeamRole("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			result := tt.role.IsAssignable()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ==================== Permission Hierarchy Tests ====================

func TestPermissionHierarchy(t *testing.T) {
	// Test that roles have the correct permission hierarchy
	// Owner > Admin > Member > Viewer

	// Owner should have all admin permissions
	for _, perm := range RolePermissions[RoleAdmin] {
		// Skip PermDeleteTeam which admin doesn't have
		assert.True(t, RoleOwner.HasPermission(perm),
			"Owner should have admin permission: %s", perm)
	}

	// Admin should have all member permissions except EditOwnSnippet
	memberPerms := []Permission{PermViewTeam, PermCreateSnippet, PermPullSnippet}
	for _, perm := range memberPerms {
		assert.True(t, RoleAdmin.HasPermission(perm),
			"Admin should have member permission: %s", perm)
	}

	// Member should have all viewer permissions
	for _, perm := range RolePermissions[RoleViewer] {
		assert.True(t, RoleMember.HasPermission(perm),
			"Member should have viewer permission: %s", perm)
	}
}

// ==================== Role Constants Tests ====================

func TestRoleConstants(t *testing.T) {
	assert.Equal(t, TeamRole("owner"), RoleOwner)
	assert.Equal(t, TeamRole("admin"), RoleAdmin)
	assert.Equal(t, TeamRole("member"), RoleMember)
	assert.Equal(t, TeamRole("viewer"), RoleViewer)
}

// ==================== Permission Constants Tests ====================

func TestPermissionConstants(t *testing.T) {
	assert.Equal(t, Permission("team:view"), PermViewTeam)
	assert.Equal(t, Permission("team:manage"), PermManageTeam)
	assert.Equal(t, Permission("team:delete"), PermDeleteTeam)
	assert.Equal(t, Permission("team:members:manage"), PermManageMembers)
	assert.Equal(t, Permission("snippet:create"), PermCreateSnippet)
	assert.Equal(t, Permission("snippet:edit:any"), PermEditAnySnippet)
	assert.Equal(t, Permission("snippet:edit:own"), PermEditOwnSnippet)
	assert.Equal(t, Permission("snippet:delete"), PermDeleteSnippet)
	assert.Equal(t, Permission("snippet:pull"), PermPullSnippet)
}

// ==================== RolePermissions Map Tests ====================

func TestRolePermissions_AllRolesHavePermissions(t *testing.T) {
	roles := []TeamRole{RoleOwner, RoleAdmin, RoleMember, RoleViewer}

	for _, role := range roles {
		perms, exists := RolePermissions[role]
		assert.True(t, exists, "Role %s should exist in RolePermissions", role)
		assert.NotEmpty(t, perms, "Role %s should have at least one permission", role)
	}
}

func TestRolePermissions_OwnerHasMostPermissions(t *testing.T) {
	ownerCount := len(RolePermissions[RoleOwner])
	adminCount := len(RolePermissions[RoleAdmin])
	memberCount := len(RolePermissions[RoleMember])
	viewerCount := len(RolePermissions[RoleViewer])

	assert.GreaterOrEqual(t, ownerCount, adminCount)
	assert.GreaterOrEqual(t, adminCount, memberCount)
	assert.GreaterOrEqual(t, memberCount, viewerCount)
}
