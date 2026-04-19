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

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jstreitb/baa/internal/ui"
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
		cmd := exec.Command("bash", "-c", "curl -sSfL https://raw.githubusercontent.com/jstreitb/baa/main/install.sh | bash")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "\nOops! The update failed.\n(GitHub returned a 404 Error. Have you pushed the repository & install.sh to 'jstreitb/baa' yet?)\n")
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Make version available to the TUI to check for updates
	ui.AppVersion = version

	p := tea.NewProgram(
		ui.NewModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "baa: %v\n", err)
		os.Exit(1)
	}
}

func printCredits() {
	// Color palette from Catppuccin Macchiato
	colorGreen := lipgloss.Color("#a6da95")
	colorLavender := lipgloss.Color("#b7bdf8")
	colorSubtext0 := lipgloss.Color("#a5adcb")

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorGreen)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(colorLavender).
		Bold(true)

	contentStyle := lipgloss.NewStyle().
		Foreground(colorSubtext0)

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
