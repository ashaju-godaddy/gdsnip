package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ashaju-godaddy/gdsnip/internal/api/repository"
	"github.com/ashaju-godaddy/gdsnip/internal/api/service"
)

func main() {
	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://gdsnip:gdsnip@localhost:5432/gdsnip?sslmode=disable"
	}

	// Connect to database
	db, err := repository.NewDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("🌱 Seeding database...")
	fmt.Println()

	// Initialize repositories and services
	userRepo := repository.NewUserRepo(db)
	snippetRepo := repository.NewSnippetRepo(db)
	teamRepo := repository.NewTeamRepo(db)
	authService := service.NewAuthService(userRepo, "seed-secret-key", 168*3600*1000000000)
	snippetService := service.NewSnippetService(snippetRepo, userRepo, teamRepo)

	// Create demo user
	fmt.Println("Creating demo user...")
	user, _, err := authService.Register("demo@example.com", "demo", "password123")
	if err != nil {
		log.Fatalf("Failed to create demo user: %v", err)
	}
	fmt.Printf("✓ Created user: %s (email: %s)\n", user.Username, user.Email)
	fmt.Println()

	// Create snippets
	snippets := []struct {
		Name        string
		Description string
		Content     string
		Tags        []string
		Visibility  string
	}{
		{
			Name:        "Docker Compose PostgreSQL",
			Description: "PostgreSQL database with Docker Compose including persistent storage",
			Content: `version: "3.8"

services:
  postgres:
    image: postgres:{{PG_VERSION}}
    container_name: {{CONTAINER_NAME}}
    environment:
      POSTGRES_DB: {{DB_NAME}}
      POSTGRES_USER: {{DB_USER}}
      POSTGRES_PASSWORD: {{DB_PASSWORD}}
    ports:
      - "{{DB_PORT}}:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - {{NETWORK_NAME}}

volumes:
  postgres_data:

networks:
  {{NETWORK_NAME}}:`,
			Tags:       []string{"docker", "postgres", "database"},
			Visibility: "public",
		},
		{
			Name:        "Docker Compose Redis",
			Description: "Redis cache with Docker Compose for development",
			Content: `version: "3.8"

services:
  redis:
    image: redis:{{REDIS_VERSION}}
    container_name: {{CONTAINER_NAME}}
    command: redis-server --requirepass {{REDIS_PASSWORD}}
    ports:
      - "{{REDIS_PORT}}:6379"
    volumes:
      - redis_data:/data
    networks:
      - {{NETWORK_NAME}}

volumes:
  redis_data:

networks:
  {{NETWORK_NAME}}:`,
			Tags:       []string{"docker", "redis", "cache"},
			Visibility: "public",
		},
		{
			Name:        "Kubernetes Deployment",
			Description: "Basic Kubernetes deployment configuration",
			Content: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{APP_NAME}}
  namespace: {{NAMESPACE}}
  labels:
    app: {{APP_NAME}}
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
        image: {{IMAGE}}
        ports:
        - containerPort: {{CONTAINER_PORT}}
        env:
        - name: APP_ENV
          value: {{APP_ENV}}
        resources:
          requests:
            memory: "{{MEMORY_REQUEST}}"
            cpu: "{{CPU_REQUEST}}"
          limits:
            memory: "{{MEMORY_LIMIT}}"
            cpu: "{{CPU_LIMIT}}"`,
			Tags:       []string{"kubernetes", "k8s", "deployment"},
			Visibility: "public",
		},
		{
			Name:        "Kubernetes Service",
			Description: "Kubernetes service configuration with LoadBalancer",
			Content: `apiVersion: v1
kind: Service
metadata:
  name: {{SERVICE_NAME}}
  namespace: {{NAMESPACE}}
  labels:
    app: {{APP_NAME}}
spec:
  type: {{SERVICE_TYPE}}
  selector:
    app: {{APP_NAME}}
  ports:
  - name: http
    protocol: TCP
    port: {{SERVICE_PORT}}
    targetPort: {{TARGET_PORT}}`,
			Tags:       []string{"kubernetes", "k8s", "service"},
			Visibility: "public",
		},
		{
			Name:        "GitHub Actions CI Pipeline",
			Description: "Complete CI pipeline with linting, testing, and building",
			Content: `name: {{WORKFLOW_NAME}}

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

      - name: Set up {{LANGUAGE}}
        uses: actions/setup-{{LANGUAGE}}@v3
        with:
          {{LANGUAGE}}-version: '{{VERSION}}'

      - name: Install dependencies
        run: {{INSTALL_COMMAND}}

      - name: Run linter
        run: {{LINT_COMMAND}}

      - name: Run tests
        run: {{TEST_COMMAND}}

      - name: Build
        run: {{BUILD_COMMAND}}`,
			Tags:       []string{"github", "ci", "actions"},
			Visibility: "public",
		},
		{
			Name:        "Nginx Reverse Proxy",
			Description: "Nginx configuration for reverse proxy with SSL",
			Content: `server {
    listen 80;
    listen [::]:80;
    server_name {{DOMAIN}};
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name {{DOMAIN}};

    ssl_certificate {{SSL_CERT_PATH}};
    ssl_certificate_key {{SSL_KEY_PATH}};
    ssl_protocols TLSv1.2 TLSv1.3;

    location / {
        proxy_pass http://{{BACKEND_HOST}}:{{BACKEND_PORT}};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}`,
			Tags:       []string{"nginx", "proxy", "ssl"},
			Visibility: "public",
		},
		{
			Name:        "Dockerfile Node.js",
			Description: "Optimized Dockerfile for Node.js applications",
			Content: `FROM node:{{NODE_VERSION}}-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy source code
COPY . .

# Build application
RUN npm run build

# Production stage
FROM node:{{NODE_VERSION}}-alpine

WORKDIR /app

# Copy built application
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY package*.json ./

# Set environment
ENV NODE_ENV=production
ENV PORT={{PORT}}

# Expose port
EXPOSE {{PORT}}

# Start application
CMD ["node", "dist/index.js"]`,
			Tags:       []string{"docker", "nodejs", "dockerfile"},
			Visibility: "public",
		},
		{
			Name:        "Dockerfile Go",
			Description: "Multi-stage Dockerfile for Go applications",
			Content: `# Build stage
FROM golang:{{GO_VERSION}}-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o {{BINARY_NAME}} ./cmd/{{APP_NAME}}

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/{{BINARY_NAME}} .

# Expose port
EXPOSE {{PORT}}

# Run binary
CMD ["./{{BINARY_NAME}}"]`,
			Tags:       []string{"docker", "golang", "dockerfile"},
			Visibility: "public",
		},
		{
			Name:        "Terraform AWS EC2",
			Description: "Terraform configuration for AWS EC2 instance",
			Content: `terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "{{AWS_REGION}}"
}

resource "aws_instance" "{{INSTANCE_NAME}}" {
  ami           = "{{AMI_ID}}"
  instance_type = "{{INSTANCE_TYPE}}"
  key_name      = "{{KEY_NAME}}"

  vpc_security_group_ids = ["{{SECURITY_GROUP_ID}}"]
  subnet_id              = "{{SUBNET_ID}}"

  tags = {
    Name        = "{{INSTANCE_NAME}}"
    Environment = "{{ENVIRONMENT}}"
  }

  root_block_device {
    volume_size = {{VOLUME_SIZE}}
    volume_type = "gp3"
  }
}

output "instance_ip" {
  value = aws_instance.{{INSTANCE_NAME}}.public_ip
}`,
			Tags:       []string{"terraform", "aws", "ec2"},
			Visibility: "public",
		},
		{
			Name:        "Makefile Go Project",
			Description: "Comprehensive Makefile for Go projects",
			Content: `.PHONY: help build test clean run fmt lint install

# Variables
BINARY_NAME={{BINARY_NAME}}
GO_FILES=$(shell find . -name '*.go' -type f)
VERSION={{VERSION}}
BUILD_DIR=./bin
CMD_DIR=./cmd/{{APP_NAME}}

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags="-X main.version=$(VERSION)" $(CMD_DIR)

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

run: build ## Build and run the application
	@$(BUILD_DIR)/$(BINARY_NAME)

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download`,
			Tags:       []string{"makefile", "golang", "build"},
			Visibility: "public",
		},
		{
			Name:        "Environment Variables Template",
			Description: "Template for .env file with common environment variables",
			Content: `# Application
APP_NAME={{APP_NAME}}
APP_ENV={{APP_ENV}}
APP_PORT={{APP_PORT}}
APP_URL={{APP_URL}}

# Database
DB_HOST={{DB_HOST}}
DB_PORT={{DB_PORT}}
DB_NAME={{DB_NAME}}
DB_USER={{DB_USER}}
DB_PASSWORD={{DB_PASSWORD}}

# Redis
REDIS_HOST={{REDIS_HOST}}
REDIS_PORT={{REDIS_PORT}}
REDIS_PASSWORD={{REDIS_PASSWORD}}

# JWT
JWT_SECRET={{JWT_SECRET}}
JWT_EXPIRY={{JWT_EXPIRY}}

# Email
SMTP_HOST={{SMTP_HOST}}
SMTP_PORT={{SMTP_PORT}}
SMTP_USER={{SMTP_USER}}
SMTP_PASSWORD={{SMTP_PASSWORD}}

# AWS
AWS_ACCESS_KEY_ID={{AWS_ACCESS_KEY_ID}}
AWS_SECRET_ACCESS_KEY={{AWS_SECRET_ACCESS_KEY}}
AWS_REGION={{AWS_REGION}}`,
			Tags:       []string{"env", "config", "template"},
			Visibility: "public",
		},
		{
			Name:        "systemd Service Unit",
			Description: "systemd service file for running applications as services",
			Content: `[Unit]
Description={{SERVICE_DESCRIPTION}}
After=network.target

[Service]
Type={{SERVICE_TYPE}}
User={{SERVICE_USER}}
Group={{SERVICE_GROUP}}
WorkingDirectory={{WORKING_DIR}}
ExecStart={{EXEC_START}}
Restart=always
RestartSec=5

# Environment
Environment="{{ENV_VAR_NAME}}={{ENV_VAR_VALUE}}"

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier={{SERVICE_NAME}}

# Security
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target`,
			Tags:       []string{"systemd", "linux", "service"},
			Visibility: "public",
		},
	}

	fmt.Println("Creating snippets...")
	for i, s := range snippets {
		req := service.CreateRequest{
			Name:        s.Name,
			Description: s.Description,
			Content:     s.Content,
			Tags:        s.Tags,
			Visibility:  s.Visibility,
		}

		snippet, err := snippetService.Create(user.ID, req)
		if err != nil {
			log.Printf("Failed to create snippet '%s': %v", s.Name, err)
			continue
		}

		fmt.Printf("✓ [%2d/12] %s (%s/%s) - %d variable(s)\n",
			i+1, snippet.Name, snippet.Namespace, snippet.Slug, len(snippet.Variables))
	}

	fmt.Println()
	fmt.Println("🎉 Seeding complete!")
	fmt.Println()
	fmt.Println("Demo user credentials:")
	fmt.Println("  Email:    demo@example.com")
	fmt.Println("  Username: demo")
	fmt.Println("  Password: password123")
	fmt.Println()
	fmt.Println("To log in with the CLI:")
	fmt.Println("  gdsnip auth login")
}
