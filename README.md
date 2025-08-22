# Library Management System API

A RESTful API for managing a library's book collection with user authentication and role-based access control. Built with Go, MongoDB, and JWT authentication.

## Features

- 📚 **Book Management**: CRUD operations for books with detailed metadata
- 🔐 **User Authentication**: JWT-based authentication system
- 👥 **Role-Based Access**: Admin and user roles with different permissions
- 🛡️ **Security**: Password hashing with bcrypt, secure JWT tokens
- 📝 **Structured Logging**: Comprehensive logging with logrus
- 🏗️ **Clean Architecture**: Repository pattern with service layers
- ⚡ **Performance**: MongoDB with optimized queries and timeouts

## Tech Stack

- **Language**: Go 1.23.2
- **Database**: MongoDB
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **HTTP Router**: Gorilla Mux
- **Logging**: Logrus
- **Environment**: Docker-ready

## Project Structure

```
my-library/
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── handlers/                # HTTP request handlers
│   │   ├── book.go             # Book-related endpoints
│   │   └── user.go             # Authentication endpoints
│   ├── services/                # Business logic layer
│   │   ├── book_service.go
│   │   └── user_service.go
│   ├── repositories/            # Data access layer
│   │   ├── book_repository.go
│   │   └── user_repository.go
│   ├── models/                  # Data models
│   │   ├── book.go
│   │   ├── user.go
│   │   └── response.go
│   ├── middleware/              # HTTP middleware
│   │   ├── auth.go             # JWT authentication
│   │   └── logging.go          # Request logging
│   ├── database/               # Database connection
│   │   └── database.go
│   └── logger/                 # Logging utilities
│       └── logger.go
├── .env                        # Environment variables
├── go.mod                      # Go modules
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.23+ installed
- MongoDB instance (local or MongoDB Atlas)
- Git

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/4Noyis/my-library.git
cd my-library
```

2. **Install dependencies**
```bash
go mod download
```

3. **Environment Setup**
```bash
# Copy and configure environment variables
cp .env.example .env
```

Configure your `.env` file:
```env
MONGO_URI="mongodb://localhost:27017"  # or your MongoDB Atlas URI
JWT_SECRET="your-super-secret-jwt-key-change-in-production"
```

4. **Run the application**
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "secure_password",
    "role": "user"  // optional, defaults to "user"
}
```

**Response:**
```json
{
    "status": "success",
    "message": "User registered successfully",
    "data": {
        "id": "64f5a7b2e123456789abcdef",
        "username": "john_doe",
        "email": "john@example.com",
        "role": "user",
        "is_active": true,
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
    }
}
```

#### Login User
```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "john_doe",
    "password": "secure_password"
}
```

**Response:**
```json
{
    "status": "success",
    "message": "Login successful",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": "64f5a7b2e123456789abcdef",
            "username": "john_doe",
            "email": "john@example.com",
            "role": "user",
            "is_active": true
        }
    }
}
```

### Book Endpoints

**Note**: All book endpoints require authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer YOUR_JWT_TOKEN
```

#### Get All Books
```http
GET /api/v1/books
Authorization: Bearer YOUR_JWT_TOKEN
```

#### Get Book by ID
```http
GET /api/v1/books/{id}
Authorization: Bearer YOUR_JWT_TOKEN
```

#### Create New Book
```http
POST /api/v1/books
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
    "isbn": "978-0123456789",
    "title": "The Go Programming Language",
    "author": "Alan Donovan",
    "publisher": "Addison-Wesley",
    "published_at": "2015-11-16T00:00:00Z",
    "genre": "Programming",
    "language": "English",
    "pages": 380,
    "description": "A comprehensive guide to Go programming",
    "coverURL": "https://example.com/cover.jpg",
    "location": "A1-B2"
}
```

#### Update Book
```http
PATCH /api/v1/books/{id}
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
    "title": "Updated Title",
    "location": "A2-B3"
}
```

#### Delete Book
```http
DELETE /api/v1/books/{id}
Authorization: Bearer YOUR_JWT_TOKEN
```

## Data Models

### Book Model
```go
type Book struct {
    ID          int       `json:"id"`
    ISBN        string    `json:"isbn"`
    Title       string    `json:"title"`
    Author      string    `json:"author"`
    Publisher   string    `json:"publisher"`
    PublishedAt time.Time `json:"published_at"`
    Genre       string    `json:"genre"`
    Language    string    `json:"language"`
    Pages       int       `json:"pages"`
    Description string    `json:"description"`
    CoverURL    string    `json:"coverURL"`
    Location    string    `json:"location"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### User Model
```go
type User struct {
    ID        primitive.ObjectID `json:"id"`
    Username  string             `json:"username"`
    Email     string             `json:"email"`
    Role      string             `json:"role"`      // "admin" or "user"
    IsActive  bool               `json:"is_active"`
    CreatedAt time.Time          `json:"created_at"`
    UpdatedAt time.Time          `json:"updated_at"`
}
```

## Development

### Building
```bash
# Build the application
go build ./cmd/server

# Run with hot reload (install air first)
go install github.com/cosmtrek/air@latest
air
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run static analysis
go vet ./...

# Run linter (install golangci-lint first)
golangci-lint run
```

## Database Collections

### Books Collection
- **Database**: `library`
- **Collection**: `books`
- **ID Type**: Auto-incremented integer

### Users Collection
- **Database**: `library`
- **Collection**: `users`
- **ID Type**: MongoDB ObjectID

## Security Features

- **Password Hashing**: Uses bcrypt with default cost
- **JWT Tokens**: 24-hour expiration, HMAC-SHA256 signing
- **Role-Based Access**: Admin and user roles
- **Input Validation**: Request body validation
- **CORS Ready**: Easy to configure for frontend applications

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `MONGO_URI` | MongoDB connection string | - | Yes |
| `JWT_SECRET` | Secret key for JWT signing | Development key | Yes |
| `PORT` | Server port | 8080 | No |

## Error Handling

The API returns standardized error responses:

```json
{
    "status": "error",
    "message": "Error description",
    "data": null
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict (duplicate resource)
- `500` - Internal Server Error

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
