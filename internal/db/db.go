package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	pool   = make(map[string]*sql.DB)
	poolMu sync.Mutex
)

// Connect opens a SQLite database connection
func Connect(dataDir string) (*sql.DB, error) {
	if dataDir == "" {
		return nil, fmt.Errorf("data dir is required")
	}

	dbPath := dataDir + "/ateam.db"

	poolMu.Lock()
	defer poolMu.Unlock()

	if db, ok := pool[dbPath]; ok {
		return db, nil
	}

	dsn := "file:" + dbPath + "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)&_txlock=immediate"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	// Only allow one connection to avoid WAL issues
	db.SetMaxOpenConns(1)

	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	pool[dbPath] = db
	return db, nil
}

// Close closes all pooled database connections
func Close() {
	poolMu.Lock()
	defer poolMu.Unlock()
	for _, db := range pool {
		db.Close()
	}
	clear(pool)
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS person (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'online',
			agent_ids TEXT NOT NULL DEFAULT '[]',
			computers TEXT NOT NULL DEFAULT '[]',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS agents (
			id TEXT PRIMARY KEY,
			person_id TEXT NOT NULL,
			name TEXT NOT NULL,
			role TEXT NOT NULL,
			ai_type TEXT NOT NULL,
			model TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'online',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (person_id) REFERENCES person(id)
		)`,
		`CREATE TABLE IF NOT EXISTS channels (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			purpose TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			channel_id TEXT NOT NULL,
			author_id TEXT NOT NULL,
			author_type TEXT NOT NULL,
			content TEXT NOT NULL,
			timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (channel_id) REFERENCES channels(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_channel_id ON messages(channel_id)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'todo',
			priority TEXT NOT NULL DEFAULT 'medium',
			created_by TEXT NOT NULL,
			assignee_id TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (created_by) REFERENCES person(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return err
		}
	}
	return nil
}
