package tests

import (
	"errors"
	"testing"

	"github.com/jstreitb/baa/internal/detector"
)

type MockRunner struct {
	paths map[string]string
}

func (m *MockRunner) LookPath(file string) (string, error) {
	if path, ok := m.paths[file]; ok {
		return path, nil
	}
	return "", errors.New("executable file not found in $PATH")
}

func TestDetectInstalled(t *testing.T) {
	runner := &MockRunner{
		paths: map[string]string{
			"apt-get": "/usr/bin/apt-get",
			"brew":    "/opt/homebrew/bin/brew",
		},
	}

	det := detector.New(runner)
	found := det.DetectInstalled()

	if len(found) != 2 {
		t.Fatalf("expected 2 package managers, got %d", len(found))
	}

	if found[0].Name() != "apt" {
		t.Errorf("expected first manager to be apt, got %s", found[0].Name())
	}
	if found[1].Name() != "brew" {
		t.Errorf("expected second manager to be brew, got %s", found[1].Name())
	}
}
