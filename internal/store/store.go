package store

import (
	"database/sql"

	"github.com/zrcoder/ateam/internal/db"
	"github.com/zrcoder/ateam/internal/models"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func New(dataDir string) (*Store, error) {
	database, err := db.Connect(dataDir)
	if err != nil {
		return nil, err
	}
	s := &Store{db: database}
	if err := s.seed(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) Seed() error {
	return s.seed()
}

func (s *Store) seed() error {
	// Check if current person exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM person").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Already seeded
	}

	personID := "user-1"
	_, err = s.db.Exec(
		`INSERT INTO person (id, name, email, status) VALUES (?, ?, ?, ?)`,
		personID, "You", "you@example.com", models.StatusOnline,
	)
	if err != nil {
		return err
	}

	// Create default channel
	channelID := "channel-general"
	_, err = s.db.Exec(
		`INSERT INTO channels (id, name, purpose) VALUES (?, ?, ?)`,
		channelID, "general", "General discussions",
	)
	if err != nil {
		return err
	}

	// Create two agents
	devBotID := uuid.New().String()
	_, err = s.db.Exec(
		`INSERT INTO agents (id, person_id, name, role, ai_type, model) VALUES (?, ?, ?, ?, ?, ?)`,
		devBotID, personID, "dev-bot", "engineer", "claude", "claude-sonnet-4-7",
	)
	if err != nil {
		return err
	}

	pmBotID := uuid.New().String()
	_, err = s.db.Exec(
		`INSERT INTO agents (id, person_id, name, role, ai_type, model) VALUES (?, ?, ?, ?, ?, ?)`,
		pmBotID, personID, "pm-bot", "product-manager", "opencode", "gpt-4o",
	)
	if err != nil {
		return err
	}

	// Create two tasks
	_, err = s.db.Exec(
		`INSERT INTO tasks (id, title, created_by) VALUES (?, ?, ?)`,
		uuid.New().String(), "Setup project structure", personID,
	)
	if err != nil {
		return err
	}

	task2ID := uuid.New().String()
	_, err = s.db.Exec(
		`INSERT INTO tasks (id, title, status, created_by) VALUES (?, ?, ?, ?)`,
		task2ID, "Implement TUI layout", models.TaskStatusInProgress, personID,
	)
	return err
}

func (s *Store) AddMessage(msg *models.Message) error {
	_, err := s.db.Exec(
		`INSERT INTO messages (id, channel_id, author_id, author_type, content, timestamp) VALUES (?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.ChannelID, msg.AuthorID, msg.AuthorType, msg.Content, msg.Timestamp,
	)
	return err
}

func (s *Store) GetMessages(channelID string) ([]*models.Message, error) {
	rows, err := s.db.Query(
		`SELECT id, channel_id, author_id, author_type, content, timestamp FROM messages WHERE channel_id = ? ORDER BY timestamp ASC`,
		channelID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.ChannelID, &msg.AuthorID, &msg.AuthorType, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, rows.Err()
}

func (s *Store) AddTask(task *models.Task) error {
	_, err := s.db.Exec(
		`INSERT INTO tasks (id, title, description, status, priority, created_by, assignee_id) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		task.ID, task.Title, task.Description, task.Status, task.Priority, task.CreatedBy, task.AssigneeID,
	)
	return err
}

func (s *Store) GetTasks() ([]*models.Task, error) {
	rows, err := s.db.Query(
		`SELECT id, title, description, status, priority, created_by, COALESCE(assignee_id, '') FROM tasks`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedBy, &task.AssigneeID); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, rows.Err()
}

func (s *Store) GetAgents() ([]*models.Agent, error) {
	rows, err := s.db.Query(
		`SELECT id, person_id, name, role, ai_type, model, status FROM agents`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*models.Agent
	for rows.Next() {
		var agent models.Agent
		if err := rows.Scan(&agent.ID, &agent.PersonID, &agent.Name, &agent.Role, &agent.AIType, &agent.Model, &agent.Status); err != nil {
			return nil, err
		}
		agents = append(agents, &agent)
	}
	return agents, rows.Err()
}

func (s *Store) GetAgent(id string) (*models.Agent, error) {
	var agent models.Agent
	err := s.db.QueryRow(
		`SELECT id, person_id, name, role, ai_type, model, status FROM agents WHERE id = ?`,
		id,
	).Scan(&agent.ID, &agent.PersonID, &agent.Name, &agent.Role, &agent.AIType, &agent.Model, &agent.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (s *Store) GetCurrentPerson() *models.Person {
	var person models.Person
	err := s.db.QueryRow(
		`SELECT id, name, email, status, agent_ids, computers FROM person LIMIT 1`,
	).Scan(&person.ID, &person.Name, &person.Email, &person.Status, &person.AgentIDs, &person.Computers)
	if err != nil {
		return nil
	}
	return &person
}

func (s *Store) GetChannels() ([]*models.Channel, error) {
	rows, err := s.db.Query(`SELECT id, name, purpose FROM channels`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		var ch models.Channel
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.Purpose); err != nil {
			return nil, err
		}
		channels = append(channels, &ch)
	}
	return channels, rows.Err()
}

func (s *Store) GetPersons() ([]*models.Person, error) {
	rows, err := s.db.Query(`SELECT id, name, email, status, agent_ids, computers FROM person`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []*models.Person
	for rows.Next() {
		var p models.Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.Status, &p.AgentIDs, &p.Computers); err != nil {
			return nil, err
		}
		persons = append(persons, &p)
	}
	return persons, rows.Err()
}

// Close closes the database connection and checkpoints WAL
func (s *Store) Close() {
	if s.db != nil {
		s.db.Exec("PRAGMA wal_checkpoint(FULL)")
		s.db.Close()
	}
}
