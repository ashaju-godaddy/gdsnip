package commands

import (
	"fmt"
	"strings"

	"github.com/ashaju-godaddy/gdsnip/internal/cli/api"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/tui"
	"github.com/spf13/cobra"
)

var (
	searchTags  []string
	searchLimit int
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for public snippets",
	Long: `Search for public snippets by query string and/or tags.

Examples:
  gdsnip search docker                    # Search for "docker"
  gdsnip search postgres -t database      # Search with tag filter
  gdsnip search -t docker -t compose      # Search by multiple tags
  gdsnip search "github actions" --limit 50  # Search with custom limit`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSearch,
}

func init() {
	searchCmd.Flags().StringSliceVarP(&searchTags, "tags", "t", []string{}, "Filter by tags (can be used multiple times)")
	searchCmd.Flags().IntVar(&searchLimit, "limit", 20, "Maximum number of results (default 20, max 100)")

	AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	// Get query from args
	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	// Validate that we have at least query or tags
	if query == "" && len(searchTags) == 0 {
		return fmt.Errorf("provide a search query or tags")
	}

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Search with spinner
	var response *api.SearchResponse
	searchMsg := "Searching"
	if query != "" {
		searchMsg += fmt.Sprintf(" for '%s'", query)
	}
	if len(searchTags) > 0 {
		searchMsg += fmt.Sprintf(" with tags [%s]", strings.Join(searchTags, ", "))
	}
	searchMsg += "..."

	err = tui.RunWithSpinner(searchMsg, func() error {
		var searchErr error
		response, searchErr = client.Search(query, searchTags, searchLimit)
		return searchErr
	})

	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Display results
	if len(response.Snippets) == 0 {
		fmt.Println()
		fmt.Println(tui.FormatInfo("No snippets found"))
		return nil
	}

	// Display header
	fmt.Println()
	fmt.Printf("%s %s %d of %d\n",
		tui.FormatTitle("Search Results:"),
		tui.FormatDim("Showing"),
		len(response.Snippets),
		response.Pagination.Total,
	)
	fmt.Println()

	// Display table
	for i, snippet := range response.Snippets {
		// Snippet path
		path := fmt.Sprintf("%s/%s", snippet.Namespace, snippet.Slug)
		fmt.Printf("%s %s\n",
			tui.FormatDim(fmt.Sprintf("%d.", i+1)),
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

		// Variables
		if len(snippet.Variables) > 0 {
			fmt.Printf("   %s %d variable(s)\n",
				tui.FormatDim("Variables:"),
				len(snippet.Variables),
			)
		}

		// Pull count
		fmt.Printf("   %s %d pulls\n",
			tui.FormatDim("Pulls:"),
			snippet.PullCount,
		)

		// Separator
		if i < len(response.Snippets)-1 {
			fmt.Println()
		}
	}

	// Show hint for more info
	fmt.Println()
	fmt.Println(tui.FormatDim("To see details, run:"))
	fmt.Printf("  %s\n", tui.FormatCode("gdsnip info <namespace>/<slug>"))
	fmt.Println()
	fmt.Println(tui.FormatDim("To pull a snippet, run:"))
	fmt.Printf("  %s\n", tui.FormatCode("gdsnip pull <namespace>/<slug>"))
	fmt.Println()

	return nil
}
