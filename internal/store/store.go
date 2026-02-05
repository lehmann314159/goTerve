package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Store represents the database connection
type Store struct {
	db *sql.DB
}

// New creates a new store with a SQLite database
func New(dbPath string) (*Store, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// DB returns the underlying database connection
func (s *Store) DB() *sql.DB {
	return s.db
}

// Migrate runs database migrations
func (s *Store) Migrate() error {
	migrations := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			avatar TEXT,
			google_id TEXT,
			cefr_level TEXT DEFAULT 'A1',
			has_completed_onboarding INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Words table
		`CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			finnish TEXT UNIQUE NOT NULL,
			english TEXT NOT NULL,
			phonetic TEXT,
			cefr_level TEXT DEFAULT 'A1',
			frequency INTEGER DEFAULT 100,
			category TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Verbs table
		`CREATE TABLE IF NOT EXISTS verbs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			infinitive TEXT UNIQUE NOT NULL,
			type INTEGER NOT NULL,
			stem TEXT NOT NULL,
			translation TEXT NOT NULL,
			examples TEXT DEFAULT '[]',
			frequency INTEGER DEFAULT 100,
			cefr_level TEXT DEFAULT 'A1',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Nouns table
		`CREATE TABLE IF NOT EXISTS nouns (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nominative TEXT UNIQUE NOT NULL,
			translation TEXT NOT NULL,
			examples TEXT DEFAULT '[]',
			declension_type INTEGER NOT NULL,
			stem TEXT NOT NULL,
			cefr_level TEXT DEFAULT 'A1',
			frequency INTEGER DEFAULT 100,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// User flashcards table
		`CREATE TABLE IF NOT EXISTS user_flashcards (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			word_id INTEGER NOT NULL,
			category TEXT DEFAULT 'new',
			times_reviewed INTEGER DEFAULT 0,
			times_correct INTEGER DEFAULT 0,
			last_reviewed_at DATETIME,
			next_review_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (word_id) REFERENCES words(id),
			UNIQUE(user_id, word_id)
		)`,

		// Saved stories table
		`CREATE TABLE IF NOT EXISTS saved_stories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			story TEXT NOT NULL,
			translation TEXT NOT NULL,
			cefr_level TEXT DEFAULT 'A1',
			topic TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		// Test sessions table
		`CREATE TABLE IF NOT EXISTS test_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT,
			test_type TEXT NOT NULL,
			cefr_level TEXT DEFAULT 'A1',
			total_questions INTEGER DEFAULT 0,
			correct_answers INTEGER DEFAULT 0,
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			completed_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		// Sessions table for auth
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_words_cefr ON words(cefr_level)`,
		`CREATE INDEX IF NOT EXISTS idx_words_category ON words(category)`,
		`CREATE INDEX IF NOT EXISTS idx_verbs_infinitive ON verbs(infinitive)`,
		`CREATE INDEX IF NOT EXISTS idx_nouns_nominative ON nouns(nominative)`,
		`CREATE INDEX IF NOT EXISTS idx_user_flashcards_user ON user_flashcards(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id)`,
	}

	for _, migration := range migrations {
		if _, err := s.db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

// Seed populates the database with initial data
func (s *Store) Seed() error {
	// Seed words
	if err := s.seedWords(); err != nil {
		return err
	}

	// Seed verbs
	if err := s.seedVerbs(); err != nil {
		return err
	}

	// Seed nouns
	if err := s.seedNouns(); err != nil {
		return err
	}

	return nil
}
