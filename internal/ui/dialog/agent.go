package dialog

import (
	"github.com/zrcoder/ateam/internal/agent/runtime"
)

const AgentID = "agent"

type AgentDialog struct {
	Base
}

func NewAgentDialog(width, height int) *AgentDialog {
	a := &AgentDialog{}
	a.width = width
	a.height = height
	a.dialogW = min(width-20, 55)
	a.dialogH = 16
	_ = runtime.DiscoverProviders()
	return a
}

func (*AgentDialog) ID() string { return AgentID }

// Content implements Base.Content
func (a *AgentDialog) Content() string {
	return ""
}