package handlers

import (
	"net/http"
	"strconv"

	"github.com/ashaju-godaddy/gdsnip/internal/api/service"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/labstack/echo/v4"
)

// TeamHandler handles team endpoints
type TeamHandler struct {
	teamService    *service.TeamService
	snippetService *service.SnippetService
}

// NewTeamHandler creates a new team handler
func NewTeamHandler(teamService *service.TeamService, snippetService *service.SnippetService) *TeamHandler {
	return &TeamHandler{
		teamService:    teamService,
		snippetService: snippetService,
	}
}

// Create handles team creation
func (h *TeamHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req models.CreateTeamRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, models.NewValidationError("invalid request body", nil))
	}

	team, err := h.teamService.CreateTeam(userID, req)
	if err != nil {
		return Error(c, err)
	}

	return Success(c, http.StatusCreated, team)
}

// List lists teams the user is a member of
func (h *TeamHandler) List(c echo.Context) error {
	userID := c.Get("user_id").(string)

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	teams, total, err := h.teamService.ListUserTeams(userID, limit, offset)
	if err != nil {
		return Error(c, err)
	}

	return PaginatedSuccess(c, http.StatusOK, teams, total, limit, offset)
}

// Get retrieves a team by slug
func (h *TeamHandler) Get(c echo.Context) error {
	userID := c.Get("user_id").(string)
	slug := c.Param("slug")

	team, err := h.teamService.GetTeam(userID, slug)
	if err != nil {
		return Error(c, err)
	}

	return Success(c, http.StatusOK, team)
}

// Update updates a team
func (h *TeamHandler) Update(c echo.Context) error {
	userID := c.Get("user_id").(string)
	slug := c.Param("slug")

	var req models.UpdateTeamRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, models.NewValidationError("invalid request body", nil))
	}

	team, err := h.teamService.UpdateTeam(userID, slug, req)
	if err != nil {
		return Error(c, err)
	}

	return Success(c, http.StatusOK, team)
}

// Delete deletes a team
func (h *TeamHandler) Delete(c echo.Context) error {
	userID := c.Get("user_id").(string)
	slug := c.Param("slug")

	if err := h.teamService.DeleteTeam(userID, slug); err != nil {
		return Error(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "team deleted successfully",
	})
}

// ListMembers lists team members
func (h *TeamHandler) ListMembers(c echo.Context) error {
	userID := c.Get("user_id").(string)
	slug := c.Param("slug")

	members, err := h.teamService.ListMembers(userID, slug)
	if err != nil {
		return Error(c, err)
	}

	return Success(c, http.StatusOK, models.TeamMemberListResponse{
		Members: members,
		Total:   len(members),
	})
}

// AddMember adds a member to a team
func (h *TeamHandler) AddMember(c echo.Context) error {
	userID := c.Get("user_id").(string)
	slug := c.Param("slug")

	var req models.AddMemberRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, models.NewValidationError("invalid request body", nil))
	}

	member, err := h.teamService.AddMember(userID, slug, req)
	if err != nil {
		return Error(c, err)
	}

	return Success(c, http.StatusCreated, member)
}

// UpdateMemberRole updates a member's role
func (h *TeamHandler) UpdateMemberRole(c echo.Context) error {
	userID := c.Get("user_id").(string)
	slug := c.Param("slug")
	username := c.Param("username")

	var req models.UpdateMemberRoleRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, models.NewValidationError("invalid request body", nil))
	}

	member, err := h.teamService.UpdateMemberRole(userID, slug, username, req.Role)
	if err != nil {
		return Error(c, err)
	}

	return Success(c, http.StatusOK, member)
}

// RemoveMember removes a member from a team
func (h *TeamHandler) RemoveMember(c echo.Context) error {
	userID := c.Get("user_id").(string)
	slug := c.Param("slug")
	username := c.Param("username")

	if err := h.teamService.RemoveMember(userID, slug, username); err != nil {
		return Error(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "member removed successfully",
	})
}

// Leave allows a user to leave a team
func (h *TeamHandler) Leave(c echo.Context) error {
	userID := c.Get("user_id").(string)
	slug := c.Param("slug")

	if err := h.teamService.LeaveTeam(userID, slug); err != nil {
		return Error(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "left team successfully",
	})
}

// ListSnippets lists team snippets
func (h *TeamHandler) ListSnippets(c echo.Context) error {
	userID := c.Get("user_id").(string)
	slug := c.Param("slug")

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	snippets, total, err := h.snippetService.ListTeamSnippets(userID, slug, limit, offset)
	if err != nil {
		return Error(c, err)
	}

	return PaginatedSuccess(c, http.StatusOK, snippets, total, limit, offset)
}
