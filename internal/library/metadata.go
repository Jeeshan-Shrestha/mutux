package library

import (
	"os"
	"strings"

	"github.com/dhowden/tag"
)

// Meta is the tag information shown in the Now Playing panel.
type Meta struct {
	Title, Artist, Album, AlbumArtist, Genre, Composer string
	Year, Track, TrackTotal                            int
	Format                                             string
	Art                                                []byte // embedded cover image, nil if none
}

// ReadMeta reads tags from path. A file with no tags yields an empty Meta
// and no error, so callers can always fall back to the file name.
func ReadMeta(path string) Meta {
	f, err := os.Open(path)
	if err != nil {
		return Meta{}
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		return Meta{}
	}

	var art []byte
	if pic := m.Picture(); pic != nil {
		art = pic.Data
	}

	track, total := m.Track()
	return Meta{
		Title:       strings.TrimSpace(m.Title()),
		Artist:      strings.TrimSpace(m.Artist()),
		Album:       strings.TrimSpace(m.Album()),
		AlbumArtist: strings.TrimSpace(m.AlbumArtist()),
		Genre:       strings.TrimSpace(m.Genre()),
		Composer:    strings.TrimSpace(m.Composer()),
		Year:        m.Year(),
		Track:       track,
		TrackTotal:  total,
		Format:      string(m.FileType()),
		Art:         art,
	}
}
