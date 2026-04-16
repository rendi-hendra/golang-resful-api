# Development Workflow Guide

Complete guide for developers working on this project.

## Table of Contents

- [Development Environment](#development-environment)
- [Local Setup](#local-setup)
- [Coding Standards](#coding-standards)
- [Adding Features](#adding-features)
- [Debugging](#debugging)
- [Performance Tips](#performance-tips)
- [Git Workflow](#git-workflow)

---

## Development Environment Setup

### Required Tools

```bash
# Go
brew install go  # macOS
# or download from https://golang.org/dl

# MySQL
brew install mysql  # macOS
# or Docker: docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:5.7

# VS Code Recommended Extensions
- Go (golang.go)
- REST Client (humao.rest-client)
- MySQL (cweijan.vscode-mysql)
- Docker (ms-azuretools.vscode-docker)
```

### IDE Setup

#### VS Code Go Extension
```json
{
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "go.useLanguageServer": true,
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

#### Vim/Neovim
```bash
# Install vim-go
git clone https://github.com/fatih/vim-go.git ~/.vim/plugged/vim-go

# Install gopls
go install github.com/golang/tools/gopls@latest
```

---

## Local Setup

### 1. Environment Configuration

Create `.env.local`:
```bash
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=3306
APP_DATABASE_USER=root
APP_DATABASE_PASSWORD=root
APP_DATABASE_NAME=resful_api
APP_JWT_SECRET=dev-secret-key-change-in-production
APP_WEB_PORT=8080
APP_LOG_LEVEL=info
```

### 2. Database Setup

```bash
# Create database
mysql -u root -p
> CREATE DATABASE resful_api;
> CREATE DATABASE resful_api_test;
> EXIT;

# Run migrations
cd db/migrations
# Apply up mig files manually or use migration tool

# Verify
mysql resful_api -e "SHOW TABLES;"
```

### 3. dependencies

```bash
# Download all dependencies
go mod download

# Check for security vulnerabilities
go list -json -m all | nancy sleuth

# Update to latest minor/patch versions
go get -u ./...
```

### 4. Start Development Server

```bash
# Simple run
go run cmd/web/main.go

# With hot reload
go install github.com/cosmtrek/air@latest
air

# With debugging (Delve)
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug cmd/web/main.go
```

---

## Coding Standards

### Go Style Guide

Follow [Effective Go](https://golang.org/doc/effective_go) and [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Project Conventions

#### Package Organization
```go
// Order imports
import (
    // Std library
    "context"
    "fmt"
    
    // Third party (alphabetically)
    "github.com/gofiber/fiber/v2"
    "github.com/sirupsen/logrus"
    
    // Internal
    "github.com/rendi-hendra/resful-api/internal/model"
)
```

#### Naming Conventions
```go
// Constants
const (
    MaxRetries = 3
    DefaultTimeout = 30 * time.Second
)

// Interfaces end with "er"
type Reader interface { ... }
type UserRepository interface { ... }

// Structs use PascalCase
type UserController struct { ... }

// Functions/variables use camelCase
func getUserByID(id string) { ... }
var userRepository UserRepository
```

#### Error Handling
```go
// Always handle errors
if err != nil {
    log.WithError(err).Error("Operation failed")
    return fmt.Errorf("operation failed: %w", err)
}

// Don't ignore errors
_ = ioutil.ReadFile("file.txt")  // BAD

// Wrap errors with context
if err := operation(); err != nil {
    return fmt.Errorf("operation: %w", err)  // GOOD
}
```

#### Documentation
```go
// Document all exported functions
// GetUser retrieves a user by ID from the database.
// Returns an error if the user is not found.
func GetUser(ctx context.Context, id string) (*User, error) {
    // ...
}

// Document complex logic
// We use exponential backoff to retry failed connections.
// This reduces server load during outages.
for i := 0; i < maxRetries; i++ {
    // ...
}
```

#### Testing
```go
// Test function naming
func TestGetUser_Success(t *testing.T) { ... }
func TestGetUser_NotFound(t *testing.T) { ... }

// Use table-driven tests
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        valid   bool
    }{
        {"valid", "user@example.com", true},
        {"invalid", "invalid", false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test code
        })
    }
}
```

---

## Adding Features

### Step-by-Step: Adding a New Endpoint

**Example: Add "Get All Users" endpoint**

#### 1. Define Entity
```go
// internal/entity/user_entity.go - Already exists
type User struct {
    ID        string
    Name      string
    Email     string
    // ...
}
```

#### 2. Create Repository Method
```go
// internal/repository/user_repository_interface.go
type UserRepository interface {
    // ... existing methods
    GetAll(ctx context.Context) ([]*entity.User, error)
}

// internal/repository/user_repository.go
func (r *userRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
    var users []*entity.User
    if err := r.DB.WithContext(ctx).Find(&users).Error; err != nil {
        r.Log.WithError(err).Error("Failed to get all users")
        return nil, err
    }
    return users, nil
}
```

#### 3. Create Usecase Method
```go
// internal/usecase/user_usecase_interface.go
type UserUseCase interface {
    // ... existing methods
    GetAll(ctx context.Context) ([]*model.UserResponse, error)
}

// internal/usecase/user_usecase.go
func (u *userUseCase) GetAll(ctx context.Context) ([]*model.UserResponse, error) {
    users, err := u.Repository.GetAll(ctx)
    if err != nil {
        return nil, err
    }
    
    responses := make([]*model.UserResponse, len(users))
    for i, user := range users {
        responses[i] = u.converter.ToUserResponse(user)
    }
    return responses, nil
}
```

#### 4. Add Controller Method
```go
// internal/delivery/http/user_controller.go
func (c *UserController) GetAll(ctx *fiber.Ctx) error {
    responses, err := c.UseCase.GetAll(ctx.UserContext())
    if err != nil {
        c.Log.WithError(err).Error("Failed to get all users")
        return err
    }
    
    return ctx.JSON(model.WebResponse[[]*model.UserResponse]{
        Data: responses,
    })
}
```

#### 5. Register Route
```go
// internal/delivery/http/route/route.go
func (c *RouteConfig) SetupAuthRoute() {
    c.App.Use(c.AuthMiddleware)
    c.App.Get("/api/users", c.UserController.GetAll)  // New line
    c.App.Get("/api/users/_current", c.UserController.Current)
    c.App.Patch("/api/users/_current", c.UserController.Update)
}
```

#### 6. Write Tests
```go
// test/user_test.go
func TestGetAllUsers_Success(t *testing.T) {
    // Setup: Create multiple users
    // Execute: GET /api/users with auth
    // Assert: Status 200, returns all users
}

func TestGetAllUsers_Unauthorized(t *testing.T) {
    // Execute: GET /api/users without token
    // Assert: Status 401
}
```

---

## Debugging

### Using Print Debugging

```go
package main

import "log"

func myFunction() {
    log.Printf("Value: %+v\n", myVar)  // Print with format
}
```

### Using Structured Logging

```go
import "github.com/sirupsen/logrus"

log := logrus.New()
log.WithFields(logrus.Fields{
    "user_id": "123",
    "action": "login",
}).Info("User action")
```

### Using Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug application
dlv debug cmd/web/main.go

# In dlv prompt:
# (dlv) break main.main          - set breakpoint
# (dlv) continue                 - continue execution
# (dlv) next                     - next line
# (dlv) step                     - step into
# (dlv) print variableName       - print variable
# (dlv) quit                     - exit
```

### Database Debugging

```bash
# Connect to database
mysql -u root -p resful_api

# View data
SELECT * FROM users;

# Check query
EXPLAIN SELECT * FROM users WHERE email = 'test@example.com';

# Monitor queries (in real-time)
SHOW PROCESSLIST;
```

### Request/Response Debugging

Create `test.http` file (VS Code REST Client):
```http
### Register user
POST http://localhost:8080/api/users
Content-Type: application/json

{
  "id": "testuser",
  "password": "TestPass@123",
  "name": "Test User",
  "email": "test@example.com"
}

### Login
POST http://localhost:8080/api/users/_login
Content-Type: application/json

{
  "id": "testuser",
  "password": "TestPass@123"
}

### Get current user
GET http://localhost:8080/api/users/_current
Authorization: Bearer YOUR_ACCESS_TOKEN_HERE
```

---

## Performance Tips

### Optimization Techniques

#### 1. Database Query Optimization
```go
// SLOW: N+1 problem
users := []User{}
db.Find(&users)
for _, user := range users {
    db.Find(&user.Posts)  // Extra query per user
}

// FAST: Use preload
db.Preload("Posts").Find(&users)

// FAST: Use select specific columns
db.Select("id", "name", "email").Find(&users)
```

#### 2. Connection Pooling
```go
// GORM handles this automatically
// But you can configure it:
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
```

#### 3. Caching
```go
// Simple in-memory cache (not suitable for distributed systems)
var userCache = make(map[string]*User)

func GetCachedUser(id string) *User {
    if user, ok := userCache[id]; ok {
        return user
    }
    // Fetch from DB
    user := fetchFromDB(id)
    userCache[id] = user
    return user
}
```

#### 4. Benchmarking
```bash
# Find bottlenecks
go test -bench=. -benchmem ./internal/...

# Profile CPU
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Profile memory
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

---

## Git Workflow

### Branch Strategy (Git Flow)

```
main (production)
  ↑
  └─── develop (development)
        ↑
        ├─── feature/user-management
        ├─── bugfix/login-issue
        └─── hotfix/security-patch
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add user registration endpoint
fix: resolve JWT token validation error
docs: update API documentation
test: add integration tests for login
refactor: reorganize repository layer
style: format code with gofmt
chore: update dependencies
```

### Workflow Steps

```bash
# 1. Create feature branch from develop
git checkout develop
git pull origin develop
git checkout -b feature/my-feature

# 2. Make changes
# ... edit files ...

# 3. Stage and commit
git add .
git commit -m "feat: add new feature"

# 4. Push to remote
git push origin feature/my-feature

# 5. Create Pull Request on GitHub
# ... request review ...

# 6. Address review comments
git add .
git commit -m "fix: address review comments"
git push origin feature/my-feature

# 7. Merge to develop (after approval)
git checkout develop
git pull origin develop
git merge feature/my-feature
git push origin develop

# 8. Delete feature branch
git branch -d feature/my-feature
git push origin --delete feature/my-feature
```

### Pre-commit Checklist

Before pushing:

- [ ] Code follows project style guide
- [ ] Tests pass: `go test ./...`
- [ ] No linting errors: `go vet ./...`
- [ ] Commit message follows convention
- [ ] Changes are documented if needed
- [ ] Database migrations included if needed

---

## Common Development Tasks

### Adding a Validation Rule
```go
// internal/config/validator.go
validator.RegisterValidation("custom", func(fl validator.FieldLevel) bool {
    value := fl.Field().String()
    return len(value) > 5 && len(value) < 50
})

// Use in model
type User struct {
    Name string `validate:"custom"`
}
```

### Adding a New Middleware
```go
// internal/delivery/http/middleware/logging_middleware.go
func LoggingMiddleware(c *fiber.Ctx) error {
    start := time.Now()
    
    err := c.Next()
    
    duration := time.Since(start)
    log.Printf("%s %s took %v", c.Method(), c.Path(), duration)
    
    return err
}

// Register in route
c.App.Use(LoggingMiddleware)
```

### Database Migration Process
```bash
# Create migration files
# 1. Create up migration: db/migrations/20260416_add_field.up.sql
# 2. Create down migration: db/migrations/20260416_add_field.down.sql
# 3. Apply migrations (manual or using migration tool)
# 4. Update entity structs
# 5. Test with fresh database
```

### Environment-Specific Configuration
```bash
# Use different configs
APP_ENV=development go run cmd/web/main.go
APP_ENV=production go run cmd/web/main.go

# Load from environment variables
export APP_DATABASE_HOST=prod-db.example.com
export APP_JWT_SECRET=production-secret
go run cmd/web/main.go
```

---

## Deployment Preparation Checklist

- [ ] All tests pass
- [ ] No security vulnerabilities: `go list -json -m all | nancy sleuth`
- [ ] Code is formatted: `go fmt ./...`
- [ ] Code is vetted: `go vet ./...`
- [ ] Database migrations tested
- [ ] Environment variables documented
- [ ] Error handling is comprehensive
- [ ] Logging is appropriate
- [ ] Performance is acceptable
- [ ] Documentation is up to date

---

## Useful Go Commands Cheat Sheet

```bash
# Build & Run
go run cmd/web/main.go                 # Run directly
go build -o app cmd/web/main.go       # Build executable
go build -ldflags="-s -w" ...          # Strip debug info

# Testing
go test ./...                          # Run all tests
go test -v ./...                       # Verbose
go test -race ./...                    # Detect race conditions
go test -cover ./...                   # Coverage
go test -run TestName                  # Run specific test

# Code Quality
go fmt ./...                            # Format code
go vet ./...                            # Static analysis
go mod tidy                             # Clean up modules
go mod verify                           # Verify module integrity

# Documentation
go doc package/function                # View documentation
godoc -http=:6060                      # Start doc server

# Profiling
go test -cpuprofile=cpu.prof          # CPU profile
go test -memprofile=mem.prof          # Memory profile
go tool pprof cpu.prof                # Analyze profile

# Dependencies
go get -u ./...                         # Update dependencies
go get -u github.com/package           # Install specific package
go mod vendor                           # Create vendor directory
```

---

**Last Updated**: April 16, 2026
