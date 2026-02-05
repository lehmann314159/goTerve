package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mikebway/goTerve/internal/models"
)

// CreateTestSession creates a new test session
func (s *Store) CreateTestSession(session *models.TestSession) error {
	query := `INSERT INTO test_sessions (user_id, test_type, cefr_level, total_questions, correct_answers, started_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()

	var userID interface{}
	if session.UserID != "" {
		userID = session.UserID
	}

	result, err := s.db.Exec(query, userID, session.TestType, session.CEFRLevel,
		session.TotalQuestions, session.CorrectAnswers, now)
	if err != nil {
		return fmt.Errorf("failed to create test session: %w", err)
	}
	id, _ := result.LastInsertId()
	session.ID = int(id)
	session.StartedAt = now
	return nil
}

// UpdateTestSession updates a test session
func (s *Store) UpdateTestSession(session *models.TestSession) error {
	query := `UPDATE test_sessions SET total_questions = ?, correct_answers = ?, completed_at = ?
		WHERE id = ?`
	_, err := s.db.Exec(query, session.TotalQuestions, session.CorrectAnswers,
		session.CompletedAt, session.ID)
	if err != nil {
		return fmt.Errorf("failed to update test session: %w", err)
	}
	return nil
}

// GetTestSession retrieves a test session by ID
func (s *Store) GetTestSession(id int) (*models.TestSession, error) {
	query := `SELECT id, user_id, test_type, cefr_level, total_questions, correct_answers, started_at, completed_at
		FROM test_sessions WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var session models.TestSession
	var userID sql.NullString
	var completedAt sql.NullTime

	err := row.Scan(&session.ID, &userID, &session.TestType, &session.CEFRLevel,
		&session.TotalQuestions, &session.CorrectAnswers, &session.StartedAt, &completedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get test session: %w", err)
	}

	if userID.Valid {
		session.UserID = userID.String
	}
	if completedAt.Valid {
		session.CompletedAt = &completedAt.Time
	}
	return &session, nil
}

// GetUserTestSessions retrieves test sessions for a user
func (s *Store) GetUserTestSessions(userID string, limit int) ([]models.TestSession, error) {
	query := `SELECT id, user_id, test_type, cefr_level, total_questions, correct_answers, started_at, completed_at
		FROM test_sessions WHERE user_id = ? ORDER BY started_at DESC`
	args := []interface{}{userID}

	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get test sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.TestSession
	for rows.Next() {
		var session models.TestSession
		var uid sql.NullString
		var completedAt sql.NullTime

		err := rows.Scan(&session.ID, &uid, &session.TestType, &session.CEFRLevel,
			&session.TotalQuestions, &session.CorrectAnswers, &session.StartedAt, &completedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan test session: %w", err)
		}

		if uid.Valid {
			session.UserID = uid.String
		}
		if completedAt.Valid {
			session.CompletedAt = &completedAt.Time
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
