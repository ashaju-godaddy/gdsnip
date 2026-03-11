# GDSNIP Implementation Plan

## Status Summary

**Current Status**: MVP implementation ~95% complete - seed data and polish remaining

### Completed ✅
- **Phase 1**: Project setup, Docker, database schema, migrations
- **Phase 2**: Shared packages (models, template engine, validators) with comprehensive tests
- **Phase 3**: Full API server with authentication, CRUD, search, and pull endpoints
- **Phase 4**: CLI foundation (config, credentials, API client)
- **Phase 5**: CLI commands - ALL 6 core commands implemented
  - ✅ TUI components (styles, prompts, spinner)
  - ✅ Auth commands (register, login, logout, status)
  - ✅ Discovery commands (search, list, info)
  - ✅ Push command with auto-variable extraction
  - ✅ Pull command with custom flag parsing and interactive prompts
- **Phase 6**: README with comprehensive quickstart guide

### MVP Remaining 📋
- Seed data script (optional for demo/testing)
- Basic integration testing (optional)
- Minor polish and bug fixes

### Post-MVP Priority Order 🚀
1. **Teams & Basic RBAC** - Organization and collaboration
2. **AI Snippet Generation** - Natural language to template
3. **Multi-step Interactive TUI Wizards** - Complex snippet creation
4. **Web UI** - Browser-based search and management
5. **AI Chat with Snippet Indexing** - Conversational snippet discovery

---

## MVP Scope

The MVP focuses on the **core pull workflow**: a user can discover, pull, and push templates with variable substitution.

### Core Features (MVP)
1. **Authentication**: Register, login, logout, status
2. **Push**: Create public/private snippets with auto-variable extraction
3. **Pull**: Render templates with variables (interactive prompts for missing vars)
4. **Search**: Full-text search with tag filtering
5. **List**: View your own snippets
6. **Info**: View snippet details including variables

### Explicitly Deferred (Post-MVP)
- Update/delete snippets
- Teams and team snippets
- Snippet versioning
- Namespace resolution by slug-only
- Web UI
- AI features
- Advanced rate limiting
- Metrics/observability
- CLI features: `--json` output, shell completion, config command

---

## Completed Phases

### Phase 1: Project Setup & Foundation ✅

**Completed:**
- ✅ Go module initialized
- ✅ Directory structure created
- ✅ Taskfile with essential commands (`build`, `test`, `db:*`)
- ✅ Docker Compose with PostgreSQL 16
- ✅ Database migrations (up/down) with full schema
- ✅ Full-text search with tsvector
- ✅ Indexes for performance
- ✅ Auto-update triggers for timestamps and search vectors

---

### Phase 2: Shared Packages ✅

**Completed:**
- ✅ `pkg/models/` - User, Snippet, APIError with structured error handling
- ✅ `pkg/template/` - Template engine with extraction, validation, and rendering
- ✅ `pkg/validator/` - Slug, username, email, password validation
- ✅ **24 template engine tests** passing
- ✅ **50+ validator tests** passing

---

### Phase 3: API Server ✅

**Completed:**
- ✅ Koanf-based config loading
- ✅ Database connection with retry logic and pooling
- ✅ Repository layer for users and snippets
- ✅ `AuthService` and `SnippetService` with full test coverage
- ✅ JWT authentication middleware
- ✅ Rate limiting middleware
- ✅ All core endpoints implemented and tested

---

### Phase 4: CLI Foundation ✅

**Completed:**
- ✅ Configuration loader (Koanf with file + env support)
- ✅ Credentials manager (secure 0600 file permissions)
- ✅ Full API client with all endpoints
- ✅ Error handling and APIError parsing

---

### Phase 5: CLI Commands ✅

**Completed:**
- ✅ TUI components (styles, prompts, spinner)
- ✅ Auth commands (register, login, logout, status)
- ✅ Discovery commands (search, list, info)
- ✅ Push command with auto-variable extraction
- ✅ Pull command with interactive prompts

---

### Phase 6: README Documentation ✅

**Completed:**
- ✅ Comprehensive quickstart guide
- ✅ Command reference with examples
- ✅ Template syntax documentation
- ✅ Configuration and troubleshooting guides

---

## MVP Remaining Work

### Phase 7: Seed Data & Polish

- [ ] `scripts/seed.go` - Create demo user + 10+ realistic snippets
- [ ] Basic integration testing
- [ ] Final bug fixes

---

## Post-MVP Implementation Plan

### Phase 8: Teams & Basic RBAC

#### 8.1 Database Schema

- [ ] Create migration `002_teams.up.sql`:

```sql
-- Teams table
CREATE TABLE teams (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    slug            VARCHAR(100) UNIQUE NOT NULL,
    description     TEXT,
    created_by      UUID REFERENCES users(id) NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- Team memberships with roles
CREATE TABLE team_members (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id         UUID REFERENCES teams(id) ON DELETE CASCADE,
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(20) NOT NULL DEFAULT 'member',
    -- Roles: 'owner', 'admin', 'member', 'viewer'
    created_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(team_id, user_id)
);

-- Add team ownership to snippets
ALTER TABLE snippets ADD COLUMN team_id UUID REFERENCES teams(id) ON DELETE CASCADE;

-- Update visibility enum concept
-- visibility: 'public', 'private', 'team'

-- Indexes
CREATE INDEX idx_teams_slug ON teams(slug);
CREATE INDEX idx_teams_created_by ON teams(created_by);
CREATE INDEX idx_team_members_team ON team_members(team_id);
CREATE INDEX idx_team_members_user ON team_members(user_id);
CREATE INDEX idx_snippets_team ON snippets(team_id);
```

- [ ] Create migration `002_teams.down.sql`

#### 8.2 RBAC Permission Model

Define permission levels:

```go
// pkg/models/team.go

type TeamRole string

const (
    RoleOwner  TeamRole = "owner"   // Full control, can delete team
    RoleAdmin  TeamRole = "admin"   // Manage members, manage all snippets
    RoleMember TeamRole = "member"  // Create/edit own snippets, pull team snippets
    RoleViewer TeamRole = "viewer"  // Pull team snippets only
)

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

func (r TeamRole) HasPermission(p Permission) bool {
    for _, perm := range RolePermissions[r] {
        if perm == p {
            return true
        }
    }
    return false
}
```

#### 8.3 Models

- [ ] Create `pkg/models/team.go`:

```go
type Team struct {
    ID          string    `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Slug        string    `json:"slug" db:"slug"`
    Description string    `json:"description" db:"description"`
    CreatedBy   string    `json:"created_by" db:"created_by"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    
    // Populated on fetch
    MemberCount int       `json:"member_count,omitempty" db:"member_count"`
    Role        TeamRole  `json:"role,omitempty" db:"role"` // Current user's role
}

type TeamMember struct {
    ID        string    `json:"id" db:"id"`
    TeamID    string    `json:"team_id" db:"team_id"`
    UserID    string    `json:"user_id" db:"user_id"`
    Role      TeamRole  `json:"role" db:"role"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    
    // Populated via JOIN
    Username  string    `json:"username,omitempty" db:"username"`
    Email     string    `json:"email,omitempty" db:"email"`
}

type CreateTeamRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=100"`
    Slug        string `json:"slug" validate:"omitempty,slug"`
    Description string `json:"description" validate:"max=500"`
}

type AddMemberRequest struct {
    Username string   `json:"username" validate:"required"`
    Role     TeamRole `json:"role" validate:"required,oneof=admin member viewer"`
}
```

#### 8.4 Repository Layer

- [ ] Create `internal/api/repository/team.go`:
  - `Create(team *Team) error`
  - `GetByID(id string) (*Team, error)`
  - `GetBySlug(slug string) (*Team, error)`
  - `Update(team *Team) error`
  - `Delete(id string) error`
  - `ListByUser(userID string) ([]Team, error)`
  - `AddMember(teamID, userID string, role TeamRole) error`
  - `RemoveMember(teamID, userID string) error`
  - `UpdateMemberRole(teamID, userID string, role TeamRole) error`
  - `GetMember(teamID, userID string) (*TeamMember, error)`
  - `ListMembers(teamID string) ([]TeamMember, error)`
  - `GetUserRole(teamID, userID string) (TeamRole, error)`

#### 8.5 Service Layer

- [ ] Create `internal/api/service/team.go`:
  - `CreateTeam(userID string, req CreateTeamRequest) (*Team, error)`
    - Auto-generate slug if not provided
    - Creator becomes owner
  - `GetTeam(userID, teamSlug string) (*Team, error)`
    - Verify user is member
  - `UpdateTeam(userID, teamSlug string, req UpdateTeamRequest) (*Team, error)`
    - Check `PermManageTeam`
  - `DeleteTeam(userID, teamSlug string) error`
    - Check `PermDeleteTeam` (owner only)
  - `ListUserTeams(userID string) ([]Team, error)`
  - `AddMember(userID, teamSlug string, req AddMemberRequest) (*TeamMember, error)`
    - Check `PermManageMembers`
    - Cannot add owner role via API
  - `RemoveMember(userID, teamSlug, targetUsername string) error`
    - Check `PermManageMembers`
    - Owner cannot be removed
  - `UpdateMemberRole(userID, teamSlug, targetUsername string, role TeamRole) error`
    - Check `PermManageMembers`
    - Cannot change owner
  - `LeaveTeam(userID, teamSlug string) error`
    - Owner cannot leave (must transfer or delete)

- [ ] Update `internal/api/service/snippet.go`:
  - Add `teamSlug` parameter to `Create`
  - Add permission checks for team snippets
  - Update `Search` to include team snippets user can access
  - Update `ListByUser` to support `--team` filter

#### 8.6 API Handlers

- [ ] Create `internal/api/handlers/teams.go`:

```go
// Team CRUD
POST   /v1/teams                           // Create team
GET    /v1/teams                           // List my teams
GET    /v1/teams/:slug                     // Get team details
PUT    /v1/teams/:slug                     // Update team
DELETE /v1/teams/:slug                     // Delete team

// Team membership
GET    /v1/teams/:slug/members             // List members
POST   /v1/teams/:slug/members             // Add member
PUT    /v1/teams/:slug/members/:username   // Update role
DELETE /v1/teams/:slug/members/:username   // Remove member
POST   /v1/teams/:slug/leave               // Leave team

// Team snippets
GET    /v1/teams/:slug/snippets            // List team snippets
```

- [ ] Update snippet handlers:
  - `POST /v1/snippets` accepts optional `team_slug` field
  - Team snippets use `team` visibility by default

#### 8.7 CLI Commands

- [ ] Create `internal/cli/commands/teams.go`:

```bash
# Team management
gdsnip teams create <name> [--slug=<slug>] [--description=<desc>]
gdsnip teams list
gdsnip teams info <slug>
gdsnip teams delete <slug>

# Membership
gdsnip teams members <slug>
gdsnip teams add <slug> <username> [--role=member]
gdsnip teams remove <slug> <username>
gdsnip teams role <slug> <username> <role>
gdsnip teams leave <slug>
```

- [ ] Update `internal/cli/commands/push.go`:
  - Add `--team` flag to push snippet to team namespace
  - Example: `gdsnip push -f template.yml -n "shared config" --team platform`

- [ ] Update `internal/cli/commands/search.go`:
  - Add `--team` flag to search within team
  - Show team snippets user has access to

- [ ] Update `internal/cli/commands/list.go`:
  - Add `--team` flag to list team's snippets

#### 8.8 Testing

- [ ] Unit tests for RBAC permission logic
- [ ] Service tests for team operations
- [ ] Integration tests for team workflows:
  - Create team → add members → create snippet → verify access
  - Permission denied scenarios

**Acceptance Criteria:**
- [ ] User can create a team and becomes owner
- [ ] Owner can add/remove members with roles
- [ ] Members can create team snippets
- [ ] Viewers can only pull, not create
- [ ] Team snippets appear in search for members
- [ ] Non-members cannot access team snippets

---

### Phase 9: AI Snippet Generation

#### 9.1 AI Service Layer

- [ ] Create `internal/ai/client.go`:

```go
type Client struct {
    apiKey   string
    model    string
    baseURL  string
    http     *http.Client
}

type GenerateRequest struct {
    Description string   `json:"description"`
    Tags        []string `json:"tags,omitempty"`
    Format      string   `json:"format,omitempty"` // yaml, json, toml, dockerfile, etc.
}

type GenerateResult struct {
    Content          string   `json:"content"`
    SuggestedName    string   `json:"suggested_name"`
    SuggestedSlug    string   `json:"suggested_slug"`
    SuggestedDesc    string   `json:"suggested_description"`
    SuggestedTags    []string `json:"suggested_tags"`
    ExtractedVars    []string `json:"extracted_variables"`
    Confidence       float64  `json:"confidence"` // 0-1
}

func NewClient(apiKey string) *Client
func (c *Client) GenerateSnippet(ctx context.Context, req GenerateRequest) (*GenerateResult, error)
```

- [ ] Create `internal/ai/prompts.go`:

```go
const SystemPrompt = `You are a configuration template generator. Generate templates with these rules:

1. Use {{UPPERCASE_VARIABLE}} syntax for ALL configurable values
2. Variables must be uppercase with underscores only
3. Include sensible inline comments
4. Follow best practices for the technology
5. Use common variable names:
   - Ports: PORT, HTTP_PORT, HTTPS_PORT
   - Database: DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD
   - App: APP_NAME, APP_ENV, APP_VERSION
   - Secrets: *_SECRET, *_KEY, *_TOKEN, *_PASSWORD

Output JSON with:
{
  "content": "the template content",
  "suggested_name": "Human readable name",
  "suggested_slug": "kebab-case-slug",
  "suggested_description": "Brief description",
  "suggested_tags": ["tag1", "tag2"]
}`

func BuildGeneratePrompt(req GenerateRequest) string {
    return fmt.Sprintf(`Generate a configuration template for:
%s

%s%s`,
        req.Description,
        formatHint(req.Format),
        tagHint(req.Tags),
    )
}
```

#### 9.2 API Endpoint

- [ ] Create `internal/api/handlers/ai.go`:

```go
// POST /v1/ai/generate
type GenerateRequest struct {
    Description string   `json:"description" validate:"required,min=10,max=1000"`
    Tags        []string `json:"tags" validate:"max=10"`
    Format      string   `json:"format" validate:"omitempty,oneof=yaml json toml dockerfile nginx"`
}

func (h *AIHandler) Generate(c echo.Context) error {
    // Rate limit AI requests more strictly (e.g., 10/hour)
    // Call AI service
    // Return generated snippet (not saved yet)
}
```

- [ ] Add rate limiting specifically for AI endpoints (stricter limits)

#### 9.3 CLI Command

- [ ] Create `internal/cli/commands/generate.go`:

```go
var generateCmd = &cobra.Command{
    Use:   "generate <description>",
    Short: "Generate a snippet using AI",
    Long: `Generate a configuration template from natural language description.

The AI will create a template with appropriate variables that you can
review, edit, and save to your snippets.`,
    Example: `  # Generate a Docker Compose for PostgreSQL
  gdsnip generate "docker compose for postgres with persistent storage"

  # Generate with specific format
  gdsnip generate "kubernetes deployment for node.js app" --format=yaml

  # Generate with suggested tags
  gdsnip generate "nginx reverse proxy with SSL" --tags=nginx,proxy,ssl`,
    Args: cobra.MinimumNArgs(1),
    Run:  runGenerate,
}

func init() {
    generateCmd.Flags().StringP("format", "f", "", "Output format hint (yaml, json, dockerfile, etc.)")
    generateCmd.Flags().StringSliceP("tags", "t", nil, "Suggested tags")
    generateCmd.Flags().Bool("save", false, "Save directly without preview")
}
```

- [ ] Create `internal/cli/tui/generate.go` - Generation flow TUI:

```go
// Step 1: Show spinner during generation
// Step 2: Display generated content in preview box
// Step 3: Show extracted variables
// Step 4: Confirm or edit before saving
// Step 5: Prompt for name/visibility if saving

type generateModel struct {
    state       generateState
    spinner     spinner.Model
    content     string
    variables   []string
    suggestions GenerateSuggestions
    editMode    bool
    textarea    textarea.Model
}

type generateState int

const (
    stateGenerating generateState = iota
    statePreview
    stateEditing
    stateMetadata
    stateSaving
    stateDone
)
```

#### 9.4 TUI Flow

```
┌─────────────────────────────────────────────────────────────┐
│  🤖 AI Snippet Generator                                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ⠋ Generating template...                                   │
│                                                             │
└─────────────────────────────────────────────────────────────┘

        ↓ After generation ↓

┌─────────────────────────────────────────────────────────────┐
│  🤖 Generated Snippet                                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  version: "3.8"                                             │
│  services:                                                  │
│    postgres:                                                │
│      image: postgres:{{PG_VERSION}}                         │
│      environment:                                           │
│        POSTGRES_DB: {{DB_NAME}}                             │
│        POSTGRES_PASSWORD: {{DB_PASSWORD}}                   │
│      ...                                                    │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│  Variables: PG_VERSION, DB_NAME, DB_PASSWORD, DB_USER, PORT │
├─────────────────────────────────────────────────────────────┤
│  Suggested: "Docker PostgreSQL" • docker, postgres          │
├─────────────────────────────────────────────────────────────┤
│  [s] Save  [e] Edit  [r] Regenerate  [q] Quit               │
└─────────────────────────────────────────────────────────────┘
```

#### 9.5 Configuration

- [ ] Add AI configuration to `internal/api/config/config.go`:

```go
type AIConfig struct {
    Enabled     bool   `koanf:"enabled"`
    Provider    string `koanf:"provider"` // openai, anthropic
    APIKey      string `koanf:"api_key"`
    Model       string `koanf:"model"`    // gpt-4, claude-3-sonnet
    MaxTokens   int    `koanf:"max_tokens"`
    RateLimit   int    `koanf:"rate_limit"` // requests per hour per user
}
```

- [ ] Add CLI AI configuration to `internal/cli/config/config.go`:

```go
type Config struct {
    APIURL      string `koanf:"api_url"`
    AIEnabled   bool   `koanf:"ai_enabled"`
}
```

#### 9.6 Testing

- [ ] Mock AI client for testing
- [ ] Test prompt generation
- [ ] Test response parsing
- [ ] Test error handling (rate limits, API errors)
- [ ] Integration test with real API (manual/optional)

**Acceptance Criteria:**
- [ ] `gdsnip generate "description"` calls AI and shows preview
- [ ] User can edit generated content before saving
- [ ] User can save directly with `--save` flag
- [ ] Rate limiting prevents abuse
- [ ] Graceful error handling for API failures

---

### Phase 10: Multi-Step Interactive TUI Wizards

#### 10.1 Wizard Framework

- [ ] Create `internal/cli/tui/wizard/wizard.go`:

```go
// Generic multi-step wizard framework

type Step interface {
    Title() string
    Description() string
    View() string
    Update(msg tea.Msg) (Step, tea.Cmd)
    Validate() error
    Value() interface{}
}

type Wizard struct {
    steps       []Step
    currentStep int
    title       string
    completed   bool
    cancelled   bool
}

func New(title string, steps ...Step) *Wizard
func (w *Wizard) Init() tea.Cmd
func (w *Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (w *Wizard) View() string
func (w *Wizard) Run() (map[string]interface{}, error)
```

- [ ] Create common step types in `internal/cli/tui/wizard/steps.go`:

```go
// TextInputStep - single line text input
type TextInputStep struct {
    label       string
    placeholder string
    required    bool
    validate    func(string) error
    value       string
}

// PasswordStep - hidden input
type PasswordStep struct {
    label    string
    required bool
    value    string
}

// SelectStep - single selection from options
type SelectStep struct {
    label    string
    options  []string
    selected int
}

// MultiSelectStep - multiple selection
type MultiSelectStep struct {
    label    string
    options  []string
    selected map[int]bool
}

// ConfirmStep - yes/no confirmation
type ConfirmStep struct {
    label    string
    value    bool
}

// TextAreaStep - multi-line input
type TextAreaStep struct {
    label       string
    placeholder string
    value       string
}
```

#### 10.2 Snippet Creation Wizard

- [ ] Create `internal/cli/tui/wizard/create_snippet.go`:

```go
func NewCreateSnippetWizard() *Wizard {
    return New("Create New Snippet",
        // Step 1: Basic info
        &TextInputStep{
            label:       "Snippet Name",
            placeholder: "e.g., Docker PostgreSQL",
            required:    true,
        },
        &TextAreaStep{
            label:       "Description",
            placeholder: "Brief description of what this template does",
        },
        
        // Step 2: Template content
        &TextAreaStep{
            label:       "Template Content",
            placeholder: "Paste your template here...",
            required:    true,
        },
        
        // Step 3: Variables (auto-detected, allow editing)
        &VariableConfigStep{
            // Shows detected variables
            // Allows marking required/optional
            // Allows setting defaults
        },
        
        // Step 4: Metadata
        &MultiSelectStep{
            label:   "Tags",
            options: commonTags, // docker, kubernetes, etc.
        },
        &SelectStep{
            label:   "Visibility",
            options: []string{"public", "private", "team"},
        },
        
        // Step 5: Confirmation
        &ConfirmStep{
            label: "Create this snippet?",
        },
    )
}
```

- [ ] Add command `gdsnip create` (alias for wizard):

```go
var createCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a snippet using interactive wizard",
    Run: func(cmd *cobra.Command, args []string) {
        wizard := wizard.NewCreateSnippetWizard()
        result, err := wizard.Run()
        if err != nil {
            // Handle cancel or error
        }
        // Create snippet via API
    },
}
```

#### 10.3 Pull Wizard (Complex Snippets)

- [ ] Create `internal/cli/tui/wizard/pull.go`:

For snippets with many variables, offer a wizard mode:

```go
func NewPullWizard(snippet *models.Snippet) *Wizard {
    steps := make([]Step, 0)
    
    // Group variables by category (if we add categories)
    // Or just present them in a nice paginated way
    
    for _, v := range snippet.Variables {
        if v.Sensitive {
            steps = append(steps, &PasswordStep{
                label:    fmt.Sprintf("%s - %s", v.Name, v.Description),
                required: v.Required,
            })
        } else {
            steps = append(steps, &TextInputStep{
                label:       fmt.Sprintf("%s - %s", v.Name, v.Description),
                placeholder: v.Default,
                required:    v.Required,
            })
        }
    }
    
    return New(fmt.Sprintf("Pull: %s", snippet.FullPath()), steps...)
}
```

- [ ] Add `--wizard` flag to pull command:

```bash
gdsnip pull demo/complex-snippet --wizard
```

#### 10.4 Team Setup Wizard

- [ ] Create `internal/cli/tui/wizard/create_team.go`:

```go
func NewCreateTeamWizard() *Wizard {
    return New("Create New Team",
        &TextInputStep{
            label:    "Team Name",
            required: true,
        },
        &TextInputStep{
            label:       "Slug",
            placeholder: "auto-generated from name",
        },
        &TextAreaStep{
            label: "Description",
        },
        &ConfirmStep{
            label: "Create this team?",
        },
    )
}
```

#### 10.5 Visual Design

```
┌─────────────────────────────────────────────────────────────┐
│  Create New Snippet                            Step 1 of 5  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Snippet Name *                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Docker PostgreSQL█                                   │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  A short, descriptive name for your snippet                │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│  [←] Back    [Enter] Next    [Esc] Cancel                   │
│                                                             │
│  ○ ○ ○ ○ ○                                                  │
└─────────────────────────────────────────────────────────────┘
```

**Acceptance Criteria:**
- [ ] `gdsnip create` launches multi-step wizard
- [ ] Each step validates before proceeding
- [ ] User can go back to previous steps
- [ ] Escape cancels gracefully with confirmation
- [ ] Progress indicator shows current step
- [ ] Final confirmation shows summary before creating

---

### Phase 11: Web UI

#### 11.1 Technology Choice

**Recommended**: Next.js 14 with App Router
- Server components for initial load
- Client components for interactivity
- API routes proxy to Go backend
- Tailwind CSS for styling

Alternative: React + Vite (simpler, static deployment)

#### 11.2 Project Setup

- [ ] Create `web/` directory
- [ ] Initialize Next.js project:

```bash
cd web
npx create-next-app@latest . --typescript --tailwind --app --src-dir
```

- [ ] Configure proxy to API backend (development)
- [ ] Setup environment variables

#### 11.3 Pages

- [ ] **Home/Search** (`/`):
  - Search bar with instant results
  - Tag filter sidebar
  - Trending/popular snippets
  - Grid/list view toggle

- [ ] **Snippet Detail** (`/snippets/[namespace]/[slug]`):
  - Template preview with syntax highlighting
  - Variables table
  - "Try it" form - input variables, see rendered preview
  - Copy button for rendered output
  - Pull command display

- [ ] **User Profile** (`/users/[username]`):
  - User's public snippets
  - Stats (total snippets, total pulls)

- [ ] **Dashboard** (`/dashboard`) - Protected:
  - My snippets list
  - Quick actions (create, view analytics)
  - Team snippets (if teams enabled)

- [ ] **Create Snippet** (`/create`) - Protected:
  - Form matching CLI push
  - Template editor with syntax highlighting
  - Live variable extraction preview
  - Variable configuration

- [ ] **Auth Pages**:
  - `/login` - Login form
  - `/register` - Registration form
  - OAuth buttons (future)

#### 11.4 Components

```
web/src/components/
├── layout/
│   ├── Header.tsx
│   ├── Footer.tsx
│   └── Sidebar.tsx
├── snippets/
│   ├── SnippetCard.tsx
│   ├── SnippetGrid.tsx
│   ├── SnippetDetail.tsx
│   ├── VariablesTable.tsx
│   ├── VariablesForm.tsx
│   └── TemplatePreview.tsx
├── search/
│   ├── SearchBar.tsx
│   ├── TagFilter.tsx
│   └── SearchResults.tsx
├── editor/
│   ├── TemplateEditor.tsx
│   └── VariableConfig.tsx
├── auth/
│   ├── LoginForm.tsx
│   └── RegisterForm.tsx
└── ui/
    ├── Button.tsx
    ├── Input.tsx
    ├── Card.tsx
    └── Modal.tsx
```

#### 11.5 API Client

- [ ] Create `web/src/lib/api.ts`:

```typescript
export class APIClient {
  private baseURL: string;
  private token: string | null;
  
  constructor() {
    this.baseURL = process.env.NEXT_PUBLIC_API_URL || '/api/v1';
    this.token = null;
  }
  
  setToken(token: string) { ... }
  clearToken() { ... }
  
  // Auth
  async register(data: RegisterRequest): Promise<AuthResponse>
  async login(data: LoginRequest): Promise<AuthResponse>
  async me(): Promise<User>
  
  // Snippets
  async search(query: string, tags?: string[]): Promise<SearchResponse>
  async getSnippet(namespace: string, slug: string): Promise<Snippet>
  async pullSnippet(namespace: string, slug: string, variables: Record<string, string>): Promise<PullResponse>
  async createSnippet(data: CreateSnippetRequest): Promise<Snippet>
  async listMySnippets(): Promise<Snippet[]>
}
```

#### 11.6 Key Features

- [ ] **Live Template Preview**:
  - Enter variables on snippet page
  - See rendered output in real-time
  - Copy to clipboard

- [ ] **Syntax Highlighting**:
  - Use Prism.js or Shiki
  - Highlight `{{VARIABLES}}` specially

- [ ] **Search with URL State**:
  - `/search?q=docker&tags=postgres,database`
  - Shareable search URLs

- [ ] **Responsive Design**:
  - Mobile-friendly search and browse
  - Collapsed sidebar on mobile

#### 11.7 Deployment

- [ ] Dockerfile for web container
- [ ] Add to docker-compose.yml
- [ ] Configure nginx or Caddy for reverse proxy
- [ ] Environment-based API URL

**Acceptance Criteria:**
- [ ] Can search snippets without login
- [ ] Can view snippet details with variable preview
- [ ] Can login and view dashboard
- [ ] Can create snippets via web form
- [ ] Responsive on mobile devices

---

### Phase 12: AI Chat with Snippet Indexing (Post-MVP)

#### 12.1 Snippet Embedding & Indexing

- [ ] Add vector storage (options):
  - PostgreSQL with pgvector extension
  - Pinecone/Weaviate (external service)
  - Local embedding with SQLite + faiss

- [ ] Create embedding pipeline:

```go
// internal/ai/embeddings.go

type EmbeddingService struct {
    client    *Client
    storage   VectorStorage
}

func (s *EmbeddingService) IndexSnippet(snippet *models.Snippet) error {
    // Create embedding from name + description + content sample
    text := fmt.Sprintf("%s\n%s\n%s",
        snippet.Name,
        snippet.Description,
        truncate(snippet.Content, 500),
    )
    
    embedding, err := s.client.CreateEmbedding(text)
    if err != nil {
        return err
    }
    
    return s.storage.Store(snippet.ID, embedding)
}

func (s *EmbeddingService) SemanticSearch(query string, limit int) ([]string, error) {
    queryEmbedding, err := s.client.CreateEmbedding(query)
    if err != nil {
        return nil, err
    }
    
    return s.storage.Search(queryEmbedding, limit)
}
```

- [ ] Background job to index all snippets
- [ ] Index on snippet create/update
- [ ] Remove from index on delete

#### 12.2 Chat Interface

- [ ] Create `internal/ai/chat.go`:

```go
type ChatService struct {
    client         *Client
    embeddings     *EmbeddingService
    snippetService *service.SnippetService
}

type ChatMessage struct {
    Role    string `json:"role"` // user, assistant, system
    Content string `json:"content"`
}

type ChatRequest struct {
    Messages []ChatMessage `json:"messages"`
    UserID   string        `json:"-"`
}

type ChatResponse struct {
    Message         ChatMessage `json:"message"`
    RelevantSnippets []string   `json:"relevant_snippets,omitempty"`
    SuggestedAction  string     `json:"suggested_action,omitempty"` // pull, generate, search
}

func (s *ChatService) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
    // 1. Extract intent from latest message
    // 2. If snippet-related, do semantic search
    // 3. Include relevant snippets in context
    // 4. Generate response with Claude/GPT
    // 5. Parse response for actions
}
```

#### 12.3 CLI Chat Command

- [ ] Create `internal/cli/commands/chat.go`:

```go
var chatCmd = &cobra.Command{
    Use:   "chat",
    Short: "Chat with AI assistant about snippets",
    Long: `Interactive chat interface to discover snippets, get help,
and generate new templates through conversation.`,
    Run: runChat,
}
```

- [ ] Create `internal/cli/tui/chat.go`:

```go
type chatModel struct {
    messages    []ChatMessage
    input       textinput.Model
    viewport    viewport.Model
    typing      bool
    width       int
    height      int
}

// Render chat bubbles with different styles for user/assistant
// Show typing indicator when waiting for response
// Display suggested snippets as clickable items
// Support commands like /pull, /generate, /search
```

#### 12.4 Chat UI Example

```
┌─────────────────────────────────────────────────────────────┐
│  🤖 GDSNIP Assistant                              [Esc] Exit │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ You: I need to set up a kubernetes deployment       │   │
│  │      for a python flask app with redis              │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 🤖 I found some relevant snippets:                  │   │
│  │                                                     │   │
│  │  📦 demo/k8s-deployment - Basic K8s deployment     │   │
│  │  📦 demo/k8s-flask - Flask app deployment          │   │
│  │  📦 demo/docker-redis - Redis setup                │   │
│  │                                                     │   │
│  │ Would you like me to:                              │   │
│  │  [1] Generate a combined template                  │   │
│  │  [2] Pull k8s-flask and customize it               │   │
│  │  [3] Show more details about these snippets        │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│  > █                                                        │
│  Type a message or /help for commands                       │
└─────────────────────────────────────────────────────────────┘
```

#### 12.5 Chat Commands

```
/help           - Show available commands
/search <query> - Search snippets
/pull <path>    - Pull a snippet
/generate       - Generate new snippet
/clear          - Clear chat history
/exit           - Exit chat
```

**Acceptance Criteria:**
- [ ] Chat can answer questions about available snippets
- [ ] Semantic search finds relevant snippets
- [ ] Can pull snippets directly from chat
- [ ] Can generate new snippets through conversation
- [ ] Chat history persisted during session

---

## Updated Technical Notes

### AI Integration Considerations

1. **API Key Security**:
   - Server-side only for API key storage
   - CLI proxies through backend for AI calls
   - Never expose API key to client

2. **Cost Management**:
   - Strict rate limiting per user
   - Token counting and limits
   - Consider caching common generations

3. **Fallback Behavior**:
   - Graceful degradation when AI unavailable
   - Clear error messages about rate limits

### RBAC Implementation Notes

1. **Permission Checks**:
   - Always check at service layer
   - Cache user roles per request
   - Fail closed (deny by default)

2. **Audit Trail** (Future):
   - Log permission checks
   - Track who did what to team resources

### Web UI Notes

1. **Authentication**:
   - Store JWT in httpOnly cookie (not localStorage)
   - Refresh token flow for long sessions

2. **Performance**:
   - Server-side render search results
   - Client-side for interactive features
   - Cache snippet metadata

---

## Phase Priority Summary

| Phase | Feature | Impact | Effort | Priority |
|-------|---------|--------|--------|----------|
| 8 | Teams & Basic RBAC | High | Medium | P1 |
| 9 | AI Snippet Generation | Very High | Medium | P1 |
| 10 | Multi-step TUI Wizards | Medium | Medium | P2 |
| 11 | Web UI | High | High | P2 |
| 12 | AI Chat + Indexing | Very High | High | P3 |

**Recommended Order**: Phase 8 → Phase 9 → Phase 10 → Phase 11 → Phase 12

---

## Quick Reference: New Files to Create

### Phase 8: Teams
```
migrations/002_teams.up.sql
migrations/002_teams.down.sql
pkg/models/team.go
internal/api/repository/team.go
internal/api/service/team.go
internal/api/service/team_test.go
internal/api/handlers/teams.go
internal/cli/commands/teams.go
```

### Phase 9: AI Generation
```
internal/ai/client.go
internal/ai/prompts.go
internal/api/handlers/ai.go
internal/cli/commands/generate.go
internal/cli/tui/generate.go
```

### Phase 10: Wizards
```
internal/cli/tui/wizard/wizard.go
internal/cli/tui/wizard/steps.go
internal/cli/tui/wizard/create_snippet.go
internal/cli/tui/wizard/pull.go
internal/cli/tui/wizard/create_team.go
internal/cli/commands/create.go
```

### Phase 11: Web UI
```
web/                          # New Next.js project
web/src/app/page.tsx          # Home/Search
web/src/app/snippets/[namespace]/[slug]/page.tsx
web/src/app/dashboard/page.tsx
web/src/app/create/page.tsx
web/src/lib/api.ts
web/src/components/...
```

### Phase 12: AI Chat
```
internal/ai/embeddings.go
internal/ai/chat.go
internal/cli/commands/chat.go
internal/cli/tui/chat.go
```

---

## Definition of Done (Updated)

A feature is complete when:

1. **Code**: Implementation matches specification
2. **Tests**: Unit and integration tests pass
3. **Docs**: README updated if user-facing
4. **Manual**: Feature works end-to-end
5. **Errors**: Graceful handling with helpful messages
6. **Security**: Permission checks where applicable