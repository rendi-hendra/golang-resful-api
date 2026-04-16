# RESTful API - User Management System

A modern RESTful API built with Go, featuring user authentication, user management, and JWT-based token handling. This project demonstrates clean architecture principles with layered separation of concerns.

## Table of Contents

- [Project Overview](#project-overview)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Configuration](#configuration)
- [Testing](#testing)
- [Development Guidelines](#development-guidelines)

---

## Project Overview

This is a production-ready RESTful API for user management with the following features:

- **User Registration**: Create new user accounts with unique email and ID
- **User Authentication**: Login with credentials to receive JWT tokens
- **Token Refresh**: Refresh expired access tokens using refresh tokens
- **User Profile Management**: Retrieve and update user information
- **JWT Security**: Secure endpoints with Bearer token authentication
- **Input Validation**: Comprehensive validation for all API requests
- **Error Handling**: Standardized error responses across all endpoints

## Tech Stack

- **Language**: Go 1.25.3
- **Web Framework**: Fiber v2.52
- **Database**: MySQL with GORM ORM
- **Authentication**: JWT (golang-jwt v5)
- **Validation**: go-playground/validator v10
- **Configuration**: Viper
- **Logging**: Logrus
- **Cryptography**: golang.org/x/crypto
- **Testing**: Testify

## Project Structure

```
├── cmd/
│   └── web/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/                    # Configuration management
│   │   ├── app.go                 # App setup
│   │   ├── fiber.go               # Fiber framework config
│   │   ├── gorm.go                # Database config
│   │   ├── logrus.go              # Logger config
│   │   ├── validator.go           # Validator config
│   │   └── viper.go               # Config file management
│   ├── delivery/
│   │   └── http/                  # HTTP handlers
│   │       ├── user_controller.go # User API endpoints
│   │       ├── middleware/
│   │       │   └── auth_middleware.go  # JWT authentication
│   │       └── route/
│   │           └── route.go       # Route definitions
│   ├── entity/
│   │   └── user_entity.go         # Database entity models
│   ├── model/                     # Request/response models
│   │   ├── auth.go                # Auth models
│   │   ├── model.go               # Generic models
│   │   ├── user_model.go          # User models
│   │   └── converter/
│   │       └── user_converter.go  # Entity to model conversion
│   ├── repository/                # Data access layer
│   │   ├── repository.go          # Base repository
│   │   ├── user_repository.go     # User database operations
│   │   └── user_repository_interface.go  # Repository interface
│   ├── usecase/                   # Business logic layer
│   │   ├── user_usecase.go        # User service logic
│   │   └── user_usecase_interface.go    # Usecase interface
│   └── util/                      # Utilities
│       ├── mailer.go              # Email service
│       ├── notifier.go            # Notification service
│       ├── token_manager.go       # JWT token operations
│       └── token_util.go          # Token utilities
├── db/
│   └── migrations/                # Database migrations
│       ├── 20260329163349_create_table_users.*
│       └── 20260402014500_add_email_to_users.*
├── test/
│   ├── user_test.go               # Endpoint integration tests
│   └── mailer_test.go             # Mailer unit tests
├── config.json                    # Default configuration
├── go.mod                         # Go module definition
└── README.md                      # This file
```

## Architecture

This project follows a **layered architecture pattern** with clear separation of concerns:

```
┌─────────────────────────────────────────┐
│      HTTP Delivery Layer                │
│  (Controllers, Routes, Middleware)      │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│      Business Logic Layer (Usecase)     │
│  (Service, Validation, Business Rules)  │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│      Repository Layer (Data Access)     │
│  (Database Operations, CRUD)            │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│      Database Layer (MySQL)             │
│  (Entity Models, Persistence)           │
└─────────────────────────────────────────┘
```

### Key Components

1. **Delivery Layer** (`internal/delivery/http/`)
   - Handles HTTP request/response
   - JWT authentication middleware
   - Route management

2. **Use Case Layer** (`internal/usecase/`)
   - Implements business logic
   - Performs validation
   - Orchestrates repository operations

3. **Repository Layer** (`internal/repository/`)
   - Abstracts database operations
   - Implements CRUD operations
   - Manages entity persistence

4. **Entity Layer** (`internal/entity/`)
   - Defines database schema structures
   - Maps to database tables

5. **Configuration Layer** (`internal/config/`)
   - Centralized configuration management
   - Framework initialization
   - Service bootstrap

## Quick Start

### Prerequisites

- Go 1.25.3 or higher
- MySQL 5.7 or higher
- Environment with Web server capabilities

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/rendi-hendra/resful-api.git
   cd resful-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Configure environment**
   - Copy and modify `config.json`:
   ```json
   {
     "database": {
       "host": "localhost",
       "port": 3306,
       "user": "root",
       "password": "your_password",
       "name": "resful_api_db"
     },
     "jwt": {
       "secret": "your_jwt_secret_key"
     },
     "web": {
       "port": 8080
     }
   }
   ```

4. **Run database migrations**
   ```bash
   go run cmd/web/main.go migrate
   ```

5. **Start the application**
   ```bash
   go run cmd/web/main.go
   ```

   The API will be available at `http://localhost:8080`

---

## API Documentation

### Base URL
```
http://localhost:8080
```

### Response Format

All responses follow a standardized format:

**Success Response (200 OK)**:
```json
{
  "data": {
    "id": "user_id",
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

**Error Response**:
```json
{
  "error": "Error message description"
}
```

### Endpoints

#### 1. Register User
**POST** `/api/users`

Create a new user account.

**Request Body**:
```json
{
  "id": "john_doe",
  "password": "SecurePassword123",
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Response** (200 OK):
```json
{
  "data": {
    "id": "john_doe",
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

**Error Codes**:
- `400 Bad Request` - Invalid input or missing required fields
- `409 Conflict` - User ID or email already exists

---

#### 2. Login User
**POST** `/api/users/_login`

Authenticate user and receive JWT tokens.

**Request Body**:
```json
{
  "id": "john_doe",
  "password": "SecurePassword123"
}
```

**Response** (200 OK):
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer"
  }
}
```

**Error Codes**:
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Wrong password
- `404 Not Found` - User not found

---

#### 3. Refresh Token
**POST** `/refresh-token`

Generate a new access token using refresh token.

**Headers**:
```
Authorization: Bearer <refresh_token>
```

**Response** (200 OK):
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer"
  }
}
```

**Error Codes**:
- `400 Bad Request` - Invalid token format
- `401 Unauthorized` - Invalid or expired token

---

#### 4. Get Current User
**GET** `/api/users/_current`

Retrieve profile of authenticated user.

**Headers**:
```
Authorization: Bearer <access_token>
```

**Response** (200 OK):
```json
{
  "data": {
    "id": "john_doe",
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

**Error Codes**:
- `401 Unauthorized` - Missing or invalid token

---

#### 5. Update Current User
**PATCH** `/api/users/_current`

Update authenticated user's profile (name, email, password).

**Headers**:
```
Authorization: Bearer <access_token>
```

**Request Body** (all fields optional):
```json
{
  "name": "John Smith",
  "email": "john.smith@example.com",
  "password": "NewSecurePassword123"
}
```

**Response** (200 OK):
```json
{
  "data": {
    "id": "john_doe",
    "name": "John Smith",
    "email": "john.smith@example.com"
  }
}
```

**Error Codes**:
- `400 Bad Request` - Invalid input or validation failure
- `401 Unauthorized` - Missing or invalid token

---

## Database Schema

### Users Table

```sql
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

**Columns**:
- `id`: User's unique identifier (primary key)
- `password`: Hashed password (bcrypt)
- `name`: User's full name
- `email`: User's email address (unique)
- `token`: JWT refresh token
- `created_at`: Account creation timestamp (milliseconds)
- `updated_at`: Last update timestamp (milliseconds)

### Migrations

Two migration files manage schema evolution:

1. **20260329163349_create_table_users**: Initial schema creation
2. **20260402014500_add_email_to_users**: Added email field and constraints

---

## Configuration

Configuration is managed via `config.json` and can be extended with environment variables.

### Configuration Keys

```json
{
  "database": {
    "host": "localhost",
    "port": 3306,
    "user": "root",
    "password": "password",
    "name": "resful_api"
  },
  "jwt": {
    "secret": "your-secret-key",
    "access_token_lifetime": 3600,
    "refresh_token_lifetime": 604800
  },
  "web": {
    "port": 8080,
    "debug": false
  },
  "log": {
    "level": "info",
    "format": "json"
  }
}
```

### Environment Variable Override

Configuration values can be overridden via environment variables using the pattern:
```
APP_DATABASE_HOST=localhost
APP_JWT_SECRET=secret-key
APP_WEB_PORT=8080
```

---

## Testing

### Unit Test Scenarios

Comprehensive unit tests cover all API endpoints with the following scenarios:

#### User Registration Tests
- ✅ Successful registration with valid data
- ❌ Duplicate user ID error
- ❌ Invalid email format
- ❌ Missing required fields

#### User Login Tests
- ✅ Successful login with correct credentials
- ❌ Non-existent user
- ❌ Incorrect password
- ❌ Invalid request format

#### Token Refresh Tests
- ✅ Successful token refresh
- ❌ Invalid refresh token
- ❌ Missing token

#### Get Current User Tests
- ✅ Successful retrieval with valid token
- ❌ Invalid token
- ❌ Missing authorization header

#### User Update Tests
- ✅ Update name successfully
- ✅ Update password successfully
- ✅ Update email successfully
- ❌ Invalid email format
- ❌ Missing authorization

### Running Tests

```bash
# Run all tests
go test ./test/...

# Run with verbose output
go test -v ./test/...

# Run specific test
go test -run TestUserRegister ./test/...

# Run with coverage
go test -cover ./test/...
```

### Test Files Location

- `test/user_test.go` - User API endpoint integration tests
- `test/mailer_test.go` - Email service unit tests

### Test Setup

Each test:
1. Creates a fresh test instance
2. Initializes database connection
3. Truncates `users` table before each scenario
4. Executes isolated test cases
5. Validates HTTP status codes and response structures

---

## Development Guidelines

### Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use meaningful variable and function names
- Keep functions concise and focused
- Document exported functions with comments

### Adding New Features

1. **Add database entity** in `internal/entity/`
2. **Create repository interface** in `internal/repository/`
3. **Implement repository** in `internal/repository/`
4. **Create usecase interface** in `internal/usecase/`
5. **Implement business logic** in `internal/usecase/`
6. **Create request/response models** in `internal/model/`
7. **Add controller handler** in `internal/delivery/http/`
8. **Register routes** in `internal/delivery/http/route/`
9. **Create integration tests** in `test/`

### Git Workflow

```bash
# Create feature branch
git checkout -b feature/user-feature

# Make changes and commit
git add .
git commit -m "feat: add new feature"

# Push to remote
git push origin feature/user-feature

# Create pull request
```

### Error Handling

Always return appropriate HTTP status codes and error messages:

```go
if err != nil {
    log.Warnf("Failed operation: %+v", err)
    return fiber.NewError(fiber.StatusBadRequest, "Error message")
}
```

### Logging

Use structured logging with logrus:

```go
log.WithField("user_id", userID).Info("User created")
log.WithError(err).Error("Database error")
```

---

## Common Issues & Troubleshooting

### Database Connection Error
**Problem**: Cannot connect to MySQL
**Solution**:
1. Verify MySQL is running
2. Check database credentials in `config.json`
3. Ensure database exists

### Port Already in Use
**Problem**: Port 8080 already in use
**Solution**:
```bash
# Change port in config.json or use environment variable
APP_WEB_PORT=8081 go run cmd/web/main.go
```

### JWT Token Errors
**Problem**: "Invalid token" errors
**Solution**:
1. Ensure token is in Bearer format: `Authorization: Bearer <token>`
2. Check token expiration
3. Verify JWT secret matches in config

### Validation Errors
**Problem**: 400 Bad Request on valid-looking data
**Solution**:
1. Check field names match API documentation
2. Verify data types (string, not number for IDs)
3. Ensure email format is valid

---

## Performance Considerations

- **Connection Pooling**: Database connections are pooled via GORM
- **Middleware Chain**: Auth middleware runs only on protected routes
- **Logging**: Info level used; adjust to warning for production
- **Validation**: Performed early to fail fast

---

## Security Considerations

- **Password Hashing**: Bcrypt used for password storage
- **JWT Security**: Secret key should be strong and kept secure
- **JWT Expiry**: Access tokens expire after configured time
- **HTTPS**: Should be enforced in production
- **CORS**: Should be configured based on security requirements

---

## Future Enhancements

- [ ] Rate limiting
- [ ] API documentation with Swagger/OpenAPI
- [ ] Email verification for registration
- [ ] Password reset functionality
- [ ] User roles and permissions
- [ ] Audit logging
- [ ] Caching layer (Redis)
- [ ] Containerization (Docker)
- [ ] CI/CD pipeline
- [ ] API versioning

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add/update tests
5. Submit a pull request

---

## License

This project is proprietary. All rights reserved.

---

## Support

For issues and questions, please create an issue in the GitHub repository.

---

**Last Updated**: April 16, 2026  
**Version**: 1.0.0  
**Go Version**: 1.25.3
