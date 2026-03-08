package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set via ldflags during build
	Version = "dev"
	// BuildTime is set via ldflags during build
	BuildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "gdsnip",
	Short: "GDSNIP - A CLI-first snippet registry with variable substitution",
	Long: `GDSNIP is a command-line tool for storing, sharing, and pulling
parameterized code templates. Think of it as "npm for boilerplate code"
with variable substitution.

Store templates with {{VARIABLE}} placeholders and render them with
values on pull.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("gdsnip %s (built: %s)\n", Version, BuildTime)
		},
	})

	// Disable automatic completion command (we'll add it post-MVP)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// AddCommand adds a subcommand to the root command
func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

// GetRootCmd returns the root command (useful for testing)
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// handleError handles command errors with consistent formatting
func handleError(cmd *cobra.Command, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
