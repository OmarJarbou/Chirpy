# Chirpy API Documentation

## Overview
Chirpy is a social media platform API that allows users to create, read, and manage short messages called "chirps". The API provides user authentication, chirp management, and webhook integration capabilities.

**Base URL:** `http://localhost:8080`  
**Port:** 8080

## Authentication
The API uses JWT (JSON Web Tokens) for authentication. Protected endpoints require a valid JWT token in the Authorization header.

**Format:** `Authorization: Bearer <token>`

## Endpoints

### Health Check

#### GET /api/healthz
Check if the API is running and ready to accept requests.

**Response:**
- **Status Code:** 200 OK
- **Content-Type:** `text/plain; charset=utf-8`
- **Body:** `OK`

---

### User Management

#### POST /api/users
Create a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
- **Status Code:** 201 Created
- **Content-Type:** `application/json`
- **Body:**
```json
{
  "id": "uuid-string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

**Error Responses:**
- **400 Bad Request:** Invalid password or hashing error
- **500 Internal Server Error:** Database error during user creation OR JSON decoding error

---

#### POST /api/login
Authenticate a user and receive access tokens.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
- **Status Code:** 200 OK
- **Content-Type:** `application/json`
- **Body:**
```json
{
  "id": "uuid-string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "jwt-access-token",
  "refresh_token": "refresh-token-string"
}
```

**Error Responses:**
- **400 Bad Request:** Error making access tokens or refresh tokens
- **401 Unauthorized:** Invalid email or password
- **500 Internal Server Error:** JSON decoding error

**Notes:**
- Access token expires in 1 hour (3600 seconds)
- Refresh token expires in 60 hours (216000 seconds)

---

#### PUT /api/users
Update user information (requires authentication).

**Headers:**
- `Authorization: Bearer <jwt-token>`

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "password": "newpassword123"
}
```

**Response:**
- **Status Code:** 200 OK
- **Content-Type:** `application/json`
- **Body:**
```json
{
  "id": "uuid-string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "newemail@example.com",
  "is_chirpy_red": false
}
```

**Error Responses:**
- **400 Bad Request:** Invalid password or hashing error
- **401 Unauthorized:** Invalid or missing JWT token
- **500 Internal Server Error:** Database error during update

---

#### POST /api/refresh
Refresh an expired access token using a valid refresh token.

**Headers:**
- `Authorization: Bearer <refresh-token>`

**Response:**
- **Status Code:** 200 OK
- **Content-Type:** `application/json`
- **Body:**
```json
{
  "token": "new-jwt-access-token"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid or missing refresh token or error making token
- **401 Unauthorized:** Invalid or expired or revoked refresh token

---

#### POST /api/revoke
Revoke a refresh token (requires authentication).

**Headers:**
- `Authorization: Bearer <refresh-token>`

**Response:**
- **Status Code:** 204 No Content

**Error Responses:**
- **400 Bad Request:** Invalid or missing refresh token

---

### Chirp Management

#### POST /api/chirps
Create a new chirp (requires authentication).

**Headers:**
- `Authorization: Bearer <jwt-token>`

**Request Body:**
```json
{
  "body": "This is my chirp message"
}
```

**Response:**
- **Status Code:** 201 Created
- **Content-Type:** `application/json`
- **Body:**
```json
{
  "id": "uuid-string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "This is my chirp message",
  "user_id": "user-uuid-string"
}
```

**Error Responses:**
- **400 Bad Request:** Chirp body exceeds 140 characters or contains banned words
- **401 Unauthorized:** Invalid or missing JWT token
- **500 Internal Server Error:** Database error during creation OR JSON decoding error

**Notes:**
- Chirp body must be 140 characters or less
- Banned words (case-insensitive): "kerfuffle", "sharbert", "fornax" (replaced with "****")

---

#### GET /api/chirps
Retrieve all chirps with optional filtering and sorting.

**Query Parameters:**
- `author_id` (optional): Filter chirps by specific user ID
- `sort` (optional): Sort order - "asc" (oldest first) or "desc" (newest first)

**Examples:**
- `GET /api/chirps` - Get all chirps
- `GET /api/chirps?author_id=uuid&sort=desc` - Get chirps by specific user, newest first

**Response:**
- **Status Code:** 200 OK
- **Content-Type:** `application/json`
- **Body:**
```json
[
  {
    "id": "uuid-string",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "body": "Chirp message",
    "user_id": "user-uuid-string"
  }
]
```

**Error Responses:**
- **401 Unauthorized:** Invalid UUID format for author_id
- **500 Internal Server Error:** Database error during retrieval

---

#### GET /api/chirps/{chirpID}
Retrieve a specific chirp by ID.

**Path Parameters:**
- `chirpID`: UUID of the chirp to retrieve

**Response:**
- **Status Code:** 200 OK
- **Content-Type:** `application/json`
- **Body:**
```json
{
  "id": "uuid-string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "Chirp message",
  "user_id": "user-uuid-string"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid UUID format
- **404 Not Found:** Chirp with specified ID not found
- **500 Internal Server Error:** Database error during retrieval

---

#### DELETE /api/chirps/{chirpID}
Delete a specific chirp (requires authentication and ownership).

**Headers:**
- `Authorization: Bearer <jwt-token>`

**Path Parameters:**
- `chirpID`: UUID of the chirp to delete

**Response:**
- **Status Code:** 204 No Content

**Error Responses:**
- **400 Bad Request:** Invalid UUID format
- **401 Unauthorized:** Invalid or missing JWT token
- **403 Forbidden:** User doesn't own the chirp
- **404 Not Found:** Chirp with specified ID not found
- **500 Internal Server Error:** Database error during deletion

---

### Webhook Integration

#### POST /api/polka/webhooks
Handle webhook events from external services (requires API key authentication).

**Headers:**
- `Authorization: ApiKey <api-key>`

**Request Body:**
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "user-uuid-string"
  }
}
```

**Supported Events:**
- `user.upgraded`: Upgrades a user to Chirpy Red status

**Response:**
- **Status Code:** 204 No Content

**Error Responses:**
- **400 Bad Request:** Invalid UUID format or unsupported event
- **401 Unauthorized:** Invalid or missing API key
- **404 Not Found:** User not found during upgrade
- **500 Internal Server Error:** JSON decoding error

---

### Admin Endpoints

#### GET /admin/metrics
View file server hit statistics (admin only).

**Response:**
- **Status Code:** 200 OK
- **Content-Type:** `text/html`
- **Body:** HTML page showing visit count

---

#### POST /admin/reset
Reset file server hit counter and delete all users (admin only).

**Response:**
- **Status Code:** 200 OK
- **Content-Type:** `application/json`
- **Body:**
```json
{
  "message": "File server hits has been reset successfully"
}
```

**Error Responses:**
- **500 Internal Server Error:** Database error during reset

---

### Static File Serving

#### GET /app/*
Serve static files from the application directory.

**Notes:**
- Serves files from the root directory with `/app/` prefix stripped
- Used for serving HTML, CSS, JavaScript, and other static assets
- Example: `/app/index.html` serves the root `index.html` file

---

## Data Models

### User
```json
{
  "id": "uuid-string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

### Chirp
```json
{
  "id": "uuid-string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "Chirp message content",
  "user_id": "user-uuid-string"
}
```

### Error Response
```json
{
  "error": "Error message description"
}
```

---

## Error Handling

All API endpoints return consistent error responses with the following structure:

- **Status Code:** Appropriate HTTP status code
- **Content-Type:** `application/json`
- **Body:** JSON object with `error` field containing error message

Common HTTP status codes:
- **200 OK:** Request successful
- **201 Created:** Resource created successfully
- **204 No Content:** Request successful, no response body
- **400 Bad Request:** Invalid request data
- **401 Unauthorized:** Authentication required or failed
- **403 Forbidden:** Access denied
- **404 Not Found:** Resource not found
- **500 Internal Server Error:** Server error

---

## Rate Limiting

Currently, the API does not implement rate limiting. All endpoints are accessible without request frequency restrictions.

---

## Security Features

1. **Password Hashing:** Passwords are hashed using bcrypt before storage
2. **JWT Authentication:** Secure token-based authentication
3. **Refresh Token Rotation:** Secure token refresh mechanism
4. **Content Filtering:** Automatic filtering of banned words in chirps
5. **API Key Protection:** Webhook endpoints protected by API keys
6. **Authorization Checks:** Users can only modify their own resources

---

## Dependencies

- **Go Version:** 1.24.5
- **Database:** PostgreSQL with sqlc for type-safe queries
- **Authentication:** JWT v5 for token management
- **Password Hashing:** bcrypt for secure password storage
- **UUID Generation:** Google UUID library
- **Environment Management:** godotenv for configuration

---

## Development Notes

- The API runs on port 8080 by default
- Database connection is configured via `DB_URL` environment variable
- JWT secret key is configured via `CHIRPY_SECRET_KEY` environment variable
- Polka webhook API key is configured via `POLKA_KEY` environment variable
- All timestamps are in ISO 8601 format (UTC)
- UUIDs are used for all ID fields to ensure uniqueness and security 
