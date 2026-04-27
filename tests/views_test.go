package tests

import (
	"testing"

	"github.com/jstreitb/baa/internal/ui"
)

func TestLastErrorLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single line",
			input:    "error: foo",
			expected: "error: foo",
		},
		{
			name:     "multiple lines",
			input:    "line 1\nline 2\nline 3",
			expected: "line 3",
		},
		{
			name:     "trailing newlines",
			input:    "line 1\nline 2\n\n\n",
			expected: "line 2",
		},
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   \n  \t  \n",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ui.LastErrorLine(tt.input); got != tt.expected {
				t.Errorf("LastErrorLine(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
