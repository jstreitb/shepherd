package ui

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jstreitb/baa/internal/pkgmanager"
	"github.com/jstreitb/baa/internal/ui/components"
	"github.com/jstreitb/baa/internal/utils"
)

// ─── Application States ────────────────────────────────────────────────────

// State represents a screen in the TUI state machine.
type State int

const (
	StateInit     State = iota // Silently detect package managers.
	StateLogin                 // Prompt for sudo password.
	StateUpdating              // Run updates sequentially.
	StateFailed                // A manager failed; offer retry.
	StateSummary               // Show final results.
)

// ─── Internal Messages ─────────────────────────────────────────────────────

// detectDoneMsg carries the list of detected PackageManagers.
type detectDoneMsg struct{ managers []pkgmanager.PackageManager }

// outputLineMsg delivers a single live output line from a running command.
type outputLineMsg string

// outputDoneMsg signals that the output channel has been closed.
type outputDoneMsg struct{}

// managerDoneMsg carries the result of a finished manager update.
type managerDoneMsg struct{ result pkgmanager.UpdateResult }

// interactiveDoneMsg is sent when an interactive retry completes.
type interactiveDoneMsg struct{ err error }

// checkUpdateMsg carries the latest version from GitHub, or empty if check failed/none.
type checkUpdateMsg string

// AppVersion is injected by main.go from ldflags.
var AppVersion = "dev"

// ─── Model ─────────────────────────────────────────────────────────────────

// Model is the top-level Bubbletea model for BAA.
type Model struct {
	state  State
	width  int
	height int

	// Detection results.
	managers []pkgmanager.PackageManager

	// Update info
	latestVersion string

	// Login screen.
	textInput textinput.Model
	password  []byte

	// Update screen.
	currentMgr int
	lastLine   string
	outputCh   chan string
	resultCh   chan pkgmanager.UpdateResult
	animation  components.Animation
	spinner    spinner.Model

	// Results.
	results  []pkgmanager.UpdateResult
	quitting bool
}

// NewModel constructs the initial Model ready for tea.NewProgram.
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "sudo password"
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '✱'
	ti.CharLimit = 128
	ti.Width = 40
	ti.Focus()

	return Model{
		state:     StateInit,
		textInput: ti,
		animation: components.NewAnimation(),
		spinner:   components.NewSpinner(ColorMauve),
	}
}

// Init starts detection and the text input cursor blink.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		detectManagersCmd,
		checkUpdateCmd(),
		m.spinner.Tick,
	)
}

func checkUpdateCmd() tea.Cmd {
	return func() tea.Msg {
		if AppVersion == "test-update" {
			// Fake an update for testing purposes
			return checkUpdateMsg("2.0.0-PRO-EDITION")
		}
		if AppVersion == "dev" {
			return checkUpdateMsg("")
		}

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // Don't follow redirect, just read Location
			},
			Timeout: 2 * time.Second,
		}
		res, err := client.Get("https://github.com/jstreitb/baa/releases/latest")
		if err != nil {
			return checkUpdateMsg("")
		}
		defer res.Body.Close()

		loc, err := res.Location()
		if err == nil && loc != nil && strings.Contains(loc.Path, "/tag/") {
			parts := strings.Split(loc.Path, "/")
			latest := parts[len(parts)-1]
			latest = strings.TrimPrefix(latest, "v")
			curr := strings.TrimPrefix(AppVersion, "v")

			if latest != "" && latest != curr {
				return checkUpdateMsg(latest)
			}
		}
		return checkUpdateMsg("")
	}
}

// detectManagersCmd runs package manager detection in a background goroutine.
func detectManagersCmd() tea.Msg {
	return detectDoneMsg{managers: pkgmanager.DetectInstalled()}
}

// ─── Update (top-level dispatcher) ─────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, keys.Quit) && m.state != StateLogin && m.state != StateFailed {
			m.quitting = true
			// Zero password before exit.
			utils.ZeroBytes(m.password)
			return m, tea.Quit
		}
	}

	switch m.state {
	case StateInit:
		return m.updateInit(msg)
	case StateLogin:
		return m.updateLogin(msg)
	case StateUpdating:
		return m.updateUpdating(msg)
	case StateFailed:
		return m.updateFailed(msg)
	case StateSummary:
		return m.updateSummary(msg)
	}
	return m, nil
}

// ─── State: Init ────────────────────────────────────────────────────────────

func (m Model) updateInit(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case checkUpdateMsg:
		m.latestVersion = string(msg)
		return m, nil
	case detectDoneMsg:
		m.managers = msg.managers
		if len(m.managers) == 0 {
			// Nothing to do — skip straight to summary.
			m.state = StateSummary
			return m, nil
		}
		m.state = StateLogin
		return m, textinput.Blink
	}
	// Forward spinner messages while detecting.
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// ─── State: Login ───────────────────────────────────────────────────────────

func (m Model) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case checkUpdateMsg:
		m.latestVersion = string(msg)
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, keys.Submit):
			// Capture password as []byte (never as string in our domain).
			m.password = []byte(m.textInput.Value())
			m.textInput.SetValue("") // Clear display immediately.
			m.textInput.Blur()
			m.state = StateUpdating
			m.currentMgr = 0
			return m, tea.Batch(
				m.startManagerUpdate(),
				components.AnimTickCmd(200*time.Millisecond),
			)
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// ─── State: Updating ────────────────────────────────────────────────────────

func (m Model) updateUpdating(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case outputLineMsg:
		m.lastLine = string(msg)
		return m, waitForOutput(m.outputCh)

	case outputDoneMsg:
		// Stream ended; result will arrive shortly.
		return m, nil

	case managerDoneMsg:
		m.results = append(m.results, msg.result)
		if !msg.result.Success {
			// Offer interactive retry.
			m.state = StateFailed
			return m, nil
		}
		return m, m.advanceToNextManager()

	case components.AnimTickMsg:
		m.animation.NextFrame()
		return m, components.AnimTickCmd(200 * time.Millisecond)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

// ─── State: Failed (interactive retry prompt) ──────────────────────────────

func (m Model) updateFailed(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Retry):
			// Drop to raw terminal for interactive resolution.
			return m, m.interactiveRetry()
		case key.Matches(msg, keys.Skip):
			return m, m.advanceToNextManager()
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			utils.ZeroBytes(m.password)
			return m, tea.Quit
		}

	case interactiveDoneMsg:
		// Overwrite the last (failed) result.
		if len(m.results) > 0 {
			last := &m.results[len(m.results)-1]
			if msg.err == nil {
				last.Success = true
				last.Error = ""
			} else {
				last.Error = utils.SanitizeError(msg.err.Error(), 120)
			}
		}
		return m, m.advanceToNextManager()
	}
	return m, nil
}

// ─── State: Summary ─────────────────────────────────────────────────────────

func (m Model) updateSummary(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Quit) {
			m.quitting = true
			utils.ZeroBytes(m.password)
			return m, tea.Quit
		}
		_ = msg
	}
	return m, nil
}

// ─── Helpers ────────────────────────────────────────────────────────────────

// startManagerUpdate launches the current manager's update in a goroutine
// and returns Cmds to listen for output and the final result.
func (m *Model) startManagerUpdate() tea.Cmd {
	mgr := m.managers[m.currentMgr]
	m.outputCh = make(chan string, 128)
	m.resultCh = make(chan pkgmanager.UpdateResult, 1)
	m.lastLine = ""

	pw := make([]byte, len(m.password))
	copy(pw, m.password)

	go func() {
		result := utils.RunManagerUpdate(pw, mgr, m.outputCh)
		utils.ZeroBytes(pw) // Zero the copy.
		m.resultCh <- result
	}()

	return tea.Batch(
		waitForOutput(m.outputCh),
		waitForResult(m.resultCh),
		m.spinner.Tick,
	)
}

// advanceToNextManager moves to the next manager or finishes.
func (m *Model) advanceToNextManager() tea.Cmd {
	m.currentMgr++
	if m.currentMgr < len(m.managers) {
		m.state = StateUpdating
		return tea.Batch(
			m.startManagerUpdate(),
			components.AnimTickCmd(200*time.Millisecond),
		)
	}
	// All done — zero password and show summary.
	utils.ZeroBytes(m.password)
	m.state = StateSummary
	return nil
}

// interactiveRetry suspends the TUI and re-runs the failed manager's
// commands with full terminal I/O so the user can answer any prompts.
func (m *Model) interactiveRetry() tea.Cmd {
	mgr := m.managers[m.currentMgr]
	var parts []string
	for _, cmd := range mgr.Commands() {
		parts = append(parts, strings.Join(cmd, " "))
	}
	script := strings.Join(parts, " && ")

	var c *exec.Cmd
	if mgr.NeedsSudo() {
		c = exec.Command("sudo", "bash", "-c", script)
	} else {
		c = exec.Command("bash", "-c", script)
	}

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return interactiveDoneMsg{err: err}
	})
}

// waitForOutput returns a Cmd that blocks until a line arrives on ch.
func waitForOutput(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-ch
		if !ok {
			return outputDoneMsg{}
		}
		return outputLineMsg(line)
	}
}

// waitForResult returns a Cmd that blocks until the update result arrives.
func waitForResult(ch <-chan pkgmanager.UpdateResult) tea.Cmd {
	return func() tea.Msg {
		return managerDoneMsg{result: <-ch}
	}
}

// View renders the current state.
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	switch m.state {
	case StateInit:
		return viewInit(m)
	case StateLogin:
		return viewLogin(m)
	case StateUpdating:
		return viewUpdating(m)
	case StateFailed:
		return viewFailed(m)
	case StateSummary:
		return viewSummary(m)
	}
	return ""
}

// managerNames returns a formatted string of detected manager names.
func (m Model) managerNames() string {
	names := make([]string, len(m.managers))
	for i, mgr := range m.managers {
		names[i] = fmt.Sprintf("• %s", mgr.Name())
	}
	return strings.Join(names, "  ")
}
