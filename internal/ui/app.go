package ui

import (
	"fmt"
	"strings"

	"github.com/zrcoder/ateam/internal/ui/dialog"

	"github.com/zrcoder/ateam/internal/store"

	"github.com/zrcoder/ateam/internal/models"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	store           *store.Store
	width           int
	height          int
	currentChannel  string
	messages        []MessageRow
	selectedMessage int
	textarea        textarea.Model
	dialogOverlay   *dialog.Overlay
	dialogVisible   bool
	viewport        viewport.Model
	viewportReady   bool
}

type MessageRow struct {
	author     string
	authorType string
	content    string
	time       string
}

func NewModel(s *store.Store) *Model {
	ta := textarea.New()
	ta.Placeholder = "Type message..."
	ta.Focus()
	ta.SetHeight(3)

	ta.KeyMap.InsertNewline = key.NewBinding(
		key.WithKeys("shift+enter"),
		key.WithHelp("shift+enter", "newline"),
	)

	return &Model{
		store:           s,
		currentChannel:  "channel-general",
		messages:        []MessageRow{},
		selectedMessage: 0,
		textarea:        ta,
		dialogOverlay:   dialog.NewOverlay(),
		dialogVisible:   false,
	}
}

func (m *Model) loadMessages() {
	messages, err := m.store.GetMessages(m.currentChannel)
	if err != nil {
		return
	}
	m.messages = make([]MessageRow, 0, len(messages))
	for _, msg := range messages {
		author := "You"
		if msg.AuthorType == models.AuthorTypeAgent {
			agent, _ := m.store.GetAgent(msg.AuthorID)
			if agent != nil {
				author = agent.Name
			}
		}
		m.messages = append(m.messages, MessageRow{
			author:     author,
			authorType: msg.AuthorType,
			content:    msg.Content,
			time:       msg.Timestamp.Format("15:04"),
		})
	}
	if m.viewportReady {
		m.viewport.SetContent(m.buildMessagesContent())
	}
}

func (m *Model) buildMessagesContent() string {
	var b strings.Builder
	for _, msg := range m.messages {
		isUser := msg.authorType == models.AuthorTypePerson

		marker := "*"
		markerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
		if msg.authorType == models.AuthorTypeAgent {
			marker = "A"
			markerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6")).Bold(true)
		}

		authorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD")).
			Bold(true)

		timeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))

		contentStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))

		headerPrefix := "  "
		contentPrefix := "    "

		if isUser {
			headerPrefix = strings.Repeat(" ", max(0, m.width-30))
			contentPrefix = strings.Repeat(" ", max(0, m.width-20))
		}

		b.WriteString(headerPrefix)
		b.WriteString(markerStyle.Render(marker))
		b.WriteString(" ")
		b.WriteString(authorStyle.Render(msg.author))
		b.WriteString(" ")
		b.WriteString(timeStyle.Render(msg.time))
		b.WriteString("\n")

		for line := range strings.SplitSeq(msg.content, "\n") {
			b.WriteString(contentPrefix)
			b.WriteString(contentStyle.Render(line))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m *Model) Init() tea.Cmd {
	m.loadMessages()
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(m.width)

		m.viewport = viewport.New(viewport.WithWidth(m.width), viewport.WithHeight(m.height-9))
		m.viewport.YPosition = 1
		m.viewportReady = true
		m.viewport.SetContent(m.buildMessagesContent())
		m.viewport.GotoBottom()

		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.Key().Code == tea.KeyEnter && (msg.Key().Mod&tea.ModShift) == 0 {
			if !m.dialogVisible {
				m.handleInput()
			}
			return m, nil
		}

		if msg.Key().Code == tea.KeyEsc {
			if m.dialogVisible {
				m.dialogOverlay.CloseDialog(dialog.HelpID)
				m.dialogVisible = false
				return m, nil
			}
		}

		if msg.Key().Text == "/" && m.textarea.Value() == "" {
			m.dialogOverlay.OpenDialog(dialog.NewCommandsDialog(m.width, m.height))
			m.dialogVisible = true
			return m, nil
		}
	}

	if m.dialogVisible {
		if a := m.dialogOverlay.Update(msg); a != nil {
			if _, ok := a.(dialog.ActionClose); ok {
				m.dialogOverlay.CloseDialog(dialog.HelpID)
				m.dialogVisible = false
			}
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m *Model) handleInput() {
	text := strings.TrimSpace(m.textarea.Value())
	if text == "" {
		return
	}

	if strings.HasPrefix(text, "/") {
		m.handleCommand(text)
	} else {
		m.sendMessage(text)
	}
	m.textarea.Reset()
}

func (m *Model) handleCommand(text string) {
	parts := strings.SplitN(text, " ", 2)
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "/tasks":
		m.showTasks()
	case "/newtask":
		if len(parts) > 1 {
			title := strings.TrimSpace(parts[1])
			if title != "" {
				task, err := models.NewTask(title, m.store.GetCurrentPerson().ID)
				if err != nil {
					m.showError("Failed to create task: " + err.Error())
					return
				}
				m.store.AddTask(task)
				m.loadMessages()
			} else {
				m.showError("/newtask requires a title")
			}
		} else {
			m.showError("/newtask requires a title, usage: /newtask <title>")
		}
	case "/help":
		m.dialogOverlay.OpenDialog(dialog.NewCommandsDialog(m.width, m.height))
		m.dialogVisible = true
	}
}

func (m *Model) showError(msg string) {
	errMsg, err := models.NewMessage(m.currentChannel, "system", models.AuthorTypeSystem, "Error: "+msg)
	if err != nil {
		return
	}
	m.store.AddMessage(errMsg)
	m.loadMessages()
	if m.viewportReady {
		m.viewport.GotoBottom()
	}
}

func (m *Model) sendMessage(text string) {
	msg, err := models.NewMessage(
		m.currentChannel,
		m.store.GetCurrentPerson().ID,
		models.AuthorTypePerson,
		text,
	)
	if err != nil {
		m.showError("Failed to send message: " + err.Error())
		return
	}
	m.store.AddMessage(msg)
	m.loadMessages()
	if m.viewportReady {
		m.viewport.GotoBottom()
	}
}

func (m *Model) showTasks() {
	tasks, _ := m.store.GetTasks()
	var b strings.Builder
	b.WriteString("\n  tasks\n")
	b.WriteString("  ─────\n")
	if len(tasks) == 0 {
		b.WriteString("  (no tasks)\n")
	}
	for _, t := range tasks {
		var status string
		switch t.Status {
		case models.TaskStatusInProgress:
			status = "[~]"
		case models.TaskStatusDone:
			status = "[x]"
		default:
			status = "[ ]"
		}
		fmt.Fprintf(&b, "  %s %s\n", status, t.Title)
	}
	msg, err := models.NewMessage(m.currentChannel, "system", models.AuthorTypeSystem, b.String())
	if err != nil {
		m.showError("Failed to show tasks: " + err.Error())
		return
	}
	m.store.AddMessage(msg)
	m.loadMessages()
	if m.viewportReady {
		m.viewport.GotoBottom()
	}
}

func (m Model) View() tea.View {
	title := m.renderTitle()
	border := m.renderBorder()
	messages := m.viewport.View()
	inputContent := m.renderInput()
	statusBar := m.renderStatusBar()

	mainContent := lipgloss.JoinVertical(lipgloss.Left,
		title,
		border,
		messages,
		border,
		inputContent,
		border,
		statusBar,
	)

	var content string
	if m.dialogVisible {
		dialogContent := m.dialogOverlay.View(m.width, m.height)
		dialogLayer := lipgloss.NewLayer(dialogContent).
			X((m.width - 55) / 2).
			Y((m.height - 16) / 2).
			Z(1)

		content = lipgloss.NewCompositor(
			lipgloss.NewLayer(mainContent),
			dialogLayer,
		).Render()
	} else {
		content = mainContent
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m Model) renderTitle() string {
	tasks, _ := m.store.GetTasks()
	taskCount := len(tasks)
	left := " ATEAM "
	right := fmt.Sprintf(" tasks: %d ", taskCount)

	leftLen := lipgloss.Width(left)
	rightLen := lipgloss.Width(right)
	available := max(m.width-leftLen-rightLen, 0)
	spacer := strings.Repeat(" ", available)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8BE9FD")).
		Bold(true)

	rightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4"))

	return titleStyle.Render(left) + spacer + rightStyle.Render(right)
}

func (m Model) renderBorder() string {
	border := strings.Repeat("─", m.width)
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#44475A"))
	return borderStyle.Render(border)
}

func (m Model) renderInput() string {
	return m.textarea.View()
}

func (m Model) renderStatusBar() string {
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272A4"))
	return statusStyle.Render(" / commands  ·  shift+enter newline  ·  ctrl+u scroll up  ·  ctrl+d scroll down  ·  ctrl+c quit")
}
