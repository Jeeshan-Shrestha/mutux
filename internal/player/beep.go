package player

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

// every track is resampled to this rate so the speaker is only initialised once
const outRate beep.SampleRate = 44100

const volStep = 5 // percent per key press

// Info describes the track that is currently loaded.
type Info struct {
	SampleRate int
	Channels   int
	Length     time.Duration
}

type Player struct {
	mu sync.Mutex // guards everything below except ended/gen

	ready bool // speaker initialised

	src    beep.StreamSeekCloser
	format beep.Format
	ctrl   *beep.Ctrl
	vol    *effects.Volume
	volPct int
	loaded bool

	gen   atomic.Int64 // bumped on every Play/Stop, ignores stale callbacks
	ended atomic.Bool  // set when the current track played to the end
}

func New() *Player { return &Player{volPct: 80} }

// Supported reports whether the file extension can be decoded.
func Supported(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp3", ".flac", ".wav", ".ogg":
		return true
	}
	return false
}

func decode(path string) (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, beep.Format{}, err
	}

	var (
		s    beep.StreamSeekCloser
		fmt_ beep.Format
	)
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp3":
		s, fmt_, err = mp3.Decode(f)
	case ".flac":
		s, fmt_, err = flac.Decode(f)
	case ".wav":
		s, fmt_, err = wav.Decode(f)
	case ".ogg":
		s, fmt_, err = vorbis.Decode(f)
	default:
		f.Close()
		return nil, beep.Format{}, fmt.Errorf("unsupported format %q", filepath.Ext(path))
	}
	if err != nil {
		f.Close()
		return nil, beep.Format{}, err
	}
	return s, fmt_, nil
}

func volumeValue(pct int) (float64, bool) {
	if pct <= 0 {
		return 0, true
	}
	return math.Log2(float64(pct) / 100), false
}

// Play stops whatever is playing and starts path from the beginning.
func (p *Player) Play(path string) (Info, error) {
	src, format, err := decode(path)
	if err != nil {
		return Info{}, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.ready {
		if err := speaker.Init(outRate, outRate.N(time.Second/10)); err != nil {
			src.Close()
			return Info{}, err
		}
		p.ready = true
	}

	p.stopLocked()
	gen := p.gen.Add(1)
	p.ended.Store(false)

	var s beep.Streamer = src
	if format.SampleRate != outRate {
		s = beep.Resample(4, format.SampleRate, outRate, s)
	}

	v, silent := volumeValue(p.volPct)
	ctrl := &beep.Ctrl{Streamer: s}
	vol := &effects.Volume{Streamer: ctrl, Base: 2, Volume: v, Silent: silent}

	p.src, p.format, p.ctrl, p.vol, p.loaded = src, format, ctrl, vol, true

	speaker.Play(beep.Seq(vol, beep.Callback(func() {
		if p.gen.Load() == gen {
			p.ended.Store(true)
		}
	})))

	return Info{
		SampleRate: int(format.SampleRate),
		Channels:   format.NumChannels,
		Length:     format.SampleRate.D(src.Len()),
	}, nil
}

// stopLocked requires p.mu to be held.
func (p *Player) stopLocked() {
	if !p.loaded {
		return
	}
	p.gen.Add(1)
	speaker.Clear()
	speaker.Lock()
	p.src.Close()
	speaker.Unlock()
	p.src, p.ctrl, p.vol, p.loaded = nil, nil, nil, false
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stopLocked()
	p.ended.Store(false)
}

// Toggle flips pause/resume. Returns false if nothing is loaded.
func (p *Player) Toggle() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.loaded {
		return false
	}
	speaker.Lock()
	p.ctrl.Paused = !p.ctrl.Paused
	speaker.Unlock()
	return true
}

func (p *Player) Paused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.loaded {
		return false
	}
	speaker.Lock()
	defer speaker.Unlock()
	return p.ctrl.Paused
}

func (p *Player) Loaded() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.loaded
}

// Ended is true once the current track has played to the end.
func (p *Player) Ended() bool { return p.ended.Load() }

func (p *Player) Position() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.loaded {
		return 0
	}
	speaker.Lock()
	defer speaker.Unlock()
	return p.format.SampleRate.D(p.src.Position())
}

// Seek moves the playhead by delta, clamped to the track.
func (p *Player) Seek(delta time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.loaded {
		return
	}
	speaker.Lock()
	defer speaker.Unlock()

	pos := p.src.Position() + p.format.SampleRate.N(delta)
	if pos < 0 {
		pos = 0
	}
	if max := p.src.Len() - 1; pos > max {
		pos = max
	}
	_ = p.src.Seek(pos)
}

func (p *Player) Volume() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.volPct
}

// AdjustVolume changes the volume by steps of 5%, clamped to 0..100.
func (p *Player) AdjustVolume(steps int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.volPct += steps * volStep
	if p.volPct < 0 {
		p.volPct = 0
	}
	if p.volPct > 100 {
		p.volPct = 100
	}
	if !p.loaded {
		return
	}
	v, silent := volumeValue(p.volPct)
	speaker.Lock()
	p.vol.Volume, p.vol.Silent = v, silent
	speaker.Unlock()
}
