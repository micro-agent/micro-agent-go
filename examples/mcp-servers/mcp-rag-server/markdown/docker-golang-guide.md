# Dockerizing Go Applications with Docker Compose

## Table of Contents

1. [Introduction](#1-introduction)
2. [Basic Go Application Setup](#2-basic-go-application-setup)
3. [Creating a Basic Dockerfile](#3-creating-a-basic-dockerfile)
4. [Multi-stage Docker Build](#4-multi-stage-docker-build)
5. [Docker Compose Setup](#5-docker-compose-setup)
6. [Go Web Application with Database](#6-go-web-application-with-database)
7. [Advanced Docker Compose Configuration](#7-advanced-docker-compose-configuration)
8. [Development vs Production Environments](#8-development-vs-production-environments)
9. [Microservices Architecture](#9-microservices-architecture)
10. [Best Practices and Optimization](#10-best-practices-and-optimization)
11. [Debugging and Troubleshooting](#11-debugging-and-troubleshooting)
12. [Deployment Strategies](#12-deployment-strategies)

---

## 1. Introduction

Docker provides an excellent way to containerize Go applications, ensuring consistency across development, testing, and production environments. This guide covers everything from basic containerization to complex multi-service architectures using Docker Compose.

### Benefits of Dockerizing Go Applications

- **Consistent Environment**: Same runtime across all environments
- **Easy Deployment**: Package application with all dependencies
- **Scalability**: Easy horizontal scaling with container orchestration
- **Isolation**: Applications run in isolated containers
- **Development Efficiency**: Quick setup for new developers

---

## 2. Basic Go Application Setup

Let's start with a simple Go application structure:

```
my-go-app/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   └── health.go
│   └── server/
│       └── server.go
├── pkg/
│   └── config/
│       └── config.go
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
└── README.md
```

### Sample Go Application

**go.mod**:
```go
module my-go-app

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/joho/godotenv v1.5.1
    gorm.io/driver/postgres v1.5.2
    gorm.io/gorm v1.25.4
)
```

**cmd/api/main.go**:
```go
package main

import (
    "log"
    "my-go-app/internal/server"
    "my-go-app/pkg/config"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    srv := server.New(cfg)
    if err := srv.Start(); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

**internal/server/server.go**:
```go
package server

import (
    "fmt"
    "my-go-app/internal/handlers"
    "my-go-app/pkg/config"

    "github.com/gin-gonic/gin"
)

type Server struct {
    config *config.Config
    router *gin.Engine
}

func New(cfg *config.Config) *Server {
    return &Server{
        config: cfg,
        router: gin.Default(),
    }
}

func (s *Server) Start() error {
    s.setupRoutes()
    return s.router.Run(fmt.Sprintf(":%s", s.config.Port))
}

func (s *Server) setupRoutes() {
    api := s.router.Group("/api/v1")
    {
        api.GET("/health", handlers.HealthCheck)
        api.GET("/users", handlers.GetUsers)
        api.POST("/users", handlers.CreateUser)
    }
}
```

**internal/handlers/health.go**:
```go
package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":    "healthy",
        "timestamp": time.Now().UTC(),
        "service":   "my-go-app",
    })
}

func GetUsers(c *gin.Context) {
    // Mock users data
    users := []map[string]interface{}{
        {"id": 1, "name": "John Doe", "email": "john@example.com"},
        {"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
    }
    c.JSON(http.StatusOK, gin.H{"users": users})
}

func CreateUser(c *gin.Context) {
    var user map[string]interface{}
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user["id"] = 3 // Mock ID generation
    c.JSON(http.StatusCreated, gin.H{"user": user})
}
```

**pkg/config/config.go**:
```go
package config

import (
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    DatabaseURL string
    Environment string
}

func Load() (*Config, error) {
    // Load .env file if it exists
    godotenv.Load()

    return &Config{
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost/mydb?sslmode=disable"),
        Environment: getEnv("ENVIRONMENT", "development"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

---

## 3. Creating a Basic Dockerfile

### Basic Dockerfile

```dockerfile
# Basic Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main cmd/api/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN adduser -D -s /bin/sh appuser

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to appuser
USER appuser

EXPOSE 8080

CMD ["./main"]
```

### .dockerignore

```dockerignore
# Git
.git
.gitignore

# Documentation
README.md
*.md

# Docker
Dockerfile*
docker-compose*.yml

# Development files
.env
.env.local
.vscode/
.idea/

# Build artifacts
main
*.exe

# Test files
*_test.go

# Temporary files
tmp/
temp/
```

### Building and Running

```bash
# Build the Docker image
docker build -t my-go-app:latest .

# Run the container
docker run -p 8080:8080 my-go-app:latest

# Test the application
curl http://localhost:8080/api/v1/health
```

---

## 4. Multi-stage Docker Build

### Optimized Multi-stage Dockerfile

```dockerfile
# Multi-stage Dockerfile for Go application
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates (needed for private repos and SSL)
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser for security
RUN adduser -D -g '' appuser

WORKDIR /build

# Copy go mod files for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/api/main.go

# Final stage - minimal image
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /build/main /app/main

# Use an unprivileged user
USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/main"]
```

### Build with Build Args

```dockerfile
# Dockerfile with build arguments
FROM golang:1.21-alpine AS builder

ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with version information
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s \
    -X main.Version=${VERSION} \
    -X main.BuildTime=${BUILD_TIME} \
    -X main.GitCommit=${GIT_COMMIT}" \
    -o main cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D appuser

WORKDIR /app

COPY --from=builder /build/main .

USER appuser

EXPOSE 8080

CMD ["./main"]
```

Build with arguments:
```bash
docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --build-arg GIT_COMMIT=$(git rev-parse HEAD) \
  -t my-go-app:1.0.0 .
```

---

## 5. Docker Compose Setup

### Basic docker-compose.yml

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - ENVIRONMENT=development
    volumes:
      - .:/app
    depends_on:
      - postgres
      - redis
    networks:
      - app-network

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - app-network

volumes:
  postgres_data:
  redis_data:

networks:
  app-network:
    driver: bridge
```

### Environment-specific Configuration

**docker-compose.override.yml** (for development):
```yaml
version: '3.8'

services:
  app:
    build:
      target: builder
    environment:
      - GIN_MODE=debug
      - LOG_LEVEL=debug
    volumes:
      - .:/app
      - /app/vendor
    command: go run cmd/api/main.go
    
  postgres:
    environment:
      POSTGRES_DB: myapp_dev
    ports:
      - "5432:5432"
```

**docker-compose.prod.yml** (for production):
```yaml
version: '3.8'

services:
  app:
    build:
      target: production
    environment:
      - GIN_MODE=release
      - LOG_LEVEL=info
    restart: unless-stopped
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
        max_attempts: 3
    
  postgres:
    environment:
      POSTGRES_DB: myapp_prod
    restart: unless-stopped
    
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
```

---

## 6. Go Web Application with Database

### Database Integration

**internal/database/database.go**:
```go
package database

import (
    "fmt"
    "my-go-app/pkg/config"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type DB struct {
    *gorm.DB
}

func New(cfg *config.Config) (*DB, error) {
    db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    return &DB{db}, nil
}

func (db *DB) Migrate() error {
    return db.AutoMigrate(&User{})
}

type User struct {
    ID    uint   `gorm:"primaryKey" json:"id"`
    Name  string `gorm:"not null" json:"name"`
    Email string `gorm:"uniqueIndex;not null" json:"email"`
}
```

### Updated Handlers with Database

**internal/handlers/users.go**:
```go
package handlers

import (
    "net/http"
    "my-go-app/internal/database"

    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    db *database.DB
}

func NewUserHandler(db *database.DB) *UserHandler {
    return &UserHandler{db: db}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
    var users []database.User
    if err := h.db.Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var user database.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := h.db.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{"user": user})
}
```

### Database Initialization Script

**init.sql**:
```sql
-- Create database if not exists
CREATE DATABASE IF NOT EXISTS myapp;

-- Create user table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample data
INSERT INTO users (name, email) VALUES 
    ('John Doe', 'john@example.com'),
    ('Jane Smith', 'jane@example.com')
ON CONFLICT (email) DO NOTHING;
```

### Environment Variables

**.env**:
```bash
# Application
PORT=8080
ENVIRONMENT=development
GIN_MODE=debug

# Database
DATABASE_URL=postgres://user:password@postgres:5432/myapp?sslmode=disable
DB_HOST=postgres
DB_PORT=5432
DB_NAME=myapp
DB_USER=user
DB_PASSWORD=password

# Redis
REDIS_URL=redis://redis:6379/0

# Logging
LOG_LEVEL=debug
```

---

## 7. Advanced Docker Compose Configuration

### Complete Production Setup

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        VERSION: ${VERSION:-latest}
        BUILD_TIME: ${BUILD_TIME}
        GIT_COMMIT: ${GIT_COMMIT}
    image: my-go-app:${VERSION:-latest}
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - DATABASE_URL=postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable
      - REDIS_URL=redis://redis:6379/0
      - ENVIRONMENT=${ENVIRONMENT:-production}
      - GIN_MODE=release
    volumes:
      - ./logs:/app/logs
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 30s
    networks:
      - app-network
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_INITDB_ARGS: "--auth-host=scram-sha-256"
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
      - ./postgresql.conf:/etc/postgresql/postgresql.conf
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes --requirepass ${REDIS_PASSWORD}
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
      - ./redis.conf:/etc/redis/redis.conf
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "--raw", "incr", "ping"]
      interval: 10s
      timeout: 3s
      retries: 5
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/conf.d:/etc/nginx/conf.d
      - ./ssl:/etc/nginx/ssl
      - ./logs/nginx:/var/log/nginx
    depends_on:
      - app
    restart: unless-stopped
    networks:
      - app-network

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--web.enable-lifecycle'
    networks:
      - app-network

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/dashboards:/etc/grafana/provisioning/dashboards
      - ./grafana/datasources:/etc/grafana/provisioning/datasources
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_PASSWORD}
    networks:
      - app-network

volumes:
  postgres_data:
  redis_data:
  prometheus_data:
  grafana_data:

networks:
  app-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### Nginx Configuration

**nginx/nginx.conf**:
```nginx
events {
    worker_connections 1024;
}

http {
    upstream app {
        server app:8080;
    }

    server {
        listen 80;
        server_name localhost;

        location / {
            proxy_pass http://app;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /health {
            access_log off;
            proxy_pass http://app/api/v1/health;
        }
    }
}
```

---

## 8. Development vs Production Environments

### Development Environment

**docker-compose.dev.yml**:
```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
      - "2345:2345"  # Delve debugger port
    environment:
      - GIN_MODE=debug
      - LOG_LEVEL=debug
      - HOT_RELOAD=true
    volumes:
      - .:/app
      - go-modules:/go/pkg/mod
    command: >
      sh -c "
        go mod download &&
        go install github.com/cosmtrek/air@latest &&
        air
      "
    networks:
      - app-network

volumes:
  go-modules:
```

**Dockerfile.dev**:
```dockerfile
FROM golang:1.21-alpine

RUN apk add --no-cache git

WORKDIR /app

# Install air for hot reloading
RUN go install github.com/cosmtrek/air@latest

# Install delve for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080 2345

CMD ["air"]
```

**.air.toml** (for hot reloading):
```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/api"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
```

### Production Environment

**docker-compose.prod.yml**:
```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
      target: production
    environment:
      - GIN_MODE=release
      - LOG_LEVEL=info
    restart: unless-stopped
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
        max_attempts: 3
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  postgres:
    environment:
      POSTGRES_DB: myapp_prod
    volumes:
      - postgres_prod_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres_prod_data:
```

### Running Different Environments

```bash
# Development
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# With environment file
docker-compose --env-file .env.prod up -d
```

---

## 9. Microservices Architecture

### Multi-service Setup

```yaml
version: '3.8'

services:
  # User Service
  user-service:
    build:
      context: ./services/user-service
    ports:
      - "8081:8080"
    environment:
      - SERVICE_NAME=user-service
      - DATABASE_URL=postgres://user:password@postgres:5432/users
    depends_on:
      - postgres
      - consul
    networks:
      - microservices

  # Order Service
  order-service:
    build:
      context: ./services/order-service
    ports:
      - "8082:8080"
    environment:
      - SERVICE_NAME=order-service
      - DATABASE_URL=postgres://user:password@postgres:5432/orders
      - USER_SERVICE_URL=http://user-service:8080
    depends_on:
      - postgres
      - consul
    networks:
      - microservices

  # API Gateway
  api-gateway:
    build:
      context: ./services/api-gateway
    ports:
      - "8080:8080"
    environment:
      - USER_SERVICE_URL=http://user-service:8080
      - ORDER_SERVICE_URL=http://order-service:8080
    depends_on:
      - user-service
      - order-service
    networks:
      - microservices

  # Service Discovery
  consul:
    image: consul:latest
    ports:
      - "8500:8500"
    environment:
      - CONSUL_BIND_INTERFACE=eth0
    networks:
      - microservices

  # Message Broker
  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: password
    networks:
      - microservices

  # Monitoring
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"
      - "14268:14268"
    environment:
      - COLLECTOR_OTLP_ENABLED=true
    networks:
      - microservices

networks:
  microservices:
    driver: bridge
```

### Service Discovery Example

**internal/discovery/consul.go**:
```go
package discovery

import (
    "fmt"
    "strconv"

    "github.com/hashicorp/consul/api"
)

type ConsulClient struct {
    client *api.Client
}

func NewConsulClient(address string) (*ConsulClient, error) {
    config := api.DefaultConfig()
    config.Address = address
    
    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }
    
    return &ConsulClient{client: client}, nil
}

func (c *ConsulClient) RegisterService(name, address string, port int) error {
    registration := &api.AgentServiceRegistration{
        ID:      fmt.Sprintf("%s-%s-%d", name, address, port),
        Name:    name,
        Port:    port,
        Address: address,
        Check: &api.AgentServiceCheck{
            HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
            Timeout:                        "10s",
            Interval:                       "10s",
            DeregisterCriticalServiceAfter: "60s",
        },
    }
    
    return c.client.Agent().ServiceRegister(registration)
}

func (c *ConsulClient) DiscoverService(serviceName string) ([]string, error) {
    services, _, err := c.client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return nil, err
    }
    
    var addresses []string
    for _, service := range services {
        addr := service.Service.Address + ":" + strconv.Itoa(service.Service.Port)
        addresses = append(addresses, addr)
    }
    
    return addresses, nil
}
```

---

## 10. Best Practices and Optimization

### Optimized Dockerfile

```dockerfile
# Multi-stage build with caching optimization
FROM golang:1.21-alpine AS dependencies

WORKDIR /app

# Install dependencies in separate layer for better caching
COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy dependencies from previous stage
COPY --from=dependencies /go/pkg /go/pkg
COPY --from=dependencies /app/go.mod /app/go.sum ./

# Copy source code
COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/api/main.go

# Final minimal image
FROM scratch

# Import CA certificates and timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
COPY --from=builder /app/main /main

# Non-root user
USER 65534:65534

EXPOSE 8080

ENTRYPOINT ["/main"]
```

### Health Checks

**internal/health/checker.go**:
```go
package health

import (
    "context"
    "database/sql"
    "time"
)

type Checker struct {
    db *sql.DB
}

func New(db *sql.DB) *Checker {
    return &Checker{db: db}
}

func (c *Checker) Check(ctx context.Context) map[string]interface{} {
    status := map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now().UTC(),
        "checks": map[string]interface{}{},
    }

    // Database health check
    if err := c.checkDatabase(ctx); err != nil {
        status["status"] = "unhealthy"
        status["checks"].(map[string]interface{})["database"] = map[string]interface{}{
            "status": "unhealthy",
            "error":  err.Error(),
        }
    } else {
        status["checks"].(map[string]interface{})["database"] = map[string]interface{}{
            "status": "healthy",
        }
    }

    return status
}

func (c *Checker) checkDatabase(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    return c.db.PingContext(ctx)
}
```

### Logging Configuration

**pkg/logger/logger.go**:
```go
package logger

import (
    "io"
    "os"

    "github.com/sirupsen/logrus"
)

func New(level string, format string) *logrus.Logger {
    logger := logrus.New()
    
    // Set log level
    switch level {
    case "debug":
        logger.SetLevel(logrus.DebugLevel)
    case "info":
        logger.SetLevel(logrus.InfoLevel)
    case "warn":
        logger.SetLevel(logrus.WarnLevel)
    case "error":
        logger.SetLevel(logrus.ErrorLevel)
    default:
        logger.SetLevel(logrus.InfoLevel)
    }
    
    // Set formatter
    if format == "json" {
        logger.SetFormatter(&logrus.JSONFormatter{})
    } else {
        logger.SetFormatter(&logrus.TextFormatter{
            FullTimestamp: true,
        })
    }
    
    // Set output
    logger.SetOutput(io.MultiWriter(os.Stdout))
    
    return logger
}
```

### Metrics and Monitoring

**internal/metrics/prometheus.go**:
```go
package metrics

import (
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}

func PrometheusMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()))
        
        c.Next()
        
        timer.ObserveDuration()
        httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), string(rune(c.Writer.Status()))).Inc()
    })
}

func Handler() gin.HandlerFunc {
    h := promhttp.Handler()
    return gin.WrapH(h)
}
```

---

## 11. Debugging and Troubleshooting

### Debugging with Delve

**Dockerfile.debug**:
```dockerfile
FROM golang:1.21-alpine

RUN apk add --no-cache git

WORKDIR /app

# Install Delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with debug symbols
RUN go build -gcflags="all=-N -l" -o main cmd/api/main.go

EXPOSE 8080 2345

CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main"]
```

**docker-compose.debug.yml**:
```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.debug
    ports:
      - "8080:8080"
      - "2345:2345"
    volumes:
      - .:/app
    security_opt:
      - "apparmor:unconfined"
    cap_add:
      - SYS_PTRACE
```

### Logging and Debugging

```bash
# View logs
docker-compose logs app
docker-compose logs -f app  # Follow logs

# Execute commands in container
docker-compose exec app sh
docker-compose exec postgres psql -U user -d myapp

# Debug specific service
docker-compose run --rm app go run cmd/api/main.go

# Check container health
docker-compose ps
docker inspect <container_id>
```

### Common Issues and Solutions

1. **Port conflicts**:
```bash
# Check port usage
netstat -tulpn | grep :8080
lsof -i :8080

# Use different ports
docker-compose up --force-recreate
```

2. **Database connection issues**:
```bash
# Check database connectivity
docker-compose exec app ping postgres
docker-compose exec postgres pg_isready
```

3. **Memory issues**:
```bash
# Monitor container resources
docker stats
docker-compose exec app top
```

---

## 12. Deployment Strategies

### Production Deployment

**Makefile**:
```makefile
.PHONY: build push deploy

VERSION ?= $(shell git describe --tags --always --dirty)
REGISTRY ?= your-registry.com
IMAGE_NAME ?= my-go-app

build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ) \
		--build-arg GIT_COMMIT=$(shell git rev-parse HEAD) \
		-t $(REGISTRY)/$(IMAGE_NAME):$(VERSION) \
		-t $(REGISTRY)/$(IMAGE_NAME):latest .

push: build
	docker push $(REGISTRY)/$(IMAGE_NAME):$(VERSION)
	docker push $(REGISTRY)/$(IMAGE_NAME):latest

deploy:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

deploy-staging:
	docker-compose -f docker-compose.yml -f docker-compose.staging.yml up -d

clean:
	docker-compose down -v
	docker system prune -f
```

### CI/CD Pipeline (GitHub Actions)

**.github/workflows/deploy.yml**:
```yaml
name: Build and Deploy

on:
  push:
    branches: [main]
    tags: ['v*']

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
    - name: Checkout
      uses: actions/checkout@v3

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2

    - name: Log in to Container Registry
      uses: docker/login-action@v2
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v4
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}

    - name: Build and push Docker image
      uses: docker/build-push-action@v4
      with:
        context: .
        platforms: linux/amd64,linux/arm64
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
        build-args: |
          VERSION=${{ github.ref_name }}
          BUILD_TIME=${{ github.event.head_commit.timestamp }}
          GIT_COMMIT=${{ github.sha }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
    - name: Deploy to production
      uses: appleboy/ssh-action@v0.1.5
      with:
        host: ${{ secrets.HOST }}
        username: ${{ secrets.USERNAME }}
        key: ${{ secrets.PRIVATE_KEY }}
        script: |
          cd /opt/my-go-app
          docker-compose pull
          docker-compose up -d
          docker system prune -f
```

### Scaling with Docker Swarm

**docker-stack.yml**:
```yaml
version: '3.8'

services:
  app:
    image: my-go-app:latest
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:password@postgres:5432/myapp
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
        max_attempts: 3
      placement:
        constraints:
          - node.role == worker
    networks:
      - app-network

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    deploy:
      replicas: 1
      placement:
        constraints:
          - node.role == manager
    networks:
      - app-network

volumes:
  postgres_data:

networks:
  app-network:
    driver: overlay
```

Deploy to swarm:
```bash
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-stack.yml myapp

# Scale service
docker service scale myapp_app=5

# Update service
docker service update --image my-go-app:v2.0.0 myapp_app
```

---

## Summary

This comprehensive guide covers:

1. **Basic Setup**: Simple Go app containerization
2. **Multi-stage Builds**: Optimized Docker images
3. **Docker Compose**: Service orchestration
4. **Database Integration**: PostgreSQL and Redis setup
5. **Advanced Configuration**: Production-ready setups
6. **Development Environment**: Hot reloading and debugging
7. **Microservices**: Multi-service architecture
8. **Best Practices**: Security, logging, monitoring
9. **Debugging**: Tools and techniques
10. **Deployment**: CI/CD and scaling strategies

Key benefits of this approach:
- **Consistent environments** across development and production
- **Easy scaling** with container orchestration
- **Simplified deployment** with infrastructure as code
- **Better resource utilization** with containerization
- **Improved security** with isolated environments