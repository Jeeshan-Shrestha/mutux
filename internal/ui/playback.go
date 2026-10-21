package ui

import (
	"fmt"
	"time"

	"github.com/Jeeshan-Shrestha/mutux.git/internal/library"
	"github.com/Jeeshan-Shrestha/mutux.git/internal/player"
)

// playQueue starts the queue (a snapshot of the songs list) at idx.
func (m *Model) playQueue(queue []library.Entry, idx int) {
	m.queue = append([]library.Entry(nil), queue...)
	m.play(idx)
}

// play loads queue[idx]; on failure the error is shown and nothing changes.
func (m *Model) play(idx int) bool {
	if idx < 0 || idx >= len(m.queue) {
		return false
	}
	e := m.queue[idx]

	if !player.Supported(e.Path) {
		m.playErr = "can't play " + e.Name + ": unsupported format"
		return false
	}
	info, err := m.player.Play(e.Path)
	if err != nil {
		m.playErr = fmt.Sprintf("can't play %s: %v", e.Name, err)
		return false
	}

	m.playErr = ""
	m.qIdx = idx
	m.nowPath = e.Path
	m.info = info
	m.pos = 0
	m.meta = library.ReadMeta(e.Path)
	m.art = loadArt(e.Path, m.meta.Art)
	m.size = e.Size
	return true
}

// next plays the following playable track; at the end of the queue it stops.
func (m *Model) next() {
	for i := m.qIdx + 1; i < len(m.queue); i++ {
		if m.play(i) {
			return
		}
	}
	if m.player.Loaded() {
		m.player.Stop()
		m.pos = 0
	}
}

// prev restarts the current track, or goes back one if it just started.
func (m *Model) prev() {
	if m.player.Loaded() && m.player.Position() > 3*time.Second {
		m.play(m.qIdx)
		return
	}
	for i := m.qIdx - 1; i >= 0; i-- {
		if m.play(i) {
			return
		}
	}
}

// title is the tag title, or the file name when untagged.
func (m Model) title() string {
	if m.meta.Title != "" {
		return m.meta.Title
	}
	if m.qIdx >= 0 && m.qIdx < len(m.queue) {
		return m.queue[m.qIdx].Name
	}
	return ""
}

func fmtDur(d time.Duration) string {
	s := int(d.Round(time.Second) / time.Second)
	if s >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", s/3600, s/60%60, s%60)
	}
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}
