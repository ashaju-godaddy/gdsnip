package service

import (
	"testing"
	"time"

	"github.com/ashaju-godaddy/gdsnip/internal/api/repository"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a test DB (requires running PostgreSQL)
func setupTestDB(t *testing.T) (*repository.UserRepo, func()) {
	// This test requires a running database
	// Skip if not available
	dbURL := "postgres://gdsnip:gdsnip@localhost:5432/gdsnip?sslmode=disable"

	db, err := repository.NewDB(dbURL)
	if err != nil {
		t.Skip("Database not available, skipping integration test")
		return nil, nil
	}

	// Clean up tables (order matters due to foreign keys)
	_, _ = db.Exec("DELETE FROM snippets")
	_, _ = db.Exec("DELETE FROM users")

	userRepo := repository.NewUserRepo(db)

	cleanup := func() {
		_, _ = db.Exec("DELETE FROM snippets")
		_, _ = db.Exec("DELETE FROM users")
		db.Close()
	}

	return userRepo, cleanup
}

func TestAuthService_Register(t *testing.T) {
	userRepo, cleanup := setupTestDB(t)
	if cleanup == nil {
		return
	}
	defer cleanup()

	authService := NewAuthService(userRepo, "test-secret-key-that-is-very-long", 24*time.Hour)

	t.Run("successful registration", func(t *testing.T) {
		user, token, err := authService.Register("test@example.com", "testuser", "password123")
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.NotEmpty(t, token)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "testuser", user.Username)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.PasswordHash)
	})

	t.Run("duplicate email", func(t *testing.T) {
		_, _, err := authService.Register("test@example.com", "anotheruser", "password123")
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "DUPLICATE_EMAIL", apiErr.Code)
	})

	t.Run("duplicate username", func(t *testing.T) {
		_, _, err := authService.Register("another@example.com", "testuser", "password123")
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "DUPLICATE_USERNAME", apiErr.Code)
	})

	t.Run("invalid email", func(t *testing.T) {
		_, _, err := authService.Register("invalid-email", "newuser", "password123")
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "VALIDATION_ERROR", apiErr.Code)
	})

	t.Run("invalid username - too short", func(t *testing.T) {
		_, _, err := authService.Register("valid@example.com", "ab", "password123")
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "VALIDATION_ERROR", apiErr.Code)
	})

	t.Run("invalid password - too short", func(t *testing.T) {
		_, _, err := authService.Register("valid@example.com", "validuser", "short")
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "VALIDATION_ERROR", apiErr.Code)
	})
}

func TestAuthService_Login(t *testing.T) {
	userRepo, cleanup := setupTestDB(t)
	if cleanup == nil {
		return
	}
	defer cleanup()

	authService := NewAuthService(userRepo, "test-secret-key-that-is-very-long", 24*time.Hour)

	// Register a user first
	email := "login@example.com"
	username := "loginuser"
	password := "password123"
	_, _, err := authService.Register(email, username, password)
	require.NoError(t, err)

	t.Run("successful login", func(t *testing.T) {
		user, token, err := authService.Login(email, password)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.NotEmpty(t, token)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, username, user.Username)
	})

	t.Run("wrong password", func(t *testing.T) {
		_, _, err := authService.Login(email, "wrongpassword")
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "UNAUTHORIZED", apiErr.Code)
	})

	t.Run("non-existent email", func(t *testing.T) {
		_, _, err := authService.Login("nonexistent@example.com", password)
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "UNAUTHORIZED", apiErr.Code)
	})
}

func TestAuthService_JWT(t *testing.T) {
	userRepo, cleanup := setupTestDB(t)
	if cleanup == nil {
		return
	}
	defer cleanup()

	jwtSecret := "test-secret-key-that-is-very-long"
	authService := NewAuthService(userRepo, jwtSecret, 24*time.Hour)

	t.Run("generate and parse JWT", func(t *testing.T) {
		userID := "test-user-id"
		username := "testuser"

		token, err := authService.generateJWT(userID, username)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := authService.ParseJWT(token)
		require.NoError(t, err)
		require.NotNil(t, claims)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, username, claims.Username)
	})

	t.Run("parse invalid token", func(t *testing.T) {
		_, err := authService.ParseJWT("invalid-token")
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "INVALID_TOKEN", apiErr.Code)
	})

	t.Run("parse expired token", func(t *testing.T) {
		// Create service with very short expiry
		shortExpiryService := NewAuthService(userRepo, jwtSecret, 1*time.Nanosecond)
		token, err := shortExpiryService.generateJWT("test-id", "testuser")
		require.NoError(t, err)

		// Wait for token to expire
		time.Sleep(2 * time.Millisecond)

		_, err = shortExpiryService.ParseJWT(token)
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok)
		assert.Equal(t, "TOKEN_EXPIRED", apiErr.Code)
	})
}

func TestAuthService_GetCurrentUser(t *testing.T) {
	userRepo, cleanup := setupTestDB(t)
	if cleanup == nil {
		return
	}
	defer cleanup()

	authService := NewAuthService(userRepo, "test-secret-key-that-is-very-long", 24*time.Hour)

	// Register a user
	user, _, err := authService.Register("current@example.com", "currentuser", "password123")
	require.NoError(t, err)

	t.Run("get existing user", func(t *testing.T) {
		fetchedUser, err := authService.GetCurrentUser(user.ID)
		require.NoError(t, err)
		require.NotNil(t, fetchedUser)
		assert.Equal(t, user.ID, fetchedUser.ID)
		assert.Equal(t, user.Email, fetchedUser.Email)
		assert.Equal(t, user.Username, fetchedUser.Username)
	})

	t.Run("get non-existent user", func(t *testing.T) {
		// Use a valid UUID format that doesn't exist
		_, err := authService.GetCurrentUser("00000000-0000-0000-0000-000000000000")
		require.Error(t, err)
		apiErr, ok := err.(*models.APIError)
		require.True(t, ok, "error type: %T, value: %v", err, err)
		assert.Equal(t, "USER_NOT_FOUND", apiErr.Code)
	})
}

func TestPasswordHashing(t *testing.T) {
	userRepo, cleanup := setupTestDB(t)
	if cleanup == nil {
		return
	}
	defer cleanup()

	authService := NewAuthService(userRepo, "test-secret-key-that-is-very-long", 24*time.Hour)

	password := "my-secure-password"

	t.Run("hash and check password", func(t *testing.T) {
		hash, err := authService.hashPassword(password)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)

		// Check correct password
		assert.True(t, authService.checkPassword(hash, password))

		// Check wrong password
		assert.False(t, authService.checkPassword(hash, "wrong-password"))
	})
}
