// Package ui implements the Bubbletea TUI for Shepherd.
package ui

import "github.com/charmbracelet/lipgloss"

// ─── Catppuccin Macchiato Palette ───────────────────────────────────────────
// https://github.com/catppuccin/catppuccin

var (
	ColorBase     = lipgloss.Color("#24273a")
	ColorSurface0 = lipgloss.Color("#363a4f")
	ColorSurface1 = lipgloss.Color("#494d64")
	ColorOverlay0 = lipgloss.Color("#6e738d")
	ColorText     = lipgloss.Color("#cad3f5")
	ColorSubtext0 = lipgloss.Color("#a5adcb")
	ColorSubtext1 = lipgloss.Color("#b8c0e0")
	ColorLavender = lipgloss.Color("#b7bdf8")
	ColorBlue     = lipgloss.Color("#8aadf4")
	ColorSapphire = lipgloss.Color("#7dc4e4")
	ColorGreen    = lipgloss.Color("#a6da95")
	ColorYellow   = lipgloss.Color("#eed49f")
	ColorPeach    = lipgloss.Color("#f5a97f")
	ColorRed      = lipgloss.Color("#ed8796")
	ColorMauve    = lipgloss.Color("#c6a0f6")
	ColorPink     = lipgloss.Color("#f5bde6")
	ColorFlamingo = lipgloss.Color("#f0c6c6")
	ColorRosewater= lipgloss.Color("#f4dbd6")
)

// ─── Re-usable Styles ──────────────────────────────────────────────────────

var (
	// TitleStyle is used for primary headings.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorLavender).
			MarginBottom(1)

	// SubtitleStyle is for secondary descriptive text.
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext0).
			Italic(true)

	// StatusStyle highlights the current operation.
	StatusStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSapphire).
			PaddingLeft(2)

	// SuccessStyle for positive results.
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	// ErrorStyle for failures.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	// WarningStyle for cautions.
	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorYellow)

	// BoxStyle wraps content in a rounded border.
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSurface1).
			Padding(1, 2)

	// AccentBoxStyle wraps content in a colored border for emphasis.
	AccentBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorMauve).
			Padding(1, 3)

	// LogLineStyle dims live log output so it doesn't distract.
	LogLineStyle = lipgloss.NewStyle().
			Foreground(ColorOverlay0).
			MaxWidth(80)

	// HelpStyle for bottom-bar hints.
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSurface1).
			Italic(true).
			MarginTop(1)

	// ResultLabelStyle for the summary screen labels.
	ResultLabelStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext1).
			Width(12)

	// DurationStyle for timing display.
	DurationStyle = lipgloss.NewStyle().
			Foreground(ColorOverlay0)
)
