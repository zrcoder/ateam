package models

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	Status     string      `json:"status"`
	Computers  []Computer  `json:"computers"`
	AgentIDs   []string    `json:"agent_ids"`
}

type Agent struct {
	ID       string `json:"id"`
	PersonID string `json:"person_id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	AIType   string `json:"ai_type"`
	Model    string `json:"model"`
	Status   string `json:"status"`
}

type Computer struct {
	ID       string `json:"id"`
	PersonID string `json:"person_id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Status   string `json:"status"`
}

type Channel struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Purpose string `json:"purpose"`
}

type Message struct {
	ID          string    `json:"id"`
	ChannelID   string    `json:"channel_id"`
	AuthorID    string    `json:"author_id"`
	AuthorType  string    `json:"author_type"`
	Content     string    `json:"content"`
	Timestamp   time.Time `json:"timestamp"`
}

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedBy   string    `json:"created_by"`
	AssigneeID  string    `json:"assignee_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewPerson(name, email string) *Person {
	return &Person{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Status:    "online",
		Computers: []Computer{},
		AgentIDs:  []string{},
	}
}

func NewAgent(personID, name, role, aiType, model string) *Agent {
	return &Agent{
		ID:       uuid.New().String(),
		PersonID: personID,
		Name:     name,
		Role:     role,
		AIType:   aiType,
		Model:    model,
		Status:   "online",
	}
}

func NewComputer(personID, name, hostname string) *Computer {
	return &Computer{
		ID:       uuid.New().String(),
		PersonID: personID,
		Name:     name,
		Hostname: hostname,
		Status:   "online",
	}
}

func NewChannel(name, purpose string) *Channel {
	return &Channel{
		ID:      uuid.New().String(),
		Name:    name,
		Purpose: purpose,
	}
}

func NewMessage(channelID, authorID, authorType, content string) *Message {
	return &Message{
		ID:          uuid.New().String(),
		ChannelID:   channelID,
		AuthorID:    authorID,
		AuthorType:  authorType,
		Content:     content,
		Timestamp:  time.Now(),
	}
}

func NewTask(title, createdBy string) *Task {
	now := time.Now()
	return &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: "",
		Status:      "todo",
		Priority:    "medium",
		CreatedBy:   createdBy,
		AssigneeID:  "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

const (
	StatusOnline  = "online"
	StatusAway    = "away"
	StatusOffline = "offline"

	TaskStatusTodo       = "todo"
	TaskStatusInProgress = "in-progress"
	TaskStatusDone       = "done"

	TaskPriorityLow    = "low"
	TaskPriorityMedium = "medium"
	TaskPriorityHigh   = "high"

	AuthorTypePerson = "person"
	AuthorTypeAgent  = "agent"
)

type StoreInterface interface {
	GetCurrentPerson() *Person
	GetAgents() []*Agent
	GetAgent(id string) *Agent
	GetPeople() []*Person
	GetChannels() []*Channel
	GetMessages(channelID string) []*Message
	AddMessage(msg *Message)
	GetTasks() []*Task
	AddTask(task *Task)
}