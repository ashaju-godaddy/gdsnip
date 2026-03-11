package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/ashaju-godaddy/gdsnip/internal/cli/api"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/config"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/tui"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/ashaju-godaddy/gdsnip/pkg/template"
	"github.com/ashaju-godaddy/gdsnip/pkg/validator"
	"github.com/spf13/cobra"
)

var (
	pushFile        string
	pushName        string
	pushSlug        string
	pushDescription string
	pushTags        []string
	pushPublic      bool
	pushTeam        string
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Create a new snippet",
	Long: `Create a new snippet from a file. Variables in {{UPPERCASE}} format
will be automatically extracted.

Examples:
  gdsnip push -f template.yml -n "Docker Postgres"
  gdsnip push -f template.yml -n "My Template" --public
  gdsnip push -f template.yml -n "My Template" -s custom-slug -t docker -t database
  gdsnip push -f template.yml -n "Private Template" -d "My private template"
  gdsnip push -f template.yml -n "Shared Config" --team platform`,
	RunE: runPush,
}

func init() {
	pushCmd.Flags().StringVarP(&pushFile, "file", "f", "", "Template file to push (required)")
	pushCmd.Flags().StringVarP(&pushName, "name", "n", "", "Snippet name (required)")
	pushCmd.Flags().StringVarP(&pushSlug, "slug", "s", "", "Custom slug (auto-generated if not provided)")
	pushCmd.Flags().StringVarP(&pushDescription, "description", "d", "", "Snippet description")
	pushCmd.Flags().StringSliceVarP(&pushTags, "tags", "t", []string{}, "Tags (can be used multiple times)")
	pushCmd.Flags().BoolVar(&pushPublic, "public", false, "Make snippet public (default: private)")
	pushCmd.Flags().StringVar(&pushTeam, "team", "", "Push snippet to a team namespace")

	pushCmd.MarkFlagRequired("file")
	pushCmd.MarkFlagRequired("name")

	AddCommand(pushCmd)
}

func runPush(cmd *cobra.Command, args []string) error {
	// Check if logged in
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	// Read file content
	content, err := os.ReadFile(pushFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)

	// Extract variables
	fmt.Println()
	fmt.Println(tui.FormatTitle("Creating Snippet"))
	fmt.Println()

	variableNames := template.ExtractVariables(contentStr)

	// Display extracted variables
	if len(variableNames) > 0 {
		fmt.Printf("%s Found %d variable(s):\n", tui.FormatInfo("ℹ"), len(variableNames))
		for _, name := range variableNames {
			fmt.Printf("  %s %s\n", tui.FormatDim("·"), tui.FormatHighlight(name))
		}
		fmt.Println()
	} else {
		fmt.Println(tui.FormatWarning("No variables found in template"))
		fmt.Println()
	}

	// Auto-generate slug if not provided
	slug := pushSlug
	if slug == "" {
		slug = validator.GenerateSlug(pushName)
		fmt.Printf("%s Auto-generated slug: %s\n", tui.FormatDim("ℹ"), tui.FormatCode(slug))
		fmt.Println()
	}

	// Determine visibility
	visibility := "private"
	if pushPublic {
		visibility = "public"
	} else if pushTeam != "" {
		visibility = "team" // Default to team visibility for team snippets
	}

	// Display summary
	fmt.Println(tui.FormatLabel("Summary:"))
	fmt.Printf("  %s %s\n", tui.FormatDim("Name:"), pushName)
	fmt.Printf("  %s %s\n", tui.FormatDim("Slug:"), tui.FormatCode(slug))
	if pushDescription != "" {
		fmt.Printf("  %s %s\n", tui.FormatDim("Description:"), pushDescription)
	}
	if pushTeam != "" {
		fmt.Printf("  %s %s\n", tui.FormatDim("Team:"), tui.FormatHighlight(pushTeam))
	}
	fmt.Printf("  %s %s\n", tui.FormatDim("Visibility:"), visibility)
	if len(pushTags) > 0 {
		fmt.Printf("  %s %s\n", tui.FormatDim("Tags:"), strings.Join(pushTags, ", "))
	}
	fmt.Printf("  %s %d variable(s)\n", tui.FormatDim("Variables:"), len(variableNames))
	fmt.Println()

	// Confirm
	if !tui.Confirm("Create this snippet?") {
		fmt.Println(tui.FormatInfo("Cancelled"))
		return nil
	}

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Create snippet request
	req := &api.CreateSnippetRequest{
		Name:        pushName,
		Slug:        slug,
		Description: pushDescription,
		Content:     contentStr,
		Tags:        pushTags,
		Visibility:  visibility,
		TeamSlug:    pushTeam,
	}

	// Create snippet with spinner
	var snippet *models.Snippet
	err = tui.RunWithSpinner("Creating snippet...", func() error {
		var createErr error
		snippet, createErr = client.CreateSnippet(req)
		return createErr
	})

	if err != nil {
		return fmt.Errorf("failed to create snippet: %w", err)
	}

	// Get namespace for display
	var namespace string
	if pushTeam != "" {
		namespace = pushTeam
	} else {
		creds, _ := config.LoadCredentials()
		namespace = creds.Username
	}
	path := fmt.Sprintf("%s/%s", namespace, snippet.Slug)

	// Display success
	fmt.Println()
	fmt.Println(tui.FormatSuccess("Snippet created successfully!"))
	fmt.Println()
	fmt.Printf("  %s %s\n", tui.FormatLabel("Path"), tui.FormatHighlight(path))
	fmt.Printf("  %s %s\n", tui.FormatLabel("Visibility"), snippet.Visibility)
	fmt.Println()
	fmt.Println(tui.FormatInfo("To pull this snippet, run:"))
	pullCmdStr := fmt.Sprintf("gdsnip pull %s", path)

	// Add example variables if available
	if len(variableNames) > 0 && len(variableNames) <= 2 {
		for _, name := range variableNames {
			pullCmdStr += fmt.Sprintf(" --%s=value", name)
		}
	}

	fmt.Printf("  %s\n", tui.FormatCode(pullCmdStr))
	fmt.Println()

	return nil
}
