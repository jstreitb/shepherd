package pkgmanager

import (
	"os/exec"

	"github.com/jstreitb/baa/internal/pkgmanager/providers"
)

// DetectInstalled probes the system for known package managers using exec.LookPath.
// It returns a slice of PackageManager implementations for every manager found
// on the current PATH, preserving a deterministic update order:
// apt → dnf → zypper → pacman → nix → brew → flatpak → snap.
func DetectInstalled() []PackageManager {
	type entry struct {
		bin     string
		factory func() PackageManager
	}

	candidates := []entry{
		{"apt-get", func() PackageManager { return &providers.Apt{} }},
		{"dnf", func() PackageManager { return &providers.Dnf{} }},
		{"zypper", func() PackageManager { return &providers.Zypper{} }},
		{"pacman", func() PackageManager { return &providers.Pacman{} }},
		{"nix-env", func() PackageManager { return &providers.Nix{} }},
		{"brew", func() PackageManager { return &providers.Brew{} }},
		{"flatpak", func() PackageManager { return &providers.Flatpak{} }},
		{"snap", func() PackageManager { return &providers.Snap{} }},
	}

	var found []PackageManager
	for _, c := range candidates {
		if _, err := exec.LookPath(c.bin); err == nil {
			found = append(found, c.factory())
		}
	}
	return found
}
