package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

var (
	rowSelFocused = lipgloss.NewStyle().
			Background(colFocus).
			Foreground(lipgloss.Color("0"))
	rowSelBlurred = lipgloss.NewStyle().
			Background(lipgloss.Color("238"))
)

// renderRows draws the visible slice of items starting at offset.
// Every row is truncated and padded to exactly w columns, so alignment
// never breaks and the highlight spans the full pane width.
func renderRows(items []string, cursor, offset, w, h int, focused bool) string {
	if w < 1 || h < 1 {
		return ""
	}

	lines := make([]string, 0, h)
	for i := offset; i < len(items) && len(lines) < h; i++ {
		text := ansi.Truncate(" "+items[i], w, "…")
		if pad := w - ansi.StringWidth(text); pad > 0 {
			text += strings.Repeat(" ", pad)
		}

		if i == cursor {
			if focused {
				text = rowSelFocused.Render(text)
			} else {
				text = rowSelBlurred.Render(text)
			}
		}
		lines = append(lines, text)
	}
	return strings.Join(lines, "\n")
}
