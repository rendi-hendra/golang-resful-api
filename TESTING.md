# Testing Guide

## Table of Contents

- [Overview](#overview)
- [Testing Setup](#testing-setup)
- [Test Structure](#test-structure)
- [Running Tests](#running-tests)
- [Test Scenarios](#test-scenarios)
- [Best Practices](#best-practices)
- [Continuous Integration](#continuous-integration)

---

## Overview

This document provides comprehensive guidance for setting up and running tests for the RESTful API project. The project uses Go's built-in `testing` package along with the `testify` assertion library.

### Testing Pyramid

```
        ▲
       /|\       Unit Tests (Utilities, Converters)
      / | \      Quick execution, high coverage
     /  |  \
    ├───┼───┤    Integration Tests (Repository Layer)
    |   |   |    Database involved, moderate speed
    ├───┼───┤
    | Acceptance Tests (HTTP Endpoints) |
    | Full stack, slower but comprehensive  |
    └─────────────────────────────┘
```

---

## Testing Setup

### Prerequisites

1. **Go**: Version 1.25.3 or higher
2. **MySQL**: Test database instance
3. **Dependencies**: Already in `go.mod`

```bash
go get github.com/stretchr/testify
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
```

### Test Database Setup

Create a separate test database:

```sql
CREATE DATABASE resful_api_test;

USE resful_api_test;

CREATE TABLE users (
  id VARCHAR(255) PRIMARY KEY,
  password VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  token LONGTEXT,
  created_at BIGINT,
  updated_at BIGINT
);
```

### Test Configuration

Create `config.test.json` for test environment:

```json
{
  "database": {
    "host": "localhost",
    "port": 3306,
    "user": "root",
    "password": "your_password",
    "name": "resful_api_test"
  },
  "jwt": {
    "secret": "test-secret-key-12345"
  },
  "web": {
    "port": 8081
  },
  "log": {
    "level": "error"
  }
}
```

---

## Test Structure

### Directory Layout

```
test/
├── user_test.go         # HTTP integration tests
├── mailer_test.go       # Email service tests
└── README.md             # Test documentation
```

### Test File Naming

- Unit tests: `filename_test.go`
- Integration tests: `filename_integration_test.go`
- Helper tests: `helper_test.go`

### Test Function Naming

```go
// Format: Test<FunctionName><Scenario>
func TestRegisterUser_Success(t *testing.T) { ... }
func TestRegisterUser_DuplicateID(t *testing.T) { ... }
func TestLoginUser_InvalidCredentials(t *testing.T) { ... }
```

---

## Running Tests

### Basic Test Execution

```bash
# Run all tests in current directory
go test

# Run all tests with verbose output
go test -v

# Run all tests in a package
go test ./test/...

# Run with race detection (detect concurrent bugs)
go test -race

# Run with coverage
go test -cover
```

### Running Specific Tests

```bash
# Run single test
go test -run TestRegisterUser_Success

# Run tests matching pattern
go test -run TestRegisterUser

# Run excluding pattern
go test -run '^(?!.*Integration)'

# Run only integration tests
go test -run Integration
```

### Coverage Analysis

```bash
# Generate coverage report
go test -cover ./test/...

# Detailed coverage profile
go test -coverprofile=coverage.out ./test/...

# HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Check coverage percentage
go tool cover -func=coverage.out
```

### Benchmark Tests

```bash
# Run benchmarks
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkRegisterUser

# Save benchmark results
go test -bench=. -benchmem > bench.txt
```

---

## Test Scenarios

### 1. User Registration Tests

#### Scenario 1: Successful Registration

```go
func TestRegisterUser_Success(t *testing.T) {
    // Setup: Create fresh app instance
    // Execute: POST /api/users with valid data
    // Assert: 
    //   - Status code: 200 OK
    //   - Response contains user ID, name, email
    //   - User exists in database
}
```

**What to validate**:
- Response status is 200
- Response body contains correct user data
- User is persisted in database
- Email is unique
- Password is hashed (not plain text)

#### Scenario 2: Duplicate User ID

```go
func TestRegisterUser_DuplicateID(t *testing.T) {
    // Setup: 
    //   - Create first user with ID "test_user"
    //   - Clear data before test
    // Execute: Register second user with same ID
    // Assert: Status code 409 Conflict
}
```

#### Scenario 3: Invalid Email Format

```go
func TestRegisterUser_InvalidEmail(t *testing.T) {
    // Execute: POST with email "invalid-format"
    // Assert: Status code 400 Bad Request
}
```

#### Scenario 4: Missing Required Fields

```go
func TestRegisterUser_MissingFields(t *testing.T) {
    // Execute: POST with incomplete body
    // Assert: Status code 400 Bad Request
}
```

### 2. User Login Tests

#### Scenario 1: Successful Login

```go
func TestLoginUser_Success(t *testing.T) {
    // Setup:
    //   - Register user first
    //   - Clear data before test
    // Execute: POST /api/users/_login with correct credentials
    // Assert:
    //   - Status 200 OK
    //   - Response contains access_token
    //   - Response contains refresh_token
    //   - Tokens are valid JWT
}
```

#### Scenario 2: User Not Found

```go
func TestLoginUser_NotFound(t *testing.T) {
    // Execute: Login with non-existent user ID
    // Assert: Status 404 Not Found
}
```

#### Scenario 3: Wrong Password

```go
func TestLoginUser_WrongPassword(t *testing.T) {
    // Setup: Register user
    // Execute: Login with wrong password
    // Assert: Status 401 Unauthorized
}
```

#### Scenario 4: Invalid Request

```go
func TestLoginUser_InvalidRequest(t *testing.T) {
    // Execute: POST with malformed JSON
    // Assert: Status 400 Bad Request
}
```

### 3. Token Refresh Tests

#### Scenario 1: Successful Token Refresh

```go
func TestRefreshToken_Success(t *testing.T) {
    // Setup:
    //   - Register and login user to get refresh token
    // Execute: POST /refresh-token with valid refresh_token
    // Assert:
    //   - Status 200 OK
    //   - Response contains new access_token
    //   - New access_token is valid
}
```

#### Scenario 2: Invalid Refresh Token

```go
func TestRefreshToken_InvalidToken(t *testing.T) {
    // Execute: POST with invalid token
    // Assert: Status 401 Unauthorized
}
```

#### Scenario 3: Expired Refresh Token

```go
func TestRefreshToken_ExpiredToken(t *testing.T) {
    // Setup: Use manipulated expired token
    // Execute: POST /refresh-token
    // Assert: Status 401 Unauthorized
}
```

### 4. Get Current User Tests

#### Scenario 1: Successful Get Current User

```go
func TestGetCurrentUser_Success(t *testing.T) {
    // Setup:
    //   - Register user and get access token
    // Execute: GET /api/users/_current with Authorization header
    // Assert:
    //   - Status 200 OK
    //   - Response contains correct user data
    //   - Data matches registered user
}
```

#### Scenario 2: Missing Authorization Header

```go
func TestGetCurrentUser_MissingHeader(t *testing.T) {
    // Execute: GET without Authorization header
    // Assert: Status 401 Unauthorized
}
```

#### Scenario 3: Invalid Token

```go
func TestGetCurrentUser_InvalidToken(t *testing.T) {
    // Execute: GET with malformed token
    // Assert: Status 401 Unauthorized
}
```

### 5. Update Current User Tests

#### Scenario 1: Update Name Successfully

```go
func TestUpdateUser_UpdateName_Success(t *testing.T) {
    // Setup:
    //   - Register user and get token
    // Execute: PATCH /api/users/_current with new name
    // Assert:
    //   - Status 200 OK
    //   - Response shows updated name
    //   - Database reflects change
}
```

#### Scenario 2: Update Password Successfully

```go
func TestUpdateUser_UpdatePassword_Success(t *testing.T) {
    // Setup: Register user
    // Execute: PATCH with new password
    // Assert:
    //   - Status 200 OK
    //   - Can login with new password
    //   - Cannot login with old password
}
```

#### Scenario 3: Update Email Successfully

```go
func TestUpdateUser_UpdateEmail_Success(t *testing.T) {
    // Setup: Register user
    // Execute: PATCH with new email
    // Assert:
    //   - Status 200 OK
    //   - Email updated in database
    //   - Old email is now available
}
```

#### Scenario 4: Invalid Email Format

```go
func TestUpdateUser_InvalidEmail(t *testing.T) {
    // Execute: PATCH with email "invalid-format"
    // Assert: Status 400 Bad Request
}
```

#### Scenario 5: Duplicate Email

```go
func TestUpdateUser_DuplicateEmail(t *testing.T) {
    // Setup:
    //   - Register two users
    //   - Get token for first user
    // Execute: Try to update first user's email to second user's email
    // Assert: Status 409 Conflict
}
```

#### Scenario 6: Unauthorized Update

```go
func TestUpdateUser_Unauthorized(t *testing.T) {
    // Execute: PATCH without authorization header
    // Assert: Status 401 Unauthorized
}
```

---

## Test Implementation Example

### Basic Test Structure

```go
package test

import (
    "testing"
    "bytes"
    "encoding/json"
    "net/http"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "your-module/internal/model"
)

// TestRegisterUser_Success tests successful user registration
func TestRegisterUser_Success(t *testing.T) {
    // Setup
    app := setupTestApp(t)
    truncateUserTable(t)
    
    // Prepare request
    body := model.RegisterUserRequest{
        ID:       "testuser",
        Password: "TestPass@123",
        Name:     "Test User",
        Email:    "test@example.com",
    }
    jsonBody, _ := json.Marshal(body)
    
    // Create request
    req, _ := http.NewRequest(
        "POST",
        "http://localhost:8081/api/users",
        bytes.NewBuffer(jsonBody),
    )
    req.Header.Set("Content-Type", "application/json")
    
    // Execute
    resp, err := app.Test(req)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)
    
    // Parse response
    var response model.WebResponse[*model.UserResponse]
    json.NewDecoder(resp.Body).Decode(&response)
    
    assert.Equal(t, "testuser", response.Data.ID)
    assert.Equal(t, "Test User", response.Data.Name)
    assert.Equal(t, "test@example.com", response.Data.Email)
}

// Helper functions

func setupTestApp(t *testing.T) *fiber.App {
    // Initialize app with test config
    // Setup database
    // Setup routes
    return app
}

func truncateUserTable(t *testing.T) {
    // Connect to test database
    // Execute: TRUNCATE TABLE users
}
```

---

## Best Practices

### 1. Test Independence
- Each test should be independent
- Don't rely on test execution order
- Clean up data before each test
- Use helper functions for setup

### 2. Descriptive Names
```go
// Good
func TestRegisterUser_DuplicateEmail_Returns409(t *testing.T) { }

// Bad
func TestRegister(t *testing.T) { }
```

### 3. Arrange-Act-Assert Pattern
```go
func TestExample(t *testing.T) {
    // Arrange: Setup test data
    user := &User{ID: "test", Email: "test@example.com"}
    
    // Act: Execute the test
    err := repo.Create(context.Background(), user)
    
    // Assert: Verify results
    assert.NoError(t, err)
}
```

### 4. Use Table-Driven Tests
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        isValid bool
    }{
        {"valid email", "test@example.com", true},
        {"invalid email", "invalid", false},
        {"empty email", "", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := validateEmail(tt.email)
            assert.Equal(t, tt.isValid, result)
        })
    }
}
```

### 5. Mock External Dependencies
```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func TestUseCase_WithMock(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    usecase := NewUserUseCase(mockRepo)
    err := usecase.Register(context.Background(), &RegisterRequest{})
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### 6. Test Error Cases
```go
func TestLoginUser_WrongPassword_Returns401(t *testing.T) {
    // Setup test user
    // Attempt login with wrong password
    // Assert status is 401 and error is returned
}
```

### 7. Verify Side Effects
```go
func TestUpdateUser_VerifyPasswordHashed(t *testing.T) {
    // Update user password
    // Retrieve user from DB
    // Assert password is hashed and different from original
}
```

---

## Test Data Management

### Fixtures

Create reusable test data:

```go
func createTestUser(t *testing.T, db *gorm.DB) *User {
    user := &User{
        ID:       "testuser",
        Password: hashPassword("TestPass@123"),
        Name:     "Test User",
        Email:    "test@example.com",
    }
    
    if err := db.Create(user).Error; err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }
    
    return user
}
```

### Database Reset Strategy

```go
func resetDatabase(t *testing.T, db *gorm.DB) {
    // Truncate all test tables
    db.Exec("TRUNCATE TABLE users")
}

func TestExample(t *testing.T) {
    db := connectTestDB(t)
    defer resetDatabase(t, db)
    
    // Test code
}
```

---

## Continuous Integration

### GitHub Actions Example

Create `.github/workflows/test.yml`:

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      mysql:
        image: mysql:5.7
        env:
          MYSQL_ROOT_PASSWORD: root
        options: >-
          --health-cmd="mysqladmin ping"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=3

    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.25.3
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v2
```

---

## Troubleshooting Tests

### Common Issues

#### Connection Refused
```
Error: dial tcp 127.0.0.1:3306: connection refused
```
**Solution**: Ensure MySQL test instance is running
```bash
mysql -u root -p
CREATE DATABASE resful_api_test;
```

#### Port Already in Use
```
Error: listen tcp 127.0.0.1:8081: bind: address already in use
```
**Solution**: Kill existing process or use different port
```bash
lsof -i :8081
kill -9 <PID>
```

#### Flaky Tests
**Solution**: Add timeouts and better synchronization
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

#### Test Isolation Issues
**Solution**: Always reset database state
```go
func TestExample(t *testing.T) {
    truncateUserTable(t) // Always run first
    // Test code
}
```

---

## Performance Testing

### Load Testing with Apache Bench

```bash
# Test registration endpoint
ab -n 1000 -c 10 -T application/json \
   -p user.json \
   http://localhost:8081/api/users

# Test login endpoint
ab -n 100 -c 5 -T application/json \
   -p login.json \
   http://localhost:8081/api/users/_login
```

### Benchmarking Go Functions

```go
func BenchmarkPasswordHash(b *testing.B) {
    password := "TestPassword@123"
    
    for i := 0; i < b.N; i++ {
        hashPassword(password)
    }
}
```

---

**Last Updated**: April 16, 2026
