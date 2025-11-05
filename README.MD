# MPB Blog Platform

Backend API for MPB Blog Platform, built with **Go + Fiber + PostgreSQL + JWT**.  
Implements secure user registration, login, JWT-based authentication, and blog post management with event-driven architecture.

---

## 🚀 Features

- User registration and authentication (`POST /api/auth/register`, `POST /api/auth/login`)
- JWT authentication middleware
- Blog posts CRUD operations with authorization
- Event-driven architecture with Watermill
- Input validation via generics + go-playground/validator
- Integrated Swagger documentation
- Database migrations with Goose
- Docker support with multi-stage builds
- CI/CD with GitHub Actions

---

## ⚙️ Tech Stack

| Component       | Library |
|-----------------|----------|
| Web Framework   | [Fiber v2](https://github.com/gofiber/fiber) |
| DB Access       | [sqlx](https://github.com/jmoiron/sqlx) |
| Database        | PostgreSQL, Redis |
| JWT Auth        | [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) |
| Validation      | [go-playground/validator/v10](https://github.com/go-playground/validator) |
| Documentation   | [fiber-swagger](https://github.com/swaggo/fiber-swagger) |
| Migrations      | [goose](https://github.com/pressly/goose) |
| Events          | [watermill](https://github.com/ThreeDotsLabs/watermill) |

---

## 📋 Prerequisites

- Go 1.25+
- PostgreSQL 15+
- Docker and Docker Compose (optional)

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
   export JWT_SECRET="your-secret-key-here"
   export JWT_TTL="24h"
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

- `GET /api/posts` - List all posts
- `GET /api/posts/{id}` - Get post by ID
- `POST /api/posts` - Create new post (requires auth)
- `PUT /api/posts/{id}` - Update post (requires auth, owner only)
- `DELETE /api/posts/{id}` - Delete post (requires auth, owner only)

### Documentation

- `GET /swagger/*` - Swagger UI

---

## 🔒 Environment Variables

| Variable    | Description                          | Default                                    |
|-------------|--------------------------------------|--------------------------------------------|
| `DSN`       | PostgreSQL connection string        | `postgres://mpb:mpb_pas@localhost:5432/mpb_db?sslmode=disable` |
| `JWT_SECRET` | Secret key for JWT signing          | Required                                   |
| `JWT_TTL`   | JWT token TTL                       | `24h`                                      |

---

## 📦 Project Structure

```
mpb/
├── cmd/
│   ├── main.go          # Application entry point
│   └── migrate.go       # Migration tool
├── configs/
│   └── config.go        # Configuration loading
├── internal/
│   ├── auth/            # Authentication module
│   ├── posts/           # Posts module
│   └── user/            # User model
├── migrations/          # Database migrations
├── pkg/
│   ├── db/              # Database connection
│   ├── errors_constant/ # Error constants
│   ├── middleware/      # HTTP middleware
│   └── security/        # Security utilities
├── docs/                # Swagger documentation
├── Dockerfile           # Docker build configuration
├── docker-compose.yml   # Docker Compose configuration
├── Makefile             # Build automation
└── go.mod               # Go dependencies
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License.

---

## 🔗 Useful Links

- [Fiber Documentation](https://docs.gofiber.io/)
- [Swagger/OpenAPI](https://swagger.io/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
