package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	// Slug pattern: lowercase alphanumeric with hyphens, must start and end with alphanumeric
	slugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)

	// Username pattern: lowercase alphanumeric with underscores, 3-50 chars
	usernamePattern = regexp.MustCompile(`^[a-z0-9_]{3,50}$`)

	// Basic email pattern (simplified, real validation should use net/mail or similar)
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// GenerateSlug generates a URL-safe slug from a name
func GenerateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' {
			result.WriteRune(r)
		}
	}
	slug = result.String()

	// Remove consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Limit length to 100 characters
	if len(slug) > 100 {
		slug = slug[:100]
		// Trim trailing hyphen if cutting created one
		slug = strings.TrimRight(slug, "-")
	}

	return slug
}

// ValidateSlug validates a slug format
func ValidateSlug(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	if len(slug) > 100 {
		return fmt.Errorf("slug must be at most 100 characters")
	}

	if len(slug) < 1 {
		return fmt.Errorf("slug must be at least 1 character")
	}

	// Check for consecutive hyphens
	if strings.Contains(slug, "--") {
		return fmt.Errorf("slug cannot contain consecutive hyphens")
	}

	// Single character slugs are allowed as long as they're alphanumeric
	if len(slug) == 1 {
		if !slugPattern.MatchString(slug + "x") { // Hack: add char to pass pattern check
			matched := regexp.MustCompile(`^[a-z0-9]$`).MatchString(slug)
			if !matched {
				return fmt.Errorf("slug must contain only lowercase letters, numbers, and hyphens")
			}
		}
		return nil
	}

	if !slugPattern.MatchString(slug) {
		return fmt.Errorf("slug must contain only lowercase letters, numbers, and hyphens, and must start and end with a letter or number")
	}

	return nil
}

// ValidateUsername validates a username format
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}

	if len(username) > 50 {
		return fmt.Errorf("username must be at most 50 characters")
	}

	if !usernamePattern.MatchString(username) {
		return fmt.Errorf("username must contain only lowercase letters, numbers, and underscores")
	}

	return nil
}

// ValidateEmail validates an email format (basic validation)
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if len(email) > 255 {
		return fmt.Errorf("email must be at most 255 characters")
	}

	if !emailPattern.MatchString(email) {
		return fmt.Errorf("email format is invalid")
	}

	return nil
}

// ValidatePassword validates password requirements
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must be at most 128 characters")
	}

	return nil
}

// ValidateVisibility validates snippet visibility values
func ValidateVisibility(visibility string) error {
	switch visibility {
	case "public", "private", "team":
		return nil
	default:
		return fmt.Errorf("visibility must be 'public', 'private', or 'team'")
	}
}

// ValidateTags validates snippet tags
func ValidateTags(tags []string) error {
	if len(tags) > 20 {
		return fmt.Errorf("snippets can have at most 20 tags")
	}

	for _, tag := range tags {
		if len(tag) == 0 {
			return fmt.Errorf("tags cannot be empty")
		}
		if len(tag) > 50 {
			return fmt.Errorf("tags must be at most 50 characters")
		}
	}

	return nil
}
