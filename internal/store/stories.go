package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mikebway/goTerve/internal/models"
)

// CreateSavedStory saves a story for a user
func (s *Store) CreateSavedStory(story *models.SavedStory) error {
	query := `INSERT INTO saved_stories (user_id, story, translation, cefr_level, topic, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()
	result, err := s.db.Exec(query, story.UserID, story.Story, story.Translation,
		story.CEFRLevel, story.Topic, now)
	if err != nil {
		return fmt.Errorf("failed to save story: %w", err)
	}
	id, _ := result.LastInsertId()
	story.ID = int(id)
	story.CreatedAt = now
	return nil
}

// GetUserStories retrieves all stories for a user
func (s *Store) GetUserStories(userID string, limit int) ([]models.SavedStory, error) {
	query := `SELECT id, user_id, story, translation, cefr_level, topic, created_at
		FROM saved_stories WHERE user_id = ? ORDER BY created_at DESC`
	args := []interface{}{userID}

	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get stories: %w", err)
	}
	defer rows.Close()

	var stories []models.SavedStory
	for rows.Next() {
		var story models.SavedStory
		var topic sql.NullString

		err := rows.Scan(&story.ID, &story.UserID, &story.Story, &story.Translation,
			&story.CEFRLevel, &topic, &story.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan story: %w", err)
		}

		if topic.Valid {
			story.Topic = topic.String
		}
		stories = append(stories, story)
	}

	return stories, nil
}

// GetStoryByID retrieves a specific story by ID
func (s *Store) GetStoryByID(id int) (*models.SavedStory, error) {
	query := `SELECT id, user_id, story, translation, cefr_level, topic, created_at
		FROM saved_stories WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var story models.SavedStory
	var topic sql.NullString

	err := row.Scan(&story.ID, &story.UserID, &story.Story, &story.Translation,
		&story.CEFRLevel, &topic, &story.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get story: %w", err)
	}

	if topic.Valid {
		story.Topic = topic.String
	}
	return &story, nil
}
