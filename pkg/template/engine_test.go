package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== ExtractVariables Tests ====================

func TestExtractVariables_MultipleVariables(t *testing.T) {
	content := `
	Database: {{DB_NAME}}
	Password: {{DB_PASSWORD}}
	Port: {{PORT}}
	`

	vars := ExtractVariables(content)

	assert.Len(t, vars, 3)
	assert.Contains(t, vars, "DB_NAME")
	assert.Contains(t, vars, "DB_PASSWORD")
	assert.Contains(t, vars, "PORT")
}

func TestExtractVariables_Deduplication(t *testing.T) {
	content := `
	{{DB_NAME}} is the database name
	Connect to {{DB_NAME}} on port {{PORT}}
	{{DB_NAME}} should be unique
	`

	vars := ExtractVariables(content)

	assert.Len(t, vars, 2)
	assert.Contains(t, vars, "DB_NAME")
	assert.Contains(t, vars, "PORT")
}

func TestExtractVariables_NoVariables(t *testing.T) {
	content := "This is plain text with no variables"

	vars := ExtractVariables(content)

	assert.Empty(t, vars)
}

func TestExtractVariables_InvalidPatterns(t *testing.T) {
	// lowercase should not match
	content1 := "{{lowercase}}"
	vars1 := ExtractVariables(content1)
	assert.Empty(t, vars1)

	// starting with number should not match
	content2 := "{{123START}}"
	vars2 := ExtractVariables(content2)
	assert.Empty(t, vars2)

	// hyphens should not match
	content3 := "{{MY-VAR}}"
	vars3 := ExtractVariables(content3)
	assert.Empty(t, vars3)
}

func TestExtractVariables_ValidPatterns(t *testing.T) {
	content := `
	{{SIMPLE}}
	{{WITH_UNDERSCORE}}
	{{WITH_123_NUMBERS}}
	{{A}}
	{{VAR_NAME_2}}
	`

	vars := ExtractVariables(content)

	assert.Len(t, vars, 5)
	assert.Contains(t, vars, "SIMPLE")
	assert.Contains(t, vars, "WITH_UNDERSCORE")
	assert.Contains(t, vars, "WITH_123_NUMBERS")
	assert.Contains(t, vars, "A")
	assert.Contains(t, vars, "VAR_NAME_2")
}

func TestExtractVariables_Ordering(t *testing.T) {
	content := "{{Z}} {{B}} {{A}} {{M}}"

	vars := ExtractVariables(content)

	// Should be sorted alphabetically
	assert.Equal(t, []string{"A", "B", "M", "Z"}, vars)
}

// ==================== Validate Tests ====================

func TestValidate_AllRequiredProvided(t *testing.T) {
	definitions := []Variable{
		{Name: "DB_NAME", Required: true},
		{Name: "DB_PASSWORD", Required: true},
		{Name: "PORT", Required: false, Default: "5432"},
	}

	provided := map[string]string{
		"DB_NAME":     "mydb",
		"DB_PASSWORD": "secret",
	}

	err := Validate(definitions, provided)

	assert.NoError(t, err)
}

func TestValidate_MissingRequired(t *testing.T) {
	definitions := []Variable{
		{Name: "DB_NAME", Required: true},
		{Name: "DB_PASSWORD", Required: true},
		{Name: "PORT", Required: true},
	}

	provided := map[string]string{
		"DB_NAME": "mydb",
		// Missing DB_PASSWORD and PORT
	}

	err := Validate(definitions, provided)

	require.Error(t, err)
	// Should list ALL missing variables, not just the first one
	assert.Contains(t, err.Error(), "DB_PASSWORD")
	assert.Contains(t, err.Error(), "PORT")
	assert.Contains(t, err.Error(), "missing required variables")
}

func TestValidate_DefaultsFillIn(t *testing.T) {
	definitions := []Variable{
		{Name: "DB_NAME", Required: true, Default: "defaultdb"},
		{Name: "PORT", Required: false, Default: "5432"},
	}

	provided := map[string]string{}

	// Should pass because DB_NAME has a default
	err := Validate(definitions, provided)

	assert.NoError(t, err)
}

func TestValidate_EmptyStringCountsAsMissing(t *testing.T) {
	definitions := []Variable{
		{Name: "DB_PASSWORD", Required: true},
	}

	provided := map[string]string{
		"DB_PASSWORD": "", // Empty string should count as missing
	}

	err := Validate(definitions, provided)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "DB_PASSWORD")
}

// ==================== Render Tests ====================

func TestRender_BasicSubstitution(t *testing.T) {
	content := "Database: {{DB_NAME}}, Port: {{PORT}}"
	definitions := []Variable{
		{Name: "DB_NAME", Required: true},
		{Name: "PORT", Required: true},
	}
	provided := map[string]string{
		"DB_NAME": "mydb",
		"PORT":    "5432",
	}

	result, err := Render(content, definitions, provided)

	require.NoError(t, err)
	assert.Equal(t, "Database: mydb, Port: 5432", result.Content)
	assert.Empty(t, result.Warnings)
}

func TestRender_WithDefaults(t *testing.T) {
	content := "Database: {{DB_NAME}}, Port: {{PORT}}"
	definitions := []Variable{
		{Name: "DB_NAME", Required: true},
		{Name: "PORT", Required: false, Default: "5432"},
	}
	provided := map[string]string{
		"DB_NAME": "mydb",
		// PORT not provided, should use default
	}

	result, err := Render(content, definitions, provided)

	require.NoError(t, err)
	assert.Equal(t, "Database: mydb, Port: 5432", result.Content)
}

func TestRender_OverrideDefaults(t *testing.T) {
	content := "Port: {{PORT}}"
	definitions := []Variable{
		{Name: "PORT", Required: false, Default: "5432"},
	}
	provided := map[string]string{
		"PORT": "3000", // Override default
	}

	result, err := Render(content, definitions, provided)

	require.NoError(t, err)
	assert.Equal(t, "Port: 3000", result.Content)
}

func TestRender_PartialVariables(t *testing.T) {
	content := "Required: {{REQUIRED}}, Optional: {{OPTIONAL}}"
	definitions := []Variable{
		{Name: "REQUIRED", Required: true},
		{Name: "OPTIONAL", Required: false, Default: "default-value"},
	}
	provided := map[string]string{
		"REQUIRED": "value",
		// OPTIONAL omitted, should use default
	}

	result, err := Render(content, definitions, provided)

	require.NoError(t, err)
	assert.Equal(t, "Required: value, Optional: default-value", result.Content)
}

func TestRender_ExtraVariablesWarning(t *testing.T) {
	content := "Database: {{DB_NAME}}"
	definitions := []Variable{
		{Name: "DB_NAME", Required: true},
	}
	provided := map[string]string{
		"DB_NAME": "mydb",
		"EXTRA1":  "unused1",
		"EXTRA2":  "unused2",
	}

	result, err := Render(content, definitions, provided)

	require.NoError(t, err)
	assert.Equal(t, "Database: mydb", result.Content)
	// Should warn about unused variables
	assert.NotEmpty(t, result.Warnings)
	assert.Contains(t, result.Warnings[0], "EXTRA")
}

func TestRender_MissingRequiredError(t *testing.T) {
	content := "Database: {{DB_NAME}}"
	definitions := []Variable{
		{Name: "DB_NAME", Required: true},
	}
	provided := map[string]string{} // Missing required variable

	result, err := Render(content, definitions, provided)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "DB_NAME")
}

func TestRender_PreservesNonMatchingBraces(t *testing.T) {
	// Go templates use {{.Field}} which shouldn't match our pattern
	content := "GDSNIP var: {{MY_VAR}}, Go template: {{.Field}}"
	definitions := []Variable{
		{Name: "MY_VAR", Required: true},
	}
	provided := map[string]string{
		"MY_VAR": "value",
	}

	result, err := Render(content, definitions, provided)

	require.NoError(t, err)
	// Should replace MY_VAR but leave {{.Field}} alone
	assert.Equal(t, "GDSNIP var: value, Go template: {{.Field}}", result.Content)
}

func TestRender_MultilineContent(t *testing.T) {
	content := `version: '3.8'
services:
  postgres:
    image: postgres:{{PG_VERSION}}
    environment:
      POSTGRES_PASSWORD: {{DB_PASSWORD}}
    ports:
      - "{{PORT}}:5432"`

	definitions := []Variable{
		{Name: "PG_VERSION", Required: false, Default: "15"},
		{Name: "DB_PASSWORD", Required: true},
		{Name: "PORT", Required: false, Default: "5432"},
	}
	provided := map[string]string{
		"DB_PASSWORD": "supersecret",
	}

	result, err := Render(content, definitions, provided)

	require.NoError(t, err)
	assert.Contains(t, result.Content, "postgres:15")
	assert.Contains(t, result.Content, "POSTGRES_PASSWORD: supersecret")
	assert.Contains(t, result.Content, `"5432:5432"`)
}

// ==================== MergeExtractedWithDefined Tests ====================

func TestMergeExtractedWithDefined_AllDefined(t *testing.T) {
	extracted := []string{"DB_NAME", "PORT"}
	defined := []Variable{
		{Name: "DB_NAME", Description: "Database name", Required: true},
		{Name: "PORT", Description: "Port number", Required: false, Default: "5432"},
	}

	merged := MergeExtractedWithDefined(extracted, defined)

	assert.Len(t, merged, 2)
	assert.Equal(t, "Database name", merged[0].Description)
	assert.Equal(t, "Port number", merged[1].Description)
}

func TestMergeExtractedWithDefined_PartialDefined(t *testing.T) {
	extracted := []string{"DB_NAME", "PORT", "PASSWORD"}
	defined := []Variable{
		{Name: "DB_NAME", Description: "Database name", Required: true},
		// PORT and PASSWORD not defined
	}

	merged := MergeExtractedWithDefined(extracted, defined)

	assert.Len(t, merged, 3)
	// DB_NAME should have description
	assert.Equal(t, "Database name", merged[0].Description)
	// PORT and PASSWORD should be basic variables
	assert.True(t, merged[1].Required) // Default to required
	assert.Empty(t, merged[1].Description)
	assert.True(t, merged[2].Required)
	assert.Empty(t, merged[2].Description)
}

func TestMergeExtractedWithDefined_NothingDefined(t *testing.T) {
	extracted := []string{"VAR1", "VAR2"}
	defined := []Variable{}

	merged := MergeExtractedWithDefined(extracted, defined)

	assert.Len(t, merged, 2)
	// All should be basic required variables
	for _, v := range merged {
		assert.True(t, v.Required)
		assert.Empty(t, v.Description)
	}
}

// ==================== GetMissingVariables Tests ====================

func TestGetMissingVariables_SomeMissing(t *testing.T) {
	definitions := []Variable{
		{Name: "DB_NAME", Description: "Database name", Required: true},
		{Name: "DB_PASSWORD", Description: "Database password", Required: true, Sensitive: true},
		{Name: "PORT", Required: false, Default: "5432"},
	}
	provided := map[string]string{
		"DB_NAME": "mydb",
		// DB_PASSWORD is missing
	}

	missing := GetMissingVariables(definitions, provided)

	assert.Len(t, missing, 1)
	assert.Equal(t, "DB_PASSWORD", missing[0].Name)
	assert.Equal(t, "Database password", missing[0].Description)
	assert.True(t, missing[0].Sensitive)
}

func TestGetMissingVariables_NoneMissing(t *testing.T) {
	definitions := []Variable{
		{Name: "DB_NAME", Required: true},
		{Name: "PORT", Required: false, Default: "5432"},
	}
	provided := map[string]string{
		"DB_NAME": "mydb",
	}

	missing := GetMissingVariables(definitions, provided)

	assert.Empty(t, missing)
}

func TestGetMissingVariables_AllMissing(t *testing.T) {
	definitions := []Variable{
		{Name: "VAR1", Required: true},
		{Name: "VAR2", Required: true},
		{Name: "VAR3", Required: true},
	}
	provided := map[string]string{}

	missing := GetMissingVariables(definitions, provided)

	assert.Len(t, missing, 3)
}
