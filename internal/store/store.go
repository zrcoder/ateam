package store

import (
	"sync"

	"ateam/internal/models"
)

type Store struct {
	mu        sync.RWMutex
	People    map[string]*models.Person
	Agents    map[string]*models.Agent
	Computers map[string]*models.Computer
	Channels  map[string]*models.Channel
	Messages  map[string][]*models.Message
	Tasks     map[string]*models.Task
}

func New() *Store {
	s := &Store{
		People:    make(map[string]*models.Person),
		Agents:    make(map[string]*models.Agent),
		Computers: make(map[string]*models.Computer),
		Channels:  make(map[string]*models.Channel),
		Messages:  make(map[string][]*models.Message),
		Tasks:     make(map[string]*models.Task),
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

	computer, err := models.NewComputer(person.ID, "My MacBook", "localhost")
	if err != nil {
		panic("failed to create computer: " + err.Error())
	}
	s.Computers[computer.ID] = computer
	person.Computers = append(person.Computers, *computer)

	agent1, err := models.NewAgent(person.ID, "dev-bot", "engineer", "claude", "claude-sonnet-4-7")
	if err != nil {
		panic("failed to create agent: " + err.Error())
	}
	s.Agents[agent1.ID] = agent1
	person.AgentIDs = append(person.AgentIDs, agent1.ID)

	agent2, err := models.NewAgent(person.ID, "pm-bot", "product-manager", "opencode", "gpt-4o")
	if err != nil {
		panic("failed to create agent: " + err.Error())
	}
	s.Agents[agent2.ID] = agent2
	person.AgentIDs = append(person.AgentIDs, agent2.ID)

	// Helper to create messages with panic on error (seed data should always be valid)
	newMsg := func(channelID, authorID, authorType, content string) *models.Message {
		msg, err := models.NewMessage(channelID, authorID, authorType, content)
		if err != nil {
			panic("failed to create message: " + err.Error())
		}
		return msg
	}

	mockMessages := []*models.Message{
		newMsg("channel-general", person.ID, models.AuthorTypePerson, "Welcome to **ateam**! This is the beginning of your channel."),
		newMsg("channel-general", person.ID, models.AuthorTypePerson, "Makes sense. This feels like a production-ready TUI app already."),
	}

	s.Messages["channel-general"] = mockMessages

	task1, err := models.NewTask("Setup project structure", person.ID)
	if err != nil {
		panic("failed to create task: " + err.Error())
	}
	s.Tasks[task1.ID] = task1
	task2, err := models.NewTask("Implement TUI layout", person.ID)
	if err != nil {
		panic("failed to create task: " + err.Error())
	}
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
