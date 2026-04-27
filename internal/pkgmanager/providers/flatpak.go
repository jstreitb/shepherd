package providers

// Flatpak implements PackageManager for Flatpak.
// Flatpak does not require root privileges for user-scope updates.
type Flatpak struct{}

func (f *Flatpak) Name() string    { return "flatpak" }
func (f *Flatpak) NeedsSudo() bool { return false }

func (f *Flatpak) Commands() [][]string {
	return [][]string{
		{"flatpak", "update", "-y", "--noninteractive"},
	}
}

func (f *Flatpak) Env() []string { return nil }
