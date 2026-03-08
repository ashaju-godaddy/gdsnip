package handlers

import (
	"net/http"
	"strconv"

	"github.com/ashaju-godaddy/gdsnip/internal/api/service"
	"github.com/labstack/echo/v4"
)

// SnippetHandler handles snippet endpoints
type SnippetHandler struct {
	snippetService *service.SnippetService
}

// NewSnippetHandler creates a new snippet handler
func NewSnippetHandler(snippetService *service.SnippetService) *SnippetHandler {
	return &SnippetHandler{snippetService: snippetService}
}

// Create handles snippet creation
func (h *SnippetHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req service.CreateRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, err)
	}

	snippet, err := h.snippetService.Create(userID, req)
	if err != nil {
		return Error(c, err)
	}

	return Success(c, http.StatusCreated, snippet)
}

// Get retrieves a snippet by namespace and slug
func (h *SnippetHandler) Get(c echo.Context) error {
	namespace := c.Param("namespace")
	slug := c.Param("slug")

	snippet, err := h.snippetService.GetByPath(namespace, slug)
	if err != nil {
		return Error(c, err)
	}

	return Success(c, http.StatusOK, snippet)
}

// Pull renders a snippet with provided variables
func (h *SnippetHandler) Pull(c echo.Context) error {
	namespace := c.Param("namespace")
	slug := c.Param("slug")

	// Get user ID if authenticated (optional for public snippets)
	userID := ""
	if uid := c.Get("user_id"); uid != nil {
		userID = uid.(string)
	}

	var req service.PullRequest
	if err := c.Bind(&req); err != nil {
		return Error(c, err)
	}

	result, err := h.snippetService.Pull(namespace, slug, userID, req.Variables)
	if err != nil {
		return Error(c, err)
	}

	return Success(c, http.StatusOK, result)
}

// Search searches for public snippets
func (h *SnippetHandler) Search(c echo.Context) error {
	query := c.QueryParam("q")
	tags := c.QueryParams()["tags"]

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	snippets, total, err := h.snippetService.Search(query, tags, limit, offset)
	if err != nil {
		return Error(c, err)
	}

	return PaginatedSuccess(c, http.StatusOK, snippets, total, limit, offset)
}

// ListMine lists snippets owned by the current user
func (h *SnippetHandler) ListMine(c echo.Context) error {
	userID := c.Get("user_id").(string)

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	snippets, total, err := h.snippetService.ListMine(userID, limit, offset)
	if err != nil {
		return Error(c, err)
	}

	return PaginatedSuccess(c, http.StatusOK, snippets, total, limit, offset)
}
