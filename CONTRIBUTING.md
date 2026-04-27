# Contributing to BAA

First off — thank you for considering contributing to BAA! Every contribution matters, whether it's a bug fix, a new package manager, or a typo correction.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Adding a New Package Manager](#adding-a-new-package-manager)
- [Code Architecture](#code-architecture)
- [Code Style](#code-style)
- [Commit Messages](#commit-messages)
- [Pull Requests](#pull-requests)

## Getting Started

```bash
# Fork and clone
git clone https://github.com/<your-username>/baa.git
cd baa

# Install dependencies
go mod download

# Build
go build -o baa ./cmd/baa/

# Run Tests
go test -v ./...
```

**Requirements:**
- Go 1.21+
- A Linux system (for testing package managers)

## Development Workflow

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feat/my-feature
   ```
2. Make your changes in the appropriate `internal/` packages.
3. Ensure the project builds and passes checks (including your new tests):
   ```bash
   go vet ./...
   go test -v ./...
   ```
4. Commit with a descriptive message.
5. Push and open a Pull Request.

## Adding a New Package Manager

BAA makes it easy to add support for new package managers. The architecture dictates that providers are decoupled from the core application. Here's how:

1. **Create a new file** in `internal/pkgmanager/providers/` (e.g., `yay.go`).
2. **Implement the `PackageManager` interface:**

```go
package providers

import "github.com/jstreitb/baa/internal/pkgmanager"

// Yay implements PackageManager for yay.
type Yay struct{}

func (y *Yay) Name() string       { return "yay" }
func (y *Yay) NeedsSudo() bool    { return false }
func (y *Yay) Commands() [][]string {
    return [][]string{
        {"yay", "-Syu", "--noconfirm"},
    }
}
func (y *Yay) Env() []string { return nil }
```

3. **Register it in the detector** by adding an entry to the `candidates` slice inside `internal/detector/detector.go`:

```go
{"yay", func() pkgmanager.PackageManager { return &providers.Yay{} }},
```

4. **Add a test** for your provider in `tests/providers_test.go` using the table-driven test format.
5. **Test** on a system that has the package manager installed.
6. **Open a PR** with your changes!

## Code Architecture

BAA follows a strict modular Clean Architecture. Please maintain these boundaries:
- **`cmd/baa/`**: The Composition Root. Only wiring dependencies together.
- **`internal/ui/`**: Bubble Tea components. The `Model` is a thin shell. State execution is handled by `Orchestrator`, and visual projections by `ViewData`.
- **`internal/executor/`**: Secure shell execution and memory clearing.
- **`tests/`**: All unit and integration tests live here to ensure they test public APIs and avoid circular dependencies.

## Code Style

- Follow standard Go formatting (`gofmt`).
- Keep functions small and focused on a single responsibility.
- Write unit tests for your changes in the `tests/` directory.
- Use meaningful variable names and document exported types.
- **Security-sensitive code**: Passwords must be handled using the `SecurePassword` struct in `internal/executor/secure.go` and zeroed explicitly via `ZeroBytes`.

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add yay package manager support
fix: handle apt lock file gracefully
docs: update README with new architecture
test: add unit tests for executor module
refactor: separate UI state from bubble tea model
```

## Pull Requests

- Fill out the PR template completely.
- Keep PRs focused — one feature or fix per PR.
- Ensure `go build ./...`, `go vet ./...`, and `go test ./...` all pass.
- Add a brief description of what you changed and why, specifically calling out any architectural decisions.

---

Thank you for helping make BAA better! 🐑
