package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/ashaju-godaddy/gdsnip/internal/cli/api"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/tui"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <namespace>/<slug>",
	Short: "Show snippet details",
	Long: `Display detailed information about a snippet including variables, tags, and usage.

Examples:
  gdsnip info demo/docker-pg
  gdsnip info username/my-snippet`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func init() {
	AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	// Parse namespace/slug
	parts := strings.Split(args[0], "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid snippet path. Format: <namespace>/<slug>")
	}

	namespace := parts[0]
	slug := parts[1]

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Get snippet with spinner
	var snippet *models.Snippet
	err = tui.RunWithSpinner(fmt.Sprintf("Loading %s/%s...", namespace, slug), func() error {
		var getErr error
		snippet, getErr = client.GetSnippet(namespace, slug)
		return getErr
	})

	if err != nil {
		return fmt.Errorf("failed to get snippet: %w", err)
	}

	// Display snippet info
	fmt.Println()
	fmt.Println(tui.FormatTitle(snippet.Name))

	// Path
	path := fmt.Sprintf("%s/%s", snippet.Namespace, snippet.Slug)
	fmt.Printf("%s %s\n", tui.FormatLabel("Path"), tui.FormatHighlight(path))

	// Description
	if snippet.Description != "" {
		fmt.Printf("%s %s\n", tui.FormatLabel("Description"), snippet.Description)
	}

	// Visibility
	visibilityIcon := "🔒 Private"
	if snippet.Visibility == "public" {
		visibilityIcon = "🌐 Public"
	}
	fmt.Printf("%s %s\n", tui.FormatLabel("Visibility"), visibilityIcon)

	// Version and stats
	fmt.Printf("%s v%d  %s %d pulls\n",
		tui.FormatLabel("Version"),
		snippet.Version,
		tui.FormatDim("·"),
		snippet.PullCount,
	)

	// Dates
	fmt.Printf("%s %s\n", tui.FormatLabel("Created"), formatDate(snippet.CreatedAt))
	fmt.Printf("%s %s\n", tui.FormatLabel("Updated"), formatDate(snippet.UpdatedAt))

	// Tags
	if len(snippet.Tags) > 0 {
		fmt.Printf("%s ", tui.FormatLabel("Tags"))
		for i, tag := range snippet.Tags {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(tui.FormatCode(tag))
		}
		fmt.Println()
	}

	// Variables
	if len(snippet.Variables) > 0 {
		fmt.Println()
		fmt.Println(tui.FormatTitle("Variables"))
		fmt.Println()

		// Table header
		fmt.Printf("  %-20s %-10s %-15s %s\n",
			tui.FormatLabel("Name"),
			tui.FormatLabel("Required"),
			tui.FormatLabel("Default"),
			tui.FormatLabel("Description"),
		)
		fmt.Println(tui.FormatDim(strings.Repeat("─", 80)))

		// Table rows
		for _, v := range snippet.Variables {
			required := "No"
			if v.Required {
				required = tui.FormatError("Yes")
			} else {
				required = tui.FormatDim("No")
			}

			defaultVal := tui.FormatDim("-")
			if v.Default != "" {
				defaultVal = tui.FormatCode(v.Default)
			}

			desc := tui.FormatDim("-")
			if v.Description != "" {
				desc = v.Description
			}

			fmt.Printf("  %-20s %-10s %-15s %s\n",
				tui.FormatHighlight(v.Name),
				required,
				defaultVal,
				desc,
			)
		}
	}

	// Usage example
	fmt.Println()
	fmt.Println(tui.FormatTitle("Usage"))
	fmt.Println()
	fmt.Println(tui.FormatDim("Pull this snippet:"))
	pullCmd := fmt.Sprintf("gdsnip pull %s", path)
	fmt.Printf("  %s\n", tui.FormatCode(pullCmd))

	// Add variable examples if available
	if len(snippet.Variables) > 0 {
		fmt.Println()
		fmt.Println(tui.FormatDim("With variables:"))
		varExample := fmt.Sprintf("gdsnip pull %s", path)
		exampleCount := 0
		for _, v := range snippet.Variables {
			if exampleCount < 2 { // Show max 2 examples
				varExample += fmt.Sprintf(" --%s=value", v.Name)
				exampleCount++
			}
		}
		fmt.Printf("  %s\n", tui.FormatCode(varExample))
	}

	fmt.Println()

	return nil
}

// formatDate formats a time.Time to a human-readable string
func formatDate(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours < 1 {
			minutes := int(diff.Minutes())
			if minutes < 1 {
				return "just now"
			}
			return fmt.Sprintf("%d minute(s) ago", minutes)
		}
		return fmt.Sprintf("%d hour(s) ago", hours)
	}

	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d day(s) ago", days)
	}

	if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / 24 / 7)
		return fmt.Sprintf("%d week(s) ago", weeks)
	}

	return t.Format("2006-01-02")
}
