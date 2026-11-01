package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8) NOT NULL DEFAULT '',
    title   VARCHAR NOT NULL DEFAULT '',
    comment TEXT    NOT NULL DEFAULT '',
    repeat  VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return fmt.Errorf("ping db: %w", err)
	}

	if install {
		if _, err := conn.Exec(schema); err != nil {
			conn.Close()
			return fmt.Errorf("apply schema: %w", err)
		}
	}

	db = conn
	return nil
}