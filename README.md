# mutux

A terminal music player written in Go. It browses `~/Music` and plays mp3, flac, wav and ogg files.

## Dependencies

- Go 1.26.1 (an older Go 1.21+ will download it automatically)
- Linux only: a C compiler (`gcc`), `pkg-config` and the ALSA development headers (audio output uses cgo)
- A truecolor terminal, for album art

Go libraries (fetched automatically by `go build`): Bubble Tea, Lip Gloss, beep (audio playback) and dhowden/tag (metadata).

## Install

```sh
sudo dnf copr enable nara-nara/mutux
sudo dnf install mutux
mutux
```
for fedora

```sh
git clone https://github.com/Jeeshan-Shrestha/mutux.git mutux
cd mutux
./install.sh
```

The script installs the Linux dependencies above (apt, dnf, pacman or zypper, using sudo), builds the app and copies it to `~/.local/bin/mutux`. Make sure that folder is on your `PATH`, then run:

```sh
mutux
```

To build by hand instead:

```sh
go build -o mutux
./mutux
```

## Keys

| Key | Action |
|---|---|
| `j` / `k` | Move down / up |
| `g` / `G` | Jump to top / bottom |
| `l` / `h` | Open folder / go back (Directories panel) |
| `tab` / `shift+tab` or `ctrl+l` / `ctrl+h` | Switch panel |
| `enter` | Play song |
| `space` | Pause / resume |
| `n` / `p` | Next / previous song |
| `<` / `>` | Seek back / forward 5s |
| `q` | Quit |

## Album art

Shown from the song's embedded cover, or a `cover.jpg`, `folder.jpg` etc. in the same folder.
