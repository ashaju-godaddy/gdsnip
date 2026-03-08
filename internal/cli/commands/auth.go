package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ashaju-godaddy/gdsnip/internal/cli/api"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/config"
	"github.com/ashaju-godaddy/gdsnip/internal/cli/tui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  "Commands for authentication: register, login, logout, status",
}

var authRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new account",
	Long:  "Create a new GDSNIP account with email, username, and password",
	RunE:  runRegister,
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your account",
	Long:  "Authenticate with your email and password to get an access token",
	RunE:  runLogin,
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from your account",
	Long:  "Clear stored credentials and logout from the current session",
	RunE:  runLogout,
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  "Display information about the currently logged in user",
	RunE:  runStatus,
}

func init() {
	// Add subcommands to auth
	authCmd.AddCommand(authRegisterCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)

	// Register auth command with root
	AddCommand(authCmd)
}

func runRegister(cmd *cobra.Command, args []string) error {
	fmt.Println(tui.FormatTitle("Register New Account"))

	// Prompt for email
	email, err := tui.PromptString("Email", "", true)
	if err != nil {
		return fmt.Errorf("failed to get email: %w", err)
	}

	// Prompt for username
	username, err := tui.PromptString("Username", "", true)
	if err != nil {
		return fmt.Errorf("failed to get username: %w", err)
	}

	// Prompt for password
	fmt.Println()
	fmt.Print(tui.FormatLabel("Password"))
	fmt.Print(" " + tui.FormatError("*"))
	fmt.Print(" " + tui.FormatDim("(min 8 characters)"))
	fmt.Println()
	password, err := promptPassword()
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	// Confirm password
	fmt.Println()
	fmt.Print(tui.FormatLabel("Confirm Password"))
	fmt.Print(" " + tui.FormatError("*"))
	fmt.Println()
	confirmPassword, err := promptPassword()
	if err != nil {
		return fmt.Errorf("failed to confirm password: %w", err)
	}

	// Validate passwords match
	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Register with spinner
	var response *api.AuthResponse
	err = tui.RunWithSpinner("Creating account...", func() error {
		var registerErr error
		response, registerErr = client.Register(email, username, password)
		return registerErr
	})

	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	// Save credentials
	creds := &config.Credentials{
		Token:    response.Token,
		UserID:   response.User.ID,
		Username: response.User.Username,
		Email:    response.User.Email,
		SavedAt:  time.Now(),
	}

	if err := config.SaveCredentials(creds); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	// Display success
	fmt.Println()
	fmt.Println(tui.FormatSuccess("Account created successfully!"))
	fmt.Println()
	fmt.Printf("  %s %s\n", tui.FormatLabel("Username"), response.User.Username)
	fmt.Printf("  %s %s\n", tui.FormatLabel("Email"), response.User.Email)
	fmt.Println()
	fmt.Println(tui.FormatInfo("You are now logged in. Start by searching for snippets:"))
	fmt.Printf("  %s\n", tui.FormatCode("gdsnip search docker"))
	fmt.Println()

	return nil
}

func runLogin(cmd *cobra.Command, args []string) error {
	fmt.Println(tui.FormatTitle("Login"))

	// Prompt for email
	email, err := tui.PromptString("Email", "", true)
	if err != nil {
		return fmt.Errorf("failed to get email: %w", err)
	}

	// Prompt for password
	fmt.Println()
	fmt.Print(tui.FormatLabel("Password"))
	fmt.Print(" " + tui.FormatError("*"))
	fmt.Println()
	password, err := promptPassword()
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	// Create API client
	client, err := api.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Login with spinner
	var response *api.AuthResponse
	err = tui.RunWithSpinner("Logging in...", func() error {
		var loginErr error
		response, loginErr = client.Login(email, password)
		return loginErr
	})

	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Save credentials
	creds := &config.Credentials{
		Token:    response.Token,
		UserID:   response.User.ID,
		Username: response.User.Username,
		Email:    response.User.Email,
		SavedAt:  time.Now(),
	}

	if err := config.SaveCredentials(creds); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	// Display success
	fmt.Println()
	fmt.Println(tui.FormatSuccess("Logged in successfully!"))
	fmt.Println()
	fmt.Printf("  %s %s\n", tui.FormatLabel("Username"), response.User.Username)
	fmt.Printf("  %s %s\n", tui.FormatLabel("Email"), response.User.Email)
	fmt.Println()

	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	// Check if logged in
	if !config.IsLoggedIn() {
		fmt.Println(tui.FormatInfo("Not currently logged in"))
		return nil
	}

	// Get current user info
	creds, err := config.LoadCredentials()
	if err == nil {
		fmt.Printf("Logging out %s...\n", tui.FormatHighlight(creds.Username))
	}

	// Clear credentials
	if err := config.ClearCredentials(); err != nil {
		return fmt.Errorf("failed to clear credentials: %w", err)
	}

	fmt.Println(tui.FormatSuccess("Logged out successfully"))
	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Check if logged in
	if !config.IsLoggedIn() {
		fmt.Println(tui.FormatInfo("Not logged in"))
		fmt.Println()
		fmt.Println("To login, run:")
		fmt.Printf("  %s\n", tui.FormatCode("gdsnip auth login"))
		fmt.Println()
		fmt.Println("To register, run:")
		fmt.Printf("  %s\n", tui.FormatCode("gdsnip auth register"))
		fmt.Println()
		return nil
	}

	// Load credentials
	creds, err := config.LoadCredentials()
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	// Display status
	fmt.Println(tui.FormatTitle("Authentication Status"))
	fmt.Println()
	fmt.Printf("  %s %s\n", tui.FormatLabel("Status"), tui.FormatSuccess("Logged in"))
	fmt.Printf("  %s %s\n", tui.FormatLabel("Username"), creds.Username)
	fmt.Printf("  %s %s\n", tui.FormatLabel("Email"), creds.Email)
	fmt.Printf("  %s %s\n", tui.FormatLabel("User ID"), tui.FormatDim(creds.UserID))
	fmt.Printf("  %s %s\n", tui.FormatLabel("Logged in at"), tui.FormatDim(creds.SavedAt.Format("2006-01-02 15:04:05")))
	fmt.Println()

	return nil
}

// promptPassword is a helper for password prompts with hidden input
func promptPassword() (string, error) {
	fmt.Print("> ")

	// Read password with hidden input
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	fmt.Println() // New line after hidden input
	return strings.TrimSpace(string(bytePassword)), nil
}
