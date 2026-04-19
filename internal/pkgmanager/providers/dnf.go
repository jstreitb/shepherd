package providers

// Dnf implements PackageManager for Fedora/RHEL's dnf.
type Dnf struct{}

func (d *Dnf) Name() string     { return "dnf" }
func (d *Dnf) NeedsSudo() bool  { return true }

func (d *Dnf) Commands() [][]string {
	return [][]string{
		{"dnf", "upgrade", "-y"},
	}
}

func (d *Dnf) Env() []string { return nil }
