package db

import (
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaSQL string

type Database struct {
	Conn *sql.DB
}

func Connect(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{Conn: db}, nil
}

func (db *Database) Init() error {
	_, err := db.Conn.Exec(schemaSQL)
	return err
}

func (db *Database) Close() error {
	return db.Conn.Close()
}

// LogAudit records an event to the audit log
func (db *Database) LogAudit(event string, userID int, appName, sourceIP, details string) {
	db.Conn.Exec(
		"INSERT INTO audit_log (event, user_id, app_name, source_ip, details) VALUES (?, ?, ?, ?, ?)",
		event, userID, appName, sourceIP, details,
	)
}
