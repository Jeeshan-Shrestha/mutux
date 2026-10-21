package ui

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// artState holds the decoded cover and the last render, since re-scaling
// the image on every tick would be wasteful.
type artState struct {
	img image.Image

	w, h     int
	rendered string
}

var coverNames = []string{"cover", "folder", "front", "album", "albumart"}

// loadArt prefers the embedded picture, then a cover image beside the file.
// It returns a state with a nil image when nothing usable is found.
func loadArt(audioPath string, embedded []byte) *artState {
	if len(embedded) > 0 {
		if img, _, err := image.Decode(bytes.NewReader(embedded)); err == nil {
			return &artState{img: img}
		}
	}

	dir := filepath.Dir(audioPath)
	for _, name := range coverNames {
		for _, ext := range []string{".jpg", ".jpeg", ".png"} {
			b, err := os.ReadFile(filepath.Join(dir, name+ext))
			if err != nil {
				continue
			}
			if img, _, err := image.Decode(bytes.NewReader(b)); err == nil {
				return &artState{img: img}
			}
		}
	}
	return &artState{}
}

// render draws the cover into a w x h cell area using half-block
// characters: each cell shows two vertically stacked pixels.
func (a *artState) render(w, h int) string {
	if a == nil {
		return ""
	}
	if a.img == nil {
		return faintStyle.Render(" no album art")
	}
	if w < 1 || h < 1 {
		return ""
	}
	if a.w == w && a.h == h {
		return a.rendered
	}

	b := a.img.Bounds()
	sw, sh := b.Dx(), b.Dy()
	if sw == 0 || sh == 0 {
		return ""
	}

	// fit inside w x 2h pixels, keeping aspect ratio
	pw, ph := w, h*2
	nw, nh := pw, sh*pw/sw
	if nh > ph {
		nw, nh = sw*ph/sh, ph
	}
	if nw < 1 {
		nw = 1
	}
	if nh < 1 {
		nh = 1
	}

	// box-average the source pixels that fall into each target pixel
	px := func(x, y int) (r, g, bl uint32) {
		x0, x1 := b.Min.X+x*sw/nw, b.Min.X+(x+1)*sw/nw
		y0, y1 := b.Min.Y+y*sh/nh, b.Min.Y+(y+1)*sh/nh
		if x1 <= x0 {
			x1 = x0 + 1
		}
		if y1 <= y0 {
			y1 = y0 + 1
		}
		var sr, sg, sb, n uint32
		for yy := y0; yy < y1; yy++ {
			for xx := x0; xx < x1; xx++ {
				cr, cg, cb, _ := a.img.At(xx, yy).RGBA()
				sr, sg, sb, n = sr+cr>>8, sg+cg>>8, sb+cb>>8, n+1
			}
		}
		return sr / n, sg / n, sb / n
	}

	pad := strings.Repeat(" ", (w-nw)/2)
	var sb strings.Builder
	rows := (nh + 1) / 2
	for row := 0; row < rows; row++ {
		sb.WriteString(pad)
		for x := 0; x < nw; x++ {
			tr, tg, tb := px(x, row*2)
			fmt.Fprintf(&sb, "\x1b[38;2;%d;%d;%dm", tr, tg, tb)
			if row*2+1 < nh {
				br, bg, bb := px(x, row*2+1)
				fmt.Fprintf(&sb, "\x1b[48;2;%d;%d;%dm", br, bg, bb)
			}
			sb.WriteString("▀")
			sb.WriteString(ansi.ResetStyle)
		}
		if row < rows-1 {
			sb.WriteByte('\n')
		}
	}

	a.w, a.h, a.rendered = w, h, sb.String()
	return a.rendered
}
