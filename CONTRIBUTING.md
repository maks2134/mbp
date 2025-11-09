# Contributing to MPB Blog Platform

Thank you for your interest in contributing to MPB Blog Platform! This document provides guidelines and instructions for contributing.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Documentation](#documentation)

## 📜 Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Follow the project's coding standards

## 🚀 Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (recommended)
- Git

### Setting Up Development Environment

1. **Fork and Clone**
   ```bash
   git clone https://github.com/your-username/mpb.git
   cd mpb
   ```

2. **Set Up Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start Services**
   ```bash
   docker-compose up -d
   ```

4. **Run Migrations**
   ```bash
   make migrate-up
   # or
   go run ./cmd/migrate.go
   ```

5. **Run the Application**
   ```bash
   make run
   # or
   go run ./cmd/main.go
   ```

6. **Verify Setup**
   - Check Swagger: http://localhost:8000/swagger/index.html
   - Test health endpoint (if available)

## 🔄 Development Workflow

### Branch Naming

Use descriptive branch names with prefixes:

- `feature/` - New features (e.g., `feature/user-profile`)
- `fix/` - Bug fixes (e.g., `fix/auth-token-expiry`)
- `refactor/` - Code refactoring (e.g., `refactor/post-service`)
- `docs/` - Documentation updates (e.g., `docs/api-endpoints`)
- `test/` - Test additions/updates (e.g., `test/post-handler`)

### Workflow Steps

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**
   - Write code following the coding standards
   - Add tests for new functionality
   - Update documentation if needed

3. **Test Your Changes**
   ```bash
   make test
   make lint
   ```

4. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "feat: add user profile endpoint"
   ```

5. **Push and Create PR**
   ```bash
   git push origin feature/your-feature-name
   ```

## 💻 Coding Standards

### Go Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` and `golint` (or `golangci-lint`)
- Maximum line length: 120 characters
- Use meaningful variable and function names

### Project Structure

Follow the existing project structure:

```
internal/
  ├── {module}/
  │   ├── handler.go      # HTTP handlers
  │   ├── service.go       # Business logic
  │   ├── repository.go    # Data access
  │   ├── model.go         # Domain models
  │   ├── routes.go        # Route definitions
  │   └── dto/             # Data Transfer Objects
```

### Error Handling

- Use custom error constants from `pkg/errors_constant`
- Always handle errors explicitly (no silent failures)
- Provide meaningful error messages
- Log errors with context

```go
// Good
if err != nil {
    return fmt.Errorf("failed to create post: %w", err)
}

// Bad
if err != nil {
    return err
}
```

### Database Access

- Use `sqlx` for all database operations
- Always use prepared statements or parameterized queries
- Handle transactions properly
- Close rows and connections

### Event Publishing

- Publish events after successful operations
- Use structured event types from `events.go`
- Handle event publishing errors gracefully

```go
event := PostLikedEvent{
    PostID: postID,
    UserID: userID,
    Likes:  likes,
}
if err := s.publishEvent("post.liked", event); err != nil {
    s.logger.Error("failed to publish event", err, nil)
}
```

### Testing

- Write unit tests for services and repositories
- Write integration tests for handlers
- Aim for >80% code coverage
- Use table-driven tests where appropriate

```go
func TestCreatePost(t *testing.T) {
    tests := []struct {
        name    string
        input   CreatePostRequest
        wantErr bool
    }{
        {
            name: "valid post",
            input: CreatePostRequest{
                Title: "Test Post",
                Description: "Test Description",
            },
            wantErr: false,
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## 📝 Commit Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions/changes
- `chore`: Build process or auxiliary tool changes

### Examples

```bash
feat(posts): add like/unlike endpoints
fix(auth): handle token expiration correctly
docs(readme): update installation instructions
refactor(service): simplify post creation logic
test(handler): add tests for post endpoints
```

### Commit Message Best Practices

- Use imperative mood ("add" not "added")
- Keep subject line under 50 characters
- Capitalize first letter
- No period at the end
- Reference issues: `fix: resolve #123`

## 🔍 Pull Request Process

### Before Submitting

1. **Update Documentation**
   - Update README if needed
   - Add/update API documentation
   - Update architecture docs if architecture changed

2. **Run Tests**
   ```bash
   make test
   make test-coverage
   ```

3. **Run Linter**
   ```bash
   make lint
   ```

4. **Check Build**
   ```bash
   make build
   ```

5. **Test Manually**
   - Test your changes locally
   - Verify with Docker Compose

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests pass locally
```

### Review Process

1. Maintainers will review your PR
2. Address any feedback
3. Ensure CI checks pass
4. PR will be merged after approval

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./internal/posts/...

# Run with verbose output
go test -v ./internal/posts/...
```

### Writing Tests

- Test file naming: `*_test.go`
- Use `testing` package
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Test both success and error cases

## 📚 Documentation

### Code Documentation

- Document all exported functions, types, and packages
- Use Go doc comments
- Include examples for complex functions

```go
// CreatePost creates a new blog post for the authenticated user.
// It validates the input, creates the post in the database,
// and publishes a post.created event.
//
// Example:
//   post, err := service.CreatePost(ctx, userID, "Title", "Description", "tag")
func CreatePost(ctx context.Context, userID int, title, description, tag string) (*Post, error) {
    // implementation
}
```

### API Documentation

- Update Swagger annotations when adding/modifying endpoints
- Regenerate Swagger docs: `make swagger`
- Keep examples up to date

## 🐛 Reporting Bugs

### Before Reporting

1. Check existing issues
2. Verify it's a bug, not a feature request
3. Try to reproduce the issue

### Bug Report Template

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce:
1. ...
2. ...

**Expected behavior**
What you expected to happen.

**Environment:**
- OS: [e.g., macOS 14.0]
- Go version: [e.g., 1.25]
- Version: [e.g., commit hash]

**Additional context**
Any other relevant information.
```

## 💡 Feature Requests

1. Check if the feature already exists
2. Open an issue with the `enhancement` label
3. Describe the use case and benefits
4. Wait for discussion before implementing

## 📞 Getting Help

- Open an issue for bugs or questions
- Check existing documentation
- Review closed issues for similar problems

## 🙏 Thank You!

Your contributions make this project better. Thank you for taking the time to contribute!

