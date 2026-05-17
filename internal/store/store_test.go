package store

import (
	"sync"
	"testing"

	"github.com/zrcoder/ateam/internal/models"
)

func TestStore_AddMessage(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	msg, err := models.NewMessage("channel-1", "author-1", models.AuthorTypePerson, "Hello")
	if err != nil {
		t.Fatalf("failed to create message: %v", err)
	}

	s.AddMessage(msg)

	if len(s.Messages["channel-1"]) != 1 {
		t.Errorf("expected 1 message, got %d", len(s.Messages["channel-1"]))
	}
	if s.Messages["channel-1"][0].Content != "Hello" {
		t.Errorf("expected 'Hello', got %s", s.Messages["channel-1"][0].Content)
	}
}

func TestStore_GetMessages(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	msg1, _ := models.NewMessage("channel-1", "author-1", models.AuthorTypePerson, "Hello")
	msg2, _ := models.NewMessage("channel-1", "author-2", models.AuthorTypeAgent, "Hi there")
	s.AddMessage(msg1)
	s.AddMessage(msg2)

	messages := s.GetMessages("channel-1")
	if len(messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(messages))
	}
}

func TestStore_AddTask(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	task, err := models.NewTask("Test Task", "creator-1")
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	s.AddTask(task)

	if len(s.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(s.Tasks))
	}
	if s.Tasks[task.ID].Title != "Test Task" {
		t.Errorf("expected 'Test Task', got %s", s.Tasks[task.ID].Title)
	}
}

func TestStore_GetTasks(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	task1, _ := models.NewTask("Task 1", "creator-1")
	task2, _ := models.NewTask("Task 2", "creator-1")
	s.AddTask(task1)
	s.AddTask(task2)

	tasks := s.GetTasks()
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestStore_AddAgent(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	person := &models.Person{
		ID:       "person-1",
		Name:     "Test Person",
		Email:    "test@example.com",
		Status:   models.StatusOnline,
		AgentIDs: []string{},
	}
	s.People[person.ID] = person

	agent, err := models.NewAgent("person-1", "dev-bot", "engineer", "claude", "claude-sonnet")
	if err != nil {
		t.Fatalf("failed to create agent: %v", err)
	}

	s.AddAgent(agent)

	if len(s.Agents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(s.Agents))
	}
	if len(person.AgentIDs) != 1 {
		t.Errorf("expected person to have 1 agent, got %d", len(person.AgentIDs))
	}
}

func TestStore_GetAgents(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	agent1, _ := models.NewAgent("person-1", "dev-bot", "engineer", "claude", "model")
	agent2, _ := models.NewAgent("person-1", "pm-bot", "product-manager", "opencode", "model")
	s.AddAgent(agent1)
	s.AddAgent(agent2)

	agents := s.GetAgents()
	if len(agents) != 2 {
		t.Errorf("expected 2 agents, got %d", len(agents))
	}
}

func TestStore_GetAgent(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	agent, _ := models.NewAgent("person-1", "dev-bot", "engineer", "claude", "model")
	s.AddAgent(agent)

	result := s.GetAgent(agent.ID)
	if result == nil {
		t.Error("expected to find agent")
	}
	if result.Name != "dev-bot" {
		t.Errorf("expected 'dev-bot', got %s", result.Name)
	}

	notFound := s.GetAgent("non-existent")
	if notFound != nil {
		t.Error("expected nil for non-existent agent")
	}
}

func TestStore_GetCurrentPerson(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	person := &models.Person{
		ID:    "user-1",
		Name:  "Test User",
		Email: "test@example.com",
	}
	s.People["user-1"] = person

	result := s.GetCurrentPerson()
	if result == nil {
		t.Error("expected to find person")
	}
	if result.Name != "Test User" {
		t.Errorf("expected 'Test User', got %s", result.Name)
	}
}

func TestStore_GetChannels(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	ch1, _ := models.NewChannel("general", "General discussions")
	ch2, _ := models.NewChannel("random", "Random talk")
	s.Channels[ch1.ID] = ch1
	s.Channels[ch2.ID] = ch2

	channels := s.GetChannels()
	if len(channels) != 2 {
		t.Errorf("expected 2 channels, got %d", len(channels))
	}
}

func TestStore_GetPeople(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	p1 := &models.Person{ID: "person-1", Name: "Person 1", Email: "p1@example.com"}
	p2 := &models.Person{ID: "person-2", Name: "Person 2", Email: "p2@example.com"}
	s.People[p1.ID] = p1
	s.People[p2.ID] = p2

	people := s.GetPeople()
	if len(people) != 2 {
		t.Errorf("expected 2 people, got %d", len(people))
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			msg, _ := models.NewMessage("channel-1", "author-1", models.AuthorTypePerson, "Test message")
			s.AddMessage(msg)
		}(i)
	}

	wg.Wait()

	if len(s.Messages["channel-1"]) != 100 {
		t.Errorf("expected 100 messages, got %d", len(s.Messages["channel-1"]))
	}
}

func TestStore_GetMessages_EmptyChannel(t *testing.T) {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
	}

	messages := s.GetMessages("non-existent")
	if len(messages) != 0 {
		t.Errorf("expected 0 messages, got %d", len(messages))
	}
}
