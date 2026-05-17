package ui

import (
	"ateam/internal/models"
	"ateam/internal/store"
	"testing"
)

func TestNewModel(t *testing.T) {
	s := store.New()
	m := NewModel(s)

	if m == nil {
		t.Fatal("expected non-nil model")
	}
	if m.store != s {
		t.Error("expected store to be set")
	}
	if m.currentChannel != "channel-general" {
		t.Errorf("expected 'channel-general', got %s", m.currentChannel)
	}
	if m.dialogVisible {
		t.Error("expected dialog to be hidden initially")
	}
}

func TestModel_handleCommand(t *testing.T) {
	s := store.New()
	m := NewModel(s)

	m.handleCommand("/tasks")

	if len(m.messages) == 0 {
		t.Error("expected messages to be loaded after /tasks command")
	}
}

func TestModel_handleCommand_EmptyNewtask(t *testing.T) {
	s := store.New()
	m := NewModel(s)

	initialTasks := len(s.GetTasks())

	m.handleCommand("/newtask")

	afterTasks := len(s.GetTasks())
	if afterTasks != initialTasks {
		t.Error("expected no new task for empty /newtask command")
	}
}

func TestModel_handleCommand_Help(t *testing.T) {
	s := store.New()
	m := NewModel(s)
	m.width = 80
	m.height = 24

	m.handleCommand("/help")

	if !m.dialogVisible {
		t.Error("expected dialog to be visible after /help")
	}
}

func TestModel_sendMessage(t *testing.T) {
	s := store.New()
	m := NewModel(s)

	initialCount := len(s.GetMessages(m.currentChannel))

	m.sendMessage("Test message")

	afterCount := len(s.GetMessages(m.currentChannel))
	if afterCount <= initialCount {
		t.Error("expected message to be added")
	}
}

func TestModel_showError(t *testing.T) {
	s := store.New()
	m := NewModel(s)
	m.viewportReady = true

	initialCount := len(s.GetMessages(m.currentChannel))

	m.showError("Test error")

	afterCount := len(s.GetMessages(m.currentChannel))
	if afterCount <= initialCount {
		t.Error("expected error message to be added")
	}
}

func TestModel_handleInput_Empty(t *testing.T) {
	s := store.New()
	m := NewModel(s)

	m.textarea.Reset()
	m.handleInput()
}

func TestModel_handleInput_Whitespace(t *testing.T) {
	s := store.New()
	m := NewModel(s)

	m.textarea.SetValue("   ")
	m.handleInput()
}

func TestMessageRow(t *testing.T) {
	row := MessageRow{
		author:     "Test Author",
		authorType: models.AuthorTypePerson,
		content:    "Test content",
		time:       "12:00",
	}

	if row.author != "Test Author" {
		t.Errorf("expected 'Test Author', got %s", row.author)
	}
	if row.authorType != models.AuthorTypePerson {
		t.Errorf("expected 'person', got %s", row.authorType)
	}
}

func TestNewModel_WithNilStore(t *testing.T) {
	m := NewModel(nil)
	if m == nil {
		t.Error("expected non-nil model even with nil store")
	}
}