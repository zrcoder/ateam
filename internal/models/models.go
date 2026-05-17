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

func NewPerson(name, email string) (*Person, error) {
	if err := ValidatePerson(name, email); err != nil {
		return nil, err
	}
	return &Person{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Status:    StatusOnline,
		Computers: []Computer{},
		AgentIDs:  []string{},
	}, nil
}

func NewAgent(personID, name, role, aiType, model string) (*Agent, error) {
	if err := ValidateAgent(name); err != nil {
		return nil, err
	}
	if personID == "" {
		return nil, ValidationError{Field: "personID", Message: "person ID cannot be empty"}
	}
	return &Agent{
		ID:       uuid.New().String(),
		PersonID: personID,
		Name:     name,
		Role:     role,
		AIType:   aiType,
		Model:    model,
		Status:   StatusOnline,
	}, nil
}

func NewComputer(personID, name, hostname string) (*Computer, error) {
	if personID == "" {
		return nil, ValidationError{Field: "personID", Message: "person ID cannot be empty"}
	}
	if name == "" {
		return nil, ErrEmptyName
	}
	return &Computer{
		ID:       uuid.New().String(),
		PersonID: personID,
		Name:     name,
		Hostname: hostname,
		Status:   StatusOnline,
	}, nil
}

func NewChannel(name, purpose string) (*Channel, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	return &Channel{
		ID:      uuid.New().String(),
		Name:    name,
		Purpose: purpose,
	}, nil
}

func NewMessage(channelID, authorID, authorType, content string) (*Message, error) {
	if channelID == "" {
		return nil, ValidationError{Field: "channelID", Message: "channel ID cannot be empty"}
	}
	if authorID == "" {
		return nil, ValidationError{Field: "authorID", Message: "author ID cannot be empty"}
	}
	if err := ValidateMessage(content); err != nil {
		return nil, err
	}
	if authorType != AuthorTypePerson && authorType != AuthorTypeAgent && authorType != AuthorTypeSystem {
		return nil, ValidationError{Field: "authorType", Message: "author type must be person, agent, or system"}
	}
	return &Message{
		ID:          uuid.New().String(),
		ChannelID:   channelID,
		AuthorID:    authorID,
		AuthorType:  authorType,
		Content:     content,
		Timestamp:   time.Now(),
	}, nil
}

func NewTask(title, createdBy string) (*Task, error) {
	if err := ValidateTask(title); err != nil {
		return nil, err
	}
	now := time.Now()
	return &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: "",
		Status:      TaskStatusTodo,
		Priority:    TaskPriorityMedium,
		CreatedBy:   createdBy,
		AssigneeID:  "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
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
	AuthorTypeAgent = "agent"
	AuthorTypeSystem = "system"
)

// Validation errors
var (
	ErrEmptyName    = ValidationError{Field: "name", Message: "name cannot be empty"}
	ErrEmptyEmail   = ValidationError{Field: "email", Message: "email cannot be empty"}
	ErrInvalidEmail = ValidationError{Field: "email", Message: "email format is invalid"}
	ErrEmptyTitle   = ValidationError{Field: "title", Message: "title cannot be empty"}
	ErrEmptyContent = ValidationError{Field: "content", Message: "content cannot be empty"}
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// ValidatePerson validates person data
func ValidatePerson(name, email string) error {
	if name == "" {
		return ErrEmptyName
	}
	if email == "" {
		return ErrEmptyEmail
	}
	if !isValidEmail(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidateAgent validates agent data
func ValidateAgent(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	return nil
}

// ValidateTask validates task data
func ValidateTask(title string) error {
	if title == "" {
		return ErrEmptyTitle
	}
	if len(title) > 200 {
		return ValidationError{Field: "title", Message: "title cannot exceed 200 characters"}
	}
	return nil
}

// ValidateMessage validates message data
func ValidateMessage(content string) error {
	if content == "" {
		return ErrEmptyContent
	}
	if len(content) > 10000 {
		return ValidationError{Field: "content", Message: "content cannot exceed 10000 characters"}
	}
	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	atIndex := -1
	for i, c := range email {
		if c == '@' {
			if atIndex != -1 {
				return false
			}
			atIndex = i
		}
	}
	if atIndex < 1 || atIndex == len(email)-1 {
		return false
	}
	dotAfterAt := false
	for i := atIndex + 1; i < len(email); i++ {
		if email[i] == '.' {
			if i == atIndex+1 || i == len(email)-1 {
				return false
			}
			dotAfterAt = true
		}
	}
	return dotAfterAt
}

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