# API Specification & Examples

## Quick Reference

### Base URL
```
http://localhost:8080
```

### Authentication
All protected endpoints require JWT Bearer token in the `Authorization` header:
```
Authorization: Bearer <access_token>
```

---

## API Endpoints Summary

| Method | Endpoint | Protected | Description |
|--------|----------|-----------|-------------|
| POST | `/api/users` | No | Register new user |
| POST | `/api/users/_login` | No | Login user |
| POST | `/refresh-token` | No | Refresh access token |
| GET | `/api/users/_current` | Yes | Get current user profile |
| PATCH | `/api/users/_current` | Yes | Update current user |

---

## Detailed Endpoint Documentation

### 1. Register New User

**Endpoint**: `POST /api/users`

**Protected**: No

**Description**: Create a new user account with unique ID and email.

#### Request

**Headers**:
```
Content-Type: application/json
```

**Body**:
```json
{
  "id": "john_doe",
  "password": "SecurePass@123",
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Field Validations**:
- `id`: Required, unique, string (1-255 characters)
- `password`: Required, string (minimum 6 characters)
- `name`: Required, string (1-255 characters)
- `email`: Required, valid email format, unique

#### Response

**Success (200 OK)**:
```json
{
  "data": {
    "id": "john_doe",
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

**Error (400 Bad Request)** - Invalid input:
```json
{
  "error": "invalid request format or validation error"
}
```

**Error (409 Conflict)** - Duplicate user:
```json
{
  "error": "user id or email already exists"
}
```

#### Examples

**cURL**:
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

**JavaScript/Fetch**:
```javascript
const response = await fetch('http://localhost:8080/api/users', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    id: 'john_doe',
    password: 'SecurePass@123',
    name: 'John Doe',
    email: 'john@example.com'
  })
});

const data = await response.json();
console.log(data);
```

**Python**:
```python
import requests

url = 'http://localhost:8080/api/users'
payload = {
    'id': 'john_doe',
    'password': 'SecurePass@123',
    'name': 'John Doe',
    'email': 'john@example.com'
}

response = requests.post(url, json=payload)
print(response.json())
```

**Go**:
```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func main() {
    payload := map[string]string{
        "id":       "john_doe",
        "password": "SecurePass@123",
        "name":     "John Doe",
        "email":    "john@example.com",
    }
    
    data, _ := json.Marshal(payload)
    
    resp, _ := http.Post(
        "http://localhost:8080/api/users",
        "application/json",
        bytes.NewBuffer(data),
    )
    defer resp.Body.Close()
}
```

---

### 2. Login User

**Endpoint**: `POST /api/users/_login`

**Protected**: No

**Description**: Authenticate user and receive JWT tokens.

#### Request

**Headers**:
```
Content-Type: application/json
```

**Body**:
```json
{
  "id": "john_doe",
  "password": "SecurePass@123"
}
```

**Field Validations**:
- `id`: Required, string
- `password`: Required, string

#### Response

**Success (200 OK)**:
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJqb2huX2RvZSIsImV4cCI6MTcxMzI4NTQ4MCwiaWF0IjoxNzEzMjgxODgwfQ.1234567890",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJqb2huX2RvZSIsImV4cCI6MTcxMzg4Njc2MCwiaWF0IjoxNzEzMjgxODgwfQ.abcdefghij",
    "token_type": "Bearer"
  }
}
```

**Error (400 Bad Request)**:
```json
{
  "error": "invalid request format"
}
```

**Error (401 Unauthorized)** - Wrong password:
```json
{
  "error": "invalid credentials"
}
```

**Error (404 Not Found)**:
```json
{
  "error": "user not found"
}
```

#### Examples

**cURL**:
```bash
curl -X POST http://localhost:8080/api/users/_login \
  -H "Content-Type: application/json" \
  -d '{
    "id": "john_doe",
    "password": "SecurePass@123"
  }'
```

**JavaScript/Fetch**:
```javascript
const response = await fetch('http://localhost:8080/api/users/_login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    id: 'john_doe',
    password: 'SecurePass@123'
  })
});

const { data } = await response.json();
localStorage.setItem('accessToken', data.access_token);
localStorage.setItem('refreshToken', data.refresh_token);
```

---

### 3. Refresh Token

**Endpoint**: `POST /refresh-token`

**Protected**: No

**Description**: Generate new access token using refresh token.

#### Request

**Headers**:
```
Authorization: Bearer <refresh_token>
Content-Type: application/json
```

**Body**:
```json
{}
```

#### Response

**Success (200 OK)**:
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJqb2huX2RvZSIsImV4cCI6MTcxMzI4NTQ4MCwiaWF0IjoxNzEzMjgxODgwfQ.new_token_1234567890",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJqb2huX2RvZSIsImV4cCI6MTcxMzg4Njc2MCwiaWF0IjoxNzEzMjgxODgwfQ.new_refresh_token",
    "token_type": "Bearer"
  }
}
```

**Error (400 Bad Request)**:
```json
{
  "error": "invalid token format"
}
```

**Error (401 Unauthorized)**:
```json
{
  "error": "invalid or expired token"
}
```

#### Examples

**cURL**:
```bash
curl -X POST http://localhost:8080/refresh-token \
  -H "Authorization: Bearer YOUR_REFRESH_TOKEN" \
  -H "Content-Type: application/json"
```

**JavaScript/Fetch**:
```javascript
const refreshToken = localStorage.getItem('refreshToken');

const response = await fetch('http://localhost:8080/refresh-token', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${refreshToken}`,
    'Content-Type': 'application/json'
  }
});

const { data } = await response.json();
localStorage.setItem('accessToken', data.access_token);
```

---

### 4. Get Current User Profile

**Endpoint**: `GET /api/users/_current`

**Protected**: Yes (requires access token)

**Description**: Retrieve the profile of the authenticated user.

#### Request

**Headers**:
```
Authorization: Bearer <access_token>
```

#### Response

**Success (200 OK)**:
```json
{
  "data": {
    "id": "john_doe",
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

**Error (401 Unauthorized)**:
```json
{
  "error": "unauthorized or invalid token"
}
```

#### Examples

**cURL**:
```bash
curl -X GET http://localhost:8080/api/users/_current \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**JavaScript/Fetch**:
```javascript
const accessToken = localStorage.getItem('accessToken');

const response = await fetch('http://localhost:8080/api/users/_current', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${accessToken}`
  }
});

const { data } = await response.json();
console.log(data);
// {
//   id: 'john_doe',
//   name: 'John Doe',
//   email: 'john@example.com'
// }
```

**Python**:
```python
import requests

access_token = 'YOUR_ACCESS_TOKEN'
headers = {
    'Authorization': f'Bearer {access_token}'
}

response = requests.get(
    'http://localhost:8080/api/users/_current',
    headers=headers
)
print(response.json())
```

---

### 5. Update Current User

**Endpoint**: `PATCH /api/users/_current`

**Protected**: Yes (requires access token)

**Description**: Update the profile of the authenticated user.

#### Request

**Headers**:
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Body** (all fields optional):
```json
{
  "name": "John Smith",
  "email": "john.smith@example.com",
  "password": "NewSecurePass@456"
}
```

**Field Validations**:
- `name`: Optional, string (1-255 characters)
- `email`: Optional, valid email format, must be unique
- `password`: Optional, string (minimum 6 characters)

#### Response

**Success (200 OK)**:
```json
{
  "data": {
    "id": "john_doe",
    "name": "John Smith",
    "email": "john.smith@example.com"
  }
}
```

**Error (400 Bad Request)** - Invalid input:
```json
{
  "error": "invalid email format or validation error"
}
```

**Error (401 Unauthorized)**:
```json
{
  "error": "unauthorized or invalid token"
}
```

**Error (409 Conflict)** - Email already exists:
```json
{
  "error": "email already in use"
}
```

#### Examples

**Update Name**:
```bash
curl -X PATCH http://localhost:8080/api/users/_current \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith"
  }'
```

**Update Email**:
```bash
curl -X PATCH http://localhost:8080/api/users/_current \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.smith@example.com"
  }'
```

**Update Password**:
```bash
curl -X PATCH http://localhost:8080/api/users/_current \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "NewSecurePass@456"
  }'
```

**JavaScript/Fetch**:
```javascript
const accessToken = localStorage.getItem('accessToken');

const response = await fetch('http://localhost:8080/api/users/_current', {
  method: 'PATCH',
  headers: {
    'Authorization': `Bearer ${accessToken}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'John Smith',
    email: 'john.smith@example.com'
  })
});

const { data } = await response.json();
console.log('Updated user:', data);
```

---

## Complete Authentication Flow Example

### Step 1: Register New User

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "id": "alice123",
    "password": "SecurePass@789",
    "name": "Alice Johnson",
    "email": "alice@example.com"
  }'
```

Response:
```json
{
  "data": {
    "id": "alice123",
    "name": "Alice Johnson",
    "email": "alice@example.com"
  }
}
```

### Step 2: Login

```bash
curl -X POST http://localhost:8080/api/users/_login \
  -H "Content-Type: application/json" \
  -d '{
    "id": "alice123",
    "password": "SecurePass@789"
  }'
```

Response (save these tokens):
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiJs...",
    "token_type": "Bearer"
  }
}
```

### Step 3: Access Protected Endpoint

```bash
curl -X GET http://localhost:8080/api/users/_current \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

Response:
```json
{
  "data": {
    "id": "alice123",
    "name": "Alice Johnson",
    "email": "alice@example.com"
  }
}
```

### Step 4: Refresh Access Token (when expired)

```bash
curl -X POST http://localhost:8080/refresh-token \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiJs..."
```

Response:
```json
{
  "data": {
    "access_token": "new_access_token_here...",
    "refresh_token": "new_refresh_token_here...",
    "token_type": "Bearer"
  }
}
```

### Step 5: Update Profile

```bash
curl -X PATCH http://localhost:8080/api/users/_current \
  -H "Authorization: Bearer new_access_token_here..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Smith",
    "password": "NewPass@123"
  }'
```

Response:
```json
{
  "data": {
    "id": "alice123",
    "name": "Alice Smith",
    "email": "alice@example.com"
  }
}
```

---

## Error Handling

### HTTP Status Codes

| Code | Meaning | Example |
|------|---------|---------|
| 200 | OK | Successful request |
| 400 | Bad Request | Invalid input or validation error |
| 401 | Unauthorized | Missing or invalid authentication |
| 404 | Not Found | User/resource not found |
| 409 | Conflict | Resource already exists (duplicate) |
| 500 | Internal Server Error | Unexpected server error |

### Common Error Scenarios

#### Missing Required Fields
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "id": "bob123"
  }'
```

Response (400):
```json
{
  "error": "missing required fields"
}
```

#### Invalid Email Format
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "id": "bob123",
    "password": "pass123",
    "name": "Bob",
    "email": "invalid-email"
  }'
```

Response (400):
```json
{
  "error": "invalid email format"
}
```

#### Duplicate ID
```bash
# Registering with existing ID
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "id": "alice123",
    "password": "NewPass@123",
    "name": "Another Alice",
    "email": "another@example.com"
  }'
```

Response (409):
```json
{
  "error": "user id already exists"
}
```

#### Missing Token
```bash
curl -X GET http://localhost:8080/api/users/_current \
  -H "Authorization: "
```

Response (401):
```json
{
  "error": "missing or invalid authorization"
}
```

---

## Rate Limiting & Best Practices

### Recommendations

1. **Rate Limiting** (recommended to implement):
   - 100 requests per minute for login endpoint
   - 1000 requests per minute for other endpoints

2. **Token Management**:
   - Store tokens securely (HttpOnly cookies or secure storage)
   - Refresh tokens before expiry
   - Implement token revocation on logout

3. **Security**:
   - Always use HTTPS in production
   - Use strong JWT secrets (min 32 characters)
   - Implement CORS properly
   - Validate all inputs

4. **Error Handling**:
   - Don't expose internal details in error messages
   - Log errors server-side
   - Return generic error messages to clients

---

**Last Updated**: April 16, 2026
