package tests

import (
	"testing"
	"github.com/jstreitb/baa/internal/sanitize"
)

func TestSanitizeError(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		maxLen   int
		expected string
	}{
		{
			name:     "Removes sudo password prompt",
			raw:      "[sudo] password for user: \nError: failed to sync",
			maxLen:   100,
			expected: "user: \nError: failed to sync",
		},
		{
			name:     "Removes try again noise",
			raw:      "Sorry, try again.\nConnection refused",
			maxLen:   100,
			expected: "Connection refused",
		},
		{
			name:     "Truncates long strings",
			raw:      "This is a very long error message that should be truncated",
			maxLen:   10,
			expected: "… truncated",
		},
		{
			name:     "Trims space",
			raw:      "   error message   ",
			maxLen:   100,
			expected: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitize.SanitizeError(tt.raw, tt.maxLen)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
