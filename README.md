# mutux

A terminal music player written in Go. It browses `~/Music` and plays mp3, flac, wav and ogg files.

## Build

Needs Go and, on Linux, the ALSA development libraries (`libasound2-dev` on Debian/Ubuntu, `alsa-lib-devel` on Fedora).

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

Shown from the song's embedded cover, or a `cover.jpg`, `folder.jpg` etc. in the same folder. It needs a truecolor terminal.
