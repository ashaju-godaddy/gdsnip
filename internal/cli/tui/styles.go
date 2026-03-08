package tui

import "github.com/charmbracelet/lipgloss"

// Color palette
var (
	ColorPrimary   = lipgloss.Color("#00ADD8") // Go blue
	ColorSuccess   = lipgloss.Color("#00C851") // Green
	ColorError     = lipgloss.Color("#FF4444") // Red
	ColorWarning   = lipgloss.Color("#FFBB33") // Yellow
	ColorInfo      = lipgloss.Color("#33B5E5") // Light blue
	ColorDim       = lipgloss.Color("#666666") // Gray
	ColorHighlight = lipgloss.Color("#FF6F00") // Orange
)

// Text styles
var (
	// TitleStyle for command titles and headers
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	// SuccessStyle for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	// ErrorStyle for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	// WarningStyle for warning messages
	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	// InfoStyle for informational messages
	InfoStyle = lipgloss.NewStyle().
			Foreground(ColorInfo)

	// DimStyle for secondary/muted text
	DimStyle = lipgloss.NewStyle().
			Foreground(ColorDim)

	// CodeStyle for code snippets and file names
	CodeStyle = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Background(lipgloss.Color("#1A1A1A")).
			Padding(0, 1)

	// LabelStyle for field labels
	LabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	// HighlightStyle for highlighted text
	HighlightStyle = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Bold(true)
)

// Box styles
var (
	// BoxStyle for bordered containers
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			MarginBottom(1)

	// SuccessBoxStyle for success messages
	SuccessBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSuccess).
			Padding(1, 2).
			MarginBottom(1)

	// ErrorBoxStyle for error messages
	ErrorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorError).
			Padding(1, 2).
			MarginBottom(1)
)

// Table styles
var (
	// TableHeaderStyle for table headers
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPrimary).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(ColorDim)

	// TableCellStyle for table cells
	TableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)
)

// Helper functions for common formatting

// FormatTitle formats a title with consistent styling
func FormatTitle(text string) string {
	return TitleStyle.Render(text)
}

// FormatSuccess formats a success message
func FormatSuccess(text string) string {
	return SuccessStyle.Render("✓ " + text)
}

// FormatError formats an error message
func FormatError(text string) string {
	return ErrorStyle.Render("✗ " + text)
}

// FormatWarning formats a warning message
func FormatWarning(text string) string {
	return WarningStyle.Render("⚠ " + text)
}

// FormatInfo formats an informational message
func FormatInfo(text string) string {
	return InfoStyle.Render("ℹ " + text)
}

// FormatCode formats code or file paths
func FormatCode(text string) string {
	return CodeStyle.Render(text)
}

// FormatLabel formats a field label
func FormatLabel(text string) string {
	return LabelStyle.Render(text + ":")
}

// FormatDim formats dimmed/secondary text
func FormatDim(text string) string {
	return DimStyle.Render(text)
}

// FormatHighlight formats highlighted text
func FormatHighlight(text string) string {
	return HighlightStyle.Render(text)
}
