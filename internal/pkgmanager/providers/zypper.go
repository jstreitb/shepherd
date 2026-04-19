package providers

// Zypper implements PackageManager for openSUSE's zypper.
type Zypper struct{}

func (z *Zypper) Name() string     { return "zypper" }
func (z *Zypper) NeedsSudo() bool  { return true }

func (z *Zypper) Commands() [][]string {
	return [][]string{
		{"zypper", "--non-interactive", "update"},
	}
}

func (z *Zypper) Env() []string { return nil }
