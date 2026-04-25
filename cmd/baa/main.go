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
	"github.com/jstreitb/baa/internal/pkgmanager"
	"github.com/jstreitb/baa/internal/ui"
	"github.com/jstreitb/baa/internal/utils"
)

var version = "dev" // Overridden by GoReleaser using ldflags

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	updateMe := flag.Bool("update", false, "update baa to the latest version")
	uninstallMe := flag.Bool("uninstall", false, "uninstall baa from the system")
	showHelp := flag.Bool("help", false, "show help message and exit")
	showCredits := flag.Bool("credits", false, "show credits and exit")
	showDetected := flag.Bool("detect", false, "show detected package managers and exit")
	quiet := flag.Bool("quiet", false, "skip the TUI; run updates non-interactively and print a one-line summary per manager")
	flag.BoolVar(showCredits, "c", *showCredits, "show credits and exit")
	flag.BoolVar(showDetected, "d", *showDetected, "show detected package managers and exit")
	flag.BoolVar(quiet, "q", *quiet, "skip the TUI; run updates non-interactively and print a one-line summary per manager")

	flag.Usage = func() {
		fmt.Printf("BAA — A universal, autonomous Linux package manager updater.\n\n")
		fmt.Printf("Usage:\n  baa [flags]\n\n")
		fmt.Printf("Flags:\n")
		fmt.Printf("  --update     Update baa to the latest version\n")
		fmt.Printf("  --uninstall  Uninstall baa from the system\n")
		fmt.Printf("  --version    Print version and exit\n")
		fmt.Printf("  --credits    Show credits and exit\n")
		fmt.Printf("  --detect     Show detected package managers and exit\n")
		fmt.Printf("  --quiet      Skip the TUI; run updates non-interactively (set BAA_SUDO_PASSWORD for sudo managers)\n")
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

	if *showDetected {
		found := pkgmanager.DetectInstalled()
		if len(found) == 0 {
			fmt.Println("baa has detected no package managers.")
			os.Exit(0)
		}
		fmt.Println("baa has detected the following managers:")
		for _, manager := range found {
			fmt.Println("-", manager.Name())
		}
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

	if *quiet {
		os.Exit(runQuiet())
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

// runQuiet executes every detected package manager's update sequence
// without the TUI, printing one summary line per manager. Sudo-needing
// managers read their password from BAA_SUDO_PASSWORD; if that is
// unset and at least one manager needs sudo, runQuiet exits non-zero
// before running anything (so a script does not partially update).
//
// Returns the exit code: 0 if every manager succeeded, 1 otherwise.
func runQuiet() int {
	managers := pkgmanager.DetectInstalled()
	if len(managers) == 0 {
		fmt.Fprintln(os.Stderr, "baa: no package managers detected")
		return 1
	}

	var password []byte
	needsSudo := false
	for _, mgr := range managers {
		if mgr.NeedsSudo() {
			needsSudo = true
			break
		}
	}
	if needsSudo {
		pw := os.Getenv("BAA_SUDO_PASSWORD")
		if pw == "" {
			fmt.Fprintln(os.Stderr,
				"baa: --quiet requires BAA_SUDO_PASSWORD when sudo-needing managers are present")
			return 1
		}
		password = []byte(pw)
		defer utils.ZeroBytes(password)
	}

	failures := 0
	for _, mgr := range managers {
		fmt.Printf("▶ %s... ", mgr.Name())
		// outputCh is drained but its contents are dropped to keep
		// quiet mode actually quiet; only the final summary is shown.
		outputCh := make(chan string, 64)
		go func() {
			for range outputCh {
			}
		}()
		result := utils.RunManagerUpdate(password, mgr, outputCh)
		if result.Success {
			fmt.Printf("ok (%s)\n", result.Duration.Round(1e9))
		} else {
			failures++
			fmt.Printf("failed (%s)\n", result.Duration.Round(1e9))
			if result.Error != "" {
				fmt.Fprintf(os.Stderr, "  %s\n", result.Error)
			}
		}
	}

	if failures > 0 {
		return 1
	}
	return 0
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
