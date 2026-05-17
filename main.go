package main

import (
	"os"

	"github.com/zrcoder/ateam/internal/ui"

	"github.com/zrcoder/ateam/internal/store"

	tea "charm.land/bubbletea/v2"
)

func main() {
	s := store.New()
	p := tea.NewProgram(ui.NewModel(s), tea.WithoutSignalHandler())

	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
