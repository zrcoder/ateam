package dialog

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
)

const HelpID = "help"

type KeyMap struct {
	Newline    key.Binding
	ScrollUp   key.Binding
	ScrollDown key.Binding
}

var DefaultKeyMap = KeyMap{
	Newline: key.NewBinding(
		key.WithKeys("shift+enter"),
		key.WithHelp("shift+enter", "newline"),
	),
	ScrollUp: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "scroll up"),
	),
	ScrollDown: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "scroll down"),
	),
}

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Newline, km.ScrollUp, km.ScrollDown}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Newline},
		{km.ScrollUp, km.ScrollDown},
	}
}

type HelpDialog struct {
	Base
	keyHelp help.Model
	keyMap  KeyMap
}

func NewHelpDialog(width, height int) *HelpDialog {
	h := &HelpDialog{}
	h.keyMap = DefaultKeyMap
	h.keyHelp = help.New()
	h.keyHelp.ShowAll = true
	h.Base = NewBase("Help", width, height)
	return h
}

func (*HelpDialog) ID() string { return HelpID }

func (h *HelpDialog) View(width, height int) string {

	content := lipgloss.JoinVertical(lipgloss.Left,

		ContentStyle.Render(`
/tasks           list all tasks
/newtask <title> create a new task
/help            show this help`),
		"",
		TitleStyle.Render("Key bindings"),
		h.keyHelp.View(h.keyMap),
	)
	return h.Base.Content(content)
}

type ActionClose struct{}
