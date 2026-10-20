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
	entries []library.Entry // sub-directories of cwd
	cursor  int
	offset  int
	songs   []library.Entry // audio files of cwd
	songCur int
	songOff int
	err     error

	width, height int
	focus         pane
}

func New(root string) Model {
	all, err := library.ReadDir(root)
	dirs, songs := splitEntries(all)
	return Model{root: root, cwd: root, entries: dirs, songs: songs, err: err}
}

// splitEntries separates directories from audio files, keeping order.
func splitEntries(all []library.Entry) (dirs, songs []library.Entry) {
	for _, e := range all {
		if e.IsDir {
			dirs = append(dirs, e)
		} else {
			songs = append(songs, e)
		}
	}
	return dirs, songs
}

func (m Model) Init() tea.Cmd { return nil }
