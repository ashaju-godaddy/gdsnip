# GDSNIP - Project Details

## Overview

**GDSNIP** is a developer productivity tool that provides a CLI-first snippet registry for storing, sharing, and pulling parameterized code/config templates. Think "npm for boilerplate code" — developers can instantly pull pre-configured templates with variable substitution instead of copy-pasting from Confluence, Slack, or old projects.

---

## Problem Statement

### Current Pain Points

1. **Scattered Boilerplate**: Developers copy-paste from Confluence, Slack messages, old repos, or Google searches
2. **Context Switching**: Leaving IDE to find that Docker Compose template breaks flow
3. **Stale Templates**: That Confluence page hasn't been updated in 2 years
4. **No Parameterization**: Copy-paste requires manual find-replace for project-specific values
5. **Tribal Knowledge**: Senior devs have templates bookmarked; new hires struggle
6. **Team Silos**: Each team has their own templates with no discoverability

### User Stories

> "As a developer, I want to quickly pull a Docker Postgres setup without leaving my terminal, so I can focus on building features instead of writing boilerplate."

> "As a platform team member, I want to publish standardized service templates that other teams can use, so we maintain consistency across the organization."

> "As a new hire, I want to discover what templates exist and how to use them, so I can onboard faster."

---

## Solution

A CLI tool + API that enables:

```bash
# Search for templates
gdsnip search kafka

# Pull a template with variables
gdsnip pull payments/kafka-setup --BROKER_URL=localhost:9092 --TOPIC=orders

# Output is rendered and ready to use
# Save directly to file
gdsnip pull docker-pg --DB_PASSWORD=secret -o docker-compose.yml

# Push your own templates
gdsnip push -f my-template.yml --name my-service --public
```

---

## Product Scope

### MVP Features

| Feature | Description | Priority |
|---------|-------------|----------|
| User Authentication | Register, login, JWT-based auth | P0 |
| Pull Snippets | Fetch and render templates with variables | P0 |
| Push Snippets | Create new snippets via CLI | P0 |
| Search Snippets | Find snippets by keyword | P0 |
| Variable Substitution | Replace `{{VAR}}` placeholders with provided values | P0 |
| Public/Private Visibility | Control who can see your snippets | P0 |
| Namespaces | User namespaces (e.g., `johndoe/docker-pg`) | P0 |

### Stretch Features

| Feature | Description | Priority |
|---------|-------------|----------|
| Teams | Team namespaces and shared snippets | P1 |
| Versioning | Snippet version history | P1 |
| Tags | Categorize snippets with tags | P1 |
| Web UI | Browse snippets in browser | P2 |
| Forking | Fork and modify others' public snippets | P2 |
| Analytics | Track pull counts | P2 |
| VS Code Extension | Pull snippets from IDE | P3 |

---

## User Flows

### Flow 1: First-Time User Registration

```
1. User runs: gdsnip auth register
2. CLI presents interactive form (bubbletea):
   - Email
   - Username
   - Password
3. CLI calls POST /v1/auth/register
4. API creates user, returns JWT
5. CLI stores token in ~/.gdsnip/credentials.json
6. CLI displays: "✔ Welcome, {username}! Run 'gdsnip search' to explore snippets."
```

### Flow 2: Login (Returning User)

```
1. User runs: gdsnip auth login
2. CLI presents interactive form:
   - Email
   - Password
3. CLI calls POST /v1/auth/login
4. API validates credentials, returns JWT
5. CLI stores token in ~/.gdsnip/credentials.json
6. CLI displays: "✔ Logged in as {username}"
```

### Flow 3: Search Snippets

```
1. User runs: gdsnip search docker
2. CLI calls GET /v1/snippets?q=docker
3. API returns matching public snippets
4. CLI displays formatted table:

   NAME                 NAMESPACE    DESCRIPTION                  PULLS
   docker-pg            public       PostgreSQL with Docker       234
   docker-redis         public       Redis with persistence       189
   docker-kafka         payments     Kafka + Zookeeper setup      45

5. User sees: "Use: gdsnip pull <namespace>/<name>"
```

### Flow 4: Pull Snippet (Core Feature)

```
1. User runs: gdsnip pull johndoe/docker-pg --DB_PASSWORD=secret
2. CLI parses snippet path and variables
3. CLI calls POST /v1/snippets/johndoe/docker-pg/pull
   Body: { "variables": { "DB_PASSWORD": "secret" } }
4. API:
   a. Fetches snippet from database
   b. Validates required variables are provided
   c. Renders template with variable substitution
   d. Returns rendered content
5. CLI displays:
   
   ✔ Found: johndoe/docker-pg (v1)
   
   ──────────────────────────────────────────
   version: '3.8'
   services:
     postgres:
       image: postgres:15
       environment:
         POSTGRES_PASSWORD: secret
   ...
   ──────────────────────────────────────────
   
6. If -o flag provided, saves to file:
   ✔ Saved to docker-compose.yml
```

### Flow 5: Pull with Missing Variables

```
1. User runs: gdsnip pull docker-pg (missing required variable)
2. CLI calls API
3. API returns error: "Missing required variables: DB_PASSWORD"
4. CLI displays:
   
   ✗ Missing required variables:
   
     DB_PASSWORD (required) - Database password
   
   Usage: gdsnip pull docker-pg --DB_PASSWORD=<value>

5. Alternatively, CLI prompts interactively (bubbletea):
   
   ? DB_PASSWORD (required): ********
```

### Flow 6: Push New Snippet

```
1. User creates template file: my-service.yml
   
   apiVersion: v1
   kind: Service
   metadata:
     name: {{SERVICE_NAME}}
   spec:
     port: {{PORT}}

2. User runs: gdsnip push -f my-service.yml --name k8s-service --description "K8s service template" --public

3. CLI:
   a. Reads file content
   b. Extracts variables: SERVICE_NAME, PORT
   c. Prompts for variable metadata (bubbletea):
      
      Found variables: SERVICE_NAME, PORT
      ? Add descriptions? (y/N)

4. CLI calls POST /v1/snippets
   Body: {
     "name": "k8s-service",
     "slug": "k8s-service",
     "content": "apiVersion: v1...",
     "variables": [
       {"name": "SERVICE_NAME", "required": true},
       {"name": "PORT", "required": true, "default": "80"}
     ],
     "visibility": "public"
   }

5. API creates snippet, returns snippet details
6. CLI displays:
   
   ✔ Created: johndoe/k8s-service
   
   Pull command:
     gdsnip pull johndoe/k8s-service --SERVICE_NAME=<value> --PORT=<value>
```

### Flow 7: View Snippet Info

```
1. User runs: gdsnip info payments/kafka-setup
2. CLI calls GET /v1/snippets/payments/kafka-setup
3. CLI displays (styled with lipgloss):

   ┌─────────────────────────────────────────────────────────────┐
   │  payments/kafka-setup                                       │
   ├─────────────────────────────────────────────────────────────┤
   │  Description:  Kafka configuration for payments team        │
   │  Visibility:   team                                         │
   │  Version:      3                                            │
   │  Created:      2024-01-10 by @alice                         │
   │  Pulls:        45                                           │
   ├─────────────────────────────────────────────────────────────┤
   │  Variables:                                                 │
   │    BROKER_URL     (required)   Kafka broker URL             │
   │    TOPIC_PREFIX   (optional)   Default: "payments"          │
   │    REPLICAS       (optional)   Default: 3                   │
   ├─────────────────────────────────────────────────────────────┤
   │  Tags:  kafka, docker, messaging                            │
   └─────────────────────────────────────────────────────────────┘

   Usage:
     gdsnip pull payments/kafka-setup --BROKER_URL=<value>
```

### Flow 8: List My Snippets

```
1. User runs: gdsnip list
2. CLI calls GET /v1/users/me/snippets
3. CLI displays:

   Your snippets:

   NAME           VISIBILITY   PULLS   UPDATED
   docker-pg      public       234     2 days ago
   my-service     private      12      1 week ago
   redis-cache    public       89      3 weeks ago

   Total: 3 snippets
```

---

## Technical Specifications

### Tech Stack

| Component | Technology | Justification |
|-----------|------------|---------------|
| CLI | Go + Cobra + Bubbletea | Single binary, interactive TUI, great UX |
| API | Go + Echo | Lightweight, performant, shared language with CLI |
| Database | PostgreSQL | Reliable, JSONB for variables, full-text search |
| Config | Koanf | Flexible, supports multiple sources |
| Auth | JWT | Stateless, simple |
| Template Engine | Custom Go | Simple `{{VAR}}` substitution |
| Task Runner | Task | Simple, cross-platform alternative to Make |

### System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                           CLIENTS                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌─────────────────┐                  ┌─────────────────┐      │
│   │   GDSNIP CLI    │                  │   Web UI        │      │
│   │ (Cobra+Bubble)  │                  │   (Optional)    │      │
│   │                 │                  │                 │      │
│   └────────┬────────┘                  └────────┬────────┘      │
│            │                                    │               │
└────────────┼────────────────────────────────────┼───────────────┘
             │              HTTPS                 │
             └─────────────────┬──────────────────┘
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────┐
│                         API SERVER                               │
│                        (Go + Echo)                               │
│                                                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │   Auth   │  │ Snippets │  │  Search  │  │  Users   │        │
│  │ Handler  │  │ Handler  │  │ Handler  │  │ Handler  │        │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘        │
│       │             │             │             │               │
│       └─────────────┴─────────────┴─────────────┘               │
│                             │                                    │
│                    ┌────────┴────────┐                          │
│                    │ Template Engine │                          │
│                    │   (pkg/template)│                          │
│                    └────────┬────────┘                          │
│                             │                                    │
└─────────────────────────────┼───────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        DATA LAYER                                │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                     PostgreSQL                           │   │
│   │                                                          │   │
│   │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌───────────┐  │   │
│   │  │  users  │  │snippets │  │  teams  │  │team_members│  │   │
│   │  └─────────┘  └─────────┘  └─────────┘  └───────────┘  │   │
│   │                                                          │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Data Models

### User

```go
type User struct {
    ID           string    `json:"id"`
    Email        string    `json:"email"`
    Username     string    `json:"username"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### Snippet

```go
type Snippet struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`
    Slug        string     `json:"slug"`
    Description string     `json:"description"`
    Content     string     `json:"content"`
    Variables   []Variable `json:"variables"`
    Tags        []string   `json:"tags"`
    Visibility  string     `json:"visibility"`  // public, private, team
    OwnerType   string     `json:"owner_type"`  // user, team
    OwnerID     string     `json:"owner_id"`
    Namespace   string     `json:"namespace"`   // Computed: username or team slug
    Version     int        `json:"version"`
    CreatedBy   string     `json:"created_by"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
```

### Variable

```go
type Variable struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    Required    bool   `json:"required"`
    Default     string `json:"default,omitempty"`
    Sensitive   bool   `json:"sensitive,omitempty"`
}
```

### Team (Stretch)

```go
type Team struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Slug        string    `json:"slug"`
    Description string    `json:"description"`
    CreatedBy   string    `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
}
```

---

## Database Schema

```sql
-- Users
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    username        VARCHAR(50) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- Snippets
CREATE TABLE snippets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    slug            VARCHAR(100) NOT NULL,
    description     TEXT,
    content         TEXT NOT NULL,
    variables       JSONB DEFAULT '[]',
    tags            TEXT[] DEFAULT '{}',
    visibility      VARCHAR(20) DEFAULT 'private',
    owner_type      VARCHAR(20) NOT NULL,
    owner_id        UUID NOT NULL,
    created_by      UUID REFERENCES users(id),
    version         INTEGER DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(owner_type, owner_id, slug)
);

-- Teams (Stretch)
CREATE TABLE teams (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) UNIQUE NOT NULL,
    slug            VARCHAR(100) UNIQUE NOT NULL,
    description     TEXT,
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Team Members (Stretch)
CREATE TABLE team_members (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id         UUID REFERENCES teams(id) ON DELETE CASCADE,
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(20) DEFAULT 'member',
    joined_at       TIMESTAMP DEFAULT NOW(),
    UNIQUE(team_id, user_id)
);

-- Pull Analytics (Stretch)
CREATE TABLE snippet_pulls (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snippet_id      UUID REFERENCES snippets(id),
    user_id         UUID REFERENCES users(id),
    pulled_at       TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_snippets_slug ON snippets(slug);
CREATE INDEX idx_snippets_owner ON snippets(owner_type, owner_id);
CREATE INDEX idx_snippets_visibility ON snippets(visibility);
CREATE INDEX idx_snippets_tags ON snippets USING GIN(tags);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

---

## API Specification

### Base URL

```
Production: https://api.gdsnip.dev/v1
Development: http://localhost:8080/v1
```

### Authentication

All protected endpoints require:
```
Authorization: Bearer <jwt_token>
```

### Endpoints

#### Auth

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/auth/register` | Create new account | No |
| POST | `/auth/login` | Login, returns JWT | No |
| GET | `/auth/me` | Get current user | Yes |
| POST | `/auth/logout` | Invalidate token | Yes |

#### Snippets

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/snippets` | Search/list snippets | No |
| GET | `/snippets/:namespace/:slug` | Get snippet details | No* |
| POST | `/snippets/:namespace/:slug/pull` | Pull rendered snippet | No* |
| POST | `/snippets` | Create new snippet | Yes |
| PUT | `/snippets/:id` | Update snippet | Yes |
| DELETE | `/snippets/:id` | Delete snippet | Yes |

*Auth required for private/team snippets

#### Users

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/users/:username` | Get user public profile | No |
| GET | `/users/:username/snippets` | Get user's public snippets | No |
| GET | `/users/me/snippets` | Get my snippets | Yes |

#### Teams (Stretch)

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/teams` | List my teams | Yes |
| POST | `/teams` | Create team | Yes |
| GET | `/teams/:slug` | Get team details | Yes |
| GET | `/teams/:slug/snippets` | Get team snippets | Yes |
| POST | `/teams/:slug/members` | Add team member | Yes |

---

### API Request/Response Examples

#### POST /v1/auth/register

**Request:**
```json
{
  "email": "john@example.com",
  "username": "johndoe",
  "password": "securepassword123"
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john@example.com",
      "username": "johndoe",
      "created_at": "2024-01-15T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**Error Response (400):**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Username already taken"
  }
}
```

#### POST /v1/auth/login

**Request:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "johndoe"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### GET /v1/snippets?q=docker&tags=postgres

**Response (200):**
```json
{
  "success": true,
  "data": {
    "snippets": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "name": "Docker PostgreSQL",
        "slug": "docker-pg",
        "namespace": "johndoe",
        "full_path": "johndoe/docker-pg",
        "description": "Docker Compose setup for PostgreSQL",
        "tags": ["docker", "postgres", "database"],
        "visibility": "public",
        "version": 2,
        "variables": [
          {"name": "DB_PASSWORD", "required": true, "sensitive": true},
          {"name": "DB_NAME", "required": false, "default": "mydb"},
          {"name": "PORT", "required": false, "default": "5432"}
        ],
        "created_at": "2024-01-10T08:00:00Z",
        "updated_at": "2024-01-14T12:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 20
  }
}
```

#### GET /v1/snippets/johndoe/docker-pg

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Docker PostgreSQL",
    "slug": "docker-pg",
    "namespace": "johndoe",
    "full_path": "johndoe/docker-pg",
    "description": "Docker Compose setup for PostgreSQL with persistent volume",
    "content": "version: '3.8'\nservices:\n  postgres:\n    image: postgres:{{PG_VERSION}}\n    environment:\n      POSTGRES_DB: {{DB_NAME}}\n      POSTGRES_PASSWORD: {{DB_PASSWORD}}\n    ports:\n      - '{{PORT}}:5432'\n    volumes:\n      - pgdata:/var/lib/postgresql/data\n\nvolumes:\n  pgdata:",
    "variables": [
      {"name": "PG_VERSION", "description": "PostgreSQL version", "required": false, "default": "15"},
      {"name": "DB_NAME", "description": "Database name", "required": false, "default": "mydb"},
      {"name": "DB_PASSWORD", "description": "Database password", "required": true, "sensitive": true},
      {"name": "PORT", "description": "Host port", "required": false, "default": "5432"}
    ],
    "tags": ["docker", "postgres", "database"],
    "visibility": "public",
    "version": 2,
    "owner": {
      "type": "user",
      "username": "johndoe"
    },
    "pull_command": "gdsnip pull johndoe/docker-pg --DB_PASSWORD=<value>",
    "created_at": "2024-01-10T08:00:00Z",
    "updated_at": "2024-01-14T12:00:00Z"
  }
}
```

#### POST /v1/snippets/johndoe/docker-pg/pull

**Request:**
```json
{
  "variables": {
    "DB_PASSWORD": "supersecret",
    "DB_NAME": "production",
    "PORT": "5433"
  }
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "snippet": "johndoe/docker-pg",
    "version": 2,
    "rendered": "version: '3.8'\nservices:\n  postgres:\n    image: postgres:15\n    environment:\n      POSTGRES_DB: production\n      POSTGRES_PASSWORD: supersecret\n    ports:\n      - '5433:5432'\n    volumes:\n      - pgdata:/var/lib/postgresql/data\n\nvolumes:\n  pgdata:",
    "filename_suggestion": "docker-compose.yml"
  }
}
```

**Error Response (400) - Missing Variables:**
```json
{
  "success": false,
  "error": {
    "code": "MISSING_VARIABLES",
    "message": "Missing required variables",
    "details": {
      "missing": [
        {"name": "DB_PASSWORD", "description": "Database password"}
      ]
    }
  }
}
```

#### POST /v1/snippets

**Request:**
```json
{
  "name": "Kubernetes Service",
  "slug": "k8s-service",
  "description": "Basic Kubernetes service manifest",
  "content": "apiVersion: v1\nkind: Service\nmetadata:\n  name: {{SERVICE_NAME}}\n  namespace: {{NAMESPACE}}\nspec:\n  selector:\n    app: {{APP_NAME}}\n  ports:\n    - port: {{PORT}}\n      targetPort: {{TARGET_PORT}}",
  "variables": [
    {"name": "SERVICE_NAME", "description": "Name of the service", "required": true},
    {"name": "NAMESPACE", "description": "Kubernetes namespace", "required": false, "default": "default"},
    {"name": "APP_NAME", "description": "App selector label", "required": true},
    {"name": "PORT", "description": "Service port", "required": false, "default": "80"},
    {"name": "TARGET_PORT", "description": "Container port", "required": false, "default": "8080"}
  ],
  "tags": ["kubernetes", "k8s", "service"],
  "visibility": "public"
}
```

**Response (201):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "name": "Kubernetes Service",
    "slug": "k8s-service",
    "namespace": "johndoe",
    "full_path": "johndoe/k8s-service",
    "pull_command": "gdsnip pull johndoe/k8s-service --SERVICE_NAME=<value> --APP_NAME=<value>",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

---

## CLI Specification

### Installation

```bash
# macOS/Linux
curl -fsSL https://gdsnip.dev/install.sh | bash

# Or with Go
go install github.com/yourorg/gdsnip/cmd/cli@latest

# Verify installation
gdsnip --version
```

### Command Reference

```
gdsnip - A CLI snippet registry for developers

Usage:
  gdsnip [command]

Available Commands:
  auth        Authentication commands
  pull        Pull and render a snippet
  push        Create or update a snippet
  search      Search for snippets
  list        List your snippets
  info        Get detailed snippet information
  help        Help about any command

Flags:
  -h, --help       help for gdsnip
  -v, --version    version for gdsnip
      --json       Output in JSON format
      --no-color   Disable colored output

Use "gdsnip [command] --help" for more information about a command.
```

### Command Details

#### gdsnip auth

```
Manage authentication

Usage:
  gdsnip auth [command]

Available Commands:
  register    Create a new account
  login       Login to your account
  logout      Logout from your account
  status      Show current authentication status

Examples:
  gdsnip auth register
  gdsnip auth login
  gdsnip auth status
```

#### gdsnip pull

```
Pull and render a snippet with variable substitution

Usage:
  gdsnip pull <namespace/slug> [flags]

Flags:
  -o, --output string     Output file path
      --var-file string   JSON file containing variables
      --version int       Specific version to pull (default: latest)
      --VAR_NAME string   Set variable value (can be repeated)

Examples:
  # Pull a public snippet
  gdsnip pull docker-pg --DB_PASSWORD=secret

  # Pull from a specific namespace
  gdsnip pull johndoe/docker-pg --DB_PASSWORD=secret

  # Save to file
  gdsnip pull docker-pg --DB_PASSWORD=secret -o docker-compose.yml

  # Use variable file
  gdsnip pull docker-pg --var-file=vars.json

  # Pull specific version
  gdsnip pull docker-pg --DB_PASSWORD=secret --version=2
```

#### gdsnip push

```
Create or update a snippet

Usage:
  gdsnip push [flags]

Flags:
  -f, --file string          Template file path (required)
  -n, --name string          Snippet name
  -s, --slug string          URL-friendly slug (auto-generated from name if not provided)
  -d, --description string   Snippet description
  -t, --tags strings         Tags (comma-separated)
      --public               Make snippet public
      --team string          Push to team namespace

Examples:
  # Push a new snippet
  gdsnip push -f template.yml --name "My Template" --public

  # Push to a team
  gdsnip push -f template.yml --name "Team Template" --team payments

  # Interactive mode
  gdsnip push
```

#### gdsnip search

```
Search for snippets

Usage:
  gdsnip search <query> [flags]

Flags:
  -t, --tags strings    Filter by tags
      --team string     Search within a team
      --limit int       Maximum results (default 20)

Examples:
  gdsnip search docker
  gdsnip search kafka --tags=messaging,docker
  gdsnip search --team payments
```

#### gdsnip list

```
List your snippets

Usage:
  gdsnip list [flags]

Flags:
      --team string   List team snippets
      --public        Show only public snippets
      --private       Show only private snippets

Examples:
  gdsnip list
  gdsnip list --team payments
  gdsnip list --public
```

#### gdsnip info

```
Get detailed information about a snippet

Usage:
  gdsnip info <namespace/slug>

Examples:
  gdsnip info docker-pg
  gdsnip info johndoe/docker-pg
  gdsnip info payments/kafka-setup
```

### CLI Output Examples

#### gdsnip search docker

```
$ gdsnip search docker

Found 5 snippets:

  NAME               NAMESPACE    DESCRIPTION                      PULLS   UPDATED
  docker-pg          public       PostgreSQL with Docker           234     2 days ago
  docker-redis       public       Redis with persistence           189     1 week ago
  docker-kafka       payments     Kafka + Zookeeper setup          45      3 days ago
  docker-mongo       johndoe      MongoDB replica set              67      2 weeks ago
  docker-nginx       public       Nginx reverse proxy              156     5 days ago

Use: gdsnip pull <namespace>/<name>
```

#### gdsnip pull docker-pg --DB_PASSWORD=secret

```
$ gdsnip pull docker-pg --DB_PASSWORD=secret

✔ Found: public/docker-pg (v2)
✔ Variables resolved
✔ Template rendered

──────────────────────────────────────────────────────────────────
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: mydb
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secret
    ports:
      - '5432:5432'
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
──────────────────────────────────────────────────────────────────

? Save to file? › (y/N)
```

#### gdsnip info johndoe/docker-pg

```
$ gdsnip info johndoe/docker-pg

┌───────────────────────────────────────────────────────────────────┐
│  johndoe/docker-pg                                                │
├───────────────────────────────────────────────────────────────────┤
│  Name:         Docker PostgreSQL                                  │
│  Description:  Docker Compose setup for PostgreSQL                │
│  Visibility:   public                                             │
│  Version:      2                                                  │
│  Created:      2024-01-10 by johndoe                              │
│  Updated:      2024-01-14                                         │
│  Pulls:        234                                                │
├───────────────────────────────────────────────────────────────────┤
│  Variables:                                                       │
│                                                                   │
│    Name           Required   Default   Description                │
│    ─────────────────────────────────────────────────────────────  │
│    PG_VERSION     no         15        PostgreSQL version         │
│    DB_NAME        no         mydb      Database name              │
│    DB_PASSWORD    yes        -         Database password          │
│    PORT           no         5432      Host port                  │
│                                                                   │
├───────────────────────────────────────────────────────────────────┤
│  Tags:  docker, postgres, database                                │
└───────────────────────────────────────────────────────────────────┘

Usage:
  gdsnip pull johndoe/docker-pg --DB_PASSWORD=<value>
```

#### gdsnip auth status

```
$ gdsnip auth status

✔ Logged in as johndoe
  Email: john@example.com
  Since: 2024-01-15 10:30:00
```

---

## Template Syntax

### Variable Format

Variables use double curly braces:

```
{{VARIABLE_NAME}}
```

### Rules

1. Variable names must be UPPERCASE
2. Variable names can contain letters, numbers, and underscores
3. Variable names must start with a letter
4. Valid: `{{DB_HOST}}`, `{{PORT_8080}}`, `{{MY_VAR}}`
5. Invalid: `{{db_host}}`, `{{8080_PORT}}`, `{{MY-VAR}}`

### Example Template

```yaml
# docker-compose.yml template
version: '3.8'

services:
  {{SERVICE_NAME}}:
    image: {{IMAGE_NAME}}:{{IMAGE_TAG}}
    environment:
      DATABASE_URL: postgres://{{DB_USER}}:{{DB_PASSWORD}}@{{DB_HOST}}:{{DB_PORT}}/{{DB_NAME}}
    ports:
      - "{{HOST_PORT}}:{{CONTAINER_PORT}}"
    restart: {{RESTART_POLICY}}
```

### Variable Definition

When pushing, variables are defined with metadata:

```json
{
  "variables": [
    {
      "name": "SERVICE_NAME",
      "description": "Name of the Docker service",
      "required": true
    },
    {
      "name": "IMAGE_TAG",
      "description": "Docker image tag",
      "required": false,
      "default": "latest"
    },
    {
      "name": "DB_PASSWORD",
      "description": "Database password",
      "required": true,
      "sensitive": true
    }
  ]
}
```

---

## Project Structure

```
gdsnip/
├── cmd/
│   ├── cli/
│   │   └── main.go                  # CLI entry point
│   └── api/
│       └── main.go                  # API entry point
│
├── pkg/                              # Shared packages (CLI + API)
│   ├── models/
│   │   ├── user.go                  # User model
│   │   ├── snippet.go               # Snippet model
│   │   └── team.go                  # Team model
│   ├── template/
│   │   ├── engine.go                # Template rendering
│   │   └── engine_test.go           # Template tests
│   └── validator/
│       └── variables.go             # Variable validation
│
├── internal/
│   ├── cli/
│   │   ├── commands/
│   │   │   ├── root.go              # Root command
│   │   │   ├── auth.go              # Auth commands
│   │   │   ├── pull.go              # Pull command
│   │   │   ├── push.go              # Push command
│   │   │   ├── search.go            # Search command
│   │   │   ├── list.go              # List command
│   │   │   └── info.go              # Info command
│   │   ├── tui/
│   │   │   ├── login.go             # Login TUI (bubbletea)
│   │   │   ├── register.go          # Register TUI
│   │   │   ├── prompt.go            # Variable prompts
│   │   │   └── styles.go            # Lipgloss styles
│   │   ├── api/
│   │   │   └── client.go            # API HTTP client
│   │   └── config/
│   │       └── config.go            # CLI config (koanf)
│   │
│   └── api/
│       ├── server/
│       │   └── server.go            # Echo server setup
│       ├── handlers/
│       │   ├── auth.go              # Auth handlers
│       │   ├── snippets.go          # Snippet handlers
│       │   ├── users.go             # User handlers
│       │   └── teams.go             # Team handlers
│       ├── middleware/
│       │   ├── auth.go              # JWT middleware
│       │   ├── ratelimit.go         # Rate limiting
│       │   └── cors.go              # CORS middleware
│       ├── repository/
│       │   ├── db.go                # Database connection
│       │   ├── user.go              # User queries
│       │   ├── snippet.go           # Snippet queries
│       │   └── team.go              # Team queries
│       └── config/
│           └── config.go            # API config (koanf)
│
├── migrations/
│   ├── 001_initial.up.sql           # Initial schema
│   └── 001_initial.down.sql         # Rollback
│
├── scripts/
│   └── seed.go                      # Seed data
│
├── docker/
│   ├── Dockerfile.api               # API Docker image
│   └── Dockerfile.cli               # CLI Docker image
│
├── .github/
│   └── workflows/
│       └── ci.yml                   # CI/CD pipeline
│
├── Taskfile.yml                      # Task runner commands
├── docker-compose.yml                # Local dev environment
├── go.mod                            # Go modules
├── go.sum
└── README.md
```

---

## Core Implementation

### Shared Template Engine

```go
// pkg/template/engine.go
package template

import (
	"fmt"
	"regexp"
	"strings"
)

var variablePattern = regexp.MustCompile(`\{\{([A-Z_][A-Z0-9_]*)\}\}`)

type Variable struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
	Sensitive   bool   `json:"sensitive,omitempty"`
}

type RenderResult struct {
	Content  string
	Warnings []string
}

// ExtractVariables finds all {{VAR}} patterns in template
func ExtractVariables(content string) []string {
	matches := variablePattern.FindAllStringSubmatch(content, -1)
	
	seen := make(map[string]bool)
	var vars []string
	
	for _, match := range matches {
		varName := match[1]
		if !seen[varName] {
			seen[varName] = true
			vars = append(vars, varName)
		}
	}
	
	return vars
}

// Validate checks all required variables are provided
func Validate(definitions []Variable, provided map[string]string) error {
	var missing []string
	
	for _, v := range definitions {
		if v.Required {
			val, exists := provided[v.Name]
			if !exists || val == "" {
				if v.Default == "" {
					missing = append(missing, v.Name)
				}
			}
		}
	}
	
	if len(missing) > 0 {
		return fmt.Errorf("missing required variables: %s", strings.Join(missing, ", "))
	}
	
	return nil
}

// Render substitutes variables in template
func Render(content string, definitions []Variable, provided map[string]string) (*RenderResult, error) {
	// Validate first
	if err := Validate(definitions, provided); err != nil {
		return nil, err
	}
	
	// Build lookup map with defaults
	values := make(map[string]string)
	for _, def := range definitions {
		if def.Default != "" {
			values[def.Name] = def.Default
		}
	}
	for k, v := range provided {
		values[k] = v
	}
	
	// Replace variables
	result := variablePattern.ReplaceAllStringFunc(content, func(match string) string {
		varName := match[2 : len(match)-2] // Strip {{ and }}
		if val, ok := values[varName]; ok {
			return val
		}
		return match // Leave unreplaced if not found
	})
	
	return &RenderResult{Content: result}, nil
}
```

### Shared Models

```go
// pkg/models/snippet.go
package models

import (
	"time"
	
	"github.com/yourorg/gdsnip/pkg/template"
)

type Snippet struct {
	ID          string              `json:"id" db:"id"`
	Name        string              `json:"name" db:"name"`
	Slug        string              `json:"slug" db:"slug"`
	Description string              `json:"description" db:"description"`
	Content     string              `json:"content" db:"content"`
	Variables   []template.Variable `json:"variables" db:"variables"`
	Tags        []string            `json:"tags" db:"tags"`
	Visibility  string              `json:"visibility" db:"visibility"`
	OwnerType   string              `json:"owner_type" db:"owner_type"`
	OwnerID     string              `json:"owner_id" db:"owner_id"`
	Namespace   string              `json:"namespace"`
	Version     int                 `json:"version" db:"version"`
	CreatedBy   string              `json:"created_by" db:"created_by"`
	CreatedAt   time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" db:"updated_at"`
}

// FullPath returns namespace/slug format
func (s *Snippet) FullPath() string {
	return s.Namespace + "/" + s.Slug
}

// PullCommand returns example CLI usage
func (s *Snippet) PullCommand() string {
	cmd := "gdsnip pull " + s.FullPath()
	for _, v := range s.Variables {
		if v.Required && v.Default == "" {
			cmd += " --" + v.Name + "=<value>"
		}
	}
	return cmd
}
```

```go
// pkg/models/user.go
package models

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserPublic struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (u *User) Public() UserPublic {
	return UserPublic{
		ID:       u.ID,
		Username: u.Username,
	}
}
```

### API Server (Echo)

```go
// cmd/api/main.go
package main

import (
	"log"
	
	"github.com/yourorg/gdsnip/internal/api/config"
	"github.com/yourorg/gdsnip/internal/api/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Starting server on :%s", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
```

```go
// internal/api/server/server.go
package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yourorg/gdsnip/internal/api/config"
	"github.com/yourorg/gdsnip/internal/api/handlers"
	customMiddleware "github.com/yourorg/gdsnip/internal/api/middleware"
	"github.com/yourorg/gdsnip/internal/api/repository"
)

type Server struct {
	echo   *echo.Echo
	config *config.Config
}

func New(cfg *config.Config) (*Server, error) {
	// Initialize DB
	db, err := repository.NewDB(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	
	// Initialize repositories
	userRepo := repository.NewUserRepo(db)
	snippetRepo := repository.NewSnippetRepo(db)
	
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
	snippetHandler := handlers.NewSnippetHandler(snippetRepo, userRepo)
	
	// Setup Echo
	e := echo.New()
	e.HideBanner = true
	
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	
	// Public routes
	v1 := e.Group("/v1")
	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)
	v1.GET("/snippets", snippetHandler.Search)
	v1.GET("/snippets/:namespace/:slug", snippetHandler.Get)
	
	// Protected routes
	protected := v1.Group("")
	protected.Use(customMiddleware.JWT(cfg.JWTSecret))
	protected.GET("/auth/me", authHandler.Me)
	protected.POST("/snippets", snippetHandler.Create)
	protected.POST("/snippets/:namespace/:slug/pull", snippetHandler.Pull)
	protected.PUT("/snippets/:id", snippetHandler.Update)
	protected.DELETE("/snippets/:id", snippetHandler.Delete)
	protected.GET("/users/me/snippets", snippetHandler.ListMine)
	
	return &Server{echo: e, config: cfg}, nil
}

func (s *Server) Start() error {
	return s.echo.Start(":" + s.config.Port)
}
```

### API Config (Koanf)

```go
// internal/api/config/config.go
package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Port        string `koanf:"port"`
	DatabaseURL string `koanf:"database_url"`
	JWTSecret   string `koanf:"jwt_secret"`
	JWTExpiry   string `koanf:"jwt_expiry"`
}

func Load() (*Config, error) {
	k := koanf.New(".")
	
	// Load from config file if exists
	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		// Ignore if file doesn't exist
	}
	
	// Load from environment variables (override file config)
	// GDSNIP_PORT, GDSNIP_DATABASE_URL, etc.
	k.Load(env.Provider("GDSNIP_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "GDSNIP_"))
	}), nil)
	
	// Defaults
	cfg := &Config{
		Port:      "8080",
		JWTExpiry: "24h",
	}
	
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}
```

### API Handlers (Echo)

```go
// internal/api/handlers/snippets.go
package handlers

import (
	"net/http"
	
	"github.com/labstack/echo/v4"
	"github.com/yourorg/gdsnip/internal/api/repository"
	"github.com/yourorg/gdsnip/pkg/models"
	"github.com/yourorg/gdsnip/pkg/template"
)

type SnippetHandler struct {
	snippetRepo *repository.SnippetRepo
	userRepo    *repository.UserRepo
}

func NewSnippetHandler(sr *repository.SnippetRepo, ur *repository.UserRepo) *SnippetHandler {
	return &SnippetHandler{
		snippetRepo: sr,
		userRepo:    ur,
	}
}

// POST /v1/snippets
func (h *SnippetHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(string)
	
	var req struct {
		Name        string              `json:"name" validate:"required"`
		Slug        string              `json:"slug" validate:"required"`
		Description string              `json:"description"`
		Content     string              `json:"content" validate:"required"`
		Variables   []template.Variable `json:"variables"`
		Tags        []string            `json:"tags"`
		Visibility  string              `json:"visibility"`
		TeamSlug    string              `json:"team_slug"`
	}
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"message": err.Error()},
		})
	}
	
	snippet := &models.Snippet{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Content:     req.Content,
		Variables:   req.Variables,
		Tags:        req.Tags,
		Visibility:  req.Visibility,
		OwnerType:   "user",
		OwnerID:     userID,
		CreatedBy:   userID,
		Version:     1,
	}
	
	if err := h.snippetRepo.Create(snippet); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"message": err.Error()},
		})
	}
	
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    snippet,
	})
}

// GET /v1/snippets/:namespace/:slug
func (h *SnippetHandler) Get(c echo.Context) error {
	namespace := c.Param("namespace")
	slug := c.Param("slug")
	
	snippet, err := h.snippetRepo.GetByPath(namespace, slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"message": "snippet not found"},
		})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    snippet,
	})
}

// POST /v1/snippets/:namespace/:slug/pull
func (h *SnippetHandler) Pull(c echo.Context) error {
	namespace := c.Param("namespace")
	slug := c.Param("slug")
	
	var req struct {
		Variables map[string]string `json:"variables"`
	}
	
	if err := c.Bind(&req); err != nil {
		req.Variables = make(map[string]string)
	}
	
	// Get snippet
	snippet, err := h.snippetRepo.GetByPath(namespace, slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"message": "snippet not found"},
		})
	}
	
	// Render template
	result, err := template.Render(snippet.Content, snippet.Variables, req.Variables)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"message": err.Error()},
		})
	}
	
	// Track pull (async)
	userID, _ := c.Get("user_id").(string)
	go h.snippetRepo.TrackPull(snippet.ID, userID)
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"snippet":  snippet.FullPath(),
			"version":  snippet.Version,
			"rendered": result.Content,
		},
	})
}

// GET /v1/snippets?q=docker&tags=postgres
func (h *SnippetHandler) Search(c echo.Context) error {
	query := c.QueryParam("q")
	tags := c.QueryParams()["tags"]
	limit := 20
	
	snippets, err := h.snippetRepo.Search(query, tags, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"message": err.Error()},
		})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"snippets": snippets,
			"total":    len(snippets),
		},
	})
}

// GET /v1/users/me/snippets
func (h *SnippetHandler) ListMine(c echo.Context) error {
	userID := c.Get("user_id").(string)
	
	snippets, err := h.snippetRepo.ListByOwner("user", userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"message": err.Error()},
		})
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"snippets": snippets,
			"total":    len(snippets),
		},
	})
}
```

### Auth Middleware (Echo + JWT)

```go
// internal/api/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"
	
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWT(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   map[string]string{"message": "missing authorization header"},
				})
			}
			
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   map[string]string{"message": "invalid authorization header"},
				})
			}
			
			tokenString := parts[1]
			
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   map[string]string{"message": "invalid token"},
				})
			}
			
			claims := token.Claims.(jwt.MapClaims)
			c.Set("user_id", claims["user_id"])
			c.Set("username", claims["username"])
			
			return next(c)
		}
	}
}
```

### CLI Entry Point

```go
// cmd/cli/main.go
package main

import (
	"os"
	
	"github.com/yourorg/gdsnip/internal/cli/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
```

### CLI Root Command (Cobra)

```go
// internal/cli/commands/root.go
package commands

import (
	"github.com/spf13/cobra"
	"github.com/yourorg/gdsnip/internal/cli/config"
)

var (
	cfgFile   string
	jsonOutput bool
	noColor   bool
)

var rootCmd = &cobra.Command{
	Use:   "gdsnip",
	Short: "A CLI snippet registry for developers",
	Long: `GDSNIP is a developer productivity tool that provides a CLI-first 
snippet registry for storing, sharing, and pulling parameterized 
code/config templates.

Examples:
  gdsnip search docker
  gdsnip pull johndoe/docker-pg --DB_PASSWORD=secret
  gdsnip push -f template.yml --name my-template --public`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gdsnip/config.yml)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	
	// Add subcommands
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(infoCmd)
}

func initConfig() {
	config.Init(cfgFile)
}
```

### CLI Pull Command

```go
// internal/cli/commands/pull.go
package commands

import (
	"fmt"
	"os"
	"strings"
	
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/yourorg/gdsnip/internal/cli/api"
	"github.com/yourorg/gdsnip/internal/cli/tui"
)

var (
	outputFile string
	varFile    string
	version    int
)

var pullCmd = &cobra.Command{
	Use:   "pull <namespace/slug>",
	Short: "Pull and render a snippet",
	Long: `Pull a snippet from the registry and render it with provided variables.
	
Variables are passed as flags in the format --VAR_NAME=value`,
	Example: `  gdsnip pull docker-pg --DB_PASSWORD=secret
  gdsnip pull payments/kafka-setup --BROKER_URL=localhost:9092
  gdsnip pull docker-pg -o docker-compose.yml`,
	Args: cobra.ExactArgs(1),
	Run:  runPull,
}

func init() {
	pullCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	pullCmd.Flags().StringVar(&varFile, "var-file", "", "JSON file containing variables")
	pullCmd.Flags().IntVar(&version, "version", 0, "Specific version to pull")
}

func runPull(cmd *cobra.Command, args []string) {
	snippetPath := args[0]
	
	// Parse --VAR=value flags
	variables := parseVariableFlags(os.Args)
	
	// Create API client
	client := api.NewClient()
	
	// Styles
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	
	fmt.Printf("%s Fetching %s...\n", dimStyle.Render("⏳"), snippetPath)
	
	// Pull snippet
	result, err := client.PullSnippet(snippetPath, variables)
	if err != nil {
		// Check if missing variables
		if missingVars := api.ParseMissingVarsError(err); missingVars != nil {
			fmt.Println(errorStyle.Render("✗ Missing required variables:"))
			fmt.Println()
			for _, v := range missingVars {
				fmt.Printf("  %s (required)\n", v.Name)
				if v.Description != "" {
					fmt.Printf("    %s\n", dimStyle.Render(v.Description))
				}
			}
			fmt.Println()
			
			// Prompt for missing variables interactively
			for _, v := range missingVars {
				value := tui.PromptVariable(v.Name, v.Description, v.Sensitive)
				variables[v.Name] = value
			}
			
			// Retry with filled variables
			result, err = client.PullSnippet(snippetPath, variables)
			if err != nil {
				fmt.Println(errorStyle.Render("✗ Error: " + err.Error()))
				os.Exit(1)
			}
		} else {
			fmt.Println(errorStyle.Render("✗ Error: " + err.Error()))
			os.Exit(1)
		}
	}
	
	fmt.Printf("%s Found: %s (v%d)\n", successStyle.Render("✔"), result.Snippet, result.Version)
	fmt.Printf("%s Variables resolved\n", successStyle.Render("✔"))
	fmt.Printf("%s Template rendered\n", successStyle.Render("✔"))
	
	// Print rendered content
	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println(result.Rendered)
	fmt.Println(strings.Repeat("─", 60))
	
	// Save to file if requested
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(result.Rendered), 0644)
		if err != nil {
			fmt.Println(errorStyle.Render("✗ Failed to write file: " + err.Error()))
			os.Exit(1)
		}
		fmt.Printf("%s Saved to %s\n", successStyle.Render("✔"), outputFile)
	} else {
		// Prompt to save
		if tui.ConfirmSave() {
			filename := tui.PromptFilename(result.FilenameSuggestion)
			err := os.WriteFile(filename, []byte(result.Rendered), 0644)
			if err != nil {
				fmt.Println(errorStyle.Render("✗ Failed to write file: " + err.Error()))
				os.Exit(1)
			}
			fmt.Printf("%s Saved to %s\n", successStyle.Render("✔"), filename)
		}
	}
}

func parseVariableFlags(args []string) map[string]string {
	variables := make(map[string]string)
	
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
			arg = strings.TrimPrefix(arg, "--")
			parts := strings.SplitN(arg, "=", 2)
			
			// Check if it's an uppercase variable name
			if len(parts) == 2 && strings.ToUpper(parts[0]) == parts[0] {
				variables[parts[0]] = parts[1]
			}
		}
	}
	
	return variables
}
```

### CLI TUI Components (Bubbletea)

```go
// internal/cli/tui/prompt.go
package tui

import (
	"fmt"
	
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
)

type promptModel struct {
	textInput textinput.Model
	label     string
	done      bool
}

func newPromptModel(label string, sensitive bool) promptModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 40
	
	if sensitive {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}
	
	return promptModel{
		textInput: ti,
		label:     label,
		done:      false,
	}
}

func (m promptModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m promptModel) View() string {
	return fmt.Sprintf(
		"? %s: %s",
		focusedStyle.Render(m.label),
		m.textInput.View(),
	)
}

// PromptVariable prompts user for a variable value
func PromptVariable(name, description string, sensitive bool) string {
	label := name
	if description != "" {
		label = fmt.Sprintf("%s (%s)", name, description)
	}
	
	m := newPromptModel(label, sensitive)
	p := tea.NewProgram(m)
	
	finalModel, err := p.Run()
	if err != nil {
		return ""
	}
	
	return finalModel.(promptModel).textInput.Value()
}

// ConfirmSave asks user if they want to save to file
func ConfirmSave() bool {
	// Simple y/n prompt using bubbletea
	// Implementation similar to above
	return false // Simplified for brevity
}

// PromptFilename asks for output filename
func PromptFilename(suggestion string) string {
	// Implementation similar to PromptVariable
	return suggestion
}
```

```go
// internal/cli/tui/login.go
package tui

import (
	"fmt"
	"strings"
	
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loginModel struct {
	focusIndex int
	inputs     []textinput.Model
	done       bool
	cancelled  bool
}

type LoginResult struct {
	Email    string
	Password string
}

func NewLoginModel() loginModel {
	m := loginModel{
		inputs: make([]textinput.Model, 2),
	}
	
	// Email input
	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "email@example.com"
	m.inputs[0].Focus()
	m.inputs[0].CharLimit = 64
	m.inputs[0].Width = 40
	m.inputs[0].Prompt = ""
	
	// Password input
	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "password"
	m.inputs[1].CharLimit = 64
	m.inputs[1].Width = 40
	m.inputs[1].EchoMode = textinput.EchoPassword
	m.inputs[1].EchoCharacter = '•'
	m.inputs[1].Prompt = ""
	
	return m
}

func (m loginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelled = true
			return m, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
			if msg.Type == tea.KeyEnter && m.focusIndex == len(m.inputs)-1 {
				m.done = true
				return m, tea.Quit
			}
			
			// Move focus
			if msg.Type == tea.KeyUp || msg.Type == tea.KeyShiftTab {
				m.focusIndex--
			} else {
				m.focusIndex++
			}
			
			if m.focusIndex > len(m.inputs)-1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}
			
			for i := range m.inputs {
				if i == m.focusIndex {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
		}
	}
	
	// Update all inputs
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	
	return m, tea.Batch(cmds...)
}

func (m loginModel) View() string {
	var b strings.Builder
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))
	
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
	
	b.WriteString(titleStyle.Render("Login to GDSNIP"))
	b.WriteString("\n\n")
	
	labels := []string{"Email", "Password"}
	
	for i, input := range m.inputs {
		b.WriteString(labelStyle.Render(labels[i]))
		b.WriteString("\n")
		b.WriteString(input.View())
		b.WriteString("\n\n")
	}
	
	b.WriteString(labelStyle.Render("Press Enter to submit, Esc to cancel"))
	
	return b.String()
}

func (m loginModel) Result() *LoginResult {
	if m.cancelled || !m.done {
		return nil
	}
	return &LoginResult{
		Email:    m.inputs[0].Value(),
		Password: m.inputs[1].Value(),
	}
}

// RunLogin runs the login TUI and returns credentials
func RunLogin() *LoginResult {
	m := NewLoginModel()
	p := tea.NewProgram(m)
	
	finalModel, err := p.Run()
	if err != nil {
		return nil
	}
	
	return finalModel.(loginModel).Result()
}
```

```go
// internal/cli/tui/styles.go
package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryColor   = lipgloss.Color("205")
	SecondaryColor = lipgloss.Color("241")
	SuccessColor   = lipgloss.Color("42")
	ErrorColor     = lipgloss.Color("196")
	WarningColor   = lipgloss.Color("214")
	
	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor)
	
	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor)
	
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor)
	
	DimStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor)
	
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2)
	
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(PrimaryColor).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder())
)
```

### CLI Config (Koanf)

```go
// internal/cli/config/config.go
package config

import (
	"os"
	"path/filepath"
	
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

type Config struct {
	APIURL string `koanf:"api_url"`
}

type Credentials struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
}

func Init(cfgFile string) {
	// Default config path
	if cfgFile == "" {
		home, _ := os.UserHomeDir()
		cfgFile = filepath.Join(home, ".gdsnip", "config.yml")
	}
	
	// Load config file if exists
	if _, err := os.Stat(cfgFile); err == nil {
		k.Load(file.Provider(cfgFile), yaml.Parser())
	}
	
	// Load from environment (GDSNIP_API_URL, etc.)
	k.Load(env.Provider("GDSNIP_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "GDSNIP_"))
	}), nil)
}

func Get() *Config {
	cfg := &Config{
		APIURL: "http://localhost:8080/v1",
	}
	k.Unmarshal("", cfg)
	return cfg
}

func GetAPIURL() string {
	return Get().APIURL
}

// Credentials management
func configDir() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".gdsnip")
	os.MkdirAll(dir, 0700)
	return dir
}

func SaveCredentials(creds *Credentials) error {
	path := filepath.Join(configDir(), "credentials.json")
	data, _ := json.Marshal(creds)
	return os.WriteFile(path, data, 0600)
}

func LoadCredentials() (*Credentials, error) {
	path := filepath.Join(configDir(), "credentials.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	
	return &creds, nil
}

func ClearCredentials() error {
	path := filepath.Join(configDir(), "credentials.json")
	return os.Remove(path)
}

func GetToken() (string, error) {
	creds, err := LoadCredentials()
	if err != nil {
		return "", err
	}
	return creds.AccessToken, nil
}
```

### CLI API Client

```go
// internal/cli/api/client.go
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	
	"github.com/yourorg/gdsnip/internal/cli/config"
	"github.com/yourorg/gdsnip/pkg/template"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

type PullResponse struct {
	Snippet            string `json:"snippet"`
	Version            int    `json:"version"`
	Rendered           string `json:"rendered"`
	FilenameSuggestion string `json:"filename_suggestion"`
}

type SearchResult struct {
	Snippets []SnippetSummary `json:"snippets"`
	Total    int              `json:"total"`
}

type SnippetSummary struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Namespace   string   `json:"namespace"`
	FullPath    string   `json:"full_path"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Visibility  string   `json:"visibility"`
	Version     int      `json:"version"`
	PullCount   int      `json:"pull_count"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

func NewClient() *Client {
	token, _ := config.GetToken()
	
	return &Client{
		baseURL:    config.GetAPIURL(),
		httpClient: &http.Client{},
		token:      token,
	}
}

func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		data, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(data)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}
	
	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	
	return c.httpClient.Do(req)
}

func (c *Client) PullSnippet(path string, variables map[string]string) (*PullResponse, error) {
	resp, err := c.doRequest("POST", "/snippets/"+path+"/pull", map[string]interface{}{
		"variables": variables,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to API: %w", err)
	}
	defer resp.Body.Close()
	
	var result struct {
		Success bool         `json:"success"`
		Data    PullResponse `json:"data"`
		Error   *APIError    `json:"error"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	
	if !result.Success {
		return nil, result.Error
	}
	
	return &result.Data, nil
}

func (c *Client) Search(query string, tags []string) (*SearchResult, error) {
	path := fmt.Sprintf("/snippets?q=%s", query)
	for _, tag := range tags {
		path += "&tags=" + tag
	}
	
	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Success bool         `json:"success"`
		Data    SearchResult `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	
	return &result.Data, nil
}

func (c *Client) Register(email, username, password string) (*AuthResponse, error) {
	resp, err := c.doRequest("POST", "/auth/register", map[string]string{
		"email":    email,
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Success bool         `json:"success"`
		Data    AuthResponse `json:"data"`
		Error   *APIError    `json:"error"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	
	if !result.Success {
		return nil, result.Error
	}
	
	return &result.Data, nil
}

func (c *Client) Login(email, password string) (*AuthResponse, error) {
	resp, err := c.doRequest("POST", "/auth/login", map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Success bool         `json:"success"`
		Data    AuthResponse `json:"data"`
		Error   *APIError    `json:"error"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	
	if !result.Success {
		return nil, result.Error
	}
	
	return &result.Data, nil
}

type AuthResponse struct {
	User  UserInfo `json:"user"`
	Token string   `json:"token"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

func (e *APIError) Error() string {
	return e.Message
}

// ParseMissingVarsError extracts missing variables from API error
func ParseMissingVarsError(err error) []template.Variable {
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.Code != "MISSING_VARIABLES" {
		return nil
	}
	
	missing, ok := apiErr.Details["missing"].([]interface{})
	if !ok {
		return nil
	}
	
	var vars []template.Variable
	for _, m := range missing {
		if mv, ok := m.(map[string]interface{}); ok {
			vars = append(vars, template.Variable{
				Name:        mv["name"].(string),
				Description: mv["description"].(string),
				Required:    true,
			})
		}
	}
	
	return vars
}
```

---

## Taskfile

```yaml
# Taskfile.yml
version: '3'

vars:
  BINARY_NAME: gdsnip
  API_BINARY: gdsnip-api
  VERSION: 0.1.0

tasks:
  default:
    cmds:
      - task --list

  # Development
  dev:api:
    desc: Run API with hot reload
    cmds:
      - air -c .air.toml
    dir: .

  dev:cli:
    desc: Build and install CLI locally
    cmds:
      - go build -o ./bin/{{.BINARY_NAME}} ./cmd/cli
      - cp ./bin/{{.BINARY_NAME}} /usr/local/bin/

  # Build
  build:
    desc: Build all binaries
    cmds:
      - task: build:cli
      - task: build:api

  build:cli:
    desc: Build CLI binary
    cmds:
      - go build -ldflags="-X main.version={{.VERSION}}" -o ./bin/{{.BINARY_NAME}} ./cmd/cli

  build:api:
    desc: Build API binary
    cmds:
      - go build -ldflags="-X main.version={{.VERSION}}" -o ./bin/{{.API_BINARY}} ./cmd/api

  build:all:
    desc: Build for all platforms
    cmds:
      - GOOS=linux GOARCH=amd64 go build -o ./dist/{{.BINARY_NAME}}-linux-amd64 ./cmd/cli
      - GOOS=linux GOARCH=arm64 go build -o ./dist/{{.BINARY_NAME}}-linux-arm64 ./cmd/cli
      - GOOS=darwin GOARCH=amd64 go build -o ./dist/{{.BINARY_NAME}}-darwin-amd64 ./cmd/cli
      - GOOS=darwin GOARCH=arm64 go build -o ./dist/{{.BINARY_NAME}}-darwin-arm64 ./cmd/cli
      - GOOS=windows GOARCH=amd64 go build -o ./dist/{{.BINARY_NAME}}-windows-amd64.exe ./cmd/cli

  # Database
  db:up:
    desc: Start PostgreSQL
    cmds:
      - docker-compose up -d postgres

  db:down:
    desc: Stop PostgreSQL
    cmds:
      - docker-compose down

  db:migrate:
    desc: Run database migrations
    cmds:
      - go run ./scripts/migrate.go up

  db:migrate:down:
    desc: Rollback database migrations
    cmds:
      - go run ./scripts/migrate.go down

  db:seed:
    desc: Seed database with sample data
    cmds:
      - go run ./scripts/seed.go

  db:reset:
    desc: Reset database (migrate down, up, seed)
    cmds:
      - task: db:migrate:down
      - task: db:migrate
      - task: db:seed

  # Testing
  test:
    desc: Run all tests
    cmds:
      - go test ./...

  test:cover:
    desc: Run tests with coverage
    cmds:
      - go test -coverprofile=coverage.out ./...
      - go tool cover -html=coverage.out -o coverage.html

  test:integration:
    desc: Run integration tests
    cmds:
      - go test ./tests/integration/... -tags=integration

  # Docker
  docker:build:
    desc: Build Docker images
    cmds:
      - docker build -t gdsnip-api -f docker/Dockerfile.api .
      - docker build -t gdsnip-cli -f docker/Dockerfile.cli .

  docker:up:
    desc: Start all services
    cmds:
      - docker-compose up -d

  docker:down:
    desc: Stop all services
    cmds:
      - docker-compose down

  docker:logs:
    desc: View logs
    cmds:
      - docker-compose logs -f

  # Utilities
  clean:
    desc: Clean build artifacts
    cmds:
      - rm -rf bin/ dist/ coverage.out coverage.html

  lint:
    desc: Run linter
    cmds:
      - golangci-lint run

  fmt:
    desc: Format code
    cmds:
      - go fmt ./...

  tidy:
    desc: Tidy go modules
    cmds:
      - go mod tidy

  # Full setup
  setup:
    desc: Full project setup
    cmds:
      - task: tidy
      - task: db:up
      - sleep 3
      - task: db:migrate
      - task: db:seed
      - echo "Setup complete! Run 'task dev:api' to start the server"
```

---

## Docker Configuration

### docker-compose.yml

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: gdsnip
      POSTGRES_PASSWORD: gdsnip
      POSTGRES_DB: gdsnip
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U gdsnip"]
      interval: 5s
      timeout: 5s
      retries: 5

  api:
    build:
      context: .
      dockerfile: docker/Dockerfile.api
    ports:
      - "8080:8080"
    environment:
      GDSNIP_PORT: "8080"
      GDSNIP_DATABASE_URL: postgres://gdsnip:gdsnip@postgres:5432/gdsnip?sslmode=disable
      GDSNIP_JWT_SECRET: your-secret-key-change-in-production
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  pgdata:
```

### docker/Dockerfile.api

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o gdsnip-api ./cmd/api

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/gdsnip-api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./gdsnip-api"]
```

### docker/Dockerfile.cli

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gdsnip ./cmd/cli

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/gdsnip /usr/local/bin/

ENTRYPOINT ["gdsnip"]
```

---

## Go Dependencies

```go
// go.mod
module github.com/yourorg/gdsnip

go 1.21

require (
	// API Framework
	github.com/labstack/echo/v4 v4.11.4
	
	// Database
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9
	
	// Auth
	github.com/golang-jwt/jwt/v5 v5.2.0
	golang.org/x/crypto v0.18.0
	
	// Config
	github.com/knadh/koanf/v2 v2.0.1
	github.com/knadh/koanf/parsers/yaml v0.1.0
	github.com/knadh/koanf/parsers/json v0.1.0
	github.com/knadh/koanf/providers/env v0.1.0
	github.com/knadh/koanf/providers/file v0.1.0
	
	// CLI
	github.com/spf13/cobra v1.8.0
	github.com/charmbracelet/bubbletea v0.25.0
	github.com/charmbracelet/bubbles v0.18.0
	github.com/charmbracelet/lipgloss v0.9.1
	
	// Utilities
	github.com/google/uuid v1.6.0
)
```

---

## Seed Data

### scripts/seed.go

```go
package main

// Seed snippets for demo
var seedSnippets = []SeedSnippet{
	{
		Name:        "Docker PostgreSQL",
		Slug:        "docker-pg",
		Description: "Docker Compose setup for PostgreSQL with persistent volume",
		Tags:        []string{"docker", "postgres", "database"},
		Content: `version: '3.8'
services:
  postgres:
    image: postgres:{{PG_VERSION}}
    environment:
      POSTGRES_DB: {{DB_NAME}}
      POSTGRES_USER: {{DB_USER}}
      POSTGRES_PASSWORD: {{DB_PASSWORD}}
    ports:
      - '{{PORT}}:5432'
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:`,
		Variables: []Variable{
			{Name: "PG_VERSION", Default: "15", Description: "PostgreSQL version"},
			{Name: "DB_NAME", Default: "mydb", Description: "Database name"},
			{Name: "DB_USER", Default: "postgres", Description: "Database user"},
			{Name: "DB_PASSWORD", Required: true, Sensitive: true, Description: "Database password"},
			{Name: "PORT", Default: "5432", Description: "Host port"},
		},
	},
	{
		Name:        "Docker Redis",
		Slug:        "docker-redis",
		Description: "Redis with persistence and password",
		Tags:        []string{"docker", "redis", "cache"},
		Content: `version: '3.8'
services:
  redis:
    image: redis:{{REDIS_VERSION}}
    command: redis-server --appendonly yes --requirepass {{REDIS_PASSWORD}}
    ports:
      - '{{PORT}}:6379'
    volumes:
      - redis_data:/data

volumes:
  redis_data:`,
		Variables: []Variable{
			{Name: "REDIS_VERSION", Default: "7-alpine", Description: "Redis version"},
			{Name: "REDIS_PASSWORD", Required: true, Sensitive: true, Description: "Redis password"},
			{Name: "PORT", Default: "6379", Description: "Host port"},
		},
	},
	{
		Name:        "Kubernetes Deployment",
		Slug:        "k8s-deployment",
		Description: "Basic Kubernetes deployment manifest",
		Tags:        []string{"kubernetes", "k8s", "deployment"},
		Content: `apiVersion: apps/v1
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
          resources:
            requests:
              memory: "{{MEMORY_REQUEST}}"
              cpu: "{{CPU_REQUEST}}"
            limits:
              memory: "{{MEMORY_LIMIT}}"
              cpu: "{{CPU_LIMIT}}"`,
		Variables: []Variable{
			{Name: "APP_NAME", Required: true, Description: "Application name"},
			{Name: "NAMESPACE", Default: "default", Description: "Kubernetes namespace"},
			{Name: "REPLICAS", Default: "3", Description: "Number of replicas"},
			{Name: "IMAGE", Required: true, Description: "Container image"},
			{Name: "TAG", Default: "latest", Description: "Image tag"},
			{Name: "PORT", Default: "8080", Description: "Container port"},
			{Name: "MEMORY_REQUEST", Default: "128Mi", Description: "Memory request"},
			{Name: "CPU_REQUEST", Default: "100m", Description: "CPU request"},
			{Name: "MEMORY_LIMIT", Default: "256Mi", Description: "Memory limit"},
			{Name: "CPU_LIMIT", Default: "200m", Description: "CPU limit"},
		},
	},
	{
		Name:        "GitHub Actions CI",
		Slug:        "github-action-ci",
		Description: "Basic GitHub Actions CI workflow",
		Tags:        []string{"github", "ci", "actions"},
		Content: `name: {{WORKFLOW_NAME}}

on:
  push:
    branches: [{{BRANCH}}]
  pull_request:
    branches: [{{BRANCH}}]

jobs:
  build:
    runs-on: {{RUNNER}}
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '{{GO_VERSION}}'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run tests
        run: go test ./...
      
      - name: Build
        run: go build -v ./...`,
		Variables: []Variable{
			{Name: "WORKFLOW_NAME", Default: "CI", Description: "Workflow name"},
			{Name: "BRANCH", Default: "main", Description: "Branch to trigger on"},
			{Name: "RUNNER", Default: "ubuntu-latest", Description: "GitHub runner"},
			{Name: "GO_VERSION", Default: "1.21", Description: "Go version"},
		},
	},
	{
		Name:        "Go Dockerfile",
		Slug:        "dockerfile-go",
		Description: "Multi-stage Dockerfile for Go applications",
		Tags:        []string{"docker", "go", "golang", "dockerfile"},
		Content: `# Build stage
FROM golang:{{GO_VERSION}}-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o {{BINARY_NAME}} ./{{MAIN_PATH}}

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/{{BINARY_NAME}} .

EXPOSE {{PORT}}

CMD ["./{{BINARY_NAME}}"]`,
		Variables: []Variable{
			{Name: "GO_VERSION", Default: "1.21", Description: "Go version"},
			{Name: "BINARY_NAME", Required: true, Description: "Output binary name"},
			{Name: "MAIN_PATH", Default: "cmd/main.go", Description: "Path to main.go"},
			{Name: "PORT", Default: "8080", Description: "Exposed port"},
		},
	},
	{
		Name:        "Nginx Reverse Proxy",
		Slug:        "nginx-proxy",
		Description: "Nginx reverse proxy configuration",
		Tags:        []string{"nginx", "proxy", "docker"},
		Content: `server {
    listen {{LISTEN_PORT}};
    server_name {{SERVER_NAME}};

    location / {
        proxy_pass http://{{UPSTREAM_HOST}}:{{UPSTREAM_PORT}};
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}`,
		Variables: []Variable{
			{Name: "LISTEN_PORT", Default: "80", Description: "Nginx listen port"},
			{Name: "SERVER_NAME", Default: "localhost", Description: "Server name"},
			{Name: "UPSTREAM_HOST", Required: true, Description: "Upstream host"},
			{Name: "UPSTREAM_PORT", Default: "8080", Description: "Upstream port"},
		},
	},
	{
		Name:        "Docker Compose Full Stack",
		Slug:        "docker-fullstack",
		Description: "Full stack Docker Compose with API, DB, and Redis",
		Tags:        []string{"docker", "fullstack", "postgres", "redis"},
		Content: `version: '3.8'

services:
  api:
    build: .
    ports:
      - "{{API_PORT}}:{{API_PORT}}"
    environment:
      PORT: {{API_PORT}}
      DATABASE_URL: postgres://{{DB_USER}}:{{DB_PASSWORD}}@postgres:5432/{{DB_NAME}}
      REDIS_URL: redis://:{{REDIS_PASSWORD}}@redis:6379
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: {{DB_USER}}
      POSTGRES_PASSWORD: {{DB_PASSWORD}}
      POSTGRES_DB: {{DB_NAME}}
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass {{REDIS_PASSWORD}}
    volumes:
      - redis_data:/data

volumes:
  pgdata:
  redis_data:`,
		Variables: []Variable{
			{Name: "API_PORT", Default: "8080", Description: "API port"},
			{Name: "DB_USER", Default: "app", Description: "Database user"},
			{Name: "DB_PASSWORD", Required: true, Sensitive: true, Description: "Database password"},
			{Name: "DB_NAME", Default: "app", Description: "Database name"},
			{Name: "REDIS_PASSWORD", Required: true, Sensitive: true, Description: "Redis password"},
		},
	},
	{
		Name:        "Kubernetes Service",
		Slug:        "k8s-service",
		Description: "Kubernetes Service manifest",
		Tags:        []string{"kubernetes", "k8s", "service"},
		Content: `apiVersion: v1
kind: Service
metadata:
  name: {{SERVICE_NAME}}
  namespace: {{NAMESPACE}}
spec:
  type: {{SERVICE_TYPE}}
  selector:
    app: {{APP_NAME}}
  ports:
    - port: {{PORT}}
      targetPort: {{TARGET_PORT}}
      protocol: TCP`,
		Variables: []Variable{
			{Name: "SERVICE_NAME", Required: true, Description: "Service name"},
			{Name: "NAMESPACE", Default: "default", Description: "Namespace"},
			{Name: "SERVICE_TYPE", Default: "ClusterIP", Description: "Service type (ClusterIP, NodePort, LoadBalancer)"},
			{Name: "APP_NAME", Required: true, Description: "App selector label"},
			{Name: "PORT", Default: "80", Description: "Service port"},
			{Name: "TARGET_PORT", Default: "8080", Description: "Target container port"},
		},
	},
	{
		Name:        "Go .gitignore",
		Slug:        "gitignore-go",
		Description: "Comprehensive .gitignore for Go projects",
		Tags:        []string{"go", "git", "gitignore"},
		Content: `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
{{BINARY_NAME}}

# Build directory
bin/
dist/

# Test binary
*.test

# Coverage
coverage.out
coverage.html

# Dependency directories
vendor/

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Environment
.env
.env.local
*.env

# Logs
*.log
logs/`,
		Variables: []Variable{
			{Name: "BINARY_NAME", Default: "app", Description: "Your binary name to ignore"},
		},
	},
	{
		Name:        "Environment File Template",
		Slug:        "env-template",
		Description: "Standard environment file template",
		Tags:        []string{"env", "config", "template"},
		Content: `# Application
APP_NAME={{APP_NAME}}
APP_ENV={{APP_ENV}}
PORT={{PORT}}

# Database
DATABASE_URL=postgres://{{DB_USER}}:{{DB_PASSWORD}}@{{DB_HOST}}:{{DB_PORT}}/{{DB_NAME}}

# Redis
REDIS_URL=redis://:{{REDIS_PASSWORD}}@{{REDIS_HOST}}:{{REDIS_PORT}}

# JWT
JWT_SECRET={{JWT_SECRET}}
JWT_EXPIRY={{JWT_EXPIRY}}

# Logging
LOG_LEVEL={{LOG_LEVEL}}`,
		Variables: []Variable{
			{Name: "APP_NAME", Required: true, Description: "Application name"},
			{Name: "APP_ENV", Default: "development", Description: "Environment (development, staging, production)"},
			{Name: "PORT", Default: "8080", Description: "Application port"},
			{Name: "DB_USER", Default: "postgres", Description: "Database user"},
			{Name: "DB_PASSWORD", Required: true, Sensitive: true, Description: "Database password"},
			{Name: "DB_HOST", Default: "localhost", Description: "Database host"},
			{Name: "DB_PORT", Default: "5432", Description: "Database port"},
			{Name: "DB_NAME", Required: true, Description: "Database name"},
			{Name: "REDIS_PASSWORD", Default: "", Sensitive: true, Description: "Redis password"},
			{Name: "REDIS_HOST", Default: "localhost", Description: "Redis host"},
			{Name: "REDIS_PORT", Default: "6379", Description: "Redis port"},
			{Name: "JWT_SECRET", Required: true, Sensitive: true, Description: "JWT signing secret"},
			{Name: "JWT_EXPIRY", Default: "24h", Description: "JWT expiry duration"},
			{Name: "LOG_LEVEL", Default: "info", Description: "Log level (debug, info, warn, error)"},
		},
	},
}
```

---

## Environment Variables

### API

| Variable | Description | Default |
|----------|-------------|---------|
| `GDSNIP_PORT` | API server port | `8080` |
| `GDSNIP_DATABASE_URL` | PostgreSQL connection string | (required) |
| `GDSNIP_JWT_SECRET` | Secret for signing JWTs | (required) |
| `GDSNIP_JWT_EXPIRY` | Token expiration time | `24h` |

### CLI

| Variable | Description | Default |
|----------|-------------|---------|
| `GDSNIP_API_URL` | API base URL | `http://localhost:8080/v1` |

---

## Success Metrics

### Demo Checklist

- [ ] User can register and login via CLI (interactive TUI)
- [ ] User can search for snippets
- [ ] User can pull a snippet with variables
- [ ] Missing variables prompt interactive input
- [ ] Rendered output is displayed with styling
- [ ] User can save output to file
- [ ] User can push a new snippet
- [ ] User can view snippet info
- [ ] 10+ seed snippets available
- [ ] API is deployed and accessible
- [ ] CLI binary is distributable

### Demo Script

1. **Hook (30 sec)**: "How many times have you copy-pasted a Docker Compose file from Confluence?"
2. **Problem (1 min)**: Show the pain — searching Confluence, finding outdated template, manual find-replace
3. **Solution (1 min)**: `gdsnip pull docker-pg --DB_PASSWORD=secret`
4. **Demo (3 min)**:
   - `gdsnip search kafka`
   - `gdsnip info payments/kafka-setup`
   - `gdsnip pull payments/kafka-setup --BROKER_URL=localhost:9092 -o docker-compose.yml`
   - `gdsnip push -f my-template.yml --name my-service --public`
5. **Impact (1 min)**: "No more stale Confluence pages. One command. Done."
6. **Future (30 sec)**: Teams, versioning, VS Code extension

---

## References

- [Go CLI Best Practices](https://clig.dev/)
- [Cobra Documentation](https://cobra.dev/)
- [Bubbletea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Echo Framework](https://echo.labstack.com/)
- [Koanf Documentation](https://github.com/knadh/koanf)
- [Task Documentation](https://taskfile.dev/)
- [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html)