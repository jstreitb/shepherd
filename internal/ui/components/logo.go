// Package components provides reusable TUI sub-components.
package components

import "github.com/charmbracelet/lipgloss"

// SheepLogo returns a clean, minimal title for the application.
func SheepLogo(accent lipgloss.Color) string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(accent).
		Render("🐑 baa")
}
