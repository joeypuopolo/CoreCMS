package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Connect inits and returns a SQLite db connection
func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "cms.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create tables if not exist
	createTables(db)

	return db
}

// Implement a history table
// - id, table, column, CRUD function, datetime, content

func createTables(db *sql.DB) {
	pageTable := `
	CREATE TABLE IF NOT EXISTS pages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		url TEXT NOT NULL,
		hidden INTEGER DEFAULT 1,
		active INTEGER DEFAULT 0,
		link STRING,
		link_new_tab INT DEFAULT 0,
		parent_page INT DEFAULT -1,
		settings TEXT,
		template_id INTEGER DEFAULT -1,
	    FOREIGN KEY (template_id) REFERENCES templates (id)
	);`

	// CREATE TABLE IF NOT EXISTS templates ( id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL, active INTEGER DEFAULT 1, parent_template_id INTEGER, all_codeblocks TEXT, FOREIGN KEY (parent_template_id) REFERENCES templates (id));

	templateTable := `
	CREATE TABLE IF NOT EXISTS templates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		active INTEGER DEFAULT 1,
		parent_template_id INTEGER,
		FOREIGN KEY (parent_template_id) REFERENCES templates (id)
	);`

	codeBlockTable := `
	CREATE TABLE IF NOT EXISTS code_blocks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL UNIQUE,
		active INTEGER DEFAULT 1,
		description TEXT,
		content TEXT NOT NULL
	);`

	codeblocksOrdering := `
	CREATE TABLE IF NOT EXISTS codeblocks_ordering (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		page_id INTEGER DEFAULT -1,
		template_id INTEGER DEFAULT -1,
		codeblock_id INTEGER,
		ordering INTEGER,
		active INTEGER,
		FOREIGN KEY (codeblock_id) REFERENCES codeblocks(id)
	);
	`
	_, err := db.Exec(pageTable + codeBlockTable + templateTable + codeblocksOrdering)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}
