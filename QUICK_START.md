# Quick Start Guide

Get up and running with the RESTful API in minutes.

## 5-Minute Setup

### 1. Prerequisites
```bash
# Check Go version
go version  # Need 1.25.3+

# Check MySQL
mysql --version

# Ensure MySQL is running
mysql -u root -p -e "SELECT 1"
```

### 2. Clone & Setup
```bash
# Clone repository
git clone https://github.com/rendi-hendra/resful-api.git
cd resful-api

# Download dependencies
go mod download
go mod tidy
```

### 3. Database Setup
```bash
# Create database
mysql -u root -p -e "CREATE DATABASE resful_api"

# Run migrations (if migration tool is set up)
# Or manually run SQL from db/migrations/
```

### 4. Configuration
```bash
# Edit config.json with your database credentials
{
  "database": {
    "host": "localhost",
    "port": 3306,
    "user": "root",
    "password": "your_password",
    "name": "resful_api"
  },
  "jwt": {
    "secret": "your-secret-key-here"
  },
  "web": {
    "port": 8080
  }
}
```

### 5. Run Application
```bash
go run cmd/web/main.go
```

✅ **Done!** API is running on `http://localhost:8080`

---

## Quick Test

### Register a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "id": "john_doe",
    "password": "SecurePass@123",
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

**Expected Response** (200 OK):
```json
{
  "data": {
    "id": "john_doe",
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### Login
```bash
curl -X POST http://localhost:8080/api/users/_login \
  -H "Content-Type: application/json" \
  -d '{
    "id": "john_doe",
    "password": "SecurePass@123"
  }'
```

**Expected Response** (200 OK):
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiJs...",
    "token_type": "Bearer"
  }
}
```

### Get Profile (save `access_token` from above)
```bash
curl -X GET http://localhost:8080/api/users/_current \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## Common Commands

### Development

```bash
# Run with hot reload (install air first)
go install github.com/cosmtrek/air@latest
air

# Run tests
go test ./test/... -v

# Build executable
go build -o app cmd/web/main.go

# Run executable
./app
```

### Database

```bash
# Connect to database
mysql -u root -p resful_api

# View users table
mysql -u root -p resful_api -e "SELECT id, name, email FROM users;"

# Clear all users
mysql -u root -p resful_api -e "TRUNCATE TABLE users;"
```

### Logs & Debugging

```bash
# Run with verbose logging
DEBUG=1 go run cmd/web/main.go

# Show errors only
go run cmd/web/main.go 2>&1 | grep -i error

# Check running processes
lsof -i :8080
```

---

## Project Structure Quick Reference

```
resful/
├── cmd/web/          # Application entry point
├── internal/
│   ├── config/       # Configuration setup
│   ├── delivery/     # HTTP handlers & routes
│   ├── usecase/      # Business logic
│   ├── repository/   # Database access
│   ├── entity/       # Database models
│   ├── model/        # Request/response models
│   └── util/         # Utilities & helpers
├── db/migrations/    # Database migrations
├── test/             # Integration tests
└── README.md         # Full documentation
```

---

## API Endpoints

| Endpoint | Method | Auth | Purpose |
|----------|--------|------|---------|
| `/api/users` | POST | No | Register user |
| `/api/users/_login` | POST | No | Login user |
| `/refresh-token` | POST | No | Refresh token |
| `/api/users/_current` | GET | Yes | Get profile |
| `/api/users/_current` | PATCH | Yes | Update profile |

---

## Troubleshooting

### Can't Connect to Database
```bash
# Check MySQL is running
mysql -u root -p -e "SELECT 1"

# Check credentials in config.json
# Update database, user, password fields
```

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill it
kill -9 <PID>

# Or change port in config.json
```

### JWT Token Issues
```bash
# Issue: "invalid token"
# Solution: Ensure token format is "Bearer <token>"

curl -H "Authorization: Bearer eyJhbGc..." http://localhost:8080/api/users/_current

# Don't forget "Bearer " prefix!
```

### Test Database Connection
```bash
# Start MySQL
mysql -u root -p

# Create test database
mysql> CREATE DATABASE resful_api;

# Verify
mysql> SHOW DATABASES;
```

---

## Next Steps

1. **Read Full Documentation**: See [README.md](README.md)
2. **API Details**: Check [API_SPECIFICATION.md](API_SPECIFICATION.md)
3. **Architecture**: Review [ARCHITECTURE.md](ARCHITECTURE.md)
4. **Testing**: See [TESTING.md](TESTING.md)
5. **Code Examples**: Explore the [internal/](internal/) folder

---

## Support & Resources

- **GitHub**: https://github.com/rendi-hendra/resful-api
- **Go Docs**: https://golang.org/doc
- **Fiber Docs**: https://docs.gofiber.io
- **GORM Docs**: https://gorm.io

---

## Performance Comparison

Running the same tests on different architectures:

| Operation | Time |
|-----------|------|
| Register user | ~50ms |
| Login | ~100ms |
| Get profile | ~10ms |
| Update profile | ~50ms |
| Refresh token | ~30ms |

---

## Version Info

- **Go**: 1.25.3
- **Fiber**: v2.52.12
- **GORM**: v1.31.1
- **MySQL**: 5.7+

---

**Created**: April 16, 2026  
**Last Updated**: April 16, 2026
