# Ateam - AI Agents & Humans Collaborative Building Platform

## 1. Project Overview

- **Project Name**: ateam
- **Type**: Terminal UI (TUI) Application
- **Core Functionality**: A collaborative platform for coordinating multiple AI agents and real humans to build software together, featuring channels (like Slack), people, agents, computers, and tasks.
- **Target Users**: Software development teams wanting to integrate AI agents into their workflow

## 2. Technology Stack

- **Language**: Go 1.21+
- **TUI Framework**: charm.land/bubbletea/v2
- **Styling**: charm.land/lipgloss/v2
- **Viewport**: charm.land/bubbles/v2/viewport
- **Textarea**: charm.land/bubbles/v2/textarea
- **Architecture**: Modular with tea.Model pattern

## 3. UI/UX Specification

### 3.1 Window Structure

```
┌─────────────────────────────────────────────────────────────────┐
│ ATEAM                                        tasks: 2          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│    * You - 10:30                                                  │
│        Welcome to ateam!                                          │
│                                                                  │
│    A dev-bot - 10:31                                              │
│        Hello! I'm your AI assistant.                              │
│                                                                  │
├─────────────────────────────────────────────────────────────────┤
│ > Type message...                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Layout** (top to bottom):
1. Title bar - "ATEAM" and task count
2. Message area - scrollable messages with viewport
3. Input area - multi-line textarea for typing messages/commands

### 3.2 Visual Design

- Uses Dracula color scheme (via lipgloss ANSI colors)
- User messages: right-aligned with green marker (*)
- Agent messages: left-aligned with pink marker (A)
- System messages: left-aligned
- Terminal-native styling (respects user's terminal settings)

### 3.3 Components

**Title Bar**:
- "ATEAM" branding (cyan, bold)
- Task count indicator (muted)

**Message Area**:
- Author marker (* for person, A for agent)
- Author name (cyan, bold)
- Timestamp (muted)
- Message content
- Scrollable via viewport (arrow keys, Page Up/Down, Ctrl+U/D)
- Initial scroll position: bottom (latest messages)

**Input Area**:
- Multi-line textarea
- `Shift+Enter` or `Ctrl+J` for newline
- `Enter` to send message
- `/` at start of empty input triggers help dialog

**Help Dialog**:
- Floating overlay centered on screen
- Shows available commands and keyboard shortcuts
- `Esc` to close

### 3.4 Commands

| Command | Description |
|---------|-------------|
| `/tasks` | List all tasks |
| `/newtask <title>` | Create a new task |
| `/help` | Show help dialog |

### 3.5 Keyboard Navigation

- `↑/↓` or `j/k` - Navigate messages (vim-style)
- `PageUp/PageDown` - Scroll by page
- `Ctrl+U` - Scroll up half page (vim-style)
- `Ctrl+D` - Scroll down half page (vim-style)
- `Shift+Enter` - Newline in input
- `Ctrl+J` - Newline (vim-style)
- `Enter` - Send message
- `/` - Trigger help dialog (when textarea empty)
- `Esc` - Close dialog
- `Ctrl+C` - Quit application

## 4. Data Models

### 4.1 Core Entities

```go
type Person struct {
    ID        string
    Name      string
    Email     string
    Status    string // "online", "away", "offline"
    Computers []Computer
    AgentIDs  []string
}

type Agent struct {
    ID        string
    PersonID  string // Owner
    Name      string
    Role      string // "product-manager", "architect", "engineer", "reviewer"
    AIType    string // "opencode", "claude", "gpt", etc
    Model     string // "claude-sonnet-4-7", etc
    Status    string
}

type Computer struct {
    ID       string
    PersonID string
    Name     string
    Hostname string
    Status   string
}

type Channel struct {
    ID      string
    Name    string
    Purpose string
}

type Message struct {
    ID         string
    ChannelID  string
    AuthorID   string
    AuthorType string // "person", "agent", "system"
    Content    string
    Timestamp  time.Time
}

type Task struct {
    ID          string
    Title       string
    Description string
    Status      string // "todo", "in-progress", "done"
    Priority    string // "low", "medium", "high"
    CreatedBy   string
    AssigneeID  string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## 5. Architecture

```
ateam/
├── main.go                    # Entry point
├── go.mod                    # Module definition
├── go.sum                    # Dependencies
├── internal/
│   ├── models/
│   │   └── models.go         # Data models
│   ├── store/
│   │   └── store.go         # In-memory storage with seed data
│   └── ui/
│       ├── app.go           # Main app model (tea.Model)
│       └── dialog/
│           ├── dialog.go    # Dialog interface and overlay
│           └── help.go      # Help dialog implementation
```

## 6. Features

### 6.1 Implemented

- [x] Application launches and displays TUI
- [x] Title bar with "ATEAM" branding and task count
- [x] User can type and send messages (multi-line support)
- [x] Messages appear in message area
- [x] Task count shown in title bar
- [x] `/tasks` command lists tasks
- [x] `/newtask <title>` creates a task
- [x] Help dialog via `/help` command or typing `/`
- [x] Scrollable message area with viewport
- [x] Initial view shows latest messages (scrolled to bottom)
- [x] User messages right-aligned, agent messages left-aligned
- [x] Ctrl+C quits application
- [x] Esc closes dialog
- [x] Dialog floats on top of content using lipgloss Layer/Compositor

### 6.2 Prototype Scope (Single Person)

MVP features:
- Single hardcoded Person (the user)
- One default channel (#general)
- Message sending and viewing
- Task creation and listing
- Agent listing
- In-memory storage with mock data for testing

## 7. Dialog System

Help dialog uses lipgloss Layer/Compositor for floating overlay:

```go
dialogLayer := lipgloss.NewLayer(dialogContent).
    X((width-55)/2).
    Y((height-16)/2).
    Z(1)

content = lipgloss.NewCompositor(
    lipgloss.NewLayer(mainContent),
    dialogLayer,
).Render()
```

## 8. Mock Data

Store includes 50+ mock messages for testing scroll functionality:
- Simulated conversation between user and agents
- Numbered test messages (00-49) for scroll validation