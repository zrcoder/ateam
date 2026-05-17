package main

import (
	"os"

	"ateam/internal/store"
	"ateam/internal/ui"

	"charm.land/bubbletea/v2"
)

func main() {
	s := store.New()
	p := tea.NewProgram(ui.NewModel(s), tea.WithoutSignalHandler())

	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}