package pkgmanager

// Snap implements PackageManager for Canonical's Snap.
type Snap struct{}

func (s *Snap) Name() string     { return "snap" }
func (s *Snap) NeedsSudo() bool   { return true }

func (s *Snap) Commands() [][]string {
	return [][]string{
		{"snap", "refresh"},
	}
}

func (s *Snap) Env() []string { return nil }
