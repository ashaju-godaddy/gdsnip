package service

import (
	"fmt"

	"github.com/ashaju-godaddy/gdsnip/internal/api/repository"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/ashaju-godaddy/gdsnip/pkg/validator"
	"github.com/google/uuid"
)

// TeamService handles team business logic
type TeamService struct {
	teamRepo *repository.TeamRepo
	userRepo *repository.UserRepo
}

// NewTeamService creates a new team service
func NewTeamService(teamRepo *repository.TeamRepo, userRepo *repository.UserRepo) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// CreateTeam creates a new team with the creator as owner
func (s *TeamService) CreateTeam(userID string, req models.CreateTeamRequest) (*models.Team, error) {
	// Validate name
	if req.Name == "" {
		return nil, models.NewValidationError("name is required", nil)
	}
	if len(req.Name) > 100 {
		return nil, models.NewValidationError("name must be at most 100 characters", nil)
	}

	// Auto-generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = validator.GenerateSlug(req.Name)
	} else {
		if err := validator.ValidateSlug(slug); err != nil {
			return nil, models.NewValidationError(err.Error(), nil)
		}
	}

	// Validate description length
	if len(req.Description) > 500 {
		return nil, models.NewValidationError("description must be at most 500 characters", nil)
	}

	// Create team
	team := &models.Team{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		CreatedBy:   userID,
	}

	if err := s.teamRepo.Create(team); err != nil {
		return nil, err
	}

	// Add creator as owner
	member := &models.TeamMember{
		ID:     uuid.New().String(),
		TeamID: team.ID,
		UserID: userID,
		Role:   models.RoleOwner,
	}
	if err := s.teamRepo.AddMember(member); err != nil {
		// Rollback team creation on failure
		_ = s.teamRepo.Delete(team.ID)
		return nil, models.NewInternalError("failed to add owner to team")
	}

	// Set member count and role for response
	team.MemberCount = 1
	team.Role = models.RoleOwner

	return team, nil
}

// GetTeam retrieves a team by slug for a user
func (s *TeamService) GetTeam(userID, teamSlug string) (*models.Team, error) {
	team, err := s.teamRepo.GetBySlugWithUserRole(teamSlug, userID)
	if err != nil {
		return nil, err
	}

	// Check if user is a member
	if team.Role == "" {
		return nil, models.NewForbiddenError("you are not a member of this team")
	}

	return team, nil
}

// UpdateTeam updates a team's information
func (s *TeamService) UpdateTeam(userID, teamSlug string, req models.UpdateTeamRequest) (*models.Team, error) {
	// Get team and check permissions
	team, err := s.teamRepo.GetBySlug(teamSlug)
	if err != nil {
		return nil, err
	}

	role, err := s.teamRepo.GetUserRole(team.ID, userID)
	if err != nil {
		return nil, models.NewForbiddenError("you are not a member of this team")
	}

	if !role.HasPermission(models.PermManageTeam) {
		return nil, models.NewForbiddenError("you do not have permission to manage this team")
	}

	// Apply updates
	if req.Name != nil {
		if *req.Name == "" {
			return nil, models.NewValidationError("name cannot be empty", nil)
		}
		if len(*req.Name) > 100 {
			return nil, models.NewValidationError("name must be at most 100 characters", nil)
		}
		team.Name = *req.Name
	}

	if req.Description != nil {
		if len(*req.Description) > 500 {
			return nil, models.NewValidationError("description must be at most 500 characters", nil)
		}
		team.Description = *req.Description
	}

	if err := s.teamRepo.Update(team); err != nil {
		return nil, err
	}

	// Fetch updated team with role
	return s.teamRepo.GetBySlugWithUserRole(teamSlug, userID)
}

// DeleteTeam deletes a team (owner only)
func (s *TeamService) DeleteTeam(userID, teamSlug string) error {
	team, err := s.teamRepo.GetBySlug(teamSlug)
	if err != nil {
		return err
	}

	role, err := s.teamRepo.GetUserRole(team.ID, userID)
	if err != nil {
		return models.NewForbiddenError("you are not a member of this team")
	}

	if !role.HasPermission(models.PermDeleteTeam) {
		return models.NewForbiddenError("only the team owner can delete the team")
	}

	return s.teamRepo.Delete(team.ID)
}

// ListUserTeams lists all teams the user is a member of
func (s *TeamService) ListUserTeams(userID string, limit, offset int) ([]models.Team, int, error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.teamRepo.ListByUser(userID, limit, offset)
}

// AddMember adds a user to a team
func (s *TeamService) AddMember(userID, teamSlug string, req models.AddMemberRequest) (*models.TeamMember, error) {
	// Validate request
	if req.Username == "" {
		return nil, models.NewValidationError("username is required", nil)
	}

	// Validate role (cannot add as owner via API)
	if !req.Role.IsValid() {
		return nil, models.NewInvalidRoleError(string(req.Role))
	}
	if req.Role == models.RoleOwner {
		return nil, models.NewValidationError("cannot add a member as owner; ownership must be transferred", nil)
	}

	// Get team
	team, err := s.teamRepo.GetBySlug(teamSlug)
	if err != nil {
		return nil, err
	}

	// Check permissions
	callerRole, err := s.teamRepo.GetUserRole(team.ID, userID)
	if err != nil {
		return nil, models.NewForbiddenError("you are not a member of this team")
	}

	if !callerRole.HasPermission(models.PermManageMembers) {
		return nil, models.NewForbiddenError("you do not have permission to manage members")
	}

	// Get user to add
	targetUser, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	// Create membership
	member := &models.TeamMember{
		ID:       uuid.New().String(),
		TeamID:   team.ID,
		UserID:   targetUser.ID,
		Role:     req.Role,
		Username: targetUser.Username,
		Email:    targetUser.Email,
	}

	if err := s.teamRepo.AddMember(member); err != nil {
		return nil, err
	}

	return member, nil
}

// RemoveMember removes a user from a team
func (s *TeamService) RemoveMember(userID, teamSlug, targetUsername string) error {
	team, err := s.teamRepo.GetBySlug(teamSlug)
	if err != nil {
		return err
	}

	// Check permissions
	callerRole, err := s.teamRepo.GetUserRole(team.ID, userID)
	if err != nil {
		return models.NewForbiddenError("you are not a member of this team")
	}

	if !callerRole.HasPermission(models.PermManageMembers) {
		return models.NewForbiddenError("you do not have permission to manage members")
	}

	// Get target member
	targetMember, err := s.teamRepo.GetMemberByUsername(team.ID, targetUsername)
	if err != nil {
		return err
	}

	// Cannot remove owner
	if targetMember.Role == models.RoleOwner {
		return models.NewCannotRemoveOwnerError()
	}

	return s.teamRepo.RemoveMember(team.ID, targetMember.UserID)
}

// UpdateMemberRole updates a member's role
func (s *TeamService) UpdateMemberRole(userID, teamSlug, targetUsername string, newRole models.TeamRole) (*models.TeamMember, error) {
	// Validate role
	if !newRole.IsValid() {
		return nil, models.NewInvalidRoleError(string(newRole))
	}
	if newRole == models.RoleOwner {
		return nil, models.NewValidationError("cannot change role to owner; ownership must be transferred separately", nil)
	}

	team, err := s.teamRepo.GetBySlug(teamSlug)
	if err != nil {
		return nil, err
	}

	// Check permissions
	callerRole, err := s.teamRepo.GetUserRole(team.ID, userID)
	if err != nil {
		return nil, models.NewForbiddenError("you are not a member of this team")
	}

	if !callerRole.HasPermission(models.PermManageMembers) {
		return nil, models.NewForbiddenError("you do not have permission to manage members")
	}

	// Get target member
	targetMember, err := s.teamRepo.GetMemberByUsername(team.ID, targetUsername)
	if err != nil {
		return nil, err
	}

	// Cannot change owner's role
	if targetMember.Role == models.RoleOwner {
		return nil, models.NewValidationError("cannot change the team owner's role", nil)
	}

	// Update role
	if err := s.teamRepo.UpdateMemberRole(team.ID, targetMember.UserID, newRole); err != nil {
		return nil, err
	}

	// Return updated member
	return s.teamRepo.GetMemberByUsername(team.ID, targetUsername)
}

// LeaveTeam allows a user to leave a team
func (s *TeamService) LeaveTeam(userID, teamSlug string) error {
	team, err := s.teamRepo.GetBySlug(teamSlug)
	if err != nil {
		return err
	}

	// Get user's role
	role, err := s.teamRepo.GetUserRole(team.ID, userID)
	if err != nil {
		return err // Not a member
	}

	// Owner cannot leave
	if role == models.RoleOwner {
		return models.NewOwnerCannotLeaveError()
	}

	return s.teamRepo.RemoveMember(team.ID, userID)
}

// ListMembers lists all members of a team
func (s *TeamService) ListMembers(userID, teamSlug string) ([]models.TeamMember, error) {
	team, err := s.teamRepo.GetBySlug(teamSlug)
	if err != nil {
		return nil, err
	}

	// Verify user is a member
	_, err = s.teamRepo.GetUserRole(team.ID, userID)
	if err != nil {
		return nil, models.NewForbiddenError("you are not a member of this team")
	}

	return s.teamRepo.ListMembers(team.ID)
}

// CheckPermission checks if a user has a specific permission in a team
func (s *TeamService) CheckPermission(userID, teamID string, permission models.Permission) error {
	role, err := s.teamRepo.GetUserRole(teamID, userID)
	if err != nil {
		return models.NewForbiddenError("you are not a member of this team")
	}

	if !role.HasPermission(permission) {
		return models.NewForbiddenError(fmt.Sprintf("you do not have permission: %s", permission))
	}

	return nil
}

// GetTeamBySlug retrieves a team by slug (no membership check - for namespace resolution)
func (s *TeamService) GetTeamBySlug(slug string) (*models.Team, error) {
	return s.teamRepo.GetBySlug(slug)
}

// IsMember checks if a user is a member of a team
func (s *TeamService) IsMember(teamID, userID string) (bool, error) {
	return s.teamRepo.IsMember(teamID, userID)
}
