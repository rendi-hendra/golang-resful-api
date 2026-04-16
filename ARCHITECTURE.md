# Architecture & Implementation Guide

## Table of Contents

- [System Architecture](#system-architecture)
- [Component Details](#component-details)
- [Data Flow](#data-flow)
- [Authentication Flow](#authentication-flow)
- [Error Handling Strategy](#error-handling-strategy)
- [Dependency Injection](#dependency-injection)
- [Database Design](#database-design)

---

## System Architecture

### Architectural Pattern: Clean Architecture

This project implements **Clean Architecture** principles with clear separation of concerns. The architecture ensures:

1. **Independence of Frameworks**: Core business logic doesn't depend on frameworks
2. **Testability**: Each layer can be tested independently
3. **Independence of Database**: Data access logic is abstracted
4. **Independence of Business Rules**: Logic doesn't depend on UI or external systems
5. **Independence of UI**: Can change UI without affecting business logic

### Layered Structure

```
┌─────────────────────────────────────────────────────┐
│                   Fiber HTTP Server                 │
└─────────────────────────────────────────────────────┘
                         ▲
                         │
┌─────────────────────────────────────────────────────┐
│            HTTP Delivery / Presentation Layer        │
│  - Controllers                                       │
│  - Routes                                           │
│  - Middleware (Auth, CORS, Logging)                │
│  - Request/Response Models                          │
└─────────────────────────────────────────────────────┘
                         ▲
                         │ Uses
                         │
┌─────────────────────────────────────────────────────┐
│              Business Logic / Usecase Layer          │
│  - User Service                                     │
│  - Validation                                       │
│  - Business Rules                                   │
│  - Token Management                                 │
└─────────────────────────────────────────────────────┘
                         ▲
                         │ Depends On
                         │
┌─────────────────────────────────────────────────────┐
│           Data Access / Repository Layer             │
│  - User Repository                                  │
│  - Query Building                                   │
│  - Data Persistence                                 │
└─────────────────────────────────────────────────────┘
                         ▲
                         │ Uses
                         │
┌─────────────────────────────────────────────────────┐
│        Database / Entity Layer (MySQL + GORM)       │
│  - User Entity                                      │
│  - Table Mappings                                   │
│  - Constraints                                      │
└─────────────────────────────────────────────────────┘
```

---

## Component Details

### 1. Delivery Layer (`internal/delivery/http/`)

**Responsibility**: Handle HTTP request/response and routing

#### UserController
```
Functions:
├── Register(ctx) → Create new user
├── Login(ctx) → Authenticate and return tokens
├── RefreshToken(ctx) → Generate new access token
├── Current(ctx) → Get authenticated user profile
└── Update(ctx) → Update user information
```

**Request Flow**:
1. HTTP request arrives
2. Middleware processes (auth, logging)
3. Route maps to controller
4. Controller parses request body
5. Controller calls usecase
6. Controller returns JSON response

#### Middleware
- **AuthMiddleware**: Validates JWT token, extracts user info
- **Logging**: Requests/responses logging
- **Error Handling**: Converts errors to HTTP responses

#### Route Configuration
- `POST /api/users` → Register (public)
- `POST /api/users/_login` → Login (public)
- `POST /refresh-token` → Refresh token (public)
- `GET /api/users/_current` → Get profile (protected)
- `PATCH /api/users/_current` → Update profile (protected)

### 2. Usecase Layer (`internal/usecase/`)

**Responsibility**: Implement business logic and validation

#### UserUseCase Interface
```go
type UserUseCase interface {
    Create(ctx context.Context, request *RegisterUserRequest) (*UserResponse, error)
    Login(ctx context.Context, request *LoginUserRequest) (*TokenResponse, error)
    Refresh(ctx context.Context, request *RefreshTokenRequest) (*TokenResponse, error)
    Current(ctx context.Context, request *GetUserRequest) (*UserResponse, error)
    Update(ctx context.Context, request *UpdateUserRequest) (*UserResponse, error)
}
```

#### Business Logic Flow

**Create User (Registration)**:
1. Validate input (email format, required fields)
2. Check if user ID already exists
3. Check if email already exists
4. Hash password
5. Save to database
6. Return user response

**Login**:
1. Validate credentials format
2. Find user by ID
3. Verify password hash
4. Generate access and refresh tokens
5. Store refresh token
6. Return tokens

**Refresh Token**:
1. Validate refresh token
2. Extract user ID from token
3. Verify token matches stored token
4. Generate new access and refresh tokens
5. Return tokens

**Get Current User**:
1. Get user ID from context
2. Retrieve user from database
3. Convert entity to response model
4. Return user data

**Update User**:
1. Get user ID from context
2. Validate input if provided
3. Update fields (name, email, password)
4. Hash password if changed
5. Save to database
6. Return updated user

### 3. Repository Layer (`internal/repository/`)

**Responsibility**: Abstract database operations

#### UserRepository Interface
```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

#### GORM Operations
- **Create**: `db.WithContext(ctx).Create(&user)`
- **Read**: `db.WithContext(ctx).Where("id = ?", id).First(&user)`
- **Update**: `db.WithContext(ctx).Save(&user)`
- **Delete**: `db.WithContext(ctx).Delete(&user)`

### 4. Entity Layer (`internal/entity/`)

**Responsibility**: Define database schema

```go
type User struct {
    ID        string `gorm:"column:id;primaryKey"`
    Password  string `gorm:"column:password"`
    Name      string `gorm:"column:name"`
    Email     string `gorm:"column:email;unique"`
    Token     string `gorm:"column:token"`
    CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
    UpdatedAt int64  `gorm:"column:updated_at;autoUpdateTime:milli"`
}
```

### 5. Configuration Layer (`internal/config/`)

**Responsibility**: Setup and initialize all services

#### Bootstrap Process
1. **Viper**: Load configuration from file
2. **Logrus**: Initialize logger
3. **Database**: Connect to MySQL
4. **Validator**: Setup struct validation
5. **Fiber**: Create HTTP server
6. **Routes**: Register endpoints
7. **Middleware**: Setup middleware chain

```go
config.Bootstrap(&BootstrapConfig{
    DB:       database,
    App:      fiberApp,
    Log:      logger,
    Validate: validator,
    Config:   viperConfig,
})
```

---

## Data Flow

### User Registration Flow

```
HTTP Request: POST /api/users
    │
    ├── Route Handler
    │   └── UserController.Register()
    │       │
    │       ├── Parse Request Body
    │       ├── Call UseCase.Create()
    │       │   │
    │       │   ├── Validate Input
    │       │   ├── Check Duplicate ID
    │       │   ├── Check Duplicate Email
    │       │   ├── Hash Password
    │       │   ├── Call Repository.Create()
    │       │   │   └── GORM: db.Create(&user)
    │       │   │       └── MySQL: INSERT INTO users
    │       │   │
    │       │   └── Return User Response
    │       │
    │       └── Return JSON Response
    │
HTTP Response: 200 OK {data: {id, name, email}}
```

### User Login Flow

```
HTTP Request: POST /api/users/_login
    │
    ├── Route Handler
    │   ├── Parse Request Body
    │   └── Call UseCase.Login()
    │       │
    │       ├── Validate Credentials Format
    │       ├── Call Repository.FindByID()
    │       │   └── GORM: db.Where("id = ?", id).First()
    │       │       └── MySQL: SELECT FROM users WHERE id = ?
    │       │
    │       ├── Verify Password Hash
    │       ├── Generate JWT Tokens
    │       └── Save Refresh Token
    │
HTTP Response: 200 OK {access_token, refresh_token}
```

### Protected Endpoint Flow

```
HTTP Request: GET /api/users/_current
    │
    ├── AuthMiddleware
    │   ├── Extract Token from Header
    │   ├── Validate Token Signature
    │   ├── Check Token Expiry
    │   └── Extract User ID (into context)
    │
    ├── Route Handler
    │   └── UserController.Current()
    │       │
    │       ├── Get User ID from Context
    │       └── Call Repository.FindByID()
    │           └── GORM: db.Where("id = ?", id).First()
    │
HTTP Response: 200 OK {data: {id, name, email}}
```

---

## Authentication Flow

### JWT Token Structure

**Access Token** (short-lived, ~1 hour):
- Contains user ID
- Used for API requests
- Includes expiration time

**Refresh Token** (long-lived, ~7 days):
- Stored in database
- Used to generate new access tokens
- Can be revoked

### Token Generation Process

```go
claims := jwt.MapClaims{
    "sub": userID,
    "exp": time.Now().Add(1 * time.Hour),
    "iat": time.Now(),
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, err := token.SignedString([]byte(secret))
```

### Token Validation Process

```
1. Extract token from "Authorization: Bearer <token>" header
2. Parse token using secret key
3. Verify signature
4. Check expiration time
5. Extract claims (user ID)
6. Return validated user info
```

---

## Error Handling Strategy

### Error Types

#### Validation Errors (400 Bad Request)
```
- Invalid email format
- Missing required fields
- Invalid data types
- Invalid password format
```

#### Authentication Errors (401 Unauthorized)
```
- Missing authorization header
- Invalid token format
- Expired token
- Wrong password
```

#### Conflict Errors (409 Conflict)
```
- Duplicate user ID
- Duplicate email address
```

#### Not Found Errors (404 Not Found)
```
- User not found
```

### Error Response Format

```json
{
  "error": "Descriptive error message"
}
```

### Error Handling Flow

```go
func (c *UserController) Register(ctx *fiber.Ctx) error {
    // 1. Parse request
    err := ctx.BodyParser(request)
    if err != nil {
        c.Log.Warnf("Parse error: %v", err)
        return fiber.ErrBadRequest
    }
    
    // 2. Call usecase
    response, err := c.UseCase.Create(ctx.UserContext(), request)
    if err != nil {
        c.Log.Warnf("Create error: %v", err)
        // 3. Return appropriate error
        return err
    }
    
    // 4. Return success response
    return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}
```

---

## Dependency Injection

The project uses constructor-based dependency injection:

```go
// Repository depends on DB
type UserRepository struct {
    DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{DB: db}
}

// UseCase depends on Repository and Logger
type UserUseCase struct {
    Repository UserRepository
    Log        *logrus.Logger
    Validator  *validator.Validate
}

func NewUserUseCase(
    repo UserRepository,
    log *logrus.Logger,
    validator *validator.Validate,
) *UserUseCase {
    return &UserUseCase{
        Repository: repo,
        Log:        log,
        Validator:  validator,
    }
}

// Controller depends on UseCase and Logger
type UserController struct {
    UseCase UserUseCase
    Log     *logrus.Logger
}

func NewUserController(useCase UserUseCase, logger *logrus.Logger) *UserController {
    return &UserController{
        UseCase: useCase,
        Log:     logger,
    }
}
```

**Benefits**:
- Easy to test (inject mocks)
- Loose coupling between components
- Clear dependencies
- Easy to manage lifecycle

---

## Database Design

### Users Table

```sql
CREATE TABLE users (
  id VARCHAR(255) PRIMARY KEY COMMENT 'User identifier',
  password VARCHAR(255) NOT NULL COMMENT 'Hashed password (bcrypt)',
  name VARCHAR(255) NOT NULL COMMENT 'User full name',
  email VARCHAR(255) NOT NULL UNIQUE COMMENT 'User email (unique)',
  token LONGTEXT COMMENT 'JWT refresh token',
  created_at BIGINT COMMENT 'Creation timestamp (milliseconds)',
  updated_at BIGINT COMMENT 'Last update timestamp (milliseconds)',
  
  INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### Indexes

- **Primary Key**: `id` - Fast lookup by user ID
- **Unique Index**: `email` - Prevent duplicate emails
- **Regular Index**: `email` - Improve WHERE clause performance

### Constraints

- **PRIMARY KEY**: `id` - Uniqueness and identification
- **UNIQUE**: `email` - Only one email per user
- **NOT NULL**: `password`, `name` - Required fields

### Design Rationale

1. **VARCHAR(255) for IDs**: Allows flexible ID formats (username, UUID, etc.)
2. **LONGTEXT for token**: Accommodates long JWT tokens
3. **BIGINT timestamps**: Millisecond precision for modern applications
4. **UTF8MB4 charset**: Support for Unicode characters and emojis
5. **InnoDB engine**: ACID compliance and reliability

---

## Performance Considerations

### Query Optimization

1. **Database Indexes**: Email field indexed for quick lookups
2. **Connection Pooling**: GORM manages connection pool
3. **WithContext**: Proper context propagation for timeouts

### Middleware Optimization

1. **Auth Middleware**: Only runs on protected routes
2. **Early Validation**: Fails fast on invalid input
3. **Structured Logging**: Efficient log output

### Caching Opportunities

Future enhancements:
- Redis cache for user profiles
- Token validation cache
- Query result caching

---

## Security Practices

### Password Security

- **Hashing**: bcrypt with salt
- **No Plain Text**: Passwords never stored plain
- **Bcrypt Config**: Appropriate cost factor (10-12)

### JWT Security

- **Secret Key**: Strong random key
- **Expiration**: Short-lived access tokens
- **Refresh Tokens**: Stored server-side
- **HTTPS Only**: Should only transmit over HTTPS

### Input Validation

- **Email Format**: Validated before database insert
- **Required Fields**: Checked before processing
- **Type Validation**: Struct validator for types

### SQL Injection Prevention

- **Parameterized Queries**: GORM uses prepared statements
- **No String Concatenation**: Always use parameters

---

## Testing Strategy

### Unit Testing

- **Repository Layer**: Mock database responses
- **Usecase Layer**: Mock repository calls
- **Controller Layer**: Mock usecase calls

### Integration Testing

- **Full HTTP Stack**: Test actual endpoints
- **Real Database**: Use test database
- **Data Isolation**: Truncate tables between tests

### Test Scenarios

For each endpoint, test:
1. **Happy Path**: Success case
2. **Validation Errors**: Invalid input
3. **Business Logic Errors**: Data conflicts
4. **Authentication/Authorization**: Permission checks
5. **Edge Cases**: Boundary conditions

---

## Deployment Considerations

### Environment Variables

Override config values via environment:
```bash
APP_DATABASE_HOST=prod-db.example.com
APP_JWT_SECRET=production-secret-key
APP_WEB_PORT=8080
```

### Production Checklist

- [ ] Use HTTPS only
- [ ] Set strong JWT secret
- [ ] Configure proper database credentials
- [ ] Enable logging at info level
- [ ] Set appropriate CORS headers
- [ ] Enable rate limiting
- [ ] Use environment-specific configs
- [ ] Monitor error rates
- [ ] Setup alerting

### Horizontal Scaling

- Stateless design allows multiple instances
- Use external database (already used)
- Use load balancer to distribute requests
- Share JWT secret across instances

---

**Last Updated**: April 16, 2026
