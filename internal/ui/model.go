package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Jeeshan-Shrestha/mutux.git/internal/library"
)

type pane int

const (
	paneDirs pane = iota
	paneSongs
	paneQueue
	paneCount
)

type Model struct {
	root    string
	cwd     string
	entries []library.Entry
	cursor  int
	offset  int
	err     error

	width, height int
	focus         pane
}

func New(root string) Model {
	entries, err := library.ReadDir(root)
	return Model{root: root, cwd: root, entries: entries, err: err}
}

func (m Model) Init() tea.Cmd { return nil }
