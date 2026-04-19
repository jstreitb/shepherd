package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ─── View: Init ─────────────────────────────────────────────────────────────

func viewInit(m Model) string {
	return centerVertically(m,
		lipgloss.JoinVertical(lipgloss.Center,
			TitleStyle.Render("🐑 baa"),
			"",
			m.spinner.View()+"  "+SubtitleStyle.Render("Detecting package managers…"),
		),
	)
}

// ─── View: Login ────────────────────────────────────────────────────────────
// Minimalist: title, description, detected list, password field, "Next".

func viewLogin(m Model) string {
	title := TitleStyle.Render("🐑 baa")

	desc := SubtitleStyle.Render("Enter your password to update your system.")

	var names []string
	for i, mgr := range m.managers {
		if i == 3 {
			names = append(names, "...")
			break
		}
		names = append(names, mgr.Name())
	}
	detected := lipgloss.NewStyle().
		Foreground(ColorOverlay0).
		Render("Detected:  " + strings.Join(names, "  ·  "))

	field := m.textInput.View()

	next := HelpStyle.Render("Next →   (Press Enter)")

	var updateHint string
	if m.latestVersion != "" {
		updateHint = WarningStyle.Render(fmt.Sprintf("Update available: v%s! Run `baa --update` later", m.latestVersion))
	}

	return centerVertically(m,
		lipgloss.JoinVertical(lipgloss.Center,
			title,
			"",
			desc,
			detected,
			"",
			"",
			field,
			"",
			"",
			next,
			"",
			updateHint,
		),
	)
}

// ─── View: Updating ─────────────────────────────────────────────────────────

func viewUpdating(m Model) string {
	mgr := m.managers[m.currentMgr]

	header := StatusStyle.Render(
		fmt.Sprintf("%s  Updating %s…", m.spinner.View(), mgr.Name()),
	)
	progress := SubtitleStyle.Render(
		fmt.Sprintf("%d / %d", m.currentMgr+1, len(m.managers)),
	)

	animFrame := BoxStyle.
		Foreground(ColorMauve).
		Render(m.animation.Frame())

	logLine := m.lastLine
	if len(logLine) > 72 {
		logLine = logLine[:69] + "..."
	}
	log := LogLineStyle.Render("  " + logLine)

	// Previously completed managers.
	var done []string
	for _, r := range m.results {
		icon := SuccessStyle.Render("✓")
		if !r.Success {
			icon = ErrorStyle.Render("✗")
		}
		done = append(done, fmt.Sprintf("  %s  %-12s  %s",
			icon, r.Manager,
			DurationStyle.Render(r.Duration.Round(time.Millisecond).String()),
		))
	}

	parts := []string{
		lipgloss.JoinHorizontal(lipgloss.Top, header, "    ", progress),
		"",
		animFrame,
		"",
		log,
	}
	if len(done) > 0 {
		parts = append(parts, "", strings.Join(done, "\n"))
	}

	return centerVertically(m, lipgloss.JoinVertical(lipgloss.Center, parts...))
}

// ─── View: Failed ───────────────────────────────────────────────────────────

func viewFailed(m Model) string {
	last := m.results[len(m.results)-1]

	header := ErrorStyle.Render(fmt.Sprintf("✗  %s failed", last.Manager))

	errMsg := ""
	if last.Error != "" {
		errMsg = lipgloss.NewStyle().
			Foreground(ColorPeach).
			Render(last.Error)
	}

	options := lipgloss.JoinVertical(lipgloss.Left,
		"",
		lipgloss.NewStyle().Foreground(ColorText).Render("  r   Retry interactively"),
		lipgloss.NewStyle().Foreground(ColorSubtext0).Render("  s   Skip"),
		lipgloss.NewStyle().Foreground(ColorOverlay0).Render("  q   Quit"),
	)

	return centerVertically(m,
		lipgloss.JoinVertical(lipgloss.Center,
			header,
			"",
			errMsg,
			"",
			options,
		),
	)
}

// ─── View: Summary ──────────────────────────────────────────────────────────

func viewSummary(m Model) string {
	if len(m.managers) == 0 {
		return centerVertically(m,
			lipgloss.JoinVertical(lipgloss.Center,
				TitleStyle.Render("🐑 baa"),
				"",
				WarningStyle.Render("No supported package managers found."),
				SubtitleStyle.Render("Supported: apt, pacman, flatpak, snap"),
				"",
				HelpStyle.Render("Press q to exit"),
			),
		)
	}

	divider := lipgloss.NewStyle().Foreground(ColorSurface1).
		Render(strings.Repeat("─", 36))

	allOK := true
	var rows []string
	for _, r := range m.results {
		icon := SuccessStyle.Render("✓")
		if !r.Success {
			icon = ErrorStyle.Render("✗")
			allOK = false
		}
		rows = append(rows, fmt.Sprintf("  %s  %-12s  %s",
			icon, r.Manager,
			DurationStyle.Render(r.Duration.Round(time.Millisecond).String()),
		))
		if !r.Success && r.Error != "" {
			rows = append(rows,
				lipgloss.NewStyle().
					Foreground(ColorPeach).
					PaddingLeft(5).
					Render("↳ "+last_error_line(r.Error)),
			)
		}
	}

	status := SuccessStyle.Render("All updates complete 🎉")
	if !allOK {
		status = WarningStyle.Render("Some updates had issues.")
	}
	
	var updateHint string
	if m.latestVersion != "" {
		updateHint = WarningStyle.Render(fmt.Sprintf("Update available: v%s! Run `baa --update`", m.latestVersion))
	}

	return centerVertically(m,
		lipgloss.JoinVertical(lipgloss.Center,
			TitleStyle.Render("🐑 Update Complete"),
			"",
			divider,
			"",
			strings.Join(rows, "\n"),
			"",
			divider,
			"",
			status,
			"",
			updateHint,
			"",
			HelpStyle.Render("Press q to exit"),
		),
	)
}

// ─── Helpers ────────────────────────────────────────────────────────────────

func centerVertically(m Model, content string) string {
	h := lipgloss.Height(content)
	pad := 0
	if m.height > h {
		pad = (m.height - h) / 2
	}
	return lipgloss.NewStyle().
		Width(m.width).
		PaddingTop(pad).
		Align(lipgloss.Center).
		Render(content)
}

// last_error_line returns the last non-empty line of an error string,
// which is usually the most relevant part.
func last_error_line(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if l := strings.TrimSpace(lines[i]); l != "" {
			return l
		}
	}
	return s
}
