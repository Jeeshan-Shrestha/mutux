package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"github.com/Jeeshan-Shrestha/mutux.git/internal/ui"
	"os"
	"path/filepath"
)

func musicDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Music"), nil
}

func main() {
	dir, err := musicDir()
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.New(dir))
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
