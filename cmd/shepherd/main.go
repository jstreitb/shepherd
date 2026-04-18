// Shepherd — A universal, autonomous Linux package manager updater.
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
	"github.com/jstreitb/shepherd/internal/ui"
)

var version = "dev" // Overridden by GoReleaser using ldflags

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	updateMe := flag.Bool("update", false, "update shepherd to the latest version")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *updateMe {
		cmd := exec.Command("bash", "-c", "curl -sSfL https://raw.githubusercontent.com/jstreitb/shepherd/main/install.sh | bash")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "\nOops! The update failed.\n(GitHub returned a 404 Error. Have you pushed the repository & install.sh to 'jstreitb/shepherd' yet?)\n")
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
		fmt.Fprintf(os.Stderr, "shepherd: %v\n", err)
		os.Exit(1)
	}
}
