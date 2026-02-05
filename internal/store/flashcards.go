package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mikebway/goTerve/internal/models"
)

// GetUserFlashcard retrieves a flashcard for a user and word
func (s *Store) GetUserFlashcard(userID string, wordID int) (*models.UserFlashcard, error) {
	query := `SELECT id, user_id, word_id, category, times_reviewed, times_correct,
		last_reviewed_at, next_review_at, created_at, updated_at
		FROM user_flashcards WHERE user_id = ? AND word_id = ?`
	row := s.db.QueryRow(query, userID, wordID)

	var fc models.UserFlashcard
	var lastReviewed, nextReview sql.NullTime

	err := row.Scan(&fc.ID, &fc.UserID, &fc.WordID, &fc.Category, &fc.TimesReviewed,
		&fc.TimesCorrect, &lastReviewed, &nextReview, &fc.CreatedAt, &fc.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get flashcard: %w", err)
	}

	if lastReviewed.Valid {
		fc.LastReviewedAt = &lastReviewed.Time
	}
	if nextReview.Valid {
		fc.NextReviewAt = &nextReview.Time
	}
	return &fc, nil
}

// CreateOrUpdateFlashcard creates or updates a user flashcard
func (s *Store) CreateOrUpdateFlashcard(fc *models.UserFlashcard) error {
	existing, err := s.GetUserFlashcard(fc.UserID, fc.WordID)
	if err != nil {
		return err
	}

	now := time.Now()
	if existing == nil {
		// Create new
		query := `INSERT INTO user_flashcards (user_id, word_id, category, times_reviewed, times_correct,
			last_reviewed_at, next_review_at, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := s.db.Exec(query, fc.UserID, fc.WordID, fc.Category, fc.TimesReviewed,
			fc.TimesCorrect, fc.LastReviewedAt, fc.NextReviewAt, now, now)
		if err != nil {
			return fmt.Errorf("failed to create flashcard: %w", err)
		}
		id, _ := result.LastInsertId()
		fc.ID = int(id)
		fc.CreatedAt = now
	} else {
		// Update existing
		query := `UPDATE user_flashcards SET category = ?, times_reviewed = ?, times_correct = ?,
			last_reviewed_at = ?, next_review_at = ?, updated_at = ?
			WHERE id = ?`
		_, err := s.db.Exec(query, fc.Category, fc.TimesReviewed, fc.TimesCorrect,
			fc.LastReviewedAt, fc.NextReviewAt, now, existing.ID)
		if err != nil {
			return fmt.Errorf("failed to update flashcard: %w", err)
		}
		fc.ID = existing.ID
		fc.CreatedAt = existing.CreatedAt
	}
	fc.UpdatedAt = now
	return nil
}

// GetNextFlashcard gets the next flashcard for review
func (s *Store) GetNextFlashcard(userID string, cefrLevel string) (*models.Word, *models.UserFlashcard, error) {
	// First try to get a word that's due for review
	query := `
		SELECT w.id, w.finnish, w.english, w.phonetic, w.cefr_level, w.frequency, w.category,
			w.created_at, w.updated_at,
			uf.id, uf.user_id, uf.word_id, uf.category, uf.times_reviewed, uf.times_correct,
			uf.last_reviewed_at, uf.next_review_at, uf.created_at, uf.updated_at
		FROM user_flashcards uf
		JOIN words w ON uf.word_id = w.id
		WHERE uf.user_id = ? AND uf.next_review_at <= ?
		ORDER BY uf.next_review_at ASC
		LIMIT 1
	`
	row := s.db.QueryRow(query, userID, time.Now())

	word, fc, err := scanWordAndFlashcard(row)
	if err == nil {
		return word, fc, nil
	}
	if err != sql.ErrNoRows {
		return nil, nil, err
	}

	// No cards due - get a new word the user hasn't seen
	levels := getCEFRLevelsUpTo(cefrLevel)

	query = `
		SELECT id, finnish, english, phonetic, cefr_level, frequency, category, created_at, updated_at
		FROM words
		WHERE cefr_level IN (`
	args := make([]interface{}, len(levels)+1)
	for i, level := range levels {
		if i > 0 {
			query += ", "
		}
		query += "?"
		args[i] = level
	}
	query += `) AND id NOT IN (SELECT word_id FROM user_flashcards WHERE user_id = ?)
		ORDER BY frequency ASC
		LIMIT 1`
	args[len(levels)] = userID

	row = s.db.QueryRow(query, args...)

	var w models.Word
	var phonetic, category *string
	err = row.Scan(&w.ID, &w.Finnish, &w.English, &phonetic, &w.CEFRLevel,
		&w.Frequency, &category, &w.CreatedAt, &w.UpdatedAt)
	if err == sql.ErrNoRows {
		// User has seen all words - get one from review pool
		return s.GetRandomFlashcardForUser(userID)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get new word: %w", err)
	}

	if phonetic != nil {
		w.Phonetic = *phonetic
	}
	if category != nil {
		w.Category = *category
	}

	return &w, nil, nil
}

// GetRandomFlashcardForUser gets a random flashcard from the user's review pool
func (s *Store) GetRandomFlashcardForUser(userID string) (*models.Word, *models.UserFlashcard, error) {
	query := `
		SELECT w.id, w.finnish, w.english, w.phonetic, w.cefr_level, w.frequency, w.category,
			w.created_at, w.updated_at,
			uf.id, uf.user_id, uf.word_id, uf.category, uf.times_reviewed, uf.times_correct,
			uf.last_reviewed_at, uf.next_review_at, uf.created_at, uf.updated_at
		FROM user_flashcards uf
		JOIN words w ON uf.word_id = w.id
		WHERE uf.user_id = ?
		ORDER BY RANDOM()
		LIMIT 1
	`
	row := s.db.QueryRow(query, userID)
	return scanWordAndFlashcard(row)
}

func scanWordAndFlashcard(row *sql.Row) (*models.Word, *models.UserFlashcard, error) {
	var word models.Word
	var fc models.UserFlashcard
	var phonetic, category *string
	var lastReviewed, nextReview sql.NullTime
	var fcCreatedAt, fcUpdatedAt time.Time

	err := row.Scan(
		&word.ID, &word.Finnish, &word.English, &phonetic, &word.CEFRLevel,
		&word.Frequency, &category, &word.CreatedAt, &word.UpdatedAt,
		&fc.ID, &fc.UserID, &fc.WordID, &fc.Category, &fc.TimesReviewed, &fc.TimesCorrect,
		&lastReviewed, &nextReview, &fcCreatedAt, &fcUpdatedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	if phonetic != nil {
		word.Phonetic = *phonetic
	}
	if category != nil {
		word.Category = *category
	}
	if lastReviewed.Valid {
		fc.LastReviewedAt = &lastReviewed.Time
	}
	if nextReview.Valid {
		fc.NextReviewAt = &nextReview.Time
	}
	fc.CreatedAt = fcCreatedAt
	fc.UpdatedAt = fcUpdatedAt

	return &word, &fc, nil
}
