package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/jstreitb/baa/internal/theme"
)

// ─── View: Init ─────────────────────────────────────────────────────────────

func viewInit(d InitViewData) string {
	return centerVertically(d.Width, d.Height,
		lipgloss.JoinVertical(lipgloss.Center,
			TitleStyle.Render("🐑 baa"),
			"",
			d.SpinnerView+"  "+SubtitleStyle.Render("Detecting package managers…"),
		),
	)
}

// ─── View: Login ────────────────────────────────────────────────────────────
// Minimalist: title, description, detected list, password field, "Next".

func viewLogin(d LoginViewData) string {
	title := TitleStyle.Render("🐑 baa")

	desc := SubtitleStyle.Render("Enter your password to update your system.")

	detected := lipgloss.NewStyle().
		Foreground(theme.ColorOverlay0).
		Render("Detected:  " + strings.Join(d.ManagerNames, "  ·  "))

	next := HelpStyle.Render("Next →   (Press Enter)")

	var updateHint string
	if d.LatestVersion != "" {
		updateHint = WarningStyle.Render(fmt.Sprintf("Update available: v%s! Run `baa --update` later", d.LatestVersion))
	}

	return centerVertically(d.Width, d.Height,
		lipgloss.JoinVertical(lipgloss.Center,
			title,
			"",
			desc,
			detected,
			"",
			"",
			d.TextInputView,
			"",
			"",
			next,
			"",
			updateHint,
		),
	)
}

// ─── View: Updating ─────────────────────────────────────────────────────────

func viewUpdating(d UpdatingViewData) string {
	header := StatusStyle.Render(
		fmt.Sprintf("%s  Updating %s…", d.SpinnerView, d.ManagerName),
	)
	progress := SubtitleStyle.Render(
		fmt.Sprintf("%d / %d", d.CurrentIndex+1, d.TotalManagers),
	)

	animFrame := BoxStyle.
		Foreground(theme.ColorMauve).
		Render(d.AnimationFrame)

	logLine := d.LastLogLine
	if len(logLine) > 72 {
		logLine = logLine[:69] + "..."
	}
	log := LogLineStyle.Render("  " + logLine)

	// Previously completed managers.
	var done []string
	for _, r := range d.PastResults {
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

	return centerVertically(d.Width, d.Height, lipgloss.JoinVertical(lipgloss.Center, parts...))
}

// ─── View: Failed ───────────────────────────────────────────────────────────

func viewFailed(d FailedViewData) string {
	header := ErrorStyle.Render(fmt.Sprintf("✗  %s failed", d.ManagerName))

	errMsg := ""
	if d.ErrorMsg != "" {
		errMsg = lipgloss.NewStyle().
			Foreground(theme.ColorPeach).
			Render(d.ErrorMsg)
	}

	options := lipgloss.JoinVertical(lipgloss.Left,
		"",
		lipgloss.NewStyle().Foreground(theme.ColorText).Render("  r   Retry interactively"),
		lipgloss.NewStyle().Foreground(theme.ColorSubtext0).Render("  s   Skip"),
		lipgloss.NewStyle().Foreground(theme.ColorOverlay0).Render("  q   Quit"),
	)

	return centerVertically(d.Width, d.Height,
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

func viewSummary(d SummaryViewData) string {
	if !d.HasManagers {
		return centerVertically(d.Width, d.Height,
			lipgloss.JoinVertical(lipgloss.Center,
				TitleStyle.Render("🐑 baa"),
				"",
				WarningStyle.Render("No supported package managers found."),
				SubtitleStyle.Render("Supported: apt, brew, dnf, flatpak, nix, pacman, snap, zypper"),
				"",
				HelpStyle.Render("Press q to exit"),
			),
		)
	}

	divider := lipgloss.NewStyle().Foreground(theme.ColorSurface1).
		Render(strings.Repeat("─", 36))

	allOK := true
	var rows []string
	for _, r := range d.Results {
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
					Foreground(theme.ColorPeach).
					PaddingLeft(5).
					Render("↳ "+LastErrorLine(r.Error)),
			)
		}
	}

	status := SuccessStyle.Render("All updates complete 🎉")
	if !allOK {
		status = WarningStyle.Render("Some updates had issues.")
	}

	var updateHint string
	if d.LatestVersion != "" {
		updateHint = WarningStyle.Render(fmt.Sprintf("Update available: v%s! Run `baa --update`", d.LatestVersion))
	}

	return centerVertically(d.Width, d.Height,
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

func centerVertically(width, height int, content string) string {
	h := lipgloss.Height(content)
	pad := 0
	if height > h {
		pad = (height - h) / 2
	}
	return lipgloss.NewStyle().
		Width(width).
		PaddingTop(pad).
		Align(lipgloss.Center).
		Render(content)
}

// LastErrorLine returns the last non-empty line of an error string,
// which is usually the most relevant part.
func LastErrorLine(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if l := strings.TrimSpace(lines[i]); l != "" {
			return l
		}
	}
	return strings.TrimSpace(s)
}
