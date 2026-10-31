package ui

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/Jeeshan-Shrestha/mutux.git/internal/library"
	"github.com/charmbracelet/x/ansi"
)

const searchLimit = 5

type searchMode int

const (
	searchOff searchMode = iota
	searchInput
	searchLoading
	searchResults
	searchDownloading
)

type result struct {
	title string
	url   string
}

type searchState struct {
	mode    searchMode
	query   string
	results []result
	picked  string // title being downloaded
	err     string
	status  string // shown in the hint bar after the search view closes

	all      []library.Entry // every song under root, loaded for /local
	loaded   bool
	matches  []library.Entry
	localCur int
}

type searchResultMsg struct {
	results []result
	err     error
}

type downloadDoneMsg struct {
	title string
	err   error
}

// runSearch asks yt-dlp for the top results. No shell is involved, so the
// query can't be interpreted as anything but a search string.
func runSearch(query string) tea.Cmd {
	return func() tea.Msg {
		out, err := exec.Command("yt-dlp", "--js-runtimes", "node",
			fmt.Sprintf("ytsearch%d:%s", searchLimit, query),
			"--print", "%(title)s|%(webpage_url)s").Output()
		if err != nil {
			return searchResultMsg{err: cmdErr(err)}
		}
		var res []result
		for _, line := range strings.Split(string(out), "\n") {
			i := strings.LastIndex(line, "|")
			if i < 0 {
				continue
			}
			res = append(res, result{title: line[:i], url: line[i+1:]})
		}
		return searchResultMsg{results: res}
	}
}

func download(r result, dir string) tea.Cmd {
	return func() tea.Msg {
		_, err := exec.Command("yt-dlp", "--js-runtimes", "node",
			"-o", dir+"/%(title)s.%(ext)s",
			"-x", "--audio-format", "mp3", r.url).Output()
		return downloadDoneMsg{title: r.title, err: cmdErr(err)}
	}
}

// cmdErr surfaces the last stderr line of a failed command.
func cmdErr(err error) error {
	if ee, ok := err.(*exec.ExitError); ok {
		lines := strings.Split(strings.TrimSpace(string(ee.Stderr)), "\n")
		if l := lines[len(lines)-1]; l != "" {
			return fmt.Errorf("%s", l)
		}
	}
	return err
}

// updateSearch handles every message related to the search view. It
// reports whether the message was consumed.
func (m Model) updateSearch(msg tea.Msg) (Model, tea.Cmd, bool) {
	s := &m.search
	switch msg := msg.(type) {
	case searchResultMsg:
		if s.mode != searchLoading {
			return m, nil, true // cancelled meanwhile
		}
		switch {
		case msg.err != nil:
			s.mode, s.err = searchInput, msg.err.Error()
		case len(msg.results) == 0:
			s.mode, s.err = searchInput, "no results"
		default:
			s.mode, s.results, s.err = searchResults, msg.results, ""
		}
		return m, nil, true

	case downloadDoneMsg:
		if msg.err != nil {
			s.status = "download failed: " + msg.err.Error()
		} else {
			s.status = "downloaded: " + msg.title
			m.reload()
		}
		if s.mode == searchDownloading {
			s.mode = searchOff
		}
		return m, nil, true

	case tea.KeyPressMsg:
		s.status = ""
		key := msg.String()

		if s.mode == searchOff {
			if key == "/" {
				s.mode, s.err, s.loaded = searchInput, "", false
				s.refreshLocal(m.root)
				return m, nil, true
			}
			return m, nil, false
		}

		if key == "ctrl+c" {
			return m, tea.Quit, true
		}
		if key == "esc" {
			s.mode = searchOff
			return m, nil, true
		}

		switch s.mode {
		case searchInput:
			switch key {
			case "enter":
				if _, ok := localQuery(s.query); ok {
					if s.localCur < len(s.matches) {
						m.playQueue(s.matches, s.localCur)
						s.mode = searchOff
					}
				} else if q := strings.TrimSpace(s.query); q != "" {
					s.mode, s.err = searchLoading, ""
					return m, runSearch(q), true
				}
			case "backspace":
				if r := []rune(s.query); len(r) > 0 {
					s.query = string(r[:len(r)-1])
				}
			case "ctrl+u":
				s.query = ""
			case "down", "ctrl+n":
				s.localCur = min(s.localCur+1, len(s.matches)-1)
			case "up", "ctrl+p":
				s.localCur = max(s.localCur-1, 0)
			default:
				s.query += msg.Text
			}
			s.refreshLocal(m.root)

		case searchResults:
			switch key {
			case "1", "2", "3", "4", "5":
				i := int(key[0] - '1')
				if i < len(s.results) {
					s.mode, s.picked = searchDownloading, s.results[i].title
					return m, download(s.results[i], m.root), true
				}
			case "/":
				s.mode, s.query, s.results = searchInput, "", nil
			}
		}
		return m, nil, true
	}
	return m, nil, false
}

// searchBody renders the results area that replaces the directory and
// song panes while the search view is open.
func (m Model) searchBody(w int) string {
	s := m.search
	var lines []string
	switch s.mode {
	case searchInput:
		if _, ok := localQuery(s.query); ok {
			return m.localBody(w)
		}
		lines = append(lines, faintStyle.Render(" type a query and press enter"),
			faintStyle.Render(" or /local <text> to search your own songs"))
		if s.err != "" {
			lines = append(lines, "", " "+s.err)
		}
	case searchLoading:
		lines = append(lines, faintStyle.Render(" searching youtube…"))
	case searchDownloading:
		lines = append(lines, " downloading:", " "+s.picked, "",
			faintStyle.Render(" esc to hide — it keeps running"))
	case searchResults:
		for i, r := range s.results {
			lines = append(lines, currentStyle.Render(fmt.Sprintf(" %d ", i+1))+r.title)
		}
		lines = append(lines, "", faintStyle.Render(" 1-5: download · /: new search · esc: close"))
	}
	for i, l := range lines {
		lines[i] = ansi.Truncate(l, w, "…")
	}
	return strings.Join(lines, "\n")
}

// searchPrompt is the vim-style command line shown in the top bar.
func (m Model) searchPrompt() string {
	if m.search.mode == searchInput {
		return "/" + m.search.query + "█"
	}
	return "/" + m.search.query
}

// localQuery reports whether q is a "local [text]" command and returns the text.
func localQuery(q string) (string, bool) {
	if q == "local" {
		return "", true
	}
	if rest, ok := strings.CutPrefix(q, "local "); ok {
		return strings.TrimSpace(rest), true
	}
	return "", false
}

// refreshLocal recomputes the live matches for a "local ..." query.
func (s *searchState) refreshLocal(root string) {
	text, ok := localQuery(s.query)
	if !ok {
		s.matches, s.localCur = nil, 0
		return
	}
	if !s.loaded {
		s.all, s.loaded = library.WalkAudio(root), true
	}
	words := strings.Fields(strings.ToLower(text))
	s.matches = s.matches[:0:0]
	for _, e := range s.all {
		name := strings.ToLower(e.Name)
		hit := true
		for _, w := range words {
			if !strings.Contains(name, w) {
				hit = false
				break
			}
		}
		if hit {
			s.matches = append(s.matches, e)
		}
	}
	s.localCur = min(s.localCur, max(len(s.matches)-1, 0))
}

func (m Model) localBody(w int) string {
	s := m.search
	if len(s.matches) == 0 {
		return faintStyle.Render(" no matching local songs")
	}
	rows := max(m.dirsRows()-2, 1)
	off := max(s.localCur-rows+1, 0)
	items := make([]string, len(s.matches))
	for i, e := range s.matches {
		rel, err := filepath.Rel(m.root, e.Path)
		if err != nil {
			rel = e.Name
		}
		items[i] = rel
	}
	return renderRows(items, s.localCur, off, w, rows, true) + "\n\n" +
		faintStyle.Render(fmt.Sprintf(" %d/%d · ↑/↓: select · enter: play · esc: close",
			s.localCur+1, len(s.matches)))
}
