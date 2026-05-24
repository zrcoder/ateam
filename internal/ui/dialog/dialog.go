package dialog

import (
	tea "charm.land/bubbletea/v2"
)

// Action represents an action taken in a dialog after handling a message.
type Action any

// Dialog is a component that can be displayed on top of the UI.
type Dialog interface {
	HandleMsg(msg tea.Msg) Action
	View() string
}

// Overlay manages dialogs as an overlay.
type Overlay struct {
	dialog Dialog
}

func NewOverlay() *Overlay {
	return &Overlay{}
}

func (d *Overlay) OpenDialog(dialog Dialog) {
	d.dialog = dialog
}

func (d *Overlay) CloseDialog() {
	d.dialog = nil
}

func (d *Overlay) Update(msg tea.Msg) tea.Msg {
	if d.dialog == nil {
		return nil
	}
	return d.dialog.HandleMsg(msg)
}

func (d *Overlay) View(width, height int) string {
	if d.dialog == nil {
		return ""
	}
	return d.dialog.View()
}
