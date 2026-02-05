package store

import (
	"fmt"

	"github.com/mikebway/goTerve/internal/models"
)

// GetWords retrieves all words, optionally filtered by CEFR level
func (s *Store) GetWords(cefrLevel string, limit int) ([]models.Word, error) {
	query := `SELECT id, finnish, english, phonetic, cefr_level, frequency, category, created_at, updated_at FROM words`
	args := []interface{}{}

	if cefrLevel != "" {
		query += ` WHERE cefr_level = ?`
		args = append(args, cefrLevel)
	}

	query += ` ORDER BY frequency ASC`

	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get words: %w", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		var phonetic, category *string

		err := rows.Scan(&word.ID, &word.Finnish, &word.English, &phonetic,
			&word.CEFRLevel, &word.Frequency, &category, &word.CreatedAt, &word.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan word: %w", err)
		}

		if phonetic != nil {
			word.Phonetic = *phonetic
		}
		if category != nil {
			word.Category = *category
		}
		words = append(words, word)
	}

	return words, nil
}

// GetWordByID retrieves a word by ID
func (s *Store) GetWordByID(id int) (*models.Word, error) {
	query := `SELECT id, finnish, english, phonetic, cefr_level, frequency, category, created_at, updated_at FROM words WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var word models.Word
	var phonetic, category *string

	err := row.Scan(&word.ID, &word.Finnish, &word.English, &phonetic,
		&word.CEFRLevel, &word.Frequency, &category, &word.CreatedAt, &word.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get word: %w", err)
	}

	if phonetic != nil {
		word.Phonetic = *phonetic
	}
	if category != nil {
		word.Category = *category
	}
	return &word, nil
}

// GetRandomWord retrieves a random word at or below the given CEFR level
func (s *Store) GetRandomWord(cefrLevel string) (*models.Word, error) {
	levels := getCEFRLevelsUpTo(cefrLevel)

	query := `SELECT id, finnish, english, phonetic, cefr_level, frequency, category, created_at, updated_at
		FROM words WHERE cefr_level IN (`
	args := make([]interface{}, len(levels))
	for i, level := range levels {
		if i > 0 {
			query += ", "
		}
		query += "?"
		args[i] = level
	}
	query += `) ORDER BY RANDOM() LIMIT 1`

	row := s.db.QueryRow(query, args...)

	var word models.Word
	var phonetic, category *string

	err := row.Scan(&word.ID, &word.Finnish, &word.English, &phonetic,
		&word.CEFRLevel, &word.Frequency, &category, &word.CreatedAt, &word.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get random word: %w", err)
	}

	if phonetic != nil {
		word.Phonetic = *phonetic
	}
	if category != nil {
		word.Category = *category
	}
	return &word, nil
}

// Helper function to get CEFR levels up to and including the given level
func getCEFRLevelsUpTo(level string) []string {
	allLevels := []string{"A1", "A2", "B1", "B2", "C1", "C2"}
	var result []string
	for _, l := range allLevels {
		result = append(result, l)
		if l == level {
			break
		}
	}
	return result
}
