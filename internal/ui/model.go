package ui

import (
	"time"

	"github.com/Jeeshan-Shrestha/mutux.git/internal/player"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/Jeeshan-Shrestha/mutux.git/internal/library"
)

type pane int

const (
	paneDirs pane = iota
	paneSongs
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

	player  *player.Player
	queue   []library.Entry // snapshot of the songs list when playback started
	qIdx    int
	nowPath string
	meta    library.Meta
	info    player.Info
	size    int64
	art     *artState
	pos     time.Duration
	playErr string

	width, height int
	focus         pane
}

func New(root string) Model {
	m := Model{root: root, player: player.New()}
	m.load(root)
	return m
}

// load reads dir into the panes and resets both cursors. On failure the
// current directory is left untouched and err is set.
func (m *Model) load(dir string) bool {
	all, err := library.ReadDir(dir)
	if err != nil {
		m.err = err
		return false
	}
	m.err = nil
	m.cwd = dir
	m.entries, m.songs = splitEntries(all)
	m.cursor, m.offset = 0, 0
	m.songCur, m.songOff = 0, 0
	return true
}

// enter descends into the highlighted directory.
func (m *Model) enter() {
	if m.cursor < 0 || m.cursor >= len(m.entries) {
		return
	}
	m.load(m.entries[m.cursor].Path)
}

// leave goes up to the parent (never above root) and re-selects the
// directory we just came out of.
func (m *Model) leave() {
	if m.cwd == m.root {
		return
	}
	from := m.cwd
	if !m.load(filepath.Dir(from)) {
		return
	}
	for i, e := range m.entries {
		if e.Path == from {
			m.cursor = i
			break
		}
	}
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

func (m Model) Init() tea.Cmd { return tick() }
