package models

import "time"

// User represents a user in the system
type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never expose password hash in JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserPublic is a safe representation of User for API responses (no sensitive fields)
type UserPublic struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// Public converts a User to UserPublic, removing sensitive fields
func (u *User) Public() UserPublic {
	return UserPublic{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	}
}

// UserMinimal is a minimal user representation for embedding in other models
type UserMinimal struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Minimal returns a minimal user representation
func (u *User) Minimal() UserMinimal {
	return UserMinimal{
		ID:       u.ID,
		Username: u.Username,
	}
}
