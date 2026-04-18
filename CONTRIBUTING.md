# Contributing to Shepherd

First off — thank you for considering contributing to Shepherd! Every contribution matters, whether it's a bug fix, a new package manager, or a typo correction.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Adding a New Package Manager](#adding-a-new-package-manager)
- [Code Style](#code-style)
- [Commit Messages](#commit-messages)
- [Pull Requests](#pull-requests)

## Getting Started

```bash
# Fork and clone
git clone https://github.com/<your-username>/shepherd.git
cd shepherd

# Install dependencies
go mod download

# Build
make build

# Run
make run
```

**Requirements:**
- Go 1.21+
- A Linux system (for testing package managers)

## Development Workflow

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feat/my-feature
   ```
2. Make your changes.
3. Ensure the project builds and passes checks:
   ```bash
   make build
   go vet ./...
   ```
4. Commit with a descriptive message.
5. Push and open a Pull Request.

## Adding a New Package Manager

Shepherd makes it easy to add support for new package managers. Here's how:

1. **Create a new file** in `internal/pkgmanager/` (e.g., `yay.go`).
2. **Implement the `PackageManager` interface:**

```go
package pkgmanager

type Yay struct{}

func (y *Yay) Name() string       { return "yay" }
func (y *Yay) NeedsSudo() bool     { return false }
func (y *Yay) Commands() [][]string {
    return [][]string{
        {"yay", "-Syu", "--noconfirm"},
    }
}
func (y *Yay) Env() []string { return nil }
```

3. **Register it in `detect.go`** by adding an entry to the `candidates` slice:

```go
{"yay", func() PackageManager { return &Yay{} }},
```

4. **Test** on a system that has the package manager installed.
5. **Open a PR** with your changes.

## Code Style

- Follow standard Go formatting (`gofmt`).
- Keep functions small and focused.
- Comment exported types and functions.
- Use meaningful variable names.
- Security-sensitive code (passwords, sudo) must follow the patterns in `internal/utils/sudo.go`.

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add yay package manager support
fix: handle apt lock file gracefully
docs: update README with new features
refactor: simplify animation frame rendering
```

## Pull Requests

- Fill out the PR template completely.
- Keep PRs focused — one feature or fix per PR.
- Ensure `go build ./...` and `go vet ./...` pass.
- Add a brief description of what you changed and why.

---

Thank you for helping make Shepherd better! 🐑
