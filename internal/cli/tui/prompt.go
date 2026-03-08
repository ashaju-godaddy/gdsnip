package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ashaju-godaddy/gdsnip/pkg/template"
	"golang.org/x/term"
)

// PromptVariable prompts the user for a variable value
func PromptVariable(v template.Variable) (string, error) {
	// Display variable name and description
	fmt.Println()
	fmt.Print(FormatLabel(v.Name))

	if v.Description != "" {
		fmt.Print(" " + FormatDim(v.Description))
	}

	// Show default value if available
	if v.Default != "" {
		fmt.Print(" " + FormatDim(fmt.Sprintf("[default: %s]", v.Default)))
	}

	// Show required indicator
	if v.Required && v.Default == "" {
		fmt.Print(" " + FormatError("*"))
	}

	fmt.Println()

	// Check if this is a password field (common convention: contains PASSWORD in name)
	isPassword := strings.Contains(strings.ToUpper(v.Name), "PASSWORD") ||
		strings.Contains(strings.ToUpper(v.Name), "SECRET") ||
		strings.Contains(strings.ToUpper(v.Name), "TOKEN") ||
		strings.Contains(strings.ToUpper(v.Name), "KEY")

	var value string
	var err error

	if isPassword {
		value, err = promptPassword()
	} else {
		value, err = promptText()
	}

	if err != nil {
		return "", err
	}

	// Use default if no value provided
	if value == "" && v.Default != "" {
		value = v.Default
		fmt.Println(FormatDim(fmt.Sprintf("  Using default: %s", v.Default)))
	}

	// Validate required fields
	if v.Required && value == "" && v.Default == "" {
		return "", fmt.Errorf("value is required for %s", v.Name)
	}

	return value, nil
}

// promptText prompts for regular text input
func promptText() (string, error) {
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(value), nil
}

// promptPassword prompts for password input (hidden)
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

// Confirm prompts for yes/no confirmation
func Confirm(message string) bool {
	fmt.Println()
	fmt.Print(FormatLabel(message) + " " + FormatDim("[y/N]") + " ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// PromptFilename prompts for a filename with a suggestion
func PromptFilename(suggestion string) (string, error) {
	fmt.Println()
	fmt.Print(FormatLabel("Output filename"))

	if suggestion != "" {
		fmt.Print(" " + FormatDim(fmt.Sprintf("[%s]", suggestion)))
	}

	fmt.Println()
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	filename, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read filename: %w", err)
	}

	filename = strings.TrimSpace(filename)

	// Use suggestion if no input provided
	if filename == "" && suggestion != "" {
		filename = suggestion
		fmt.Println(FormatDim(fmt.Sprintf("  Using: %s", suggestion)))
	}

	return filename, nil
}

// PromptString prompts for a generic string input
func PromptString(label, defaultValue string, required bool) (string, error) {
	fmt.Println()
	fmt.Print(FormatLabel(label))

	if defaultValue != "" {
		fmt.Print(" " + FormatDim(fmt.Sprintf("[%s]", defaultValue)))
	}

	if required && defaultValue == "" {
		fmt.Print(" " + FormatError("*"))
	}

	fmt.Println()
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	value = strings.TrimSpace(value)

	// Use default if no value provided
	if value == "" && defaultValue != "" {
		value = defaultValue
		fmt.Println(FormatDim(fmt.Sprintf("  Using: %s", defaultValue)))
	}

	// Validate required fields
	if required && value == "" && defaultValue == "" {
		return "", fmt.Errorf("%s is required", label)
	}

	return value, nil
}
