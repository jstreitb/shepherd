// Package pkgmanager defines the interface and types for system package managers.
package pkgmanager

import "time"

// PackageManager represents a system package manager that can perform updates.
type PackageManager interface {
	// Name returns the human-readable name of the package manager.
	Name() string

	// NeedsSudo returns true if this manager requires root privileges.
	NeedsSudo() bool

	// Commands returns the ordered list of commands to execute for a full update.
	// Each inner slice is a single command with its arguments.
	Commands() [][]string

	// Env returns additional environment variables for non-interactive execution.
	Env() []string
}

// UpdateResult holds the outcome of running a package manager update.
type UpdateResult struct {
	Manager  string        // Name of the package manager.
	Success  bool          // Whether all commands completed successfully.
	Output   string        // Combined stdout/stderr output.
	Error    string        // Sanitized error snippet on failure.
	Duration time.Duration // Wall-clock time for the entire update.
}
