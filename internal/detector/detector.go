package detector

import (
	"os/exec"

	"github.com/jstreitb/baa/internal/pkgmanager"
	"github.com/jstreitb/baa/internal/pkgmanager/providers"
)

// CommandRunner provides a way to look up executables in the system PATH.
type CommandRunner interface {
	LookPath(file string) (string, error)
}

// DefaultRunner is the default implementation using os/exec.
type DefaultRunner struct{}

func (d DefaultRunner) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

// Detector probes the system for known package managers.
type Detector struct {
	runner CommandRunner
}

// New creates a new Detector with the given CommandRunner.
func New(runner CommandRunner) *Detector {
	if runner == nil {
		runner = DefaultRunner{}
	}
	return &Detector{runner: runner}
}

// DetectInstalled probes the system for known package managers.
func (d *Detector) DetectInstalled() []pkgmanager.PackageManager {
	type entry struct {
		bin     string
		factory func() pkgmanager.PackageManager
	}

	candidates := []entry{
		{"apt-get", func() pkgmanager.PackageManager { return &providers.Apt{} }},
		{"dnf", func() pkgmanager.PackageManager { return &providers.Dnf{} }},
		{"zypper", func() pkgmanager.PackageManager { return &providers.Zypper{} }},
		{"pacman", func() pkgmanager.PackageManager { return &providers.Pacman{} }},
		{"nix-env", func() pkgmanager.PackageManager { return &providers.Nix{} }},
		{"brew", func() pkgmanager.PackageManager { return &providers.Brew{} }},
		{"flatpak", func() pkgmanager.PackageManager { return &providers.Flatpak{} }},
		{"snap", func() pkgmanager.PackageManager { return &providers.Snap{} }},
	}

	var found []pkgmanager.PackageManager
	for _, c := range candidates {
		if _, err := d.runner.LookPath(c.bin); err == nil {
			found = append(found, c.factory())
		}
	}
	return found
}
