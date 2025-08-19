# Chirpy - Social Media Platform API

A robust social media platform API built with Go that allows users to create, read, and manage short messages called "chirps". Built with modern Go practices, PostgreSQL, and JWT authentication.

## 🚀 Features

- **User Management**: Create, update, and manage user accounts
- **Authentication**: JWT-based authentication with refresh token support
- **Chirp Operations**: Create, read, and delete chirps (140 character limit)
- **Content Filtering**: Automatic filtering of banned words
- **Webhook Integration**: Support for external service integrations
- **Admin Tools**: Metrics and system management endpoints
- **Static File Serving**: Serve web application assets

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.24.5 or higher** - [Download Go](https://golang.org/dl/)
- **PostgreSQL 12 or higher** - [Download PostgreSQL](https://www.postgresql.org/download/)
- **Git** - [Download Git](https://git-scm.com/downloads)

## 🛠️ Installation & Setup

### Step 1: Clone the Repository

```bash
git clone https://github.com/yourusername/Chirpy.git
cd Chirpy
```

### Step 2: Install Go Dependencies

```bash
go mod download
```

### Step 3: Set Up PostgreSQL Database

1. **Create a new database:**
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE chirpy_db;

# Create user (optional but recommended)
CREATE USER chirpy_user WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE chirpy_db TO chirpy_user;

# Exit PostgreSQL
\q
```

2. **Install and run database migrations with Goose:**
```bash
# Install Goose migration tool
go install github.com/pressly/goose/v3/cmd/goose@latest

# Navigate to the sql/schema directory
cd sql/schema

# Run all migrations (this will execute 001_users.sql, 002_chirps.sql, etc.)
goose postgres "postgresql://chirpy_user:your_secure_password@localhost:5432/chirpy_db?sslmode=disable" up

# Verify migrations were applied successfully
goose postgres "postgresql://chirpy_user:your_secure_password@localhost:5432/chirpy_db?sslmode=disable" status
```

**Note:** The migration files (001_users.sql, 002_chirps.sql, etc.) will create the necessary tables in the correct order with proper foreign key relationships.

### Step 4: Configure Environment Variables

Create a `.env` file in the root directory:

```bash
# Database connection (use the same credentials from step 3)
DB_URL=postgresql://chirpy_user:your_secure_password@localhost:5432/chirpy_db

# JWT secret key (can be ANY string value - examples below)
CHIRPY_SECRET_KEY=your_super_secret_jwt_key_here_make_it_long_and_random

# Polka webhook API key (can be ANY string value - examples below)
POLKA_KEY=your_polka_webhook_api_key_here
```

**Examples for testing (you can use these or generate your own):**
```bash
# JWT Secret Key Examples (any of these will work):
CHIRPY_SECRET_KEY="2s+OYw/RE7AOjSPHfUtk+EfeCfHELu/a9UloliXnYLhYsGpZWCyFNou+Nq/qYiHFrNEVYzzhIVs6whWEghxRNw=="
CHIRPY_SECRET_KEY="my-super-secret-jwt-key-for-chirpy-api-2024"
CHIRPY_SECRET_KEY="random_string_12345_abcdef_ghijkl_mnopqr_stuvwx_yz"
CHIRPY_SECRET_KEY="chirpy-jwt-secret-key-very-long-and-random-string"

# Polka Webhook API Key Examples (any of these will work):
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e"
POLKA_KEY="my-webhook-api-key-123"
POLKA_KEY="polka_integration_key_2024"
POLKA_KEY="webhook_secret_key_for_external_services"
```

**Important Security Notes:**
- Use a strong, random JWT secret key (at least 32 characters)
- Never commit your `.env` file to version control
- Use different keys for development and production

**How to Generate Secure Keys:**

**Option 1: Use the examples above (for testing only)**
- Copy any of the example keys above for quick testing
- These are NOT secure for production use

**Option 2: Generate secure random keys**
```bash
# Generate a secure JWT secret (64 random characters)
openssl rand -base64 64

# Generate a secure webhook key (32 random characters)
openssl rand -hex 32

# Or use this simple command for a random string
head -c 64 /dev/urandom | base64
```

**Option 3: Use online generators (for development)**
- JWT Secret: Use any long random string (at least 32 characters)
- Webhook Key: Use any string value you prefer

### Step 5: How To Generate Database Code Using sqlc (IT'S ALREADY DONE BY US, YOU DON'T HAVE TO DO IT)

If you're using sqlc for database operations:

```bash
# Install sqlc (if not already installed)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate Go code from SQL queries "sql/queries"
sqlc generate
```

### Step 6: Build and Run

```bash
# Build the application
go build -o chirpy .

# Run the server
./chirpy
```

Or run directly with Go:

```bash
go run .
```

The API will be available at `http://localhost:8080`

## 🧪 Testing Your Setup

### Quick Health Check

```bash
curl http://localhost:8080/api/healthz
```

Expected response: `OK`

### Test User Creation

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpassword123"}'
```

Expected response: 201 Created with user details

## 📱 Getting Started with the API ([You can look at API documentation over here](api.md))

### Complete User Workflow

Here's a step-by-step guide to get you started with the Chirpy API:

#### 1. Create Your Account
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Expected Response (201 Created):**
```json
{
  "id": "uuid-string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

#### 2. Login to Get Access Tokens
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Expected Response (200 OK):**
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

**Save both tokens!** The access token expires in 1 hour, the refresh token in 60 hours.

#### 3. Create Your First Chirp
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "body": "Hello Chirpy! This is my first chirp message."
  }'
```

**Expected Response (201 Created):**
```json
{
  "id": "uuid-string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "Hello Chirpy! This is my first chirp message.",
  "user_id": "user-uuid-string"
}
```

#### 4. View All Chirps
```bash
# Get all chirps
curl http://localhost:8080/api/chirps

# Get chirps by specific user (newest first)
curl "http://localhost:8080/api/chirps?author_id=USER_UUID&sort=desc"

# Get chirps by specific user (oldest first)
curl "http://localhost:8080/api/chirps?author_id=USER_UUID&sort=asc"
```

#### 5. Get a Specific Chirp
```bash
curl http://localhost:8080/api/chirps/CHIRP_UUID
```

#### 6. Update Your Profile
```bash
curl -X PUT http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "email": "newemail@example.com",
    "password": "newpassword123"
  }'
```

#### 7. Refresh Your Access Token
When your access token expires, use the refresh token:
```bash
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer YOUR_REFRESH_TOKEN"
```

**Expected Response (200 OK):**
```json
{
  "token": "new-jwt-access-token"
}
```

#### 8. Delete a Chirp
```bash
curl -X DELETE http://localhost:8080/api/chirps/CHIRP_UUID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Expected Response (204 No Content)**

### 🔑 Authentication Notes

- **Always include the JWT token** in the Authorization header for protected endpoints
- **Format**: `Authorization: Bearer <your_jwt_token>`
- **Token expiration**: Access tokens expire in 1 hour, refresh tokens in 60 hours
- **Refresh flow**: Use `/api/refresh` to get a new access token when it expires

### 📝 Chirp Guidelines

- **Character limit**: Maximum 140 characters per chirp
- **Content filtering**: Banned words (kerfuffle, sharbert, fornax) are automatically replaced with "****"
- **Ownership**: You can only delete your own chirps

### 🧪 Testing Tips

1. **Start with the health check**: `curl http://localhost:8080/api/healthz`
2. **Test user creation first** before trying authenticated endpoints
3. **Save the tokens** from login response for subsequent requests
4. **Use the API Documentation for more details about APIs**

## 🔧 Development

### Project Structure

```
Chirpy/
├── main.go              # Main application entry point
├── users.go             # User management endpoints
├── chirps.go            # Chirp management endpoints
├── webhook.go           # Webhook handling
├── middlewares.go       # HTTP middleware functions
├── api_config.go        # API configuration and admin endpoints
├── readiness.go         # Health check endpoint
├── utils.go             # Utility functions
├── internal/            # Internal packages
│   ├── database/        # Database operations
│   └── auth/           # Authentication utilities
├── sql/                 # SQL queries and migrations
└── assets/              # Static assets
```

## 🔒 Security Considerations

### Production Deployment

1. **Use HTTPS**: Always use HTTPS in production
2. **Strong Secrets**: Use cryptographically secure random keys
3. **Environment Variables**: Never hardcode secrets
4. **Database Security**: Use dedicated database users with minimal privileges

## 🐛 Troubleshooting

### Common Issues

1. **Server Isn't Running**
   - Ensure your server is running:
```bash
# Build the application
go build -o chirpy .

# Run the server
./chirpy
```

3. **Database Connection Failed**
   - Check PostgreSQL is running
   - Verify connection string in `.env`
   - Ensure database exists and user has permissions
   - Ensure you have applied the migrations correctly

4. **Port Already in Use**
   - Change port in `main.go` or use environment variable
   - Check if another service is using port 8080

5. **JWT Token Invalid**
   - Verify `CHIRPY_SECRET_KEY` is set correctly
   - Check token expiration
   - Ensure proper Authorization header format

6. **Permission Denied**
   - Check file permissions
   - Ensure database user has proper privileges
   - Verify JWT token is valid and not expired

## 📚 Additional Resources

- **API Documentation**: See `API_DOCUMENTATION.md` for complete API reference [See API Docs](api.md)

## 🆘 Support

- **Issues**: Create an issue on GitHub
- **Documentation**: Check the API documentation files
- **Testing**: Use the provided API Documentation collection

---

**Happy Chirping! 🐦**
