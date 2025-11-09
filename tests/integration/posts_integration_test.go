package integration

import (
	"context"
	"mpb/internal/posts"
	"mpb/tests/testutils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	database := testutils.SetupTestDB(t)
	defer database.Conn.Close()

	redisClient := testutils.SetupTestRedis(t)
	defer redisClient.Close()
	defer testutils.CleanupRedis(t, redisClient.Client)

	publisher, subscriber := testutils.SetupTestPubSub(t)
	logger := testutils.SetupTestLogger(t)

	postRepo := posts.NewPostsRepository(database)
	metricsService := posts.NewMetricsService(redisClient.Client, publisher, logger)
	postService := posts.NewPostsService(postRepo, metricsService, publisher, logger)
	postsHandler := posts.NewPostsHandlers(postService, metricsService)

	app := fiber.New()
	api := app.Group("/api")
	postsGroup := api.Group("/posts")

	postsGroup.Get("/", postsHandler.GetAllPosts)
	postsGroup.Get("/:id", postsHandler.GetPost)

	t.Run("Create and Get Post", func(t *testing.T) {
		ctx := context.Background()
		userID := 1

		post, err := postService.CreatePost(ctx, userID, "Integration Test Post", "This is a test description", "test")
		require.NoError(t, err)
		assert.NotZero(t, post.ID)
		assert.Equal(t, userID, post.UserID)
		assert.Equal(t, "Integration Test Post", post.Title)

		retrievedPost, err := postService.GetPostByID(ctx, post.ID)
		require.NoError(t, err)
		assert.Equal(t, post.ID, retrievedPost.ID)
		assert.Equal(t, post.Title, retrievedPost.Title)
	})

	t.Run("List Posts", func(t *testing.T) {
		ctx := context.Background()

		posts, err := postService.ListPosts(ctx, posts.PostFilter{OnlyActive: true})
		require.NoError(t, err)
		assert.NotNil(t, posts)
		assert.Greater(t, len(posts), 0)
	})

	t.Run("Like Post", func(t *testing.T) {
		ctx := context.Background()
		userID := 1
		postID := 1

		err := metricsService.LikePost(ctx, userID, postID)
		require.NoError(t, err)

		likes, err := metricsService.GetLikes(ctx, postID)
		require.NoError(t, err)
		assert.Greater(t, likes, 0)

		liked, err := metricsService.IsLiked(ctx, userID, postID)
		require.NoError(t, err)
		assert.True(t, liked)
	})

	t.Run("Unlike Post", func(t *testing.T) {
		ctx := context.Background()
		userID := 1
		postID := 1

		err := metricsService.UnlikePost(ctx, userID, postID)
		require.NoError(t, err)

		liked, err := metricsService.IsLiked(ctx, userID, postID)
		require.NoError(t, err)
		assert.False(t, liked)
	})

	_ = subscriber
}

func TestPostsHandlerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	database := testutils.SetupTestDB(t)
	defer database.Conn.Close()

	redisClient := testutils.SetupTestRedis(t)
	defer redisClient.Close()
	defer testutils.CleanupRedis(t, redisClient.Client)

	publisher, _ := testutils.SetupTestPubSub(t)
	logger := testutils.SetupTestLogger(t)

	postRepo := posts.NewPostsRepository(database)
	metricsService := posts.NewMetricsService(redisClient.Client, publisher, logger)
	postService := posts.NewPostsService(postRepo, metricsService, publisher, logger)
	postsHandler := posts.NewPostsHandlers(postService, metricsService)

	app := fiber.New()
	api := app.Group("/api")
	postsGroup := api.Group("/posts")

	postsGroup.Get("/", postsHandler.GetAllPosts)
	postsGroup.Get("/:id", postsHandler.GetPost)

	t.Run("GET /api/posts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET /api/posts/:id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/posts/1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
