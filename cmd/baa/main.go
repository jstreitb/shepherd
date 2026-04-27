// BAA — A universal, autonomous Linux package manager updater.
//
// This is the application entry point. It initialises the Bubbletea program
// in alt-screen mode and runs the TUI until the user exits.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jstreitb/baa/internal/config"
	"github.com/jstreitb/baa/internal/detector"
	"github.com/jstreitb/baa/internal/theme"
	"github.com/jstreitb/baa/internal/ui"
	"github.com/jstreitb/baa/internal/updater"
)

var version = "dev" // Overridden by GoReleaser using ldflags

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	updateMe := flag.Bool("update", false, "update baa to the latest version")
	uninstallMe := flag.Bool("uninstall", false, "uninstall baa from the system")
	showHelp := flag.Bool("help", false, "show help message and exit")
	showCredits := flag.Bool("credits", false, "show credits and exit")
	flag.BoolVar(showCredits, "c", *showCredits, "show credits and exit")

	flag.Usage = func() {
		fmt.Printf("BAA — A universal, autonomous Linux package manager updater.\n\n")
		fmt.Printf("Usage:\n  baa [flags]\n\n")
		fmt.Printf("Flags:\n")
		fmt.Printf("  --update     Update baa to the latest version\n")
		fmt.Printf("  --uninstall  Uninstall baa from the system\n")
		fmt.Printf("  --version    Print version and exit\n")
		fmt.Printf("  --credits    Show credits and exit\n")
		fmt.Printf("  --help       Show this help message and exit\n")
	}

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *showCredits {
		printCredits()
		os.Exit(0)
	}

	if *uninstallMe {
		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error determining executable path: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Uninstalling baa from %s...\n", exe)
		if err := os.Remove(exe); err != nil {
			if os.IsPermission(err) {
				fmt.Fprintf(os.Stderr, "Permission denied. Please run with sudo:\n  sudo baa --uninstall\n")
			} else {
				fmt.Fprintf(os.Stderr, "Error uninstalling: %v\n", err)
			}
			os.Exit(1)
		}
		fmt.Println("BAA has been uninstalled successfully.")
		os.Exit(0)
	}

	if *updateMe {
		// Use process substitution instead of a pipe to keep stdin available for sudo prompts
		cmd := exec.Command("bash", "-c", "bash <(curl -sSfL https://raw.githubusercontent.com/jstreitb/baa/main/install.sh)")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "\nOops! The update failed.\n(GitHub returned a 404 Error. Have you pushed the repository & install.sh to 'jstreitb/baa' yet?)\n")
			os.Exit(1)
		}
		os.Exit(0)
	}

	cfg := config.Config{
		Version: version,
		Repo:    "jstreitb/baa",
	}

	// Decide which update checker to inject based on build version.
	// Dev builds skip the network call entirely; test builds use a static response.
	var checker updater.VersionChecker
	switch version {
	case "dev", "":
		// No update checking in development mode.
	case "test-update":
		checker = &updater.StaticChecker{Version: "2.0.0-PRO-EDITION"}
	default:
		checker = updater.NewGitHubChecker(cfg.Repo)
	}

	det := detector.New(detector.DefaultRunner{})

	p := tea.NewProgram(
		ui.NewModel(cfg, checker, det),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "baa: %v\n", err)
		os.Exit(1)
	}
}

func printCredits() {
	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.ColorGreen)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(theme.ColorLavender).
		Bold(true)

	contentStyle := lipgloss.NewStyle().
		Foreground(theme.ColorSubtext0)

	// Build sections
	title := titleStyle.Render("🐑 baa — The update herd")

	devSection := lipgloss.JoinVertical(lipgloss.Left,
		subtitleStyle.Render("Main Developer:"),
		contentStyle.Render("jstreitb"),
	)

	poweredSection := lipgloss.JoinVertical(lipgloss.Left,
		subtitleStyle.Render("Powered by:"),
		contentStyle.Render("Charmbracelet (Bubble Tea, Lip Gloss, Bubbles)"),
	)

	licenseSection := lipgloss.JoinVertical(lipgloss.Left,
		subtitleStyle.Render("License:"),
		contentStyle.Render("MIT"),
	)

	thanksSection := contentStyle.Render("Special thanks to all contributors!")

	// Join all sections with consistent spacing
	out := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		devSection,
		"",
		poweredSection,
		"",
		licenseSection,
		"",
		thanksSection,
	)

	// Print with a single wrapper
	fmt.Println(lipgloss.NewStyle().Margin(1, 2).Render(out))
}
