package commands

import (
	"fmt"

	"github.com/ashaju-godaddy/gdsnip/internal/cli/api"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/config"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/tui"
	"github.com/spf13/cobra"
)

var (
	listLimit  int
	listOffset int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your snippets",
	Long: `List all snippets owned by the current user.

Examples:
  gdsnip list                # List your snippets
  gdsnip list --limit 50     # List with custom limit`,
	RunE: runList,
}

func init() {
	listCmd.Flags().IntVar(&listLimit, "limit", 20, "Maximum number of results (default 20)")
	listCmd.Flags().IntVar(&listOffset, "offset", 0, "Offset for pagination")

	AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	// Check if logged in
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// List snippets with spinner
	var response *api.ListResponse
	err = tui.RunWithSpinner("Loading your snippets...", func() error {
		var listErr error
		response, listErr = client.ListMySnippets(listLimit, listOffset)
		return listErr
	})

	if err != nil {
		return fmt.Errorf("failed to list snippets: %w", err)
	}

	// Display results
	if len(response.Snippets) == 0 {
		fmt.Println()
		fmt.Println(tui.FormatInfo("You don't have any snippets yet"))
		fmt.Println()
		fmt.Println("To create a snippet, run:")
		fmt.Printf("  %s\n", tui.FormatCode("gdsnip push -f <file> -n <name>"))
		fmt.Println()
		return nil
	}

	// Get current user
	creds, _ := config.LoadCredentials()

	// Display header
	fmt.Println()
	fmt.Printf("%s %s %d of %d\n",
		tui.FormatTitle(fmt.Sprintf("%s's Snippets", creds.Username)),
		tui.FormatDim("Showing"),
		len(response.Snippets),
		response.Pagination.Total,
	)
	fmt.Println()

	// Display table
	for i, snippet := range response.Snippets {
		// Snippet path with visibility indicator
		path := fmt.Sprintf("%s/%s", snippet.Namespace, snippet.Slug)
		visibilityIcon := "🔒"
		if snippet.Visibility == "public" {
			visibilityIcon = "🌐"
		}

		fmt.Printf("%s %s %s\n",
			tui.FormatDim(fmt.Sprintf("%d.", i+1)),
			visibilityIcon,
			tui.FormatHighlight(path),
		)

		// Name and description
		fmt.Printf("   %s", snippet.Name)
		if snippet.Description != "" {
			fmt.Printf(" %s %s", tui.FormatDim("·"), tui.FormatDim(snippet.Description))
		}
		fmt.Println()

		// Tags
		if len(snippet.Tags) > 0 {
			fmt.Printf("   %s ", tui.FormatDim("Tags:"))
			for j, tag := range snippet.Tags {
				if j > 0 {
					fmt.Print(", ")
				}
				fmt.Print(tui.FormatCode(tag))
			}
			fmt.Println()
		}

		// Stats
		fmt.Printf("   %s %s  %s %d pulls  %s v%d\n",
			tui.FormatDim("Visibility:"),
			snippet.Visibility,
			tui.FormatDim("·"),
			snippet.PullCount,
			tui.FormatDim("·"),
			snippet.Version,
		)

		// Separator
		if i < len(response.Snippets)-1 {
			fmt.Println()
		}
	}

	// Pagination hint
	if response.Pagination.Total > len(response.Snippets) {
		fmt.Println()
		fmt.Println(tui.FormatDim(fmt.Sprintf("Showing %d of %d total snippets", len(response.Snippets), response.Pagination.Total)))
		if response.Pagination.Total > listOffset+listLimit {
			fmt.Println(tui.FormatDim("To see more, use: gdsnip list --offset " + fmt.Sprint(listOffset+listLimit)))
		}
	}

	fmt.Println()

	return nil
}
