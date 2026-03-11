package commands

import (
	"fmt"

	"github.com/ashaju-godaddy/gdsnip/internal/cli/api"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/config"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/tui"
	"github.com/ashaju-godaddy/gdsnip/pkg/models"
	"github.com/spf13/cobra"
)

var (
	listLimit  int
	listOffset int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your snippets",
	Long: `List all snippets you have access to, including personal and team snippets.

Examples:
  gdsnip list                # List all your snippets
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

	// Fetch user snippets and team snippets with spinner
	var userSnippets *api.ListResponse
	var teams *api.TeamListResponse
	var allSnippets []*models.Snippet

	err = tui.RunWithSpinner("Loading your snippets...", func() error {
		var err error

		// Get user's personal snippets
		userSnippets, err = client.ListMySnippets(100, 0) // Get more to combine with team snippets
		if err != nil {
			return err
		}

		// Get user's teams
		teams, err = client.ListTeams(50, 0)
		if err != nil {
			return err
		}

		// Combine user snippets (convert to pointers)
		for i := range userSnippets.Snippets {
			allSnippets = append(allSnippets, &userSnippets.Snippets[i])
		}

		// Fetch team snippets for each team
		for _, team := range teams.Teams {
			teamSnippets, err := client.ListTeamSnippets(team.Slug, 100, 0)
			if err != nil {
				// Continue even if we can't fetch some team snippets
				continue
			}
			for i := range teamSnippets.Snippets {
				allSnippets = append(allSnippets, &teamSnippets.Snippets[i])
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to list snippets: %w", err)
	}

	// Display results
	if len(allSnippets) == 0 {
		fmt.Println()
		fmt.Println(tui.FormatInfo("You don't have any snippets yet"))
		fmt.Println()
		fmt.Println("To create a snippet, run:")
		fmt.Printf("  %s\n", tui.FormatCode("gdsnip push -f <file> -n <name>"))
		fmt.Println()
		return nil
	}

	// Apply pagination to combined results
	totalSnippets := len(allSnippets)
	if listOffset >= totalSnippets {
		fmt.Println()
		fmt.Println(tui.FormatInfo("No more snippets to display"))
		fmt.Println()
		return nil
	}

	end := listOffset + listLimit
	if end > totalSnippets {
		end = totalSnippets
	}

	displaySnippets := allSnippets[listOffset:end]

	// Display header
	fmt.Println()
	fmt.Printf("%s %s %d of %d\n",
		tui.FormatTitle("All Your Snippets"),
		tui.FormatDim("Showing"),
		len(displaySnippets),
		totalSnippets,
	)
	fmt.Println()

	// Display table
	for i, snippet := range displaySnippets {
		// Snippet path with visibility indicator
		path := fmt.Sprintf("%s/%s", snippet.Namespace, snippet.Slug)
		visibilityIcon := "🔒"
		if snippet.Visibility == "public" {
			visibilityIcon = "🌐"
		} else if snippet.Visibility == "team" {
			visibilityIcon = "👥"
		}

		fmt.Printf("%s %s %s\n",
			tui.FormatDim(fmt.Sprintf("%d.", listOffset+i+1)),
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
		if i < len(displaySnippets)-1 {
			fmt.Println()
		}
	}

	// Pagination hint
	if totalSnippets > len(displaySnippets) {
		fmt.Println()
		fmt.Println(tui.FormatDim(fmt.Sprintf("Showing %d of %d total snippets", len(displaySnippets), totalSnippets)))
		if totalSnippets > listOffset+listLimit {
			fmt.Println(tui.FormatDim("To see more, use: gdsnip list --offset " + fmt.Sprint(listOffset+listLimit)))
		}
	}

	fmt.Println()

	return nil
}
