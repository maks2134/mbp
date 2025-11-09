# Testing Guide

This directory contains test utilities and integration tests for the MPB Blog Platform.

## Test Structure

```
tests/
├── testutils/          # Test utilities and helpers
│   ├── testutils.go    # Database, Redis, PubSub setup
│   └── logger.go       # Logger setup
├── integration/        # Integration tests
│   └── posts_integration_test.go
└── README.md           # This file
```

## Test Types

### Unit Tests

Unit tests are located alongside the code they test (e.g., `internal/posts/service_test.go`). They use mocks to isolate the code under test.

**Running unit tests:**
```bash
go test ./internal/posts/...
go test ./internal/auth/...
```

### Integration Tests

Integration tests are in the `tests/integration/` directory. They require real database and Redis connections.

**Running integration tests:**
```bash
# Run all tests including integration
go test -tags=integration ./tests/integration/...

# Skip integration tests
go test -short ./...
```

## Prerequisites for Integration Tests

1. **PostgreSQL** running on `localhost:5432`
   - Database: `mpb_test`
   - User: `mpb`
   - Password: `mpb_pas`

2. **Redis** running on `localhost:6379`

3. **Environment Variables** (optional):
   ```bash
   export TEST_DSN="postgres://mpb:mpb_pas@localhost:5432/mpb_test?sslmode=disable"
   export TEST_REDIS_ADDR="localhost:6379"
   export TEST_JWT_SECRET="test-secret-key"
   ```

## Running Tests

### Run All Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test-coverage
```

### Run Specific Test Package
```bash
go test ./internal/posts/...
go test ./internal/auth/...
```

### Run Integration Tests Only
```bash
go test -tags=integration ./tests/integration/...
```

### Run with Verbose Output
```bash
go test -v ./...
```

## Test Utilities

### testutils Package

The `testutils` package provides helper functions for setting up test environments:

- `TestConfig()` - Returns test configuration
- `SetupTestDB(t)` - Creates test database connection
- `SetupTestRedis(t)` - Creates test Redis connection
- `CleanupRedis(t, client)` - Cleans up Redis data
- `SetupTestPubSub(t)` - Creates test Watermill pub/sub
- `SetupTestLogger(t)` - Creates test logger

### Example Usage

```go
func TestMyFeature(t *testing.T) {
    database := testutils.SetupTestDB(t)
    defer database.Close()
    
    redisClient := testutils.SetupTestRedis(t)
    defer redisClient.Close()
    defer testutils.CleanupRedis(t, redisClient.Client)
    
    // Your test code here
}
```

## Writing Tests

### Unit Test Example

```go
func TestMyService(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("Find", 1).Return(&Entity{ID: 1}, nil)
    
    service := NewService(mockRepo)
    result, err := service.Get(1)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

### Integration Test Example

```go
// +build integration

func TestMyFeatureIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    database := testutils.SetupTestDB(t)
    defer database.Close()
    
    // Test with real database
}
```

## Best Practices

1. **Use Table-Driven Tests**: For multiple test cases with similar structure
2. **Mock External Dependencies**: Use mocks for unit tests
3. **Clean Up Resources**: Always clean up test data and connections
4. **Test Both Success and Error Cases**: Cover all code paths
5. **Use Descriptive Test Names**: Test names should describe what is being tested
6. **Keep Tests Fast**: Unit tests should run quickly
7. **Isolate Tests**: Tests should not depend on each other

## Coverage Goals

- **Unit Tests**: >80% coverage for services and repositories
- **Integration Tests**: Cover critical user flows
- **Handler Tests**: Test HTTP request/response handling

## Continuous Integration

Tests are automatically run in CI/CD pipeline:
- Unit tests run on every commit
- Integration tests run on pull requests
- Coverage reports are generated

## Troubleshooting

### Tests Fail with Database Connection Error

1. Ensure PostgreSQL is running
2. Check connection string in environment variables
3. Verify database `mpb_test` exists

### Tests Fail with Redis Connection Error

1. Ensure Redis is running
2. Check Redis address in environment variables
3. Verify Redis is accessible

### Integration Tests Timeout

1. Check database and Redis performance
2. Increase timeout if needed: `go test -timeout 30s`
3. Ensure test data is cleaned up properly

