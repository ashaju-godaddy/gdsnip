# GDSNIP Implementation Plan

## Status Summary

**Current Status**: MVP + Teams complete! Ready for production use 🚀

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
- **Phase 7**: Seed data & demo environment
  - ✅ Demo user (demo@example.com / demo / password123)
  - ✅ 12 production-ready snippet templates
  - ✅ Integrated with `task db:seed` and `task setup`
- **Phase 8**: Teams & Basic RBAC
  - ✅ Database schema with teams and team_members tables
  - ✅ RBAC permission model (owner, admin, member, viewer roles)
  - ✅ Full team management API and CLI commands
  - ✅ Team snippet support with visibility controls
  - ✅ Comprehensive unit tests

### Post-MVP Priority Order 🚀
1. ✅ **Teams & Basic RBAC** - Organization and collaboration (COMPLETED)
2. **AI Snippet Generation** - Natural language to template + config file conversion (autogen)
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

### Deferred Features
- Update/delete snippets
- ~~Teams and team snippets~~ ✅ (COMPLETED in Phase 8)
- Snippet versioning
- Namespace resolution by slug-only
- Web UI
- AI features (snippet generation, chat)
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

### Phase 7: Seed Data & Polish ✅

**Completed:**
- ✅ `scripts/seed.go` - Comprehensive seed script
  - Creates demo user (demo@example.com / demo / password123)
  - 12 production-ready snippet templates covering:
    - Docker (PostgreSQL, Redis, Node.js, Go)
    - Kubernetes (Deployment, Service)
    - CI/CD (GitHub Actions)
    - Infrastructure (Nginx, Terraform, systemd)
    - Build tools (Makefile)
    - Configuration (env template)
  - Integrated with `task db:seed` and `task setup`
- ✅ Verified with CLI commands (search, info working)
- ✅ All tests passing

---

## Post-MVP Implementation Plan

### Phase 8: Teams & Basic RBAC ✅

**Status**: COMPLETED

**Completed:**

#### 8.1 Database Schema ✅
- ✅ Created migration `002_teams.up.sql` with teams and team_members tables
- ✅ Created migration `002_teams.down.sql` for rollback
- ✅ Migrations tested and verified
- ✅ Database supports team ownership via existing owner_type/owner_id pattern

#### 8.2 RBAC Permission Model ✅
- ✅ Implemented 4-role system: owner, admin, member, viewer
- ✅ Granular permissions (9 permission types)
- ✅ Permission checking at service layer
- ✅ Unit tests for RBAC logic (12 tests passing)

#### 8.3-8.5 Backend Implementation ✅
- ✅ Models: `pkg/models/team.go` with Team, TeamMember, TeamRole
- ✅ Error types: Added 7 team-related error codes
- ✅ Repository: `internal/api/repository/team.go` with full CRUD
- ✅ Service: `internal/api/service/team.go` with permission enforcement
- ✅ Updated SnippetService for team snippet support

#### 8.6 API Handlers ✅
- ✅ Created `internal/api/handlers/teams.go`
- ✅ 11 endpoints implemented:
  - POST/GET/PUT/DELETE `/v1/teams` and `/v1/teams/:slug`
  - Team member management endpoints
  - Team snippet listing
- ✅ All endpoints protected with JWT middleware

#### 8.7 CLI Commands ✅
- ✅ Created `internal/cli/commands/teams.go`
- ✅ 9 subcommands: create, list, info, delete, members, add, remove, role, leave
- ✅ Updated push command with `--team` flag
- ✅ CLI API client methods added

#### 8.8 Testing ✅
- ✅ Unit tests for RBAC (pkg/models/team_test.go)
- ✅ Service tests updated for team support
- ✅ All tests passing

**Files Created/Modified:**
```
migrations/002_teams.up.sql               (new)
migrations/002_teams.down.sql             (new)
pkg/models/team.go                        (new)
pkg/models/team_test.go                   (new)
pkg/models/errors.go                      (modified - added team errors)
internal/api/repository/team.go           (new)
internal/api/service/team.go              (new)
internal/api/service/snippet.go           (modified - team support)
internal/api/service/snippet_test.go      (modified - teamRepo parameter)
internal/api/handlers/teams.go            (new)
internal/api/server/server.go             (modified - wire teams)
internal/cli/api/client.go                (modified - team methods)
internal/cli/commands/teams.go            (new)
internal/cli/commands/push.go             (modified - --team flag)
```

**Original Phase 8 Implementation Details**

<details>
<summary>Click to expand original implementation plan</summary>

#### 8.1 Database Schema (COMPLETED)

Example migration `002_teams.up.sql`:

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

**Acceptance Criteria:** ✅ ALL MET
- ✅ User can create a team and becomes owner
- ✅ Owner can add/remove members with roles
- ✅ Members can create team snippets
- ✅ Viewers can only pull, not create
- ✅ Team snippets visibility controls working
- ✅ Permission checks enforced at service layer

</details>

---

### Phase 9: AI Snippet Generation

**Overview**: This phase adds two AI-powered capabilities:
1. **Generate from scratch**: Natural language description → template
2. **Convert existing files**: Configuration file → parameterized template (autogen)

#### 9.1 Generate from Natural Language

##### 9.1.1 AI Service Layer

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

##### 9.1.2 Prompt Engineering

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

##### 9.1.3 API Endpoint

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

##### 9.1.4 CLI Command

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

##### 9.1.5 TUI Flow

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

#### 9.2 AI-Powered Template Conversion (Autogen)

Convert existing configuration files into parameterized GDSNIP templates with intelligent variable extraction and naming.

##### 9.2.1 Autogen Engine

- [ ] Create `internal/ai/autogen.go`:

```go
type AutogenEngine struct {
    client *Client
}

type AutogenRequest struct {
    Content     string   `json:"content"`
    Filename    string   `json:"filename,omitempty"`
    FileType    string   `json:"file_type,omitempty"` // yaml, json, env, nginx, etc.
    Context     string   `json:"context,omitempty"`    // User-provided context
}

type AutogenResult struct {
    Template         string              `json:"template"`
    Variables        []VariableMetadata  `json:"variables"`
    SuggestedName    string              `json:"suggested_name"`
    SuggestedSlug    string              `json:"suggested_slug"`
    SuggestedDesc    string              `json:"suggested_description"`
    SuggestedTags    []string            `json:"suggested_tags"`
    Confidence       float64             `json:"confidence"`
}

type VariableMetadata struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Example     string `json:"example"`      // Original value
    Category    string `json:"category"`     // "port", "url", "secret", "config"
    Sensitive   bool   `json:"sensitive"`
}

func (e *AutogenEngine) Convert(ctx context.Context, req AutogenRequest) (*AutogenResult, error) {
    // Use AI to intelligently extract variables with semantic understanding
}
```

**Two-Tiered Approach**:

1. **Basic Rule-Based** (fallback/offline mode):
   - Pattern matching for common patterns:
     - URLs: `https?://[^\s]+` → `{{*_URL}}`
     - Ports: `:\d{2,5}` → `{{*_PORT}}`
     - IDs: Numeric sequences → `{{*_ID}}`
     - Repeated strings → Variables
   - Limited semantic understanding

2. **AI-Enhanced** (when API available):
   - Semantic understanding of variable purpose
   - Intelligent variable naming (not just `VAR1`, `VAR2`)
   - Descriptions for each variable
   - Categorization (ports, secrets, configs)
   - Smart grouping (e.g., `DB_HOST`, `DB_PORT`, `DB_NAME`)
   - Detects sensitive fields (`*_PASSWORD`, `*_SECRET`, `*_KEY`)

##### 9.2.2 Autogen Prompts

- [ ] Add to `internal/ai/prompts.go`:

```go
const AutogenSystemPrompt = `You are an expert at converting configuration files into parameterized templates.

Your task is to:
1. Identify ALL values that should be configurable as variables
2. Replace them with {{UPPERCASE_VARIABLE}} placeholders
3. Use semantic, descriptive variable names (not VAR1, VAR2)
4. Group related variables with consistent prefixes (e.g., DB_*, AWS_*, REDIS_*)
5. Identify sensitive values (passwords, keys, tokens, secrets)

Variable naming rules:
- UPPERCASE with underscores only
- Descriptive (e.g., {{DATABASE_HOST}} not {{VAR1}})
- Grouped by component (e.g., {{POSTGRES_USER}}, {{POSTGRES_PASSWORD}}, {{POSTGRES_DB}})
- Common patterns:
  * Ports: {{SERVICE_NAME}}_PORT
  * Hosts: {{SERVICE_NAME}}_HOST
  * Secrets: {{SERVICE_NAME}}_PASSWORD, {{SERVICE_NAME}}_API_KEY
  * Names/IDs: {{RESOURCE_TYPE}}_NAME, {{APPLICATION_ID}}

For each variable, provide:
- name: The variable name
- description: Brief explanation of what it configures
- example: The original value from the file
- category: "url", "port", "secret", "database", "config", etc.
- sensitive: true if it's a password/key/token

Output JSON with:
{
  "template": "the converted template content",
  "variables": [ { "name": "...", "description": "...", "example": "...", "category": "...", "sensitive": false } ],
  "suggested_name": "Human-readable name for this template",
  "suggested_slug": "kebab-case-slug",
  "suggested_description": "What this template is for",
  "suggested_tags": ["technology", "type", "..."]
}`

func BuildAutogenPrompt(req AutogenRequest) string {
    prompt := fmt.Sprintf(`Convert this configuration file into a parameterized template:\n\n%s`, req.Content)

    if req.Filename != "" {
        prompt += fmt.Sprintf("\n\nFilename: %s", req.Filename)
    }

    if req.FileType != "" {
        prompt += fmt.Sprintf("\nFile type: %s", req.FileType)
    }

    if req.Context != "" {
        prompt += fmt.Sprintf("\n\nContext: %s", req.Context)
    }

    return prompt
}
```

##### 9.2.3 API Endpoint

- [ ] Add to `internal/api/handlers/ai.go`:

```go
// POST /v1/ai/autogen
type AutogenRequest struct {
    Content  string `json:"content" validate:"required,min=10,max=50000"`
    Filename string `json:"filename" validate:"max=255"`
    FileType string `json:"file_type" validate:"omitempty,oneof=yaml json toml env nginx dockerfile"`
    Context  string `json:"context" validate:"max=500"`
}

func (h *AIHandler) Autogen(c echo.Context) error {
    // Rate limit autogen requests (e.g., 20/hour)
    // Call autogen service
    // Return converted template with variables
}
```

##### 9.2.4 CLI Command

- [ ] Create `internal/cli/commands/autogen.go`:

```go
var autogenCmd = &cobra.Command{
    Use:   "autogen <file>",
    Short: "Convert a config file to a template",
    Long: `Convert an existing configuration file into a parameterized GDSNIP template.

The AI will intelligently identify configurable values and replace them with
semantic variables. You can review and edit the result before saving.`,
    Example: `  # Convert a YAML file
  gdsnip autogen config.yml

  # Convert and save directly
  gdsnip autogen nginx.conf --save

  # Provide context for better variable naming
  gdsnip autogen app.config --context "API gateway for microservices"

  # Read from stdin
  cat docker-compose.yml | gdsnip autogen -`,
    Args: cobra.ExactArgs(1),
    Run:  runAutogen,
}

func init() {
    autogenCmd.Flags().Bool("save", false, "Save directly without preview")
    autogenCmd.Flags().StringP("context", "c", "", "Additional context about the file")
    autogenCmd.Flags().StringP("format", "f", "", "File format (auto-detected if not specified)")
}

func runAutogen(cmd *cobra.Command, args []string) {
    filepath := args[0]

    // Read file (or stdin if "-")
    content := readFileOrStdin(filepath)

    // Detect file type from extension if not provided
    fileType := detectFileType(filepath)

    // Show spinner during conversion
    // Call API /v1/ai/autogen
    // Display preview with TUI (similar to generate)
    // Allow editing before save
}
```

- [ ] Create `internal/cli/tui/autogen.go`:

```go
// Interactive preview showing:
// - Original content vs template side-by-side
// - Variables table with descriptions
// - Ability to rename variables
// - Ability to mark variables as required/optional
// - Save flow

type autogenModel struct {
    state           autogenState
    originalContent string
    template        string
    variables       []VariableMetadata
    selectedVar     int
    viewport        viewport.Model
    editMode        bool
}
```

##### 9.2.5 Intelligent Features

**Smart Variable Detection**:
```yaml
# Input: gatewayconfig.yml
routes:
  - match:
      prefix: /v1/commerce/location
      headers:
        - name: :method
          string_match:
            safe_regex:
              regex: POST

# AI Output:
routes:
  - match:
      prefix: {{URL_PREFIX}}
      headers:
        - name: :method
          string_match:
            safe_regex:
              regex: {{HTTP_METHOD}}

Variables:
- name: URL_PREFIX
  description: "API endpoint URL path prefix"
  example: "/v1/commerce/location"
  category: "url"
  sensitive: false

- name: HTTP_METHOD
  description: "Allowed HTTP method(s) for this route"
  example: "POST"
  category: "config"
  sensitive: false
```

**Category-Based Grouping**:
- Database vars: `DB_*`, `POSTGRES_*`, `MYSQL_*`
- Network vars: `*_PORT`, `*_HOST`, `*_URL`
- Auth vars: `*_SECRET`, `*_KEY`, `*_TOKEN`, `*_PASSWORD` (marked sensitive)
- App vars: `APP_*`, `SERVICE_*`

**Context-Aware Naming**:
```
Input filename: "postgres-docker-compose.yml"
→ Suggests variables: POSTGRES_VERSION, POSTGRES_PASSWORD, POSTGRES_DB
  (not CONTAINER_IMAGE, PASSWORD, DATABASE)
```

#### 9.3 Shared Configuration

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

#### 9.4 Testing

**Unit Tests**:
- [ ] Mock AI client for testing
- [ ] Test prompt generation (generate & autogen)
- [ ] Test response parsing
- [ ] Test rule-based autogen fallback
- [ ] Test error handling (rate limits, API errors)

**Integration Tests**:
- [ ] Test with real AI API (manual/optional)
- [ ] Test autogen with various file formats (YAML, JSON, env, Dockerfile)
- [ ] Test variable extraction accuracy

**Acceptance Criteria:**

**Generate from Scratch**:
- [ ] `gdsnip generate "description"` calls AI and shows preview
- [ ] User can edit generated content before saving
- [ ] User can save directly with `--save` flag
- [ ] Rate limiting prevents abuse
- [ ] Graceful error handling for API failures

**Autogen**:
- [ ] `gdsnip autogen config.yml` converts file to template
- [ ] Extracted variables have semantic names (not VAR1, VAR2)
- [ ] Each variable includes description and example
- [ ] Sensitive fields (passwords, keys) marked correctly
- [ ] Variables grouped logically (DB_*, REDIS_*, etc.)
- [ ] Side-by-side preview shows original vs template
- [ ] User can rename variables before saving
- [ ] Works with stdin: `cat file.yml | gdsnip autogen -`
- [ ] Fallback to rule-based when AI unavailable

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

### Phase 7: Seed Data ✅ (COMPLETED)
```
scripts/seed.go
```

### Phase 8: Teams ✅ (COMPLETED)
```
migrations/002_teams.up.sql
migrations/002_teams.down.sql
pkg/models/team.go
pkg/models/team_test.go
internal/api/repository/team.go
internal/api/service/team.go
internal/api/handlers/teams.go
internal/cli/commands/teams.go
```

### Phase 9: AI Generation
```
internal/ai/client.go
internal/ai/prompts.go
internal/ai/autogen.go              # NEW: Autogen engine
internal/api/handlers/ai.go
internal/cli/commands/generate.go
internal/cli/commands/autogen.go    # NEW: Autogen CLI command
internal/cli/tui/generate.go
internal/cli/tui/autogen.go         # NEW: Autogen TUI
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