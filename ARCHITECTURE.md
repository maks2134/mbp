# Architecture Documentation

## 📋 Table of Contents

- [Overview](#overview)
- [System Architecture](#system-architecture)
- [Technology Stack](#technology-stack)
- [Layered Architecture](#layered-architecture)
- [Event-Driven Components](#event-driven-components)
- [Data Flow](#data-flow)
- [Database Schema](#database-schema)
- [Caching Strategy](#caching-strategy)
- [File Storage](#file-storage)
- [Security](#security)
- [Scalability Considerations](#scalability-considerations)

## 🎯 Overview

MPB Blog Platform is a RESTful API built with Go, implementing a blog platform with posts, comments, attachments, and real-time metrics (likes/views). The system uses a **layered architecture** combined with **event-driven** patterns for asynchronous processing.

### Key Design Principles

- **Separation of Concerns**: Clear separation between handlers, services, and repositories
- **Event-Driven**: Asynchronous processing for metrics synchronization
- **Performance**: Redis caching for hot data (likes, views)
- **Scalability**: Stateless API design, horizontal scaling ready
- **Maintainability**: Modular structure, clear interfaces

## 🏗️ System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Client Layer                         │
│                    (HTTP/REST Clients)                     │
└────────────────────────┬────────────────────────────────────┘
                          │
┌─────────────────────────▼────────────────────────────────────┐
│                    API Gateway Layer                        │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Fiber HTTP Server                        │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐            │  │
│  │  │  Auth   │  │  Posts   │  │ Comments │  ...        │  │
│  │  │ Routes  │  │  Routes  │  │  Routes  │            │  │
│  │  └──────────┘  └──────────┘  └──────────┘            │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────┬────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
┌───────▼──────┐  ┌───────▼──────┐  ┌───────▼──────┐
│   Handler    │  │   Handler    │  │   Handler    │
│    Layer     │  │    Layer     │  │    Layer     │
└───────┬──────┘  └───────┬──────┘  └───────┬──────┘
        │                 │                 │
┌───────▼──────┐  ┌───────▼──────┐  ┌───────▼──────┐
│   Service    │  │   Service    │  │   Service    │
│    Layer     │  │    Layer     │  │    Layer     │
└───────┬──────┘  └───────┬──────┘  └───────┬──────┘
        │                 │                 │
┌───────▼─────────────────▼─────────────────▼──────┐
│              Repository Layer                     │
│         (PostgreSQL via sqlx)                      │
└───────────────────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
┌───────▼──────┐  ┌───────▼──────┐  ┌───────▼──────┐
│ PostgreSQL  │  │    Redis     │  │   AWS S3     │
│  Database   │  │   (Metrics)   │  │  (Files)     │
└─────────────┘  └───────────────┘  └─────────────┘
                          │
                  ┌───────▼──────┐
                  │   Watermill  │
                  │  Event Bus   │
                  └──────────────┘
```

## 🛠️ Technology Stack

### Core Technologies

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Web Framework** | Fiber v2 | HTTP server and routing |
| **Database** | PostgreSQL 15+ | Primary data storage |
| **Cache** | Redis 7+ | Hot data (likes, views) |
| **File Storage** | AWS S3 | Post/comment attachments |
| **Event Bus** | Watermill (GoChannel) | Event-driven architecture |
| **ORM/Query Builder** | sqlx | Database access |
| **Authentication** | JWT (golang-jwt/jwt/v5) | Token-based auth |
| **Validation** | go-playground/validator | Input validation |
| **Migrations** | Goose | Database schema management |

## 📐 Layered Architecture

### 1. Handler Layer (`internal/*/handler.go`)

**Responsibility**: HTTP request/response handling

- Parse and validate HTTP requests
- Extract user context (JWT)
- Call service layer
- Format HTTP responses
- Handle HTTP-specific errors

**Example Flow**:
```go
func (h *PostsHandlers) CreatePost(c *fiber.Ctx) error {
    // 1. Parse request
    req := middleware.Body[dto.CreatePostRequest](c)
    
    // 2. Extract user context
    userID := c.Locals("user_id").(int)
    
    // 3. Call service
    post, err := h.service.CreatePost(ctx, userID, req.Title, ...)
    
    // 4. Format response
    return c.Status(201).JSON(postToResponse(post))
}
```

### 2. Service Layer (`internal/*/service.go`)

**Responsibility**: Business logic

- Implement business rules
- Coordinate between repositories
- Publish events
- Handle business-level errors

**Example Flow**:
```go
func (s *PostsService) CreatePost(ctx context.Context, userID int, ...) (*Post, error) {
    // 1. Validate business rules
    if len(title) < 3 {
        return nil, errors_constant.InvalidTitle
    }
    
    // 2. Create in database
    post, err := s.repo.Create(ctx, post)
    
    // 3. Publish event
    s.publishEvent("post.created", PostCreatedEvent{...})
    
    return post, nil
}
```

### 3. Repository Layer (`internal/*/repository.go`)

**Responsibility**: Data access

- Database queries (PostgreSQL)
- Data mapping (struct ↔ database)
- Transaction management
- Query optimization

**Example**:
```go
func (r *PostsRepository) Create(ctx context.Context, post *Post) (*Post, error) {
    query := `INSERT INTO posts (user_id, title, description, tag) 
              VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
    
    err := r.db.QueryRowxContext(ctx, query, ...).Scan(&post.ID, ...)
    return post, err
}
```

### 4. DTO Layer (`internal/*/dto/`)

**Responsibility**: Data Transfer Objects

- Separate API contracts from domain models
- Request/response validation
- API versioning support

## 🔄 Event-Driven Components

### Event Flow Architecture

```
User Action
    │
    ▼
Service Layer
    │
    ├──► Redis (immediate update)
    │
    └──► Publish Event (Watermill)
            │
            ▼
        Event Bus (GoChannel)
            │
            ▼
        Event Consumer
            │
            ▼
        PostgreSQL (async sync)
```

### Event Types

#### 1. Post Events

- **`post.viewed`**: Published when a post is viewed
  ```go
  type PostViewedEvent struct {
      PostID int `json:"post_id"`
      Views  int `json:"views"`
  }
  ```

- **`post.liked`**: Published when a post is liked
  ```go
  type PostLikedEvent struct {
      PostID int `json:"post_id"`
      UserID int `json:"user_id"`
      Likes  int `json:"likes"`
  }
  ```

- **`post.unliked`**: Published when a post is unliked
  ```go
  type PostUnlikedEvent struct {
      PostID int `json:"post_id"`
      UserID int `json:"user_id"`
      Likes  int `json:"likes"`
  }
  ```

#### 2. Event Consumer (`metrics_consumer.go`)

**Responsibility**: Synchronize Redis metrics to PostgreSQL

- Subscribes to events: `post.viewed`, `post.liked`, `post.unliked`
- Updates PostgreSQL asynchronously
- Handles errors and retries

**Flow**:
```go
1. Subscribe to events
2. Receive event message
3. Unmarshal event
4. Update PostgreSQL
5. Acknowledge message
```

### Watermill Configuration

- **Pub/Sub**: GoChannel (in-memory, single instance)
- **Publisher/Subscriber**: Same instance (required for GoChannel)
- **Topics**: `post.viewed`, `post.liked`, `post.unliked`

## 📊 Data Flow

### 1. Post Creation Flow

```
Client → Handler → Service → Repository → PostgreSQL
                              │
                              └──► Publish Event → Event Bus
```

### 2. Post View Flow

```
Client → Handler → Service → MetricsService → Redis (increment)
                              │
                              └──► Publish Event → Event Bus → Consumer → PostgreSQL
```

### 3. Post Like Flow

```
Client → Handler → MetricsService → Redis (check + increment)
                              │
                              └──► Publish Event → Event Bus → Consumer → PostgreSQL
```

### 4. File Upload Flow

```
Client → Handler → Service → S3 (upload file)
                              │
                              └──► Repository → PostgreSQL (save metadata)
```

## 🗄️ Database Schema

### Core Tables

#### `users`
- `id` (PK)
- `username` (unique)
- `email` (unique)
- `password_hash`
- `name`
- `age`
- `created_at`, `updated_at`

#### `posts`
- `id` (PK)
- `user_id` (FK → users)
- `title`
- `description`
- `tag`
- `like` (synced from Redis)
- `count_viewers` (synced from Redis)
- `is_active`
- `created_at`, `updated_at`

#### `comments`
- `id` (PK)
- `post_id` (FK → posts)
- `user_id` (FK → users)
- `content`
- `created_at`, `updated_at`

#### `post_attachments`
- `id` (PK)
- `post_id` (FK → posts)
- `file_url` (S3 URL)
- `file_name`
- `file_size`
- `mime_type`
- `created_at`

#### `comment_attachments`
- `id` (PK)
- `comment_id` (FK → comments)
- `file_url` (S3 URL)
- `file_name`
- `file_size`
- `mime_type`
- `created_at`

## 💾 Caching Strategy

### Redis Usage

**Purpose**: Store hot data for fast access

#### Data Stored

1. **Post Likes**: `post:likes:{postID}` → integer
2. **Post Views**: `post:views:{postID}` → integer
3. **User Likes**: `user:{userID}:liked:{postID}` → boolean (set)

#### Operations

- **Increment Views**: `INCR post:views:{postID}`
- **Increment Likes**: `INCR post:likes:{postID}`
- **Check User Liked**: `GET user:{userID}:liked:{postID}`
- **Set User Liked**: `SET user:{userID}:liked:{postID} 1`

#### Synchronization

- **Immediate**: Updates to Redis (user actions)
- **Async**: Events published → Consumer syncs to PostgreSQL
- **Read**: Service reads from Redis, falls back to PostgreSQL if needed

### Cache Invalidation

- Not required (metrics are append-only)
- TTL: None (persistent until post deletion)

## 📁 File Storage

### AWS S3 Integration

**Purpose**: Store post and comment attachments

#### Flow

1. **Upload**:
   ```
   Client → Handler → Service → S3 (upload) → Get URL → Save metadata to PostgreSQL
   ```

2. **Retrieve**:
   ```
   Client → Handler → Service → Repository → PostgreSQL (get metadata) → Return URL
   ```

3. **Delete**:
   ```
   Client → Handler → Service → S3 (delete) → Repository → PostgreSQL (delete metadata)
   ```

#### Configuration

- **Region**: Configurable via `AWS_REGION`
- **Bucket**: Configurable via `AWS_BUCKET`
- **Access**: IAM credentials (via AWS SDK)

## 🔒 Security

### Authentication

- **Method**: JWT (JSON Web Tokens)
- **Algorithm**: HS256
- **Storage**: Stateless (token in Authorization header)
- **Refresh**: Token refresh endpoint (if implemented)

### Authorization

- **Post Operations**: Owner-only (check `user_id`)
- **Comment Operations**: Owner-only
- **Attachment Operations**: Owner-only

### Security Measures

1. **Password Hashing**: bcrypt (via `golang.org/x/crypto`)
2. **Input Validation**: go-playground/validator
3. **SQL Injection Prevention**: Parameterized queries (sqlx)
4. **CORS**: Configured in Fiber (if needed)
5. **Rate Limiting**: Can be added via middleware

## 📈 Scalability Considerations

### Current Architecture

- **Stateless API**: Can scale horizontally
- **Database**: Single PostgreSQL instance (can be replicated)
- **Redis**: Single instance (can be clustered)
- **Event Bus**: In-memory (GoChannel) - single instance limitation

### Scaling Strategies

#### 1. Horizontal Scaling

- **API Servers**: Multiple instances behind load balancer
- **Database**: Read replicas for read-heavy operations
- **Redis**: Redis Cluster for distributed caching

#### 2. Event Bus Upgrade

- **Current**: GoChannel (in-memory, single instance)
- **Production**: Switch to Redis Pub/Sub or Kafka
  ```go
  // Replace GoChannel with Redis Pub/Sub
  pubsub := redisstream.NewRedisStream(redisClient, ...)
  ```

#### 3. Database Optimization

- **Indexes**: Add indexes on frequently queried columns
- **Connection Pooling**: Configure sqlx connection pool
- **Read Replicas**: Separate read/write operations

#### 4. Caching Strategy

- **CDN**: For static assets (S3 + CloudFront)
- **Application Cache**: Redis for frequently accessed data
- **Query Cache**: Cache expensive queries

### Performance Optimizations

1. **Database**:
   - Use indexes on `user_id`, `post_id`, `is_active`
   - Batch operations where possible
   - Connection pooling

2. **Redis**:
   - Pipeline operations for bulk updates
   - Use appropriate data structures

3. **API**:
   - Response compression (Fiber built-in)
   - Pagination for list endpoints
   - Efficient serialization

## 🔍 Monitoring & Observability

### Recommended Additions

1. **Logging**: Structured logging (e.g., zap, logrus)
2. **Metrics**: Prometheus metrics endpoint
3. **Tracing**: OpenTelemetry for distributed tracing
4. **Health Checks**: `/health` endpoint
5. **Error Tracking**: Sentry or similar

## 📚 Additional Resources

- [Fiber Documentation](https://docs.gofiber.io/)
- [Watermill Documentation](https://watermill.io/)
- [sqlx Documentation](https://jmoiron.github.io/sqlx/)
- [PostgreSQL Best Practices](https://www.postgresql.org/docs/current/)
- [Redis Best Practices](https://redis.io/docs/manual/patterns/)

---

**Last Updated**: 2025-10-31  
**Version**: 1.0

