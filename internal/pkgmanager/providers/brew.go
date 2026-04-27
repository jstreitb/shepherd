package providers

// Brew implements PackageManager for Homebrew.
type Brew struct{}

func (b *Brew) Name() string    { return "brew" }
func (b *Brew) NeedsSudo() bool { return false }

func (b *Brew) Commands() [][]string {
	return [][]string{
		{"brew", "upgrade"},
	}
}

func (b *Brew) Env() []string {
	return []string{
		"HOMEBREW_NO_AUTO_UPDATE=1",
	}
}
