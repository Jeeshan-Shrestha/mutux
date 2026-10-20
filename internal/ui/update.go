package ui

import tea "charm.land/bubbletea/v2"

// dirsRows is how many list rows fit in the directories pane:
// pane height minus 2 border rows and 1 title row.
func (m Model) dirsRows() int {
	r := (m.height - topH - botH) - 3
	if r < 1 {
		r = 1
	}
	return r
}

// clampScroll keeps cursor inside the window of rows starting at offset.
func clampScroll(cursor, offset, rows, n int) int {
	if cursor < offset {
		offset = cursor
	}
	if cursor >= offset+rows {
		offset = cursor - rows + 1
	}
	if maxOff := n - rows; offset > maxOff {
		offset = maxOff
	}
	if offset < 0 {
		offset = 0
	}
	return offset
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab", "ctrl+l":
			m.focus = (m.focus + 1) % paneCount
		case "shift+tab", "ctrl+h":
			m.focus = (m.focus + paneCount - 1) % paneCount

		case "j", "down":
			m.move(1)
		case "k", "up":
			m.move(-1)
		case "g":
			m.move(-1 << 30)
		case "G":
			m.move(1 << 30)
		}
	}

	m.offset = clampScroll(m.cursor, m.offset, m.dirsRows(), len(m.entries))
	m.songOff = clampScroll(m.songCur, m.songOff, m.dirsRows(), len(m.songs))
	return m, nil
}

// move shifts the cursor of the focused list by delta, clamped to bounds.
func (m *Model) move(delta int) {
	clamp := func(c, n int) int {
		if c >= n {
			c = n - 1
		}
		if c < 0 {
			c = 0
		}
		return c
	}
	switch m.focus {
	case paneDirs:
		m.cursor = clamp(m.cursor+delta, len(m.entries))
	case paneSongs:
		m.songCur = clamp(m.songCur+delta, len(m.songs))
	}
}
