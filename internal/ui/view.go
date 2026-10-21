package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
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

// rightHeights splits the middle band between Now Playing and Album Art.
func (m Model) rightHeights() (top, bot int) {
	midH := m.height - topH - botH
	top = midH * 40 / 100
	return top, midH - top
}

func (m Model) layout() string {
	midH := m.height - topH - botH

	leftW := m.width * 22 / 100
	midW := m.width * 33 / 100
	rightW := m.width - leftW - midW // remainder, so no gap from rounding

	rightTopH, rightBotH := m.rightHeights()

	hint := "tab/ctrl+h/l: pane | enter: play | space: pause | n/p: next/prev | </>: seek | q: quit"
	top := box("", ansi.Truncate(hint, m.width-2, "…"), m.width, topH, false)

	dirs := box(m.dirsTitle(), m.dirsBody(leftW-2), leftW, midH, m.focus == paneDirs)
	songs := box(m.songsTitle(), m.songsBody(midW-2), midW, midH, m.focus == paneSongs)

	info := box("Now Playing", m.nowPlayingBody(rightW-2), rightW, rightTopH, false)
	// body gets the pane minus the 2 border rows and 1 title row
	art := box("Album Art", m.art.render(rightW-2, rightBotH-3), rightW, rightBotH, false)
	right := lipgloss.JoinVertical(lipgloss.Left, info, art)

	mid := lipgloss.JoinHorizontal(lipgloss.Top, dirs, songs, right)

	bottom := box("", m.controlsBody(m.width-2), m.width, botH, false)

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
		items[i] = "  " + e.Name
		if m.player.Loaded() && e.Path == m.nowPath {
			items[i] = "▶ " + e.Name
		}
	}

	return renderRows(items, m.songCur, m.songOff, innerW, m.dirsRows(), m.focus == paneSongs)
}
