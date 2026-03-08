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
- Basic integration testing (optional for Day 8)
- Minor polish and bug fixes

### Post-MVP (Deferred) 🔮
- Update/Delete snippet endpoints
- Teams functionality
- Versioning system
- Web UI
- Analytics and trending
- Forking
- CI/CD pipelines
- Cloud deployment

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
- Namespace resolution by slug-only (e.g., `gdsnip pull docker-pg` without namespace)
- Web UI
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

**Key Files:**
- `Taskfile.yml`
- `docker-compose.yml`
- `migrations/001_initial.up.sql`
- `migrations/001_initial.down.sql`

---

### Phase 2: Shared Packages ✅

**Completed:**
- ✅ `pkg/models/` - User, Snippet, APIError with structured error handling
- ✅ `pkg/template/` - Template engine with extraction, validation, and rendering
  - Supports `{{UPPERCASE}}` variable syntax
  - Returns ALL missing variables (not just first)
  - Applies defaults before provided values
  - Warns about unused variables
- ✅ `pkg/validator/` - Slug, username, email, password validation
- ✅ **24 template engine tests** passing
- ✅ **50+ validator tests** passing

**Key Files:**
- `pkg/models/user.go`, `snippet.go`, `errors.go`
- `pkg/template/engine.go`, `engine_test.go`
- `pkg/validator/slug.go`, `slug_test.go`

---

### Phase 3: API Server ✅

**Completed:**

#### Configuration & Database
- ✅ Koanf-based config loading (file + environment variables)
- ✅ Database connection with retry logic and pooling
- ✅ Repository layer for users and snippets

#### Service Layer (Business Logic)
- ✅ `AuthService` - Registration, login, JWT generation/parsing
- ✅ `SnippetService` - Create, get, pull, search, list
  - Auto-slug generation
  - Auto-variable extraction from templates
  - Template rendering with validation
  - Pull count tracking (background goroutine)
  - Visibility enforcement (public/private)
- ✅ **14 AuthService tests** passing
- ✅ **18 SnippetService tests** passing (comprehensive coverage)

#### HTTP Layer
- ✅ JWT authentication middleware
- ✅ Rate limiting middleware (in-memory, 100 req/min default)
- ✅ Thin handlers delegating to services
- ✅ Structured error responses
- ✅ Server with graceful shutdown

#### API Endpoints
```
Health:
  GET  /health

Auth:
  POST /v1/auth/register
  POST /v1/auth/login
  GET  /v1/auth/me (protected)

Snippets:
  POST /v1/snippets (protected)                   # Create
  GET  /v1/snippets/:namespace/:slug              # Get details
  POST /v1/snippets/:namespace/:slug/pull         # Pull/render
  GET  /v1/snippets?q=query&tags=tag              # Search
  GET  /v1/users/me/snippets (protected)          # List mine
```

**Tested with curl** - All endpoints working correctly ✅

**Key Files:**
- `internal/api/config/config.go`
- `internal/api/repository/db.go`, `user.go`, `snippet.go`
- `internal/api/service/auth.go`, `snippet.go`
- `internal/api/middleware/auth.go`, `ratelimit.go`
- `internal/api/handlers/auth.go`, `snippets.go`, `response.go`
- `internal/api/server/server.go`
- `cmd/api/main.go`

---

### Phase 4: CLI Foundation ✅

**Completed:**
- ✅ Configuration loader (Koanf with file + env support)
- ✅ Credentials manager (secure 0600 file permissions)
- ✅ Full API client with all endpoints
- ✅ Error handling and APIError parsing

**Key Files:**
- `internal/cli/config/config.go`, `credentials.go`
- `internal/cli/api/client.go`

**Credentials Security:**
- Stored in `~/.gdsnip/credentials.json`
- Directory: 0700 permissions
- File: 0600 permissions (owner read/write only)

---

## Completed Work (Days 1-7)

### Phase 5: CLI Commands ✅ (Days 6-7)

#### TUI Components ✅
- [x] `internal/cli/tui/styles.go` - Lipgloss color palette and reusable styles
- [x] `internal/cli/tui/prompt.go` - Variable prompts with password mode
- [x] `internal/cli/tui/spinner.go` - Animated spinner for operations

#### Root Command ✅
- [x] `internal/cli/commands/root.go` - Cobra root with version injection
- [x] `cmd/cli/main.go` - Entry point

#### Auth Commands ✅
- [x] `internal/cli/commands/auth.go`:
  - `auth register` - Create account, save credentials
  - `auth login` - Authenticate, save credentials with hidden password input
  - `auth logout` - Clear credentials
  - `auth status` - Show current user or helpful "not logged in" message

#### Discovery Commands ✅
- [x] `internal/cli/commands/search.go` - Search with query and tags, styled output
- [x] `internal/cli/commands/list.go` - List user's snippets with visibility indicators
- [x] `internal/cli/commands/info.go` - Show snippet details with variables table and usage examples

#### Core Commands ✅
- [x] `internal/cli/commands/pull.go` ⭐ **Most critical - COMPLETED**:
  - Parse `namespace/slug` from args
  - **Custom variable flag parsing** (pre-Cobra for `--UPPERCASE=value`)
  - Load `--var-file` if provided
  - Call API to pull snippet
  - **Interactive prompts** for missing variables with auto-password detection
  - Display rendered content in styled box
  - Save to file with `-o` or prompt user
  - Handle warnings from template engine

- [x] `internal/cli/commands/push.go`:
  - Read template file with `-f` flag
  - Auto-extract variables
  - Display found variables with confirmation
  - Auto-generate slug if not provided
  - Create snippet via API
  - Show success with pull command example

### Phase 6: README Documentation ✅

- [x] Comprehensive quickstart guide (5 minute setup)
- [x] Detailed command reference with examples
- [x] Template syntax documentation
- [x] Configuration guide (API + CLI)
- [x] Example templates (Docker, K8s, GitHub Actions, Nginx)
- [x] Troubleshooting section
- [x] Architecture overview
- [x] Development guide with Task commands
- [x] Roadmap with completed/remaining features

---

## MVP Remaining Work (Optional Polish)

### Phase 7: Seed Data (Optional)

#### Seed Script
- [ ] `scripts/seed.go` - Create demo user + 10+ realistic snippets:
  - Docker Compose templates (postgres, redis, full-stack)
  - Kubernetes manifests (deployment, service)
  - CI/CD configs (GitHub Actions, GitLab CI)
  - Dockerfiles (multi-stage for Go, Node)
  - Config files (nginx, .env template)

**Note**: Seed data is optional for MVP. The CLI and API are fully functional without it.

### Phase 8: Testing & Polish (Optional)

#### Testing
- [ ] Basic CLI integration tests (optional)
- [ ] Manual test checklist execution

#### Minor Enhancements
- [ ] Add `.gitignore` entries for build artifacts
- [ ] Add LICENSE file
- [ ] Final bug fixes and polish

---

## Post-MVP Features (Deferred)

### Near-Term Enhancements
- [ ] Update snippet endpoint (`PUT /v1/snippets/:id`)
- [ ] Delete snippet endpoint (`DELETE /v1/snippets/:id`)
- [ ] CLI update/delete commands
- [ ] `--json` flag for scriptable output
- [ ] Shell completion generation
- [ ] Namespace resolution by slug-only (`gdsnip pull slug` → auto-resolve)

### Mid-Term Features
- [ ] **Teams**:
  - Migration: `teams`, `team_members` tables
  - API: Team CRUD, membership management
  - CLI: `gdsnip teams` commands
  - Team snippets with `team` visibility

- [ ] **Versioning**:
  - Migration: `snippet_versions` table
  - API: Version history, get specific version
  - CLI: `--version` flag for pull
  - Show version history in info

- [ ] **Web UI**:
  - React/Next.js frontend
  - Pages: Search, snippet detail, user profile, auth
  - Deploy as static site

### Long-Term Features
- [ ] **Analytics**: Track pulls, trending snippets
- [ ] **Forking**: Copy snippets to own namespace
- [ ] **Advanced Rate Limiting**: Redis-backed, per-user quotas
- [ ] **Metrics/Observability**: Prometheus, structured logging
- [ ] **CI/CD**: GitHub Actions workflows for test/build/release
- [ ] **Cloud Deployment**: Docker images, deployment guides

---

## Technical Implementation Notes

### Template Engine Pattern
- Variables MUST be `{{UPPERCASE}}` (regex: `\{\{([A-Z][A-Z0-9_]*)\}\}`)
- Validation returns ALL missing variables (not just first)
- Rendering applies defaults before provided values
- Warns about unused variables

### Service Layer Pattern
- **Handlers are thin** - they parse requests and call services
- **Business logic lives in services** - all validation, transformation, orchestration
- **Repositories are pure data access** - no business logic

### Authentication Flow
- Stateless JWT (no server-side sessions)
- CLI stores credentials in `~/.gdsnip/credentials.json` (0600)
- API validates JWT via middleware, sets `user_id` in context

### Custom Variable Flag Parsing (CLI)
The pull command accepts dynamic flags like `--DB_PASSWORD=secret`. These are:
1. Parsed from `os.Args` BEFORE Cobra processes flags
2. Extracted into `map[string]string`
3. Removed from args so Cobra doesn't reject them as unknown
4. Critical for UX - allows natural `gdsnip pull ns/slug --VAR=value` syntax

### Database Features
- **Full-text search**: `tsvector` with auto-update triggers
- **JSONB storage**: Variables stored as JSON for flexibility
- **Indexes**: Covering slug, namespace, owner, visibility, tags, search
- **Auto-timestamps**: Triggers update `updated_at` on every change

---

## Dependencies

### Core
```
github.com/labstack/echo/v4           # HTTP framework
github.com/jmoiron/sqlx               # SQL extensions
github.com/lib/pq                     # PostgreSQL driver
github.com/golang-jwt/jwt/v5          # JWT
github.com/knadh/koanf/v2             # Configuration
golang.org/x/crypto                   # bcrypt
github.com/google/uuid                # UUID generation
```

### CLI
```
github.com/spf13/cobra                # CLI framework
github.com/charmbracelet/bubbletea    # TUI framework
github.com/charmbracelet/bubbles      # TUI components
github.com/charmbracelet/lipgloss     # TUI styling
golang.org/x/term                     # Terminal control (password input)
```

### Testing
```
github.com/stretchr/testify           # Assertions
```

---

## MVP Definition of Done

The MVP is complete when:

1. ✅ **API Server**: All endpoints working, tested with curl
2. ✅ **Service Tests**: Auth + Snippet services fully tested
3. ✅ **CLI Auth**: Can register, login, logout, check status
4. ✅ **CLI Search**: Can search public snippets by query/tags
5. ✅ **CLI Pull**: Can pull snippets, prompted for missing variables, saves to file
6. ✅ **CLI Push**: Can push templates with auto-variable extraction
7. ✅ **CLI List**: Can list own snippets
8. ✅ **CLI Info**: Can view snippet details
9. ⚠️ **Seed Data**: 10+ realistic snippets available (OPTIONAL - MVP functional without it)
10. ✅ **End-to-End Test**: Can complete full user journey:
    ```bash
    gdsnip auth register          # ✅ Working
    gdsnip search docker          # ✅ Working
    gdsnip info demo/docker-pg    # ✅ Working
    gdsnip pull demo/docker-pg --DB_PASSWORD=test -o output.yml  # ✅ Working
    gdsnip push -f mytemplate.yml -n "test" --public             # ✅ Working
    gdsnip list                   # ✅ Working
    gdsnip auth logout            # ✅ Working
    ```
11. ✅ **README**: Comprehensive quickstart guide with examples

**MVP STATUS: ✅ COMPLETE** (Seed data optional for demo/testing purposes)

---

## File Checklist

### Must Have (MVP)

#### Shared Packages
- [x] `pkg/models/user.go`
- [x] `pkg/models/snippet.go`
- [x] `pkg/models/errors.go`
- [x] `pkg/template/engine.go`
- [x] `pkg/template/engine_test.go`
- [x] `pkg/validator/slug.go`
- [x] `pkg/validator/slug_test.go`

#### API
- [x] `cmd/api/main.go`
- [x] `internal/api/config/config.go`
- [x] `internal/api/repository/db.go`
- [x] `internal/api/repository/user.go`
- [x] `internal/api/repository/snippet.go`
- [x] `internal/api/service/auth.go`
- [x] `internal/api/service/snippet.go`
- [x] `internal/api/service/auth_test.go`
- [x] `internal/api/service/snippet_test.go`
- [x] `internal/api/middleware/auth.go`
- [x] `internal/api/middleware/ratelimit.go`
- [x] `internal/api/handlers/response.go`
- [x] `internal/api/handlers/auth.go`
- [x] `internal/api/handlers/snippets.go`
- [x] `internal/api/server/server.go`

#### CLI
- [x] `internal/cli/config/config.go`
- [x] `internal/cli/config/credentials.go`
- [x] `internal/cli/api/client.go`
- [x] `internal/cli/tui/styles.go`
- [x] `internal/cli/tui/prompt.go`
- [x] `internal/cli/tui/spinner.go`
- [x] `internal/cli/commands/root.go`
- [x] `internal/cli/commands/auth.go`
- [x] `internal/cli/commands/pull.go` ⭐
- [x] `internal/cli/commands/push.go`
- [x] `internal/cli/commands/search.go`
- [x] `internal/cli/commands/list.go`
- [x] `internal/cli/commands/info.go`
- [x] `cmd/cli/main.go`
- [x] `README.md` (comprehensive quickstart guide)

#### Infrastructure
- [x] `migrations/001_initial.up.sql`
- [x] `migrations/001_initial.down.sql`
- [ ] `scripts/seed.go`
- [x] `Taskfile.yml`
- [x] `docker-compose.yml`

### Post-MVP

```
internal/api/handlers/snippets.go      # Update/Delete endpoints
tests/integration/setup_test.go        # Integration test framework
tests/integration/auth_test.go
tests/integration/snippets_test.go
docs/CLI.md                            # Full command reference
docs/API.md                            # API documentation
.github/workflows/ci.yml               # CI pipeline
.github/workflows/release.yml          # Release automation
docker/Dockerfile.api                  # Docker images
docker/Dockerfile.cli
```

---

## Quick Commands Reference

```bash
# Development
task db:up                  # Start PostgreSQL
task db:migrate             # Run migrations
task db:seed                # Populate seed data (TODO)
task dev:api                # Run API server
task test                   # Run all tests

# Build
task build:cli              # Build CLI binary
task build:api              # Build API server

# Database
task db:reset               # Full reset: down → up → migrate → seed
task db:psql                # Connect to database

# Testing
task test:unit              # Unit tests only
task test:integration       # Integration tests (requires DB)
```

---

## Notes

1. **MVP First**: Focus on core pull workflow - discovery, pull with variables, push
2. **Service Layer is King**: All business logic in services, thin handlers
3. **Security**: 0600 file permissions for credentials, bcrypt cost 12, JWT validation
4. **Error Messages**: Always actionable - tell user what to do next
5. **Testing**: Co-located with implementation, not deferred
6. **Custom Flag Parsing**: Critical for pull command UX (`--VAR=value`)
7. **Seed Data**: Makes testing/demos much easier - create early
8. **Pull Command**: Most complex CLI command - interactive prompts, variable resolution
