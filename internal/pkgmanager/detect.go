package pkgmanager

import "os/exec"

// DetectInstalled probes the system for known package managers using exec.LookPath.
// It returns a slice of PackageManager implementations for every manager found
// on the current PATH, preserving a deterministic update order:
// apt → pacman → flatpak → snap.
func DetectInstalled() []PackageManager {
	type entry struct {
		bin     string
		factory func() PackageManager
	}

	candidates := []entry{
		{"apt-get", func() PackageManager { return &Apt{} }},
		{"pacman", func() PackageManager { return &Pacman{} }},
		{"flatpak", func() PackageManager { return &Flatpak{} }},
		{"snap", func() PackageManager { return &Snap{} }},
	}

	var found []PackageManager
	for _, c := range candidates {
		if _, err := exec.LookPath(c.bin); err == nil {
			found = append(found, c.factory())
		}
	}
	return found
}
