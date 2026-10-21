#!/bin/sh
# Installs mutux's system dependencies (Linux only), builds it, and copies
# the binary to ~/.local/bin.
set -e

BIN_DIR="${BIN_DIR:-$HOME/.local/bin}"

have() { command -v "$1" >/dev/null 2>&1; }

if ! have go; then
	echo "error: Go is not installed. Get it from https://go.dev/dl/ and re-run." >&2
	exit 1
fi

if [ "$(uname -s)" = "Linux" ]; then
	# need a C compiler, pkg-config and the ALSA headers (audio output uses cgo)
	if ! have gcc || ! have pkg-config || ! pkg-config --exists alsa; then
		echo "installing build dependencies (needs sudo)..."
		if have apt-get; then
			sudo apt-get update && sudo apt-get install -y gcc pkg-config libasound2-dev
		elif have dnf; then
			sudo dnf install -y gcc pkgconf-pkg-config alsa-lib-devel
		elif have pacman; then
			sudo pacman -S --needed --noconfirm gcc pkgconf alsa-lib
		elif have zypper; then
			sudo zypper install -y gcc pkg-config alsa-devel
		else
			echo "error: unknown package manager; install gcc, pkg-config and the ALSA dev headers manually." >&2
			exit 1
		fi
	fi
fi

cd "$(dirname "$0")"
go build -o mutux

mkdir -p "$BIN_DIR"
install -m 755 mutux "$BIN_DIR/mutux"
rm -f mutux

echo "installed to $BIN_DIR/mutux"
case ":$PATH:" in
	*":$BIN_DIR:"*) ;;
	*) echo "note: $BIN_DIR is not on your PATH" ;;
esac
