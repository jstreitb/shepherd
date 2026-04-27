// Package ui implements the Bubbletea TUI for BAA.
package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jstreitb/baa/internal/theme"
)

// ─── Re-usable Styles ──────────────────────────────────────────────────────

var (
	// TitleStyle is used for primary headings.
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.ColorLavender).
			MarginBottom(1)

	// SubtitleStyle is for secondary descriptive text.
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(theme.ColorSubtext0).
			Italic(true)

	// StatusStyle highlights the current operation.
	StatusStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.ColorSapphire).
			PaddingLeft(2)

	// SuccessStyle for positive results.
	SuccessStyle = lipgloss.NewStyle().
			Foreground(theme.ColorGreen).
			Bold(true)

	// ErrorStyle for failures.
	ErrorStyle = lipgloss.NewStyle().
			Foreground(theme.ColorRed).
			Bold(true)

	// WarningStyle for cautions.
	WarningStyle = lipgloss.NewStyle().
			Foreground(theme.ColorYellow)

	// BoxStyle wraps content in a rounded border.
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.ColorSurface1).
			Padding(1, 2)


	// LogLineStyle dims live log output so it doesn't distract.
	LogLineStyle = lipgloss.NewStyle().
			Foreground(theme.ColorOverlay0).
			MaxWidth(80)

	// HelpStyle for bottom-bar hints.
	HelpStyle = lipgloss.NewStyle().
			Foreground(theme.ColorSurface1).
			Italic(true).
			MarginTop(1)


	// DurationStyle for timing display.
	DurationStyle = lipgloss.NewStyle().
			Foreground(theme.ColorOverlay0)
)
