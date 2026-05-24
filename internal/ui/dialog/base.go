package dialog

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Dialog styles
var (
	DialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6272A4")).
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD")).
			Bold(true)

	ContentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))
)

const (
	Sep = "───────"
)

// Base provides common dialog functionality to be embedded in dialog implementations.
type Base struct {
	title     string
	width     int
	height    int
	dialogW   int
	dialogH   int
	closeKey  key.Binding
	helpStyle lipgloss.Style
}

func NewBase(title string, width, height int) Base {
	return Base{
		title:   title,
		width:   width,
		height:  height,
		dialogW: min(width-20, 55),
		dialogH: 16,
		closeKey: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "close"),
		),
		helpStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4")),
	}
}

// HandleMsg handles common keys like ESC to close.
func (b *Base) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, b.closeKey) {
			return ActionClose{}
		}
	}
	return nil
}

func (b *Base) DialogStyle() lipgloss.Style {
	return DialogStyle.Width(b.dialogW).Height(b.dialogH)
}

func (b *Base) Content(content string) string {
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, TitleStyle.Render(b.title), b.helpStyle.MarginLeft(30).Render("[x] esc")),
		Sep,
		ContentStyle.Render(content),
	)
}
