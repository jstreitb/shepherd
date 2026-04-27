package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jstreitb/baa/internal/executor"
	"github.com/jstreitb/baa/internal/pkgmanager"
	"github.com/jstreitb/baa/internal/ui/components"
)

// Orchestrator owns all state related to running package manager updates:
// the manager list, current index, channels, password, and results.
// It is stored as a pointer in the Model so that Bubble Tea's by-value
// copies share a single instance.
type Orchestrator struct {
	password   *executor.SecurePassword
	managers   []pkgmanager.PackageManager
	currentMgr int
	lastLine   string
	outputCh   chan string
	resultCh   chan pkgmanager.UpdateResult
	results    []pkgmanager.UpdateResult
}

// NewOrchestrator creates an empty Orchestrator.
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{}
}

// SetManagers stores the detected package managers.
func (o *Orchestrator) SetManagers(mgrs []pkgmanager.PackageManager) {
	o.managers = mgrs
}

// Managers returns the list of detected package managers.
func (o *Orchestrator) Managers() []pkgmanager.PackageManager {
	return o.managers
}

// SetPassword captures the user's password in a SecurePassword wrapper.
// The source slice is zeroed by NewSecurePassword.
func (o *Orchestrator) SetPassword(raw []byte) {
	o.password = executor.NewSecurePassword(raw)
}

// CurrentManager returns the manager currently being updated.
func (o *Orchestrator) CurrentManager() pkgmanager.PackageManager {
	if o.currentMgr < len(o.managers) {
		return o.managers[o.currentMgr]
	}
	return nil
}

// CurrentIndex returns the zero-based index of the current manager.
func (o *Orchestrator) CurrentIndex() int {
	return o.currentMgr
}

// LastLine returns the most recent output line from a running update.
func (o *Orchestrator) LastLine() string {
	return o.lastLine
}

// SetLastLine stores the latest output line.
func (o *Orchestrator) SetLastLine(line string) {
	o.lastLine = line
}

// Results returns the accumulated update results.
func (o *Orchestrator) Results() []pkgmanager.UpdateResult {
	return o.results
}

// AddResult appends a result. Returns true if the update succeeded.
func (o *Orchestrator) AddResult(r pkgmanager.UpdateResult) bool {
	o.results = append(o.results, r)
	return r.Success
}

// OverwriteLastResult updates the most recent result (used after interactive retry).
func (o *Orchestrator) OverwriteLastResult(success bool, errMsg string) {
	if len(o.results) == 0 {
		return
	}
	last := &o.results[len(o.results)-1]
	last.Success = success
	last.Error = errMsg
}

// OutputCh returns the current output channel for waitForOutput.
func (o *Orchestrator) OutputCh() <-chan string {
	return o.outputCh
}

// StartUpdate launches the current manager's update in a goroutine
// and returns Cmds to listen for output and the final result.
func (o *Orchestrator) StartUpdate() tea.Cmd {
	mgr := o.managers[o.currentMgr]
	o.outputCh = make(chan string, 128)
	o.resultCh = make(chan pkgmanager.UpdateResult, 1)
	o.lastLine = ""

	pw := o.password.Copy()

	go func() {
		result := executor.RunManagerUpdate(pw, mgr, o.outputCh)
		executor.ZeroBytes(pw)
		o.resultCh <- result
	}()

	return tea.Batch(
		waitForOutput(o.outputCh),
		waitForResult(o.resultCh),
	)
}

// Advance moves to the next manager or signals completion.
// Returns the new State and any tea.Cmd to execute.
func (o *Orchestrator) Advance() (State, tea.Cmd) {
	o.currentMgr++
	if o.currentMgr < len(o.managers) {
		return StateUpdating, tea.Batch(
			o.StartUpdate(),
			components.AnimTickCmd(200*time.Millisecond),
		)
	}
	o.password.Close()
	return StateSummary, nil
}

// InteractiveRetry constructs a tea.Cmd that suspends the TUI and
// re-runs the failed manager's commands with full terminal I/O.
func (o *Orchestrator) InteractiveRetry() tea.Cmd {
	mgr := o.managers[o.currentMgr]
	c := executor.BuildInteractiveCmd(mgr)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return interactiveDoneMsg{err: err}
	})
}

// BeginUpdating resets the orchestrator for the update phase.
func (o *Orchestrator) BeginUpdating() {
	o.currentMgr = 0
}

// Cleanup zeroes the password. Safe to call multiple times.
func (o *Orchestrator) Cleanup() {
	o.password.Close()
}
