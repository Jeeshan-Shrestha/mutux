package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	minW = 60
	minH = 18

	topH = 3 // outer height of the options bar (1 content row)
	botH = 4 // outer height of the player bar (2 content rows)
)

func (m Model) View() tea.View {
	var s string
	if m.width < minW || m.height < minH {
		s = fmt.Sprintf("terminal too small (%dx%d), need at least %dx%d",
			m.width, m.height, minW, minH)
	} else {
		s = m.layout()
	}

	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func (m Model) layout() string {
	midH := m.height - topH - botH

	leftW := m.width * 25 / 100
	midW := m.width * 40 / 100
	rightW := m.width - leftW - midW // remainder, so no gap from rounding

	rightTopH := midH * 55 / 100
	rightBotH := midH - rightTopH

	top := box("", "options  |  tab: switch pane  |  q: quit", m.width, topH, false)

	dirs := box("Directories", m.listView(), leftW, midH, m.focus == paneDirs)
	songs := box("Songs", "", midW, midH, m.focus == paneSongs)

	info := box("Now Playing",
		"album art\nformat / bitrate\nvolume\nduration\nartist",
		rightW, rightTopH, false)
	queue := box("Queue / Lyrics", "", rightW, rightBotH, m.focus == paneQueue)
	right := lipgloss.JoinVertical(lipgloss.Left, info, queue)

	mid := lipgloss.JoinHorizontal(lipgloss.Top, dirs, songs, right)

	bottom := box("", "seekbar\ncontrols", m.width, botH, false)

	return lipgloss.JoinVertical(lipgloss.Left, top, mid, bottom)
}

func (m Model) listView() string {
	if m.err != nil {
		return "error: " + m.err.Error()
	}
	if len(m.entries) == 0 {
		return "(empty)"
	}

	var b strings.Builder
	for i, e := range m.entries {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		name := e.Name
		if e.IsDir {
			name += "/"
		}
		b.WriteString(prefix)
		b.WriteString(name)
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}
