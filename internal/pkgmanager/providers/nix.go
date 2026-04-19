package providers

// Nix implements PackageManager for the Nix package manager.
type Nix struct{}

func (n *Nix) Name() string     { return "nix" }
func (n *Nix) NeedsSudo() bool  { return false }

func (n *Nix) Commands() [][]string {
	return [][]string{
		{"nix-env", "-u", "*"},
	}
}

func (n *Nix) Env() []string { return nil }
