package db

import (
	"testing"
)

func TestConnect(t *testing.T) {
	tmpDir := t.TempDir()

	db, err := Connect(tmpDir)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer Close()

	if db == nil {
		t.Fatal("expected non-nil database")
	}

	if err := db.Ping(); err != nil {
		t.Errorf("Ping() error = %v", err)
	}
}

func TestConnect_EmptyDataDir(t *testing.T) {
	_, err := Connect("")
	if err == nil {
		t.Error("expected error for empty data dir")
	}
}

func TestConnect_SamePath_ReturnsSameConnection(t *testing.T) {
	tmpDir := t.TempDir()

	db1, err := Connect(tmpDir)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	db2, err := Connect(tmpDir)
	if err != nil {
		t.Fatalf("Connect() second call error = %v", err)
	}

	if db1 != db2 {
		t.Error("expected same connection for same path")
	}

	Close()
}

func TestConnect_CreatesTables(t *testing.T) {
	tmpDir := t.TempDir()

	db, err := Connect(tmpDir)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer Close()

	// Check person table
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM person").Scan(&count)
	if err != nil {
		t.Errorf("person table error: %v", err)
	}

	// Check agents table
	err = db.QueryRow("SELECT COUNT(*) FROM agents").Scan(&count)
	if err != nil {
		t.Errorf("agents table error: %v", err)
	}

	// Check channels table
	err = db.QueryRow("SELECT COUNT(*) FROM channels").Scan(&count)
	if err != nil {
		t.Errorf("channels table error: %v", err)
	}

	// Check tasks table
	err = db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		t.Errorf("tasks table error: %v", err)
	}

	// Check messages table
	err = db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
	if err != nil {
		t.Errorf("messages table error: %v", err)
	}
}

func TestConnect_IndexCreated(t *testing.T) {
	tmpDir := t.TempDir()

	db, err := Connect(tmpDir)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer Close()

	// Check indexes exist
	var indexCount int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name LIKE 'idx_%'").Scan(&indexCount)
	if err != nil {
		t.Fatalf("query error: %v", err)
	}

	if indexCount == 0 {
		t.Error("expected indexes to be created")
	}
}

func TestConnect_WALMode(t *testing.T) {
	tmpDir := t.TempDir()

	db, err := Connect(tmpDir)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer Close()

	var mode string
	err = db.QueryRow("PRAGMA journal_mode").Scan(&mode)
	if err != nil {
		t.Fatalf("PRAGMA query error: %v", err)
	}

	if mode != "wal" {
		t.Errorf("expected WAL mode, got %s", mode)
	}
}

func TestClose(t *testing.T) {
	tmpDir := t.TempDir()

	db, err := Connect(tmpDir)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	Close()

	// After Close, pool is cleared. Connecting again should work.
	db2, err := Connect(tmpDir)
	if err != nil {
		t.Fatalf("Connect() after Close() error = %v", err)
	}
	if db == db2 {
		t.Error("expected new connection after Close()")
	}
	Close()
}

func TestConnect_ForeignKeys(t *testing.T) {
	tmpDir := t.TempDir()

	db, err := Connect(tmpDir)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer Close()

	var fk int
	err = db.QueryRow("PRAGMA foreign_keys").Scan(&fk)
	if err != nil {
		t.Fatalf("PRAGMA query error: %v", err)
	}

	if fk != 1 {
		t.Errorf("expected foreign_keys=ON, got %d", fk)
	}
}
