package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

var (
	labelStyle   = lipgloss.NewStyle().Faint(true)
	currentStyle = lipgloss.NewStyle().Bold(true).Foreground(colFocus)
	faintStyle   = lipgloss.NewStyle().Faint(true)
	barFill      = lipgloss.NewStyle().Foreground(colFocus)
)

// field renders "label  value", truncated to w. Empty values are skipped.
func field(label, value string, w int) string {
	if value == "" {
		return ""
	}
	return ansi.Truncate(labelStyle.Render(fmt.Sprintf("%-9s", label))+value, w, "…")
}

func (m Model) nowPlayingBody(w int) string {
	if !m.player.Loaded() && m.nowPath == "" {
		if m.playErr != "" {
			return ansi.Truncate(m.playErr, w, "…")
		}
		return faintStyle.Render(" nothing playing — press enter on a song")
	}

	var track string
	if m.meta.Track > 0 {
		track = fmt.Sprint(m.meta.Track)
		if m.meta.TrackTotal > 0 {
			track += fmt.Sprintf(" / %d", m.meta.TrackTotal)
		}
	}
	var year string
	if m.meta.Year > 0 {
		year = fmt.Sprint(m.meta.Year)
	}

	format := m.meta.Format
	if format == "" && m.qIdx < len(m.queue) {
		format = strings.ToUpper(strings.TrimPrefix(extOf(m.queue[m.qIdx].Name), "."))
	}
	var quality string
	if m.info.Length > 0 && m.size > 0 {
		kbps := int(float64(m.size) * 8 / m.info.Length.Seconds() / 1000)
		quality = fmt.Sprintf("%d kbps · %.1f kHz · %s", kbps,
			float64(m.info.SampleRate)/1000, channelName(m.info.Channels))
	}

	lines := []string{
		ansi.Truncate(lipgloss.NewStyle().Bold(true).Render(m.title()), w, "…"),
		field("Artist", m.meta.Artist, w),
		field("Album", m.meta.Album, w),
		field("Album by", differs(m.meta.AlbumArtist, m.meta.Artist), w),
		field("Year", year, w),
		field("Genre", m.meta.Genre, w),
		field("Track", track, w),
		field("Composer", m.meta.Composer, w),
		field("Format", format, w),
		field("Quality", quality, w),
		field("Length", fmtDur(m.info.Length), w),
		field("Size", fmt.Sprintf("%.1f MB", float64(m.size)/1e6), w),
	}
	if m.playErr != "" {
		lines = append([]string{ansi.Truncate(m.playErr, w, "…")}, lines...)
	}

	out := lines[:0]
	for _, l := range lines {
		if l != "" {
			out = append(out, l)
		}
	}
	return strings.Join(out, "\n")
}

func differs(a, b string) string {
	if a == b {
		return ""
	}
	return a
}

func channelName(n int) string {
	switch n {
	case 1:
		return "mono"
	case 2:
		return "stereo"
	}
	return fmt.Sprintf("%d ch", n)
}

func extOf(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return name[i:]
	}
	return ""
}

// controlsBody is the two-line bottom bar: seekbar and transport controls.
func (m Model) controlsBody(w int) string {
	cur, total := fmtDur(m.pos), fmtDur(m.info.Length)
	if !m.player.Loaded() {
		cur, total = "0:00", fmtDur(m.info.Length)
		if m.nowPath == "" {
			total = "0:00"
		}
	}

	bw := w - len(cur) - len(total) - 4
	if bw < 1 {
		bw = 1
	}
	filled := 0
	if m.info.Length > 0 {
		filled = int(float64(bw) * float64(m.pos) / float64(m.info.Length))
		if filled > bw {
			filled = bw
		}
	}
	bar := barFill.Render(strings.Repeat("█", filled)) +
		faintStyle.Render(strings.Repeat("░", bw-filled))
	seek := fmt.Sprintf(" %s %s %s", cur, bar, total)

	state := "■ stopped"
	switch {
	case m.player.Loaded() && m.player.Paused():
		state = "⏸ paused"
	case m.player.Loaded():
		state = "▶ playing"
	}

	controls := fmt.Sprintf(" %s   ⏮ p   ⏯ space   ⏭ n   « <   > »", currentStyle.Render(state))

	return seek + "\n" + ansi.Truncate(controls, w, "…")
}
