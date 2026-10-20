package ui

import (
	tea "charm.land/bubbletea/v2"
	music_scan "github.com/Jeeshan-Shrestha/mutux.git/internal/library"
)

type Model struct {
	root    string
	cwd     string
	entries []music_scan.Entry
	cursor  int
	err     error
}

func New(root string) Model {
	entries, err := music_scan.ReadDir(root)
	return Model{root: root, cwd: root, entries: entries, err: err}
}

func (m Model) Init() tea.Cmd {
	return nil
}
