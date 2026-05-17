# ateam

Where humans and AI agents build together

## Overview

A terminal UI application for coordinating multiple AI agents and humans to build software together. Built with Go using the charm.land/bubbletea v2 framework.

## Features

- Multi-line message input (Shift+Enter / Ctrl+J for newline)
- Scrollable message area with vim-style navigation (Ctrl+U/D, arrow keys)
- Floating help dialog (type `/` to open)
- User messages right-aligned, agent messages left-aligned
- Task management (`/tasks`, `/newtask <title>`)

## Quick Start

```bash
# Build
go build .

# Run
./ateam
```

## Controls

| Key | Action |
|-----|--------|
| Enter | Send message |
| Shift+Enter / Ctrl+J | Newline |
| ↑ / ↓ | Navigate messages |
| Ctrl+U / Ctrl+D | Scroll up/down half page |
| / | Open help dialog |
| Esc | Close dialog |
| Ctrl+C | Quit |

## Commands

- `/tasks` - List all tasks
- `/newtask <title>` - Create a new task
- `/help` - Show help

## Architecture

```
ateam/
├── main.go                    # Entry point
├── internal/
│   ├── models/               # Data models
│   ├── store/                # In-memory storage
│   └── ui/
│       ├── app.go           # Main app model
│       └── dialog/          # Dialog system
```

## Tech Stack

- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [bubbles](https://github.com/charmbracelet/bubbles) - TUI components (viewport, textarea)