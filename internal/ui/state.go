package ui

// State represents a screen in the TUI state machine.
type State int

const (
	StateInit     State = iota // Silently detect package managers.
	StateLogin                 // Prompt for sudo password.
	StateUpdating              // Run updates sequentially.
	StateFailed                // A manager failed; offer retry.
	StateSummary               // Show final results.
)

// Legal state transitions:
//
//   Init ──detectDoneMsg(>0)──► Login
//   Init ──detectDoneMsg(=0)──► Summary
//   Login ──Submit──► Updating
//   Updating ──Success + more──► Updating
//   Updating ──Success + done──► Summary
//   Updating ──Failure──► Failed
//   Failed ──Retry/Skip──► Updating

