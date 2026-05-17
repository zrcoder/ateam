package store

import (
	"testing"

	"github.com/zrcoder/ateam/internal/models"
)

func TestStore_New(t *testing.T) {
	tmpDir := t.TempDir()
	s, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer s.Close()

	if s == nil {
		t.Fatal("expected non-nil store")
	}

	// Test that we can query the database
	person := s.GetCurrentPerson()
	if person == nil {
		t.Error("expected to get current person after seed")
	}
}

func TestStore_seed(t *testing.T) {
	tmpDir := t.TempDir()
	s, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer s.Close()

	// seed is already called in New(), calling again should be idempotent
	if err := s.seed(); err != nil {
		t.Fatalf("seed() error = %v", err)
	}

	// Check that seed data exists
	person := s.GetCurrentPerson()
	if person == nil {
		t.Fatal("expected current person")
	}
	if person.Name != "You" {
		t.Errorf("person.Name = %v, want 'You'", person.Name)
	}

	tasks, err := s.GetTasks()
	if err != nil {
		t.Fatalf("GetTasks() error = %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}

	agents, err := s.GetAgents()
	if err != nil {
		t.Fatalf("GetAgents() error = %v", err)
	}
	if len(agents) != 2 {
		t.Errorf("expected 2 agents, got %d", len(agents))
	}
}

func TestStore_AddMessage(t *testing.T) {
	tmpDir := t.TempDir()
	s, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer s.Close()

	s.seed()

	msg := &models.Message{
		ID:         "msg-1",
		ChannelID:  "channel-general",
		AuthorID:   "user-1",
		AuthorType: models.AuthorTypePerson,
		Content:    "Hello",
	}

	if err := s.AddMessage(msg); err != nil {
		t.Fatalf("AddMessage() error = %v", err)
	}

	messages, _, err := s.GetMessages("channel-general", 20, 0)
	if err != nil {
		t.Fatalf("GetMessages() error = %v", err)
	}
	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}
}

func TestStore_AddTask(t *testing.T) {
	tmpDir := t.TempDir()
	s, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer s.Close()

	s.seed()

	task := &models.Task{
		ID:        "task-new",
		Title:     "New Task",
		CreatedBy: "user-1",
	}

	if err := s.AddTask(task); err != nil {
		t.Fatalf("AddTask() error = %v", err)
	}

	tasks, err := s.GetTasks()
	if err != nil {
		t.Fatalf("GetTasks() error = %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestStore_GetAgent(t *testing.T) {
	tmpDir := t.TempDir()
	s, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer s.Close()

	s.seed()

	agents, err := s.GetAgents()
	if err != nil {
		t.Fatalf("GetAgents() error = %v", err)
	}
	if len(agents) == 0 {
		t.Fatal("expected agents")
	}

	agent := agents[0]
	result, err := s.GetAgent(agent.ID)
	if err != nil {
		t.Fatalf("GetAgent() error = %v", err)
	}
	if result == nil {
		t.Error("expected to find agent")
	}
	if result.Name != agent.Name {
		t.Errorf("result.Name = %v, want %v", result.Name, agent.Name)
	}

	// Test non-existent agent
	notFound, err := s.GetAgent("non-existent")
	if err != nil {
		t.Fatalf("GetAgent() error = %v", err)
	}
	if notFound != nil {
		t.Error("expected nil for non-existent agent")
	}
}

func TestStore_Channels(t *testing.T) {
	tmpDir := t.TempDir()
	s, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer s.Close()

	s.seed()

	channels, err := s.GetChannels()
	if err != nil {
		t.Fatalf("GetChannels() error = %v", err)
	}
	if len(channels) != 1 {
		t.Errorf("expected 1 channel, got %d", len(channels))
	}
}

func TestStore_Persons(t *testing.T) {
	tmpDir := t.TempDir()
	s, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer s.Close()

	s.seed()

	persons, err := s.GetPersons()
	if err != nil {
		t.Fatalf("GetPersons() error = %v", err)
	}
	if len(persons) != 1 {
		t.Errorf("expected 1 person, got %d", len(persons))
	}
}

func TestStore_DataDir(t *testing.T) {
	// Test with invalid data dir
	_, err := New("")
	if err == nil {
		t.Error("expected error for empty data dir")
	}
}
