package executor

import (
	"os/exec"
	"strings"

	"github.com/jstreitb/baa/internal/pkgmanager"
)

// BuildInteractiveCmd constructs an *exec.Cmd suitable for tea.ExecProcess
// that runs all of a manager's commands sequentially with full terminal I/O.
// This moves the shell-command construction out of the UI layer.
func BuildInteractiveCmd(mgr pkgmanager.PackageManager) *exec.Cmd {
	var parts []string
	for _, cmd := range mgr.Commands() {
		parts = append(parts, strings.Join(cmd, " "))
	}
	script := strings.Join(parts, " && ")

	if mgr.NeedsSudo() {
		return exec.Command("sudo", "bash", "-c", script)
	}
	return exec.Command("bash", "-c", script)
}
