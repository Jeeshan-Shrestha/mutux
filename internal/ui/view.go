package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m Model) View() tea.View {
	var str_builder strings.Builder

	str_builder.WriteString(m.cwd)
	str_builder.WriteString("\n\n")

	if m.err != nil {
		str_builder.WriteString("error: ")
		str_builder.WriteString(m.err.Error())
		str_builder.WriteString("\n")
	} else if len(m.entries) == 0 {
		str_builder.WriteString("  (empty)\n")
	}

	for i, e := range m.entries {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		name := e.Name
		if e.IsDir {
			name += "/"
		}
		str_builder.WriteString(prefix)
		str_builder.WriteString(name)
		str_builder.WriteString("\n")
	}

	fmt.Fprintf(&str_builder, "\n%d/%d  j/k: move  q: quit\n", m.cursor+1, len(m.entries))

	v := tea.NewView(str_builder.String())
	v.AltScreen = true
	return v
}
