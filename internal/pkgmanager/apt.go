package pkgmanager

// Apt implements PackageManager for Debian/Ubuntu's apt-get.
// It uses DEBIAN_FRONTEND=noninteractive and dpkg --force-conf* flags
// to suppress all interactive prompts during unattended upgrades.
type Apt struct{}

func (a *Apt) Name() string     { return "apt" }
func (a *Apt) NeedsSudo() bool   { return true }

func (a *Apt) Commands() [][]string {
	return [][]string{
		{"apt-get", "update", "-qq"},
		{
			"apt-get", "dist-upgrade", "-y",
			"-o", "Dpkg::Options::=--force-confold",
			"-o", "Dpkg::Options::=--force-confdef",
		},
	}
}

func (a *Apt) Env() []string {
	return []string{
		"DEBIAN_FRONTEND=noninteractive",
		"NEEDRESTART_MODE=a",
	}
}
