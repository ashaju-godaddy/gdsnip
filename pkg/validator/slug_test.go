package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== GenerateSlug Tests ====================

func TestGenerateSlug_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "Hello World", "hello-world"},
		{"with spaces", "Docker PostgreSQL Setup", "docker-postgresql-setup"},
		{"with underscores", "my_awesome_snippet", "my-awesome-snippet"},
		{"mixed case", "MyAwesomeSnippet", "myawesomesnippet"},
		{"with special chars", "Hello@World!", "helloworld"},
		{"multiple spaces", "Hello    World", "hello-world"},
		{"trailing spaces", "Hello World  ", "hello-world"},
		{"leading spaces", "  Hello World", "hello-world"},
		{"consecutive hyphens", "Hello---World", "hello-world"},
		{"numbers", "Project 123", "project-123"},
		{"single word", "Docker", "docker"},
		{"with numbers and underscores", "API_v2_Setup", "api-v2-setup"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateSlug_LongName(t *testing.T) {
	longName := "This is a very long snippet name that exceeds the maximum allowed length for a slug and should be truncated to exactly one hundred characters"
	slug := GenerateSlug(longName)

	assert.LessOrEqual(t, len(slug), 100)
	assert.NotEmpty(t, slug)
	// Should not end with hyphen
	assert.NotEqual(t, "-", slug[len(slug)-1:])
}

func TestGenerateSlug_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"only special chars", "@#$%", ""},
		{"only spaces", "   ", ""},
		{"only hyphens", "---", ""},
		{"single char", "a", "a"},
		{"single number", "1", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ==================== ValidateSlug Tests ====================

func TestValidateSlug_Valid(t *testing.T) {
	validSlugs := []string{
		"hello-world",
		"my-snippet",
		"docker-pg",
		"k8s-deployment",
		"api-v2",
		"123-test",
		"test-123",
		"a",
		"1",
		"my-really-long-slug-with-many-words",
	}

	for _, slug := range validSlugs {
		t.Run(slug, func(t *testing.T) {
			err := ValidateSlug(slug)
			assert.NoError(t, err)
		})
	}
}

func TestValidateSlug_Invalid(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		hasError bool
	}{
		{"empty", "", true},
		{"uppercase", "Hello-World", true},
		{"space", "hello world", true},
		{"underscore", "hello_world", true},
		{"special char", "hello@world", true},
		{"starts with hyphen", "-hello", true},
		{"ends with hyphen", "hello-", true},
		{"consecutive hyphens", "hello--world", true},
		{"too long", string(make([]byte, 101)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSlug(tt.slug)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== ValidateUsername Tests ====================

func TestValidateUsername_Valid(t *testing.T) {
	validUsernames := []string{
		"john",
		"john_doe",
		"john123",
		"alice_in_wonderland",
		"user_123",
		"abc",
	}

	for _, username := range validUsernames {
		t.Run(username, func(t *testing.T) {
			err := ValidateUsername(username)
			assert.NoError(t, err)
		})
	}
}

func TestValidateUsername_Invalid(t *testing.T) {
	tests := []struct {
		name     string
		username string
		hasError bool
	}{
		{"empty", "", true},
		{"too short", "ab", true},
		{"too long", string(make([]byte, 51)), true},
		{"uppercase", "JohnDoe", true},
		{"space", "john doe", true},
		{"hyphen", "john-doe", true},
		{"special char", "john@doe", true},
		{"starts with number", "123john", false}, // This is actually valid
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== ValidateEmail Tests ====================

func TestValidateEmail_Valid(t *testing.T) {
	validEmails := []string{
		"user@example.com",
		"john.doe@example.com",
		"user+tag@example.co.uk",
		"user_name@example.com",
		"123@example.com",
	}

	for _, email := range validEmails {
		t.Run(email, func(t *testing.T) {
			err := ValidateEmail(email)
			assert.NoError(t, err)
		})
	}
}

func TestValidateEmail_Invalid(t *testing.T) {
	invalidEmails := []string{
		"",
		"invalid",
		"@example.com",
		"user@",
		"user@.com",
		"user @example.com",
		"user@example",
	}

	for _, email := range invalidEmails {
		t.Run(email, func(t *testing.T) {
			err := ValidateEmail(email)
			assert.Error(t, err)
		})
	}
}

// ==================== ValidatePassword Tests ====================

func TestValidatePassword_Valid(t *testing.T) {
	validPasswords := []string{
		"password123",
		"MySecureP@ssw0rd!",
		"12345678",
		"a1b2c3d4e5f6",
	}

	for _, password := range validPasswords {
		t.Run("valid password", func(t *testing.T) {
			err := ValidatePassword(password)
			assert.NoError(t, err)
		})
	}
}

func TestValidatePassword_Invalid(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"empty", ""},
		{"too short", "pass"},
		{"7 chars", "1234567"},
		{"too long", string(make([]byte, 129))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			assert.Error(t, err)
		})
	}
}

// ==================== ValidateVisibility Tests ====================

func TestValidateVisibility(t *testing.T) {
	tests := []struct {
		name       string
		visibility string
		wantErr    bool
	}{
		{"public", "public", false},
		{"private", "private", false},
		{"team", "team", false},
		{"invalid", "invalid", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVisibility(tt.visibility)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ==================== ValidateTags Tests ====================

func TestValidateTags_Valid(t *testing.T) {
	validTags := [][]string{
		{},
		{"docker"},
		{"docker", "postgres", "database"},
		{"kubernetes", "k8s", "deployment"},
	}

	for _, tags := range validTags {
		t.Run("valid tags", func(t *testing.T) {
			err := ValidateTags(tags)
			assert.NoError(t, err)
		})
	}
}

func TestValidateTags_Invalid(t *testing.T) {
	tests := []struct {
		name string
		tags []string
	}{
		{"too many", make([]string, 21)},
		{"empty tag", []string{"docker", "", "postgres"}},
		{"too long tag", []string{string(make([]byte, 51))}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill "too many" test case with valid tags
			if tt.name == "too many" {
				for i := range tt.tags {
					tt.tags[i] = "tag"
				}
			}
			err := ValidateTags(tt.tags)
			assert.Error(t, err)
		})
	}
}
