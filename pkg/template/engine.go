package template

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Variable pattern: {{UPPERCASE_NAME}}
// Must start with a letter, can contain letters, numbers, and underscores
var variablePattern = regexp.MustCompile(`\{\{([A-Z][A-Z0-9_]*)\}\}`)

// Variable represents a template variable with metadata
type Variable struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
	Sensitive   bool   `json:"sensitive,omitempty"`
}

// RenderResult contains the rendered content and any warnings
type RenderResult struct {
	Content  string   `json:"content"`
	Warnings []string `json:"warnings,omitempty"`
}

// ExtractVariables finds all {{VAR}} patterns in template content.
// Returns a deduplicated, sorted list of variable names.
func ExtractVariables(content string) []string {
	matches := variablePattern.FindAllStringSubmatch(content, -1)

	// Use map for deduplication
	seen := make(map[string]bool)
	var vars []string

	for _, match := range matches {
		if len(match) > 1 {
			varName := match[1]
			if !seen[varName] {
				seen[varName] = true
				vars = append(vars, varName)
			}
		}
	}

	// Sort for consistent ordering
	sort.Strings(vars)

	return vars
}

// Validate checks that all required variables are provided.
// Returns an error listing ALL missing required variables (not just the first one).
func Validate(definitions []Variable, provided map[string]string) error {
	var missing []string

	for _, v := range definitions {
		if v.Required {
			val, exists := provided[v.Name]
			// Required variables must be present and non-empty, unless they have a default
			if (!exists || val == "") && v.Default == "" {
				missing = append(missing, v.Name)
			}
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required variables: %s", strings.Join(missing, ", "))
	}

	return nil
}

// Render substitutes variables in the template content.
// Returns RenderResult with the rendered content and any warnings.
func Render(content string, definitions []Variable, provided map[string]string) (*RenderResult, error) {
	// Validate first
	if err := Validate(definitions, provided); err != nil {
		return nil, err
	}

	// Build lookup map with defaults
	values := make(map[string]string)

	// First, apply defaults
	for _, def := range definitions {
		if def.Default != "" {
			values[def.Name] = def.Default
		}
	}

	// Then, apply provided values (overrides defaults)
	for k, v := range provided {
		values[k] = v
	}

	var warnings []string
	usedVars := make(map[string]bool)

	// Replace variables in content
	result := variablePattern.ReplaceAllStringFunc(content, func(match string) string {
		// Extract variable name (strip {{ and }})
		varName := match[2 : len(match)-2]
		usedVars[varName] = true

		if val, ok := values[varName]; ok {
			return val
		}
		// If not found in values, leave it as-is (shouldn't happen if Validate passed)
		return match
	})

	// Check for unused provided variables
	for k := range provided {
		if !usedVars[k] {
			warnings = append(warnings, fmt.Sprintf("variable %s was provided but not used in template", k))
		}
	}

	return &RenderResult{
		Content:  result,
		Warnings: warnings,
	}, nil
}

// MergeExtractedWithDefined merges auto-extracted variables with user-defined metadata.
// This is used when pushing snippets - we auto-extract variables from content,
// then merge with any metadata the user provides.
func MergeExtractedWithDefined(extracted []string, defined []Variable) []Variable {
	// Create a map of defined variables for quick lookup
	definedMap := make(map[string]Variable)
	for _, v := range defined {
		definedMap[v.Name] = v
	}

	var result []Variable

	// For each extracted variable, use defined metadata if available, else create basic variable
	for _, name := range extracted {
		if def, ok := definedMap[name]; ok {
			// Use user-defined metadata
			result = append(result, def)
		} else {
			// Create basic variable with just the name
			result = append(result, Variable{
				Name:     name,
				Required: true, // Default to required for safety
			})
		}
	}

	return result
}

// GetMissingVariables returns a list of Variable definitions for missing required variables.
// This is useful for CLI to prompt users interactively for missing variables.
func GetMissingVariables(definitions []Variable, provided map[string]string) []Variable {
	var missing []Variable

	for _, v := range definitions {
		if v.Required {
			val, exists := provided[v.Name]
			if (!exists || val == "") && v.Default == "" {
				missing = append(missing, v)
			}
		}
	}

	return missing
}
