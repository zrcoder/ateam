package models

import (
	"strings"
	"testing"
)

func TestValidatePerson(t *testing.T) {
	err := ValidatePerson("", "test@example.com")
	if err == nil {
		t.Error("expected error for empty name")
	}

	err = ValidatePerson("Test", "")
	if err == nil {
		t.Error("expected error for empty email")
	}

	err = ValidatePerson("Test", "invalid-email")
	if err == nil {
		t.Error("expected error for invalid email")
	}

	err = ValidatePerson("Test", "test@example.com")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateAgent(t *testing.T) {
	err := ValidateAgent("")
	if err == nil {
		t.Error("expected error for empty name")
	}

	err = ValidateAgent("valid-name")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateTask(t *testing.T) {
	err := ValidateTask("")
	if err == nil {
		t.Error("expected error for empty title")
	}

	var longTitle strings.Builder
	for range 201 {
		longTitle.WriteString("a")
	}
	err = ValidateTask(longTitle.String())
	if err == nil {
		t.Error("expected error for title too long")
	}

	err = ValidateTask("valid title")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateMessage(t *testing.T) {
	err := ValidateMessage("")
	if err == nil {
		t.Error("expected error for empty content")
	}

	var longContent strings.Builder
	for range 10001 {
		longContent.WriteString("a")
	}
	err = ValidateMessage(longContent.String())
	if err == nil {
		t.Error("expected error for content too long")
	}

	err = ValidateMessage("valid message")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewPerson(t *testing.T) {
	person, err := NewPerson("Test User", "test@example.com")
	if err != nil {
		t.Fatalf("NewPerson() unexpected error: %v", err)
	}
	if person.Name != "Test User" {
		t.Errorf("NewPerson().Name = %v, want %v", person.Name, "Test User")
	}
	if person.Email != "test@example.com" {
		t.Errorf("NewPerson().Email = %v, want %v", person.Email, "test@example.com")
	}
	if person.ID == "" {
		t.Error("NewPerson().ID should not be empty")
	}
	if person.Status != StatusOnline {
		t.Errorf("NewPerson().Status = %v, want %v", person.Status, StatusOnline)
	}
}

func TestNewPersonValidation(t *testing.T) {
	_, err := NewPerson("", "test@example.com")
	if err == nil {
		t.Error("NewPerson() expected error for empty name")
	}

	_, err = NewPerson("Test", "")
	if err == nil {
		t.Error("NewPerson() expected error for empty email")
	}
}

func TestNewAgent(t *testing.T) {
	agent, err := NewAgent("person-1", "dev-bot", "engineer", "claude", "claude-sonnet")
	if err != nil {
		t.Fatalf("NewAgent() unexpected error: %v", err)
	}
	if agent.Name != "dev-bot" {
		t.Errorf("NewAgent().Name = %v, want %v", agent.Name, "dev-bot")
	}
	if agent.PersonID != "person-1" {
		t.Errorf("NewAgent().PersonID = %v, want %v", agent.PersonID, "person-1")
	}
	if agent.ID == "" {
		t.Error("NewAgent().ID should not be empty")
	}
}

func TestNewAgentValidation(t *testing.T) {
	_, err := NewAgent("", "dev-bot", "engineer", "claude", "model")
	if err == nil {
		t.Error("NewAgent() expected error for empty personID")
	}

	_, err = NewAgent("person-1", "", "engineer", "claude", "model")
	if err == nil {
		t.Error("NewAgent() expected error for empty name")
	}
}

func TestNewComputer(t *testing.T) {
	comp, err := NewComputer("person-1", "My Mac", "localhost")
	if err != nil {
		t.Fatalf("NewComputer() unexpected error: %v", err)
	}
	if comp.Name != "My Mac" {
		t.Errorf("NewComputer().Name = %v, want %v", comp.Name, "My Mac")
	}
	if comp.PersonID != "person-1" {
		t.Errorf("NewComputer().PersonID = %v, want %v", comp.PersonID, "person-1")
	}
}

func TestNewComputerValidation(t *testing.T) {
	_, err := NewComputer("", "My Mac", "localhost")
	if err == nil {
		t.Error("NewComputer() expected error for empty personID")
	}
}

func TestNewChannel(t *testing.T) {
	ch, err := NewChannel("general", "General discussions")
	if err != nil {
		t.Fatalf("NewChannel() unexpected error: %v", err)
	}
	if ch.Name != "general" {
		t.Errorf("NewChannel().Name = %v, want %v", ch.Name, "general")
	}
}

func TestNewChannelValidation(t *testing.T) {
	_, err := NewChannel("", "General discussions")
	if err == nil {
		t.Error("NewChannel() expected error for empty name")
	}
}

func TestNewMessage(t *testing.T) {
	msg, err := NewMessage("channel-1", "author-1", AuthorTypePerson, "Hello")
	if err != nil {
		t.Fatalf("NewMessage() unexpected error: %v", err)
	}
	if msg.Content != "Hello" {
		t.Errorf("NewMessage().Content = %v, want %v", msg.Content, "Hello")
	}
	if msg.ChannelID != "channel-1" {
		t.Errorf("NewMessage().ChannelID = %v, want %v", msg.ChannelID, "channel-1")
	}
}

func TestNewMessageValidation(t *testing.T) {
	_, err := NewMessage("", "author-1", AuthorTypePerson, "Hello")
	if err == nil {
		t.Error("expected error for empty channelID")
	}

	_, err = NewMessage("channel-1", "", AuthorTypePerson, "Hello")
	if err == nil {
		t.Error("expected error for empty authorID")
	}

	_, err = NewMessage("channel-1", "author-1", AuthorTypePerson, "")
	if err == nil {
		t.Error("expected error for empty content")
	}

	_, err = NewMessage("channel-1", "author-1", "invalid", "Hello")
	if err == nil {
		t.Error("expected error for invalid authorType")
	}
}

func TestNewTask(t *testing.T) {
	task, err := NewTask("Test Task", "creator-1")
	if err != nil {
		t.Fatalf("NewTask() unexpected error: %v", err)
	}
	if task.Title != "Test Task" {
		t.Errorf("NewTask().Title = %v, want %v", task.Title, "Test Task")
	}
	if task.Status != TaskStatusTodo {
		t.Errorf("NewTask().Status = %v, want %v", task.Status, TaskStatusTodo)
	}
	if task.Priority != TaskPriorityMedium {
		t.Errorf("NewTask().Priority = %v, want %v", task.Priority, TaskPriorityMedium)
	}
}

func TestNewTaskValidation(t *testing.T) {
	_, err := NewTask("", "creator-1")
	if err == nil {
		t.Error("NewTask() expected error for empty title")
	}
}

func TestValidationError(t *testing.T) {
	err := ValidationError{Field: "name", Message: "cannot be empty"}
	if err.Error() != "name: cannot be empty" {
		t.Errorf("ValidationError.Error() = %v, want %v", err.Error(), "name: cannot be empty")
	}
}
