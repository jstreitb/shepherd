package sanitize

import "strings"

// SanitizeError truncates and cleans an error message for safe display
// in the TUI. It strips sudo password prompts and limits length.
func SanitizeError(raw string, maxLen int) string {
	// Remove common sudo noise.
	replacer := strings.NewReplacer(
		"[sudo] password for ", "",
		"Sorry, try again.", "",
	)
	clean := replacer.Replace(raw)
	clean = strings.TrimSpace(clean)

	if len(clean) > maxLen {
		clean = "…" + clean[len(clean)-maxLen:]
	}
	return clean
}
