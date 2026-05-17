package main

import (
	"log"

	"github.com/zrcoder/ateam/internal/store"
	"github.com/zrcoder/ateam/internal/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	s, err := store.NewDefault()
	if err != nil {
		log.Fatal("Failed to create store:", err)
	}
	defer s.Close()

	p := tea.NewProgram(ui.NewModel(s), tea.WithoutSignalHandler())
	if _, err := p.Run(); err != nil {
		log.Fatal("Failed to run program:", err)
	}
}
