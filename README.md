# GDSNIP - CLI-First Snippet Registry

> A developer productivity tool for storing, sharing, and pulling parameterized code templates.

GDSNIP is a CLI-first snippet registry that lets you instantly pull pre-configured templates with variable substitution. Think "npm for boilerplate code" — no more copy-pasting from Confluence or old projects.

```bash
# Search for templates
gdsnip search docker

# Pull a template with variables
gdsnip pull demo/docker-pg --DB_PASSWORD=secret --PORT=5432

# Interactive prompts for missing variables
gdsnip pull demo/k8s-deployment

# Save directly to file
gdsnip pull demo/docker-pg --DB_PASSWORD=secret -o docker-compose.yml

# Push your own templates
gdsnip push -f my-template.yml -n "My Template" --public
```

## Features

- 🎯 **Variable Substitution**: Templates with `{{UPPERCASE}}` placeholders
- 💬 **Interactive Prompts**: Auto-prompt for missing required variables
- 🔍 **Full-Text Search**: Find snippets by keywords and tags
- 🔒 **Public/Private**: Control visibility of your snippets
- 🎨 **Beautiful CLI**: Styled output with Lipgloss colors
- ⚡ **Fast**: PostgreSQL with full-text search and indexing
- 🔐 **Secure**: JWT authentication with bcrypt password hashing

## Quick Start (5 Minutes)

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) (for migrations)
- [Task](https://taskfile.dev) (optional, for convenience)

### 1. Clone and Setup

```bash
# Clone repository
git clone https://github.com/ashaju-godaddy/gdsnip.git
cd gdsnip

# Full setup (starts DB, runs migrations, seeds data, builds binaries)
task setup
```

This single command will:
- Start PostgreSQL in Docker (port 5432)
- Run database migrations
- Seed sample data
- Build both CLI and API binaries

### 2. Start the API Server

```bash
# In one terminal
task dev:api
```

API server starts on `http://localhost:8080`

### 3. Register and Login

```bash
# In another terminal
./bin/gdsnip auth register
# Enter: email, username, password

# Check authentication status
./bin/gdsnip auth status
```

Credentials are securely stored in `~/.gdsnip/credentials.json` (0600 permissions).

### 4. Create Your First Snippet

**Create a template file:**

```bash
cat > postgres-template.yml <<'EOF'
version: "3.8"
services:
  postgres:
    image: postgres:{{VERSION}}
    environment:
      POSTGRES_DB: {{DB_NAME}}
      POSTGRES_USER: {{DB_USER}}
      POSTGRES_PASSWORD: {{DB_PASSWORD}}
    ports:
      - "{{PORT}}:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
EOF
```

**Push to GDSNIP:**

```bash
./bin/gdsnip push \
  -f postgres-template.yml \
  -n "PostgreSQL Docker Compose" \
  -d "PostgreSQL setup with configurable credentials" \
  -t docker -t postgres -t database \
  --public
```

Output:
```
Creating Snippet

ℹ Found 5 variable(s):
  · VERSION
  · DB_NAME
  · DB_USER
  · DB_PASSWORD
  · PORT

ℹ Auto-generated slug: postgresql-docker-compose

Summary:
  Name: PostgreSQL Docker Compose
  Slug:  postgresql-docker-compose
  Description: PostgreSQL setup with configurable credentials
  Visibility: public
  Tags: docker, postgres, database
  Variables: 5 variable(s)

Create this snippet? [y/N] y

✓ Snippet created successfully!

  Path: yourusername/postgresql-docker-compose
  Visibility: public

ℹ To pull this snippet, run:
   gdsnip pull yourusername/postgresql-docker-compose --VERSION=value --DB_NAME=value
```

### 5. Pull and Render Snippets

**Option A: Provide all variables inline**

```bash
./bin/gdsnip pull yourusername/postgresql-docker-compose \
  --VERSION=16 \
  --DB_NAME=myapp \
  --DB_USER=admin \
  --DB_PASSWORD=secret123 \
  --PORT=5432 \
  -o docker-compose.yml
```

**Option B: Interactive prompts**

```bash
./bin/gdsnip pull yourusername/postgresql-docker-compose

# Will prompt for each variable:
VERSION *
> 16

DB_NAME *
> myapp

DB_USER *
> admin

DB_PASSWORD * (hidden input)
> ********

PORT *
> 5432

✓ Snippet rendered successfully!
Save to file? [y/N] y
Output filename  [postgresql-docker-compose.yml]
> docker-compose.yml

✓ Saved to:  docker-compose.yml
```

**Option C: Use variable file**

```bash
# Create vars.json
cat > vars.json <<'EOF'
{
  "VERSION": "16",
  "DB_NAME": "myapp",
  "DB_USER": "admin",
  "DB_PASSWORD": "secret123",
  "PORT": "5432"
}
EOF

# Pull with variable file
./bin/gdsnip pull yourusername/postgresql-docker-compose \
  --var-file vars.json \
  -o docker-compose.yml
```

### 6. Search and Discover

```bash
# Search for docker-related snippets
./bin/gdsnip search docker

# Search with multiple tags
./bin/gdsnip search postgres -t database -t docker

# List your own snippets
./bin/gdsnip list

# View detailed snippet info
./bin/gdsnip info yourusername/postgresql-docker-compose
```

## CLI Commands Reference

### Authentication

```bash
gdsnip auth register          # Create new account
gdsnip auth login             # Login with email/password
gdsnip auth logout            # Clear stored credentials
gdsnip auth status            # Show current user info
```

### Snippet Management

```bash
# Push (create new snippet)
gdsnip push -f <file> -n "<name>" [options]
  -f, --file         Template file (required)
  -n, --name         Snippet name (required)
  -s, --slug         Custom slug (auto-generated if omitted)
  -d, --description  Description
  -t, --tags         Tags (repeatable: -t docker -t postgres)
  --public           Make public (default: private)

# Pull (render snippet)
gdsnip pull <namespace>/<slug> [options]
  --VARNAME=value    Provide variable values
  --var-file <file>  Load variables from JSON file
  -o, --output       Output file (prompts if omitted)

# Search
gdsnip search [query] [options]
  -t, --tags         Filter by tags (repeatable)
  --limit            Max results (default: 20, max: 100)

# List your snippets
gdsnip list [options]
  --limit            Max results (default: 20)
  --offset           Pagination offset

# View snippet details
gdsnip info <namespace>/<slug>
```

### Utility

```bash
gdsnip version               # Show version
gdsnip help                  # Show help
gdsnip <command> --help      # Command-specific help
```

## Template Variable Syntax

Variables **must** be `{{UPPERCASE}}` format:

```yaml
# ✓ Valid
DATABASE_URL: {{DATABASE_URL}}
API_KEY: {{API_KEY}}
PORT: {{PORT}}
SERVICE_NAME: {{SERVICE_NAME}}

# ✗ Invalid
database_url: {{database_url}}   # lowercase not allowed
api-key: {{api-key}}             # hyphens not allowed
serviceName: {{serviceName}}     # camelCase not allowed
```

### Password Variables

Variables containing these keywords are automatically treated as passwords with hidden input:
- `PASSWORD`
- `SECRET`
- `TOKEN`
- `KEY`

Examples: `DB_PASSWORD`, `API_SECRET`, `JWT_TOKEN`, `SIGNING_KEY`

## Example Templates

### Docker Compose - Redis

```yaml
version: "3.8"
services:
  redis:
    image: redis:{{VERSION}}
    command: redis-server --requirepass {{PASSWORD}}
    ports:
      - "{{PORT}}:6379"
    volumes:
      - redis_data:/data

volumes:
  redis_data:
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{APP_NAME}}
  namespace: {{NAMESPACE}}
spec:
  replicas: {{REPLICAS}}
  selector:
    matchLabels:
      app: {{APP_NAME}}
  template:
    metadata:
      labels:
        app: {{APP_NAME}}
    spec:
      containers:
      - name: {{APP_NAME}}
        image: {{IMAGE}}:{{TAG}}
        ports:
        - containerPort: {{PORT}}
        env:
        - name: ENV
          value: {{ENVIRONMENT}}
```

### GitHub Actions Workflow

```yaml
name: {{WORKFLOW_NAME}}

on:
  push:
    branches: [ {{BRANCH}} ]
  pull_request:
    branches: [ {{BRANCH}} ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '{{GO_VERSION}}'

      - name: Run tests
        run: go test ./...
```

### Nginx Configuration

```nginx
server {
    listen {{PORT}};
    server_name {{DOMAIN}};

    location / {
        proxy_pass http://{{UPSTREAM_HOST}}:{{UPSTREAM_PORT}};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Configuration

### API Server

Configure via environment variables or `config.yml`:

```bash
# Required
export GDSNIP_DATABASE_URL="postgres://gdsnip:gdsnip@localhost:5432/gdsnip?sslmode=disable"
export GDSNIP_JWT_SECRET="your-secret-key-minimum-32-characters-long"

# Optional
export GDSNIP_PORT=8080                    # Default: 8080
export GDSNIP_JWT_EXPIRY=168h              # Default: 168h (7 days)
export GDSNIP_RATE_LIMIT=100               # Default: 100 req/min
```

Or use `config.yml`:

```yaml
port: 8080
database_url: postgres://gdsnip:gdsnip@localhost:5432/gdsnip?sslmode=disable
jwt_secret: your-secret-key-minimum-32-characters-long
jwt_expiry: 168h
rate_limit: 100
```

### CLI

Configure via environment or `~/.gdsnip/config.yml`:

```bash
export GDSNIP_API_URL="http://localhost:8080/v1"
```

Or create `~/.gdsnip/config.yml`:

```yaml
api_url: http://localhost:8080/v1
```

Credentials are stored in `~/.gdsnip/credentials.json` (0600 permissions).

## Development

### Project Structure

```
gdsnip/
├── cmd/
│   ├── api/                 # API server entry point
│   └── cli/                 # CLI entry point
├── internal/
│   ├── api/                 # API implementation
│   │   ├── config/          # Config loading (Koanf)
│   │   ├── handlers/        # HTTP handlers (thin, delegate to services)
│   │   ├── middleware/      # JWT auth, rate limiting
│   │   ├── repository/      # Database access (sqlx)
│   │   ├── server/          # Server setup & routing
│   │   └── service/         # Business logic
│   └── cli/                 # CLI implementation
│       ├── api/             # API client
│       ├── commands/        # Cobra commands
│       ├── config/          # Config & credentials
│       └── tui/             # Terminal UI (Lipgloss)
├── pkg/                     # Shared packages
│   ├── models/              # Data models (User, Snippet, APIError)
│   ├── template/            # Template engine (extract, validate, render)
│   └── validator/           # Input validation (slugs, emails, etc.)
├── migrations/              # Database migrations (golang-migrate)
├── scripts/                 # Utility scripts (seed data)
└── tests/                   # Integration tests
```

### Task Commands

```bash
# Development
task dev:api              # Run API with hot reload (requires air)
task dev:cli              # Build CLI binary to ./bin/gdsnip

# Build
task build                # Build both CLI and API
task build:cli            # Build CLI with version injection
task build:api            # Build API server

# Database
task db:up                # Start PostgreSQL container
task db:down              # Stop containers
task db:migrate           # Run migrations up
task db:migrate:down      # Rollback last migration
task db:seed              # Populate seed data
task db:reset             # Full reset (down + up + migrate + seed)
task db:psql              # Connect with psql

# Testing
task test                 # Run all tests
task test:unit            # Unit tests only (no DB)
task test:integration     # Integration tests (requires DB)
task test:cover           # Generate coverage report

# Code Quality
task fmt                  # Format code (go fmt)
task vet                  # Run go vet
task lint                 # Run golangci-lint
task tidy                 # Tidy go modules

# Utilities
task clean                # Clean build artifacts
task install              # Install CLI to /usr/local/bin
task uninstall            # Uninstall CLI
```

### Running Tests

```bash
# All tests (unit + service + integration)
task test

# Specific packages
go test ./pkg/template/... -v              # Template engine tests
go test ./pkg/validator/... -v             # Validator tests
go test ./internal/api/service/... -v      # Service tests

# With coverage
task test:cover
open coverage.html
```

**Current Test Coverage:**
- ✅ 24+ template engine tests
- ✅ 50+ validator tests
- ✅ 14 AuthService tests
- ✅ 18 SnippetService tests
- **Total: 100+ tests passing**

## Architecture

### Service Layer Pattern

GDSNIP follows a clean three-layer architecture:

1. **Handlers** (HTTP layer): Parse requests, call services, return responses
2. **Services** (Business logic): Validation, orchestration, business rules
3. **Repositories** (Data access): Database queries, no business logic

This keeps code testable and maintainable.

### Template Engine

- **Extraction**: Regex `\{\{([A-Z][A-Z0-9_]*)\}\}` finds all variables
- **Validation**: Returns ALL missing required variables (not just first)
- **Rendering**: Applies defaults before provided values
- **Warnings**: Notifies about unused provided variables

### Authentication

- **Stateless JWT** tokens (HS256 signing)
- **Bcrypt** password hashing (cost 12)
- Tokens expire after 168 hours (7 days) by default
- CLI stores credentials in `~/.gdsnip/credentials.json` with 0600 permissions

### Search

- **PostgreSQL full-text search** using `tsvector`
- Indexed on name, description, and tags
- Auto-updated via database triggers

## API Endpoints

```
Health:
  GET  /health

Authentication:
  POST /v1/auth/register       # Create account
  POST /v1/auth/login          # Get JWT token
  GET  /v1/auth/me             # Get current user (protected)

Snippets:
  POST /v1/snippets                        # Create snippet (protected)
  GET  /v1/snippets/:namespace/:slug       # Get snippet details (public)
  POST /v1/snippets/:namespace/:slug/pull  # Render template (public)
  GET  /v1/snippets?q=query&tags=tag       # Search public snippets
  GET  /v1/users/me/snippets               # List my snippets (protected)
```

All responses use this format:

```json
{
  "success": true,
  "data": { ... },
  "pagination": { "total": 10, "limit": 20, "offset": 0 }
}
```

Errors:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Descriptive error message",
    "details": { ... }
  }
}
```

## Troubleshooting

### API won't start

```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check database connection
task db:psql

# View API logs
task dev:api
```

### CLI can't connect to API

```bash
# Check API URL
echo $GDSNIP_API_URL

# Test API health
curl http://localhost:8080/health

# Check credentials
cat ~/.gdsnip/credentials.json
```

### Authentication errors

```bash
# Clear credentials and re-login
gdsnip auth logout
gdsnip auth login
```

### Database migration errors

```bash
# Check current version
migrate -path ./migrations -database "$GDSNIP_DATABASE_URL" version

# Force version (if stuck)
migrate -path ./migrations -database "$GDSNIP_DATABASE_URL" force <version>

# Reset database
task db:reset
```

## Roadmap

### ✅ MVP (Completed)
- ✅ Authentication (register, login, logout, status)
- ✅ Push snippets with auto-variable extraction
- ✅ Pull snippets with interactive prompts
- ✅ Search by query and tags
- ✅ List and view snippet details
- ✅ Public/private visibility
- ✅ Secure credential storage
- ✅ Styled CLI with Lipgloss

### 📋 Near-Term
- [ ] Update/delete snippets
- [ ] Variable file validation
- [ ] Shell completion (bash, zsh, fish)
- [ ] `--json` output for scripting
- [ ] Config command for settings management

### 🔮 Mid-Term
- [ ] Teams and team snippets
- [ ] Snippet versioning
- [ ] Fork snippets
- [ ] Namespace resolution by slug only
- [ ] Variable defaults in templates: `{{VAR:default}}`
- [ ] Snippet categories/collections

### 🚀 Long-Term
- [ ] Web UI (search, view, manage snippets)
- [ ] VS Code extension
- [ ] Analytics and trending
- [ ] Binary distribution (Homebrew, apt, snap)
- [ ] Self-hosted deployment guides
- [ ] Cloud-hosted service




