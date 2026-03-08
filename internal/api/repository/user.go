package repository

import (
	"database/sql"
	"fmt"

	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// UserRepo handles user database operations
type UserRepo struct {
	db *sqlx.DB
}

// NewUserRepo creates a new user repository
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create inserts a new user into the database
func (r *UserRepo) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, email, username, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`

	_, err := r.db.Exec(query, user.ID, user.Email, user.Username, user.PasswordHash)
	if err != nil {
		// Check for unique constraint violations
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if pqErr.Constraint == "users_email_key" {
					return models.NewDuplicateEmailError(user.Email)
				}
				if pqErr.Constraint == "users_username_key" {
					return models.NewDuplicateUsernameError(user.Username)
				}
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepo) GetByID(id string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewUserNotFoundError(id)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewUserNotFoundError(email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	err := r.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NewUserNotFoundError(username)
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}
