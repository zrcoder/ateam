package dialog

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

const HelpID = "help"

type HelpDialog struct {
	width  int
	height int
	keyMap struct {
		Close key.Binding
	}
}

func NewHelpDialog(width, height int) *HelpDialog {
	h := &HelpDialog{width: width, height: height}
	h.keyMap.Close = CloseKey
	return h
}

func (*HelpDialog) ID() string { return HelpID }

func (h *HelpDialog) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if key.Matches(msg, h.keyMap.Close) {
			return ActionClose{}
		}
	}
	return nil
}

func (h *HelpDialog) View(width, height int) string {
	content := `commands
───────
/tasks           list all tasks
/newtask <title> create a new task
/help            show this help

keys
────
shift+enter    newline
enter          send message
ctrl+j         newline (vim-style)
esc            close help
ctrl+c         quit`

	dialogWidth := min(width-20, 55)
	dialogHeight := 16

	dialogStyle := DialogStyle.
		Width(dialogWidth).
		Height(dialogHeight)

	return dialogStyle.Render(
		TitleStyle.Render("help") + "\n" +
			ContentStyle.Render(content) + "\n" +
			HelpStyle.Render("esc to close"),
	)
}

type ActionClose struct{}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}