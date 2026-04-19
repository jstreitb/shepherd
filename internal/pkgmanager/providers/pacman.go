package providers

// Pacman implements PackageManager for Arch Linux's pacman.
type Pacman struct{}

func (p *Pacman) Name() string     { return "pacman" }
func (p *Pacman) NeedsSudo() bool   { return true }

func (p *Pacman) Commands() [][]string {
	return [][]string{
		{"pacman", "-Syu", "--noconfirm"},
	}
}

func (p *Pacman) Env() []string { return nil }
