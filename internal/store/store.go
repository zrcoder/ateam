package store

import (
	"fmt"
	"sync"

	"ateam/internal/models"
)

type Store struct {
	mu       sync.RWMutex
	People   map[string]*models.Person
	Agents   map[string]*models.Agent
	Computers map[string]*models.Computer
	Channels map[string]*models.Channel
	Messages map[string][]*models.Message
	Tasks    map[string]*models.Task
}

func New() *Store {
	s := &Store{
		People:   make(map[string]*models.Person),
		Agents:   make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels: make(map[string]*models.Channel),
		Messages: make(map[string][]*models.Message),
		Tasks:    make(map[string]*models.Task),
	}
	s.seed()
	return s
}

func (s *Store) seed() {
	person := &models.Person{
		ID:        "user-1",
		Name:      "You",
		Email:     "you@example.com",
		Status:    models.StatusOnline,
		Computers: []models.Computer{},
		AgentIDs:  []string{},
	}
	s.People[person.ID] = person

	generalChannel := &models.Channel{
		ID:      "channel-general",
		Name:    "general",
		Purpose: "General discussions",
	}
	s.Channels[generalChannel.ID] = generalChannel

	computer := models.NewComputer(person.ID, "My MacBook", "localhost")
	s.Computers[computer.ID] = computer
	person.Computers = append(person.Computers, *computer)

	agent1 := models.NewAgent(person.ID, "dev-bot", "engineer", "claude", "claude-sonnet-4-7")
	s.Agents[agent1.ID] = agent1
	person.AgentIDs = append(person.AgentIDs, agent1.ID)

	agent2 := models.NewAgent(person.ID, "pm-bot", "product-manager", "opencode", "gpt-4o")
	s.Agents[agent2.ID] = agent2
	person.AgentIDs = append(person.AgentIDs, agent2.ID)

	mockMessages := []*models.Message{
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "Welcome to **ateam**! This is the beginning of your channel."),
		models.NewMessage("channel-general", agent1.ID, models.AuthorTypeAgent, "Hello! I'm your AI coding assistant. I can help you write code, review PRs, and more."),
		models.NewMessage("channel-general", agent2.ID, models.AuthorTypeAgent, "I'm the product manager bot. I can help you track tasks and prioritize work."),
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "Great! Let's get started with the project setup."),
		models.NewMessage("channel-general", agent1.ID, models.AuthorTypeAgent, "I've created the initial project structure with Go modules and basic dependencies."),
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "Can you explain the architecture decisions made so far?"),
		models.NewMessage("channel-general", agent1.ID, models.AuthorTypeAgent, "Sure! We decided on a clean architecture with:\n- models/ for data structures\n- store/ for in-memory storage\n- ui/ for the TUI components\n- Each package is independent and testable"),
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "That's a solid approach. How are we handling state management?"),
		models.NewMessage("channel-general", agent2.ID, models.AuthorTypeAgent, "We're using bubbletea's Model pattern for state management. Each component has its own state and updates through messages."),
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "I see. And what about the message rendering?"),
		models.NewMessage("channel-general", agent1.ID, models.AuthorTypeAgent, "Messages are rendered with lipgloss for styling. We use viewport for scrolling when content exceeds the visible area."),
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "Good call on viewport. How does the dialog system work?"),
		models.NewMessage("channel-general", agent1.ID, models.AuthorTypeAgent, "Dialogs are composable layers using lipgloss.NewLayer and lipgloss.NewCompositor. This allows us to overlay dialogs on top of the main content."),
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "Perfect! This TUI framework is really powerful."),
		models.NewMessage("channel-general", agent2.ID, models.AuthorTypeAgent, "Don't forget to test the scroll functionality! Use the viewport's built-in keybindings or Ctrl+U/D for half-page scrolling."),
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "Got it. What else should I be aware of?"),
		models.NewMessage("channel-general", agent1.ID, models.AuthorTypeAgent, "The textarea supports multi-line input with Shift+Enter or Ctrl+J for newlines, and Enter to send."),
		models.NewMessage("channel-general", agent2.ID, models.AuthorTypeAgent, "Type / to trigger the help dialog at any time. It shows all available commands and keybindings."),
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "This is great! The UX is really smooth."),
		models.NewMessage("channel-general", agent1.ID, models.AuthorTypeAgent, "Thanks! We aim for a vim-like experience with intuitive keyboard shortcuts."),
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "I'll make sure to add more content here to test scrolling."),
		models.NewMessage("channel-general", agent2.ID, models.AuthorTypeAgent, "When you add more messages, you'll see the viewport scrolling in action. The most recent messages appear at the bottom."),
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "Scroll position is preserved until new messages are added."),
		models.NewMessage("channel-general", agent1.ID, models.AuthorTypeAgent, "This allows users to scroll back through history while new messages can be posted without disrupting their view."),
		models.NewMessage("channel-general", agent2.ID, models.AuthorTypeAgent, "It's similar to how most chat applications handle scroll position after loading more messages."),
		models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, "Makes sense. This feels like a production-ready TUI app already."),
	}

	for i := 0; i < 50; i++ {
		msg := models.NewMessage("channel-general", person.ID, models.AuthorTypePerson, fmt.Sprintf("This is test message number %02d for scrolling validation.", i))
		mockMessages = append(mockMessages, msg)
	}

	s.Messages["channel-general"] = mockMessages

	task1 := models.NewTask("Setup project structure", person.ID)
	s.Tasks[task1.ID] = task1
	task2 := models.NewTask("Implement TUI layout", person.ID)
	task2.Status = models.TaskStatusInProgress
	s.Tasks[task2.ID] = task2
}

func (s *Store) AddMessage(msg *models.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages[msg.ChannelID] = append(s.Messages[msg.ChannelID], msg)
}

func (s *Store) GetMessages(channelID string) []*models.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Messages[channelID]
}

func (s *Store) AddTask(task *models.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Tasks[task.ID] = task
}

func (s *Store) GetTasks() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]*models.Task, 0, len(s.Tasks))
	for _, t := range s.Tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *Store) AddAgent(agent *models.Agent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Agents[agent.ID] = agent
	if person, ok := s.People[agent.PersonID]; ok {
		person.AgentIDs = append(person.AgentIDs, agent.ID)
	}
}

func (s *Store) GetAgents() []*models.Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	agents := make([]*models.Agent, 0, len(s.Agents))
	for _, a := range s.Agents {
		agents = append(agents, a)
	}
	return agents
}

func (s *Store) GetAgent(id string) *models.Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Agents[id]
}

func (s *Store) GetCurrentPerson() *models.Person {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.People["user-1"]
}

func (s *Store) GetChannels() []*models.Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	channels := make([]*models.Channel, 0, len(s.Channels))
	for _, c := range s.Channels {
		channels = append(channels, c)
	}
	return channels
}

func (s *Store) GetPeople() []*models.Person {
	s.mu.RLock()
	defer s.mu.RUnlock()
	people := make([]*models.Person, 0, len(s.People))
	for _, p := range s.People {
		people = append(people, p)
	}
	return people
}