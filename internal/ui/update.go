package ui

import tea "charm.land/bubbletea/v2"

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
	return m, nil
}
