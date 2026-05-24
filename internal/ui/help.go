package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
	"github.com/zrcoder/ateam/internal/ui/dialog"
)

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
		{km.ScrollUp, km.ScrollDown},
		{km.Newline},
	}
}

type HelpDialog struct {
	*dialog.Base
	keyHelp help.Model
	keyMap  KeyMap
}

func NewHelpDialog(width, height int) *HelpDialog {
	h := &HelpDialog{}
	h.keyMap = DefaultKeyMap
	h.keyHelp = help.New()
	h.keyHelp.ShowAll = true
	h.Base = dialog.NewBase("Help", width, height)
	return h
}

func (h *HelpDialog) View() string {

	content := lipgloss.JoinVertical(lipgloss.Left,

		dialog.ContentStyle.Render(`
/tasks           list all tasks
/newtask <title> create a new task
/help            show this help`),
		"",
		h.Divider(),
		"Key bindings",
		h.keyHelp.View(h.keyMap),
	)
	return h.WrapView(content)
}
