package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mikebway/goTerve/internal/models"
)

// CreateUser creates a new user
func (s *Store) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, name, avatar, google_id, cefr_level, has_completed_onboarding, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	_, err := s.db.Exec(query, user.ID, user.Email, user.Name, user.Avatar, user.GoogleID,
		user.CEFRLevel, user.HasCompletedOnboarding, now, now)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

// GetUserByID retrieves a user by ID
func (s *Store) GetUserByID(id string) (*models.User, error) {
	query := `SELECT id, email, name, avatar, google_id, cefr_level, has_completed_onboarding, created_at, updated_at FROM users WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var user models.User
	var avatar, googleID sql.NullString
	var hasCompleted int

	err := row.Scan(&user.ID, &user.Email, &user.Name, &avatar, &googleID,
		&user.CEFRLevel, &hasCompleted, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Avatar = avatar.String
	user.GoogleID = googleID.String
	user.HasCompletedOnboarding = hasCompleted == 1
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, name, avatar, google_id, cefr_level, has_completed_onboarding, created_at, updated_at FROM users WHERE email = ?`
	row := s.db.QueryRow(query, email)

	var user models.User
	var avatar, googleID sql.NullString
	var hasCompleted int

	err := row.Scan(&user.ID, &user.Email, &user.Name, &avatar, &googleID,
		&user.CEFRLevel, &hasCompleted, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Avatar = avatar.String
	user.GoogleID = googleID.String
	user.HasCompletedOnboarding = hasCompleted == 1
	return &user, nil
}

// UpdateUser updates a user
func (s *Store) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET name = ?, avatar = ?, cefr_level = ?, has_completed_onboarding = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	_, err := s.db.Exec(query, user.Name, user.Avatar, user.CEFRLevel,
		user.HasCompletedOnboarding, now, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	user.UpdatedAt = now
	return nil
}
