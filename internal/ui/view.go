package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	minW = 60
	minH = 18

	topH = 3
	botH = 4
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

	dirs := box(m.dirsTitle(), m.dirsBody(leftW-2), leftW, midH, m.focus == paneDirs)
	songs := box(m.songsTitle(), m.songsBody(midW-2), midW, midH, m.focus == paneSongs)

	info := box("Now Playing",
		"album art\nformat / bitrate\nvolume\nduration\nartist",
		rightW, rightTopH, false)
	queue := box("Queue / Lyrics", "", rightW, rightBotH, m.focus == paneQueue)
	right := lipgloss.JoinVertical(lipgloss.Left, info, queue)

	mid := lipgloss.JoinHorizontal(lipgloss.Top, dirs, songs, right)

	bottom := box("", "seekbar\ncontrols", m.width, botH, false)

	return lipgloss.JoinVertical(lipgloss.Left, top, mid, bottom)
}

func (m Model) dirsTitle() string {
	if len(m.entries) == 0 {
		return "Directories"
	}
	return fmt.Sprintf("Directories (%d/%d)", m.cursor+1, len(m.entries))
}

func (m Model) dirsBody(innerW int) string {
	if m.err != nil {
		return "error: " + m.err.Error()
	}
	if len(m.entries) == 0 {
		return " (empty)"
	}

	items := make([]string, len(m.entries))
	for i, e := range m.entries {
		items[i] = e.Name + "/"
	}

	return renderRows(items, m.cursor, m.offset, innerW, m.dirsRows(), m.focus == paneDirs)
}

func (m Model) songsTitle() string {
	if len(m.songs) == 0 {
		return "Songs"
	}
	return fmt.Sprintf("Songs (%d/%d)", m.songCur+1, len(m.songs))
}

func (m Model) songsBody(innerW int) string {
	if m.err != nil {
		return ""
	}
	if len(m.songs) == 0 {
		return " (no songs)"
	}

	items := make([]string, len(m.songs))
	for i, e := range m.songs {
		items[i] = e.Name
	}

	return renderRows(items, m.songCur, m.songOff, innerW, m.dirsRows(), m.focus == paneSongs)
}
