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

// scrollToCursor keeps the cursor inside the visible window.
func (m *Model) scrollToCursor() {
	rows := m.dirsRows()

	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+rows {
		m.offset = m.cursor - rows + 1
	}

	maxOff := len(m.entries) - rows
	if maxOff < 0 {
		maxOff = 0
	}
	if m.offset > maxOff {
		m.offset = maxOff
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.focus = (m.focus + 1) % paneCount
		case "shift+tab":
			m.focus = (m.focus + paneCount - 1) % paneCount

		case "j", "down":
			if m.focus == paneDirs && m.cursor < len(m.entries)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.focus == paneDirs && m.cursor > 0 {
				m.cursor--
			}
		case "g":
			if m.focus == paneDirs {
				m.cursor = 0
			}
		case "G":
			if m.focus == paneDirs && len(m.entries) > 0 {
				m.cursor = len(m.entries) - 1
			}
		}
	}

	m.scrollToCursor()
	return m, nil
}
