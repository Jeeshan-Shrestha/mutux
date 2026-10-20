package ui

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

var (
	colBorder = lipgloss.Color("240")
	colFocus  = lipgloss.Color("212")

	titleStyle = lipgloss.NewStyle().Bold(true)
)

/* box draws a bordered pane that is exactly outerW x outerH.
*  the inner block is sized and clipped first (no border on it), then the
* border is wrapped around it, so the math is: inner = outer - 2.*/
func box(title, body string, outerW, outerH int, focused bool) string {
	iw, ih := outerW-2, outerH-2
	if iw < 1 || ih < 1 {
		return ""
	}

	content := body
	if title != "" {
		content = titleStyle.Render(ansi.Truncate(title, iw, "…"))
		if body != "" {
			content += "\n" + body
		}
	}

	inner := lipgloss.NewStyle().
		Width(iw).Height(ih).
		MaxWidth(iw).MaxHeight(ih).
		Render(content)

	bc := colBorder
	if focused {
		bc = colFocus
	}

	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(bc).
		Render(inner)
}
