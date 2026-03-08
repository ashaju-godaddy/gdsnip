package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ashaju-godaddy/gdsnip/internal/cli/api"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/tui"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/ashaju-godaddy/gdsnip/pkg/template"
	"github.com/spf13/cobra"
)

var (
	pullOutput  string
	pullVarFile string
)

var pullCmd = &cobra.Command{
	Use:   "pull <namespace>/<slug>",
	Short: "Pull and render a snippet",
	Long: `Pull a snippet and render it with variables. Variables can be provided via:
  - Command-line flags: --VAR_NAME=value
  - Variable file: --var-file vars.json
  - Interactive prompts: For missing required variables

Examples:
  gdsnip pull demo/docker-pg --DB_PASSWORD=secret
  gdsnip pull demo/docker-pg --var-file vars.json
  gdsnip pull demo/docker-pg --DB_PASSWORD=secret -o output.yml
  gdsnip pull demo/docker-pg  # Interactive prompts for variables`,
	Args:              cobra.MinimumNArgs(1), // At least snippet path required
	DisableFlagParsing: true, // We'll parse flags manually
	RunE:              runPull,
}

func init() {
	// Note: We don't register flags here because we parse them manually
	// This allows us to accept dynamic --UPPERCASE=value flags
	AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) error {
	// Parse arguments and flags manually
	if len(args) < 1 {
		return fmt.Errorf("snippet path is required. Usage: gdsnip pull <namespace>/<slug>")
	}

	snippetPath := ""
	variables := make(map[string]string)
	output := ""
	varFile := ""

	// Parse all arguments
	for i, arg := range args {
		if i == 0 && !strings.HasPrefix(arg, "-") {
			// First arg without dash is the snippet path
			snippetPath = arg
			continue
		}

		// Parse flags
		if strings.HasPrefix(arg, "--") {
			flag := strings.TrimPrefix(arg, "--")
			parts := strings.SplitN(flag, "=", 2)

			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]

				// Check if this is a known flag or a variable
				switch key {
				case "output", "o":
					output = value
				case "var-file":
					varFile = value
				default:
					// This is a variable flag (must be UPPERCASE)
					if isUppercaseVariable(key) {
						variables[key] = value
					} else {
						return fmt.Errorf("invalid flag: --%s (variable flags must be UPPERCASE)", key)
					}
				}
			} else {
				// Flag without value
				switch flag {
				case "help", "h":
					cmd.Help()
					return nil
				default:
					return fmt.Errorf("flag --%s requires a value", flag)
				}
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) == 2 {
			// Short flags
			flag := arg[1:]
			switch flag {
			case "o":
				// Next arg should be the output file
				if i+1 < len(args) {
					output = args[i+1]
				} else {
					return fmt.Errorf("flag -o requires a value")
				}
			case "h":
				cmd.Help()
				return nil
			default:
				return fmt.Errorf("unknown flag: -%s", flag)
			}
		}
	}

	// Validate snippet path
	if snippetPath == "" {
		return fmt.Errorf("snippet path is required")
	}

	// Parse namespace/slug
	parts := strings.Split(snippetPath, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid snippet path. Format: <namespace>/<slug>")
	}

	namespace := parts[0]
	slug := parts[1]

	// Load variables from file if provided
	if varFile != "" {
		fileVars, err := loadVariablesFromFile(varFile)
		if err != nil {
			return fmt.Errorf("failed to load variables from file: %w", err)
		}

		// Merge file variables (command-line flags take precedence)
		for key, value := range fileVars {
			if _, exists := variables[key]; !exists {
				variables[key] = value
			}
		}
	}

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Pull snippet
	fmt.Println()
	fmt.Println(tui.FormatTitle(fmt.Sprintf("Pulling %s/%s", namespace, slug)))
	fmt.Println()

	var response *api.PullResponse
	var pullErr error

	err = tui.RunWithSpinner("Loading snippet...", func() error {
		response, pullErr = client.PullSnippet(namespace, slug, variables)
		return pullErr
	})

	// Check if we have missing variables error
	if pullErr != nil {
		apiErr, ok := pullErr.(*models.APIError)
		if ok && apiErr.Code == "MISSING_VARIABLES" {
			// Extract missing variables
			missingVars := extractMissingVariables(apiErr)

			if len(missingVars) > 0 {
				fmt.Println()
				fmt.Println(tui.FormatWarning(fmt.Sprintf("Missing %d required variable(s)", len(missingVars))))
				fmt.Println()

				// Prompt for missing variables
				for _, v := range missingVars {
					value, promptErr := tui.PromptVariable(v)
					if promptErr != nil {
						return fmt.Errorf("failed to get variable %s: %w", v.Name, promptErr)
					}
					variables[v.Name] = value
				}

				// Retry pull with complete variables
				fmt.Println()
				err = tui.RunWithSpinner("Rendering template...", func() error {
					response, pullErr = client.PullSnippet(namespace, slug, variables)
					return pullErr
				})

				if pullErr != nil {
					return fmt.Errorf("pull failed: %w", pullErr)
				}
			}
		} else {
			return fmt.Errorf("pull failed: %w", pullErr)
		}
	}

	// Display warnings if any
	if len(response.Warnings) > 0 {
		fmt.Println()
		for _, warning := range response.Warnings {
			fmt.Println(tui.FormatWarning(warning))
		}
	}

	// Display success
	fmt.Println()
	fmt.Println(tui.FormatSuccess("Snippet rendered successfully!"))
	fmt.Println()

	// Determine output file
	outputFile := output
	if outputFile == "" {
		// Suggest filename based on snippet slug
		suggestion := slug
		if ext := filepath.Ext(response.Snippet.Name); ext != "" {
			suggestion += ext
		}

		// Ask user if they want to save
		if tui.Confirm("Save to file?") {
			filename, err := tui.PromptFilename(suggestion)
			if err != nil {
				return fmt.Errorf("failed to get filename: %w", err)
			}
			outputFile = filename
		}
	}

	// Save to file or display content
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(response.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Println()
		fmt.Printf("%s Saved to: %s\n", tui.FormatSuccess("✓"), tui.FormatCode(outputFile))
		fmt.Println()
	} else {
		// Display content
		fmt.Println(tui.FormatTitle("Rendered Content:"))
		fmt.Println()
		fmt.Println(tui.BoxStyle.Render(response.Content))
		fmt.Println()
	}

	// Display snippet info
	fmt.Println(tui.FormatDim(fmt.Sprintf("Snippet: %s/%s (v%d, %d pulls)",
		response.Snippet.Namespace,
		response.Snippet.Slug,
		response.Snippet.Version,
		response.Snippet.PullCount,
	)))
	fmt.Println()

	return nil
}

// isUppercaseVariable checks if a flag name is a valid uppercase variable
func isUppercaseVariable(name string) bool {
	if len(name) == 0 {
		return false
	}

	for _, ch := range name {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return false
		}
	}

	// Must start with a letter
	firstChar := rune(name[0])
	return firstChar >= 'A' && firstChar <= 'Z'
}

// loadVariablesFromFile loads variables from a JSON file
func loadVariablesFromFile(filename string) (map[string]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var variables map[string]string
	if err := json.Unmarshal(data, &variables); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return variables, nil
}

// extractMissingVariables extracts missing variables from API error details
func extractMissingVariables(apiErr *models.APIError) []template.Variable {
	if apiErr.Details == nil {
		return nil
	}

	details, ok := apiErr.Details.(map[string]interface{})
	if !ok {
		return nil
	}

	missing, ok := details["missing"]
	if !ok {
		return nil
	}

	// The missing field is a slice of variables
	missingSlice, ok := missing.([]interface{})
	if !ok {
		return nil
	}

	var variables []template.Variable
	for _, item := range missingSlice {
		varMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		v := template.Variable{
			Name: getStringField(varMap, "name"),
			Description: getStringField(varMap, "description"),
			Default: getStringField(varMap, "default"),
			Required: getBoolField(varMap, "required"),
		}

		variables = append(variables, v)
	}

	return variables
}

// Helper functions to extract fields from interface{} maps
func getStringField(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getBoolField(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
