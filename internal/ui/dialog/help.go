package dialog

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const HelpID = "help"

type CommandsDialog struct {
	width  int
	height int
	keyMap struct {
		Close key.Binding
	}
}

func NewCommandsDialog(width, height int) *CommandsDialog {
	h := &CommandsDialog{width: width, height: height}
	h.keyMap.Close = CloseKey
	return h
}

func (*CommandsDialog) ID() string { return HelpID }

func (h *CommandsDialog) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if key.Matches(msg, h.keyMap.Close) {
			return ActionClose{}
		}
	}
	return nil
}

func (h *CommandsDialog) View(width, height int) string {
	content := `───────
/tasks           list all tasks
/newtask <title> create a new task
/help            show this help`

	dialogWidth := min(width-20, 55)
	dialogHeight := 16

	dialogStyle := DialogStyle.
		Width(dialogWidth).
		Height(dialogHeight)

	return dialogStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			TitleStyle.Render("commands"),
			ContentStyle.Render(content),
			HelpStyle.Render("esc to close"),
		),
	)
}

type ActionClose struct{}
