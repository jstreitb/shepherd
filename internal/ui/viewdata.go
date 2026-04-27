package ui

import "time"

// ─── Layout ─────────────────────────────────────────────────────────────────

// Layout holds terminal dimensions shared by all views.
type Layout struct {
	Width  int
	Height int
}

// ─── View Data Structs ──────────────────────────────────────────────────────
// Each view function receives only the data it needs, preventing
// accidental coupling to unrelated Model fields.

// InitViewData carries everything the init/detection view needs.
type InitViewData struct {
	Layout
	SpinnerView string
}

// LoginViewData carries everything the login/password view needs.
type LoginViewData struct {
	Layout
	ManagerNames  []string
	TextInputView string
	LatestVersion string
}

// UpdatingViewData carries everything the update-in-progress view needs.
type UpdatingViewData struct {
	Layout
	ManagerName    string
	CurrentIndex   int
	TotalManagers  int
	AnimationFrame string
	LastLogLine    string
	SpinnerView    string
	PastResults    []ViewResult
}

// FailedViewData carries everything the failure/retry view needs.
type FailedViewData struct {
	Layout
	ManagerName string
	ErrorMsg    string
}

// SummaryViewData carries everything the summary view needs.
type SummaryViewData struct {
	Layout
	HasManagers   bool
	Results       []ViewResult
	LatestVersion string
}

// ViewResult is a presentation-ready slice of a pkgmanager.UpdateResult.
type ViewResult struct {
	Manager  string
	Success  bool
	Error    string
	Duration time.Duration
}
