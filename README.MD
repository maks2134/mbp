# 🔐 MPB Blog Platform

Backend API for MPB Blog Platform, built with **Go + Fiber + PostgreSQL + JWT**.  
Implements secure user registration, login, JWT-based authentication, and blog post management with event-driven architecture.

---

## 🚀 Features

- **Authentication & Authorization**
  - User registration and login with JWT tokens
  - JWT-based authentication middleware
  - Token refresh mechanism
  
- **Blog Posts Management**
  - Full CRUD operations for posts
  - Post ownership validation
  - Post attachments (images, files) via S3
  - Real-time metrics: likes and views tracking
  - Redis-based caching for hot data
  
- **Comments System**
  - Create, update, and delete comments on posts
  - Comment attachments support
  - Nested comment structure
  
- **Event-Driven Architecture**
  - Watermill for event publishing/subscribing
  - Asynchronous metrics synchronization (Redis → PostgreSQL)
  - Event types: `post.viewed`, `post.liked`, `post.unliked`
  
- **Infrastructure**
  - PostgreSQL for persistent data storage
  - Redis for real-time metrics and caching
  - AWS S3 for file storage
  - Docker & Docker Compose for easy deployment
  - Database migrations with Goose
  - CI/CD with GitHub Actions
  - Integrated Swagger documentation
  - Input validation via generics + go-playground/validator

---

## ⚙️ Tech Stack

| Component       | Technology/Library |
|-----------------|-------------------|
| Web Framework   | [Fiber v2](https://github.com/gofiber/fiber) |
| DB Access       | [sqlx](https://github.com/jmoiron/sqlx) |
| Database        | PostgreSQL 15+ |
| Cache/Metrics    | Redis 7+ |
| File Storage     | AWS S3 |
| JWT Auth         | [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) |
| Validation       | [go-playground/validator/v10](https://github.com/go-playground/validator) |
| Documentation    | [fiber-swagger](https://github.com/swaggo/fiber-swagger) |
| Migrations       | [goose](https://github.com/pressly/goose) |
| Event Bus        | [watermill](https://github.com/ThreeDotsLabs/watermill) |
| Containerization | Docker & Docker Compose |

---

## 📋 Prerequisites

- **Go** 1.25+
- **PostgreSQL** 15+
- **Redis** 7+ (for metrics and caching)
- **Docker** and **Docker Compose** (recommended for development)
- **AWS Account** (optional, for S3 file storage)

---

## 🚀 Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd mpb
   ```

2. **Create `.env` file:**
   ```bash
   cp .env_example .env
   # Edit .env with your values
   ```

3. **Start services:**
   ```bash
   docker-compose up -d
   ```

   This will:
   - Start PostgreSQL database
   - Start Redis server
   - Run migrations automatically
   - Start the application on `http://localhost:8000`

4. **View logs:**
   ```bash
   docker-compose logs -f app
   ```

5. **Access Swagger documentation:**
   - Open `http://localhost:8000/swagger/index.html`

### Manual Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Start PostgreSQL** (or use docker-compose):
   ```bash
   docker-compose up -d postgres
   ```

3. **Set environment variables:**
   ```bash
   export DSN="postgres://mpb:mpb_pas@localhost:5432/mpb_db?sslmode=disable"
   export REDIS_ADDR="localhost:6379"
   export JWT_SECRET="your-secret-key-here"
   export JWT_TTL="24h"
   # Optional: for S3 file storage
   export AWS_REGION="us-east-1"
   export AWS_BUCKET="your-bucket-name"
   ```

4. **Run migrations:**
   ```bash
   make migrate-up
   # or
   go run ./cmd/migrate.go
   ```

5. **Run the application:**
   ```bash
   make run
   # or
   go run ./cmd/main.go
   ```

---

## 🛠️ Development

### Available Make Commands

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the application
make test          # Run tests
make test-coverage # Run tests with coverage report
make lint          # Run linter (requires golangci-lint)
make clean         # Clean build artifacts
make dev           # Start development environment
```

### Database Migrations

```bash
make migrate-up     # Run migrations
make migrate-down   # Rollback last migration
make migrate-status # Show migration status
```

### Generate Swagger Documentation

```bash
make swagger
```

---

## 🐳 Docker

### Build Docker Image

```bash
docker build -t mpb:latest .
```

### Run with Docker

```bash
docker run -d \
  -p 8000:8000 \
  -e DSN="postgres://mpb:mpb_pas@host.docker.internal:5432/mpb_db?sslmode=disable" \
  -e JWT_SECRET="your-secret-key" \
  -e JWT_TTL="24h" \
  mpb:latest
```

### Docker Compose

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f app

# Rebuild and restart
docker-compose up -d --build
```

---

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### CI/CD

The project includes GitHub Actions workflows:

- **CI Pipeline** (`.github/workflows/ci.yml`):
  - Linting with golangci-lint
  - Running tests with PostgreSQL service
  - Building the application
  - Building and pushing Docker images
  - Security scanning with Gosec and Trivy

- **CodeQL Analysis** (`.github/workflows/codeql.yml`):
  - Automated code security analysis

---

## 📝 API Endpoints

### Authentication

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login and get JWT token

### Posts (Protected)

- `GET /api/posts` - List all active posts
- `GET /api/posts/{id}` - Get post by ID (increments view count)
- `POST /api/posts` - Create new post (requires auth)
- `PUT /api/posts/{id}` - Update post (requires auth, owner only)
- `DELETE /api/posts/{id}` - Delete post (requires auth, owner only)
- `POST /api/posts/{id}/like` - Like a post (requires auth)
- `DELETE /api/posts/{id}/unlike` - Unlike a post (requires auth)

### Post Attachments

- `POST /api/posts/{id}/attachments` - Upload attachment to post (requires auth, owner only)
- `GET /api/posts/{id}/attachments` - Get all attachments for a post
- `DELETE /api/posts/{id}/attachments/{attachmentId}` - Delete attachment (requires auth, owner only)

### Comments

- `GET /api/posts/{id}/comments` - Get all comments for a post
- `POST /api/posts/{id}/comments` - Create comment on post (requires auth)
- `PUT /api/comments/{id}` - Update comment (requires auth, owner only)
- `DELETE /api/comments/{id}` - Delete comment (requires auth, owner only)

### Comment Attachments

- `POST /api/comments/{id}/attachments` - Upload attachment to comment (requires auth, owner only)
- `GET /api/comments/{id}/attachments` - Get all attachments for a comment
- `DELETE /api/comments/{id}/attachments/{attachmentId}` - Delete attachment (requires auth, owner only)

### Documentation

- `GET /swagger/*` - Swagger UI

---

## 🔒 Environment Variables

| Variable      | Description                          | Default                                    | Required |
|---------------|--------------------------------------|--------------------------------------------|----------|
| `DSN`         | PostgreSQL connection string         | `postgres://mpb:mpb_pas@localhost:5432/mpb_db?sslmode=disable` | Yes |
| `REDIS_ADDR`  | Redis server address                 | `localhost:6379`                           | Yes |
| `JWT_SECRET`  | Secret key for JWT signing           | -                                          | Yes |
| `JWT_TTL`     | JWT token TTL                        | `24h`                                      | No |
| `AWS_REGION`  | AWS region for S3                    | -                                          | No* |
| `AWS_BUCKET`  | AWS S3 bucket name                   | -                                          | No* |

\* Required only if using S3 file storage

---

## 📦 Project Structure

```
mpb/
├── cmd/
│   ├── main.go              # Application entry point
│   ├── migrate.go            # Migration tool
│   └── healthcheck/         # Health check utility
├── configs/
│   └── config.go            # Configuration loading
├── internal/
│   ├── app/                 # Application initialization
│   ├── auth/                # Authentication module
│   │   └── dto/             # Auth DTOs
│   ├── posts/               # Posts module
│   │   ├── dto/             # Post DTOs
│   │   ├── events.go        # Event definitions
│   │   ├── metrics_service.go    # Redis metrics service
│   │   └── metrics_consumer.go   # Event consumer for sync
│   ├── comments/            # Comments module
│   │   └── dto/             # Comment DTOs
│   ├── post_attachments/    # Post attachments module
│   ├── comments_attachments/ # Comment attachments module
│   └── user/                # User model
├── migrations/              # Database migrations (Goose)
├── pkg/
│   ├── db/                  # Database connection (sqlx)
│   ├── redis/               # Redis client wrapper
│   ├── s3/                  # AWS S3 client
│   ├── errors_constant/     # Error constants
│   ├── middleware/          # HTTP middleware (JWT, validation)
│   └── security/            # Security utilities (hashing)
├── docs/                    # Swagger documentation
├── Dockerfile               # Docker build configuration
├── docker-compose.yml       # Docker Compose configuration
├── Makefile                 # Build automation
└── go.mod                   # Go dependencies
```

---

## 🏗️ Architecture

The application follows a **layered architecture** with **event-driven** components:

- **Handler Layer**: HTTP request/response handling
- **Service Layer**: Business logic
- **Repository Layer**: Data access (PostgreSQL)
- **Event Layer**: Asynchronous event processing (Watermill)
- **Cache Layer**: Redis for hot data (likes, views)

### Data Flow

1. **Metrics (Likes/Views)**: 
   - User action → Redis (immediate) → Event published → PostgreSQL (async sync)

2. **Posts/Comments**:
   - User action → PostgreSQL → Response

3. **Attachments**:
   - User upload → S3 → Metadata saved to PostgreSQL

For detailed architecture documentation, see [ARCHITECTURE.md](ARCHITECTURE.md).

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## 📄 License

This project is licensed under the MIT License.

---

## 🔗 Useful Links

- [Fiber Documentation](https://docs.gofiber.io/)
- [Swagger/OpenAPI](https://swagger.io/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
