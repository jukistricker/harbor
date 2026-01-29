package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// Đổi tên từ Store thành Client cho đúng ý nghĩa hạ tầng
type Client struct {
	Conn *sql.DB
}

// Đổi tên từ NewStore thành NewClient
func NewClient(dbPath string) (*Client, error) {
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_sync=NORMAL&_cache_size=-10000&_busy_timeout=5000", dbPath)
	
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	pragmas := []string{
		"PRAGMA mmap_size = 30000000000;", 
		"PRAGMA temp_store = MEMORY;",      
		"PRAGMA locking_mode = EXCLUSIVE;", 
		"PRAGMA page_size = 4096;",         
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return nil, fmt.Errorf("failed to set pragma %s: %w", p, err)
		}
	}

	db.SetMaxOpenConns(1) 

	return &Client{Conn: db}, nil
}

// Đổi tên receiver từ s *Store thành c *Client
func (c *Client) InitSchema(sqlFile string) error {
	content, err := os.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("failed to read sql file: %w", err)
	}

	_, err = c.Conn.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}