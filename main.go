package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"github.com/Jeeshan-Shrestha/mutux.git/internal/ui"
	"os"
)

func main() {
	p := tea.NewProgram(ui.New("/home/nara/Music/"))
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
