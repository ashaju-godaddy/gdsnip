# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GDSNIP is a CLI-first snippet registry for storing, sharing, and pulling parameterized code templates. It's essentially "npm for boilerplate code" with variable substitution.

## Architecture

### Three-Layer Structure

1. **Shared Packages** (`pkg/`)
   - `pkg/template`: Template engine (extraction, validation, rendering)
   - `pkg/models`: Data models shared between CLI and API
   - `pkg/validator`: Input validation (slugs, emails, etc.)

2. **API Server** (`internal/api/`)
   - **Repository Layer**: Direct database access (sqlx + PostgreSQL)
   - **Service Layer**: Business logic (auth, snippet operations)
   - **Handler Layer**: Thin HTTP handlers (Echo framework)
   - Middleware: JWT auth, rate limiting

3. **CLI** (`internal/cli/`)
   - Commands: Cobra-based (auth, pull, push, search, list, info)
   - TUI: Bubbletea/Lipgloss for interactive prompts
   - API Client: HTTP client for backend communication

### Key Design Decisions

**Template Engine Pattern**:
- Variables MUST be `{{UPPERCASE}}` (regex: `\{\{([A-Z][A-Z0-9_]*)\}\}`)
- Validation returns ALL missing required variables (not just first)
- Rendering applies defaults before provided values
- Warns about unused variables

**Service Layer Pattern**:
- Handlers are thin wrappers that delegate to services
- Services contain ALL business logic
- Repositories are pure data access (no business logic)
- This keeps handlers testable and logic centralized

**Authentication Flow**:
- Stateless JWT (no server-side session storage)
- CLI stores credentials in `~/.gdsnip/credentials.json` (mode 0600)
- API validates JWT via middleware, sets `user_id` in context

## Essential Commands

### Development Workflow
```bash
# First-time setup
task setup                  # Start DB, run migrations, seed data

# Development
task dev:api               # Run API with hot reload (requires air)
task dev:cli               # Build CLI binary to ./bin/gdsnip

# Database
task db:up                 # Start PostgreSQL container
task db:migrate            # Run migrations up
task db:migrate:down       # Rollback last migration
task db:seed               # Populate seed data
task db:reset              # Full reset: migrate down, up, seed
task db:psql               # Connect to database with psql

# Testing
task test                  # Run all tests
task test:unit             # Unit tests only (no DB required)
task test:integration      # Integration tests (requires DB)
task test:cover            # Generate coverage report

# Single package test
go test ./pkg/template/... -v
go test ./internal/api/service/... -v -run TestAuthService

# Build
task build                 # Build both CLI and API
task build:cli             # Build CLI with version injection
task build:api             # Build API server
```

### Database Migrations

Migrations use `golang-migrate`. Pattern:
- `migrations/001_initial.up.sql` - Apply changes
- `migrations/001_initial.down.sql` - Rollback changes

The database includes:
- **Full-text search** via PostgreSQL `tsvector` (indexed on name, description, tags)
- **Auto-updated timestamps** via triggers
- **JSONB storage** for variable metadata

## Testing Philosophy

**Co-located Tests**: Tests live in the same phase as implementation (not deferred).

**Test Structure**:
- `pkg/template/engine_test.go` - Comprehensive template engine tests (24+ tests)
- `pkg/validator/slug_test.go` - Validation edge cases (50+ tests)
- API tests: Service layer tests (business logic) + handler tests (HTTP)
- CLI tests: Mock API responses with `httptest.NewServer`

**Integration Tests**: Located in `tests/integration/`, require running PostgreSQL.

## Important Patterns

### Variable Flag Parsing (CLI)
The pull command accepts dynamic flags like `--DB_PASSWORD=secret`. These are parsed from `os.Args` BEFORE Cobra, since Cobra would reject unknown flags. See `internal/cli/commands/pull.go`.

### Slug Generation
Slugs are auto-generated from names:
- `GenerateSlug("Hello World")` → `"hello-world"`
- Validation enforces: lowercase, alphanumeric, hyphens only
- Must start/end with alphanumeric
- No consecutive hyphens (`--`)

### Error Handling (API)
All errors use structured `APIError`:
```go
models.NewMissingVariablesError("Missing vars", map[string]interface{}{
    "missing": []Variable{...},
})
```
The `HTTPStatus()` method maps error codes to status codes.

### Namespace Resolution (Pull)
- `gdsnip pull namespace/slug` → Direct lookup
- `gdsnip pull slug` → Search public snippets by slug
  - Returns snippet if exactly one match
  - Errors if multiple matches (ambiguous)
  - Suggests using full path

## Configuration

### API Config (Koanf)
Sources (priority order):
1. Environment variables (`GDSNIP_*`)
2. `config.yml` file
3. Defaults

Required:
- `GDSNIP_DATABASE_URL`: PostgreSQL connection string
- `GDSNIP_JWT_SECRET`: Signing secret (min 32 chars in production)

### CLI Config
- API URL: `GDSNIP_API_URL` (default: `http://localhost:8080/v1`)
- Credentials stored: `~/.gdsnip/credentials.json`

## Project Status

**Completed**:
- ✅ Phase 1: Project setup, Docker, migrations
- ✅ Phase 2: Shared packages (template engine, models, validators)

**In Progress**:
- 🚧 Phase 3: API server (config, repos, services, handlers)

**Pending**:
- Phase 4: CLI foundation (config, API client, TUI)
- Phase 5: CLI commands (auth, pull, push, search)

## Key Files Reference

- `PLAN.MD`: Detailed implementation plan with acceptance criteria
- `PROJECT_DETAILS.MD`: Full product specification with API docs
- `Taskfile.yml`: All development tasks
- `migrations/001_initial.up.sql`: Database schema
- `pkg/template/engine.go`: Core template rendering logic

## Gotchas

1. **Template variables**: Must be UPPERCASE. `{{lowercase}}` will NOT match.
2. **Password validation**: Minimum 8 chars (see `pkg/validator/slug.go`)
3. **Migrations**: Always create both `.up.sql` and `.down.sql`
4. **JWT expiry**: Default is 168h (7 days), configurable via `GDSNIP_JWT_EXPIRY`
5. **Rate limiting**: Default 100 req/min per IP (in-memory, resets on restart)
