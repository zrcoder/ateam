package main

import (
	"os"

	"github.com/zrcoder/ateam/internal/store"
	"github.com/zrcoder/ateam/internal/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	dataDir, err := getDataDir()
	if err != nil {
		dataDir = "." // fallback to current directory
	}

	s, err := store.New(dataDir)
	if err != nil {
		os.Exit(1)
	}
	defer s.Close()

	if err := s.Seed(); err != nil {
		os.Exit(1)
	}

	p := tea.NewProgram(ui.NewModel(s), tea.WithoutSignalHandler())

	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}

func getDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dataDir := home + "/.ateam"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}
	return dataDir, nil
}
