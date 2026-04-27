package ui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jstreitb/baa/internal/config"
	"github.com/jstreitb/baa/internal/detector"
	"github.com/jstreitb/baa/internal/pkgmanager"
	"github.com/jstreitb/baa/internal/sanitize"
	"github.com/jstreitb/baa/internal/theme"
	"github.com/jstreitb/baa/internal/ui/components"
	"github.com/jstreitb/baa/internal/updater"
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

// ─── Model ─────────────────────────────────────────────────────────────────

// Model is the top-level Bubbletea model for BAA.
// It is a thin UI shell that delegates execution concerns to the Orchestrator.
type Model struct {
	state  State
	width  int
	height int

	// Update info
	latestVersion string

	// Login screen.
	textInput textinput.Model

	// Display components.
	animation components.Animation
	spinner   spinner.Model

	quitting bool

	// Injected dependencies.
	cfg      config.Config
	checker  updater.VersionChecker
	detector *detector.Detector

	// Execution orchestration (pointer: shared across Bubble Tea copies).
	orch *Orchestrator
}

// NewModel constructs the initial Model ready for tea.NewProgram.
func NewModel(cfg config.Config, checker updater.VersionChecker, det *detector.Detector) Model {
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
		spinner:   components.NewSpinner(theme.ColorMauve),
		cfg:       cfg,
		checker:   checker,
		detector:  det,
		orch:      NewOrchestrator(),
	}
}

// Init starts detection and the text input cursor blink.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.detectManagersCmd(),
		m.checkUpdateCmd(),
		m.spinner.Tick,
	)
}

func (m *Model) checkUpdateCmd() tea.Cmd {
	return func() tea.Msg {
		if m.checker == nil {
			return checkUpdateMsg("")
		}
		latest, err := m.checker.LatestVersion(context.Background(), m.cfg.Version)
		if err != nil {
			return checkUpdateMsg("")
		}
		return checkUpdateMsg(latest)
	}
}

// detectManagersCmd runs package manager detection in a background goroutine.
func (m *Model) detectManagersCmd() tea.Cmd {
	return func() tea.Msg {
		if m.detector == nil {
			return detectDoneMsg{managers: nil}
		}
		return detectDoneMsg{managers: m.detector.DetectInstalled()}
	}
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
			m.orch.Cleanup()
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
		m.orch.SetManagers(msg.managers)
		if len(msg.managers) == 0 {
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
			// Capture password securely.
			raw := []byte(m.textInput.Value())
			m.textInput.SetValue("") // Clear display immediately.
			m.textInput.Blur()
			m.orch.SetPassword(raw)
			m.orch.BeginUpdating()
			m.state = StateUpdating
			return m, tea.Batch(
				m.orch.StartUpdate(),
				m.spinner.Tick,
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
		m.orch.SetLastLine(string(msg))
		return m, waitForOutput(m.orch.OutputCh())

	case outputDoneMsg:
		// Stream ended; result will arrive shortly.
		return m, nil

	case managerDoneMsg:
		ok := m.orch.AddResult(msg.result)
		if !ok {
			// Offer interactive retry.
			m.state = StateFailed
			return m, nil
		}
		state, cmd := m.orch.Advance()
		m.state = state
		return m, cmd

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
			return m, m.orch.InteractiveRetry()
		case key.Matches(msg, keys.Skip):
			state, cmd := m.orch.Advance()
			m.state = state
			return m, cmd
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			m.orch.Cleanup()
			return m, tea.Quit
		}

	case interactiveDoneMsg:
		if msg.err == nil {
			m.orch.OverwriteLastResult(true, "")
		} else {
			m.orch.OverwriteLastResult(false, sanitize.SanitizeError(msg.err.Error(), 120))
		}
		state, cmd := m.orch.Advance()
		m.state = state
		return m, cmd
	}
	return m, nil
}

// ─── State: Summary ─────────────────────────────────────────────────────────

func (m Model) updateSummary(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Quit) {
			m.quitting = true
			m.orch.Cleanup()
			return m, tea.Quit
		}
	}
	return m, nil
}

// ─── View ───────────────────────────────────────────────────────────────────

// View renders the current state by projecting Model into typed ViewData.
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	switch m.state {
	case StateInit:
		return viewInit(m.initViewData())
	case StateLogin:
		return viewLogin(m.loginViewData())
	case StateUpdating:
		return viewUpdating(m.updatingViewData())
	case StateFailed:
		return viewFailed(m.failedViewData())
	case StateSummary:
		return viewSummary(m.summaryViewData())
	}
	return ""
}

// ─── ViewData Projections ──────────────────────────────────────────────────

func (m Model) layout() Layout {
	return Layout{Width: m.width, Height: m.height}
}

func (m Model) initViewData() InitViewData {
	return InitViewData{
		Layout:      m.layout(),
		SpinnerView: m.spinner.View(),
	}
}

func (m Model) loginViewData() LoginViewData {
	mgrs := m.orch.Managers()
	names := make([]string, 0, len(mgrs))
	for i, mgr := range mgrs {
		if i == 3 {
			names = append(names, "...")
			break
		}
		names = append(names, mgr.Name())
	}
	return LoginViewData{
		Layout:        m.layout(),
		ManagerNames:  names,
		TextInputView: m.textInput.View(),
		LatestVersion: m.latestVersion,
	}
}

func (m Model) updatingViewData() UpdatingViewData {
	mgr := m.orch.CurrentManager()
	name := ""
	if mgr != nil {
		name = mgr.Name()
	}
	return UpdatingViewData{
		Layout:         m.layout(),
		ManagerName:    name,
		CurrentIndex:   m.orch.CurrentIndex(),
		TotalManagers:  len(m.orch.Managers()),
		AnimationFrame: m.animation.Frame(),
		LastLogLine:    m.orch.LastLine(),
		SpinnerView:    m.spinner.View(),
		PastResults:    toViewResults(m.orch.Results()),
	}
}

func (m Model) failedViewData() FailedViewData {
	results := m.orch.Results()
	d := FailedViewData{Layout: m.layout()}
	if len(results) > 0 {
		last := results[len(results)-1]
		d.ManagerName = last.Manager
		d.ErrorMsg = last.Error
	}
	return d
}

func (m Model) summaryViewData() SummaryViewData {
	return SummaryViewData{
		Layout:        m.layout(),
		HasManagers:   len(m.orch.Managers()) > 0,
		Results:       toViewResults(m.orch.Results()),
		LatestVersion: m.latestVersion,
	}
}

// toViewResults converts domain results to presentation-only structs.
func toViewResults(results []pkgmanager.UpdateResult) []ViewResult {
	out := make([]ViewResult, len(results))
	for i, r := range results {
		out[i] = ViewResult{
			Manager:  r.Manager,
			Success:  r.Success,
			Error:    r.Error,
			Duration: r.Duration,
		}
	}
	return out
}

// ─── Channel Helpers ────────────────────────────────────────────────────────

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

