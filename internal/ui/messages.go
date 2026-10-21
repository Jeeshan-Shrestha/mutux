package ui

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

type tickMsg time.Time

// tick drives the seekbar and auto-advance.
func tick() tea.Cmd {
	return tea.Tick(250*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}
