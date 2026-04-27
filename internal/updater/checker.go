package updater

import (
	"context"
	"net/http"
	"strings"
	"time"
)

// VersionChecker resolves the latest available version of the application.
type VersionChecker interface {
	LatestVersion(ctx context.Context, currentVersion string) (string, error)
}

// ─── GitHubChecker ──────────────────────────────────────────────────────────

// GitHubChecker queries the GitHub releases API to find the latest version.
// It contains no test-specific logic — use StaticChecker for testing.
type GitHubChecker struct {
	Repo string
}

// NewGitHubChecker creates a GitHubChecker for the given "owner/repo" string.
func NewGitHubChecker(repo string) *GitHubChecker {
	return &GitHubChecker{Repo: repo}
}

func (c *GitHubChecker) LatestVersion(ctx context.Context, currentVersion string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirect, just read Location
		},
		Timeout: 2 * time.Second,
	}

	url := "https://github.com/" + c.Repo + "/releases/latest"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	loc, err := res.Location()
	if err == nil && loc != nil && strings.Contains(loc.Path, "/tag/") {
		parts := strings.Split(loc.Path, "/")
		latest := parts[len(parts)-1]
		latest = strings.TrimPrefix(latest, "v")
		curr := strings.TrimPrefix(currentVersion, "v")

		if latest != "" && latest != curr {
			return latest, nil
		}
	}
	return "", nil
}

// ─── StaticChecker ──────────────────────────────────────────────────────────

// StaticChecker always returns a fixed version string.
// Use this for testing or demo modes instead of polluting GitHubChecker.
type StaticChecker struct {
	Version string
}

func (s *StaticChecker) LatestVersion(_ context.Context, _ string) (string, error) {
	return s.Version, nil
}
