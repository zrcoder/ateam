package dialog

import (
	"slices"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
)

// Action represents an action taken in a dialog after handling a message.
type Action any

// Dialog is a component that can be displayed on top of the UI.
type Dialog interface {
	ID() string
	HandleMsg(msg tea.Msg) Action
	View(width, height int) string
}

// Overlay manages dialogs as an overlay.
type Overlay struct {
	dialogs []Dialog
}

func NewOverlay() *Overlay {
	return &Overlay{dialogs: []Dialog{}}
}

func (d *Overlay) HasDialogs() bool {
	return len(d.dialogs) > 0
}

func (d *Overlay) OpenDialog(dialog Dialog) {
	d.dialogs = append(d.dialogs, dialog)
}

func (d *Overlay) CloseDialog(id string) {
	for i, dialog := range d.dialogs {
		if dialog.ID() == id {
			d.dialogs = slices.Delete(d.dialogs, i, i+1)
			return
		}
	}
}

func (d *Overlay) Update(msg tea.Msg) tea.Msg {
	if len(d.dialogs) == 0 {
		return nil
	}
	return d.dialogs[len(d.dialogs)-1].HandleMsg(msg)
}

func (d *Overlay) View(width, height int) string {
	if len(d.dialogs) == 0 {
		return ""
	}
	return d.dialogs[len(d.dialogs)-1].View(width, height)
}

// DrawCenter draws the given string view centered in the screen area.
func DrawCenter(scr uv.Screen, area uv.Rectangle, view string) {
	width, height := lipgloss.Size(view)
	x := area.Min.X + (area.Dx()-width)/2
	y := area.Min.Y + (area.Dy()-height)/2
	r := uv.Rect(x, y, width, height)
	uv.NewStyledString(view).Draw(scr, r)
}
