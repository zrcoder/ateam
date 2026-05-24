package ui

import (
	"github.com/zrcoder/ateam/internal/agent/runtime"
	"github.com/zrcoder/ateam/internal/ui/dialog"
)

type AgentDialog struct {
	*dialog.Base
	providers []runtime.Provider
}

func NewAgentDialog(width, height int) *AgentDialog {
	a := &AgentDialog{}
	a.Base = dialog.NewBase("Agent", width, height)
	a.providers = runtime.DiscoverProviders()
	return a
}

// View implements Dialog.View
func (a *AgentDialog) View(width, height int) string {
	return a.WrapView("")
}
