package tests

import (
	"reflect"
	"testing"

	"github.com/jstreitb/baa/internal/pkgmanager/providers"
)

type providerTest struct {
	name     string
	provider interface {
		Name() string
		NeedsSudo() bool
		Commands() [][]string
		Env() []string
	}
	wantName string
	wantSudo bool
	wantCmds [][]string
	wantEnv  []string
}

func TestProviders(t *testing.T) {
	tests := []providerTest{
		{
			name:     "Apt",
			provider: &providers.Apt{},
			wantName: "apt",
			wantSudo: true,
			wantCmds: [][]string{
				{"apt-get", "update", "-qq"},
				{
					"apt-get", "dist-upgrade", "-y",
					"-o", "Dpkg::Options::=--force-confold",
					"-o", "Dpkg::Options::=--force-confdef",
				},
			},
			wantEnv: []string{
				"DEBIAN_FRONTEND=noninteractive",
				"NEEDRESTART_MODE=a",
			},
		},
		{
			name:     "Brew",
			provider: &providers.Brew{},
			wantName: "brew",
			wantSudo: false,
			wantCmds: [][]string{
				{"brew", "upgrade"},
			},
			wantEnv: []string{
				"HOMEBREW_NO_AUTO_UPDATE=1",
			},
		},
		{
			name:     "Dnf",
			provider: &providers.Dnf{},
			wantName: "dnf",
			wantSudo: true,
			wantCmds: [][]string{
				{"dnf", "upgrade", "-y"},
			},
			wantEnv: nil,
		},
		{
			name:     "Flatpak",
			provider: &providers.Flatpak{},
			wantName: "flatpak",
			wantSudo: false,
			wantCmds: [][]string{
				{"flatpak", "update", "-y", "--noninteractive"},
			},
			wantEnv: nil,
		},
		{
			name:     "Nix",
			provider: &providers.Nix{},
			wantName: "nix",
			wantSudo: false,
			wantCmds: [][]string{
				{"nix-env", "-u", "*"},
			},
			wantEnv: nil,
		},
		{
			name:     "Pacman",
			provider: &providers.Pacman{},
			wantName: "pacman",
			wantSudo: true,
			wantCmds: [][]string{
				{"pacman", "-Syu", "--noconfirm"},
			},
			wantEnv: nil,
		},
		{
			name:     "Snap",
			provider: &providers.Snap{},
			wantName: "snap",
			wantSudo: true,
			wantCmds: [][]string{
				{"snap", "refresh"},
			},
			wantEnv: nil,
		},
		{
			name:     "Zypper",
			provider: &providers.Zypper{},
			wantName: "zypper",
			wantSudo: true,
			wantCmds: [][]string{
				{"zypper", "--non-interactive", "update"},
			},
			wantEnv: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.provider.Name(); got != tt.wantName {
				t.Errorf("Name() = %v, want %v", got, tt.wantName)
			}
			if got := tt.provider.NeedsSudo(); got != tt.wantSudo {
				t.Errorf("NeedsSudo() = %v, want %v", got, tt.wantSudo)
			}
			if got := tt.provider.Commands(); !reflect.DeepEqual(got, tt.wantCmds) {
				t.Errorf("Commands() = %v, want %v", got, tt.wantCmds)
			}
			if got := tt.provider.Env(); !reflect.DeepEqual(got, tt.wantEnv) {
				t.Errorf("Env() = %v, want %v", got, tt.wantEnv)
			}
		})
	}
}
