package commands

import (
	"fmt"
	"strings"

	"github.com/ashaju-godaddy/gdsnip/internal/cli/api"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/config"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/tui"
	"github.com/spf13/cobra"
)

var (
	teamDescription string
	teamSlug        string
	memberRole      string
)

var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Team management commands",
	Long:  "Commands for managing teams: create, list, info, delete, and membership",
}

var teamsCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new team",
	Long: `Create a new team. You will become the owner.

Examples:
  gdsnip teams create "Platform Team"
  gdsnip teams create "Platform Team" --slug=platform
  gdsnip teams create "Platform Team" --description="Shared infrastructure templates"`,
	Args: cobra.ExactArgs(1),
	RunE: runTeamsCreate,
}

var teamsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your teams",
	Long:  "List all teams you are a member of",
	RunE:  runTeamsList,
}

var teamsInfoCmd = &cobra.Command{
	Use:   "info <slug>",
	Short: "Show team details",
	Long:  "Display detailed information about a team",
	Args:  cobra.ExactArgs(1),
	RunE:  runTeamsInfo,
}

var teamsDeleteCmd = &cobra.Command{
	Use:   "delete <slug>",
	Short: "Delete a team",
	Long:  "Delete a team (owner only). This cannot be undone.",
	Args:  cobra.ExactArgs(1),
	RunE:  runTeamsDelete,
}

var teamsMembersCmd = &cobra.Command{
	Use:   "members <slug>",
	Short: "List team members",
	Long:  "Display all members of a team with their roles",
	Args:  cobra.ExactArgs(1),
	RunE:  runTeamsMembers,
}

var teamsAddCmd = &cobra.Command{
	Use:   "add <team-slug> <username>",
	Short: "Add a member to a team",
	Long: `Add a user to a team with a specified role.

Roles:
  admin   - Can manage members and all snippets
  member  - Can create snippets and pull team snippets (default)
  viewer  - Can only pull team snippets

Examples:
  gdsnip teams add platform johndoe
  gdsnip teams add platform johndoe --role=admin`,
	Args: cobra.ExactArgs(2),
	RunE: runTeamsAdd,
}

var teamsRemoveCmd = &cobra.Command{
	Use:   "remove <team-slug> <username>",
	Short: "Remove a member from a team",
	Long:  "Remove a user from a team. Cannot remove the owner.",
	Args:  cobra.ExactArgs(2),
	RunE:  runTeamsRemove,
}

var teamsRoleCmd = &cobra.Command{
	Use:   "role <team-slug> <username> <role>",
	Short: "Change a member's role",
	Long: `Change a team member's role.

Roles: admin, member, viewer

Examples:
  gdsnip teams role platform johndoe admin
  gdsnip teams role platform johndoe viewer`,
	Args: cobra.ExactArgs(3),
	RunE: runTeamsRole,
}

var teamsLeaveCmd = &cobra.Command{
	Use:   "leave <slug>",
	Short: "Leave a team",
	Long:  "Leave a team you are a member of. Owners cannot leave their teams.",
	Args:  cobra.ExactArgs(1),
	RunE:  runTeamsLeave,
}

var teamsSnippetsCmd = &cobra.Command{
	Use:   "snippets <slug>",
	Short: "List team snippets",
	Long:  "Display all snippets belonging to a team",
	Args:  cobra.ExactArgs(1),
	RunE:  runTeamsSnippets,
}

func init() {
	// Create command flags
	teamsCreateCmd.Flags().StringVar(&teamSlug, "slug", "", "Custom slug (auto-generated if not provided)")
	teamsCreateCmd.Flags().StringVar(&teamDescription, "description", "", "Team description")

	// Add member flag
	teamsAddCmd.Flags().StringVar(&memberRole, "role", "member", "Role for the new member (admin, member, viewer)")

	// Add subcommands
	teamsCmd.AddCommand(teamsCreateCmd)
	teamsCmd.AddCommand(teamsListCmd)
	teamsCmd.AddCommand(teamsInfoCmd)
	teamsCmd.AddCommand(teamsDeleteCmd)
	teamsCmd.AddCommand(teamsMembersCmd)
	teamsCmd.AddCommand(teamsAddCmd)
	teamsCmd.AddCommand(teamsRemoveCmd)
	teamsCmd.AddCommand(teamsRoleCmd)
	teamsCmd.AddCommand(teamsLeaveCmd)
	teamsCmd.AddCommand(teamsSnippetsCmd)

	AddCommand(teamsCmd)
}

func runTeamsCreate(cmd *cobra.Command, args []string) error {
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	req := &api.CreateTeamRequest{
		Name:        args[0],
		Slug:        teamSlug,
		Description: teamDescription,
	}

	var team *api.Team
	err = tui.RunWithSpinner("Creating team...", func() error {
		var createErr error
		team, createErr = client.CreateTeam(req)
		return createErr
	})

	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	fmt.Println()
	fmt.Println(tui.FormatSuccess("Team created successfully!"))
	fmt.Println()
	fmt.Printf("  %s %s\n", tui.FormatLabel("Name"), team.Name)
	fmt.Printf("  %s %s\n", tui.FormatLabel("Slug"), tui.FormatCode(team.Slug))
	fmt.Printf("  %s %s\n", tui.FormatLabel("Role"), "owner")
	fmt.Println()
	fmt.Println(tui.FormatInfo("To push a snippet to this team:"))
	fmt.Printf("  %s\n", tui.FormatCode(fmt.Sprintf("gdsnip push -f template.yml -n \"My Template\" --team %s", team.Slug)))
	fmt.Println()

	return nil
}

func runTeamsList(cmd *cobra.Command, args []string) error {
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	var response *api.TeamListResponse
	err = tui.RunWithSpinner("Loading teams...", func() error {
		var listErr error
		response, listErr = client.ListTeams(20, 0)
		return listErr
	})

	if err != nil {
		return fmt.Errorf("failed to list teams: %w", err)
	}

	if len(response.Teams) == 0 {
		fmt.Println()
		fmt.Println(tui.FormatInfo("You're not a member of any teams yet"))
		fmt.Println()
		fmt.Println("To create a team, run:")
		fmt.Printf("  %s\n", tui.FormatCode("gdsnip teams create <name>"))
		fmt.Println()
		return nil
	}

	fmt.Println()
	fmt.Printf("%s %s %d of %d\n",
		tui.FormatTitle("Your Teams"),
		tui.FormatDim("Showing"),
		len(response.Teams),
		response.Pagination.Total,
	)
	fmt.Println()

	for i, team := range response.Teams {
		fmt.Printf("%s %s\n",
			tui.FormatDim(fmt.Sprintf("%d.", i+1)),
			tui.FormatHighlight(team.Slug),
		)
		fmt.Printf("   %s", team.Name)
		if team.Description != "" {
			fmt.Printf(" %s %s", tui.FormatDim("-"), tui.FormatDim(team.Description))
		}
		fmt.Println()
		fmt.Printf("   %s %s  %s %d members\n",
			tui.FormatDim("Role:"),
			team.Role,
			tui.FormatDim("|"),
			team.MemberCount,
		)
		if i < len(response.Teams)-1 {
			fmt.Println()
		}
	}
	fmt.Println()

	return nil
}

func runTeamsInfo(cmd *cobra.Command, args []string) error {
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	var team *api.Team
	err = tui.RunWithSpinner("Loading team...", func() error {
		var getErr error
		team, getErr = client.GetTeam(args[0])
		return getErr
	})

	if err != nil {
		return fmt.Errorf("failed to get team: %w", err)
	}

	fmt.Println()
	fmt.Println(tui.FormatTitle(team.Name))
	fmt.Println()
	fmt.Printf("  %s %s\n", tui.FormatLabel("Slug"), tui.FormatCode(team.Slug))
	if team.Description != "" {
		fmt.Printf("  %s %s\n", tui.FormatLabel("Description"), team.Description)
	}
	fmt.Printf("  %s %s\n", tui.FormatLabel("Your Role"), team.Role)
	fmt.Printf("  %s %d\n", tui.FormatLabel("Members"), team.MemberCount)
	fmt.Printf("  %s %s\n", tui.FormatLabel("Created"), tui.FormatDim(team.CreatedAt.Format("2006-01-02")))
	fmt.Println()

	return nil
}

func runTeamsDelete(cmd *cobra.Command, args []string) error {
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	slug := args[0]

	fmt.Println()
	fmt.Println(tui.FormatWarning(fmt.Sprintf("You are about to delete team '%s'", slug)))
	fmt.Println(tui.FormatDim("This will also delete all team snippets. This action cannot be undone."))
	fmt.Println()

	if !tui.Confirm("Are you sure you want to delete this team?") {
		fmt.Println(tui.FormatInfo("Cancelled"))
		return nil
	}

	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	err = tui.RunWithSpinner("Deleting team...", func() error {
		return client.DeleteTeam(slug)
	})

	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	fmt.Println()
	fmt.Println(tui.FormatSuccess("Team deleted successfully"))
	fmt.Println()

	return nil
}

func runTeamsMembers(cmd *cobra.Command, args []string) error {
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	var members []api.TeamMember
	err = tui.RunWithSpinner("Loading members...", func() error {
		var listErr error
		members, listErr = client.ListTeamMembers(args[0])
		return listErr
	})

	if err != nil {
		return fmt.Errorf("failed to list members: %w", err)
	}

	fmt.Println()
	fmt.Printf("%s %d members\n",
		tui.FormatTitle(fmt.Sprintf("Team: %s", args[0])),
		len(members),
	)
	fmt.Println()

	for _, member := range members {
		roleDisplay := member.Role
		switch member.Role {
		case "owner":
			roleDisplay = tui.FormatHighlight("owner")
		case "admin":
			roleDisplay = tui.FormatCode("admin")
		}
		fmt.Printf("  %s %s %s\n",
			tui.FormatHighlight(member.Username),
			tui.FormatDim("-"),
			roleDisplay,
		)
	}
	fmt.Println()

	return nil
}

func runTeamsAdd(cmd *cobra.Command, args []string) error {
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	teamSlugArg := args[0]
	username := args[1]

	// Validate role
	role := strings.ToLower(memberRole)
	if role != "admin" && role != "member" && role != "viewer" {
		return fmt.Errorf("invalid role '%s'; must be admin, member, or viewer", memberRole)
	}

	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	err = tui.RunWithSpinner(fmt.Sprintf("Adding %s to team...", username), func() error {
		return client.AddTeamMember(teamSlugArg, username, role)
	})

	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	fmt.Println()
	fmt.Println(tui.FormatSuccess(fmt.Sprintf("Added %s to team as %s", username, role)))
	fmt.Println()

	return nil
}

func runTeamsRemove(cmd *cobra.Command, args []string) error {
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	teamSlugArg := args[0]
	username := args[1]

	if !tui.Confirm(fmt.Sprintf("Remove %s from team %s?", username, teamSlugArg)) {
		fmt.Println(tui.FormatInfo("Cancelled"))
		return nil
	}

	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	err = tui.RunWithSpinner(fmt.Sprintf("Removing %s from team...", username), func() error {
		return client.RemoveTeamMember(teamSlugArg, username)
	})

	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	fmt.Println()
	fmt.Println(tui.FormatSuccess(fmt.Sprintf("Removed %s from team", username)))
	fmt.Println()

	return nil
}

func runTeamsRole(cmd *cobra.Command, args []string) error {
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	teamSlugArg := args[0]
	username := args[1]
	role := strings.ToLower(args[2])

	if role != "admin" && role != "member" && role != "viewer" {
		return fmt.Errorf("invalid role '%s'; must be admin, member, or viewer", role)
	}

	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	err = tui.RunWithSpinner(fmt.Sprintf("Updating %s's role...", username), func() error {
		return client.UpdateTeamMemberRole(teamSlugArg, username, role)
	})

	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	fmt.Println()
	fmt.Println(tui.FormatSuccess(fmt.Sprintf("Updated %s's role to %s", username, role)))
	fmt.Println()

	return nil
}

func runTeamsLeave(cmd *cobra.Command, args []string) error {
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	teamSlugArg := args[0]

	if !tui.Confirm(fmt.Sprintf("Leave team %s?", teamSlugArg)) {
		fmt.Println(tui.FormatInfo("Cancelled"))
		return nil
	}

	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	err = tui.RunWithSpinner("Leaving team...", func() error {
		return client.LeaveTeam(teamSlugArg)
	})

	if err != nil {
		return fmt.Errorf("failed to leave team: %w", err)
	}

	fmt.Println()
	fmt.Println(tui.FormatSuccess("Left team successfully"))
	fmt.Println()

	return nil
}

func runTeamsSnippets(cmd *cobra.Command, args []string) error {
	if !config.IsLoggedIn() {
		return fmt.Errorf("not logged in. Run 'gdsnip auth login' first")
	}

	teamSlug := args[0]

	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	var response *api.ListResponse
	err = tui.RunWithSpinner("Loading team snippets...", func() error {
		var listErr error
		response, listErr = client.ListTeamSnippets(teamSlug, 20, 0)
		return listErr
	})

	if err != nil {
		return fmt.Errorf("failed to list team snippets: %w", err)
	}

	if len(response.Snippets) == 0 {
		fmt.Println()
		fmt.Println(tui.FormatInfo(fmt.Sprintf("Team '%s' doesn't have any snippets yet", teamSlug)))
		fmt.Println()
		fmt.Println("To push a snippet to this team, run:")
		fmt.Printf("  %s\n", tui.FormatCode(fmt.Sprintf("gdsnip push -f <file> -n <name> --team %s", teamSlug)))
		fmt.Println()
		return nil
	}

	// Display header
	fmt.Println()
	fmt.Printf("%s %s %d of %d\n",
		tui.FormatTitle(fmt.Sprintf("Team: %s", teamSlug)),
		tui.FormatDim("Showing"),
		len(response.Snippets),
		response.Pagination.Total,
	)
	fmt.Println()

	// Display snippets
	for i, snippet := range response.Snippets {
		// Snippet path with visibility indicator
		path := fmt.Sprintf("%s/%s", snippet.Namespace, snippet.Slug)
		visibilityIcon := "👥"
		if snippet.Visibility == "public" {
			visibilityIcon = "🌐"
		} else if snippet.Visibility == "private" {
			visibilityIcon = "🔒"
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
	}

	fmt.Println()

	return nil
}
